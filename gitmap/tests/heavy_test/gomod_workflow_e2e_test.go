package heavy_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cmd"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func initTestRepo(t *testing.T) (string, func()) {
	t.Helper()

	dir := t.TempDir()
	t.Chdir(dir)

	run := func(args ...string) {
		t.Helper()
		c := exec.Command(args[0], args[1:]...)
		c.Dir = dir
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		if err := c.Run(); err != nil {
			t.Fatalf("command %v failed: %v", args, err)
		}
	}

	run("git", "init")
	run("git", "config", "user.email", "test@test.com")
	run("git", "config", "user.name", "Test")

	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("# test\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	run("git", "add", "-A")
	run("git", "commit", "-m", "initial")

	return dir, func() {}
}

func TestCreateGoModBranches_Integration(t *testing.T) {
	_, cleanup := initTestRepo(t)
	defer cleanup()

	slug := "github-com-org-repo"
	backup, feature := cmd.CreateGoModBranches(slug)

	expectedBackup := constants.GoModBackupPrefix + slug
	expectedFeature := constants.GoModFeaturePrefix + slug

	if backup != expectedBackup {
		t.Errorf("backup = %q, want %q", backup, expectedBackup)
	}

	if feature != expectedFeature {
		t.Errorf("feature = %q, want %q", feature, expectedFeature)
	}

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

	head, _ := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	current := strings.TrimSpace(string(head))
	if current != expectedFeature {
		t.Errorf("expected to be on %q, got %q", expectedFeature, current)
	}
}

func TestGoModCurrentBranch_Integration(t *testing.T) {
	_, cleanup := initTestRepo(t)
	defer cleanup()

	branch := cmd.GoModCurrentBranch()
	if branch != "main" && branch != "master" {
		t.Errorf("expected main or master, got %q", branch)
	}
}

func TestIsWorkTreeDirty_Clean(t *testing.T) {
	_, cleanup := initTestRepo(t)
	defer cleanup()

	if cmd.IsWorkTreeDirty() {
		t.Error("expected clean work tree after initial commit")
	}
}

func TestIsWorkTreeDirty_Dirty(t *testing.T) {
	_, cleanup := initTestRepo(t)
	defer cleanup()

	_ = os.WriteFile("dirty.txt", []byte("uncommitted"), 0o644)

	if !cmd.IsWorkTreeDirty() {
		t.Error("expected dirty work tree after creating untracked file")
	}
}

func TestCommitGoModChanges_Integration(t *testing.T) {
	_, cleanup := initTestRepo(t)
	defer cleanup()

	_ = os.WriteFile("go.mod", []byte("module github.com/new/path\n"), 0o644)

	cmd.CommitGoModChanges("github.com/old/path", "github.com/new/path", 3)

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

	_ = os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module github.com/old/mod\n\ngo 1.21\n"), 0o644)
	_ = os.MkdirAll(filepath.Join(dir, "pkg"), 0o755)
	_ = os.WriteFile(filepath.Join(dir, "pkg", "main.go"), []byte("package main\n\nimport \"github.com/old/mod/pkg\"\n"), 0o644)

	_ = cmd.ReplaceInGoMod("github.com/old/mod", "github.com/new/mod")
	count := cmd.ReplaceModulePath("github.com/old/mod", "github.com/new/mod", false, nil)

	gomod, _ := os.ReadFile(filepath.Join(dir, "go.mod"))
	if !strings.Contains(string(gomod), "github.com/new/mod") {
		t.Error("go.mod not updated")
	}

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
