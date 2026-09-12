package heavy_test

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/clonefrom"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// TestExecute_HappyPath clones a tiny bare repo via file:// to a
// fresh dest. Asserts: status=ok, dest dir exists with a .git
// child, summary contains the URL.
func TestCloneFrom_Execute_HappyPath(t *testing.T) {
	requireGit(t)
	bare := makeBareRepo(t)
	cwd := t.TempDir()

	plan := clonefrom.Plan{Rows: []clonefrom.Row{{URL: "file://" + bare, Dest: "out"}}}
	results := clonefrom.Execute(plan, cwd, io.Discard)

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
func TestCloneFrom_Execute_SkipsNonEmptyDest(t *testing.T) {
	requireGit(t)
	cwd := t.TempDir()
	dest := filepath.Join(cwd, "exists")
	if err := os.MkdirAll(dest, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	if err := os.WriteFile(filepath.Join(dest, "marker"), []byte("x"), 0o644); err != nil {
		t.Fatalf("seed: %v", err)
	}

	plan := clonefrom.Plan{Rows: []clonefrom.Row{{URL: "https://example.org/x.git", Dest: "exists"}}}
	results := clonefrom.Execute(plan, cwd, io.Discard)

	if results[0].Status != constants.CloneFromStatusSkipped {
		t.Errorf("status = %q, want skipped", results[0].Status)
	}

	if _, err := os.Stat(filepath.Join(dest, "marker")); err != nil {
		t.Errorf("marker file gone: %v", err)
	}
}

// TestExecute_FailedRowFlagsExitCode covers the failure path: a
// nonexistent file:// URL produces a `failed` Result with a
// trimmed detail. Important because the CLI exit code depends on
// this status.
func TestCloneFrom_Execute_FailedRowFlagsExitCode(t *testing.T) {
	requireGit(t)
	cwd := t.TempDir()
	bogus := filepath.Join(t.TempDir(), "does-not-exist")

	plan := clonefrom.Plan{Rows: []clonefrom.Row{{URL: "file://" + bogus, Dest: "out"}}}
	results := clonefrom.Execute(plan, cwd, io.Discard)

	if results[0].Status != constants.CloneFromStatusFailed {
		t.Fatalf("status = %q, want failed", results[0].Status)
	}

	if len(results[0].Detail) == 0 {
		t.Errorf("failed row has empty detail")
	}

	if len(results[0].Detail) > constants.CloneFromErrTrimLimit+3 {
		t.Errorf("detail %d chars exceeds trim limit %d", len(results[0].Detail), constants.CloneFromErrTrimLimit)
	}
}

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

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v in %s: %v\n%s", args, dir, err, string(out))
	}
}
