package cmd_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cmd"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/gitutil"
)

func TestPullWorkflowE2E(t *testing.T) {
	tempDir := t.TempDir()

	// Create test git repo
	repoDir := filepath.Join(tempDir, "sample-repo")
	_ = os.MkdirAll(repoDir, 0755)
	if err := exec.Command("git", "init", repoDir).Run(); err != nil {
		t.Skip("git not available in environment")
	}

	_ = exec.Command("git", "-C", repoDir, "config", "user.name", "Tester").Run()
	_ = exec.Command("git", "-C", repoDir, "config", "user.email", "test@example.com").Run()
	_ = os.WriteFile(filepath.Join(repoDir, "file.txt"), []byte("hello"), 0644)
	_ = exec.Command("git", "-C", repoDir, "add", ".").Run()
	_ = exec.Command("git", "-C", repoDir, "commit", "-m", "init").Run()

	// Make dirty
	_ = os.WriteFile(filepath.Join(repoDir, "dirty.txt"), []byte("uncommitted"), 0644)

	diag := gitutil.InspectDirtyState(repoDir)
	if !diag.IsDirty {
		t.Fatal("expected repo to be dirty")
	}

	recipes := gitutil.GenerateRemediationRecipes(repoDir, diag)
	if len(recipes) != 3 {
		t.Fatalf("expected 3 remediation options, got %d", len(recipes))
	}

	targets := cmd.ResolvePullDirectoryTargets(tempDir)
	if len(targets) != 1 {
		t.Fatalf("expected 1 discovered target in tempDir, got %d", len(targets))
	}
}
