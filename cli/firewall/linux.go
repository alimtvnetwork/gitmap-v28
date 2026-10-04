package firewall

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func allowPortLinux(port int, name string) error {
	if _, err := exec.LookPath("ufw"); err == nil {
		cmd := exec.Command("ufw", "allow", fmt.Sprintf("%d/tcp", port))
		out, errCmd := cmd.CombinedOutput()
		if errCmd != nil {
			return fmt.Errorf("ufw allow failed: %s", strings.TrimSpace(string(out)))
		}

		return nil
	}

	if _, err := exec.LookPath("firewall-cmd"); err == nil {
		cmd := exec.Command("firewall-cmd", "--permanent", fmt.Sprintf("--add-port=%d/tcp", port))
		out, errCmd := cmd.CombinedOutput()
		if errCmd != nil {
			return fmt.Errorf("firewall-cmd add-port failed: %s", strings.TrimSpace(string(out)))
		}
		_ = exec.Command("firewall-cmd", "--reload").Run()

		return nil
	}

	if _, err := exec.LookPath("iptables"); err == nil {
		cmd := exec.Command("iptables", "-I", "INPUT", "-p", "tcp", "--dport", strconv.Itoa(port), "-j", "ACCEPT")
		out, errCmd := cmd.CombinedOutput()
		if errCmd != nil {
			return fmt.Errorf("iptables allow failed: %s", strings.TrimSpace(string(out)))
		}

		return nil
	}

	return fmt.Errorf("no supported Linux firewall tool found (ufw, firewalld, iptables)")
}

func blockPortLinux(port int, name string) error {
	if _, err := exec.LookPath("ufw"); err == nil {
		cmd := exec.Command("ufw", "deny", fmt.Sprintf("%d/tcp", port))
		out, errCmd := cmd.CombinedOutput()
		if errCmd != nil {
			return fmt.Errorf("ufw deny failed: %s", strings.TrimSpace(string(out)))
		}

		return nil
	}

	if _, err := exec.LookPath("firewall-cmd"); err == nil {
		cmd := exec.Command("firewall-cmd", "--permanent", fmt.Sprintf("--remove-port=%d/tcp", port))
		out, errCmd := cmd.CombinedOutput()
		if errCmd != nil {
			return fmt.Errorf("firewall-cmd remove-port failed: %s", strings.TrimSpace(string(out)))
		}
		_ = exec.Command("firewall-cmd", "--reload").Run()

		return nil
	}

	if _, err := exec.LookPath("iptables"); err == nil {
		cmd := exec.Command("iptables", "-I", "INPUT", "-p", "tcp", "--dport", strconv.Itoa(port), "-j", "DROP")
		out, errCmd := cmd.CombinedOutput()
		if errCmd != nil {
			return fmt.Errorf("iptables drop failed: %s", strings.TrimSpace(string(out)))
		}

		return nil
	}

	return fmt.Errorf("no supported Linux firewall tool found (ufw, firewalld, iptables)")
}

func removePortRuleLinux(port int, name string) error {
	if _, err := exec.LookPath("ufw"); err == nil {
		cmd := exec.Command("ufw", "delete", "allow", fmt.Sprintf("%d/tcp", port))
		_, _ = cmd.CombinedOutput()

		return nil
	}

	if _, err := exec.LookPath("firewall-cmd"); err == nil {
		cmd := exec.Command("firewall-cmd", "--permanent", fmt.Sprintf("--remove-port=%d/tcp", port))
		_, _ = cmd.CombinedOutput()
		_ = exec.Command("firewall-cmd", "--reload").Run()

		return nil
	}

	if _, err := exec.LookPath("iptables"); err == nil {
		cmd := exec.Command("iptables", "-D", "INPUT", "-p", "tcp", "--dport", strconv.Itoa(port), "-j", "ACCEPT")
		_, _ = cmd.CombinedOutput()

		return nil
	}

	return nil
}

func isPortAllowedLinux(port int) (bool, error) {
	portStr := strconv.Itoa(port)
	if _, err := exec.LookPath("ufw"); err == nil {
		out, errCmd := exec.Command("ufw", "status").CombinedOutput()
		if errCmd == nil {
			text := string(out)
			if strings.Contains(text, portStr) && strings.Contains(text, "ALLOW") {
				return true, nil
			}
		}
	}

	if _, err := exec.LookPath("firewall-cmd"); err == nil {
		out, errCmd := exec.Command("firewall-cmd", "--list-ports").CombinedOutput()
		if errCmd == nil {
			if strings.Contains(string(out), fmt.Sprintf("%d/tcp", port)) {
				return true, nil
			}
		}
	}

	if _, err := exec.LookPath("iptables"); err == nil {
		out, errCmd := exec.Command("iptables", "-L", "INPUT", "-n").CombinedOutput()
		if errCmd == nil {
			if strings.Contains(string(out), fmt.Sprintf("dpt:%d", port)) {
				return true, nil
			}
		}
	}

	return false, nil
}

func enablePublicSSHLinux(port int) error {
	return allowPortLinux(port, fmt.Sprintf("OpenSSH-Public-%d", port))
}

func listRulesLinux() ([]Rule, error) {
	var rules []Rule
	if _, err := exec.LookPath("ufw"); err == nil {
		out, errCmd := exec.Command("ufw", "status").CombinedOutput()
		if errCmd == nil {
			lines := strings.Split(string(out), "\n")
			for _, line := range lines {
				trimmed := strings.TrimSpace(line)
				if strings.Contains(trimmed, "ALLOW") || strings.Contains(trimmed, "DENY") {
					parts := strings.Fields(trimmed)
					if len(parts) >= 2 {
						action := "ALLOW"
						if strings.Contains(parts[1], "DENY") {
							action = "BLOCK"
						}
						rules = append(rules, Rule{
							Name:      parts[0],
							Action:    action,
							Direction: "IN",
							Protocol:  "TCP",
							Enabled:   true,
						})
					}
				}
			}
		}
	}

	return rules, nil
}
