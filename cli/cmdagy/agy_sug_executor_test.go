package cmdagy

import (
	"testing"
	"time"
)

func TestTriggerSystemShutdown_DryRunSafety(t *testing.T) {
	origExecutor := OSShutdownExecutorFn
	defer func() { OSShutdownExecutorFn = origExecutor }()

	executorCalled := false
	OSShutdownExecutorFn = func(goos string) error {
		executorCalled = true
		return nil
	}

	err := triggerSystemShutdown(2, true)
	if err != nil {
		t.Fatalf("expected nil error under dry-run, got %v", err)
	}
	if executorCalled {
		t.Fatalf("CRITICAL SAFETY VIOLATION: OSShutdownExecutorFn was called under dry-run mode!")
	}
}

func TestGetShutdownCommandStr(t *testing.T) {
	winCmd := getShutdownCommandStr("windows")
	if winCmd != "shutdown /s /t 60" {
		t.Errorf("expected Windows shutdown command 'shutdown /s /t 60', got %q", winCmd)
	}

	macCmd := getShutdownCommandStr("darwin")
	if macCmd != `osascript -e 'tell app "System Events" to shut down'` {
		t.Errorf("expected Darwin shutdown command, got %q", macCmd)
	}

	linuxCmd := getShutdownCommandStr("linux")
	if linuxCmd != "shutdown -h +1" {
		t.Errorf("expected Linux shutdown command 'shutdown -h +1', got %q", linuxCmd)
	}
}

func TestParseSUGInterval(t *testing.T) {
	d := parseSUGInterval([]string{"-t", "10m"}, false)
	if d != 10*time.Minute {
		t.Errorf("expected 10m, got %v", d)
	}

	// Normal run enforces 2m minimum
	short := parseSUGInterval([]string{"-t", "30s"}, false)
	if short != 2*time.Minute {
		t.Errorf("expected 2m minimum for normal run, got %v", short)
	}

	// Dry-run allows short 10s intervals
	dryShort := parseSUGInterval([]string{"-t", "30s"}, true)
	if dryShort != 30*time.Second {
		t.Errorf("expected 30s for dry-run, got %v", dryShort)
	}
}

func TestSUGFlagsParsing(t *testing.T) {
	if !isSUGDryRun([]string{"--dry-run"}) {
		t.Errorf("expected --dry-run to be detected")
	}
	if !isSUGDryRun([]string{"-n"}) {
		t.Errorf("expected -n to be detected")
	}
	if isSUGDryRun([]string{"run"}) {
		t.Errorf("did not expect dry-run without flag")
	}

	if !hasSUGOnceFlag([]string{"--once"}) {
		t.Errorf("expected --once to be detected")
	}
	if !hasSUGOnceFlag([]string{"-1"}) {
		t.Errorf("expected -1 to be detected")
	}
	if hasSUGOnceFlag([]string{"run"}) {
		t.Errorf("did not expect once without flag")
	}
}
