package constants

// Command names and aliases for cmdundo.
const (
	CmdUndo            = "undo"
	CmdUndoAlias       = "ud"
	CmdUndoFlatten     = "undo-flatten"
	CmdRepoUndoFlatten = "repo-undo-flatten"
	CmdFixRepo         = "fix-repo"
)

// Flag literals and aliases for cmdundo.
const (
	FlagListLong       = "--list"
	FlagListShort1     = "-list"
	FlagListShort2     = "-l"
	FlagDryRunLong     = "--dry-run"
	FlagDryRunShort1   = "-dry-run"
	FlagDryRunShort2   = "-DryRun"
	FlagSnapshotLong   = "--snapshot"
	FlagSnapshotShort1 = "-snapshot"
	FlagSnapshotShort2 = "-s"
	FlagHelpLong       = "--help"
	FlagHelpShort      = "-h"
	FlagHelpWord       = "help"
)

// Directory paths, subtrees, and markers.
const (
	GitMapDir          = ".gitmap"
	SubdirBackup       = "backup"
	SubdirFiles        = "files"
	ManifestName       = "manifest.json"
	SuffixMarkedForDel = "-marked-for-deletion"
	SentinelCurrentDir = "."
	MarkerLatest       = "*"
	MarkerOlder        = " "
	ModeWrite          = "write"
	ModeDryRun         = "dry-run"
)

// ANSI terminal color formatting codes.
const (
	ColorGreen = "\033[32m"
	ColorReset = "\033[0m"
)

// Output format strings for interactive execution.
const (
	MsgNoSnapshotsFmt            = "undo: no snapshots found under %s\n"
	MsgListHeaderFmt             = "undo: snapshots under %s (%d total)\n"
	MsgListRowFmt                = "  %s %s  (%d files)\n"
	MsgRestoreHeaderFmt          = "undo: restoring snapshot %s — %d file(s) [mode: %s]\n"
	MsgRestoreRowFmt             = "  [%s] %s\n"
	MsgRestoreSummaryFmt         = "undo: restored %d file(s), %d failure(s)\n"
	MsgRestoreErrFmt             = "undo: restore failed for %s: %v\n"
	MsgFlattenRestoringFmt       = "Restoring %s from %s...\n"
	MsgFlattenRestoredSuccessFmt = "\n%s✓ Repository %s successfully restored from %s!%s\n"
	MsgFlattenRestoredPathFmt    = "  Restored Path: %s\n\n"
	MsgFlattenRemovingTreeFmt    = "Removing flattened working tree %s...\n"
	MsgFlattenUsage              = "Usage: gitmap repo undo-flatten [repo]\n\nRestore a flattened repository from its <repo>-marked-for-deletion backup archive.\n"
)

// Diagnostic error codes, operations, and format strings.
const (
	ErrCodeE9000                 = "E9000"
	ErrNoSnapshotFmt             = "undo: no snapshot to restore under %s\n"
	ErrBadFlagFmt                = "undo: unknown flag %q\n"
	ErrManifestMissingFmt        = "undo: manifest unreadable at %s: %v\n"
	ErrManifestBadFmt            = "undo: manifest malformed at %s: %v\n"
	ErrArchiveNotFoundFmt        = "archive %q not found for repository %q (expected %s)"
	OpResolveRepoPathFmt         = "resolve repository path %q"
	OpRestoreArchiveFmt          = "restore archive %s to %s"
	OpResolveCwd                 = "resolve current working directory"
	OpRemoveExistingFlattenedFmt = "remove existing flattened repository %s"
)

// Exit codes for process termination.
const (
	ExitOk          = 0
	ExitBadFlag     = 6
	ExitWriteFailed = 7
	ExitBadConfig   = 8
)
