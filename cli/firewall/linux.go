package firewall

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func hasCommand(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func runUfwAllow(port int) error {
	cmd := exec.Command("ufw", "allow", fmt.Sprintf("%d/tcp", port))
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ufw allow failed: %s", strings.TrimSpace(string(out)))
	}
	return nil
}

func runFirewallCmdAllow(port int) error {
	cmd := exec.Command("firewall-cmd", "--permanent", fmt.Sprintf("--add-port=%d/tcp", port))
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("firewall-cmd add-port failed: %s", strings.TrimSpace(string(out)))
	}
	_ = exec.Command("firewall-cmd", "--reload").Run()
	return nil
}

func runIptablesAllow(port int) error {
	cmd := exec.Command("iptables", "-I", "INPUT", "-p", "tcp", "--dport", strconv.Itoa(port), "-j", "ACCEPT")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("iptables allow failed: %s", strings.TrimSpace(string(out)))
	}
	return nil
}

func allowPortLinux(port int, name string) error {
	if hasCommand("ufw") {
		return runUfwAllow(port)
	}
	if hasCommand("firewall-cmd") {
		return runFirewallCmdAllow(port)
	}
	if hasCommand("iptables") {
		return runIptablesAllow(port)
	}
	return fmt.Errorf("no supported Linux firewall tool found (ufw, firewalld, iptables)")
}

func runUfwDeny(port int) error {
	cmd := exec.Command("ufw", "deny", fmt.Sprintf("%d/tcp", port))
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ufw deny failed: %s", strings.TrimSpace(string(out)))
	}
	return nil
}

func runFirewallCmdDeny(port int) error {
	cmd := exec.Command("firewall-cmd", "--permanent", fmt.Sprintf("--remove-port=%d/tcp", port))
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("firewall-cmd remove-port failed: %s", strings.TrimSpace(string(out)))
	}
	_ = exec.Command("firewall-cmd", "--reload").Run()
	return nil
}

func runIptablesDeny(port int) error {
	cmd := exec.Command("iptables", "-I", "INPUT", "-p", "tcp", "--dport", strconv.Itoa(port), "-j", "DROP")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("iptables drop failed: %s", strings.TrimSpace(string(out)))
	}
	return nil
}

func blockPortLinux(port int, name string) error {
	if hasCommand("ufw") {
		return runUfwDeny(port)
	}
	if hasCommand("firewall-cmd") {
		return runFirewallCmdDeny(port)
	}
	if hasCommand("iptables") {
		return runIptablesDeny(port)
	}
	return fmt.Errorf("no supported Linux firewall tool found (ufw, firewalld, iptables)")
}

func removePortRuleLinux(port int, name string) error {
	if hasCommand("ufw") {
		cmd := exec.Command("ufw", "delete", "allow", fmt.Sprintf("%d/tcp", port))
		_, _ = cmd.CombinedOutput()
		return nil
	}
	if hasCommand("firewall-cmd") {
		cmd := exec.Command("firewall-cmd", "--permanent", fmt.Sprintf("--remove-port=%d/tcp", port))
		_, _ = cmd.CombinedOutput()
		_ = exec.Command("firewall-cmd", "--reload").Run()
		return nil
	}
	if hasCommand("iptables") {
		cmd := exec.Command("iptables", "-D", "INPUT", "-p", "tcp", "--dport", strconv.Itoa(port), "-j", "ACCEPT")
		_, _ = cmd.CombinedOutput()
		return nil
	}
	return nil
}

func checkUfwPortAllowed(port int) (bool, error) {
	out, err := exec.Command("ufw", "status").CombinedOutput()
	if err != nil {
		return false, nil
	}
	text := string(out)
	isAllowed := strings.Contains(text, strconv.Itoa(port)) && strings.Contains(text, "ALLOW")
	return isAllowed, nil
}

func checkFirewallCmdPortAllowed(port int) (bool, error) {
	out, err := exec.Command("firewall-cmd", "--list-ports").CombinedOutput()
	if err != nil {
		return false, nil
	}
	isAllowed := strings.Contains(string(out), fmt.Sprintf("%d/tcp", port))
	return isAllowed, nil
}

func checkIptablesPortAllowed(port int) (bool, error) {
	out, err := exec.Command("iptables", "-L", "INPUT", "-n").CombinedOutput()
	if err != nil {
		return false, nil
	}
	isAllowed := strings.Contains(string(out), fmt.Sprintf("dpt:%d", port))
	return isAllowed, nil
}

func isPortAllowedLinux(port int) (bool, error) {
	if hasCommand("ufw") {
		return checkUfwPortAllowed(port)
	}
	if hasCommand("firewall-cmd") {
		return checkFirewallCmdPortAllowed(port)
	}
	if hasCommand("iptables") {
		return checkIptablesPortAllowed(port)
	}
	return false, nil
}

func enablePublicSSHLinux(port int) error {
	return allowPortLinux(port, fmt.Sprintf("OpenSSH-Public-%d", port))
}

func parseUfwRuleLine(trimmed string) (Rule, bool) {
	if !strings.Contains(trimmed, "ALLOW") && !strings.Contains(trimmed, "DENY") {
		return Rule{}, false
	}
	parts := strings.Fields(trimmed)
	if len(parts) < 2 {
		return Rule{}, false
	}
	action := "ALLOW"
	if strings.Contains(parts[1], "DENY") {
		action = "BLOCK"
	}
	return Rule{
		Name:      parts[0],
		Action:    action,
		Direction: "IN",
		Protocol:  "TCP",
		Enabled:   true,
	}, true
}

func listRulesLinux() ([]Rule, error) {
	var rules []Rule
	if !hasCommand("ufw") {
		return rules, nil
	}
	out, err := exec.Command("ufw", "status").CombinedOutput()
	if err != nil {
		return rules, nil
	}
	for _, line := range strings.Split(string(out), "\n") {
		rule, ok := parseUfwRuleLine(strings.TrimSpace(line))
		if ok {
			rules = append(rules, rule)
		}
	}
	return rules, nil
}
