package cmdhistory

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// runTestGit runs git in dir and fails the test on error.
func runTestGit(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_NAME=test", "GIT_AUTHOR_EMAIL=test@example.com",
		"GIT_COMMITTER_NAME=test", "GIT_COMMITTER_EMAIL=test@example.com")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("git %v: %v", args, err)
	}

	return strings.TrimSpace(string(out))
}

func TestPurgeHelpNamesRangeAndDryRun(t *testing.T) {
	if !strings.Contains(purgeHelpText, "--dry-run") {
		t.Fatal("help missing dry-run")
	}
	if !strings.Contains(purgeHelpText, "--range") {
		t.Fatal("help missing range")
	}
	if !strings.Contains(purgeHelpText, ".cursor/skills/execute-parent-task-with-n-steps-v7") {
		t.Fatal("help missing example path")
	}
}

func TestExpandCommitRangeOneCommit(t *testing.T) {
	// WHY temp repo (31-Oct-26): the old test pinned a historical
	// SHA from this repo, which does not exist in CI's shallow
	// checkout (fetch-depth 1), so rev-list failed there while
	// passing locally. A two-commit temp repo tests the same
	// one-commit range without depending on repo history.
	dir := t.TempDir()
	runTestGit(t, dir, "init", "-q")
	runTestGit(t, dir, "config", "user.name", "test")
	runTestGit(t, dir, "config", "user.email", "test@example.com")
	if err := os.WriteFile(filepath.Join(dir, "a.txt"), []byte("one\n"), 0644); err != nil {
		t.Fatal(err)
	}
	runTestGit(t, dir, "add", ".")
	runTestGit(t, dir, "commit", "-q", "-m", "first")
	if err := os.WriteFile(filepath.Join(dir, "a.txt"), []byte("two\n"), 0644); err != nil {
		t.Fatal(err)
	}
	runTestGit(t, dir, "add", ".")
	runTestGit(t, dir, "commit", "-q", "-m", "second")
	head := runTestGit(t, dir, "rev-parse", "HEAD")

	shas, err := expandCommitRange(dir, "HEAD^..HEAD")
	if err != nil {
		t.Fatal(err)
	}
	if len(shas) != 1 {
		t.Fatalf("count = %d", len(shas))
	}
	if shas[0] != head {
		t.Fatalf("sha = %s, want %s", shas[0], head)
	}
}

func TestExpandCommitRangeRejectsSha(t *testing.T) {
	_, err := expandCommitRange(".", "abc123")
	if err == nil {
		t.Fatal("expected an error")
	}
}
