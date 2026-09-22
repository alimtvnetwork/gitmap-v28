package cmdos

import (
	"testing"
)

func TestHasDryRunFlag(t *testing.T) {
	if !hasDryRunFlag([]string{"--dry-run"}) {
		t.Errorf("expected true for --dry-run")
	}

	if !hasDryRunFlag([]string{"-n"}) {
		t.Errorf("expected true for -n")
	}

	if hasDryRunFlag([]string{"--other"}) {
		t.Errorf("expected false for --other")
	}
}

func TestRunOSUpdateHelpRouting(t *testing.T) {
	if err := runOSUpdateCommand(false, []string{"--help"}); err != nil {
		t.Errorf("expected nil error for update help, got: %v", err)
	}

	if err := runOSUpdateCommand(true, []string{"-h"}); err != nil {
		t.Errorf("expected nil error for upgrade help, got: %v", err)
	}
}
