package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestRecreateRepo_NotInGitRepo(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "gitmap-not-git-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	origWd, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	defer os.Chdir(origWd)

	err = runRecreateRepo([]string{"new-repo-name"})
	if err == nil {
		t.Fatalf("expected error when not inside a git repo, got nil")
	}
	if !strings.Contains(err.Error(), "not a git repository") {
		t.Fatalf("unexpected error message: %v", err)
	}
}

func TestRecreateRepo_MissingArgs(t *testing.T) {
	err := runRecreateRepo([]string{})
	if err == nil {
		t.Fatalf("expected error when no args provided, got nil")
	}
	if !strings.Contains(err.Error(), "usage: gitmap recreate-repo") {
		t.Fatalf("unexpected error message: %v", err)
	}
}

func TestRecreateRepo_BackupBranchCreation(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "gitmap-recreate-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Init git and make an initial commit
	cmdInit := exec.Command("git", "init", tempDir)
	if err := cmdInit.Run(); err != nil {
		t.Fatalf("git init failed: %v", err)
	}

	testFile := filepath.Join(tempDir, "file.txt")
	_ = os.WriteFile(testFile, []byte("hello"), 0644)

	cmdAdd := exec.Command("git", "-C", tempDir, "add", ".")
	_ = cmdAdd.Run()

	cmdCommit := exec.Command("git", "-C", tempDir, "commit", "-m", "initial test commit")
	_ = cmdCommit.Run()

	origWd, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	defer os.Chdir(origWd)

	// We test that backup branch is created even if gh is mocked/fails
	_ = runRecreateRepo([]string{"dummy-test-repo"})

	cmdBranch := exec.Command("git", "branch")
	out, bErr := cmdBranch.Output()
	if bErr != nil {
		t.Fatalf("git branch failed: %v", bErr)
	}

	if !strings.Contains(string(out), "backup/recreate-") {
		t.Errorf("expected backup/recreate- branch to be created, got: %s", string(out))
	}
}
