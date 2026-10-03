package cmdssh

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/term"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// SSHExecutor is the command factory used for executing SSH commands.
var SSHExecutor = exec.CommandContext

type InteractiveSSHClient struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func executeClientCmd(cmd *exec.Cmd, op string, ctx map[string]any) error {
	if err := cmd.Run(); err != nil {
		return &apperror.AppError{
			Op:       op,
			Code:     "E_INTERNAL_ERROR",
			Type:     apperror.ErrorTypeExecution,
			Severity: apperror.SeverityError,
			Cause:    err,
			Caller:   apperror.CaptureCaller(apperror.DefaultCallerSkip),
			Stack:    apperror.CaptureStackTrace(apperror.DefaultStackTraceSkip),
			Ctx:      ctx,
		}
	}

	return nil
}

func (c *InteractiveSSHClient) Run(ctx context.Context, target string) error {
	cmd := SSHExecutor(ctx, "ssh", target)
	cmd.Stdin = c.Stdin
	cmd.Stdout = c.Stdout
	cmd.Stderr = c.Stderr

	return executeClientCmd(cmd, "InteractiveSSHClient.Run", map[string]any{"target": target})
}

func isCustomSSHPort(port int) bool {
	return port > 0 && port != 22
}

func appendHostKeyCheckingDefault(cmdArgs []string, args []string) []string {
	for _, a := range args {
		if strings.Contains(a, "StrictHostKeyChecking") {
			return cmdArgs
		}
	}
	return append(cmdArgs, "-o", "StrictHostKeyChecking=accept-new")
}

func buildSSHArgs(target SSHTarget, args []string) []string {
	var cmdArgs []string
	if isCustomSSHPort(target.Port) {
		cmdArgs = append(cmdArgs, "-p", strconv.Itoa(target.Port))
	}
	cmdArgs = appendHostKeyCheckingDefault(cmdArgs, args)
	cmdArgs = append(cmdArgs, target.String())
	cmdArgs = append(cmdArgs, args...)

	return cmdArgs
}

func attachAskPass(cmd *exec.Cmd, password string) func() {
	if password == "" {
		return func() {}
	}
	scriptPath, cleanup, err := CreateAskPassScript()
	if err != nil {
		return func() {}
	}
	cmd.Env = BuildAskPassEnv(os.Environ(), scriptPath, password)
	return cleanup
}

func handleHostKeyRecovery(ctx context.Context, target *SSHTarget, errStr string) {
	fmt.Fprintf(os.Stderr, "\n[ssh] Detected changed host key for %s. Auto-pruning stale host key and re-trusting...\n", target.IP)
	if offPath, lineNum, hasLine := ParseOffendingKnownHostsLine(errStr); hasLine {
		_ = RemoveKnownHostsLine(offPath, lineNum)
	}
	_ = PruneHostFromKnownHosts(target.IP, target.Port)
	autoTrustTargetHostFn(ctx, target)
}

func runSSHOnce(ctx context.Context, target SSHTarget, args []string, password string) (error, string) {
	cmdArgs := buildSSHArgs(target, args)
	cmd := SSHExecutor(ctx, "ssh", cmdArgs...)
	cleanup := attachAskPass(cmd, password)
	defer cleanup()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	var errBuf bytes.Buffer
	cmd.Stderr = io.MultiWriter(os.Stderr, &errBuf)
	err := executeClientCmd(cmd, "SpawnSSH", map[string]any{"target": target.String(), "args": args})
	return err, errBuf.String()
}

func isConnectionOrDaemonError(err error, errStr string) bool {
	if err == nil {
		return false
	}

	lowerErr := strings.ToLower(err.Error())
	lowerBuf := strings.ToLower(errStr)

	if strings.Contains(lowerErr, "connection timed out") || strings.Contains(lowerBuf, "connection timed out") {
		return true
	}

	if strings.Contains(lowerErr, "connection refused") || strings.Contains(lowerBuf, "connection refused") {
		return true
	}

	if strings.Contains(lowerErr, "255") || strings.Contains(lowerBuf, "255") {
		return true
	}

	return false
}

func printSSHDiagnosticFooter(target SSHTarget) {
	fmt.Fprintf(os.Stderr, "\n%s[gitmap] SSH Connection Failed to %s%s\n", constants.ColorYellow, target.IP, constants.ColorReset)
	fmt.Fprintf(os.Stderr, "  • OpenSSH Server might not be running or enabled on target machine.\n")
	fmt.Fprintf(os.Stderr, "  • To diagnose: run '%sgitmap ssh troubleshoot %s%s'\n", constants.ColorCyan, target.IP, constants.ColorReset)
	fmt.Fprintf(os.Stderr, "  • On target machine, run '%sgitmap ssh enable%s' to automatically install & start sshd.\n\n", constants.ColorCyan, constants.ColorReset)
}

func checkAndPrintDiagnosticFooter(target SSHTarget, err error, errStr string) {
	if !isConnectionOrDaemonError(err, errStr) {
		return
	}

	printSSHDiagnosticFooter(target)
}

func retrySSHAfterHostKeyRecovery(ctx context.Context, target SSHTarget, args []string, password string, errStr string) error {
	handleHostKeyRecovery(ctx, &target, errStr)
	retryErr, retryStr := runSSHOnce(ctx, target, args, password)
	checkAndPrintDiagnosticFooter(target, retryErr, retryStr)

	return retryErr
}

func SpawnSSHWithPassword(ctx context.Context, target SSHTarget, args []string, password string) error {
	autoTrustTargetHostFn(ctx, &target)
	err, errStr := runSSHOnce(ctx, target, args, password)
	if err == nil {
		return nil
	}

	if isHostKeyChangedError(errStr) {
		return retrySSHAfterHostKeyRecovery(ctx, target, args, password, errStr)
	}

	checkAndPrintDiagnosticFooter(target, err, errStr)

	return err
}

func SpawnSSH(ctx context.Context, target SSHTarget, args []string) error {
	return SpawnSSHWithPassword(ctx, target, args, "")
}

func PromptSSHPassword(ctx context.Context, prompt string, fd int) (string, error) {
	fmt.Print(prompt)
	password, err := term.ReadPassword(fd)
	fmt.Println()
	if err != nil {
		return "", &apperror.AppError{
			Op:       "PromptSSHPassword",
			Code:     "E_INTERNAL_ERROR",
			Type:     apperror.ErrorTypeExecution,
			Severity: apperror.SeverityError,
			Caller:   apperror.CaptureCaller(apperror.DefaultCallerSkip),
			Stack:    apperror.CaptureStackTrace(apperror.DefaultStackTraceSkip),
			Cause:    err,
		}
	}

	return strings.TrimRight(string(password), "\r\n"), nil
}

// NewAutoAcceptHostKeyConfig creates an ssh.ClientConfig that auto-accepts host keys on first join.
func NewAutoAcceptHostKeyConfig(user string, auth []ssh.AuthMethod) *ssh.ClientConfig {
	return &ssh.ClientConfig{
		User:            user,
		Auth:            auth,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}
}
