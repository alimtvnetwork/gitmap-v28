package cmdos

import (
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/power"
)

func TestOSDisplayHelp(t *testing.T) {
	if err := runOS([]string{"display", "help"}); err != nil {
		t.Fatalf("runOS display help failed: %v", err)
	}

	if err := runOS([]string{"display", "-h"}); err != nil {
		t.Fatalf("runOS display -h failed: %v", err)
	}
}

func TestOSDisplayStatusAndAliases(t *testing.T) {
	restore := power.SetRunnerForTesting(func(name string, args ...string) ([]byte, error) {
		return []byte("OK"), nil
	})
	defer restore()

	if err := runOS([]string{"display"}); err != nil {
		t.Errorf("runOS display failed: %v", err)
	}

	if err := runOS([]string{"display", "status"}); err != nil {
		t.Errorf("runOS display status failed: %v", err)
	}

	if err := runOS([]string{"disp"}); err != nil {
		t.Errorf("runOS disp failed: %v", err)
	}

	if err := runOS([]string{"screen"}); err != nil {
		t.Errorf("runOS screen failed: %v", err)
	}
}

func TestOSDisplayOperations(t *testing.T) {
	restore := power.SetRunnerForTesting(func(name string, args ...string) ([]byte, error) {
		return []byte("OK"), nil
	})
	defer restore()

	if err := runOS([]string{"display", "never-sleep"}); err != nil {
		t.Errorf("runOS display never-sleep failed: %v", err)
	}

	if err := runOS([]string{"display", "set", "0"}); err != nil {
		t.Errorf("runOS display set 0 failed: %v", err)
	}

	if err := runOS([]string{"display", "15"}); err != nil {
		t.Errorf("runOS display 15 failed: %v", err)
	}

	if err := runOS([]string{"display", "reset"}); err != nil {
		t.Errorf("runOS display reset failed: %v", err)
	}
}

func TestOSDisplayUnknownSubcommand(t *testing.T) {
	err := runOS([]string{"display", "unknownxyz"})
	if err == nil {
		t.Fatal("expected error for unknown os display subcommand, got nil")
	}

	if !strings.Contains(err.Error(), "unknown os display subcommand") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestFormatDisplayTimeout(t *testing.T) {
	if s := formatDisplayTimeout(0, false); s != "Never (inhibited)" {
		t.Errorf("expected Never (inhibited), got %s", s)
	}

	if s := formatDisplayTimeout(15, false); s != "15 minutes" {
		t.Errorf("expected 15 minutes, got %s", s)
	}

	if s := formatDisplayTimeout(20, true); s != "Never (inhibited)" {
		t.Errorf("expected Never (inhibited), got %s", s)
	}
}

func TestDetectDisplayEnvironment(t *testing.T) {
	srv := detectDisplayServer()
	if srv == "" {
		t.Error("detectDisplayServer returned empty string")
	}

	sess := detectDesktopSession()
	if sess == "" {
		t.Error("detectDesktopSession returned empty string")
	}
}
