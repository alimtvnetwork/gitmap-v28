package cmdos

import (
	"testing"
)

func TestParseDevCleanOptions(t *testing.T) {
	args := []string{"--dry-run", "-y", "--json", "--verbose", "--only", "go,npm"}
	opts := parseDevCleanOptions(args)

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
	if len(opts.OnlyCategories) != 2 || opts.OnlyCategories[0] != "go" || opts.OnlyCategories[1] != "npm" {
		t.Errorf("expected OnlyCategories [go, npm], got %v", opts.OnlyCategories)
	}
}

func TestIsHelpDevClean(t *testing.T) {
	if !isHelpDevClean([]string{"--help"}) {
		t.Errorf("expected --help to be recognized")
	}
	if !isHelpDevClean([]string{"-h"}) {
		t.Errorf("expected -h to be recognized")
	}
	if !isHelpDevClean([]string{"help"}) {
		t.Errorf("expected help to be recognized")
	}
	if isHelpDevClean([]string{"--dry-run"}) {
		t.Errorf("expected --dry-run not to be recognized as help")
	}
}
