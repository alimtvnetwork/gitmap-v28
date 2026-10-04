package firewall

import (
	"fmt"
	"os/exec"
)

func allowPortDarwin(port int, name string) error {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("echo 'pass in proto tcp to any port %d' | pfctl -ef -", port))
	_, err := cmd.CombinedOutput()

	return err
}

func blockPortDarwin(port int, name string) error {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("echo 'block in proto tcp to any port %d' | pfctl -ef -", port))
	_, err := cmd.CombinedOutput()

	return err
}

func removePortRuleDarwin(port int, name string) error {
	return nil
}

func isPortAllowedDarwin(port int) (bool, error) {
	return true, nil
}

func enablePublicSSHDarwin(port int) error {
	return allowPortDarwin(port, fmt.Sprintf("OpenSSH-Public-%d", port))
}

func listRulesDarwin() ([]Rule, error) {
	return nil, nil
}
