package cmd

// Locks the canonicalizePMPath contract documented in clonepmsync.go:
//
//   1. filepath.Clean is ALWAYS applied (collapses mixed `/` and `\`
//      separators, removes redundant `.` and trailing slashes).
//   2. filepath.EvalSymlinks resolves symlinks AND Windows 8.3 short
//      names when the path exists on disk.
//   3. Soft-fail: a non-existent path, broken symlink, or permission
//      error returns the cleaned absolute form rather than erroring.
//      A projects.json entry is always preferable to a swallowed
//      clone — see clonepmsync.go header for the rationale.
//
// The fix that this test guards against regression for landed alongside
// the Windows path-handling audit of the clone -> projects.json sync
// (see chat for context). Before the fix, mixed-separator inputs and
// 8.3-short-name ancestors produced duplicate sidebar entries because
// `filepath.Abs` alone leaves both forms distinct.

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestCanonicalizePMPathCleansMixedSeparators(t *testing.T) {
	dir := t.TempDir()

	// Build a path with redundant `.` segments and a trailing slash.
	// On every OS, filepath.Clean must collapse these. The leaf dir
	// exists so EvalSymlinks succeeds and we observe the cleaned form.
	leaf := filepath.Join(dir, "repo")
	if err := os.MkdirAll(leaf, 0o755); err != nil {
		t.Fatalf("mkdir leaf: %v", err)
	}

	dirty := dir + string(filepath.Separator) + "." +
		string(filepath.Separator) + "repo" + string(filepath.Separator)

	got := canonicalizePMPath(dirty)
	want := leaf

	// EvalSymlinks may resolve the tempdir's parent (e.g. /var ->
	// /private/var on macOS). Compare the basename + that the result
	// is clean (no redundant separators) rather than exact match.
	if filepath.Clean(got) != got {
		t.Fatalf("canonicalizePMPath(%q) = %q (not clean)", dirty, got)
	}

	if filepath.Base(got) != filepath.Base(want) {
		t.Fatalf("canonicalizePMPath(%q) basename = %q, want %q",
			dirty, filepath.Base(got), filepath.Base(want))
	}
}

func TestCanonicalizePMPathCollapsesMixedSlashes(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("mixed `/` + `\\` separators only meaningful on Windows")
	}

	dir := t.TempDir()
	leaf := filepath.Join(dir, "repo")
	if err := os.MkdirAll(leaf, 0o755); err != nil {
		t.Fatalf("mkdir leaf: %v", err)
	}

	// e.g. `D:\Users\me\AppData\Local\Temp\xyz/repo` — common when a
	// JSON manifest's forward-slash RelativePath is joined onto a
	// Windows abs target.
	mixed := dir + "/repo"

	got := canonicalizePMPath(mixed)

	if strings.Contains(got, "/") {
		t.Fatalf("canonicalizePMPath(%q) = %q still contains `/`",
			mixed, got)
	}
}

func TestCanonicalizePMPathSoftFailsOnMissingPath(t *testing.T) {
	// Build a path that definitely does NOT exist. EvalSymlinks will
	// error; canonicalizePMPath must fall back to the cleaned form
	// rather than returning empty / blocking the caller.
	missing := filepath.Join(t.TempDir(), "no", "such", "dir")

	got := canonicalizePMPath(missing)
	want := filepath.Clean(missing)

	if got != want {
		t.Fatalf("canonicalizePMPath(%q) = %q, want %q (cleaned fallback)",
			missing, got, want)
	}

	if got == "" {
		t.Fatalf("canonicalizePMPath returned empty string — would " +
			"swallow projects.json entry on first-clone of new folder")
	}
}

func TestCanonicalizePMPathResolvesSymlink(t *testing.T) {
	dir := t.TempDir()

	target := filepath.Join(dir, "real-repo")
	if err := os.MkdirAll(target, 0o755); err != nil {
		t.Fatalf("mkdir target: %v", err)
	}

	link := filepath.Join(dir, "link-repo")
	if err := os.Symlink(target, link); err != nil {
		// Windows requires Developer Mode or admin to create
		// symlinks; treat as environment-not-available, NOT a
		// regression.
		t.Skipf("symlink unavailable in this environment: %v", err)
	}

	got := canonicalizePMPath(link)

	// EvalSymlinks may rewrite the tempdir ancestor too (macOS
	// /var -> /private/var); basename is the stable invariant —
	// the LINK's basename must NOT survive, the TARGET's must.
	if filepath.Base(got) != filepath.Base(target) {
		t.Fatalf("canonicalizePMPath(%q) basename = %q, want %q "+
			"(symlink not resolved)",
			link, filepath.Base(got), filepath.Base(target))
	}
}

func TestCanonicalizePMPathSoftFailsOnBrokenSymlink(t *testing.T) {
	dir := t.TempDir()

	link := filepath.Join(dir, "broken-link")
	dangling := filepath.Join(dir, "no-such-target")
	if err := os.Symlink(dangling, link); err != nil {
		t.Skipf("symlink unavailable in this environment: %v", err)
	}

	got := canonicalizePMPath(link)

	// EvalSymlinks errors on a broken link; soft-fall to the cleaned
	// link path itself (NOT the dangling target).
	want := filepath.Clean(link)
	if got != want {
		t.Fatalf("canonicalizePMPath(broken symlink %q) = %q, want %q",
			link, got, want)
	}
}
