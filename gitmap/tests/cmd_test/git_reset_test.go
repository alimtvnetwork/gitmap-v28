package cmd_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cmd"
)

func TestGitResetExecution(t *testing.T) {
	tmpDir := t.TempDir()
	origWd, _ := os.Getwd()
	defer os.Chdir(origWd)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir to tmpDir: %v", err)
	}

	exec.Command("git", "init").Run()
	exec.Command("git", "config", "user.name", "Test User").Run()
	exec.Command("git", "config", "user.email", "test@example.com").Run()

	file1 := filepath.Join(tmpDir, "file1.txt")
	os.WriteFile(file1, []byte("first commit"), 0644)
	exec.Command("git", "add", "-A").Run()
	exec.Command("git", "commit", "-m", "commit 1").Run()

	outSha1, _ := exec.Command("git", "rev-parse", "HEAD").Output()
	sha1 := string(outSha1)[:8]

	file2 := filepath.Join(tmpDir, "file2.txt")
	os.WriteFile(file2, []byte("second commit"), 0644)
	exec.Command("git", "add", "-A").Run()
	exec.Command("git", "commit", "-m", "commit 2").Run()

	// Reset back to sha1 with --no-push
	errReset := cmd.RunGitReset([]string{sha1, "--no-push"})
	if errReset != nil {
		t.Fatalf("git-reset failed: %v", errReset)
	}

	outHead, _ := exec.Command("git", "rev-parse", "HEAD").Output()
	headSha := string(outHead)[:8]

	if headSha != sha1 {
		t.Fatalf("expected HEAD to be %s after reset, got %s", sha1, headSha)
	}
}

func TestRmGitExecution(t *testing.T) {
	tmpDir := t.TempDir()
	origWd, _ := os.Getwd()
	defer os.Chdir(origWd)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir to tmpDir: %v", err)
	}

	exec.Command("git", "init").Run()
	exec.Command("git", "config", "user.name", "Test User").Run()
	exec.Command("git", "config", "user.email", "test@example.com").Run()

	file1 := filepath.Join(tmpDir, "file1.txt")
	os.WriteFile(file1, []byte("commit 1"), 0644)
	exec.Command("git", "add", "-A").Run()
	exec.Command("git", "commit", "-m", "commit 1").Run()

	file2 := filepath.Join(tmpDir, "file2.txt")
	os.WriteFile(file2, []byte("commit 2 to remove"), 0644)
	exec.Command("git", "add", "-A").Run()
	exec.Command("git", "commit", "-m", "commit 2 to remove").Run()

	outSha2, _ := exec.Command("git", "rev-parse", "HEAD").Output()
	sha2 := string(outSha2)[:8]

	file3 := filepath.Join(tmpDir, "file3.txt")
	os.WriteFile(file3, []byte("commit 3"), 0644)
	exec.Command("git", "add", "-A").Run()
	exec.Command("git", "commit", "-m", "commit 3").Run()

	// Remove commit 2 using rm-git with --no-push
	errRm := cmd.RunRmGit([]string{sha2, "--no-push"})
	if errRm != nil {
		t.Fatalf("rm-git failed: %v", errRm)
	}

	logOut, _ := exec.Command("git", "log", "--oneline").Output()
	logStr := string(logOut)

	if strings.Contains(logStr, "commit 2 to remove") {
		t.Fatalf("expected commit 2 to be dropped, but git log contains it:\n%s", logStr)
	}
}
