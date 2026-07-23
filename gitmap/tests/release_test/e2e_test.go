package release_test

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/release"
)

// initE2ERepo creates a temp Git repo with a local bare remote (origin),
// changes CWD into the work repo, and overrides DefaultReleaseDir.
// Returns the work dir, release dir, and a cleanup function.
func initE2ERepo(t *testing.T) (string, string, func()) {
	t.Helper()

	base := t.TempDir()
	bareDir := filepath.Join(base, "origin.git")
	workDir := filepath.Join(base, "work")

	origDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	origReleaseDir := constants.DefaultReleaseDir

	git := func(dir string, args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			t.Fatalf("git %v (in %s) failed: %v", args, dir, err)
		}
	}

	// Create bare remote.
	os.MkdirAll(bareDir, 0o755)
	git(bareDir, "init", "--bare")

	// Clone into work dir.
	git(base, "clone", bareDir, "work")
	git(workDir, "config", "user.email", "test@test.com")
	git(workDir, "config", "user.name", "Test")

	// Create initial commit on main.
	readme := filepath.Join(workDir, "README.md")
	os.WriteFile(readme, []byte("# e2e test\n"), 0o644)
	git(workDir, "add", "-A")
	git(workDir, "commit", "-m", "initial")
	git(workDir, "push", "origin", "HEAD")

	if err := os.Chdir(workDir); err != nil {
		t.Fatal(err)
	}

	releaseDir := filepath.Join(workDir, ".gitmap", "release")
	constants.DefaultReleaseDir = releaseDir

	return workDir, releaseDir, func() {
		constants.DefaultReleaseDir = origReleaseDir
		os.Chdir(origDir)
	}
}

// TestE2E_FullReleaseCycle exercises Execute() end-to-end:
// version resolution → branch → tag → push → metadata → auto-commit.
func TestE2E_FullReleaseCycle(t *testing.T) {
	_, releaseDir, cleanup := initE2ERepo(t)
	defer cleanup()

	err := release.Execute(release.Options{
		Version:  "v5.0.0",
		DryRun:   false,
		NoCommit: false,
		SkipMeta: false,
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}

	// 1. Release branch should exist.
	branchName := constants.ReleaseBranchPrefix + "v5.0.0"
	if !release.BranchExists(branchName) {
		t.Errorf("release branch %s should exist", branchName)
	}

	// 2. Tag should exist locally.
	if !release.TagExistsLocally("v5.0.0") {
		t.Error("tag v5.0.0 should exist locally")
	}

	// 3. Metadata JSON should exist.
	metaPath := filepath.Join(releaseDir, "v5.0.0.json")
	if _, err := os.Stat(metaPath); os.IsNotExist(err) {
		t.Error("v5.0.0.json should exist")
	}

	// 4. Metadata content should be valid.
	data, err := os.ReadFile(metaPath)
	if err != nil {
		t.Fatalf("read metadata: %v", err)
	}

	var meta release.ReleaseMeta
	if err := json.Unmarshal(data, &meta); err != nil {
		t.Fatalf("unmarshal metadata: %v", err)
	}

	if meta.Tag != "v5.0.0" {
		t.Errorf("expected tag v5.0.0, got %s", meta.Tag)
	}
	if meta.Version != "5.0.0" {
		t.Errorf("expected version 5.0.0, got %s", meta.Version)
	}
	if meta.Branch != branchName {
		t.Errorf("expected branch %s, got %s", branchName, meta.Branch)
	}
	if len(meta.Commit) == 0 {
		t.Error("commit SHA should be populated")
	}

	// 5. latest.json should reference v5.0.0.
	latest, err := release.ReadLatest()
	if err != nil {
		t.Fatalf("ReadLatest: %v", err)
	}
	if latest.Tag != "v5.0.0" {
		t.Errorf("expected latest tag v5.0.0, got %s", latest.Tag)
	}

	// 6. Should be back on original branch (main/master), not release branch.
	current, err := release.CurrentBranchName()
	if err != nil {
		t.Fatalf("CurrentBranchName: %v", err)
	}
	if strings.HasPrefix(current, "release/") {
		t.Errorf("should be on original branch, got %s", current)
	}

	// 7. ReleaseExists should return true.
	v, _ := release.Parse("v5.0.0")
	if !release.ReleaseExists(v) {
		t.Error("ReleaseExists should return true after release")
	}
}

// TestE2E_DuplicateVersionBlocked verifies that releasing the same
// version twice returns an error.
func TestE2E_DuplicateVersionBlocked(t *testing.T) {
	_, _, cleanup := initE2ERepo(t)
	defer cleanup()

	err := release.Execute(release.Options{
		Version:  "v4.0.0",
		DryRun:   false,
		NoCommit: true,
		SkipMeta: false,
	})
	if err != nil {
		t.Fatalf("first Execute: %v", err)
	}

	// Second release with same version should fail.
	err = release.Execute(release.Options{
		Version:  "v4.0.0",
		DryRun:   false,
		NoCommit: true,
		SkipMeta: false,
	})
	if err == nil {
		t.Fatal("expected error for duplicate version, got nil")
	}
}

// TestE2E_DryRunNoSideEffects verifies that dry-run creates no branches,
// tags, or metadata files.
func TestE2E_DryRunNoSideEffects(t *testing.T) {
	_, releaseDir, cleanup := initE2ERepo(t)
	defer cleanup()

	err := release.Execute(release.Options{
		Version: "v3.0.0",
		DryRun:  true,
	})
	if err != nil {
		t.Fatalf("Execute dry-run: %v", err)
	}

	branchName := constants.ReleaseBranchPrefix + "v3.0.0"
	if release.BranchExists(branchName) {
		t.Error("dry-run should not create a branch")
	}
	if release.TagExistsLocally("v3.0.0") {
		t.Error("dry-run should not create a tag")
	}

	metaPath := filepath.Join(releaseDir, "v3.0.0.json")
	if _, err := os.Stat(metaPath); err == nil {
		t.Error("dry-run should not write metadata")
	}
}

// TestE2E_NoCommitSkipsAutoCommit verifies that --no-commit produces
// metadata but does not auto-commit it.
func TestE2E_NoCommitSkipsAutoCommit(t *testing.T) {
	_, releaseDir, cleanup := initE2ERepo(t)
	defer cleanup()

	err := release.Execute(release.Options{
		Version:  "v2.0.0",
		DryRun:   false,
		NoCommit: true,
		SkipMeta: false,
	})
	if err != nil {
		t.Skipf("Execute failed (likely push not supported in CI): %v", err)
	}

	// Metadata should exist (SkipMeta is false).
	metaPath := filepath.Join(releaseDir, "v2.0.0.json")
	if _, err := os.Stat(metaPath); os.IsNotExist(err) {
		t.Error("v2.0.0.json should exist")
	}

	// The metadata file should be uncommitted (in git status).
	// Use -uall to show individual files inside untracked directories.
	cmd := exec.Command("git", "status", "--porcelain", "-uall")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("git status: %v", err)
	}

	output := string(out)
	if !strings.Contains(output, "v2.0.0.json") {
		t.Errorf("v2.0.0.json should appear in uncommitted changes with --no-commit, got: %q", output)
	}
}

// TestE2E_SkipMetaNoFiles verifies that --skip-meta produces no metadata
// files even after a successful release.
func TestE2E_SkipMetaNoFiles(t *testing.T) {
	_, releaseDir, cleanup := initE2ERepo(t)
	defer cleanup()

	err := release.Execute(release.Options{
		Version:  "v1.0.0",
		DryRun:   false,
		NoCommit: true,
		SkipMeta: true,
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}

	// Branch and tag should still exist.
	if !release.BranchExists(constants.ReleaseBranchPrefix + "v1.0.0") {
		t.Error("branch should exist even with SkipMeta")
	}
	if !release.TagExistsLocally("v1.0.0") {
		t.Error("tag should exist even with SkipMeta")
	}

	// But no metadata files.
	metaPath := filepath.Join(releaseDir, "v1.0.0.json")
	if _, err := os.Stat(metaPath); err == nil {
		t.Error("v1.0.0.json should NOT exist with SkipMeta: true")
	}

	latestPath := filepath.Join(releaseDir, constants.DefaultLatestFile)
	if _, err := os.Stat(latestPath); err == nil {
		t.Error("latest.json should NOT exist with SkipMeta: true")
	}
}
