// Package e2e provides shared fixture builders and run helpers for
// the commit-in end-to-end test suites (Steps 9–12 of the commit-in
// implementation plan).
//
// All builders shell out to the real `git` binary so the tests cover
// the same code paths users hit in production. Tests that need a fast
// fake-runner path should stay in their respective sub-packages
// (walk, replay, etc.) — this package is intentionally heavy.
package e2e

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// Repo is a real on-disk git repo built inside t.TempDir(). Cleanup is
// automatic via t.Cleanup; callers never need to RemoveAll.
type Repo struct {
	t    *testing.T
	Path string
}

// NewRepo creates an empty repo at <tempDir>/<name> with a deterministic
// initial branch (`main`) and identity. Skips the test if `git` is not
// in PATH so CI without git degrades cleanly.
func NewRepo(t *testing.T, name string) *Repo {
	t.Helper()
	return NewRepoIn(t, t.TempDir(), name)
}

// NewRepoIn is like NewRepo but creates the repo inside an explicit
// parent directory. Used by sibling-discovery tests that need several
// repos sharing one parent (so `all` / `-N` keyword scope works).
func NewRepoIn(t *testing.T, parent, name string) *Repo {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skipf("git not available: %v", err)
	}
	dir := filepath.Join(parent, name)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", dir, err)
	}
	r := &Repo{t: t, Path: dir}
	r.git("init", "-b", "main")
	r.git("config", "user.email", "e2e@gitmap.test")
	r.git("config", "user.name", "E2E Bot")
	r.git("config", "commit.gpgsign", "false")
	return r
}

// Commit writes `path` with `body`, stages it, and commits with
// `message` at `when` (both author and committer date). Returns the
// new commit SHA. Uses plumbing (write-tree + commit-tree + update-ref)
// so the harness works inside sandboxes that block porcelain `git
// add` / `git commit`.
func (r *Repo) Commit(path, body, message string, when time.Time) string {
	r.t.Helper()
	full := filepath.Join(r.Path, path)
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		r.t.Fatalf("mkdir parent of %s: %v", full, err)
	}
	if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
		r.t.Fatalf("write %s: %v", full, err)
	}
	r.git("update-index", "--add", path)
	tree := r.gitOut("write-tree")
	stamp := when.UTC().Format(time.RFC3339)
	args := []string{"commit-tree", tree, "-m", message}
	if parent, ok := r.headShaOpt(); ok {
		args = append(args, "-p", parent)
	}
	cmd := exec.Command("git", args...)
	cmd.Dir = r.Path
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_DATE="+stamp,
		"GIT_COMMITTER_DATE="+stamp,
		"GIT_AUTHOR_NAME=E2E Bot",
		"GIT_AUTHOR_EMAIL=e2e@gitmap.test",
		"GIT_COMMITTER_NAME=E2E Bot",
		"GIT_COMMITTER_EMAIL=e2e@gitmap.test",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		r.t.Fatalf("git commit-tree: %v\n%s", err, out)
	}
	sha := strings.TrimSpace(string(out))
	r.git("update-ref", "refs/heads/main", sha)
	r.git("symbolic-ref", "HEAD", "refs/heads/main")
	return sha
}

// gitOut runs a git subcommand and returns trimmed stdout; fatals on
// error. Used for plumbing reads (write-tree, rev-parse).
func (r *Repo) gitOut(args ...string) string {
	r.t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = r.Path
	out, err := cmd.Output()
	if err != nil {
		r.t.Fatalf("git %s: %v", strings.Join(args, " "), err)
	}
	return strings.TrimSpace(string(out))
}

// headShaOpt returns (sha, true) when refs/heads/main exists, else
// ("", false). Used by Commit to decide whether to pass `-p <parent>`.
func (r *Repo) headShaOpt() (string, bool) {
	cmd := exec.Command("git", "rev-parse", "--verify", "--quiet", "refs/heads/main")
	cmd.Dir = r.Path
	out, err := cmd.Output()
	if err != nil {
		return "", false
	}
	sha := strings.TrimSpace(string(out))
	return sha, sha != ""
}

// git runs a git subcommand inside r.Path and fatals on error. Used
// for setup-only commands where output is uninteresting.
func (r *Repo) git(args ...string) {
	r.t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = r.Path
	if out, err := cmd.CombinedOutput(); err != nil {
		r.t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
}

// MustExist fatals if `rel` (relative to r.Path) is missing.
func (r *Repo) MustExist(rel string) {
	r.t.Helper()
	if _, err := os.Stat(filepath.Join(r.Path, rel)); err != nil {
		r.t.Fatalf("expected %s to exist: %v", rel, err)
	}
}

// String renders for %v formatting — useful in test failure output.
func (r *Repo) String() string { return fmt.Sprintf("Repo(%s)", r.Path) }
