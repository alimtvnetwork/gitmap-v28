package cmdos

import (
	"testing"
)

func TestParseTerminalCleanOptions(t *testing.T) {
	args := []string{"--dry-run", "-y", "--json", "--verbose", "--no-reseed", "--only", "powershell,bash"}
	opts := parseTerminalCleanOptions(args)

	if !opts.IsDryRun {
		t.Errorf("expected IsDryRun true")
	}
	if !opts.HasAutoYes {
		t.Errorf("expected HasAutoYes true")
	}
	if !opts.IsJSON {
		t.Errorf("expected IsJSON true")
	}
	if !opts.IsVerbose {
		t.Errorf("expected IsVerbose true")
	}
	if opts.ShouldReseed {
		t.Errorf("expected ShouldReseed false with --no-reseed")
	}
	if len(opts.OnlyShells) != 2 || opts.OnlyShells[0] != "powershell" || opts.OnlyShells[1] != "bash" {
		t.Errorf("expected OnlyShells [powershell, bash], got %v", opts.OnlyShells)
	}
}

func TestIsHelpTerminalClean(t *testing.T) {
	if !isHelpTerminalClean([]string{"--help"}) {
		t.Errorf("expected --help to be recognized")
	}
	if !isHelpTerminalClean([]string{"-h"}) {
		t.Errorf("expected -h to be recognized")
	}
	if !isHelpTerminalClean([]string{"help"}) {
		t.Errorf("expected help to be recognized")
	}
	if isHelpTerminalClean([]string{"--dry-run"}) {
		t.Errorf("expected --dry-run not to be recognized as help")
	}
}

func TestIsTerminalTarget(t *testing.T) {
	validTargets := []string{
		"terminal", "term", "console", "history", "shell-history",
	}
	for _, target := range validTargets {
		if !isTerminalTarget(target) {
			t.Errorf("expected isTerminalTarget(%q) to be true", target)
		}
	}
	if isTerminalTarget("dev") {
		t.Errorf("expected dev to be false")
	}
}
