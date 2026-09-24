package cmdssh

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"golang.org/x/crypto/ssh"
)

// SSHInstallExecOptions represents parsed options for remote installer execution.
type SSHInstallExecOptions struct {
	SetupPath     string
	InstallerArgs []string
	Except        string
	ExceptOS      string
	TargetOS      string
	Target        string
	DestDir       string
	IsSilent      bool
	IsDryRun      bool
	IsForceAll    bool
	IsShowHelp    bool
}

// NodeInstallExecResult records the result of executing an installer on a remote node.
type NodeInstallExecResult struct {
	Alias      string `json:"alias"`
	Host       string `json:"host"`
	OS         string `json:"os"`
	RemotePath string `json:"remotePath"`
	ExitCode   int    `json:"exitCode"`
	Stdout     string `json:"stdout"`
	DurationMs int64  `json:"durationMs"`
	Status     string `json:"status"`
	Error      string `json:"error,omitempty"`
}

func printSSHInstallExecHelp() {
	fmt.Printf("\n%sUsage:%s gitmap ssh install-exec <setup-file> [installer-args...] [flags]\n\n", constants.ColorCyan, constants.ColorReset)
	fmt.Println("Deploy and run installer setups (.exe, .msi, .sh, binary) on remote SSH machines.")
	fmt.Println()
	fmt.Println("Aliases:")
	fmt.Println("  gitmap ssh install-exec")
	fmt.Println("  gitmap ssh in-exec")
	fmt.Println("  gitmap ssh setup-exec")
	fmt.Println("  gitmap ssh install-run")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  -t, --target string     Target machine alias or IP (default: all matching online nodes)")
	fmt.Println("      --except string     Exclude machines by ID (1, worker-1), alias, or IP")
	fmt.Println("      --except-os string  Exclude machines by OS (e.g. unix, win, linux, ubuntu, darwin)")
	fmt.Println("      --os string         Target machines by OS (e.g. win, unix, linux, darwin)")
	fmt.Println("  -f, --force-all         Deploy installer to all machines regardless of default OS target")
	fmt.Println("      --dest string       Remote directory for setup (default: Windows Temp directory, /tmp on Unix)")
	fmt.Println("  -s, --silent            Append silent unattended switch if not already provided")
	fmt.Println("  -n, --dry-run           Simulate file transfer and execution without running")
	fmt.Println("  -h, --help              Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  gitmap ssh install-exec ./setup.exe /SILENT (auto-targets Windows nodes)")
	fmt.Println("  gitmap ssh install-exec ./setup.exe --force-all (targets all nodes regardless of OS)")
	fmt.Println("  gitmap ssh install-exec ./bootstrap.sh (auto-targets Unix nodes)")
	fmt.Println("  gitmap ssh install-exec ./agent-installer.exe /qn --except worker-1,10.20.0.15")
	fmt.Println("  gitmap ssh install-exec ./app.msi --os win --silent")
	fmt.Println()
}

// ParseInstallExecArgs parses arguments and flags for install-exec.
func ParseInstallExecArgs(args []string) SSHInstallExecOptions {
	opts := SSHInstallExecOptions{
		IsSilent: true,
	}
	var positional []string

	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "-h" || a == "--help" {
			opts.IsShowHelp = true
			return opts
		}
		if a == "-n" || a == "--dry-run" {
			opts.IsDryRun = true
			continue
		}
		if a == "-s" || a == "--silent" {
			opts.IsSilent = true
			continue
		}
		if a == "--no-silent" || a == "--gui" || a == "--interactive" {
			opts.IsSilent = false
			continue
		}
		if a == "-f" || a == "--force-all" || a == "force-all" || a == "--all-os" {
			opts.IsForceAll = true
			continue
		}
		if (a == "--except" || a == "--excep") && i+1 < len(args) {
			opts.Except = args[i+1]
			i++
			continue
		}
		if strings.HasPrefix(a, "--except=") || strings.HasPrefix(a, "--excep=") {
			opts.Except = strings.TrimPrefix(strings.TrimPrefix(a, "--except="), "--excep=")
			continue
		}
		if (a == "--except-os" || a == "--exclude-os" || a == "--skip-os") && i+1 < len(args) {
			opts.ExceptOS = args[i+1]
			i++
			continue
		}
		if strings.HasPrefix(a, "--except-os=") {
			opts.ExceptOS = strings.TrimPrefix(a, "--except-os=")
			continue
		}
		if (a == "--os" || a == "--target-os") && i+1 < len(args) {
			opts.TargetOS = args[i+1]
			i++
			continue
		}
		if strings.HasPrefix(a, "--os=") {
			opts.TargetOS = strings.TrimPrefix(a, "--os=")
			continue
		}
		if (a == "-t" || a == "--target") && i+1 < len(args) {
			opts.Target = args[i+1]
			i++
			continue
		}
		if strings.HasPrefix(a, "--target=") {
			opts.Target = strings.TrimPrefix(a, "--target=")
			continue
		}
		if (a == "-d" || a == "--dest") && i+1 < len(args) {
			opts.DestDir = args[i+1]
			i++
			continue
		}
		if strings.HasPrefix(a, "--dest=") {
			opts.DestDir = strings.TrimPrefix(a, "--dest=")
			continue
		}

		positional = append(positional, a)
	}

	if len(positional) > 0 {
		opts.SetupPath = positional[0]
		opts.InstallerArgs = positional[1:]
	}

	if opts.TargetOS != "" || opts.ExceptOS != "" || opts.IsForceAll || opts.SetupPath == "" {
		return opts
	}

	ext := strings.ToLower(filepath.Ext(opts.SetupPath))
	if isWindowsInstallerExt(ext) {
		opts.TargetOS = constants.OSTargetWin
	} else if isUnixInstallerExt(ext) {
		opts.TargetOS = constants.OSTargetUnix
	}

	return opts
}

func isWindowsInstallerExt(ext string) bool {
	switch ext {
	case ".exe", ".msi", ".bat", ".cmd", ".ps1":
		return true
	default:
		return false
	}
}

func isUnixInstallerExt(ext string) bool {
	switch ext {
	case ".sh", ".bash", ".deb", ".rpm", ".bin", ".run", ".appimage":
		return true
	default:
		return false
	}
}

// BuildRemoteInstallerExecCmd constructs the OS-specific remote execution command.
func BuildRemoteInstallerExecCmd(osType, remotePath string, installerArgs []string, isSilent bool) (string, string) {
	return BuildRemoteInstallerExecCmdWithPayload(osType, remotePath, installerArgs, isSilent, nil)
}

// BuildRemoteInstallerExecCmdWithPayload constructs the OS-specific remote execution command using payload inspection.
func BuildRemoteInstallerExecCmdWithPayload(osType, remotePath string, installerArgs []string, isSilent bool, payload []byte) (string, string) {
	isWin := isWindowsOS(osType)
	ext := strings.ToLower(filepath.Ext(remotePath))
	argsStr := strings.Join(installerArgs, " ")

	if isSilent {
		silentFlag := resolveDefaultSilentInstallerArgs(ext, payload)
		argsStr = injectSilentFlag(argsStr, silentFlag, installerArgs)
	}

	if isWin {
		return buildWindowsInstallerExecCmd(ext, remotePath, argsStr)
	}

	// Linux / Unix / macOS
	cmd := fmt.Sprintf("chmod +x '%s' && '%s' %s", remotePath, remotePath, argsStr)
	return cmd, "bash"
}

func injectSilentFlag(argsStr, silentFlag string, args []string) string {
	if silentFlag == "" || containsSilentArg(args) {
		return argsStr
	}
	if argsStr == "" {
		return silentFlag
	}
	return silentFlag + " " + argsStr
}

func containsSilentArg(args []string) bool {
	for _, a := range args {
		if isSilentFlagMatch(a) {
			return true
		}
	}
	return false
}

func isSilentFlagMatch(arg string) bool {
	low := strings.ToLower(arg)
	switch low {
	case "/s", "-s", "--silent", "-silent", "/silent", "/verysilent", "-quiet", "--quiet", "/quiet", "/qn", "-qn":
		return true
	default:
		return false
	}
}

func resolveDefaultSilentInstallerArgs(ext string, payload []byte) string {
	if ext == ".msi" {
		return "/qn"
	}
	if ext != ".exe" {
		return ""
	}
	if bytes.Contains(payload, []byte("NullsoftInst")) {
		return "/S"
	}
	if bytes.Contains(payload, []byte("Inno Setup")) {
		return "/VERYSILENT /SUPPRESSMSGBOXES /NORESTART /SP-"
	}
	if bytes.Contains(payload, []byte("WixBundle")) || bytes.Contains(payload, []byte("Burn")) {
		return "/quiet /norestart"
	}
	return "/S"
}

func buildWindowsInstallerExecCmd(ext, remotePath, argsStr string) (string, string) {
	switch ext {
	case ".msi":
		if argsStr == "" {
			return fmt.Sprintf("$p = '%s'; $proc = Start-Process msiexec.exe -ArgumentList @('/i', $p) -Wait -PassThru; exit $proc.ExitCode", remotePath), "ps"
		}
		return fmt.Sprintf("$p = '%s'; $proc = Start-Process msiexec.exe -ArgumentList @('/i', $p, '%s') -Wait -PassThru; exit $proc.ExitCode", remotePath, argsStr), "ps"
	case ".ps1":
		return fmt.Sprintf("powershell -NoProfile -ExecutionPolicy Bypass -File \"%s\" %s", remotePath, argsStr), "cmd"
	case ".bat", ".cmd":
		return fmt.Sprintf("\"%s\" %s", remotePath, argsStr), "cmd"
	default:
		if argsStr == "" {
			return fmt.Sprintf("$p = '%s'; $proc = Start-Process -FilePath $p -Wait -PassThru; exit $proc.ExitCode", remotePath), "ps"
		}
		return fmt.Sprintf("$p = '%s'; $proc = Start-Process -FilePath $p -ArgumentList '%s' -Wait -PassThru; exit $proc.ExitCode", remotePath, argsStr), "ps"
	}
}

// RunSSHInstallExecCLI handles the 'gitmap ssh install-exec' command.
func RunSSHInstallExecCLI(args []string) error {
	opts := ParseInstallExecArgs(args)
	if opts.IsShowHelp {
		printSSHInstallExecHelp()
		return nil
	}
	if opts.SetupPath == "" {
		printSSHInstallExecHelp()
		return apperror.NewValidationError("missing required <setup-file> argument")
	}

	fileInfo, errStat := os.Stat(opts.SetupPath)
	if errStat != nil {
		return apperror.WrapSimple(errStat, fmt.Sprintf("local setup file not found: %s", opts.SetupPath))
	}
	if fileInfo.IsDir() {
		return apperror.NewValidationError(fmt.Sprintf("target path is a directory, expected executable file: %s", opts.SetupPath))
	}

	fileBytes, errRead := os.ReadFile(opts.SetupPath)
	if errRead != nil {
		return apperror.WrapSimple(errRead, fmt.Sprintf("failed to read setup file: %s", opts.SetupPath))
	}

	conns, errLoad := fetchAllSSHConnections()
	if errLoad != nil {
		return errLoad
	}

	if opts.Target != "" && opts.Target != "all" {
		conns = filterConnectionsByTarget(conns, opts.Target)
	}
	if opts.Except != "" {
		conns = FilterSSHConnectionsByExcept(conns, opts.Except)
	}
	if opts.ExceptOS != "" || opts.TargetOS != "" {
		conns = FilterSSHConnectionsByOS(conns, opts.TargetOS, opts.ExceptOS)
	}

	if len(conns) == 0 {
		return renderNoMatchingMachinesMessage(opts)
	}

	fileName := filepath.Base(opts.SetupPath)
	fmt.Println()
	fmt.Printf("  %s● Deploying installer '%s' across %d machine(s)%s\n", constants.ColorCyan, fileName, len(conns), constants.ColorReset)
	if opts.IsDryRun {
		fmt.Printf("    %s[dry-run] Would stream %d bytes and execute with args: %v%s\n\n", constants.ColorDim, len(fileBytes), opts.InstallerArgs, constants.ColorReset)
		for _, c := range conns {
			fmt.Printf("    - %s (%s, os: %s)\n", c.Alias, c.IPAddress, c.OS)
		}
		fmt.Println()
		return nil
	}

	return executeInstallerAcrossFleet(conns, fileName, fileBytes, opts)
}

func executeInstallerAcrossFleet(conns []db.SSHConnection, fileName string, data []byte, opts SSHInstallExecOptions) error {
	results := make([]NodeInstallExecResult, len(conns))
	var wg sync.WaitGroup

	for i, c := range conns {
		wg.Add(1)
		go func(idx int, conn db.SSHConnection) {
			defer wg.Done()
			results[idx] = executeInstallerOnSingleNode(conn, fileName, data, opts)
		}(i, c)
	}

	wg.Wait()
	renderInstallExecResultsTable(fileName, results)
	return nil
}

func executeInstallerOnSingleNode(c db.SSHConnection, fileName string, data []byte, opts SSHInstallExecOptions) NodeInstallExecResult {
	start := time.Now()
	res := NodeInstallExecResult{
		Alias: c.Alias,
		Host:  fmt.Sprintf("%s:22", c.IPAddress),
		OS:    c.OS,
	}

	isOnline, _ := CheckConnLiveness(context.Background(), c.IPAddress, 22, 0)
	if !isOnline {
		res.Status = "○ OFFLINE"
		res.Error = "machine is offline"
		return res
	}

	client, isConnected := connectSSHClient(c, fmt.Sprintf("[%s]", c.Alias))
	if !isConnected {
		res.Status = "✗ AUTH FAILED"
		res.Error = "authentication failed"
		return res
	}
	defer client.Close()

	probedOS := probeRemoteOSType(client)
	if probedOS != "" && probedOS != c.OS {
		c.OS = probedOS
		persistUpdatedNodeOS(c)
	}
	if probedOS != "" {
		res.OS = probedOS
	}

	remoteDestPath := resolveRemoteTempPath(c.OS, fileName, opts.DestDir)
	res.RemotePath = remoteDestPath

	killStaleInstallerProcess(client, fileName, c.OS)

	// 1. Stream binary over pure SSH channel
	errStream := StreamFileToRemote(client, remoteDestPath, data, c.OS)
	if errStream != nil {
		res.Status = "✗ UPLOAD FAILED"
		res.Error = errStream.Error()
		return res
	}

	// 2. Build and execute installation command
	execCmd, execShell := BuildRemoteInstallerExecCmdWithPayload(c.OS, remoteDestPath, opts.InstallerArgs, opts.IsSilent, data)
	stdout, errExec := crypto.RunCommand(client, execCmd, execShell)
	dur := time.Since(start).Milliseconds()
	res.DurationMs = dur
	res.Stdout = strings.TrimSpace(stdout)

	if errExec != nil {
		res.ExitCode = resolveProcessExitCode(errExec)
		res.Status = fmt.Sprintf("✗ FAILED (%d)", res.ExitCode)
		res.Error = errExec.Error()
	} else {
		res.ExitCode = 0
		res.Status = "✓ INSTALLED (0)"
	}

	return res
}

func killStaleInstallerProcess(client *ssh.Client, fileName, osType string) {
	if client == nil || fileName == "" {
		return
	}
	base := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	if isWindowsOS(osType) {
		cmd := fmt.Sprintf("cmd.exe /c taskkill /F /IM \"%s*\" 2>nul", base)
		_, _ = crypto.RunCommand(client, cmd, "")
		return
	}
	cmd := fmt.Sprintf("pkill -9 -f '%s' || true", fileName)
	_, _ = crypto.RunCommand(client, cmd, "bash")
}

func persistUpdatedNodeOS(c db.SSHConnection) {
	ctx := context.Background()
	dbConn, err := store.OpenDefault()
	if err == nil {
		defer dbConn.Close()
		_ = db.UpdateConnectionOS(ctx, dbConn.SQL(), c.Alias, c.OS)
	}
	globalConn, gErr := store.OpenGlobalDefault()
	if gErr == nil {
		defer globalConn.Close()
		_ = db.UpdateConnectionOS(ctx, globalConn.SQL(), c.Alias, c.OS)
	}
}

func resolveRemoteTempPath(osType, fileName, customDest string) string {
	if customDest != "" && isWindowsOS(osType) {
		return strings.TrimRight(customDest, "\\/") + "\\" + fileName
	}
	if customDest != "" {
		return strings.TrimRight(customDest, "\\/") + "/" + fileName
	}

	if isWindowsOS(osType) {
		return "C:\\Windows\\Temp\\" + fileName
	}
	return "/tmp/" + fileName
}

func renderInstallExecResultsTable(fileName string, results []NodeInstallExecResult) {
	fmt.Println()
	fmt.Printf("  %-16s %-20s %-10s %-26s %-16s %s\n",
		"ALIAS", "HOST (IP)", "OS", "REMOTE PATH", "RESULT", "DURATION")
	fmt.Println("  " + strings.Repeat("-", 100))

	for _, r := range results {
		statusColor := constants.ColorGreen
		if strings.Contains(r.Status, "FAILED") || strings.Contains(r.Status, "OFFLINE") {
			statusColor = constants.ColorRed
		}
		fmt.Printf("  %-16s %-20s %-10s %-26s %s%-16s%s %dms\n",
			r.Alias,
			r.Host,
			r.OS,
			r.RemotePath,
			statusColor, r.Status, constants.ColorReset,
			r.DurationMs,
		)
		if r.Error != "" && !strings.Contains(r.Status, "✓") {
			fmt.Printf("    %sError: %s%s\n", constants.ColorRed, r.Error, constants.ColorReset)
		}
		if r.Stdout != "" && !strings.Contains(r.Status, "✓") {
			fmt.Printf("    %sOutput: %s%s\n", constants.ColorDim, r.Stdout, constants.ColorReset)
		}
	}
	fmt.Println()
}

func renderNoMatchingMachinesMessage(opts SSHInstallExecOptions) error {
	allConns, _ := fetchAllSSHConnections()
	if len(allConns) == 0 {
		fmt.Println("No registered SSH machines found in registry.")
		fmt.Println("Enroll machines first using: gitmap sjc <user> <ip(alias)...> --pass <password> or gitmap sj add <user@ip> [alias]")
		return nil
	}
	fmt.Printf("No matching SSH machines to deploy setup to (filtered by except: %q, except-os: %q, os: %q).\n", opts.Except, opts.ExceptOS, opts.TargetOS)
	return nil
}
