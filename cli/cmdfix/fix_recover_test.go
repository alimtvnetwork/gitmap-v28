package cmdfix

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
)

func TestIsNoInitialCommitError(t *testing.T) {
	cases := []struct {
		output string
		want   bool
	}{
		{"You do not have the initial commit yet", true},
		{"fatal: your current branch 'main' does not have any commits yet", true},
		{"working tree clean", false},
		{"pathspec 'foo' did not match any files", false},
	}

	for _, tc := range cases {
		got := isNoInitialCommitError(tc.output)
		if got != tc.want {
			t.Errorf("isNoInitialCommitError(%q) = %v, want %v", tc.output, got, tc.want)
		}
	}
}

func TestIsBenignPullNoTracking(t *testing.T) {
	step := gitutil.RemediationStep{
		Name: "git",
		Args: []string{"-C", "/path/to/repo", "pull"},
	}

	cases := []struct {
		output string
		want   bool
	}{
		{"There is no tracking information for the current branch.", true},
		{"fatal: No remote repository specified. Please, specify either a URL or a remote name", true},
		{"fatal: No configured push destination.", true},
		{"fatal: Authentication failed", false},
	}

	for _, tc := range cases {
		got := isBenignPullNoTracking(step, tc.output)
		if got != tc.want {
			t.Errorf("isBenignPullNoTracking(%q) = %v, want %v", tc.output, got, tc.want)
		}
	}
}

func TestExtractRepoPath(t *testing.T) {
	argsWithC := []string{"-C", "/home/user/my-repo", "status"}
	if path := extractRepoPath(argsWithC); path != "/home/user/my-repo" {
		t.Errorf("expected /home/user/my-repo, got %q", path)
	}

	argsWithoutC := []string{"status", "--short"}
	if path := extractRepoPath(argsWithoutC); path != "" {
		t.Errorf("expected empty string, got %q", path)
	}
}

func TestSanitizeRecipeForRepo_NoInitialCommit(t *testing.T) {
	tempDir := t.TempDir()
	_ = exec.Command("git", "init", tempDir).Run()

	stashRecipe := gitutil.GenerateStashRecipe(tempDir)
	sanitized := sanitizeRecipeForRepo(tempDir, stashRecipe)

	if sanitized.Title == stashRecipe.Title {
		t.Errorf("expected stash recipe to be sanitized for repo with no initial commit")
	}
	if !gitutil.HasInitialCommit(tempDir) && sanitized.Title != "Option 1 (Initial Commit)" {
		t.Errorf("expected Option 1 (Initial Commit), got %q", sanitized.Title)
	}
}

func configureTestGitUser(repoPath string) {
	_ = exec.Command("git", "-C", repoPath, "config", "user.name", "Test User").Run()
	_ = exec.Command("git", "-C", repoPath, "config", "user.email", "test@example.com").Run()
}

func TestRecoverNoInitialCommit_CreatesCommitWhenDirty(t *testing.T) {
	tempDir := t.TempDir()
	_ = exec.Command("git", "init", tempDir).Run()
	configureTestGitUser(tempDir)

	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("hello"), 0o644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	step := gitutil.RemediationStep{Name: "git", Args: []string{"-C", tempDir, "stash", "-u"}}
	recovered, msg := recoverNoInitialCommit(step, "You do not have the initial commit yet")
	if !recovered || msg != "recovered (created initial commit)" {
		t.Fatalf("expected recoverNoInitialCommit to succeed, got %v, %q", recovered, msg)
	}
	if !gitutil.HasInitialCommit(tempDir) {
		t.Fatalf("expected repository to now have an initial commit")
	}
}

func TestRecoverNoInitialCommit_CleanRepo(t *testing.T) {
	tempDir := t.TempDir()
	_ = exec.Command("git", "init", tempDir).Run()
	configureTestGitUser(tempDir)

	step := gitutil.RemediationStep{Name: "git", Args: []string{"-C", tempDir, "stash", "-u"}}
	recovered, msg := recoverNoInitialCommit(step, "You do not have the initial commit yet")
	if !recovered || msg != "clean (no commits yet, working tree clean)" {
		t.Fatalf("expected recoverNoInitialCommit to succeed for clean repo, got %v, %q", recovered, msg)
	}
}
