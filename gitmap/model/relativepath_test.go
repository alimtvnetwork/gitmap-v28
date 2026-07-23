package model

// Locks the CleanRelativePath contract documented in relativepath.go:
//
//   1. Empty input returns "" (NEVER ".") — empty signals "no manifest
//      row" to downstream callers (reclone_confirm.go / reclone_summary.go).
//   2. Forward slashes become OS-native via filepath.FromSlash.
//   3. Redundant `.`, doubled separators, and trailing separators are
//      collapsed via filepath.Clean.
//   4. Already-clean OS-native input is a no-op (idempotent).
//
// These four properties match the failure modes cataloged in the
// relativepath.go header comment.

import (
	"path/filepath"
	"runtime"
	"testing"
)

func TestCleanRelativePathPreservesEmpty(t *testing.T) {
	if got := CleanRelativePath(""); got != "" {
		t.Fatalf("expected empty string preserved, got %q", got)
	}
}

func TestCleanRelativePathConvertsForwardSlashesToOSNative(t *testing.T) {
	got := CleanRelativePath("acme/widget")
	want := filepath.Join("acme", "widget")

	if got != want {
		t.Fatalf("CleanRelativePath(%q) = %q, want %q", "acme/widget", got, want)
	}
}

func TestCleanRelativePathCollapsesRedundantDot(t *testing.T) {
	got := CleanRelativePath("./acme/widget")
	want := filepath.Join("acme", "widget")

	if got != want {
		t.Fatalf("CleanRelativePath(%q) = %q, want %q", "./acme/widget", got, want)
	}
}

func TestCleanRelativePathCollapsesDoubledSeparators(t *testing.T) {
	got := CleanRelativePath("acme//widget")
	want := filepath.Join("acme", "widget")

	if got != want {
		t.Fatalf("CleanRelativePath(%q) = %q, want %q", "acme//widget", got, want)
	}
}

func TestCleanRelativePathStripsTrailingSeparator(t *testing.T) {
	got := CleanRelativePath("acme/widget/")
	want := filepath.Join("acme", "widget")

	if got != want {
		t.Fatalf("CleanRelativePath(%q) = %q, want %q", "acme/widget/", got, want)
	}
}

func TestCleanRelativePathIsIdempotent(t *testing.T) {
	once := CleanRelativePath("acme/widget/sub")
	twice := CleanRelativePath(once)

	if once != twice {
		t.Fatalf("expected idempotent normalization, %q != %q", once, twice)
	}
}

func TestCleanRelativePathPreservesWindowsBackslashesOnWindows(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("backslash semantics only meaningful on Windows")
	}

	got := CleanRelativePath(`acme\widget`)
	want := `acme\widget`

	if got != want {
		t.Fatalf("CleanRelativePath(%q) = %q, want %q", `acme\widget`, got, want)
	}
}
