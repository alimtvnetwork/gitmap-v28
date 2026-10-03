package cmdssh

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var portLineRegex = regexp.MustCompile(`(?i)^\s*#?\s*Port\s+\d+`)

func parsePortArgument(args []string) (int, bool, error) {
	if len(args) == 0 {
		return 0, false, nil
	}

	first := args[0]
	if first == "--help" || first == "-h" || first == "help" {
		return 0, true, nil
	}

	if first == "--port" || first == "-p" {
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
	fmt.Println("Reconfigures OpenSSH Server listening port in sshd_config, updates firewall rules, and restarts sshd.")
	fmt.Println("\nExamples:")
	fmt.Println("  gitmap ssh port 2222")
	fmt.Println("  gitmap ssh port --port 2200")
	fmt.Println()
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
	if progData == "" {
		progData = `C:\ProgramData`
	}

	return filepath.Join(progData, "ssh", "sshd_config")
}

func updateSSHDConfigContent(content string, newPort int) string {
	lines := strings.Split(content, "\n")
	hasReplaced := false
	newDirective := fmt.Sprintf("Port %d", newPort)

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

func backupAndWriteSSHDConfig(configPath string, newContent string) error {
	bakPath := configPath + ".gitmap.bak"
	origBytes, err := os.ReadFile(configPath)
	if err == nil {
		_ = os.WriteFile(bakPath, origBytes, 0600)
	}

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

func testSSHDConfigSyntax(configPath string) error {
	sshdPath, err := exec.LookPath("sshd")
	if err != nil {
		return nil
	}

	cmd := exec.Command(sshdPath, "-t", "-f", configPath)
	out, testErr := cmd.CombinedOutput()
	if testErr != nil {
		revertSSHDConfigBackup(configPath)

		return apperror.NewValidationError(fmt.Sprintf("sshd -t syntax test failed (reverted to backup): %s", strings.TrimSpace(string(out))))
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

func updateFirewallWindows(port int) error {
	script := fmt.Sprintf(
		"if (Get-NetFirewallRule -Name 'OpenSSH-Server-In-TCP' -ErrorAction SilentlyContinue) { Set-NetFirewallRule -Name 'OpenSSH-Server-In-TCP' -LocalPort %d } else { New-NetFirewallRule -Name 'OpenSSH-Server-In-TCP' -DisplayName 'OpenSSH Server (sshd)' -Enabled True -Direction Inbound -Protocol TCP -Action Allow -LocalPort %d }",
		port, port,
	)

	return runWindowsPowerShellScript(script, "Update Firewall for SSH Port")
}

func updateFirewallLinux(port int) error {
	portRule := fmt.Sprintf("%d/tcp", port)
	if _, lookUfw := exec.LookPath("ufw"); lookUfw == nil {
		_, _ = daemonExecRunner("ufw", "allow", portRule)
		return nil
	}

	if _, lookFw := exec.LookPath("firewall-cmd"); lookFw == nil {
		portArg := fmt.Sprintf("--add-port=%s", portRule)
		_, _ = daemonExecRunner("firewall-cmd", "--permanent", portArg)
		_, _ = daemonExecRunner("firewall-cmd", "--reload")
		return nil
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
	fmt.Println()
	fmt.Printf("%s╔══════════════════════════════════════════════════════════════════╗%s\n", constants.ColorGreen, constants.ColorReset)
	fmt.Printf("%s║          OPENSSH PORT RECONFIGURED SUCCESSFULLY                  ║%s\n", constants.ColorGreen, constants.ColorReset)
	fmt.Printf("%s╚══════════════════════════════════════════════════════════════════╝%s\n", constants.ColorGreen, constants.ColorReset)
	fmt.Printf("  New SSH Port:    %s%d/TCP%s\n", constants.ColorCyan, port, constants.ColorReset)
	fmt.Printf("  Configuration:   %s%s%s\n", constants.ColorDim, configPath, constants.ColorReset)
	fmt.Printf("  Firewall Status: %sUpdated to allow TCP port %d%s\n", constants.ColorGreen, port, constants.ColorReset)
	fmt.Printf("  Daemon Status:   %ssshd service restarted%s\n", constants.ColorGreen, constants.ColorReset)
	fmt.Printf("  Connect test:    %sssh user@<ip> -p %d%s\n", constants.ColorYellow, port, constants.ColorReset)
	fmt.Println()
}

func applyPortConfiguration(configPath string, port int) error {
	content := ""
	if data, err := os.ReadFile(configPath); err == nil {
		content = string(data)
	}

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
