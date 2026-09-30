// Package cmdnodes coordinates fleet-wide operations across unified cluster and SSH nodes.
package cmdnodes

import (
	"context"
	"fmt"
	"net"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdclone"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func filterRemoteConnections(conns []db.SSHConnection, opts NodesCloneOptions) []db.SSHConnection {
	var remote []db.SSHConnection
	for _, c := range conns {
		if isLocalMachineConnection(c) {
			continue
		}
		if isConnectionExcluded(c, opts) {
			continue
		}
		remote = append(remote, c)
	}
	return remote
}

func isConnectionExcluded(c db.SSHConnection, opts NodesCloneOptions) bool {
	if opts.TargetFilter != "" && !strings.EqualFold(c.Alias, opts.TargetFilter) && c.IPAddress != opts.TargetFilter {
		return true
	}
	if opts.ExcludeFilter != "" && (strings.EqualFold(c.Alias, opts.ExcludeFilter) || c.IPAddress == opts.ExcludeFilter) {
		return true
	}
	return false
}

func isLocalMachineConnection(c db.SSHConnection) bool {
	if c.IPAddress == "127.0.0.1" || strings.EqualFold(c.IPAddress, "localhost") {
		return true
	}
	if strings.EqualFold(c.Alias, "local") || strings.EqualFold(c.Alias, "localhost") {
		return true
	}
	return isLocalHostIP(c.IPAddress)
}

func isLocalHostIP(ip string) bool {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return false
	}
	for _, a := range addrs {
		ipNet, ok := a.(*net.IPNet)
		if ok && ipNet.IP.String() == ip {
			return true
		}
	}
	return false
}

func resolveRepoAuthURL(rawURL string) string {
	rec := model.ScanRecord{HTTPSUrl: rawURL, RepoName: filepath.Base(rawURL)}
	resolved, err := cmdclone.ResolveRepoAuth(rec)
	if err == nil && resolved.HTTPSUrl != "" {
		return resolved.HTTPSUrl
	}
	return rawURL
}

func resolveRemoteArgItem(arg, detectedFile, fileName string) string {
	if arg == detectedFile || filepath.Base(arg) == fileName {
		return fileName
	}
	if strings.HasPrefix(arg, "http://") || strings.HasPrefix(arg, "https://") {
		return resolveRepoAuthURL(arg)
	}
	return arg
}

func resolveRemoteArgs(passArgs []string, detectedFile, fileName string) string {
	if len(passArgs) == 0 {
		return fileName
	}
	var remoteArgs []string
	for _, a := range passArgs {
		remoteArgs = append(remoteArgs, resolveRemoteArgItem(a, detectedFile, fileName))
	}
	return strings.TrimSpace(strings.Join(remoteArgs, " "))
}

func isWindowsNode(conn db.SSHConnection) bool {
	if strings.EqualFold(conn.OS, "windows") || strings.EqualFold(conn.OS, "win") {
		return true
	}
	if strings.EqualFold(conn.OSGroup, "windows") {
		return true
	}
	if strings.EqualFold(conn.Username, "Administrator") {
		return true
	}
	return false
}

func appendJSONFlag(args string) string {
	if strings.Contains(args, "--json") || strings.Contains(args, " -j") || args == "-j" {
		return args
	}
	if args == "" {
		return "--json"
	}
	return args + " --json"
}

func buildWindowsWorkDirExecString(kindStr, args, targetDir string) string {
	workDir := `D:\work`
	if targetDir != "" {
		workDir = targetDir
	}
	execArgs := appendJSONFlag(args)
	return fmt.Sprintf("Set-Location \"%s\"; gitmap %s %s", workDir, kindStr, execArgs)
}

func buildUnixWorkDirExecString(kindStr, args, targetDir string) string {
	workDir := "~/work"
	if targetDir != "" {
		workDir = targetDir
	}
	execArgs := appendJSONFlag(args)
	return fmt.Sprintf("cd %s && gitmap %s %s", workDir, kindStr, execArgs)
}

func buildRemoteExecString(opts NodesCloneOptions, fileName string, isWin bool) string {
	cmdArgs := resolveRemoteArgs(opts.PassArgs, opts.DetectedFile, fileName)
	kindStr := string(opts.Kind)
	if isWin {
		return buildWindowsWorkDirExecString(kindStr, cmdArgs, opts.TargetDir)
	}
	return buildUnixWorkDirExecString(kindStr, cmdArgs, opts.TargetDir)
}

func isOfflineError(errStr string) bool {
	low := strings.ToLower(errStr)
	return strings.Contains(low, "network unreachable") ||
		strings.Contains(low, "offline") ||
		strings.Contains(low, "timeout") ||
		strings.Contains(low, "i/o timeout") ||
		strings.Contains(low, "connectex") ||
		strings.Contains(low, "connection refused") ||
		strings.Contains(low, "actively refused") ||
		strings.Contains(low, "machine is off") ||
		strings.Contains(low, "no route to host") ||
		strings.Contains(low, "host is down") ||
		strings.Contains(low, "unreachable host") ||
		strings.Contains(low, "getaddrinfo") ||
		strings.Contains(low, "name resolution")
}

func classifyConnectFailure(res RemoteCloneNodeResult, errConnect error, start time.Time) RemoteCloneNodeResult {
	res.Duration = time.Since(start)
	res.DurationMs = res.Duration.Milliseconds()
	errStr := errConnect.Error()
	if isOfflineError(errStr) {
		res.Status = "offline"
		res.Error = "node offline or unreachable"
		return res
	}
	res.Status = "auth_failed"
	res.Error = "credentials rejected by remote host"
	return res
}

func runRemoteNodeWorker(conn db.SSHConnection, opts NodesCloneOptions, fileBytes []byte, fileName string) RemoteCloneNodeResult {
	start := time.Now()
	res := RemoteCloneNodeResult{Alias: conn.Alias, Host: conn.IPAddress, Role: "worker"}
	client, errConnect := cmdssh.ConnectSSHClientWithErr(conn)
	if errConnect != nil {
		return classifyConnectFailure(res, errConnect, start)
	}
	defer client.Close()
	return executeRemoteCloneSession(client, conn, opts, fileBytes, fileName, res, start)
}

func stageRemoteFileIfNeeded(client *ssh.Client, conn db.SSHConnection, opts NodesCloneOptions, fileBytes []byte, fileName string) error {
	if !opts.HasFile || len(fileBytes) == 0 {
		return nil
	}
	_, err := StageFileToRemoteNode(client, conn, fileName, fileBytes)
	return err
}

func executeRemoteCloneSession(client *ssh.Client, conn db.SSHConnection, opts NodesCloneOptions, fileBytes []byte, fileName string, res RemoteCloneNodeResult, start time.Time) RemoteCloneNodeResult {
	if errStage := stageRemoteFileIfNeeded(client, conn, opts, fileBytes, fileName); errStage != nil {
		res.Status = "failed"
		res.Error = "manifest staging failed: " + errStage.Error()
		return res
	}
	return runRemoteExecOverSSH(client, conn, opts, fileName, res, start)
}

func isBashMissingError(out string, err error) bool {
	if err != nil && strings.Contains(err.Error(), "'bash' is not recognized") {
		return true
	}
	if strings.Contains(out, "'bash' is not recognized") || strings.Contains(out, "bash: command not found") {
		return true
	}
	return false
}

func persistCorrectedNodeOS(alias, osType string) {
	dbConn, err := store.OpenDefault()
	if err != nil {
		return
	}
	defer dbConn.Close()
	_ = db.UpdateConnectionOS(context.Background(), dbConn.Conn(), alias, osType)
}

func retryWithPowerShell(client *ssh.Client, opts NodesCloneOptions, fileName, alias string) (string, error) {
	cmdStr := buildRemoteExecString(opts, fileName, true)
	out, err := crypto.RunCommand(client, cmdStr, "ps")
	if err == nil {
		go persistCorrectedNodeOS(alias, "windows")
	}
	return out, err
}

func populateFromJSON(res *RemoteCloneNodeResult, payload cmdclone.DirectCloneJSONResponse) {
	res.Details = payload.Message
	if payload.Success {
		res.Status = "success"
		return
	}
	res.Status = "failed"
	res.Error = payload.Message
}

func isJSONFlagUnsupportedError(out string) bool {
	return strings.Contains(out, "flag provided but not defined: -json") ||
		strings.Contains(out, "flag provided but not defined: --json")
}

func buildLegacyWindowsExec(kindStr, cmdArgs, targetDir string) string {
	workDir := `D:\work`
	if targetDir != "" {
		workDir = targetDir
	}
	if cmdArgs == "" {
		return fmt.Sprintf("Set-Location \"%s\"; gitmap %s", workDir, kindStr)
	}
	return fmt.Sprintf("Set-Location \"%s\"; gitmap %s %s", workDir, kindStr, cmdArgs)
}

func buildLegacyUnixExec(kindStr, cmdArgs, targetDir string) string {
	workDir := "~/work"
	if targetDir != "" {
		workDir = targetDir
	}
	if cmdArgs == "" {
		return fmt.Sprintf("cd %s && gitmap %s", workDir, kindStr)
	}
	return fmt.Sprintf("cd %s && gitmap %s %s", workDir, kindStr, cmdArgs)
}

func buildLegacyExecString(kindStr, cmdArgs, targetDir string, isWin bool) string {
	if isWin {
		return buildLegacyWindowsExec(kindStr, cmdArgs, targetDir)
	}
	return buildLegacyUnixExec(kindStr, cmdArgs, targetDir)
}

func retryWithoutJSON(client *ssh.Client, opts NodesCloneOptions, fileName string, isWin bool, shell string) (string, error) {
	cmdArgs := resolveRemoteArgs(opts.PassArgs, opts.DetectedFile, fileName)
	cmdStr := buildLegacyExecString(string(opts.Kind), cmdArgs, opts.TargetDir, isWin)
	return crypto.RunCommand(client, cmdStr, shell)
}

func extractLegacyCloneDetails(out string) string {
	if strings.Contains(out, "already exists on disk") {
		return "already exists on disk"
	}
	if strings.Contains(out, "cloned") {
		return "cloned successfully"
	}
	return ""
}

func populateExecutionResult(res RemoteCloneNodeResult, out string, errRun error, start time.Time) RemoteCloneNodeResult {
	res.Duration = time.Since(start)
	res.DurationMs = res.Duration.Milliseconds()
	res.Stdout = strings.TrimSpace(out)
	if payload, hasJSON := cmdclone.ParseCloneJSONResponse(out); hasJSON {
		populateFromJSON(&res, payload)
		return res
	}
	if errRun != nil {
		res.Status = "failed"
		res.Error = errRun.Error()
		return res
	}
	res.Status = "success"
	res.Details = extractLegacyCloneDetails(out)
	return res
}

func runRemoteExecOverSSH(client *ssh.Client, conn db.SSHConnection, opts NodesCloneOptions, fileName string, res RemoteCloneNodeResult, start time.Time) RemoteCloneNodeResult {
	isWin := isWindowsNode(conn)
	shell := "bash"
	if isWin {
		shell = "ps"
	}
	cmdStr := buildRemoteExecString(opts, fileName, isWin)
	out, errRun := crypto.RunCommand(client, cmdStr, shell)
	if errRun != nil && !isWin && isBashMissingError(out, errRun) {
		out, errRun = retryWithPowerShell(client, opts, fileName, conn.Alias)
	}
	if isJSONFlagUnsupportedError(out) {
		out, errRun = retryWithoutJSON(client, opts, fileName, isWin, shell)
	}
	return populateExecutionResult(res, out, errRun, start)
}

func executeFleetNodesParallel(conns []db.SSHConnection, opts NodesCloneOptions, fileBytes []byte, fileName string) []RemoteCloneNodeResult {
	var wg sync.WaitGroup
	var mu sync.Mutex
	results := make([]RemoteCloneNodeResult, 0, len(conns))
	for _, c := range conns {
		wg.Add(1)
		go func(conn db.SSHConnection) {
			defer wg.Done()
			r := runRemoteNodeWorker(conn, opts, fileBytes, fileName)
			mu.Lock()
			results = append(results, r)
			mu.Unlock()
		}(c)
	}
	wg.Wait()
	return results
}
