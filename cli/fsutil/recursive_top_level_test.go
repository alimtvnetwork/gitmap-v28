package fsutil

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRecursiveTopLevelScanner(t *testing.T) {
	tempDir := t.TempDir()

	// Top level repo
	topRepo := filepath.Join(tempDir, "top-repo")
	os.MkdirAll(filepath.Join(topRepo, ".git"), 0755)

	// Nested sub-repo inside topRepo (should be pruned)
	nestedRepo := filepath.Join(topRepo, "submodules", "nested-repo")
	os.MkdirAll(filepath.Join(nestedRepo, ".git"), 0755)

	// Another top level repo
	secondRepo := filepath.Join(tempDir, "second-repo")
	os.MkdirAll(filepath.Join(secondRepo, ".git"), 0755)

	discovered, err := DiscoverTopLevelGitRepos(tempDir)
	if err != nil {
		t.Fatalf("DiscoverTopLevelGitRepos failed: %v", err)
	}

	if len(discovered) != 2 {
		t.Fatalf("expected 2 top-level repos, got %d (%v)", len(discovered), discovered)
	}

	// Test IsInsideWorkDir
	if !IsInsideWorkDir(topRepo, tempDir) {
		t.Fatal("expected topRepo to be inside tempDir")
	}

	// Test IsValidGitDir
	if !IsValidGitDir(topRepo) {
		t.Fatal("expected topRepo to be a valid git dir")
	}
}
