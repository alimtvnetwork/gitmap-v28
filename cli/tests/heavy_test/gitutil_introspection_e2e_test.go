package heavy_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
)

func TestGitIntrospectionSuite(t *testing.T) {
	tempDir := t.TempDir()
	cmdInit := exec.Command("git", "init", tempDir)
	if err := cmdInit.Run(); err != nil {
		t.Skip("git not available in test environment")
	}

	dummyFile := filepath.Join(tempDir, "README.md")
	_ = os.WriteFile(dummyFile, []byte("# Test"), 0644)

	_ = exec.Command("git", "-C", tempDir, "config", "user.name", "Tester").Run()
	_ = exec.Command("git", "-C", tempDir, "config", "user.email", "test@example.com").Run()
	_ = exec.Command("git", "-C", tempDir, "add", ".").Run()
	_ = exec.Command("git", "-C", tempDir, "commit", "-m", "initial commit").Run()

	sha := gitutil.GetLastCommitSHA(tempDir)
	if len(sha) == 0 || sha == "-" {
		t.Fatalf("unexpected sha: %s", sha)
	}

	branch := gitutil.GetActiveBranch(tempDir)
	if branch == "" || branch == "-" {
		t.Fatalf("unexpected branch: %s", branch)
	}
}
