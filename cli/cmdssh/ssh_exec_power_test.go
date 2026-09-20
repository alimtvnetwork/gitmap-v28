package cmdssh

import (
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

func TestIsPowerCommand(t *testing.T) {
	cases := []struct {
		first    string
		expected bool
	}{
		{"shutdown", true},
		{"reboot", true},
		{"poweroff", true},
		{"halt", true},
		{"restart", true},
		{"status", false},
		{"uname", false},
		{"", false},
	}

	for _, c := range cases {
		if got := isPowerCommand(c.first); got != c.expected {
			t.Errorf("isPowerCommand(%q) = %v, want %v", c.first, got, c.expected)
		}
	}
}

func TestIsRebootAction(t *testing.T) {
	cases := []struct {
		first    string
		args     []string
		expected bool
	}{
		{"reboot", []string{"reboot"}, true},
		{"restart", []string{"restart"}, true},
		{"shutdown", []string{"shutdown", "-r", "now"}, true},
		{"shutdown", []string{"shutdown", "/r"}, true},
		{"shutdown", []string{"shutdown", "now"}, false},
		{"shutdown", []string{"shutdown", "-h", "now"}, false},
	}

	for _, c := range cases {
		if got := isRebootAction(c.first, c.args); got != c.expected {
			t.Errorf("isRebootAction(%q, %v) = %v, want %v", c.first, c.args, got, c.expected)
		}
	}
}

func TestResolvePowerCommand_Windows(t *testing.T) {
	c := db.SSHConnection{OS: "windows"}
	shell, cmd := resolvePowerCommand(c, []string{"shutdown", "now"})
	if shell != "ps" || cmd != "shutdown.exe /s /t 0 /f" {
		t.Errorf("expected Windows shutdown command, got (%q, %q)", shell, cmd)
	}

	shellR, cmdR := resolvePowerCommand(c, []string{"reboot"})
	if shellR != "ps" || cmdR != "shutdown.exe /r /t 0 /f" {
		t.Errorf("expected Windows reboot command, got (%q, %q)", shellR, cmdR)
	}
}

func TestResolvePowerCommand_Unix(t *testing.T) {
	c := db.SSHConnection{OS: "linux"}
	shell, cmd := resolvePowerCommand(c, []string{"shutdown", "now"})
	if shell != "bash" || !strings.Contains(cmd, "sudo -n poweroff") {
		t.Errorf("expected Unix shutdown with sudo -n, got (%q, %q)", shell, cmd)
	}

	shellR, cmdR := resolvePowerCommand(c, []string{"reboot"})
	if shellR != "bash" || !strings.Contains(cmdR, "sudo -n reboot") {
		t.Errorf("expected Unix reboot with sudo -n, got (%q, %q)", shellR, cmdR)
	}
}
