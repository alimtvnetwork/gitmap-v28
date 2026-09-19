package cmdautomation

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestVersionSyncClean(t *testing.T) {
	tempDir := t.TempDir()
	setupCleanVersionFiles(t, tempDir, "6.262.0")
	monad := RunVersionSync(VersionSyncOptions{Dir: tempDir})
	isFail := monad.IsFailure()
	if isFail {
		t.Fatalf("expected success, got error: %v", monad.Err)
	}
	res := monad.Value
	isClean := res.IsClean
	if !isClean {
		t.Errorf("expected clean sync, got mismatch count: %d", res.MismatchCount)
	}
}

func setupCleanVersionFiles(t *testing.T, dir, ver string) {
	writeTestFile(t, dir, "version.json", `{"Version": "`+ver+`"}`)
	writeTestFile(t, dir, "package.json", `{"version": "`+ver+`"}`)
	writeTestFile(t, dir, "cli/constants/constants.go", `package constants; var Version = "`+ver+`"`)
	writeTestFile(t, dir, "changelog.md", `## [v`+ver+`] 2026-09-19 Release v`+ver+`\n`)
}

func writeTestFile(t *testing.T, baseDir, relPath, content string) {
	full := filepath.Join(baseDir, relPath)
	dir := filepath.Dir(full)
	_ = os.MkdirAll(dir, 0755)
	err := os.WriteFile(full, []byte(content), 0644)
	isErr := err != nil
	if isErr {
		t.Fatalf("failed to write test file %s: %v", relPath, err)
	}
}

func TestVersionSyncMismatchAndFix(t *testing.T) {
	tempDir := t.TempDir()
	setupMismatchVersionFiles(t, tempDir)
	monad := RunVersionSync(VersionSyncOptions{Dir: tempDir})
	isClean := monad.Value.IsClean
	if isClean {
		t.Errorf("expected mismatch, but got clean")
	}
	fixMonad := RunVersionSync(VersionSyncOptions{Dir: tempDir, IsFixMode: true})
	isFixSuccess := fixMonad.Value.FixedCount > 0
	if isFixSuccess == false {
		t.Errorf("expected fixed count > 0, got %d", fixMonad.Value.FixedCount)
	}
}

func setupMismatchVersionFiles(t *testing.T, dir string) {
	writeTestFile(t, dir, "version.json", `{"Version": "6.262.0"}`)
	writeTestFile(t, dir, "package.json", `{"version": "6.261.0"}`)
	writeTestFile(t, dir, "cli/constants/constants.go", `package constants; var Version = "6.260.0"`)
	writeTestFile(t, dir, "changelog.md", `## [v6.262.0] 2026-09-19 Release\n`)
}

func TestComputeNextSemVer(t *testing.T) {
	verifySemVerBump(t, "1.2.3", "major", "2.0.0")
	verifySemVerBump(t, "1.2.3", "minor", "1.3.0")
	verifySemVerBump(t, "1.2.3", "patch", "1.2.4")
}

func verifySemVerBump(t *testing.T, cur, tier, expected string) {
	actual := computeNextSemVer(cur, tier)
	isMatch := actual == expected
	if isMatch == false {
		t.Errorf("bump %s on %s: expected %s, got %s", tier, cur, expected, actual)
	}
}

func TestReleaseBumpDryRun(t *testing.T) {
	tempDir := t.TempDir()
	setupCleanVersionFiles(t, tempDir, "6.262.0")
	opts := ReleaseBumpOptions{Dir: tempDir, BumpType: "minor", IsDryRun: true}
	monad := RunReleaseBump(opts)
	isFail := monad.IsFailure()
	if isFail {
		t.Fatalf("expected success, got error: %v", monad.Err)
	}
	isMatch := monad.Value.NewVersion == "6.263.0"
	if isMatch == false {
		t.Errorf("expected 6.263.0, got %s", monad.Value.NewVersion)
	}
}

func TestReleaseBumpApply(t *testing.T) {
	tempDir := t.TempDir()
	setupCleanVersionFiles(t, tempDir, "6.262.0")
	opts := ReleaseBumpOptions{Dir: tempDir, BumpType: "patch", IsDryRun: false}
	monad := RunReleaseBump(opts)
	isSuccess := monad.Value.IsSuccess
	if isSuccess == false {
		t.Fatalf("expected successful bump apply")
	}
	isVerMatch := monad.Value.NewVersion == "6.262.1"
	if isVerMatch == false {
		t.Errorf("expected new version 6.262.1, got %s", monad.Value.NewVersion)
	}
}

func TestMilestonesReleaseNotes(t *testing.T) {
	issues := createSampleMilestoneIssues()
	notes := formatMilestoneReleaseNotes("14", issues)
	hasAdded := strings.Contains(notes, "### Added")
	if hasAdded == false {
		t.Errorf("expected '### Added' in release notes")
	}
	hasFixed := strings.Contains(notes, "### Fixed")
	if hasFixed == false {
		t.Errorf("expected '### Fixed' in release notes")
	}
}

func createSampleMilestoneIssues() []MilestoneIssue {
	return []MilestoneIssue{
		{Number: 101, Title: "feat: add version sync command", Category: "Added", IsClosed: true},
		{Number: 102, Title: "fix: resolve nested if in release bumper", Category: "Fixed", IsClosed: true},
		{Number: 103, Title: "refactor: clean up milestone types", Category: "Changed", IsClosed: true},
	}
}

func TestMilestonesLocalScan(t *testing.T) {
	tempDir := t.TempDir()
	planPath := "01-test-plan.md"
	writeTestFile(t, tempDir, ".ai-memory/plans/completed/"+planPath, "# Test Plan Feature\n\nSome details")
	monad := RunMilestones(MilestonesOptions{Dir: tempDir, MilestoneId: "all"})
	isFail := monad.IsFailure()
	if isFail {
		t.Fatalf("expected success, got %v", monad.Err)
	}
	hasItems := monad.Value.TotalItems > 0
	if hasItems == false {
		t.Errorf("expected at least 1 milestone item found")
	}
}
