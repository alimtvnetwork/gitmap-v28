package cmdssh

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// ClusterRunResult holds the execution outcome from a remote cluster node.
type ClusterRunResult struct {
	Host     store.SSHHost
	ExitCode int
	Stdout   string
	Stderr   string
	Duration time.Duration
	Err      error
}

var (
	executeNodeFn  = ExecuteNodeCommand
	clusterPrintMu sync.Mutex
)

func escapeSingleQuotes(s string) string {
	return strings.ReplaceAll(s, "'", "'\\''")
}

// FormatClusterCommand formats a remote command, optionally wrapping with sudo -S.
func FormatClusterCommand(shellCmd string, password string, isSudo bool) string {
	escapedCmd := escapeSingleQuotes(shellCmd)
	if isSudo {
		escapedPass := escapeSingleQuotes(password)
		return fmt.Sprintf("echo '%s' | sudo -S bash -c '%s'", escapedPass, escapedCmd)
	}
	return fmt.Sprintf("bash -c '%s'", escapedCmd)
}

// FormatClusterScriptCommand formats execution of a remote script file.
func FormatClusterScriptCommand(scriptPath string, password string, isSudo bool) string {
	if isSudo {
		escapedPass := escapeSingleQuotes(password)
		return fmt.Sprintf("echo '%s' | sudo -S bash %s", escapedPass, scriptPath)
	}
	return fmt.Sprintf("bash %s", scriptPath)
}

// FormatClusterScriptCleanup formats remote script deletion.
func FormatClusterScriptCleanup(scriptPath string, password string, isSudo bool) string {
	if isSudo {
		escapedPass := escapeSingleQuotes(password)
		return fmt.Sprintf("echo '%s' | sudo -S rm -f %s", escapedPass, scriptPath)
	}
	return fmt.Sprintf("rm -f %s", scriptPath)
}

func resolveHostAlias(host store.SSHHost) string {
	if host.Alias != "" {
		return host.Alias
	}
	if host.IP != "" {
		return host.IP
	}
	return host.ID
}

func resolveNodePassword(host store.SSHHost) (string, error) {
	if host.EncryptedPassword == "" {
		return "", nil
	}
	return DecryptSSHPassword(host.EncryptedPassword)
}

// BuildNodeSSHArgs builds the arguments slice for invoking ssh against a host.
func BuildNodeSSHArgs(host store.SSHHost, remoteCmd string) []string {
	args := []string{"-o", "StrictHostKeyChecking=no"}
	if host.Port > 0 && host.Port != 22 {
		args = append(args, "-p", strconv.Itoa(host.Port))
	}
	user := host.Username
	if user == "" {
		user = "root"
	}
	target := fmt.Sprintf("%s@%s", user, host.IP)
	args = append(args, target, remoteCmd)
	return args
}

func buildDecryptErrorResult(host store.SSHHost, err error, start time.Time) ClusterRunResult {
	alias := resolveHostAlias(host)
	fmt.Fprintf(os.Stderr, "[%s] Decrypt password failed: %v\n", alias, err)
	return ClusterRunResult{
		Host:     host,
		ExitCode: 1,
		Duration: time.Since(start),
		Err:      apperror.WrapSimple(err, "DecryptSSHPassword"),
	}
}

func buildPipeErrorResult(host store.SSHHost, err error, start time.Time) ClusterRunResult {
	return ClusterRunResult{
		Host:     host,
		ExitCode: 1,
		Duration: time.Since(start),
		Err:      apperror.WrapSimple(err, "createCommandPipes"),
	}
}

func buildStartErrorResult(host store.SSHHost, err error, start time.Time) ClusterRunResult {
	alias := resolveHostAlias(host)
	fmt.Fprintf(os.Stderr, "[%s] Process start failed: %v\n", alias, err)
	return ClusterRunResult{
		Host:     host,
		ExitCode: 1,
		Duration: time.Since(start),
		Err:      apperror.WrapSimple(err, "cmd.Start"),
	}
}

func buildClusterRunResult(host store.SSHHost, exitCode int, stdout, stderr string, d time.Duration, err error) ClusterRunResult {
	return ClusterRunResult{
		Host:     host,
		ExitCode: exitCode,
		Stdout:   stdout,
		Stderr:   stderr,
		Duration: d,
		Err:      err,
	}
}

func resolveProcessExitCode(err error) int {
	if err == nil {
		return 0
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return exitErr.ExitCode()
	}
	return 1
}

func createCommandPipes(cmd *exec.Cmd) (io.ReadCloser, io.ReadCloser, error) {
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, apperror.WrapSimple(err, "StdoutPipe")
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return nil, nil, apperror.WrapSimple(err, "StderrPipe")
	}
	return stdoutPipe, stderrPipe, nil
}

func printPrefixedLine(alias string, line string, isStderr bool) {
	clusterPrintMu.Lock()
	defer clusterPrintMu.Unlock()
	if isStderr {
		fmt.Fprintf(os.Stderr, "[%s] %s\n", alias, line)
		return
	}
	fmt.Fprintf(os.Stdout, "[%s] %s\n", alias, line)
}

func scanPipeToOutput(r io.Reader, buf *strings.Builder, alias string, isStderr bool, wg *sync.WaitGroup) {
	defer wg.Done()
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		printPrefixedLine(alias, line, isStderr)
		buf.WriteString(line)
		buf.WriteString("\n")
	}
}

func streamPipesLive(stdoutPipe, stderrPipe io.ReadCloser, alias string) (string, string) {
	var stdoutBuf, stderrBuf strings.Builder
	var wg sync.WaitGroup
	wg.Add(2)
	go scanPipeToOutput(stdoutPipe, &stdoutBuf, alias, false, &wg)
	go scanPipeToOutput(stderrPipe, &stderrBuf, alias, true, &wg)
	wg.Wait()
	return stdoutBuf.String(), stderrBuf.String()
}

func executeAndStreamCmd(cmd *exec.Cmd, host store.SSHHost, alias string, startTime time.Time) ClusterRunResult {
	stdoutPipe, stderrPipe, err := createCommandPipes(cmd)
	if err != nil {
		return buildPipeErrorResult(host, err, startTime)
	}
	if err := cmd.Start(); err != nil {
		return buildStartErrorResult(host, err, startTime)
	}
	stdout, stderr := streamPipesLive(stdoutPipe, stderrPipe, alias)
	waitErr := cmd.Wait()
	exitCode := resolveProcessExitCode(waitErr)
	return buildClusterRunResult(host, exitCode, stdout, stderr, time.Since(startTime), waitErr)
}

func runNodeDirectCmd(ctx context.Context, host store.SSHHost, alias, formattedCmd, password string, startTime time.Time) ClusterRunResult {
	sshArgs := BuildNodeSSHArgs(host, formattedCmd)
	cmd := SSHExecutor(ctx, "ssh", sshArgs...)
	cleanup := attachAskPass(cmd, password)
	defer cleanup()
	return executeAndStreamCmd(cmd, host, alias, startTime)
}

func runNodeSSHProcess(ctx context.Context, host store.SSHHost, alias, shellCmd, password string, isSudo bool, startTime time.Time) ClusterRunResult {
	formattedCmd := FormatClusterCommand(shellCmd, password, isSudo)
	return runNodeDirectCmd(ctx, host, alias, formattedCmd, password, startTime)
}

// ExecuteNodeCommand executes a shell command on an SSH host with live output streaming.
func ExecuteNodeCommand(ctx context.Context, host store.SSHHost, shellCmd string, isSudo bool) ClusterRunResult {
	startTime := time.Now()
	alias := resolveHostAlias(host)
	password, err := resolveNodePassword(host)
	if err != nil {
		return buildDecryptErrorResult(host, err, startTime)
	}
	return runNodeSSHProcess(ctx, host, alias, shellCmd, password, isSudo, startTime)
}

func resolveConcurrencyLimit(concurrency int) int {
	if concurrency > 0 {
		return concurrency
	}
	return 4
}

func runClusterWorker(ctx context.Context, host store.SSHHost, shellCmd string, isSudo bool, idx int, results []ClusterRunResult, sem chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	sem <- struct{}{}
	defer func() { <-sem }()
	results[idx] = executeNodeFn(ctx, host, shellCmd, isSudo)
}

func resolveResultStatus(res ClusterRunResult) string {
	if res.ExitCode == 0 && res.Err == nil {
		return "SUCCESS"
	}
	return "FAILED"
}

func formatSummaryRow(res ClusterRunResult) string {
	status := resolveResultStatus(res)
	alias := resolveHostAlias(res.Host)
	role := res.Host.ClusterRole
	if role == "" {
		role = "worker"
	}
	durStr := res.Duration.Round(time.Millisecond).String()
	return fmt.Sprintf("%-16s %-10s %-16s %-10s %-10d %-10s\n", alias, role, res.Host.IP, status, res.ExitCode, durStr)
}

// RenderClusterSummaryTable formats results into a clean CLI table string.
func RenderClusterSummaryTable(results []ClusterRunResult) string {
	var sb strings.Builder
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("%-16s %-10s %-16s %-10s %-10s %-10s\n", "NODE", "ROLE", "IP", "STATUS", "EXIT CODE", "DURATION"))
	sb.WriteString(strings.Repeat("-", 76))
	sb.WriteString("\n")
	for _, res := range results {
		sb.WriteString(formatSummaryRow(res))
	}
	return sb.String()
}

// PrintClusterSummaryTable prints the summary table to stdout.
func PrintClusterSummaryTable(results []ClusterRunResult) {
	table := RenderClusterSummaryTable(results)
	fmt.Print(table)
}

// DispatchClusterRun executes a command across multiple hosts with bounded concurrency.
func DispatchClusterRun(ctx context.Context, hosts []store.SSHHost, shellCmd string, isSudo bool, concurrency int) []ClusterRunResult {
	limit := resolveConcurrencyLimit(concurrency)
	results := make([]ClusterRunResult, len(hosts))
	var wg sync.WaitGroup
	sem := make(chan struct{}, limit)
	for idx, host := range hosts {
		wg.Add(1)
		go runClusterWorker(ctx, host, shellCmd, isSudo, idx, results, sem, &wg)
	}
	wg.Wait()
	PrintClusterSummaryTable(results)
	return results
}
