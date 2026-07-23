package scanner

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// makeRepo creates a fake repo by mkdir'ing path/.git so the scanner's
// repo-detection rule fires. Returns the absolute repo path.
func makeRepo(t *testing.T, root, rel string) string {
	t.Helper()
	full := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Join(full, constants.ExtGit), 0o755); err != nil {
		t.Fatalf("makeRepo %s: %v", rel, err)
	}

	return full
}

// TestScanDirFindsAllRepos verifies the parallel walker discovers every
// .git-bearing directory regardless of nesting depth.
func TestScanDirFindsAllRepos(t *testing.T) {
	root := t.TempDir()
	want := []string{
		"a",
		"b",
		"deep/nested/c",
		"side/d",
		"side/sub/e",
	}
	for _, r := range want {
		makeRepo(t, root, r)
	}

	got, err := ScanDir(root, nil)
	if err != nil {
		t.Fatalf("ScanDir: %v", err)
	}
	if len(got) != len(want) {
		t.Fatalf("repo count: got %d (%v), want %d", len(got), got, len(want))
	}

	gotRel := make([]string, len(got))
	for i, r := range got {
		gotRel[i] = filepath.ToSlash(r.RelativePath)
	}
	sort.Strings(gotRel)
	sort.Strings(want)
	for i := range want {
		if gotRel[i] != want[i] {
			t.Errorf("repo[%d]: got %q want %q", i, gotRel[i], want[i])
		}
	}
}

// TestScanDirRespectsExcludes confirms excluded dir names are not
// descended into and any repos beneath them are invisible.
func TestScanDirRespectsExcludes(t *testing.T) {
	root := t.TempDir()
	makeRepo(t, root, "keep")
	makeRepo(t, root, "node_modules/skip")
	makeRepo(t, root, "vendor/skip")

	got, err := ScanDir(root, []string{"node_modules", "vendor"})
	if err != nil {
		t.Fatalf("ScanDir: %v", err)
	}
	if len(got) != 1 || filepath.ToSlash(got[0].RelativePath) != "keep" {
		t.Fatalf("expected only 'keep', got %+v", got)
	}
}

// TestScanDirDoesNotDescendIntoRepos asserts that once a .git is found
// the subtree is treated opaque — nested repos under it are NOT picked
// up. Mirrors the spec: "Do not descend further into a discovered repo."
func TestScanDirDoesNotDescendIntoRepos(t *testing.T) {
	root := t.TempDir()
	makeRepo(t, root, "outer")
	// A second .git nested under outer/ — should be ignored.
	if err := os.MkdirAll(filepath.Join(root, "outer", "submodule", constants.ExtGit), 0o755); err != nil {
		t.Fatalf("nested repo: %v", err)
	}

	got, err := ScanDir(root, nil)
	if err != nil {
		t.Fatalf("ScanDir: %v", err)
	}
	if len(got) != 1 || filepath.ToSlash(got[0].RelativePath) != "outer" {
		t.Fatalf("expected only outer, got %+v", got)
	}
}

// TestScanDirEmpty verifies an empty tree returns no repos and no error.
func TestScanDirEmpty(t *testing.T) {
	root := t.TempDir()
	got, err := ScanDir(root, nil)
	if err != nil {
		t.Fatalf("ScanDir: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected 0 repos, got %d", len(got))
	}
}

// TestScanDirManyReposParallel stress-tests the worker pool with enough
// repos to span multiple workers. Run with -race in CI.
func TestScanDirManyReposParallel(t *testing.T) {
	root := t.TempDir()
	const n = 50
	for i := 0; i < n; i++ {
		makeRepo(t, root, filepath.Join("group", filepath.FromSlash(string(rune('a'+i%5))), "repo", filepath.FromSlash(string(rune('0'+i%10))+"-"+string(rune('a'+i%26)))))
	}

	got, err := ScanDir(root, nil)
	if err != nil {
		t.Fatalf("ScanDir: %v", err)
	}
	// Some path collisions are expected when i%5/i%10/i%26 coincide;
	// just assert the walker produced a non-trivial, unique result set.
	if len(got) == 0 {
		t.Fatalf("expected some repos, got 0")
	}
	seen := make(map[string]bool, len(got))
	for _, r := range got {
		if seen[r.AbsolutePath] {
			t.Errorf("duplicate repo in result: %s", r.AbsolutePath)
		}
		seen[r.AbsolutePath] = true
	}
}
