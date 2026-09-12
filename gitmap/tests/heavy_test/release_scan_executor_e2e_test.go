package heavy_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/release"
)

func TestRelease_ExecuteCommitActions(t *testing.T) {
	repoDir, hash := setupReleaseTestRepo(t)
	commits := []release.ParsedCommit{{Hash: hash, Message: "release: v1.0.0", Version: "v1.0.0"}}

	actions, err := release.ExecuteCommitActions(repoDir, commits)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !actions[0].IsBranchCreated || !actions[0].IsTagCreated {
		t.Errorf("expected branch and tag created, got %+v", actions[0])
	}

	actions2, err := release.ExecuteCommitActions(repoDir, commits)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !actions2[0].IsBranchSkipped || !actions2[0].IsTagSkipped {
		t.Errorf("expected branch and tag skipped, got %+v", actions2[0])
	}
}

func setupReleaseTestRepo(t *testing.T) (string, string) {
	t.Helper()
	dir := t.TempDir()
	initReleaseGitRepo(t, dir)

	return dir, createReleaseInitialCommit(t, dir)
}

func initReleaseGitRepo(t *testing.T, dir string) {
	t.Helper()
	runReleaseTestGitCmd(t, dir, "init")
	runReleaseTestGitCmd(t, dir, "config", "user.name", "test")
	runReleaseTestGitCmd(t, dir, "config", "user.email", "test@test.com")
}

func createReleaseInitialCommit(t *testing.T, dir string) string {
	t.Helper()
	file := filepath.Join(dir, "test.txt")
	if err := os.WriteFile(file, []byte("hello"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	runReleaseTestGitCmd(t, dir, "add", "test.txt")
	runReleaseTestGitCmd(t, dir, "commit", "-m", "init")
	out := runReleaseTestGitCmdOutput(t, dir, "rev-parse", "HEAD")

	return strings.TrimSpace(string(out))
}

func runReleaseTestGitCmd(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("git %v failed: %v", args, err)
	}
}

func runReleaseTestGitCmdOutput(t *testing.T, dir string, args ...string) []byte {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("git %v failed: %v", args, err)
	}

	return out
}
