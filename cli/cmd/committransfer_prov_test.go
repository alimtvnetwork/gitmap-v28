package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestProvisionDestTarget_Existing(t *testing.T) {
	tempDir := t.TempDir()
	got := provisionDestTarget(tempDir, true)
	abs, _ := filepath.Abs(tempDir)
	if got != abs {
		t.Errorf("got %q, want %q", got, abs)
	}
}

func TestProvisionDestTarget_MissingLocal(t *testing.T) {
	tempDir := t.TempDir()
	targetPath := filepath.Join(tempDir, "auto-dest-repo")

	got := provisionDestTarget(targetPath, true)
	abs, _ := filepath.Abs(targetPath)
	if got != abs {
		t.Errorf("got %q, want %q", got, abs)
	}

	gitDir := filepath.Join(targetPath, ".git")
	if _, statErr := os.Stat(gitDir); statErr != nil {
		t.Errorf("expected .git directory to exist at %s: %v", gitDir, statErr)
	}
}
