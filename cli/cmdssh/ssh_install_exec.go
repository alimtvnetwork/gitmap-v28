package cmdssh

import (
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
	fmt.Println("      --dest string       Remote directory for setup (default: %TEMP% on Windows, /tmp on Unix)")
	fmt.Println("  -s, --silent            Append silent unattended switch if not already provided")
	fmt.Println("  -n, --dry-run           Simulate file transfer and execution without running")
	fmt.Println("  -h, --help              Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  gitmap ssh install-exec ./setup.exe /SILENT --except-os unix")
	fmt.Println("  gitmap ssh install-exec ./agent-installer.exe /qn --except worker-1,10.20.0.15")
	fmt.Println("  gitmap ssh install-exec ./bootstrap.sh --except-os win")
	fmt.Println("  gitmap ssh install-exec ./app.msi --os win --silent")
	fmt.Println()
}

// ParseInstallExecArgs parses arguments and flags for install-exec.
func ParseInstallExecArgs(args []string) SSHInstallExecOptions {
	var opts SSHInstallExecOptions
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

	return opts
}

// BuildRemoteInstallerExecCmd constructs the OS-specific remote execution command.
func BuildRemoteInstallerExecCmd(osType, remotePath string, installerArgs []string, isSilent bool) (string, string) {
	isWin := isWindowsOS(osType)
	ext := strings.ToLower(filepath.Ext(remotePath))
	argsStr := strings.Join(installerArgs, " ")

	if isSilent && argsStr == "" {
		switch ext {
		case ".msi":
			argsStr = "/qn"
		case ".exe":
			argsStr = "/VERYSILENT /SUPPRESSMSGBOXES /NORESTART /SP-"
		}
	}

	if isWin {
		if ext == ".msi" {
			cmd := fmt.Sprintf(`powershell -NoProfile -Command "$p = '%s'; $proc = Start-Process msiexec.exe -ArgumentList @('/i', "`"$p`"", '%s') -Wait -PassThru; exit $proc.ExitCode"`, remotePath, argsStr)
			return cmd, "ps"
		}
		if ext == ".ps1" {
			cmd := fmt.Sprintf(`powershell -NoProfile -ExecutionPolicy Bypass -File '%s' %s`, remotePath, argsStr)
			return cmd, "ps"
		}
		if ext == ".bat" || ext == ".cmd" {
			cmd := fmt.Sprintf(`cmd.exe /c "%s %s"`, remotePath, argsStr)
			return cmd, "cmd"
		}
		// Default Windows .exe execution
		cmd := fmt.Sprintf(`powershell -NoProfile -Command "$p = '%s'; $proc = Start-Process -FilePath $p -ArgumentList '%s' -Wait -PassThru; exit $proc.ExitCode"`, remotePath, argsStr)
		return cmd, "ps"
	}

	// Linux / Unix / macOS
	cmd := fmt.Sprintf(`chmod +x '%s' && '%s' %s`, remotePath, remotePath, argsStr)
	return cmd, "bash"
}

// RunSSHInstallExecCLI handles the 'gitmap ssh install-exec' command.
func RunSSHInstallExecCLI(args []string) error {
	opts := ParseInstallExecArgs(args)
	if opts.IsShowHelp || opts.SetupPath == "" {
		printSSHInstallExecHelp()
		if opts.SetupPath == "" && !opts.IsShowHelp {
			return apperror.NewValidationError("missing required <setup-file> argument")
		}
		return nil
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
		fmt.Printf("No matching SSH machines to deploy setup to (filtered by except: %q, except-os: %q, os: %q).\n", opts.Except, opts.ExceptOS, opts.TargetOS)
		return nil
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

	if c.OS == "" {
		c.OS = probeRemoteOSType(client)
		res.OS = c.OS
	}

	remoteDestPath := resolveRemoteTempPath(c.OS, fileName, opts.DestDir)
	res.RemotePath = remoteDestPath

	// 1. Stream binary over pure SSH channel
	writeCmd := buildRemoteWriteCmd(remoteDestPath, data, isWindowsOS(c.OS))
	shell := determineFallbackShell(c.OS)
	_, errWrite := crypto.RunCommand(client, writeCmd, shell)
	if errWrite != nil {
		res.Status = "✗ UPLOAD FAILED"
		res.Error = errWrite.Error()
		return res
	}

	// 2. Build and execute installation command
	execCmd, execShell := BuildRemoteInstallerExecCmd(c.OS, remoteDestPath, opts.InstallerArgs, opts.IsSilent)
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

func resolveRemoteTempPath(osType, fileName, customDest string) string {
	if customDest != "" {
		if isWindowsOS(osType) {
			return strings.TrimRight(customDest, "\\/") + "\\" + fileName
		}
		return strings.TrimRight(customDest, "\\/") + "/" + fileName
	}

	if isWindowsOS(osType) {
		return fmt.Sprintf(`C:\Windows\Temp\%s`, fileName)
	}
	return fmt.Sprintf(`/tmp/%s`, fileName)
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
		if r.Stdout != "" && !strings.Contains(r.Status, "✓") {
			fmt.Printf("    %sOutput: %s%s\n", constants.ColorDim, r.Stdout, constants.ColorReset)
		}
	}
	fmt.Println()
}
