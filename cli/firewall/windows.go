package firewall

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func runNetsh(args ...string) (string, error) {
	cmd := exec.Command("netsh", args...)
	out, err := cmd.CombinedOutput()

	return string(out), err
}

func allowPortWindows(port int, name string) error {
	_ = removePortRuleWindows(port, name)
	args := []string{
		"advfirewall", "firewall", "add", "rule",
		fmt.Sprintf("name=%s", name),
		"dir=in",
		"action=allow",
		"protocol=TCP",
		fmt.Sprintf("localport=%d", port),
		"profile=any",
	}
	out, err := runNetsh(args...)
	if err != nil && !strings.Contains(out, "Ok.") {
		return fmt.Errorf("netsh allow port %d failed: %s", port, strings.TrimSpace(out))
	}

	return nil
}

func blockPortWindows(port int, name string) error {
	_ = removePortRuleWindows(port, name)
	args := []string{
		"advfirewall", "firewall", "add", "rule",
		fmt.Sprintf("name=%s", name),
		"dir=in",
		"action=block",
		"protocol=TCP",
		fmt.Sprintf("localport=%d", port),
		"profile=any",
	}
	out, err := runNetsh(args...)
	if err != nil && !strings.Contains(out, "Ok.") {
		return fmt.Errorf("netsh block port %d failed: %s", port, strings.TrimSpace(out))
	}

	return nil
}

func removePortRuleWindows(port int, name string) error {
	args := []string{"advfirewall", "firewall", "delete", "rule", fmt.Sprintf("name=%s", name)}
	_, err := runNetsh(args...)

	return err
}

func isPortAllowedWindows(port int) (bool, error) {
	out, err := runNetsh("advfirewall", "firewall", "show", "rule", "name=all", "dir=in")
	if err != nil {
		return false, err
	}

	targetPortStr := strconv.Itoa(port)
	lines := strings.Split(out, "\n")
	var currentName string
	var currentAction string
	var currentPort string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "Rule Name:") {
			currentName = strings.TrimSpace(strings.TrimPrefix(trimmed, "Rule Name:"))
			currentAction = ""
			currentPort = ""
		} else if strings.HasPrefix(trimmed, "Action:") {
			currentAction = strings.TrimSpace(strings.TrimPrefix(trimmed, "Action:"))
		} else if strings.HasPrefix(trimmed, "LocalPort:") {
			currentPort = strings.TrimSpace(strings.TrimPrefix(trimmed, "LocalPort:"))
			if currentPort == targetPortStr && strings.EqualFold(currentAction, "Allow") {
				if strings.Contains(strings.ToLower(currentName), "ssh") {
					return true, nil
				}
			}
		}
	}

	return false, nil
}

func enablePublicSSHWindows(port int) error {
	ruleName := fmt.Sprintf("OpenSSH-Server-Public-TCP-%d", port)
	_ = removePortRuleWindows(port, ruleName)
	args := []string{
		"advfirewall", "firewall", "add", "rule",
		fmt.Sprintf("name=%s", ruleName),
		"dir=in",
		"action=allow",
		"protocol=TCP",
		fmt.Sprintf("localport=%d", port),
		"profile=any",
		"edge=yes",
	}
	out, err := runNetsh(args...)
	if err != nil && !strings.Contains(out, "Ok.") {
		return fmt.Errorf("failed to enable public SSH firewall rule on port %d: %s", port, strings.TrimSpace(out))
	}

	return nil
}

func listRulesWindows() ([]Rule, error) {
	out, err := runNetsh("advfirewall", "firewall", "show", "rule", "name=all", "dir=in")
	if err != nil {
		return nil, err
	}

	var rules []Rule
	lines := strings.Split(out, "\n")
	var current Rule

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "Rule Name:") {
			if current.Name != "" && current.Port > 0 {
				rules = append(rules, current)
			}
			current = Rule{
				Name:      strings.TrimSpace(strings.TrimPrefix(trimmed, "Rule Name:")),
				Direction: "IN",
				Protocol:  "TCP",
				Enabled:   true,
			}
		} else if strings.HasPrefix(trimmed, "Action:") {
			current.Action = strings.ToUpper(strings.TrimSpace(strings.TrimPrefix(trimmed, "Action:")))
		} else if strings.HasPrefix(trimmed, "Enabled:") {
			current.Enabled = strings.EqualFold(strings.TrimSpace(strings.TrimPrefix(trimmed, "Enabled:")), "Yes")
		} else if strings.HasPrefix(trimmed, "LocalPort:") {
			pStr := strings.TrimSpace(strings.TrimPrefix(trimmed, "LocalPort:"))
			if p, errConv := strconv.Atoi(pStr); errConv == nil {
				current.Port = p
			}
		}
	}

	if current.Name != "" && current.Port > 0 {
		rules = append(rules, current)
	}

	return rules, nil
}
