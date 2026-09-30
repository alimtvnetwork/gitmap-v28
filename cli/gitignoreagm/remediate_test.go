package gitignoreagm

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestRemediateRepoUntracksAndCommits(t *testing.T) {
	tmpDir := t.TempDir()
	if err := exec.Command("git", "init", tmpDir).Run(); err != nil {
		t.Skipf("git not available: %v", err)
	}
	_ = exec.Command("git", "-C", tmpDir, "config", "user.email", "test@example.com").Run()
	_ = exec.Command("git", "-C", tmpDir, "config", "user.name", "Test User").Run()

	resumeFile := filepath.Join(tmpDir, PrimaryIgnoreEntry)
	if err := os.WriteFile(resumeFile, []byte(`{"task":"resume"}`), 0o644); err != nil {
		t.Fatalf("write resume file: %v", err)
	}
	_ = exec.Command("git", "-C", tmpDir, "add", PrimaryIgnoreEntry).Run()
	_ = exec.Command("git", "-C", tmpDir, "commit", "-m", "initial").Run()

	if !HasUnignoredResumeTask(tmpDir) {
		t.Fatalf("expected HasUnignoredResumeTask to return true for tracked file")
	}

	res, err := RemediateRepo(tmpDir, true)
	if err != nil {
		t.Fatalf("RemediateRepo failed: %v", err)
	}
	if !res.WasUntracked || !res.WasFileDeleted || !res.WasIgnored || !res.WasCommitted {
		t.Fatalf("expected full remediation (untracked=%v, deleted=%v, ignored=%v, committed=%v)",
			res.WasUntracked, res.WasFileDeleted, res.WasIgnored, res.WasCommitted)
	}

	ignoreData, err := os.ReadFile(filepath.Join(tmpDir, ".gitignore"))
	if err != nil {
		t.Fatalf("read .gitignore: %v", err)
	}
	if !strings.Contains(string(ignoreData), PrimaryIgnoreEntry) {
		t.Fatalf("expected .gitignore to contain %s, got: %s", PrimaryIgnoreEntry, string(ignoreData))
	}

	if HasUnignoredResumeTask(tmpDir) {
		t.Fatalf("expected HasUnignoredResumeTask to be false after remediation")
	}
}
