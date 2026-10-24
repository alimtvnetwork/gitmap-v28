package constants

import (
	"testing"
)

func TestConstantsValues(t *testing.T) {
	if FlagParallel != "parallel" {
		t.Errorf("expected parallel, got %s", FlagParallel)
	}

	if TableColHeaderRepo != "REPO" {
		t.Errorf("expected REPO, got %s", TableColHeaderRepo)
	}

	if GlyphTreeBranch != "├── " {
		t.Errorf("expected ├── , got %s", GlyphTreeBranch)
	}
}
