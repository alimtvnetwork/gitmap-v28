package cmd

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestIsPathInside(t *testing.T) {
	cases := []struct {
		name   string
		child  string
		parent string
		want   bool
	}{
		{"equal", "/a/b", "/a/b", true},
		{"descendant", "/a/b/c", "/a/b", true},
		{"sibling", "/a/c", "/a/b", false},
		{"parent-of", "/a", "/a/b", false},
		{"case-fold", "/A/B", "/a/b", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := isPathInside(filepath.Clean(tc.child), filepath.Clean(tc.parent))
			if got != tc.want {
				t.Fatalf("isPathInside(%q,%q)=%v want %v",
					tc.child, tc.parent, got, tc.want)
			}
		})
	}
}

func TestEscapeCwdIfInside_NotInside(t *testing.T) {
	target := t.TempDir()
	other := cleanExistingPath(t.TempDir())
	restoreCwd(t)

	if err := os.Chdir(other); err != nil {
		t.Fatalf("chdir other: %v", err)
	}

	got, err := escapeCwdIfInside(target)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !sameCleanPath(got, other) {
		t.Fatalf("cwd should be unchanged; got %q want %q", got, other)
	}
}

// restoreCwd snapshots the current working directory and restores it
// during t.Cleanup. Critical on Windows: if a test chdir's into a
// t.TempDir() and never restores, the process CWD becomes invalid
// when that temp dir is removed during cleanup, and subsequent tests
// that walk up from CWD (schema/golden lookups) fail with "walking up
// from C:\". The cleanup also runs BEFORE the t.TempDir RemoveAll
// (LIFO order: Setenv/Cleanup registered after TempDir runs first),
// so the directory is no longer in-use when Windows tries to remove it.
func restoreCwd(t *testing.T) {
	t.Helper()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatalf("snapshot cwd: %v", err)
	}
	t.Cleanup(func() {
		if cerr := os.Chdir(orig); cerr != nil {
			t.Logf("restore cwd to %q: %v", orig, cerr)
		}
	})
}

func TestEscapeCwdIfInside_EscapesWhenInside(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("temp-dir symlink resolution differs on Windows CI; behavior covered by integration tests")
	}

	target, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatalf("evalsymlinks: %v", err)
	}
	restoreCwd(t)

	if err := os.Chdir(target); err != nil {
		t.Fatalf("chdir target: %v", err)
	}

	got, err := escapeCwdIfInside(target)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	wantParent := filepath.Dir(target)
	if !sameCleanPath(got, wantParent) {
		t.Fatalf("escape landed in %q want parent %q", got, wantParent)
	}

	cwd, _ := os.Getwd()
	resolved, _ := filepath.EvalSymlinks(cwd)
	if !sameCleanPath(resolved, wantParent) {
		t.Fatalf("os cwd %q (resolved %q) != parent %q", cwd, resolved, wantParent)
	}
}

func sameCleanPath(left, right string) bool {
	return strings.EqualFold(cleanExistingPath(left), cleanExistingPath(right))
}
