package constants

// gitmap:cmd top-level
// Temp-release command names.
const (
	CmdTempRelease      = "temp-release"
	CmdTempReleaseShort = "tr"
	SubCmdTRList        = "list"
	SubCmdTRRemove      = "remove"
)

// Temp-release branch prefix.
const TempReleaseBranchPrefix = "temp-release/"

// Temp-release limits.
const TempReleaseMaxCount = 50

// Temp-release placeholder.
const TempReleasePlaceholder = "$$"

// TempRelease table (v15: singular + TempReleaseId PK).
const TableTempRelease = "TempRelease"

// Legacy plural retained for migration detection.
const LegacyTableTempReleases = "TempReleases"

// SQL: create TempRelease table (v15).
const SQLCreateTempRelease = `CREATE TABLE IF NOT EXISTS TempRelease (
	TempReleaseId  INTEGER PRIMARY KEY AUTOINCREMENT,
	Branch         TEXT NOT NULL UNIQUE,
	VersionPrefix  TEXT NOT NULL DEFAULT '',
	SequenceNumber INTEGER NOT NULL DEFAULT 0,
	CommitSha      TEXT NOT NULL DEFAULT '',
	CommitMessage  TEXT NOT NULL DEFAULT '',
	CreatedAt      TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
)`

// SQL: temp-release operations (v15).
const (
	SQLInsertTempRelease = `INSERT INTO TempRelease (Branch, VersionPrefix, SequenceNumber, CommitSha, CommitMessage)
		VALUES (?, ?, ?, ?, ?)`

	SQLSelectAllTempReleases = `SELECT TempReleaseId, Branch, VersionPrefix, SequenceNumber, CommitSha, CommitMessage, CreatedAt
		FROM TempRelease ORDER BY SequenceNumber`

	SQLSelectMaxSeqByPrefix = `SELECT COALESCE(MAX(SequenceNumber), 0) FROM TempRelease WHERE VersionPrefix = ?`

	SQLDeleteTempRelease = `DELETE FROM TempRelease WHERE Branch = ?`

	SQLDeleteAllTempReleases = `DELETE FROM TempRelease`

	SQLCountTempReleases = `SELECT COUNT(*) FROM TempRelease`
)

// SQL: drop TempRelease table (and legacy plural).
const (
	SQLDropTempRelease  = "DROP TABLE IF EXISTS TempRelease"
	SQLDropTempReleases = "DROP TABLE IF EXISTS TempReleases" // legacy
)

// SQL: migrate Commit → CommitSha column. Operates on legacy TempReleases —
// the v15 rebuild copies the already-renamed column into TempRelease.
const SQLMigrateTRCommitSha = `ALTER TABLE TempReleases RENAME COLUMN "Commit" TO CommitSha`

// Temp-release flag descriptions.
const (
	FlagDescTRStart   = "Starting sequence number (default: auto-increment)"
	FlagDescTRDryRun  = "Preview branch names without creating"
	FlagDescTRJSON    = "Output structured JSON"
	FlagDescTRVerbose = "Detailed logging"
)

// Temp-release messages.
const (
	MsgTRCreating      = "  Creating %d temp-release branch(es)...\n"
	MsgTRCreated       = "  ✓ Created %s from %s\n"
	MsgTRPushing       = "  Pushing %d branch(es) to origin...\n"
	MsgTRPushed        = "  ✓ Pushed %d branch(es) to origin\n"
	MsgTRSeqStart      = "  → Starting sequence: %d\n"
	MsgTRSeqAuto       = "  → Starting sequence: %d (auto-detected)\n"
	MsgTRDryRunHeader  = "  Dry-run: would create %d temp-release branch(es):\n"
	MsgTRDryRunEntry   = "    %s  %s  %s\n"
	MsgTRListHeader    = "\n  Temp-release branches (%d):\n\n"
	MsgTRListRow       = "  %-35s %s  %-50s %s\n"
	MsgTRListEmpty     = "  No temp-release branches found.\n"
	MsgTRRemovePrompt  = "  Remove %s? (y/N): "
	MsgTRRemoveRange   = "  Remove %d temp-release branch(es):\n"
	MsgTRRemoveAll     = "  Remove ALL %d temp-release branch(es):\n"
	MsgTRRemoveBranch  = "    %s\n"
	MsgTRRemoveConfirm = "  Proceed? (y/N): "
	MsgTRRemoved       = "  ✓ Removed %d temp-release branch(es) (local + remote)\n"
	MsgTRRemovedOne    = "  ✓ Removed %s (local + remote)\n"
	MsgTRSkipExists    = "  ⚠ Branch already exists, skipping: %s\n"
	MsgTRSkipMissing   = "  ⚠ Branch not found, skipping: %s\n"
	MsgTRComplete      = "  Temp-release complete.\n"
	MsgTRNoneToRemove  = "  No temp-release branches to remove.\n"
)

// Temp-release error messages.
const (
	ErrTRUsage         = "Usage: gitmap temp-release <count> <version-pattern> [-s N]"
	ErrTRInvalidCount  = "count must be between 1 and %d"
	ErrTRNoPlaceholder = "version pattern must contain at least one '$' placeholder (e.g., v1.$$)"
	ErrTROverflow      = "sequence %d exceeds %d-digit format (max %d)"
	ErrTRNotEnough     = "  ⚠ Only %d commit(s) available (requested %d)\n"
	ErrTRCreate        = "failed to create temp-release: %v"
	ErrTRQuery         = "failed to query temp-releases: %v"
	ErrTRDelete        = "failed to delete temp-release: %v"
	ErrTRRemoveUsage   = "Usage: gitmap tr remove <version> | <v1> to <v2> | all"
)

// Temp-release help text.
const (
	HelpTempRelease = "  temp-release (tr) <count> <pattern> [-s N]  Create temp branches from recent commits"
)
