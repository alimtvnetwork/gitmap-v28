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
	if !res.WasUntracked || !res.WasFileDeleted || !res.WasIgnored || !res.WasCommitted || !res.WasDeleteCommit || !res.WasIgnoreCommit {
		t.Fatalf("expected full two-step remediation (untracked=%v, deleted=%v, ignored=%v, deleteCommit=%v, ignoreCommit=%v, committed=%v)",
			res.WasUntracked, res.WasFileDeleted, res.WasIgnored, res.WasDeleteCommit, res.WasIgnoreCommit, res.WasCommitted)
	}

	ignoreData, err := os.ReadFile(filepath.Join(tmpDir, ".gitignore"))
	if err != nil {
		t.Fatalf("read .gitignore: %v", err)
	}
	if !strings.Contains(string(ignoreData), PrimaryIgnoreEntry) {
		t.Fatalf("expected .gitignore to contain %s, got: %s", PrimaryIgnoreEntry, string(ignoreData))
	}

	out, _ := exec.Command("git", "-C", tmpDir, "log", "-n", "3", "--oneline").Output()
	logStr := string(out)
	if !strings.Contains(logStr, "remove "+PrimaryIgnoreEntry) {
		t.Fatalf("expected deletion commit in git log, got: %s", logStr)
	}
	if !strings.Contains(logStr, "ignore "+PrimaryIgnoreEntry) {
		t.Fatalf("expected ignore commit in git log, got: %s", logStr)
	}

	if HasUnignoredResumeTask(tmpDir) {
		t.Fatalf("expected HasUnignoredResumeTask to be false after remediation")
	}
}

func TestRemediateRepo_UntrackedFileOnDisk(t *testing.T) {
	tmpDir := t.TempDir()
	if err := exec.Command("git", "init", tmpDir).Run(); err != nil {
		t.Skipf("git not available: %v", err)
	}
	_ = exec.Command("git", "-C", tmpDir, "config", "user.email", "test@example.com").Run()
	_ = exec.Command("git", "-C", tmpDir, "config", "user.name", "Test User").Run()

	dummy := filepath.Join(tmpDir, "README.md")
	_ = os.WriteFile(dummy, []byte("test"), 0o644)
	_ = exec.Command("git", "-C", tmpDir, "add", "README.md").Run()
	_ = exec.Command("git", "-C", tmpDir, "commit", "-m", "initial").Run()

	resumeFile := filepath.Join(tmpDir, PrimaryIgnoreEntry)
	_ = os.WriteFile(resumeFile, []byte(`{"task":"resume"}`), 0o644)

	res, err := RemediateRepo(tmpDir, true)
	if err != nil {
		t.Fatalf("RemediateRepo failed: %v", err)
	}
	if res.WasDeleteCommit {
		t.Fatalf("did not expect WasDeleteCommit for uncommitted file")
	}
	if !res.WasFileDeleted || !res.WasIgnored || !res.WasIgnoreCommit {
		t.Fatalf("expected file deletion and ignore commit (deleted=%v, ignored=%v, ignoreCommit=%v)",
			res.WasFileDeleted, res.WasIgnored, res.WasIgnoreCommit)
	}
}
