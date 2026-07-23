package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// initTestRepo creates a temporary git repo with an initial commit and returns
// the directory path. The caller's working directory is changed to the repo
// via t.Chdir, which Go restores after the test (in the correct order relative
// to t.TempDir cleanup, preventing leftover cwds pointing at deleted dirs that
// poison sibling tests via os.Getwd → "/" fallbacks on macOS).
func initTestRepo(t *testing.T) (string, func()) {
	t.Helper()

	dir := t.TempDir()
	t.Chdir(dir)

	run := func(args ...string) {
		t.Helper()
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = dir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			t.Fatalf("command %v failed: %v", args, err)
		}
	}

	run("git", "init")
	run("git", "config", "user.email", "test@test.com")
	run("git", "config", "user.name", "Test")

	// Write a dummy file and commit so HEAD exists.
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("# test\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	run("git", "add", "-A")
	run("git", "commit", "-m", "initial")

	return dir, func() {}
}

func TestDeriveSlug_Integration(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"github.com/org/repo", "github-com-org-repo"},
		{"go.example.dev/pkg", "go-example-dev-pkg"},
		{"github.com/a/b@v2", "github-com-a-b-v2"},
		{"simple", "simple"},
	}

	for _, tt := range tests {
		got := deriveSlug(tt.input)
		if got != tt.expected {
			t.Errorf("deriveSlug(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestCreateGoModBranches_Integration(t *testing.T) {
	_, cleanup := initTestRepo(t)
	defer cleanup()

	slug := "github-com-org-repo"
	backup, feature := createGoModBranches(slug)

	expectedBackup := constants.GoModBackupPrefix + slug
	expectedFeature := constants.GoModFeaturePrefix + slug

	if backup != expectedBackup {
		t.Errorf("backup = %q, want %q", backup, expectedBackup)
	}
	if feature != expectedFeature {
		t.Errorf("feature = %q, want %q", feature, expectedFeature)
	}

	// Verify branches exist in git.
	out, err := exec.Command("git", "branch", "--list").Output()
	if err != nil {
		t.Fatalf("git branch --list failed: %v", err)
	}
	branches := string(out)
	if !strings.Contains(branches, expectedBackup) {
		t.Errorf("backup branch %q not found in:\n%s", expectedBackup, branches)
	}
	if !strings.Contains(branches, expectedFeature) {
		t.Errorf("feature branch %q not found in:\n%s", expectedFeature, branches)
	}

	// Verify we are on the feature branch.
	head, _ := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	current := strings.TrimSpace(string(head))
	if current != expectedFeature {
		t.Errorf("expected to be on %q, got %q", expectedFeature, current)
	}
}

func TestGoModCurrentBranch_Integration(t *testing.T) {
	_, cleanup := initTestRepo(t)
	defer cleanup()

	branch := goModCurrentBranch()
	// git init defaults to "master" or "main" depending on config.
	if branch != "main" && branch != "master" {
		t.Errorf("expected main or master, got %q", branch)
	}
}

func TestIsWorkTreeDirty_Clean(t *testing.T) {
	_, cleanup := initTestRepo(t)
	defer cleanup()

	if isWorkTreeDirty() {
		t.Error("expected clean work tree after initial commit")
	}
}

func TestIsWorkTreeDirty_Dirty(t *testing.T) {
	_, cleanup := initTestRepo(t)
	defer cleanup()

	os.WriteFile("dirty.txt", []byte("uncommitted"), 0o644)

	if !isWorkTreeDirty() {
		t.Error("expected dirty work tree after creating untracked file")
	}
}

func TestCommitGoModChanges_Integration(t *testing.T) {
	_, cleanup := initTestRepo(t)
	defer cleanup()

	// Create a change to commit.
	os.WriteFile("go.mod", []byte("module github.com/new/path\n"), 0o644)

	commitGoModChanges("github.com/old/path", "github.com/new/path", 3)

	// Verify the commit message.
	out, err := exec.Command("git", "log", "-1", "--pretty=%s").Output()
	if err != nil {
		t.Fatalf("git log failed: %v", err)
	}
	subject := strings.TrimSpace(string(out))
	if subject != "refactor: rename go module path" {
		t.Errorf("unexpected commit subject: %q", subject)
	}
}

func TestReplaceModulePath_Integration(t *testing.T) {
	dir, cleanup := initTestRepo(t)
	defer cleanup()

	// Set up go.mod and a .go file with old path.
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module github.com/old/mod\n\ngo 1.21\n"), 0o644)
	os.MkdirAll(filepath.Join(dir, "pkg"), 0o755)
	os.WriteFile(filepath.Join(dir, "pkg", "main.go"), []byte("package main\n\nimport \"github.com/old/mod/pkg\"\n"), 0o644)

	replaceInGoMod("github.com/old/mod", "github.com/new/mod")
	count := replaceModulePath("github.com/old/mod", "github.com/new/mod", false, nil)

	// go.mod should be updated.
	gomod, _ := os.ReadFile(filepath.Join(dir, "go.mod"))
	if !strings.Contains(string(gomod), "github.com/new/mod") {
		t.Error("go.mod not updated")
	}

	// .go file should be updated.
	gofile, _ := os.ReadFile(filepath.Join(dir, "pkg", "main.go"))
	if !strings.Contains(string(gofile), "github.com/new/mod") {
		t.Error("main.go not updated")
	}
	if strings.Contains(string(gofile), "github.com/old/mod") {
		t.Error("main.go still contains old path")
	}

	if count != 1 {
		t.Errorf("expected 1 file replaced, got %d", count)
	}
}
