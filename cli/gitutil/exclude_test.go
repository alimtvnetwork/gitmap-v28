package gitutil

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEnsureGitMapExcluded_FreshRepo(t *testing.T) {
	tempDir := t.TempDir()
	gitDir := filepath.Join(tempDir, ".git")
	if err := os.MkdirAll(gitDir, 0o755); err != nil {
		t.Fatalf("failed to make .git: %v", err)
	}

	EnsureGitMapExcluded(tempDir)

	excludePath := filepath.Join(gitDir, "info", "exclude")
	data, err := os.ReadFile(excludePath)
	if err != nil {
		t.Fatalf("exclude file not created: %v", err)
	}

	if !strings.Contains(string(data), ".gitmap") {
		t.Fatalf("expected .gitmap in exclude, got %q", string(data))
	}
}

func TestEnsureGitMapExcluded_Idempotent(t *testing.T) {
	tempDir := t.TempDir()
	infoDir := filepath.Join(tempDir, ".git", "info")
	if err := os.MkdirAll(infoDir, 0o755); err != nil {
		t.Fatalf("failed to make info dir: %v", err)
	}

	excludePath := filepath.Join(infoDir, "exclude")
	if err := os.WriteFile(excludePath, []byte("# existing\n.gitmap/\n"), 0o644); err != nil {
		t.Fatalf("failed to write exclude: %v", err)
	}

	EnsureGitMapExcluded(tempDir)

	data, err := os.ReadFile(excludePath)
	if err != nil {
		t.Fatalf("failed to read exclude: %v", err)
	}

	if strings.Count(string(data), ".gitmap") != 1 {
		t.Fatalf("expected exactly one occurrence of .gitmap, got %q", string(data))
	}
}
