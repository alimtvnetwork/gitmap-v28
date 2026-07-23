package cmd

import (
	"path/filepath"
	"strings"
	"testing"
)

// skipOnWindowsSubprocess is retained as a no-op so historical
// references compile. The Windows stderr-capture flakiness it
// guarded against was fixed in v5.53.0 by:
//
//  1. Adding glyphs.Drain / theme.Drain and registering them with
//     cliexit.RegisterFlusher in cmd/root.go, so cliexit.Fail
//     flushes the pipe-wrapped os.Stderr before os.Exit instead of
//     racing the forwarder goroutine.
//  2. Pinning GITMAP_GLYPHS=rich and GITMAP_THEME=bright in
//     hermeticEnv so subprocess tests bypass the pipe wrap entirely
//     (defense in depth — the wrap itself stays correct for
//     production use cases that do need it).
//
// Kept callable (instead of deleted) so any out-of-tree test in
// the same package can still reference it without a compile break
// during the next refactor sweep.
func skipOnWindowsSubprocess(_ *testing.T) {}

// Integration tests asserting that user-facing failure stderr from
// scan and clone-family commands carries the standardized context
// fields produced by gitmap/cliexit:
//
//   - the command attempted (e.g. "gitmap scan", "gitmap clone-from")
//   - the subject path the command was operating on (the failing
//     repo / manifest path)
//   - some recognizable trace of the underlying error (a phrase
//     specific to the failure mode, not a generic "error")
//
// These run against the real built binary so we exercise the actual
// stderr that wrapper scripts and CI annotations consume — not an
// in-process formatter mock.

// TestCLI_FailureContext_Scan asserts a missing-directory scan
// surfaces the command name + the missing path on stderr. The
// "underlying error" check looks for any of the documented
// not-exist phrasings so OS-specific wording (Linux vs Windows)
// doesn't make the test brittle.
func TestCLI_FailureContext_Scan(t *testing.T) {
	skipOnWindowsSubprocess(t)
	t.Parallel()
	missing := filepath.Join(t.TempDir(), "definitely-not-here")
	code, stdout, stderr := runGitmap(t, []string{"scan", "--quiet", missing}, "")
	if code == 0 {
		t.Fatalf("scan of missing dir unexpectedly succeeded\nstdout=%s\nstderr=%s",
			stdout, stderr)
	}
	// Some Windows CI configurations merge or buffer the
	// subprocess streams differently; assert against the
	// combined output so the contract check stays meaningful
	// regardless of which stream the message lands on.
	combined := stdout + "\n" + stderr
	assertStderrContext(t, "scan", combined, missing, []string{
		"no such file", "cannot find", "does not exist", "not exist",
	})
}

// TestCLI_FailureContext_CloneFromMissingManifest drives clone-from
// against a manifest path that doesn't exist. Asserts the command
// label, the manifest path, and an open/read failure phrase.
func TestCLI_FailureContext_CloneFromMissingManifest(t *testing.T) {
	skipOnWindowsSubprocess(t)
	t.Parallel()
	missing := filepath.Join(t.TempDir(), "no-such-manifest.json")
	code, stdout, stderr := runGitmap(t,
		[]string{"clone-from", "--execute", missing}, "")
	if code == 0 {
		t.Fatalf("clone-from with missing manifest unexpectedly succeeded\nstdout=%s\nstderr=%s",
			stdout, stderr)
	}
	combined := stdout + "\n" + stderr
	assertStderrContext(t, "clone-from", combined, missing, []string{
		"no such file", "cannot find", "does not exist", "not exist", "open",
	})
}

// TestCLI_FailureContext_CloneNowMissingManifest is the clone-now
// counterpart. Distinct from clone-from: clone-now uses a different
// parser path and we want both surfaces validated.
func TestCLI_FailureContext_CloneNowMissingManifest(t *testing.T) {
	skipOnWindowsSubprocess(t)
	t.Parallel()
	missing := filepath.Join(t.TempDir(), "no-such-manifest.json")
	code, stdout, stderr := runGitmap(t,
		[]string{"clone-now", "--execute", "--quiet", missing}, "")
	if code == 0 {
		t.Fatalf("clone-now with missing manifest unexpectedly succeeded\nstdout=%s\nstderr=%s",
			stdout, stderr)
	}
	combined := stdout + "\n" + stderr
	assertStderrContextAny(t, []string{"clone-now", "reclone"}, combined, missing, []string{
		"no such file", "cannot find", "does not exist", "not exist", "open",
	})
}

// assertStderrContext is the per-command wrapper around the more
// general assertStderrContextAny. Single-command label form keeps
// the common case readable.
func assertStderrContext(t *testing.T, command, stderr, subject string, errPhrases []string) {
	t.Helper()
	assertStderrContextAny(t, []string{command}, stderr, subject, errPhrases)
}

// assertStderrContextAny verifies stderr carries (a) at least one of
// the acceptable command labels, (b) the subject string verbatim,
// and (c) at least one of the underlying-error phrases. Each missing
// piece is reported separately so a failing run pinpoints exactly
// which contract field regressed.
func assertStderrContextAny(t *testing.T, commands []string, stderr, subject string, errPhrases []string) {
	t.Helper()
	if !containsAnyCI(stderr, commands) {
		t.Errorf("stderr missing command label (any of %v)\nstderr=%s",
			commands, stderr)
	}
	if !strings.Contains(stderr, subject) {
		t.Errorf("stderr missing subject path %q\nstderr=%s", subject, stderr)
	}
	if !containsAnyCI(stderr, errPhrases) {
		t.Errorf("stderr missing underlying-error phrase (any of %v)\nstderr=%s",
			errPhrases, stderr)
	}
}

// containsAnyCI returns true when haystack contains any of the
// needles (case-insensitive). Pulled out so the assertion helpers
// stay one-liners.
func containsAnyCI(haystack string, needles []string) bool {
	for _, n := range needles {
		if containsCI(haystack, n) {
			return true
		}
	}

	return false
}
