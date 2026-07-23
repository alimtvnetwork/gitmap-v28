// Package cmd — chromeprofile_copy_test.go: edge-case coverage for
// the Chrome profile copy helpers used by `gitmap cpc`.
package cmd

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func TestCopyEntryMissingSourceIsSilent(t *testing.T) {
	dst := filepath.Join(t.TempDir(), "out")
	n, err := copyEntry(filepath.Join(t.TempDir(), "nope"), dst)
	if err != nil || n != 0 {
		t.Fatalf("missing src: want (0,nil), got (%d,%v)", n, err)
	}
	if _, statErr := os.Stat(dst); !os.IsNotExist(statErr) {
		t.Fatalf("dst should not be created for missing src, got %v", statErr)
	}
}

func TestCopyEntryRegularFile(t *testing.T) {
	src := filepath.Join(t.TempDir(), "Bookmarks")
	if err := os.WriteFile(src, []byte("data"), 0o600); err != nil {
		t.Fatal(err)
	}
	dst := filepath.Join(t.TempDir(), "nested", "Bookmarks")
	n, err := copyEntry(src, dst)
	if err != nil || n != 1 {
		t.Fatalf("regular file: want (1,nil), got (%d,%v)", n, err)
	}
	got, _ := os.ReadFile(dst) //nolint:gosec // test
	if string(got) != "data" {
		t.Fatalf("dst content = %q", got)
	}
}

func TestCopyDirNestedTreeCountsAllFiles(t *testing.T) {
	src := t.TempDir()
	mustWrite(t, filepath.Join(src, "a.txt"), "1")
	mustWrite(t, filepath.Join(src, "sub", "b.txt"), "2")
	mustWrite(t, filepath.Join(src, "sub", "deep", "c.txt"), "3")

	dst := filepath.Join(t.TempDir(), "out")
	n, err := copyDir(src, dst)
	if err != nil || n != 3 {
		t.Fatalf("nested: want (3,nil), got (%d,%v)", n, err)
	}
	for _, rel := range []string{"a.txt", "sub/b.txt", "sub/deep/c.txt"} {
		if _, err := os.Stat(filepath.Join(dst, filepath.FromSlash(rel))); err != nil {
			t.Fatalf("missing %s: %v", rel, err)
		}
	}
}

func TestCopyDirEmptyDirectoryCreatesDestination(t *testing.T) {
	src := t.TempDir()
	dst := filepath.Join(t.TempDir(), "out")
	n, err := copyDir(src, dst)
	if err != nil || n != 0 {
		t.Fatalf("empty: want (0,nil), got (%d,%v)", n, err)
	}
	if info, statErr := os.Stat(dst); statErr != nil || !info.IsDir() {
		t.Fatalf("dst dir not created: %v", statErr)
	}
}

func TestCopyChromeProfileSkipsMissingEntries(t *testing.T) {
	src := t.TempDir()
	mustWrite(t, filepath.Join(src, "Bookmarks"), "{}")
	mustWrite(t, filepath.Join(src, "Preferences"), "{}")
	dst := filepath.Join(t.TempDir(), "out")
	n, err := copyChromeProfile(src, dst)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if n != 2 {
		t.Fatalf("count = %d, want 2", n)
	}
}

func TestHandleChromeFileOpenErrorSkipsLockFile(t *testing.T) {
	src := filepath.Join("any", constants.ChromeProfileLockFileName)
	stderr := captureStderr(t, func() {
		copied, err := handleChromeFileOpenError(src, "dst", errors.New("locked"))
		if err != nil || copied {
			t.Fatalf("LOCK: want (false,nil), got (%v,%v)", copied, err)
		}
	})
	if !strings.Contains(stderr, "skipped volatile Chrome lock file") {
		t.Fatalf("missing warn banner: %q", stderr)
	}
}

func TestHandleChromeFileOpenErrorPropagatesNonLockErrors(t *testing.T) {
	cause := errors.New("perm")
	_, err := handleChromeFileOpenError("Bookmarks", "dst", cause)
	var ce *chromeProfileCopyError
	if !errors.As(err, &ce) {
		t.Fatalf("want *chromeProfileCopyError, got %T", err)
	}
	if ce.Op != constants.ChromeProfileCopyOpRead || !errors.Is(ce.Err, cause) {
		t.Fatalf("unexpected wrapped err: %+v", ce)
	}
}

func TestHandleChromeFileCopyErrorSkipsLockFile(t *testing.T) {
	src := filepath.Join("Local Extension Settings", "abc", constants.ChromeProfileLockFileName)
	copied, err := handleChromeFileCopyError(src, "dst", errors.New("io"))
	if err != nil || copied {
		t.Fatalf("LOCK mid-copy: want (false,nil), got (%v,%v)", copied, err)
	}
}

func TestIsChromeVolatileLockFileDetectsExactBasename(t *testing.T) {
	cases := map[string]bool{
		"LOCK":        true,
		"a/b/LOCK":    true,
		"LOCK.txt":    false,
		"locked":      false,
		"prefix-LOCK": false,
		"LOCK/child":  false,
	}
	for in, want := range cases {
		if got := isChromeVolatileLockFile(filepath.FromSlash(in)); got != want {
			t.Errorf("%q: got %v want %v", in, got, want)
		}
	}
}

func TestUnwrapChromeProfileCopyErrorFallsBackForPlainError(t *testing.T) {
	got := unwrapChromeProfileCopyError(errors.New("boom"))
	if got.Source != constants.ChromeProfileCopyUnknown || got.Op != constants.ChromeProfileCopyOpCopy {
		t.Fatalf("fallback shape: %+v", got)
	}
}

// Unreadable non-LOCK wrapped-error contract is covered by
// TestHandleChromeFileOpenErrorPropagatesNonLockErrors above — keeping
// the assertion at the helper boundary avoids platform-specific stat /
// mkdir behavior differences between Windows and Unix runners.

func mustWrite(t *testing.T, path, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
}

// captureStderr lives in capturestderr_testhelper_test.go.
