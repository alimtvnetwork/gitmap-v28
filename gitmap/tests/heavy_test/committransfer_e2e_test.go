package heavy_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/committransfer"
)

func TestCommitTransfer_CountParityMainline_RunRight(t *testing.T) {
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

	plan, err := committransfer.BuildPlan(source, target, committransfer.Options{LogPrefix: "[t]"})
	if err != nil {
		t.Fatalf("BuildPlan: %v", err)
	}

	if got := len(plan.Commits); got != 5 {
		t.Fatalf("plan.Commits = %d, want 5", got)
	}

	if plan.MergeExcluded != 0 {
		t.Errorf("plan.MergeExcluded = %d, want 0", plan.MergeExcluded)
	}

	res, err := committransfer.Replay(plan, committransfer.Options{Yes: true, NoPush: true, LogPrefix: "[t]"})
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

func TestCommitTransfer_CountParityMergeExcluded(t *testing.T) {
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
	gitInDir(t, source, "merge", "--no-ff", "-m", "merge feature", "feature")

	plan, err := committransfer.BuildPlan(source, target, committransfer.Options{LogPrefix: "[t]"})
	if err != nil {
		t.Fatalf("BuildPlan: %v", err)
	}

	if plan.MergeExcluded < 1 {
		t.Errorf("plan.MergeExcluded = %d, want >= 1 (the --no-ff merge)", plan.MergeExcluded)
	}

	res, err := committransfer.Replay(plan, committransfer.Options{Yes: true, NoPush: true, LogPrefix: "[t]"})
	if err != nil {
		t.Fatalf("Replay: %v", err)
	}

	assertReconcile(t, plan, res)
}

func TestCommitTransfer_DirtySource_RejectsRunRight(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	source := mustInitCountRepo(t, filepath.Join(root, "src"))
	target := mustInitCountRepo(t, filepath.Join(root, "dst"))

	mustCommitCount(t, source, "file1.txt", "content1", "commit 1")
	mustCommitCount(t, source, "file2.txt", "content2", "commit 2")

	dirtyFile := filepath.Join(source, "file1.txt")
	if err := os.WriteFile(dirtyFile, []byte("uncommitted change"), 0644); err != nil {
		t.Fatalf("write dirty file: %v", err)
	}

	opts := committransfer.Options{
		LogPrefix: "[commit-right]",
		DryRun:    false,
		Yes:       true,
		NoPush:    true,
	}

	err := committransfer.RunRight(source, target, opts)
	if err == nil {
		t.Fatal("expected RunRight to fail on dirty source, but it succeeded")
	}

	if !strings.Contains(err.Error(), "uncommitted changes") {
		t.Fatalf("expected uncommitted changes error, got: %v", err)
	}
}

func assertReconcile(t *testing.T, plan committransfer.ReplayPlan, res committransfer.ReplayResult) {
	t.Helper()
	considered := len(plan.Commits) + plan.MergeExcluded
	accounted := res.Replayed + res.SkippedDrop + res.SkippedReplayed +
		res.SkippedEmpty + plan.MergeExcluded
	if considered != accounted {
		t.Errorf("reconcile failed: considered=%d accounted=%d", considered, accounted)
	}
}

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

func gitInDir(t *testing.T, dir string, args ...string) {
	t.Helper()
	full := append([]string{"-C", dir}, args...)
	if out, err := exec.Command("git", full...).CombinedOutput(); err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
}

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
