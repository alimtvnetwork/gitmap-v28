package cmdclone

import (
	"testing"
)

func TestRunCloneOnlyMissing_Flags(t *testing.T) {
	cf := parseCloneFlags([]string{"--ssh", "-w", "8"})
	cf.MissingOnly = true
	cf.NoReplace = true

	if !cf.MissingOnly {
		t.Errorf("expected MissingOnly to be true")
	}
	if !cf.NoReplace {
		t.Errorf("expected NoReplace to be true")
	}
	if !cf.UseSSH {
		t.Errorf("expected UseSSH to be true")
	}
	if cf.MaxConcurrency != 8 {
		t.Errorf("expected MaxConcurrency 8, got %d", cf.MaxConcurrency)
	}
}
