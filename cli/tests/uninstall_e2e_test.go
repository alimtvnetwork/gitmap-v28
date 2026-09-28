//go:build e2e

package tests

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdagy"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdwinutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/osclean"
)

// TestE2EAGYRestoreSnapshot verifies snapshot generation produces valid JSON with metadata.
func TestE2EAGYRestoreSnapshot(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "gitmap-snap-e2e-*")
	hasErr := err != nil
	if hasErr {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	targetFile := filepath.Join(tempDir, "agy-snapshot-test.json")
	snap, snapErr := cmdagy.ExportAGYRestoreSnapshot(targetFile)
	hasSnapErr := snapErr != nil
	if hasSnapErr {
		t.Fatalf("snapshot export failed: %v", snapErr)
	}
	verifySnapshotOutput(t, targetFile, snap)
}

func verifySnapshotOutput(t *testing.T, filePath string, snap cmdagy.AGYRestoreSnapshot) {
	data, readErr := os.ReadFile(filePath)
	hasReadErr := readErr != nil
	if hasReadErr {
		t.Fatalf("snapshot file missing: %v", readErr)
	}
	var decoded cmdagy.AGYRestoreSnapshot
	decodeErr := json.Unmarshal(data, &decoded)
	hasDecodeErr := decodeErr != nil
	if hasDecodeErr {
		t.Fatalf("invalid json generated: %v", decodeErr)
	}
	hasMatch := decoded.Timestamp == snap.Timestamp
	if !hasMatch {
		t.Errorf("timestamp mismatch: expected %d, got %d", snap.Timestamp, decoded.Timestamp)
	}
	hasFilePath := decoded.FilePath != ""
	if !hasFilePath {
		t.Errorf("expected snapshot FilePath populated, got empty")
	}
}

// TestE2EUninstallSafetyProtection verifies path safety invariants protecting d:\work and root.
func TestE2EUninstallSafetyProtection(t *testing.T) {
	assertWorkOverlap(t, "d:/work", true)
	assertWorkOverlap(t, "d:/work/sample-repo", true)
	assertWorkOverlap(t, "D:\\Work\\Anything", true)
	assertWorkOverlap(t, "C:\\Users\\test", false)

	assertPathSafety(t, "d:/work", false)
	assertPathSafety(t, "d:/work/sample-repo", false)
	assertPathSafety(t, "C:\\", false)
	assertPathSafety(t, "/", false)
	assertPathSafety(t, "", false)
	assertPathSafety(t, ".", false)
}

func assertWorkOverlap(t *testing.T, target string, expected bool) {
	t.Helper()
	isOverlap := cmdagy.IsWorkDirectoryOverlap(target)
	hasMismatch := isOverlap != expected
	if hasMismatch {
		t.Errorf("IsWorkDirectoryOverlap(%s) = %v, want %v", target, isOverlap, expected)
	}
}

func assertPathSafety(t *testing.T, target string, expected bool) {
	t.Helper()
	isSafe := cmdagy.IsPathSafeToDelete(target)
	hasMismatch := isSafe != expected
	if hasMismatch {
		t.Errorf("IsPathSafeToDelete(%s) = %v, want %v", target, isSafe, expected)
	}
}

// TestE2EAGYUninstallDryRun verifies standard and purge dry-run uninstallations execute safely.
func TestE2EAGYUninstallDryRun(t *testing.T) {
	errStandard := cmdagy.RunAGYUninstall(false, true, true, "")
	hasStdErr := errStandard != nil
	if hasStdErr {
		t.Fatalf("standard dry-run failed: %v", errStandard)
	}

	tempDir, _ := os.MkdirTemp("", "agy-dryrun-*")
	defer os.RemoveAll(tempDir)
	backupPath := filepath.Join(tempDir, "backup.json")

	errPurge := cmdagy.RunAGYUninstall(true, true, true, backupPath)
	hasPurgeErr := errPurge != nil
	if hasPurgeErr {
		t.Fatalf("full purge dry-run failed: %v", errPurge)
	}
}

// TestE2EWinUtilCopilotDryRun verifies Windows Copilot removal dry-run.
func TestE2EWinUtilCopilotDryRun(t *testing.T) {
	res, err := cmdwinutil.RunCopilotUninstall(true)
	hasErr := err != nil
	if hasErr {
		t.Fatalf("copilot uninstall dry-run returned error: %v", err)
	}
	isDryRun := res.IsDryRun
	if !isDryRun {
		t.Errorf("expected IsDryRun true, got false")
	}
	isOk := res.IsSuccess
	if !isOk {
		t.Errorf("expected IsSuccess true, got false")
	}
}

// TestE2EWinUtilEdgeDryRun verifies Microsoft Edge WinUtil removal dry-run.
func TestE2EWinUtilEdgeDryRun(t *testing.T) {
	res, err := cmdwinutil.RunEdgeUninstall(true, true)
	hasErr := err != nil
	if hasErr {
		t.Fatalf("edge uninstall dry-run returned error: %v", err)
	}
	isDryRun := res.IsDryRun
	if !isDryRun {
		t.Errorf("expected IsDryRun true, got false")
	}
	isOk := res.IsSuccess
	if !isOk {
		t.Errorf("expected IsSuccess true, got false")
	}
}

// TestE2EDevToolCleanerDryRun verifies 12-category cache cleaner dry-run and formatting.
func TestE2EDevToolCleanerDryRun(t *testing.T) {
	opts := osclean.DevCleanOptions{IsDryRun: true}
	summary := osclean.CleanEnhancedDevCaches(opts)

	isDryRun := summary.IsDryRun
	if !isDryRun {
		t.Errorf("expected summary IsDryRun true, got false")
	}
	hasTwelveCategories := len(summary.Categories) == 12
	if !hasTwelveCategories {
		t.Errorf("expected 12 categories, got %d", len(summary.Categories))
	}

	testFormatCleanSize(t)
}

func testFormatCleanSize(t *testing.T) {
	mbStr := osclean.FormatCleanSize(50 * 1024 * 1024)
	hasMB := mbStr == "50.00 MB"
	if !hasMB {
		t.Errorf("expected 50.00 MB, got %s", mbStr)
	}

	gbStr := osclean.FormatCleanSize(2 * 1024 * 1024 * 1024)
	hasGB := gbStr == "2.00 GB"
	if !hasGB {
		t.Errorf("expected 2.00 GB, got %s", gbStr)
	}
}

// TestE2EAGYRestoreFromSnapshot verifies restoring projects and conversations from snapshot JSON.
func TestE2EAGYRestoreFromSnapshot(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "gitmap-restore-e2e-*")
	hasErr := err != nil
	if hasErr {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	targetFile := filepath.Join(tempDir, "agy-snapshot-restore.json")
	snap, snapErr := cmdagy.ExportAGYRestoreSnapshot(targetFile)
	hasSnapErr := snapErr != nil
	if hasSnapErr {
		t.Fatalf("snapshot export failed: %v", snapErr)
	}

	pCount, _, resErr := cmdagy.RunAGYRestore(targetFile)
	hasResErr := resErr != nil
	if hasResErr {
		t.Fatalf("RunAGYRestore failed: %v", resErr)
	}
	isMatch := pCount == snap.TotalProjects
	if !isMatch {
		t.Errorf("project count mismatch: expected %d, got %d", snap.TotalProjects, pCount)
	}
}
