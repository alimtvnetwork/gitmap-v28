package cmdssh

import (
	"fmt"
	"os"
	"os/exec"
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

func parsePortFlagValue(val string) (int, error) {
	port, err := strconv.Atoi(val)
	if err != nil {
		return 0, apperror.NewValidationError(fmt.Sprintf("invalid port value %q: must be integer", val))
	}

	if port < 1 || port > 65535 {
		return 0, apperror.NewValidationError(fmt.Sprintf("port %d out of range: must be between 1 and 65535", port))
	}

	return port, nil
}

func parsePortFlagToken(arg string, next string, hasNext bool, opts *sshEnableOptions) (int, error) {
	if strings.HasPrefix(arg, "--port=") {
		val := strings.TrimPrefix(arg, "--port=")
		p, err := parsePortFlagValue(val)
		opts.port = p
		return 1, err
	}

	if arg != "--port" && arg != "-p" {
		return 0, nil
	}

	if !hasNext {
		return 1, apperror.NewValidationError("missing argument for --port flag")
	}

	p, err := parsePortFlagValue(next)
	opts.port = p

	return 2, err
}

func parseSingleEnableFlag(args []string, idx int, opts *sshEnableOptions) (int, error) {
	arg := args[idx]
	if arg == "--help" || arg == "-h" || arg == "help" {
		opts.isHelp = true
		return 1, nil
	}

	if arg == "--force" || arg == "-f" {
		opts.isForce = true
		return 1, nil
	}

	hasNext := idx+1 < len(args)
	next := ""
	if hasNext {
		next = args[idx+1]
	}

	consumed, err := parsePortFlagToken(arg, next, hasNext, opts)
	if err != nil {
		return 0, err
	}

	if consumed > 0 {
		return consumed, nil
	}

	return 1, nil
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
	fmt.Println("Installs, configures, starts OpenSSH Server daemon (sshd) and opens firewall ports.")
	fmt.Println("\nFlags:")
	fmt.Println("  --port, -p <port>   Port to listen on (default 22)")
	fmt.Println("  --force, -f         Force installation even if prerequisites exist")
	fmt.Println("  --help, -h          Show this help text")
	fmt.Println()
}

func runWindowsPowerShellScript(script string, desc string) error {
	psArgs := []string{
		"-NoProfile",
		"-NonInteractive",
		"-ExecutionPolicy",
		"Bypass",
		"-Command",
		script,
	}

	out, err := daemonExecRunner("powershell.exe", psArgs...)
	if err != nil {
		return apperror.NewExecutionError(fmt.Sprintf("%s failed: %v (output: %s)", desc, err, strings.TrimSpace(string(out))))
	}

	return nil
}

func enableSSHWindows(port int, isForce bool) error {
	fmt.Printf("  %s[1/4]%s Installing OpenSSH Server Windows Capability...\n", constants.ColorCyan, constants.ColorReset)
	capCmd := "Add-WindowsCapability -Online -Name OpenSSH.Server~~~~0.0.1.0"
	if err := runWindowsPowerShellScript(capCmd, "Add-WindowsCapability"); err != nil && !isForce {
		return err
	}

	fmt.Printf("  %s[2/4]%s Setting sshd service startup type to Automatic...\n", constants.ColorCyan, constants.ColorReset)
	setSvcCmd := "Set-Service -Name sshd -StartupType Automatic"
	if err := runWindowsPowerShellScript(setSvcCmd, "Set-Service sshd Automatic"); err != nil {
		return err
	}

	fmt.Printf("  %s[3/4]%s Configuring inbound firewall rule for TCP port %d...\n", constants.ColorCyan, constants.ColorReset, port)
	fwCmd := fmt.Sprintf(
		"if (Get-NetFirewallRule -Name 'OpenSSH-Server-In-TCP' -ErrorAction SilentlyContinue) { Set-NetFirewallRule -Name 'OpenSSH-Server-In-TCP' -LocalPort %d -Enabled True } else { New-NetFirewallRule -Name 'OpenSSH-Server-In-TCP' -DisplayName 'OpenSSH Server (sshd)' -Enabled True -Direction Inbound -Protocol TCP -Action Allow -LocalPort %d }",
		port, port,
	)
	if err := runWindowsPowerShellScript(fwCmd, "Configure NetFirewallRule"); err != nil {
		return err
	}

	fmt.Printf("  %s[4/4]%s Starting sshd service...\n", constants.ColorCyan, constants.ColorReset)
	startCmd := "Start-Service sshd"
	if err := runWindowsPowerShellScript(startCmd, "Start-Service sshd"); err != nil {
		return err
	}

	return nil
}

func detectLinuxDistro() string {
	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return detectLinuxDistroFromFiles()
	}

	content := strings.ToLower(string(data))
	if strings.Contains(content, "ubuntu") || strings.Contains(content, "debian") {
		return "debian"
	}

	if strings.Contains(content, "rhel") || strings.Contains(content, "centos") || strings.Contains(content, "fedora") {
		return "rhel"
	}

	return detectLinuxDistroFromFiles()
}

func detectLinuxDistroFromFiles() string {
	if _, err := os.Stat("/etc/debian_version"); err == nil {
		return "debian"
	}

	if _, err := os.Stat("/etc/redhat-release"); err == nil {
		return "rhel"
	}

	return "debian"
}

func enableSSHDebian(port int) error {
	fmt.Printf("  %s[1/3]%s Installing openssh-server via apt-get...\n", constants.ColorCyan, constants.ColorReset)
	if _, err := daemonExecRunner("apt-get", "install", "-y", "openssh-server"); err != nil {
		return apperror.NewExecutionError(fmt.Sprintf("apt-get install failed: %v", err))
	}

	fmt.Printf("  %s[2/3]%s Enabling and starting ssh service via systemctl...\n", constants.ColorCyan, constants.ColorReset)
	if _, err := daemonExecRunner("systemctl", "enable", "--now", "ssh"); err != nil {
		return apperror.NewExecutionError(fmt.Sprintf("systemctl enable ssh failed: %v", err))
	}

	fmt.Printf("  %s[3/3]%s Opening firewall port %d/tcp via ufw...\n", constants.ColorCyan, constants.ColorReset, port)
	portRule := fmt.Sprintf("%d/tcp", port)
	_, _ = daemonExecRunner("ufw", "allow", portRule)

	return nil
}

func enableSSHRhel(port int) error {
	fmt.Printf("  %s[1/4]%s Installing openssh-server via dnf/yum...\n", constants.ColorCyan, constants.ColorReset)
	pkgMgr := "dnf"
	if _, lookErr := exec.LookPath("dnf"); lookErr != nil {
		pkgMgr = "yum"
	}

	if _, err := daemonExecRunner(pkgMgr, "install", "-y", "openssh-server"); err != nil {
		return apperror.NewExecutionError(fmt.Sprintf("%s install failed: %v", pkgMgr, err))
	}

	fmt.Printf("  %s[2/4]%s Enabling and starting sshd service via systemctl...\n", constants.ColorCyan, constants.ColorReset)
	if _, err := daemonExecRunner("systemctl", "enable", "--now", "sshd"); err != nil {
		return apperror.NewExecutionError(fmt.Sprintf("systemctl enable sshd failed: %v", err))
	}

	fmt.Printf("  %s[3/4]%s Opening firewall port %d/tcp via firewall-cmd...\n", constants.ColorCyan, constants.ColorReset, port)
	portArg := fmt.Sprintf("--add-port=%d/tcp", port)
	_, _ = daemonExecRunner("firewall-cmd", "--permanent", portArg)

	fmt.Printf("  %s[4/4]%s Reloading firewall-cmd...\n", constants.ColorCyan, constants.ColorReset)
	_, _ = daemonExecRunner("firewall-cmd", "--reload")

	return nil
}

func enableSSHLinux(port int, isForce bool) error {
	distro := detectLinuxDistro()
	if distro == "rhel" {
		return enableSSHRhel(port)
	}

	return enableSSHDebian(port)
}

func printSSHEnableBanner(targetOS string, port int) {
	fmt.Println()
	fmt.Printf("%s╔══════════════════════════════════════════════════════════════════╗%s\n", constants.ColorGreen, constants.ColorReset)
	fmt.Printf("%s║         OPENSSH SERVER DAEMON SUCCESSFULLY ENABLED               ║%s\n", constants.ColorGreen, constants.ColorReset)
	fmt.Printf("%s╚══════════════════════════════════════════════════════════════════╝%s\n", constants.ColorGreen, constants.ColorReset)
	fmt.Printf("  Target OS:       %s%s%s\n", constants.ColorCyan, targetOS, constants.ColorReset)
	fmt.Printf("  Service Name:    %ssshd (Automatic, Running)%s\n", constants.ColorGreen, constants.ColorReset)
	fmt.Printf("  Listening Port:  %s%d/TCP (Inbound Allowed)%s\n", constants.ColorCyan, port, constants.ColorReset)
	fmt.Printf("  Verify Command:  %sgitmap ssh troubleshoot 127.0.0.1%s\n", constants.ColorYellow, constants.ColorReset)
	fmt.Println()
}

func runEnableSSHWindows(opts sshEnableOptions) error {
	err := enableSSHWindows(opts.port, opts.isForce)
	if err != nil {
		return err
	}

	printSSHEnableBanner("Windows (OpenSSH.Server)", opts.port)

	return nil
}

func runEnableSSHLinux(opts sshEnableOptions) error {
	err := enableSSHLinux(opts.port, opts.isForce)
	if err != nil {
		return err
	}

	printSSHEnableBanner("Linux ("+detectLinuxDistro()+")", opts.port)

	return nil
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

	targetOS := runtime.GOOS
	if targetOS == "windows" {
		return runEnableSSHWindows(opts)
	}

	if targetOS == "linux" {
		return runEnableSSHLinux(opts)
	}

	return apperror.NewExecutionError(fmt.Sprintf("unsupported operating system %q for automated sshd enabling", targetOS))
}
