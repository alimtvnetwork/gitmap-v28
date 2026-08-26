package fsutil

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDiscoverChildGitRepos(t *testing.T) {
	tempDir := t.TempDir()

	repo1 := filepath.Join(tempDir, "repo1")
	os.MkdirAll(filepath.Join(repo1, ".git"), 0755)

	repo2 := filepath.Join(tempDir, "repo2")
	os.MkdirAll(filepath.Join(repo2, ".git"), 0755)

	notRepo := filepath.Join(tempDir, "notRepo")
	os.MkdirAll(notRepo, 0755)

	repos, err := DiscoverChildGitRepos(tempDir)
	if err != nil {
		t.Fatalf("DiscoverChildGitRepos failed: %v", err)
	}

	if len(repos) != 2 {
		t.Fatalf("expected 2 child repos, got %d", len(repos))
	}
}
