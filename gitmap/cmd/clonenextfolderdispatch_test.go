// Package cmd — clonenextfolderdispatch_test.go covers the v3.117.0
// folder-arg forms of `gitmap cn`. Tests focus on classification +
// path resolution; the actual clone pipeline (runCloneNext) is
// exercised by gitmap/clonenext/* unit tests and remains untouched.
//
// We deliberately do NOT call tryFolderArgCloneNext end-to-end here
// because its success path invokes performCrossDirCloneNext →
// runCloneNext → real `git clone`. Instead each test exercises the
// classification + resolution helpers in isolation, which is where
// every interesting failure mode lives.
package cmd

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLooksLikeVersionAcceptsBumpShortcuts(t *testing.T) {
	t.Parallel()

	cases := map[string]bool{
		"v++":     true,
		"v+1":     true,
		"v+12":    true,
		"V++":     false, // intentionally case-sensitive — clonenext.ResolveTarget lowercases before its own check
		"v1.2.3":  true,
		"1.2.3":   true,
		"v3.31.0": true,
		"v+abc":   false,
		"v+":      false,
		"vabc":    false,
		"gitmap":  false,
		"":        false,
	}

	for in, want := range cases {
		if got := looksLikeVersion(in); got != want {
			t.Errorf("looksLikeVersion(%q) = %v, want %v", in, got, want)
		}
	}
}

func TestHasFolderHint(t *testing.T) {
	t.Parallel()

	cases := map[string]bool{
		"~/dev":      true,
		"~":          true,
		"./repo":     true,
		"../sibling": true,
		"/abs/path":  true,
		`C:\Windows`: true,
		"repo/sub":   true,
		"plain-name": false,
		"v++":        false,
		"":           false,
	}

	for in, want := range cases {
		if got := hasFolderHint(in); got != want {
			t.Errorf("hasFolderHint(%q) = %v, want %v", in, got, want)
		}
	}
}

func TestResolveCloneNextFolder(t *testing.T) {
	tmp := t.TempDir()
	subdir := filepath.Join(tmp, "macro-ahk-v11")
	if err := os.Mkdir(subdir, 0o755); err != nil {
		t.Fatalf("setup: %v", err)
	}
	regularFile := filepath.Join(tmp, "not-a-dir")
	if err := os.WriteFile(regularFile, []byte("x"), 0o600); err != nil {
		t.Fatalf("setup: %v", err)
	}

	t.Run("absolute path resolves", func(t *testing.T) {
		got, err := resolveCloneNextFolder(subdir)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if got != subdir {
			t.Errorf("got %q, want %q", got, subdir)
		}
	})

	t.Run("relative path resolves against cwd", func(t *testing.T) {
		// t.Chdir auto-restores in the correct order relative to any t.TempDir
		// cleanups in this subtest, avoiding "/" cwd fallbacks on macOS.
		t.Chdir(tmp)

		got, err := resolveCloneNextFolder("macro-ahk-v11")
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		// macOS sometimes resolves /var → /private/var; compare suffixes.
		if !strings.HasSuffix(got, "macro-ahk-v11") {
			t.Errorf("got %q, want suffix macro-ahk-v11", got)
		}
	})

	t.Run("non-existent path returns error", func(t *testing.T) {
		_, err := resolveCloneNextFolder(filepath.Join(tmp, "ghost"))
		if err == nil {
			t.Fatal("expected error for missing folder")
		}
	})

	t.Run("file (not dir) returns errCNFolderNotDir", func(t *testing.T) {
		_, err := resolveCloneNextFolder(regularFile)
		if !errors.Is(err, errCNFolderNotDir) {
			t.Errorf("got %v, want errCNFolderNotDir", err)
		}
	})
}

func TestIsFolderShapedHints(t *testing.T) {
	t.Parallel()

	// Path-hint tokens are folder-shaped regardless of on-disk state —
	// the dispatcher uses this to escalate "missing folder" to a hard
	// error instead of silently falling through to the alias resolver.
	hints := []string{"./missing", "../also-missing", "~/nope", "/abs/ghost", `C:\ghost`}
	for _, h := range hints {
		if !isFolderShaped(h) {
			t.Errorf("isFolderShaped(%q) = false, want true (path-hint should win regardless of stat)", h)
		}
	}

	// Tokens with no path hint AND no on-disk match must NOT be folder-
	// shaped — otherwise bare alias names like "gitmap" would be hijacked.
	if isFolderShaped("definitely-not-a-real-folder-xyz-12345") {
		t.Error("bare non-existent name must not be folder-shaped (would shadow alias resolver)")
	}
}

func TestExpandTildeUsedByResolver(t *testing.T) {
	t.Parallel()

	// expandTilde is defined in updaterepo.go in the same package; this
	// test pins that resolveCloneNextFolder actually calls it (via a
	// "~" prefix that would otherwise stat-fail). We don't assert the
	// resolved home dir contents — just that expansion happened.
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("no home dir available")
	}

	got, err := resolveCloneNextFolder("~")
	if err != nil {
		t.Fatalf("unexpected err for ~: %v", err)
	}
	if got != home {
		t.Errorf("got %q, want %q (tilde expansion broke)", got, home)
	}
}
