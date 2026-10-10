package cmdresolver

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExtractSlugFromURL(t *testing.T) {
	cases := []struct {
		url      string
		expected string
	}{
		{"https://github.com/my-org/auth-service.git", "auth-service"},
		{"git@github.com:my-org/billing-api.git", "billing-api"},
		{"https://gitlab.com/group/sub/payments", "payments"},
	}

	for _, c := range cases {
		slug := extractSlugFromURL(c.url)
		if slug != c.expected {
			t.Fatalf("expected slug %s, got %s for url %s", c.expected, slug, c.url)
		}
	}
}

func TestResolveRepoFeature_LocalNoGit(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "gitmap-repo-feature-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	res, errResolve := ResolveRepoFeature(tempDir)
	if errResolve != nil {
		t.Fatalf("ResolveRepoFeature failed: %v", errResolve)
	}

	if res.TargetType != TargetTypeLocalNoGit {
		t.Fatalf("expected target type %s, got %s", TargetTypeLocalNoGit, res.TargetType)
	}
	if !res.IsNewRepo {
		t.Fatalf("expected isNewRepo to be true")
	}

	gitDir := filepath.Join(tempDir, ".git")
	if _, errStat := os.Stat(gitDir); os.IsNotExist(errStat) {
		t.Fatalf("expected .git directory to be initialized")
	}
}

func TestResolveRepoFeature_LocalGit(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "gitmap-repo-feature-git-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Initialize git manually
	gitDir := filepath.Join(tempDir, ".git")
	_ = os.MkdirAll(gitDir, 0755)

	res, errResolve := ResolveRepoFeature(tempDir)
	if errResolve != nil {
		t.Fatalf("ResolveRepoFeature failed: %v", errResolve)
	}

	if res.TargetType != TargetTypeLocalGit {
		t.Fatalf("expected target type %s, got %s", TargetTypeLocalGit, res.TargetType)
	}
}
