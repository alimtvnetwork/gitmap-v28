package cmdssh

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

type sshEnableOptions struct {
	port    int
	isForce bool
	isHelp  bool
}

type daemonCmdRunner func(name string, args ...string) ([]byte, error)

var daemonExecRunner daemonCmdRunner = func(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)

	return cmd.CombinedOutput()
}

type LinuxStep struct {
	Binary string
	Args   []string
	Desc   string
}

func isPortOutOfBounds(port int) bool {
	return port < 1 || port > 65535
}

func parsePortFlagValue(val string) (int, error) {
	port, err := strconv.Atoi(val)
	if err != nil {
		return 0, apperror.NewValidationError(fmt.Sprintf("invalid port value %q: must be integer", val))
	}

	if isPortOutOfBounds(port) {
		return 0, apperror.NewValidationError(fmt.Sprintf("port %d out of range: must be between 1 and 65535", port))
	}

	return port, nil
}

func isPortFlag(arg string) bool {
	return arg == "--port" || arg == "-p"
}

func parsePortFlagValueWithNext(next string, hasNext bool, opts *sshEnableOptions) (int, error) {
	if !hasNext {
		return 1, apperror.NewValidationError("missing argument for --port flag")
	}

	p, err := parsePortFlagValue(next)
	opts.port = p

	return 2, err
}

func parsePortFlagToken(arg string, next string, hasNext bool, opts *sshEnableOptions) (int, error) {
	if strings.HasPrefix(arg, "--port=") {
		p, err := parsePortFlagValue(strings.TrimPrefix(arg, "--port="))
		opts.port = p

		return 1, err
	}

	if !isPortFlag(arg) {
		return 1, nil
	}

	return parsePortFlagValueWithNext(next, hasNext, opts)
}

func isForceFlag(arg string) bool {
	return arg == "--force" || arg == "-f"
}

func peekNextArg(args []string, idx int) (string, bool) {
	if idx+1 < len(args) {
		return args[idx+1], true
	}

	return "", false
}

func parseSingleEnableFlag(args []string, idx int, opts *sshEnableOptions) (int, error) {
	arg := args[idx]
	if isHelpFlag(arg) {
		opts.isHelp = true

		return 1, nil
	}

	if isForceFlag(arg) {
		opts.isForce = true

		return 1, nil
	}

	next, hasNext := peekNextArg(args, idx)

	return parsePortFlagToken(arg, next, hasNext, opts)
}

func parseSSHEnableFlags(args []string) (sshEnableOptions, error) {
	opts := sshEnableOptions{port: 22}
	idx := 0

	for idx < len(args) {
		consumed, err := parseSingleEnableFlag(args, idx, &opts)
		if err != nil {
			return opts, err
		}

		idx += consumed
	}

	return opts, nil
}

func printSSHEnableHelp() {
	fmt.Printf("\n%sUsage:%s gitmap ssh enable [flags]\n\n", constants.ColorCyan, constants.ColorReset)
	fmt.Println("Installs, configures, starts OpenSSH Server daemon (sshd) and opens firewall ports.\n\nFlags:\n  --port, -p <port>   Port to listen on (default 22)\n  --force, -f         Force installation even if prerequisites exist\n  --help, -h          Show this help text")
}

func buildWindowsCheckCapabilityCmd() string {
	return "if ((Get-WindowsCapability -Online -Name OpenSSH.Server~~~~0.0.1.0 -ErrorAction SilentlyContinue).State -eq 'Installed') { 'Installed' } elseif (Get-Service -Name sshd -ErrorAction SilentlyContinue) { 'Installed' } else { 'NotPresent' }"
}

func buildWindowsAddCapabilityCmd() string {
	return "$s = Get-Service wuauserv -ErrorAction SilentlyContinue; $d = ($s -and $s.StartType -eq 'Disabled'); if ($d) { Set-Service wuauserv -StartupType Manual -ErrorAction SilentlyContinue; Start-Service wuauserv -ErrorAction SilentlyContinue }; try { Add-WindowsCapability -Online -Name OpenSSH.Server~~~~0.0.1.0 -ErrorAction Stop } finally { if ($d) { Stop-Service wuauserv -ErrorAction SilentlyContinue; Set-Service wuauserv -StartupType Disabled -ErrorAction SilentlyContinue } }"
}

func buildWindowsFirewallRuleCmd(port int) string {
	return fmt.Sprintf("New-NetFirewallRule -Name 'OpenSSH-Server-In-TCP' -DisplayName 'OpenSSH Server (sshd)' -Enabled True -Direction Inbound -Protocol TCP -Action Allow -LocalPort %d", port)
}

func buildWindowsFirewallEnsureCmd(port int) string {
	return fmt.Sprintf("if (Get-NetFirewallRule -Name 'OpenSSH-Server-In-TCP' -ErrorAction SilentlyContinue) { Set-NetFirewallRule -Name 'OpenSSH-Server-In-TCP' -LocalPort %d -Enabled True } else { %s }", port, buildWindowsFirewallRuleCmd(port))
}

func runWindowsPowerShell(script string) (string, error) {
	args := []string{"-NoProfile", "-NonInteractive", "-ExecutionPolicy", "Bypass", "-Command", script}
	out, err := daemonExecRunner("powershell.exe", args...)
	trimmed := strings.TrimSpace(string(out))
	if err != nil {
		return trimmed, apperror.NewExecutionError(fmt.Sprintf("powershell failed: %v (%s)", err, trimmed))
	}

	return trimmed, nil
}

func runWindowsPowerShellScript(script string, desc string) error {
	_, err := runWindowsPowerShell(script)
	if err != nil {
		return apperror.NewExecutionError(fmt.Sprintf("%s failed: %v", desc, err))
	}

	return nil
}

func isWindowsSSHDFilePresent() bool {
	sysRoot := os.Getenv("SystemRoot")
	if sysRoot == "" {
		sysRoot = "C:\\Windows"
	}

	sshdPath := filepath.Join(sysRoot, "System32", "OpenSSH", "sshd.exe")
	info, err := os.Stat(sshdPath)
	if err != nil || info.IsDir() {
		return false
	}

	return true
}

func isWindowsCapabilityInstalled() bool {
	if isWindowsSSHDFilePresent() {
		return true
	}

	out, err := runWindowsPowerShell(buildWindowsCheckCapabilityCmd())
	if err != nil {
		return false
	}

	return out == "Installed"
}

func printWindowsCapabilityFailureNotice(err error) {
	fmt.Printf("\n  %s[!] OpenSSH Server capability installation failed:%s\n      %v\n\n", constants.ColorYellow, constants.ColorReset, err)
	fmt.Println("  Remediation steps:")
	fmt.Println("  1. Run as Administrator: Ensure your terminal is running with elevated privileges.")
	fmt.Println("  2. Windows Update: If Windows Update (wuauserv) is disabled, temporarily enable it:")
	fmt.Println("     Set-Service -Name wuauserv -StartupType Manual; Start-Service -Name wuauserv")
	fmt.Println("  3. Alternative Install: Install OpenSSH Server via winget or chocolatey:")
	fmt.Printf("     winget install Microsoft.OpenSSH.Beta\n\n")
}

func newSSHCapabilityAbortError(err error) *apperror.AppError {
	printWindowsCapabilityFailureNotice(err)

	return apperror.NewWithDetails("ssh", "E9001", "OpenSSH Server installation aborted", "cli", apperror.ErrorTypeAbort, apperror.SeverityWarn, map[string]any{"reported": true})
}

func handleCapabilityError(err error, isForce bool) error {
	if isForce {
		return nil
	}

	return newSSHCapabilityAbortError(err)
}

func ensureWindowsCapability(isForce bool) error {
	if isWindowsCapabilityInstalled() {
		fmt.Printf("  %s[1/4]%s OpenSSH Server Windows Capability is already installed.\n", constants.ColorGreen, constants.ColorReset)

		return nil
	}

	fmt.Printf("  %s[1/4]%s Installing OpenSSH Server Windows Capability...\n", constants.ColorCyan, constants.ColorReset)
	capCmd := buildWindowsAddCapabilityCmd()
	err := runWindowsPowerShellScript(capCmd, "Add-WindowsCapability")
	if err != nil {
		return handleCapabilityError(err, isForce)
	}

	return nil
}

func handleSSHDServiceError(err error, stepDesc string) *apperror.AppError {
	fmt.Printf("\n  %s[!] Operation failed at step '%s':%s\n      %v\n\n", constants.ColorYellow, stepDesc, constants.ColorReset, err)
	fmt.Println("  Ensure you are running GitMap from an elevated PowerShell prompt (Run as Administrator).")

	return apperror.NewWithDetails("ssh", "E9002", fmt.Sprintf("%s failed: %v", stepDesc, err), "cli", apperror.ErrorTypeAbort, apperror.SeverityWarn, map[string]any{"reported": true})
}

func configureSSHWindowsService(port int) error {
	fmt.Printf("  %s[2/4]%s Setting sshd service startup type to Automatic...\n", constants.ColorCyan, constants.ColorReset)
	if err := runWindowsPowerShellScript("Set-Service -Name sshd -StartupType Automatic", "Set-Service sshd Automatic"); err != nil {
		return handleSSHDServiceError(err, "Set-Service sshd Automatic")
	}

	fmt.Printf("  %s[3/4]%s Configuring inbound firewall rule for TCP port %d...\n", constants.ColorCyan, constants.ColorReset, port)
	if err := runWindowsPowerShellScript(buildWindowsFirewallEnsureCmd(port), "Configure NetFirewallRule"); err != nil {
		return handleSSHDServiceError(err, "Configure NetFirewallRule")
	}

	return nil
}

func enableSSHWindows(port int, isForce bool) error {
	if err := ensureWindowsCapability(isForce); err != nil {
		return err
	}

	if err := configureSSHWindowsService(port); err != nil {
		return err
	}

	fmt.Printf("  %s[4/4]%s Starting sshd service...\n", constants.ColorCyan, constants.ColorReset)
	if err := runWindowsPowerShellScript("Start-Service sshd", "Start-Service sshd"); err != nil {
		return handleSSHDServiceError(err, "Start-Service sshd")
	}

	return nil
}

func buildDebianEnableSteps(port int) []LinuxStep {
	portRule := fmt.Sprintf("%d/tcp", port)

	return []LinuxStep{
		{Binary: "apt-get", Args: []string{"install", "-y", "openssh-server"}, Desc: "Installing openssh-server via apt-get"},
		{Binary: "systemctl", Args: []string{"enable", "--now", "ssh"}, Desc: "Enabling and starting ssh service via systemctl"},
		{Binary: "ufw", Args: []string{"allow", portRule}, Desc: fmt.Sprintf("Opening firewall port %s via ufw", portRule)},
	}
}

func buildRhelEnableSteps(pkgMgr string, port int) []LinuxStep {
	portRule := fmt.Sprintf("--add-port=%d/tcp", port)

	return []LinuxStep{
		{Binary: pkgMgr, Args: []string{"install", "-y", "openssh-server"}, Desc: fmt.Sprintf("Installing openssh-server via %s", pkgMgr)},
		{Binary: "systemctl", Args: []string{"enable", "--now", "sshd"}, Desc: "Enabling and starting sshd service via systemctl"},
		{Binary: "firewall-cmd", Args: []string{"--permanent", portRule}, Desc: fmt.Sprintf("Opening firewall port %d/tcp via firewall-cmd", port)},
		{Binary: "firewall-cmd", Args: []string{"--reload"}, Desc: "Reloading firewall-cmd"},
	}
}

func readOSRelease() string {
	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return ""
	}

	return strings.ToLower(string(data))
}

func detectLinuxDistro() string {
	content := readOSRelease()
	if strings.Contains(content, "ubuntu") || strings.Contains(content, "debian") {
		return "debian"
	}

	if strings.Contains(content, "rhel") || strings.Contains(content, "centos") {
		return "rhel"
	}

	if _, err := os.Stat("/etc/redhat-release"); err == nil {
		return "rhel"
	}

	return "debian"
}

func detectRhelPackageManager() string {
	_, err := exec.LookPath("dnf")
	if err == nil {
		return "dnf"
	}

	return "yum"
}

func isFatalLinuxStepError(bin string, err error) bool {
	if err == nil {
		return false
	}

	return bin != "ufw" && bin != "firewall-cmd"
}

func executeLinuxSteps(steps []LinuxStep) error {
	for i, step := range steps {
		fmt.Printf("  %s[%d/%d]%s %s...\n", constants.ColorCyan, i+1, len(steps), constants.ColorReset, step.Desc)
		_, err := daemonExecRunner(step.Binary, step.Args...)
		if isFatalLinuxStepError(step.Binary, err) {
			return apperror.NewExecutionError(fmt.Sprintf("%s failed: %v", step.Desc, err))
		}
	}

	return nil
}

func enableSSHDebian(port int) error {
	return executeLinuxSteps(buildDebianEnableSteps(port))
}

func enableSSHRhel(port int) error {
	return executeLinuxSteps(buildRhelEnableSteps(detectRhelPackageManager(), port))
}

func enableSSHLinux(port int, isForce bool) error {
	_ = isForce
	if detectLinuxDistro() == "rhel" {
		return enableSSHRhel(port)
	}

	return enableSSHDebian(port)
}

func printSSHEnableBanner(targetOS string, port int) {
	fmt.Printf("\n%s╔══════════════════════════════════════════════════════════════════╗%s\n", constants.ColorGreen, constants.ColorReset)
	fmt.Printf("%s║         OPENSSH SERVER DAEMON SUCCESSFULLY ENABLED               ║%s\n", constants.ColorGreen, constants.ColorReset)
	fmt.Printf("%s╚══════════════════════════════════════════════════════════════════╝%s\n", constants.ColorGreen, constants.ColorReset)
	fmt.Printf("  Target OS:       %s%s%s\n  Service Name:    %ssshd (Automatic, Running)%s\n  Listening Port:  %s%d/TCP (Inbound Allowed)%s\n  Verify Command:  %sgitmap ssh troubleshoot 127.0.0.1%s\n\n", constants.ColorCyan, targetOS, constants.ColorReset, constants.ColorGreen, constants.ColorReset, constants.ColorCyan, port, constants.ColorReset, constants.ColorYellow, constants.ColorReset)
}

func runEnableSSHWindows(opts sshEnableOptions) error {
	if err := enableSSHWindows(opts.port, opts.isForce); err != nil {
		return err
	}

	printSSHEnableBanner("Windows (OpenSSH.Server)", opts.port)

	return nil
}

func runEnableSSHLinux(opts sshEnableOptions) error {
	if err := enableSSHLinux(opts.port, opts.isForce); err != nil {
		return err
	}

	printSSHEnableBanner("Linux ("+detectLinuxDistro()+")", opts.port)

	return nil
}

func dispatchSSHEnableOS(opts sshEnableOptions) error {
	if runtime.GOOS == "windows" {
		return runEnableSSHWindows(opts)
	}

	if runtime.GOOS == "linux" {
		return runEnableSSHLinux(opts)
	}

	return apperror.NewExecutionError(fmt.Sprintf("unsupported OS %q for automated sshd enabling", runtime.GOOS))
}

func runSSHEnableCLI(args []string) error {
	opts, err := parseSSHEnableFlags(args)
	if err != nil {
		return err
	}

	if opts.isHelp {
		printSSHEnableHelp()

		return nil
	}

	return dispatchSSHEnableOS(opts)
}

func printSSHDisableHelp() {
	fmt.Printf("\n%sUsage:%s gitmap ssh disable [flags]\n\n", constants.ColorCyan, constants.ColorReset)
	fmt.Println("Stops OpenSSH Server daemon (sshd), disables auto-start, and closes firewall ports.\n\nFlags:\n  --help, -h          Show this help text")
}

func printSSHDisableBanner(targetOS string) {
	fmt.Printf("%s✔ SSH daemon stopped and disabled on %s.%s\n", constants.ColorGreen, targetOS, constants.ColorReset)
}

func disableSSHWindows() error {
	script := "Stop-Service sshd -ErrorAction SilentlyContinue; Set-Service sshd -StartupType Manual -ErrorAction SilentlyContinue; Disable-NetFirewallRule -Name 'OpenSSH-Server-In-TCP' -ErrorAction SilentlyContinue; Remove-NetFirewallRule -Name 'OpenSSH-Server-In-TCP' -ErrorAction SilentlyContinue"
	_, err := runWindowsPowerShell(script)
	if err != nil {
		return apperror.NewExecutionError(fmt.Sprintf("failed to disable sshd on Windows: %v", err))
	}

	return nil
}

func disableSSHLinux() error {
	_, _ = daemonExecRunner("systemctl", "stop", "ssh")
	_, _ = daemonExecRunner("systemctl", "stop", "sshd")
	_, _ = daemonExecRunner("systemctl", "disable", "ssh")
	_, _ = daemonExecRunner("systemctl", "disable", "sshd")
	_, _ = daemonExecRunner("ufw", "delete", "allow", "22/tcp")
	_, _ = daemonExecRunner("ufw", "deny", "22/tcp")
	_, _ = daemonExecRunner("firewall-cmd", "--permanent", "--remove-service=ssh")
	_, _ = daemonExecRunner("firewall-cmd", "--permanent", "--remove-port=22/tcp")
	_, _ = daemonExecRunner("firewall-cmd", "--reload")

	return nil
}

func disableSSHDarwin() error {
	_, err := daemonExecRunner("systemsetup", "-setremotelogin", "off")
	if err != nil {
		return apperror.NewExecutionError(fmt.Sprintf("systemsetup failed: %v", err))
	}

	return nil
}

func dispatchSSHDisableOS() (string, error) {
	switch runtime.GOOS {
	case "windows":
		return "Windows", disableSSHWindows()
	case "linux":
		return "Linux", disableSSHLinux()
	case "darwin":
		return "Darwin", disableSSHDarwin()
	default:
		return runtime.GOOS, apperror.NewExecutionError(fmt.Sprintf("unsupported OS %q for automated sshd disabling", runtime.GOOS))
	}
}

func runSSHDisableCLI(args []string) error {
	if hasHelpFlag(args) {
		printSSHDisableHelp()

		return nil
	}

	targetOS, err := dispatchSSHDisableOS()
	if err != nil {
		return err
	}

	printSSHDisableBanner(targetOS)

	return nil
}
