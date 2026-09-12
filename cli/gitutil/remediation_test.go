package gitutil

import (
	"testing"
)

func TestRemediationSuite(t *testing.T) {
	diag := DirtyDiagnosis{
		IsDirty:        true,
		ModifiedCount:  2,
		UntrackedCount: 1,
		SummaryReason:  "+2 modified, +1 untracked",
	}

	recipes := GenerateRemediationRecipes("/path/to/repo", diag)
	if len(recipes) != 3 {
		t.Fatalf("expected 3 remediation recipes, got %d", len(recipes))
	}

	if !HasPotentialConflictRisk(diag) {
		t.Fatal("expected conflict risk for modified files")
	}

	badge := FormatDirtyBadges(diag)
	if badge == "clean" {
		t.Fatal("expected non-clean badge")
	}
}
