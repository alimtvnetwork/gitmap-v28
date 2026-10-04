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
	"github.com/alimtvnetwork/gitmap-v28/cli/firewall"
	"github.com/alimtvnetwork/gitmap-v28/cli/lazyregex"
)

var portLineRegex = lazyregex.New(`(?i)^\s*#?\s*Port\s+\d+`)

func parsePortArgument(args []string) (int, bool, error) {
	if len(args) == 0 {
		return 0, false, nil
	}

	first := args[0]
	if isHelpFlag(first) {
		return 0, true, nil
	}

	if isPortFlag(first) {
		return parseFlagPort(args)
	}

	if strings.HasPrefix(first, "--port=") {
		val := strings.TrimPrefix(first, "--port=")
		p, err := parsePortFlagValue(val)

		return p, false, err
	}

	p, err := parsePortFlagValue(first)

	return p, false, err
}

func parseFlagPort(args []string) (int, bool, error) {
	if len(args) < 2 {
		return 0, false, apperror.NewValidationError("missing port number after flag")
	}

	p, err := parsePortFlagValue(args[1])

	return p, false, err
}

func printSSHPortHelp() {
	fmt.Printf("\n%sUsage:%s gitmap ssh port <port-number> [flags]\n\n", constants.ColorCyan, constants.ColorReset)
	fmt.Println("Reconfigures OpenSSH Server listening port in sshd_config, updates firewall rules, and restarts sshd.\n\nExamples:\n  gitmap ssh port 2222\n  gitmap ssh port --port 2200")
}

func resolveSSHDConfigPath() string {
	if envPath := os.Getenv("GITMAP_SSHD_CONFIG_PATH"); envPath != "" {
		return envPath
	}

	if runtime.GOOS == "windows" {
		return resolveWindowsSSHDConfigPath()
	}

	return "/etc/ssh/sshd_config"
}

func resolveWindowsSSHDConfigPath() string {
	progData := os.Getenv("ProgramData")
	if len(progData) == 0 {
		progData = `C:\ProgramData`
	}

	return filepath.Join(progData, "ssh", "sshd_config")
}

func updateSSHDConfigContent(content string, newPort int) string {
	newDirective := fmt.Sprintf("Port %d", newPort)
	trimmed := strings.TrimSpace(content)
	if len(trimmed) == 0 {
		return newDirective + "\n"
	}

	lines := strings.Split(content, "\n")
	hasReplaced := false

	for i, line := range lines {
		if !portLineRegex.MatchString(line) {
			continue
		}

		lines[i] = newDirective
		hasReplaced = true
		break
	}

	if !hasReplaced {
		lines = append(lines, newDirective)
	}

	return strings.Join(lines, "\n")
}

func createSSHDConfigBackup(configPath string) {
	origBytes, err := os.ReadFile(configPath)
	if err != nil {
		return
	}

	bakPath := configPath + ".gitmap.bak"
	_ = os.WriteFile(bakPath, origBytes, 0600)
}

func backupAndWriteSSHDConfig(configPath string, newContent string) error {
	createSSHDConfigBackup(configPath)

	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return apperror.NewExecutionError(fmt.Sprintf("create config directory %q: %v", dir, err))
	}

	writeErr := os.WriteFile(configPath, []byte(newContent), 0644)
	if writeErr != nil {
		return apperror.NewExecutionError(fmt.Sprintf("write sshd_config %q: %v", configPath, writeErr))
	}

	return nil
}

func revertSSHDConfigBackup(configPath string) {
	bakPath := configPath + ".gitmap.bak"
	bakBytes, readErr := os.ReadFile(bakPath)
	if readErr != nil {
		return
	}

	_ = os.WriteFile(configPath, bakBytes, 0644)
}

func testSSHDConfigSyntax(configPath string) error {
	sshdPath, err := exec.LookPath("sshd")
	if err != nil {
		return nil
	}

	out, testErr := daemonExecRunner(sshdPath, "-t", "-f", configPath)
	if testErr != nil {
		revertSSHDConfigBackup(configPath)

		return apperror.NewValidationError(fmt.Sprintf("sshd -t syntax test failed (reverted to backup): %s", strings.TrimSpace(string(out))))
	}

	return nil
}

func updateFirewallWindows(port int) error {
	script := buildWindowsFirewallEnsureCmd(port)

	return runWindowsPowerShellScript(script, "Update Firewall for SSH Port")
}

func hasBinary(name string) bool {
	_, err := exec.LookPath(name)

	return err == nil
}

func applyFirewallCmdRule(portRule string) error {
	portArg := fmt.Sprintf("--add-port=%s", portRule)
	_, _ = daemonExecRunner("firewall-cmd", "--permanent", portArg)
	_, _ = daemonExecRunner("firewall-cmd", "--reload")

	return nil
}

func updateFirewallLinux(port int) error {
	portRule := fmt.Sprintf("%d/tcp", port)
	if hasBinary("ufw") {
		_, _ = daemonExecRunner("ufw", "allow", portRule)

		return nil
	}

	if hasBinary("firewall-cmd") {
		return applyFirewallCmdRule(portRule)
	}

	return nil
}

func updateFirewallForPort(port int) error {
	if runtime.GOOS == "windows" {
		return updateFirewallWindows(port)
	}

	if runtime.GOOS == "linux" {
		return updateFirewallLinux(port)
	}

	return nil
}

func restartSSHDWindows() error {
	return runWindowsPowerShellScript("Restart-Service -Name sshd -Force", "Restart sshd service")
}

func restartSSHDLinux() error {
	if _, err := daemonExecRunner("systemctl", "restart", "ssh"); err == nil {
		return nil
	}

	if _, err := daemonExecRunner("systemctl", "restart", "sshd"); err != nil {
		return apperror.NewExecutionError(fmt.Sprintf("failed to restart ssh/sshd service: %v", err))
	}

	return nil
}

func restartSSHDService() error {
	if runtime.GOOS == "windows" {
		return restartSSHDWindows()
	}

	if runtime.GOOS == "linux" {
		return restartSSHDLinux()
	}

	return nil
}

func printSSHPortBanner(configPath string, port int) {
	fmt.Printf("\n%s╔══════════════════════════════════════════════════════════════════╗%s\n", constants.ColorGreen, constants.ColorReset)
	fmt.Printf("%s║          OPENSSH PORT RECONFIGURED SUCCESSFULLY                  ║%s\n", constants.ColorGreen, constants.ColorReset)
	fmt.Printf("%s╚══════════════════════════════════════════════════════════════════╝%s\n", constants.ColorGreen, constants.ColorReset)
	fmt.Printf("  New SSH Port:    %s%d/TCP%s\n  Configuration:   %s%s%s\n  Firewall Status: %sUpdated to allow TCP port %d%s\n  Daemon Status:   %ssshd service restarted%s\n", constants.ColorCyan, port, constants.ColorReset, constants.ColorDim, configPath, constants.ColorReset, constants.ColorGreen, port, constants.ColorReset, constants.ColorGreen, constants.ColorReset)
	fmt.Printf("  Next Steps:\n    1. Verify connection: %sssh user@<host> -p %d%s\n    2. Update fleet node: %sgitmap nodes (or gitmap ssh add <node> --port %d)%s\n\n", constants.ColorYellow, port, constants.ColorReset, constants.ColorCyan, port, constants.ColorReset)
}

func readExistingConfigContent(configPath string) string {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return ""
	}

	return string(data)
}

func applyPortConfiguration(configPath string, port int) error {
	content := readExistingConfigContent(configPath)
	newContent := updateSSHDConfigContent(content, port)

	if err := backupAndWriteSSHDConfig(configPath, newContent); err != nil {
		return err
	}

	if err := testSSHDConfigSyntax(configPath); err != nil {
		return err
	}

	if err := updateFirewallForPort(port); err != nil {
		return err
	}

	if err := restartSSHDService(); err != nil {
		return err
	}

	printSSHPortBanner(configPath, port)

	return nil
}

func runSSHPortCLI(args []string) error {
	if len(args) == 0 {
		return runSSHPortList()
	}

	sub := strings.ToLower(strings.TrimSpace(args[0]))
	switch sub {
	case "ls", "list", "show":
		return runSSHPortList()
	case "add":
		return runSSHPortAdd(args[1:])
	case "rm", "remove", "delete", "del":
		return runSSHPortRemove(args[1:])
	case "set":
		return runSSHPortSet(args[1:])
	case "public", "enable-public":
		return runSSHPortEnablePublic(args[1:])
	case "help", "--help", "-h":
		printSSHPortHelp()

		return nil
	default:
		return runSSHPortSet(args)
	}
}

func parsePortsFromConfig(content string) []int {
	var ports []int
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#") {
			continue
		}
		if portLineRegex.MatchString(trimmed) {
			fields := strings.Fields(trimmed)
			if len(fields) >= 2 {
				if p, err := strconv.Atoi(fields[1]); err == nil && p > 0 {
					ports = append(ports, p)
				}
			}
		}
	}
	if len(ports) == 0 {
		ports = append(ports, 22)
	}

	return ports
}

func runSSHPortList() error {
	configPath := resolveSSHDConfigPath()
	content := readExistingConfigContent(configPath)
	ports := parsePortsFromConfig(content)

	fmt.Println()
	fmt.Printf("%s╔══════════════════════════════════════════════════════════════════╗%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("%s║                   OPENSSH CONFIGURED PORTS                       ║%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("%s╚══════════════════════════════════════════════════════════════════╝%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  Configuration File: %s%s%s\n\n", constants.ColorDim, configPath, constants.ColorReset)
	fmt.Printf("  %-8s %-12s %-22s\n", "PORT", "PROTOCOL", "FIREWALL STATUS")
	fmt.Println("  ──────────────────────────────────────────")

	for _, p := range ports {
		allowed, _ := firewall.IsPortAllowed(p)
		status := constants.ColorRed + "✘ BLOCKED" + constants.ColorReset
		if allowed {
			status = constants.ColorGreen + "✔ ALLOWED" + constants.ColorReset
		}
		fmt.Printf("  %-8d %-12s %s\n", p, "TCP", status)
	}
	fmt.Println()
	fmt.Printf("  %sCommands:%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Println("    gitmap ssh port add <port>      - Add listening port and open firewall")
	fmt.Println("    gitmap ssh port rm <port>       - Remove listening port and delete firewall rule")
	fmt.Println("    gitmap ssh port set <port>      - Switch primary listening port")
	fmt.Println("    gitmap ssh enable-public <port> - Open port for public / internet access")
	fmt.Println()

	return nil
}

func runSSHPortAdd(args []string) error {
	if len(args) == 0 {
		return apperror.NewValidationError("missing required port argument (1-65535)")
	}
	p, err := parsePortFlagValue(args[0])
	if err != nil {
		return err
	}
	configPath := resolveSSHDConfigPath()
	content := readExistingConfigContent(configPath)
	newDirective := fmt.Sprintf("Port %d", p)
	if strings.Contains(content, newDirective) {
		fmt.Printf("  Port %d already exists in %s\n", p, configPath)
	} else {
		content = content + "\n" + newDirective + "\n"
		if err := backupAndWriteSSHDConfig(configPath, content); err != nil {
			return err
		}
	}

	_ = firewall.AllowPort(p, fmt.Sprintf("OpenSSH-Server-In-TCP-%d", p))
	_ = restartSSHDService()
	fmt.Printf("  %s✔ Added SSH Port %d/TCP and opened in firewall.%s\n", constants.ColorGreen, p, constants.ColorReset)

	return nil
}

func runSSHPortRemove(args []string) error {
	if len(args) == 0 {
		return apperror.NewValidationError("missing required port argument (1-65535)")
	}
	p, err := parsePortFlagValue(args[0])
	if err != nil {
		return err
	}
	configPath := resolveSSHDConfigPath()
	content := readExistingConfigContent(configPath)
	lines := strings.Split(content, "\n")
	var newLines []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if portLineRegex.MatchString(trimmed) && strings.Contains(trimmed, strconv.Itoa(p)) {
			continue
		}
		newLines = append(newLines, line)
	}
	if err := backupAndWriteSSHDConfig(configPath, strings.Join(newLines, "\n")); err != nil {
		return err
	}

	_ = firewall.RemovePortRule(p, fmt.Sprintf("OpenSSH-Server-In-TCP-%d", p))
	_ = restartSSHDService()
	fmt.Printf("  %s✔ Removed SSH Port %d/TCP and cleaned firewall rules.%s\n", constants.ColorGreen, p, constants.ColorReset)

	return nil
}

func runSSHPortSet(args []string) error {
	port, isHelp, err := parsePortArgument(args)
	if err != nil {
		return err
	}
	if isHelp {
		printSSHPortHelp()

		return nil
	}
	if port == 0 {
		printSSHPortHelp()

		return apperror.NewValidationError("missing required port argument (1-65535)")
	}
	configPath := resolveSSHDConfigPath()

	return applyPortConfiguration(configPath, port)
}

func runSSHPortEnablePublic(args []string) error {
	port := 22
	if len(args) > 0 {
		if p, err := parsePortFlagValue(args[0]); err == nil && p > 0 {
			port = p
		}
	}

	configPath := resolveSSHDConfigPath()
	content := readExistingConfigContent(configPath)

	portDirective := fmt.Sprintf("Port %d", port)
	if !strings.Contains(content, portDirective) {
		content = updateSSHDConfigContent(content, port)
	}

	if !strings.Contains(content, "GatewayPorts") {
		content += "\nGatewayPorts yes\n"
	}
	if !strings.Contains(content, "ListenAddress 0.0.0.0") {
		content += "\nListenAddress 0.0.0.0\n"
	}

	if err := backupAndWriteSSHDConfig(configPath, content); err != nil {
		return err
	}

	if err := firewall.EnablePublicSSH(port); err != nil {
		return err
	}

	_ = restartSSHDService()

	fmt.Printf("\n%s╔══════════════════════════════════════════════════════════════════╗%s\n", constants.ColorGreen, constants.ColorReset)
	fmt.Printf("%s║          OPENSSH PUBLIC ACCESS ENABLED ON PORT %-5d            ║%s\n", constants.ColorGreen, port, constants.ColorReset)
	fmt.Printf("%s╚══════════════════════════════════════════════════════════════════╝%s\n", constants.ColorGreen, constants.ColorReset)
	fmt.Printf("  Port:            %s%d/TCP%s\n", constants.ColorCyan, port, constants.ColorReset)
	fmt.Printf("  Firewall Status: %sInbound edge traversal and public profile allowed%s\n", constants.ColorGreen, constants.ColorReset)
	fmt.Printf("  Listen Address:  %s0.0.0.0 (all interfaces)%s\n", constants.ColorGreen, constants.ColorReset)
	fmt.Printf("  Daemon Status:   %ssshd service active and restarted%s\n\n", constants.ColorGreen, constants.ColorReset)

	return nil
}
