package scanner

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// makeWorktreeRepo simulates a `git worktree add` linked checkout: a
// regular `.git` FILE whose contents start with `gitdir: …` rather than
// a `.git` directory. Returns the absolute repo path.
func makeWorktreeRepo(t *testing.T, root, rel, gitdirTarget string) string {
	t.Helper()
	full := filepath.Join(root, rel)
	if err := os.MkdirAll(full, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", rel, err)
	}
	contents := []byte("gitdir: " + gitdirTarget + "\n")
	if err := os.WriteFile(filepath.Join(full, constants.ExtGit), contents, 0o644); err != nil {
		t.Fatalf("write .git file %s: %v", rel, err)
	}

	return full
}

// TestScanDirDetectsWorktreeGitFile verifies that a `.git` regular file
// with `gitdir:` prefix (the layout produced by `git worktree add` and by
// absorbed submodules) is recognized as a repo root just like a `.git`
// directory.
func TestScanDirDetectsWorktreeGitFile(t *testing.T) {
	root := t.TempDir()
	makeRepo(t, root, "main-checkout")
	makeWorktreeRepo(t, root, "worktree-feature", "/var/repos/main/.git/worktrees/feature")

	got, err := ScanDir(root, nil)
	if err != nil {
		t.Fatalf("ScanDir: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("want 2 repos (dir + worktree), got %d: %+v", len(got), got)
	}
}

// TestScanDirIgnoresNonGitdirFile asserts that a stray `.git` regular
// file WITHOUT the `gitdir:` prefix is NOT treated as a repo. This
// guards against false positives from misconfigured editors or backup
// tools dropping a `.git` text file into a normal folder.
func TestScanDirIgnoresNonGitdirFile(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "not-a-repo")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, constants.ExtGit),
		[]byte("just some text, no prefix\n"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	got, err := ScanDir(root, nil)
	if err != nil {
		t.Fatalf("ScanDir: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("want 0 repos (stray .git file should not count), got %+v", got)
	}
}

// TestScanDirWorktreeStopsDescent confirms that once the `.git` FILE
// marker is recognized, the walker treats the subtree as opaque just
// like for `.git` directories — nested checkouts beneath are ignored.
func TestScanDirWorktreeStopsDescent(t *testing.T) {
	root := t.TempDir()
	makeWorktreeRepo(t, root, "wt", "/var/repos/main/.git/worktrees/x")
	// Nested standard repo under the worktree — must NOT be discovered.
	if err := os.MkdirAll(filepath.Join(root, "wt", "nested", constants.ExtGit), 0o755); err != nil {
		t.Fatalf("nested: %v", err)
	}

	got, err := ScanDir(root, nil)
	if err != nil {
		t.Fatalf("ScanDir: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("want 1 repo (worktree only), got %d: %+v", len(got), got)
	}
}
