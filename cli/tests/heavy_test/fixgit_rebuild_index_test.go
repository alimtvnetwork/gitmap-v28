package heavy_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmd"
)

func TestFixGit_CorruptZeroByteIndex(t *testing.T) {
	tempDir := t.TempDir()

	c := exec.Command("git", "init")
	c.Dir = tempDir
	if err := c.Run(); err != nil {
		t.Skipf("git not installed or init failed: %v", err)
	}

	_ = exec.Command("git", "-C", tempDir, "config", "user.email", "test@example.com").Run()
	_ = exec.Command("git", "-C", tempDir, "config", "user.name", "Test").Run()

	testFile := filepath.Join(tempDir, "sample.txt")
	_ = os.WriteFile(testFile, []byte("hello gitmap"), 0644)

	_ = exec.Command("git", "-C", tempDir, "add", "sample.txt").Run()
	_ = exec.Command("git", "-C", tempDir, "commit", "-m", "init", "--author=Test <test@example.com>").Run()

	gitDir := filepath.Join(tempDir, ".git")
	indexPath := filepath.Join(gitDir, "index")

	// Corrupt index by truncating to 0 bytes
	_ = os.WriteFile(indexPath, []byte(""), 0644)

	opts := cmd.FixGitOptions{
		TargetDir:   tempDir,
		IsIndexOnly: true,
	}

	issues, err := cmd.RemediateGitIndex(tempDir, gitDir, opts)
	if err != nil {
		t.Fatalf("remediateGitIndex failed: %v", err)
	}

	if len(issues) != 1 {
		t.Fatalf("expected 1 index issue, got %d", len(issues))
	}

	if !issues[0].IsFixed {
		t.Errorf("expected issue to be fixed, got: %+v", issues[0])
	}

	info, statErr := os.Stat(indexPath)
	if statErr != nil || info.Size() == 0 {
		t.Errorf("index was not restored: info=%v, err=%v", info, statErr)
	}
}
