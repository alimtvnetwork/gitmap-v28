package firewall

import (
	"fmt"
	"runtime"
)

// Rule represents an OS firewall rule entry.
type Rule struct {
	Name      string `json:"name"`
	Port      int    `json:"port"`
	Protocol  string `json:"protocol"`
	Direction string `json:"direction"`
	Action    string `json:"action"`
	Enabled   bool   `json:"enabled"`
}

// AllowPort configures the operating system firewall to allow inbound TCP traffic on port.
func AllowPort(port int, name string) error {
	if name == "" {
		name = fmt.Sprintf("OpenSSH-Server-In-TCP-%d", port)
	}

	switch runtime.GOOS {
	case "windows":
		return allowPortWindows(port, name)
	case "linux":
		return allowPortLinux(port, name)
	case "darwin":
		return allowPortDarwin(port, name)
	default:
		return fmt.Errorf("unsupported operating system for firewall automation: %s", runtime.GOOS)
	}
}

// BlockPort configures the operating system firewall to deny or drop inbound TCP traffic on port.
func BlockPort(port int, name string) error {
	if name == "" {
		name = fmt.Sprintf("OpenSSH-Server-Block-TCP-%d", port)
	}

	switch runtime.GOOS {
	case "windows":
		return blockPortWindows(port, name)
	case "linux":
		return blockPortLinux(port, name)
	case "darwin":
		return blockPortDarwin(port, name)
	default:
		return fmt.Errorf("unsupported operating system for firewall automation: %s", runtime.GOOS)
	}
}

// RemovePortRule removes a firewall rule for the specified port.
func RemovePortRule(port int, name string) error {
	if name == "" {
		name = fmt.Sprintf("OpenSSH-Server-In-TCP-%d", port)
	}

	switch runtime.GOOS {
	case "windows":
		return removePortRuleWindows(port, name)
	case "linux":
		return removePortRuleLinux(port, name)
	case "darwin":
		return removePortRuleDarwin(port, name)
	default:
		return fmt.Errorf("unsupported operating system for firewall automation: %s", runtime.GOOS)
	}
}

// IsPortAllowed checks whether inbound traffic on port is permitted by the OS firewall.
func IsPortAllowed(port int) (bool, error) {
	switch runtime.GOOS {
	case "windows":
		return isPortAllowedWindows(port)
	case "linux":
		return isPortAllowedLinux(port)
	case "darwin":
		return isPortAllowedDarwin(port)
	default:
		return false, nil
	}
}

// EnablePublicSSH ensures the firewall permits external public connections to the SSH port.
func EnablePublicSSH(port int) error {
	switch runtime.GOOS {
	case "windows":
		return enablePublicSSHWindows(port)
	case "linux":
		return enablePublicSSHLinux(port)
	case "darwin":
		return enablePublicSSHDarwin(port)
	default:
		return fmt.Errorf("unsupported operating system for public SSH enablement: %s", runtime.GOOS)
	}
}

// ListRules returns detected firewall rules for SSH / remote administration.
func ListRules() ([]Rule, error) {
	switch runtime.GOOS {
	case "windows":
		return listRulesWindows()
	case "linux":
		return listRulesLinux()
	case "darwin":
		return listRulesDarwin()
	default:
		return nil, nil
	}
}
