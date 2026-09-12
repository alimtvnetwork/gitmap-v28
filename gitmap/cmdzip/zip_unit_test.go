package cmdzip

import (
	"testing"
)

func TestReorderFlagsBeforeArgs(t *testing.T) {
	args := []string{"target.zip", "-o", "outdir", "file1.txt"}
	reordered := reorderFlagsBeforeArgs(args)

	if len(reordered) != 4 {
		t.Fatalf("expected 4 tokens, got %d: %v", len(reordered), reordered)
	}

	if reordered[0] != "-o" || reordered[1] != "outdir" {
		t.Errorf("expected flags moved to front, got %v", reordered)
	}

	if reordered[2] != "target.zip" || reordered[3] != "file1.txt" {
		t.Errorf("expected positionals preserved at end, got %v", reordered)
	}
}

func TestZipGroupHints(t *testing.T) {
	createHints := zipGroupCreateHints()
	if len(createHints) == 0 {
		t.Errorf("expected non-empty create hints")
	}

	listHints := zipGroupListHints()
	if len(listHints) == 0 {
		t.Errorf("expected non-empty list hints")
	}

	showHints := zipGroupShowHints()
	if len(showHints) == 0 {
		t.Errorf("expected non-empty show hints")
	}
}
