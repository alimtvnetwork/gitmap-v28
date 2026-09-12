package heavy_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmd"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
)

func TestReconcileWorkflowE2E(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("USERPROFILE", t.TempDir())

	tempDir := t.TempDir()
	repoDir := initDummyGitRepoWithRemote(t, tempDir)

	diag := gitutil.InspectDirtyState(repoDir)
	recipes := gitutil.GenerateRemediationRecipes(repoDir, diag)
	item := cmd.RemediationItem{
		RepoName:      "sample-repo",
		RepoPath:      repoDir,
		SummaryReason: diag.SummaryReason,
		Recipes:       recipes,
	}

	_ = cmd.SaveRemediationState([]cmd.RemediationItem{item})

	err := cmd.RunReconcileCmd([]string{"sample-repo", "discard"})
	if err != nil {
		t.Fatalf("runReconcileCmd failed: %v", err)
	}

	remaining := cmd.LoadRemediationState()
	if len(remaining) != 0 {
		t.Fatalf("expected 0 remaining items, got %d", len(remaining))
	}
}

func initDummyGitRepoWithRemote(t *testing.T, baseDir string) string {
	t.Helper()
	remoteDir := filepath.Join(baseDir, "remote.git")
	runCmdIn(baseDir, "git", "init", "--bare", remoteDir)

	repoDir := filepath.Join(baseDir, "local")
	runCmdIn(baseDir, "git", "clone", remoteDir, repoDir)
	runCmdIn(repoDir, "git", "config", "user.name", "Tester")
	runCmdIn(repoDir, "git", "config", "user.email", "tester@example.com")

	_ = os.WriteFile(filepath.Join(repoDir, "tracked.txt"), []byte("tracked"), 0644)
	runCmdIn(repoDir, "git", "add", "tracked.txt")
	runCmdIn(repoDir, "git", "commit", "-m", "init")
	runCmdIn(repoDir, "git", "push", "origin", "HEAD")

	_ = os.WriteFile(filepath.Join(repoDir, "untracked.txt"), []byte("dirty"), 0644)

	return repoDir
}

func runCmdIn(dir, name string, args ...string) {
	c := exec.Command(name, args...)
	c.Dir = dir
	_ = c.Run()
}
