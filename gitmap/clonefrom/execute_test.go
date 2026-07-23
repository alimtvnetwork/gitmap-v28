package clonefrom

// Executor end-to-end test using a LOCAL bare repo as the clone
// source. No network required — the test creates a bare repo in a
// tempdir, then asks Execute to clone it elsewhere via a `file://`
// URL. Skips on hosts without `git` on PATH so the rest of the
// test suite still runs (e.g. minimal CI containers).

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// TestExecute_HappyPath clones a tiny bare repo via file:// to a
// fresh dest. Asserts: status=ok, dest dir exists with a .git
// child, summary contains the URL.
func TestExecute_HappyPath(t *testing.T) {
	requireGit(t)
	bare := makeBareRepo(t)
	cwd := t.TempDir()

	plan := Plan{Rows: []Row{{URL: "file://" + bare, Dest: "out"}}}
	results := Execute(plan, cwd, io.Discard)

	if len(results) != 1 {
		t.Fatalf("results = %d, want 1", len(results))
	}
	if results[0].Status != constants.CloneFromStatusOK {
		t.Fatalf("status = %q, want ok (detail=%q)", results[0].Status, results[0].Detail)
	}
	gitDir := filepath.Join(cwd, "out", ".git")
	if _, err := os.Stat(gitDir); err != nil {
		t.Errorf("expected .git dir at %s: %v", gitDir, err)
	}
}

// TestExecute_SkipsNonEmptyDest pre-creates the dest with a file
// in it and confirms Execute marks the row skipped without
// invoking git. Idempotent re-run guarantee.
func TestExecute_SkipsNonEmptyDest(t *testing.T) {
	requireGit(t)
	cwd := t.TempDir()
	dest := filepath.Join(cwd, "exists")
	if err := os.MkdirAll(dest, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dest, "marker"), []byte("x"), 0o644); err != nil {
		t.Fatalf("seed: %v", err)
	}

	plan := Plan{Rows: []Row{{URL: "https://example.org/x.git", Dest: "exists"}}}
	results := Execute(plan, cwd, io.Discard)

	if results[0].Status != constants.CloneFromStatusSkipped {
		t.Errorf("status = %q, want skipped", results[0].Status)
	}
	// Marker file should still be there — proof we didn't recurse-delete.
	if _, err := os.Stat(filepath.Join(dest, "marker")); err != nil {
		t.Errorf("marker file gone: %v", err)
	}
}

// TestExecute_FailedRowFlagsExitCode covers the failure path: a
// nonexistent file:// URL produces a `failed` Result with a
// trimmed detail. Important because the CLI exit code depends on
// this status.
func TestExecute_FailedRowFlagsExitCode(t *testing.T) {
	requireGit(t)
	cwd := t.TempDir()
	bogus := filepath.Join(t.TempDir(), "does-not-exist")

	plan := Plan{Rows: []Row{{URL: "file://" + bogus, Dest: "out"}}}
	results := Execute(plan, cwd, io.Discard)

	if results[0].Status != constants.CloneFromStatusFailed {
		t.Fatalf("status = %q, want failed", results[0].Status)
	}
	if len(results[0].Detail) == 0 {
		t.Errorf("failed row has empty detail")
	}
	if len(results[0].Detail) > constants.CloneFromErrTrimLimit+3 { // +3 for "..."
		t.Errorf("detail %d chars exceeds trim limit %d", len(results[0].Detail), constants.CloneFromErrTrimLimit)
	}
}

// TestRenderSummary_TalliesAllStatuses asserts the header counts
// every status correctly. Pure assembly test — uses synthesized
// Result values, no git.
func TestRenderSummary_TalliesAllStatuses(t *testing.T) {
	results := []Result{
		{Status: constants.CloneFromStatusOK, Row: Row{URL: "a"}},
		{Status: constants.CloneFromStatusOK, Row: Row{URL: "b"}},
		{Status: constants.CloneFromStatusSkipped, Row: Row{URL: "c"}, Detail: "dest exists"},
		{Status: constants.CloneFromStatusFailed, Row: Row{URL: "d"}, Detail: "boom"},
	}
	var buf bytes.Buffer
	if err := RenderSummary(&buf, results, "/tmp/r.csv"); err != nil {
		t.Fatalf("RenderSummary: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "2 ok, 1 skipped, 1 failed (4 total)") {
		t.Errorf("tally line missing or wrong:\n%s", out)
	}
	if !strings.Contains(out, "report: /tmp/r.csv") {
		t.Errorf("report path line missing")
	}
}

// requireGit skips the test when git isn't on PATH. Used by every
// executor test that actually invokes git so a minimal CI image
// without git installed runs the parser/render tests and skips
// only the clone tests.
func requireGit(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skipf("git not on PATH: %v", err)
	}
}

// makeBareRepo creates a brand-new bare repo with one commit and
// returns its absolute path. The commit ensures `git clone` has
// something to populate the working tree with — cloning an empty
// bare repo prints a warning to stderr that would mask real
// failures.
func makeBareRepo(t *testing.T) string {
	t.Helper()
	work := t.TempDir()
	bare := filepath.Join(t.TempDir(), "src.git")

	runGit(t, work, "init", "-q")
	runGit(t, work, "config", "user.email", "t@e")
	runGit(t, work, "config", "user.name", "t")
	if err := os.WriteFile(filepath.Join(work, "README"), []byte("hi"), 0o644); err != nil {
		t.Fatalf("seed file: %v", err)
	}
	runGit(t, work, "add", ".")
	runGit(t, work, "commit", "-q", "-m", "init")
	runGit(t, work, "clone", "--bare", "-q", work, bare)

	return bare
}

// runGit is a tiny exec wrapper that fatals on error with the
// combined output included in the failure message.
func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v in %s: %v\n%s", args, dir, err, string(out))
	}
}
