package committransfer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestCountParityMainline_RunRight is the regression guard for the
// 212→150 commit-count mismatch. A pure-mainline source repo with N
// commits MUST produce N commits on the target after RunRight, and
// the reconciliation invariant must hold:
//
//	source-considered == replayed + skipped(all) + merge-excluded
//
// Issue: .lovable/memory/issues/2026-05-09-commit-transfer-count-mismatch.md
func TestCountParityMainline_RunRight(t *testing.T) {
	t.Parallel()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skipf("git not available: %v", err)
	}

	root := t.TempDir()
	source := mustInitCountRepo(t, filepath.Join(root, "src"))
	target := mustInitCountRepo(t, filepath.Join(root, "dst"))
	for i := 1; i <= 5; i++ {
		mustCommitCount(t, source, fmt.Sprintf("f%d.txt", i),
			fmt.Sprintf("v%d\n", i), fmt.Sprintf("commit %d", i))
	}

	plan, err := BuildPlan(source, target, Options{LogPrefix: "[t]"})
	if err != nil {
		t.Fatalf("BuildPlan: %v", err)
	}
	if got := len(plan.Commits); got != 5 {
		t.Fatalf("plan.Commits = %d, want 5", got)
	}
	if plan.MergeExcluded != 0 {
		t.Errorf("plan.MergeExcluded = %d, want 0", plan.MergeExcluded)
	}
	res, err := Replay(plan, Options{Yes: true, NoPush: true, LogPrefix: "[t]"})
	if err != nil {
		t.Fatalf("Replay: %v", err)
	}
	if res.Replayed != 5 {
		t.Errorf("res.Replayed = %d, want 5", res.Replayed)
	}
	assertReconcile(t, plan, res)
	if got := countCommits(t, target); got != 5 {
		t.Errorf("target commit count = %d, want 5", got)
	}
}

// TestCountParityMergeExcluded asserts that merge commits stripped by
// the default --no-merges path are reported in plan.MergeExcluded so
// the user can reconcile against `git log` (which counts merges).
func TestCountParityMergeExcluded(t *testing.T) {
	t.Parallel()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skipf("git not available: %v", err)
	}

	root := t.TempDir()
	source := mustInitCountRepo(t, filepath.Join(root, "src"))
	target := mustInitCountRepo(t, filepath.Join(root, "dst"))

	mustCommitCount(t, source, "base.txt", "0\n", "base")
	gitInDir(t, source, "checkout", "-b", "feature")
	mustCommitCount(t, source, "feat.txt", "1\n", "feature work")
	gitInDir(t, source, "checkout", "main")
	mustCommitCount(t, source, "main2.txt", "2\n", "main work")
	// True merge (not fast-forward).
	gitInDir(t, source, "merge", "--no-ff", "-m", "merge feature", "feature")

	plan, err := BuildPlan(source, target, Options{LogPrefix: "[t]"})
	if err != nil {
		t.Fatalf("BuildPlan: %v", err)
	}
	if plan.MergeExcluded < 1 {
		t.Errorf("plan.MergeExcluded = %d, want >= 1 (the --no-ff merge)",
			plan.MergeExcluded)
	}
	res, err := Replay(plan, Options{Yes: true, NoPush: true, LogPrefix: "[t]"})
	if err != nil {
		t.Fatalf("Replay: %v", err)
	}
	assertReconcile(t, plan, res)
}

// assertReconcile is the count-parity invariant shared by every test.
func assertReconcile(t *testing.T, plan ReplayPlan, res ReplayResult) {
	t.Helper()
	considered := len(plan.Commits) + plan.MergeExcluded
	accounted := res.Replayed + res.SkippedDrop + res.SkippedReplayed +
		res.SkippedEmpty + plan.MergeExcluded
	if considered != accounted {
		t.Errorf("reconcile failed: considered=%d accounted=%d (replayed=%d drop=%d alreadyReplayed=%d empty=%d mergeExcl=%d)",
			considered, accounted, res.Replayed, res.SkippedDrop,
			res.SkippedReplayed, res.SkippedEmpty, plan.MergeExcluded)
	}
}

// mustInitCountRepo creates a fresh repo at dir with deterministic
// identity. Kept local to this test file so it does not collide with
// other test helpers in the package.
func mustInitCountRepo(t *testing.T, dir string) string {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", dir, err)
	}
	gitInDir(t, dir, "init", "-b", "main")
	gitInDir(t, dir, "config", "user.email", "e2e@gitmap.test")
	gitInDir(t, dir, "config", "user.name", "E2E Bot")
	gitInDir(t, dir, "config", "commit.gpgsign", "false")

	return dir
}

// mustCommitCount writes path with body and commits via porcelain. Uses
// porcelain (not plumbing) so the working-tree state is realistic for
// snapshot-copy paths under test.
func mustCommitCount(t *testing.T, dir, path, body, msg string) {
	t.Helper()
	full := filepath.Join(dir, path)
	if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
		t.Fatalf("write %s: %v", full, err)
	}
	gitInDir(t, dir, "add", path)
	stamp := time.Now().UTC().Format(time.RFC3339)
	cmd := exec.Command("git", "-C", dir, "commit", "-m", msg)
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_DATE="+stamp,
		"GIT_COMMITTER_DATE="+stamp,
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git commit (%s): %v\n%s", msg, err, out)
	}
}

// gitInDir runs git -C dir <args> and fatals on error.
func gitInDir(t *testing.T, dir string, args ...string) {
	t.Helper()
	full := append([]string{"-C", dir}, args...)
	if out, err := exec.Command("git", full...).CombinedOutput(); err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
}

// countCommits returns the number of reachable commits on HEAD via
// `git rev-list --count HEAD`.
func countCommits(t *testing.T, dir string) int {
	t.Helper()
	out, err := exec.Command("git", "-C", dir, "rev-list", "--count", "HEAD").Output()
	if err != nil {
		t.Fatalf("rev-list --count: %v", err)
	}
	var n int
	if _, perr := fmt.Sscanf(strings.TrimSpace(string(out)), "%d", &n); perr != nil {
		t.Fatalf("parse rev-list count %q: %v", out, perr)
	}

	return n
}
