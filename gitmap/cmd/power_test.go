package cmd

import (
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/power"
)

func TestParseSetArgs(t *testing.T) {
	disp, sleep, err := parseSetArgs([]string{"15"})
	if err != nil {
		t.Fatalf("parseSetArgs positional failed: %v", err)
	}
	if disp != 15 || sleep != 15 {
		t.Errorf("Expected 15, 15, got %d, %d", disp, sleep)
	}

	disp, sleep, err = parseSetArgs([]string{"--display", "20", "--sleep", "45"})
	if err != nil {
		t.Fatalf("parseSetArgs flags failed: %v", err)
	}
	if disp != 20 || sleep != 45 {
		t.Errorf("Expected 20, 45, got %d, %d", disp, sleep)
	}

	_, _, err = parseSetArgs([]string{})
	if err == nil {
		t.Errorf("Expected error for empty args")
	}
}

func TestFormatTimeoutMinutes(t *testing.T) {
	if s := formatTimeoutMinutes(0); s != "Never" {
		t.Errorf("Expected Never, got %s", s)
	}

	if s := formatTimeoutMinutes(30); !strings.Contains(s, "30 minutes") {
		t.Errorf("Expected 30 minutes, got %s", s)
	}
}

func TestRunPowerSubcommands(t *testing.T) {
	restore := power.SetRunnerForTesting(func(name string, args ...string) ([]byte, error) {
		return []byte("OK"), nil
	})
	defer restore()

	if err := runPower([]string{"status"}); err != nil {
		t.Errorf("runPower status failed: %v", err)
	}

	if err := runPower([]string{"never-sleep"}); err != nil {
		t.Errorf("runPower never-sleep failed: %v", err)
	}

	if err := runPower([]string{"set", "0"}); err != nil {
		t.Errorf("runPower set 0 failed: %v", err)
	}

	if err := runPower([]string{"reset"}); err != nil {
		t.Errorf("runPower reset failed: %v", err)
	}
}
