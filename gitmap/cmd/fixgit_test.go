// Package cmd — fixgit_test.go: unit tests for gitmap fix-git.
package cmd

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestFixGit_NotAGitRepo(t *testing.T) {
	tempDir := t.TempDir()

	err := runFixGit([]string{tempDir})
	if err == nil {
		t.Fatal("expected error for non-git directory, got nil")
	}

	if !strings.Contains(err.Error(), "not a git repository") {
		t.Fatalf("unexpected error message: %v", err)
	}
}

func TestFixGit_RemediateLocks(t *testing.T) {
	tempDir := t.TempDir()
	gitDir := filepath.Join(tempDir, ".git")
	_ = os.MkdirAll(filepath.Join(gitDir, "refs", "heads"), 0755)

	lock1 := filepath.Join(gitDir, "index.lock")
	lock2 := filepath.Join(gitDir, "HEAD.lock")
	lock3 := filepath.Join(gitDir, "refs", "heads", "main.lock")

	_ = os.WriteFile(lock1, []byte("lock1"), 0644)
	_ = os.WriteFile(lock2, []byte("lock2"), 0644)
	_ = os.WriteFile(lock3, []byte("lock3"), 0644)

	opts := FixGitOptions{
		TargetDir:   tempDir,
		IsLocksOnly: true,
	}

	issues, err := remediateGitLocks(gitDir, opts)
	if err != nil {
		t.Fatalf("remediateGitLocks error: %v", err)
	}

	if len(issues) != 3 {
		t.Fatalf("expected 3 lock issues found, got %d", len(issues))
	}

	for _, lockFile := range []string{lock1, lock2, lock3} {
		if _, statErr := os.Stat(lockFile); !os.IsNotExist(statErr) {
			t.Errorf("lock file still exists: %s", lockFile)
		}
	}
}

func TestFixGit_DryRunLocks(t *testing.T) {
	tempDir := t.TempDir()
	gitDir := filepath.Join(tempDir, ".git")
	_ = os.MkdirAll(gitDir, 0755)

	lock := filepath.Join(gitDir, "index.lock")
	_ = os.WriteFile(lock, []byte("lock"), 0644)

	opts := FixGitOptions{
		TargetDir:   tempDir,
		IsLocksOnly: true,
		IsDryRun:    true,
	}

	issues, err := remediateGitLocks(gitDir, opts)
	if err != nil {
		t.Fatalf("remediateGitLocks error: %v", err)
	}

	if len(issues) != 1 {
		t.Fatalf("expected 1 lock issue, got %d", len(issues))
	}

	if issues[0].IsFixed {
		t.Error("expected issue.IsFixed to be false in dry-run mode")
	}

	if _, statErr := os.Stat(lock); os.IsNotExist(statErr) {
		t.Error("lock file was deleted during dry-run")
	}
}

func TestFixGit_CorruptZeroByteIndex(t *testing.T) {
	tempDir := t.TempDir()

	cmd := exec.Command("git", "init")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		t.Skipf("git not installed or init failed: %v", err)
	}

	testFile := filepath.Join(tempDir, "sample.txt")
	_ = os.WriteFile(testFile, []byte("hello gitmap"), 0644)

	_ = exec.Command("git", "-C", tempDir, "add", "sample.txt").Run()
	_ = exec.Command("git", "-C", tempDir, "commit", "-m", "init", "--author=Test <test@example.com>").Run()

	gitDir := filepath.Join(tempDir, ".git")
	indexPath := filepath.Join(gitDir, "index")

	// Corrupt index by truncating to 0 bytes
	_ = os.WriteFile(indexPath, []byte(""), 0644)

	opts := FixGitOptions{
		TargetDir:   tempDir,
		IsIndexOnly: true,
	}

	issues, err := remediateGitIndex(tempDir, gitDir, opts)
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

func TestFixGit_JSONOutput(t *testing.T) {
	res := FixGitResult{
		TargetDir:   "C:/repos/test",
		IsClean:     true,
		IssuesFound: 0,
		IssuesFixed: 0,
		Issues:      []FixGitIssue{},
	}

	var buf bytes.Buffer
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := renderFixGitOutput(res, true)
	_ = w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("renderFixGitOutput failed: %v", err)
	}

	_, _ = buf.ReadFrom(r)
	if !strings.Contains(buf.String(), `"isClean": true`) {
		t.Errorf("unexpected JSON output: %s", buf.String())
	}
}
