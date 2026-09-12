package clonefrom

// Tests for the dest-parent / hierarchy-preservation behavior added
// in execute_dest.go. Split out of execute_test.go so neither file
// breaches the project's 200-line cap. Reuses the requireGit /
// makeBareRepo helpers defined in execute_test.go (same package).

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// TestExecute_MkdirParentFailureIsFailedRow covers the negative
// path of prepareDestParent: when MkdirAll cannot create the parent
// (here: parent path collides with an existing FILE), the row must
// be reported as `failed` with a non-empty Detail — NOT crash and
// NOT silently swallow. Locks in the Code Red zero-swallow promise.
func TestExecute_MkdirParentFailureIsFailedRow(t *testing.T) {
	cwd := t.TempDir()
	// Plant a regular FILE where the dest's parent dir would go.
	// MkdirAll on a path whose ancestor is a file returns ENOTDIR.
	blocker := filepath.Join(cwd, "blocker")
	if err := os.WriteFile(blocker, []byte("x"), 0o644); err != nil {
		t.Fatalf("seed blocker: %v", err)
	}

	plan := Plan{Rows: []Row{{
		URL:  "file:///does/not/matter.git",
		Dest: filepath.Join("blocker", "child", "repo"),
	}}}

	results := Execute(plan, cwd, io.Discard)

	if results[0].Status != constants.CloneFromStatusFailed {
		t.Fatalf("status = %q, want failed", results[0].Status)
	}

	if !strings.Contains(results[0].Detail, "mkdir parent") {
		t.Errorf("detail = %q, want mkdir-parent diagnosis", results[0].Detail)
	}
}
