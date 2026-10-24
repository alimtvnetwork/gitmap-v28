package cmdundo

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdfixrepo"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdundo/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdundo/enums"
)

func TestParseUndoArgsList(t *testing.T) {
	for _, flag := range []string{"--list", "-list", "-l"} {
		opts, err := parseUndoArgs([]string{flag})
		if err != nil || !opts.isList {
			t.Errorf("flag %s failed to set isList, err: %v", flag, err)
		}
	}
}

func TestParseUndoArgsDryRun(t *testing.T) {
	for _, flag := range []string{"--dry-run", "-dry-run", "-DryRun"} {
		opts, err := parseUndoArgs([]string{flag})
		if err != nil || !opts.isDryRun {
			t.Errorf("flag %s failed to set isDryRun, err: %v", flag, err)
		}
	}
}

func TestParseUndoArgsSnapshot(t *testing.T) {
	for _, flag := range []string{"--snapshot", "-snapshot", "-s"} {
		opts, err := parseUndoArgs([]string{flag, "20260101T000000Z"})
		if err != nil || opts.snapshot != "20260101T000000Z" {
			t.Errorf("flag %s failed to set snapshot, err: %v", flag, err)
		}
	}
}

func TestParseUndoArgsInvalidFlag(t *testing.T) {
	_, err := parseUndoArgs([]string{"--invalid-flag"})
	if err == nil {
		t.Error("expected error for invalid flag")
	}
}

func TestSnapshotListing(t *testing.T) {
	tmpDir := t.TempDir()
	_ = os.Mkdir(filepath.Join(tmpDir, "20260101T100000Z"), 0755)
	_ = os.Mkdir(filepath.Join(tmpDir, "20260102T100000Z"), 0755)
	_ = os.WriteFile(filepath.Join(tmpDir, "ignore.txt"), []byte("txt"), 0644)

	snaps := listUndoSnapshots(tmpDir)
	if len(snaps) != 2 {
		t.Fatalf("expected 2 snapshots, got %d", len(snaps))
	}
	if snaps[0] != "20260102T100000Z" {
		t.Errorf("expected newest snapshot first, got %s", snaps[0])
	}
}

func TestSnapshotPicking(t *testing.T) {
	snaps := []string{"20260102T000000Z", "20260101T000000Z"}
	if pickUndoSnapshot(snaps, "") != "20260102T000000Z" {
		t.Error("expected default to newest snapshot")
	}
	if pickUndoSnapshot(snaps, "20260101T000000Z") != "20260101T000000Z" {
		t.Error("expected explicit snapshot match")
	}
	if pickUndoSnapshot(snaps, "missing") != "" {
		t.Error("expected empty string for missing explicit snapshot")
	}
	if pickUndoSnapshot(nil, "") != "" {
		t.Error("expected empty string for empty snapshots slice")
	}
}

func TestSnapshotMarker(t *testing.T) {
	if snapshotMarker(0) != constants.MarkerLatest {
		t.Errorf("expected %s for index 0", constants.MarkerLatest)
	}
	if snapshotMarker(1) != constants.MarkerOlder {
		t.Errorf("expected %s for index 1", constants.MarkerOlder)
	}
}

func TestResolveRestoreMode(t *testing.T) {
	if resolveRestoreMode(true) != constants.ModeDryRun {
		t.Errorf("expected dry-run mode")
	}
	if resolveRestoreMode(false) != constants.ModeWrite {
		t.Errorf("expected write mode")
	}
	if !enums.UndoModeDryRun.IsDryRun() || !enums.UndoActionRestored.IsRestored() {
		t.Error("unexpected enum predicate state")
	}
}

func TestPrintUndoSnapshotList(t *testing.T) {
	printUndoSnapshotList("/tmp/test", nil)
	printUndoSnapshotList("/tmp/test", []string{"20260101T000000Z"})
}

func TestManifestOperations(t *testing.T) {
	tmpDir := t.TempDir()
	if countUndoFiles(tmpDir) != 0 {
		t.Error("expected 0 count for missing manifest")
	}

	manifest := cmdfixrepo.FixRepoBackupManifest{
		SchemaVersion: 1,
		Files:         []string{"foo.txt", "bar.txt"},
	}
	body, _ := json.Marshal(manifest)
	_ = os.WriteFile(filepath.Join(tmpDir, constants.ManifestName), body, 0644)

	loaded, hasManifest := readUndoManifest(tmpDir)
	if !hasManifest || len(loaded.Files) != 2 {
		t.Errorf("failed to read manifest: %v, hasManifest: %v", loaded, hasManifest)
	}
	if countUndoFiles(tmpDir) != 2 {
		t.Errorf("expected 2 files counted, got %d", countUndoFiles(tmpDir))
	}
}

func TestRestoreSnapshotWorkflow(t *testing.T) {
	snapDir := t.TempDir()
	repoRoot := t.TempDir()

	manifest := cmdfixrepo.FixRepoBackupManifest{
		SchemaVersion: 1,
		Files:         []string{"test.txt"},
	}
	body, _ := json.Marshal(manifest)
	_ = os.WriteFile(filepath.Join(snapDir, constants.ManifestName), body, 0644)

	filesDir := filepath.Join(snapDir, constants.SubdirFiles)
	_ = os.MkdirAll(filesDir, 0755)
	_ = os.WriteFile(filepath.Join(filesDir, "test.txt"), []byte("restored content"), 0644)

	testDryRunRestore(t, snapDir, repoRoot)
	testWriteRestore(t, snapDir, repoRoot)
}

func testDryRunRestore(t *testing.T, snapDir, repoRoot string) {
	if err := restoreUndoSnapshot(snapDir, repoRoot, true); err != nil {
		t.Fatalf("dry run restore failed: %v", err)
	}
	if _, err := os.Stat(filepath.Join(repoRoot, "test.txt")); err == nil {
		t.Error("file should not exist after dry run")
	}
}

func testWriteRestore(t *testing.T, snapDir, repoRoot string) {
	if err := restoreUndoSnapshot(snapDir, repoRoot, false); err != nil {
		t.Fatalf("write restore failed: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(repoRoot, "test.txt"))
	if err != nil || string(data) != "restored content" {
		t.Errorf("restored file content mismatch: %q, err: %v", string(data), err)
	}
}

func TestUndoFlattenHelpers(t *testing.T) {
	testFlattenHelp(t)
	testFlattenTargetExtraction(t)
	testFlattenDirChecks(t)
}

func testFlattenHelp(t *testing.T) {
	if !isHelpRequest("-h") || !isHelpRequest("--help") || !isHelpRequest("help") {
		t.Error("expected true for help flags")
	}
	if isHelpRequest("my-repo") {
		t.Error("expected false for regular repo name")
	}
}

func testFlattenTargetExtraction(t *testing.T) {
	args := []string{"--foo", "-bar", "target-repo"}
	if extractFlattenTarget(args) != "target-repo" {
		t.Errorf("expected target-repo, got %s", extractFlattenTarget(args))
	}
	if extractFlattenTarget([]string{"--foo"}) != "" {
		t.Error("expected empty target when only flags provided")
	}
}

func testFlattenDirChecks(t *testing.T) {
	if !isInsideOrSameDir("/path/to/child", "/path/to") {
		t.Error("expected child to be inside parent")
	}
	if isInsideOrSameDir("/other/dir", "/path/to") {
		t.Error("expected other dir to not be inside parent")
	}
}

func TestRunUndoFlattenFlow(t *testing.T) {
	if err := RunUndoFlatten("-h"); err != nil {
		t.Fatalf("unexpected error on -h: %v", err)
	}

	tmpDir := t.TempDir()
	repoPath := filepath.Join(tmpDir, "sample-repo")
	archivePath := repoPath + constants.SuffixMarkedForDel
	_ = os.Mkdir(archivePath, 0755)
	_ = os.WriteFile(filepath.Join(archivePath, "sample.txt"), []byte("data"), 0644)

	if err := RunUndoFlatten(repoPath); err != nil {
		t.Fatalf("RunUndoFlatten failed: %v", err)
	}
	if !hasCanonicalPath(repoPath) {
		t.Error("expected canonical repo path to exist after restore")
	}
	if hasArchivePath(archivePath) {
		t.Error("expected archive path to be moved")
	}
}

func TestRunUndoFlattenCLI(t *testing.T) {
	if err := RunUndoFlattenCLI([]string{"-h"}); err != nil {
		t.Fatalf("unexpected error on -h: %v", err)
	}
}

func TestRunUndoHelpAndBadFlag(t *testing.T) {
	called := false
	CheckHelpFn = func(command string, args []string) {
		called = true
	}
	defer func() { CheckHelpFn = nil }()

	_ = RunUndo([]string{"--list"})
	if !called {
		t.Error("expected CheckHelpFn to be called")
	}

	err := RunUndo([]string{"--unknown-flag"})
	if err == nil {
		t.Error("expected error on unknown flag")
	}
}
