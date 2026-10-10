package cmdfix

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/cloner"
)

func TestClassifyNonRepoFolder_SpecialInfrastructureRepo(t *testing.T) {
	tmpDir := t.TempDir()
	specialRepoDir := filepath.Join(tmpDir, "repo-cache")
	if err := os.MkdirAll(specialRepoDir, 0o755); err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	diag := cloner.ClassifyNonRepoFolder(specialRepoDir, "repo-cache")
	if !diag.IsNonRepoFolder {
		t.Fatalf("expected IsNonRepoFolder=true, got false")
	}
	if !diag.IsSpecialRepo {
		t.Fatalf("expected IsSpecialRepo=true, got false")
	}
	if !strings.Contains(diag.Reason, "known infrastructure repository: repo-cache") {
		t.Fatalf("unexpected reason: %q", diag.Reason)
	}
	if !strings.HasPrefix(diag.Option1, "gitmap clone") {
		t.Fatalf("expected Option1 to start with 'gitmap clone', got %q", diag.Option1)
	}
	if !strings.Contains(diag.Option2, "init") {
		t.Fatalf("expected Option2 to contain 'init', got %q", diag.Option2)
	}
}

func TestClassifyNonRepoFolder_StandardRepoNoRemote(t *testing.T) {
	tmpDir := t.TempDir()
	standardDir := filepath.Join(tmpDir, "my-local-notes")
	if err := os.MkdirAll(standardDir, 0o755); err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	diag := cloner.ClassifyNonRepoFolder(standardDir, "my-local-notes")
	if !diag.IsNonRepoFolder {
		t.Fatalf("expected IsNonRepoFolder=true, got false")
	}
	if diag.IsSpecialRepo {
		t.Fatalf("expected IsSpecialRepo=false, got true")
	}
	if !strings.Contains(diag.Option1, "init") {
		t.Fatalf("expected Option1 to contain 'init', got %q", diag.Option1)
	}
	if !strings.Contains(diag.Option2, "gitmap rm") {
		t.Fatalf("expected Option2 to contain 'gitmap rm', got %q", diag.Option2)
	}
}

func TestClassifyNonRepoFolder_ValidGitRepoIgnored(t *testing.T) {
	tmpDir := t.TempDir()
	gitDir := filepath.Join(tmpDir, "valid-repo", ".git")
	if err := os.MkdirAll(gitDir, 0o755); err != nil {
		t.Fatalf("failed to create temp git dir: %v", err)
	}

	diag := cloner.ClassifyNonRepoFolder(filepath.Dir(gitDir), "valid-repo")
	if diag.IsNonRepoFolder {
		t.Fatalf("expected IsNonRepoFolder=false for valid git repo, got true")
	}
}

func TestClassifyNonRepoFolder_MissingDirIgnored(t *testing.T) {
	tmpDir := t.TempDir()
	nonExistent := filepath.Join(tmpDir, "does-not-exist")

	diag := cloner.ClassifyNonRepoFolder(nonExistent, "does-not-exist")
	if diag.IsNonRepoFolder {
		t.Fatalf("expected IsNonRepoFolder=false for non-existent path, got true")
	}
}

func TestResolveNonRepoRemediationItem(t *testing.T) {
	tmpDir := t.TempDir()
	repoPath := filepath.Join(tmpDir, "rc-folder")
	if err := os.MkdirAll(repoPath, 0o755); err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	diag := cloner.NonRepoDiagnosis{
		IsNonRepoFolder: true,
		IsSpecialRepo:   true,
		RepoName:        "repo-cache",
		Path:            repoPath,
		RemoteURL:       "https://github.com/example/repo-cache.git",
		HasRemote:       true,
		Reason:          "directory exists but is not a Git repository",
		Option1:         "gitmap clone repo-cache",
		Option2:         "git -C \"" + repoPath + "\" init",
	}

	item, err := buildNonRepoRemediationItem(diag, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item == nil {
		t.Fatalf("expected non-nil item")
	}
	if len(item.Recipes) < 2 {
		t.Fatalf("expected at least 2 recipes, got %d", len(item.Recipes))
	}
	if item.Recipes[0].Title != "Clone from Remote" {
		t.Fatalf("expected first recipe 'Clone from Remote', got %q", item.Recipes[0].Title)
	}
}

func TestRunFixDirect_NonRepoInitResolution(t *testing.T) {
	tmpDir := t.TempDir()
	targetDir := filepath.Join(tmpDir, "test-non-repo")
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	// Run fix directly targeting targetDir with "init" action
	err := runFixDirect([]string{targetDir, "init"}, "")
	if err != nil {
		t.Fatalf("runFixDirect failed: %v", err)
	}

	// Verify that git init was executed and .git directory exists now
	gitDirPath := filepath.Join(targetDir, ".git")
	info, statErr := os.Stat(gitDirPath)
	if statErr != nil || !info.IsDir() {
		t.Fatalf("expected .git directory to be initialized at %s, but stat returned: %v", gitDirPath, statErr)
	}
}
