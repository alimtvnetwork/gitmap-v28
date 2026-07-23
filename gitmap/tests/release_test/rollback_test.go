package release_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/release"
)

// initRollbackRepo creates a temporary git repo with an initial commit,
// changes CWD into it, and returns a cleanup function.
func initRollbackRepo(t *testing.T) (string, func()) {
	t.Helper()

	dir := t.TempDir()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	git := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			t.Fatalf("git %v failed: %v", args, err)
		}
	}

	git("init")
	git("config", "user.email", "test@test.com")
	git("config", "user.name", "Test")
	git("checkout", "-b", "main")

	err = os.WriteFile(filepath.Join(dir, "README.md"), []byte("# test\n"), 0o644)
	if err != nil {
		t.Fatal(err)
	}

	git("add", "-A")
	git("commit", "-m", "initial")

	return dir, func() { os.Chdir(orig) }
}

// branchExists checks if a local git branch exists.
func branchExists(t *testing.T, branch string) bool {
	t.Helper()
	cmd := exec.Command("git", "branch", "--list", branch)
	out, err := cmd.Output()
	if err != nil {
		return false
	}

	return len(strings.TrimSpace(string(out))) > 0
}

// tagExists checks if a local git tag exists.
func tagExists(t *testing.T, tag string) bool {
	t.Helper()
	cmd := exec.Command("git", "tag", "--list", tag)
	out, err := cmd.Output()
	if err != nil {
		return false
	}

	return len(strings.TrimSpace(string(out))) > 0
}

// currentBranch returns the current branch name.
func currentBranch(t *testing.T) string {
	t.Helper()
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("rev-parse: %v", err)
	}

	return strings.TrimSpace(string(out))
}

// TestRollback_DeletesBranchAndTag verifies that Rollback removes the
// local release branch and tag after a simulated push failure.
func TestRollback_DeletesBranchAndTag(t *testing.T) {
	_, cleanup := initRollbackRepo(t)
	defer cleanup()

	branchName := constants.ReleaseBranchPrefix + "v10.0.0"
	tag := "v10.0.0"

	// Create branch and tag (simulating steps 5–6 of the release workflow).
	err := release.CreateBranch(branchName, "HEAD")
	if err != nil {
		t.Fatalf("CreateBranch: %v", err)
	}
	err = release.CreateTag(tag, "Release "+tag)
	if err != nil {
		t.Fatalf("CreateTag: %v", err)
	}

	// Confirm both exist before rollback.
	if !branchExists(t, branchName) {
		t.Fatal("branch should exist before rollback")
	}
	if !tagExists(t, tag) {
		t.Fatal("tag should exist before rollback")
	}

	// Simulate push failure → trigger rollback.
	release.Rollback(branchName, tag, "main")

	// Verify branch and tag are deleted.
	if branchExists(t, branchName) {
		t.Error("branch should be deleted after rollback")
	}
	if tagExists(t, tag) {
		t.Error("tag should be deleted after rollback")
	}

	// Verify we're back on the original branch.
	if got := currentBranch(t); got != "main" {
		t.Errorf("expected branch main, got %s", got)
	}
}

// TestRollback_BranchOnlyWhenTagEmpty verifies that Rollback handles
// an empty tag string gracefully (only deletes the branch).
func TestRollback_BranchOnlyWhenTagEmpty(t *testing.T) {
	_, cleanup := initRollbackRepo(t)
	defer cleanup()

	branchName := constants.ReleaseBranchPrefix + "v10.1.0"

	err := release.CreateBranch(branchName, "HEAD")
	if err != nil {
		t.Fatalf("CreateBranch: %v", err)
	}

	release.Rollback(branchName, "", "main")

	if branchExists(t, branchName) {
		t.Error("branch should be deleted after rollback")
	}
	if got := currentBranch(t); got != "main" {
		t.Errorf("expected branch main, got %s", got)
	}
}

// TestRollback_TagOnlyWhenBranchEmpty verifies that Rollback handles
// an empty branch string gracefully (only deletes the tag).
func TestRollback_TagOnlyWhenBranchEmpty(t *testing.T) {
	_, cleanup := initRollbackRepo(t)
	defer cleanup()

	tag := "v10.2.0"
	err := release.CreateTag(tag, "Release "+tag)
	if err != nil {
		t.Fatalf("CreateTag: %v", err)
	}

	release.Rollback("", tag, "main")

	if tagExists(t, tag) {
		t.Error("tag should be deleted after rollback")
	}
}

// TestRollback_SwitchBackToOriginalBranch verifies that Rollback
// returns to the original branch even when starting from the release branch.
func TestRollback_SwitchBackToOriginalBranch(t *testing.T) {
	_, cleanup := initRollbackRepo(t)
	defer cleanup()

	branchName := constants.ReleaseBranchPrefix + "v10.3.0"
	tag := "v10.3.0"

	err := release.CreateBranch(branchName, "HEAD")
	if err != nil {
		t.Fatalf("CreateBranch: %v", err)
	}
	err = release.CheckoutBranch(branchName)
	if err != nil {
		t.Fatalf("CheckoutBranch: %v", err)
	}
	err = release.CreateTag(tag, "Release "+tag)
	if err != nil {
		t.Fatalf("CreateTag: %v", err)
	}

	// We're on the release branch — rollback should switch us back.
	if got := currentBranch(t); got != branchName {
		t.Fatalf("expected to be on %s, got %s", branchName, got)
	}

	release.Rollback(branchName, tag, "main")

	if got := currentBranch(t); got != "main" {
		t.Errorf("expected branch main after rollback, got %s", got)
	}
	if branchExists(t, branchName) {
		t.Error("branch should be deleted after rollback")
	}
	if tagExists(t, tag) {
		t.Error("tag should be deleted after rollback")
	}
}

// TestRollback_NoOpWhenAllEmpty verifies that Rollback handles all
// empty arguments without panicking.
func TestRollback_NoOpWhenAllEmpty(t *testing.T) {
	_, cleanup := initRollbackRepo(t)
	defer cleanup()

	// Should not panic or error.
	release.Rollback("", "", "")

	// Still on main.
	if got := currentBranch(t); got != "main" {
		t.Errorf("expected branch main, got %s", got)
	}
}
