// Package cmd — fixgit_unmerged_untracked_test.go: tests for unmerged and untracked repair.
package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestFixGit_DetectAndAbortUnmerged(t *testing.T) {
	tempDir := t.TempDir()

	cmd := exec.Command("git", "init")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		t.Skipf("git not installed: %v", err)
	}

	gitDir := filepath.Join(tempDir, ".git")
	mergeHead := filepath.Join(gitDir, "MERGE_HEAD")
	_ = os.WriteFile(mergeHead, []byte("fake-merge-sha"), 0644)

	opts := FixGitOptions{
		TargetDir: tempDir,
		IsDryRun:  true,
	}

	issues, err := remediateUnmergedConflict(tempDir, opts)
	if err != nil {
		t.Fatalf("remediateUnmergedConflict error: %v", err)
	}

	if len(issues) != 1 {
		t.Fatalf("expected 1 unmerged issue, got %d", len(issues))
	}

	if issues[0].IsFixed {
		t.Error("expected IsFixed to be false in dry run")
	}
}

func TestFixGit_UntrackedCollisionParser(t *testing.T) {
	sampleOutput := `error: The following untracked working tree files would be overwritten by merge:
	.agents/skills/bump-version-and-release/SKILL.md
	changelogs/v0.2.6.md
	frontend/src/utils/remote-url.ts
Please move or remove them before you merge.
Aborting`

	files := parseUntrackedConflictFiles(sampleOutput)
	if len(files) != 3 {
		t.Fatalf("expected 3 files parsed, got %d", len(files))
	}

	expected := ".agents/skills/bump-version-and-release/SKILL.md"
	if files[0] != expected {
		t.Errorf("expected %s, got %s", expected, files[0])
	}
}

func TestFixGit_BackupAndRemoveFile(t *testing.T) {
	tempDir := t.TempDir()
	src := filepath.Join(tempDir, "sample.txt")
	dst := filepath.Join(tempDir, "backup", "sample.txt")

	_ = os.WriteFile(src, []byte("untracked content"), 0644)

	err := copyAndRemoveFile(src, dst)
	if err != nil {
		t.Fatalf("copyAndRemoveFile failed: %v", err)
	}

	if _, statErr := os.Stat(src); !os.IsNotExist(statErr) {
		t.Error("src was not removed")
	}

	data, readErr := os.ReadFile(dst)
	if readErr != nil || string(data) != "untracked content" {
		t.Errorf("dst content invalid: %v, %s", readErr, string(data))
	}
}
