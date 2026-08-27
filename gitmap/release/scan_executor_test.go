package release

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestExecuteCommitActions(t *testing.T) {
	repoDir, hash := setupTestRepo(t)
	commits := []ParsedCommit{{Hash: hash, Message: "release: v1.0.0", Version: "v1.0.0"}}

	actions, err := ExecuteCommitActions(repoDir, commits)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !actions[0].IsBranchCreated || !actions[0].IsTagCreated {
		t.Errorf("expected branch and tag created, got %+v", actions[0])
	}

	actions2, err := ExecuteCommitActions(repoDir, commits)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !actions2[0].IsBranchSkipped || !actions2[0].IsTagSkipped {
		t.Errorf("expected branch and tag skipped, got %+v", actions2[0])
	}
}

func setupTestRepo(t *testing.T) (string, string) {
	t.Helper()
	dir := t.TempDir()
	initGitRepo(t, dir)
	return dir, createInitialCommit(t, dir)
}

func initGitRepo(t *testing.T, dir string) {
	t.Helper()
	runTestGitCmd(t, dir, "init")
	runTestGitCmd(t, dir, "config", "user.name", "test")
	runTestGitCmd(t, dir, "config", "user.email", "test@test.com")
}

func createInitialCommit(t *testing.T, dir string) string {
	t.Helper()
	file := filepath.Join(dir, "test.txt")
	os.WriteFile(file, []byte("hello"), 0644)
	runTestGitCmd(t, dir, "add", "test.txt")
	runTestGitCmd(t, dir, "commit", "-m", "init")
	out := runTestGitCmdOutput(t, dir, "rev-parse", "HEAD")
	return strings.TrimSpace(string(out))
}

func runTestGitCmd(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("git %v failed: %v", args, err)
	}
}

func runTestGitCmdOutput(t *testing.T, dir string, args ...string) []byte {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("git %v failed: %v", args, err)
	}
	return out
}
