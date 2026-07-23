package constants

// Backup + undo constants for `gitmap fix-repo` snapshots (v5.40.0+).
//
// Filesystem layout under the repo root:
//
//   .gitmap/backup/<repo-name>/v<current>/fix-repo/<UTC-timestamp>/
//     manifest.json
//     files/<rel/path>
//
// `gitmap undo` reads the latest snapshot to restore pre-rewrite copies.

const (
	// FixRepoBackupSubdir is the namespace under `.gitmap/` that holds
	// per-command backup trees. Currently only `fix-repo` writes here.
	FixRepoBackupSubdir = "backup"
	// FixRepoBackupFilesSubdir is the verbatim file-copy subtree inside
	// each timestamped snapshot directory.
	FixRepoBackupFilesSubdir = "files"
	// FixRepoBackupManifestName is the per-snapshot index file.
	FixRepoBackupManifestName = "manifest.json"
	// FixRepoBackupTimestampFmt is the canonical UTC stamp embedded in
	// every snapshot dir name. Lexical sort == chronological sort.
	FixRepoBackupTimestampFmt = "20060102T150405Z"
	// FixRepoBackupSchemaVersion is bumped when manifest.json fields
	// change in a breaking way. Readers must tolerate unknown fields.
	FixRepoBackupSchemaVersion = 1
)

// fix-repo backup user-facing strings.
const (
	FixRepoBackupMsgFmt         = "fix-repo: backed up %d file(s) → %s\n"
	FixRepoBackupErrFmt         = "fix-repo: backup failed for %s: %v\n"
	FixRepoBackupManifestErrFmt = "fix-repo: failed to write backup manifest: %v\n"
)

// `gitmap undo` user-facing strings.
const (
	UndoMsgNoSnapshotsFmt     = "undo: no snapshots found under %s\n"
	UndoMsgListHeaderFmt      = "undo: snapshots under %s (%d total)\n"
	UndoMsgListRowFmt         = "  %s %s  (%d files)\n"
	UndoMsgRestoreHeaderFmt   = "undo: restoring snapshot %s — %d file(s) [mode: %s]\n"
	UndoMsgRestoreRowFmt      = "  [%s] %s\n"
	UndoMsgRestoreSummaryFmt  = "undo: restored %d file(s), %d failure(s)\n"
	UndoMsgRestoreErrFmt      = "undo: restore failed for %s: %v\n"
	UndoErrNoSnapshotFmt      = "undo: no snapshot to restore under %s\n"
	UndoErrBadFlagFmt         = "undo: unknown flag %q\n"
	UndoErrManifestMissingFmt = "undo: manifest unreadable at %s: %v\n"
	UndoErrManifestBadFmt     = "undo: manifest malformed at %s: %v\n"
)
