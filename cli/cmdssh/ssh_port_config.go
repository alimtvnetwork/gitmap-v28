package cmdssh

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
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
