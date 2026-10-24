package cmdfix

import (
	"os/exec"
	"testing"
)

func TestResolveSavedGitIdentity_FromGitConfig(t *testing.T) {
	tempDir := t.TempDir()
	_ = exec.Command("git", "init", tempDir).Run()
	_ = exec.Command("git", "-C", tempDir, "config", "user.name", "Saved User").Run()
	_ = exec.Command("git", "-C", tempDir, "config", "user.email", "saved@example.com").Run()

	ident := resolveSavedGitIdentity(tempDir)
	if !ident.IsConfigured {
		t.Fatalf("expected identity to be configured from git config")
	}
	if ident.Name != "Saved User" {
		t.Errorf("expected 'Saved User', got %q", ident.Name)
	}
	if ident.Email != "saved@example.com" {
		t.Errorf("expected 'saved@example.com', got %q", ident.Email)
	}
}

func TestResolveSavedGitIdentity_FromEnv(t *testing.T) {
	tempDir := t.TempDir()
	_ = exec.Command("git", "init", tempDir).Run()
	// Ensure repo has no local user config
	_ = exec.Command("git", "-C", tempDir, "config", "--unset", "user.name").Run()
	_ = exec.Command("git", "-C", tempDir, "config", "--unset", "user.email").Run()

	t.Setenv("GIT_AUTHOR_NAME", "Env Author")
	t.Setenv("GIT_AUTHOR_EMAIL", "env@author.org")

	ident, hasIdent := resolveFromEnv()
	if !hasIdent || !ident.IsConfigured {
		t.Fatalf("expected identity from env")
	}
	if ident.Name != "Env Author" || ident.Email != "env@author.org" {
		t.Errorf("unexpected env identity: %+v", ident)
	}
}

func TestBuildCommitArgsWithIdentity_ConfiguredRepo(t *testing.T) {
	tempDir := t.TempDir()
	_ = exec.Command("git", "init", tempDir).Run()
	_ = exec.Command("git", "-C", tempDir, "config", "user.name", "Repo User").Run()
	_ = exec.Command("git", "-C", tempDir, "config", "user.email", "repo@example.com").Run()

	args, hasIdentity := buildCommitArgsWithIdentity(tempDir, "test commit")
	if !hasIdentity {
		t.Fatalf("expected hasIdentity to be true")
	}
	if len(args) == 0 {
		t.Fatalf("expected non-empty commit args")
	}
	// When already configured in repo, standard commit args are used
	if args[0] != "-C" || args[1] != tempDir || args[2] != "commit" {
		t.Errorf("unexpected commit args: %v", args)
	}
}
