package gitutil

import (
	"testing"
)

func TestParseDirtyLineClassifications(t *testing.T) {
	cases := []struct {
		line          string
		wantUntracked int
		wantModified  int
		wantDeleted   int
		wantStaged    int
		wantPrefix    string
	}{
		{"?? new_file.txt", 1, 0, 0, 0, "untracked: "},
		{" M modified.txt", 0, 1, 0, 0, "modified: "},
		{" D deleted.txt", 0, 0, 1, 0, "deleted: "},
		{"A  staged.txt", 0, 0, 0, 1, "staged: "},
	}

	for _, tc := range cases {
		diag := DirtyDiagnosis{IsDirty: true}
		parseDirtyLine(tc.line, &diag)
		if diag.UntrackedCount != tc.wantUntracked {
			t.Errorf("line %q UntrackedCount = %d, want %d", tc.line, diag.UntrackedCount, tc.wantUntracked)
		}
		if diag.ModifiedCount != tc.wantModified {
			t.Errorf("line %q ModifiedCount = %d, want %d", tc.line, diag.ModifiedCount, tc.wantModified)
		}
		if diag.DeletedCount != tc.wantDeleted {
			t.Errorf("line %q DeletedCount = %d, want %d", tc.line, diag.DeletedCount, tc.wantDeleted)
		}
		if diag.StagedCount != tc.wantStaged {
			t.Errorf("line %q StagedCount = %d, want %d", tc.line, diag.StagedCount, tc.wantStaged)
		}
		if len(diag.AllFiles) != 1 || diag.AllFiles[0] != tc.wantPrefix+tc.line[3:] {
			t.Errorf("line %q AllFiles = %v, want prefix %q", tc.line, diag.AllFiles, tc.wantPrefix)
		}
	}
}

func TestBuildSummaryReason(t *testing.T) {
	diag := DirtyDiagnosis{
		ModifiedCount:  2,
		UntrackedCount: 1,
		DeletedCount:   1,
	}
	reason := buildSummaryReason(&diag)
	expected := "+2 modified, +1 untracked, -1 deleted"
	if reason != expected {
		t.Errorf("buildSummaryReason = %q, want %q", reason, expected)
	}

	emptyDiag := DirtyDiagnosis{}
	emptyReason := buildSummaryReason(&emptyDiag)
	if emptyReason != "uncommitted changes" {
		t.Errorf("buildSummaryReason(empty) = %q, want 'uncommitted changes'", emptyReason)
	}
}
