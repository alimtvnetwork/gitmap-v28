package constants

import (
	"testing"
)

func TestCommandConstants(t *testing.T) {
	if CmdUndo != "undo" {
		t.Errorf("expected undo, got %s", CmdUndo)
	}
	if CmdUndoAlias != "ud" {
		t.Errorf("expected ud, got %s", CmdUndoAlias)
	}
	if CmdUndoFlatten != "undo-flatten" {
		t.Errorf("expected undo-flatten, got %s", CmdUndoFlatten)
	}
	if CmdRepoUndoFlatten != "repo-undo-flatten" {
		t.Errorf("expected repo-undo-flatten, got %s", CmdRepoUndoFlatten)
	}
	if CmdFixRepo != "fix-repo" {
		t.Errorf("expected fix-repo, got %s", CmdFixRepo)
	}
}

func TestFlagConstants(t *testing.T) {
	testFlagList(t)
	testFlagDryRun(t)
	testFlagSnapshot(t)
	testFlagHelp(t)
}

func testFlagList(t *testing.T) {
	if FlagListLong != "--list" || FlagListShort1 != "-list" || FlagListShort2 != "-l" {
		t.Error("unexpected list flag constants")
	}
}

func testFlagDryRun(t *testing.T) {
	if FlagDryRunLong != "--dry-run" || FlagDryRunShort1 != "-dry-run" || FlagDryRunShort2 != "-DryRun" {
		t.Error("unexpected dry-run flag constants")
	}
}

func testFlagSnapshot(t *testing.T) {
	if FlagSnapshotLong != "--snapshot" || FlagSnapshotShort1 != "-snapshot" || FlagSnapshotShort2 != "-s" {
		t.Error("unexpected snapshot flag constants")
	}
}

func testFlagHelp(t *testing.T) {
	if FlagHelpLong != "--help" || FlagHelpShort != "-h" || FlagHelpWord != "help" {
		t.Error("unexpected help flag constants")
	}
}

func TestPathAndMarkerConstants(t *testing.T) {
	if GitMapDir != ".gitmap" || SubdirBackup != "backup" || SubdirFiles != "files" {
		t.Error("unexpected path constants")
	}
	if ManifestName != "manifest.json" || SuffixMarkedForDel != "-marked-for-deletion" {
		t.Error("unexpected manifest or suffix constants")
	}
	if SentinelCurrentDir != "." || MarkerLatest != "*" || MarkerOlder != " " {
		t.Error("unexpected marker constants")
	}
	if ModeWrite != "write" || ModeDryRun != "dry-run" {
		t.Error("unexpected mode constants")
	}
}

func TestColorConstants(t *testing.T) {
	if ColorGreen != "\033[32m" || ColorReset != "\033[0m" {
		t.Error("unexpected color constants")
	}
}

func TestMessageAndErrorConstants(t *testing.T) {
	testMessageConstants(t)
	testErrorConstants(t)
	testExitCodeConstants(t)
}

func testMessageConstants(t *testing.T) {
	if len(MsgNoSnapshotsFmt) == 0 || len(MsgListHeaderFmt) == 0 || len(MsgListRowFmt) == 0 {
		t.Error("empty message format constants")
	}
	if len(MsgRestoreHeaderFmt) == 0 || len(MsgRestoreRowFmt) == 0 || len(MsgRestoreSummaryFmt) == 0 {
		t.Error("empty restore message format constants")
	}
	if len(MsgFlattenRestoringFmt) == 0 || len(MsgFlattenRestoredSuccessFmt) == 0 || len(MsgFlattenUsage) == 0 {
		t.Error("empty flatten message format constants")
	}
}

func testErrorConstants(t *testing.T) {
	if ErrCodeE9000 != "E9000" {
		t.Errorf("expected E9000, got %s", ErrCodeE9000)
	}
	if len(ErrNoSnapshotFmt) == 0 || len(ErrBadFlagFmt) == 0 || len(ErrManifestMissingFmt) == 0 {
		t.Error("empty error format constants")
	}
	if len(ErrArchiveNotFoundFmt) == 0 || len(OpResolveRepoPathFmt) == 0 || len(OpRestoreArchiveFmt) == 0 {
		t.Error("empty operation format constants")
	}
}

func testExitCodeConstants(t *testing.T) {
	if ExitOk != 0 || ExitBadFlag != 6 || ExitWriteFailed != 7 || ExitBadConfig != 8 {
		t.Error("unexpected exit code constants")
	}
}
