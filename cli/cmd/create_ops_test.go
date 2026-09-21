package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnsureOrProvisionDestinationRepo_Existing(t *testing.T) {
	tempDir := t.TempDir()
	got, err := EnsureOrProvisionDestinationRepo(tempDir, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	abs, _ := filepath.Abs(tempDir)
	if got != abs {
		t.Errorf("got %q, want %q", got, abs)
	}
}

func TestEnsureOrProvisionDestinationRepo_RemoteURL(t *testing.T) {
	url := "https://github.com/example/repo.git"
	got, err := EnsureOrProvisionDestinationRepo(url, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != url {
		t.Errorf("got %q, want %q", got, url)
	}
}

func TestEnsureOrProvisionDestinationRepo_MissingLocal(t *testing.T) {
	tempDir := t.TempDir()
	targetPath := filepath.Join(tempDir, "new-project-dir")

	got, err := EnsureOrProvisionDestinationRepo(targetPath, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	abs, _ := filepath.Abs(targetPath)
	if got != abs {
		t.Errorf("got %q, want %q", got, abs)
	}

	gitDir := filepath.Join(targetPath, ".git")
	if _, statErr := os.Stat(gitDir); statErr != nil {
		t.Errorf("expected .git directory to exist at %s: %v", gitDir, statErr)
	}
}
