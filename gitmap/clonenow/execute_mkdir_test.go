package clonenow

// Verifies runGitClone pre-creates the destination's parent
// directory before invoking git, so nested RelativePath values
// (e.g. `org-a/team/repo-x`) work on a fresh checkout where the
// intermediate folders don't yet exist.
//
// We cannot assert "git clone succeeded" without network + a real
// remote, but we CAN assert the parent dir exists after the call:
// MkdirAll runs unconditionally before the git invocation, so a
// later git failure (no network, bad URL) doesn't undo it.

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunGitClone_CreatesNestedParents(t *testing.T) {
	tmp := t.TempDir()
	row := Row{
		HTTPSUrl:     "https://invalid.invalid/owner/repo.git",
		RelativePath: filepath.Join("org-a", "team-b", "repo-x"),
	}
	// Intentionally let git fail (invalid host) -- we only care that
	// the parent was created BEFORE git ran.
	_, _ = runGitClone(row, row.HTTPSUrl, row.RelativePath, tmp)

	parent := filepath.Join(tmp, "org-a", "team-b")
	info, err := os.Stat(parent)
	if err != nil {
		t.Fatalf("expected parent %q to exist, got: %v", parent, err)
	}
	if !info.IsDir() {
		t.Fatalf("expected %q to be a directory", parent)
	}
}

func TestRunGitClone_PreExistingParentIsNoOp(t *testing.T) {
	tmp := t.TempDir()
	parent := filepath.Join(tmp, "already", "here")
	if err := os.MkdirAll(parent, 0o755); err != nil {
		t.Fatalf("seed parent: %v", err)
	}
	row := Row{
		HTTPSUrl:     "https://invalid.invalid/owner/repo.git",
		RelativePath: filepath.Join("already", "here", "repo-x"),
	}
	// MkdirAll on an existing dir must not error -- the call should
	// proceed to git (which then fails on the invalid host, fine).
	detail, ok := runGitClone(row, row.HTTPSUrl, row.RelativePath, tmp)
	if ok {
		t.Fatalf("expected git to fail on invalid host, got ok=true")
	}
	if strings.HasPrefix(detail, "mkdir parent:") {
		t.Fatalf("unexpected mkdir failure on existing parent: %q", detail)
	}
}
