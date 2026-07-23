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

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// TestExecute_CreatesMissingParentDirs covers the "preserve folder
// hierarchy" promise: a row whose dest is a NESTED path
// (`org-a/repo-1`) must clone successfully even when the parent
// directory `org-a/` does not exist in cwd. Without prepareDestParent,
// `git clone` fatals with "could not create work tree dir".
func TestExecute_CreatesMissingParentDirs(t *testing.T) {
	requireGit(t)
	bare := makeBareRepo(t)
	cwd := t.TempDir()

	// Two-level nesting — also exercises MkdirAll's recursion.
	nested := filepath.Join("org-a", "team-x", "repo-1")
	plan := Plan{Rows: []Row{{URL: "file://" + bare, Dest: nested}}}

	results := Execute(plan, cwd, io.Discard)

	if results[0].Status != constants.CloneFromStatusOK {
		t.Fatalf("status = %q, detail = %q (want ok)",
			results[0].Status, results[0].Detail)
	}
	gitDir := filepath.Join(cwd, nested, ".git")
	if _, err := os.Stat(gitDir); err != nil {
		t.Errorf("expected .git dir at %s: %v", gitDir, err)
	}
}

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
