package constants

// ScanFolder + VersionProbe schema (v3.7.0+, Phase 2.1).
//
// A ScanFolder is the root absolute path that `gitmap scan` was invoked
// against. One row per unique AbsolutePath (UNIQUE index). Repos record the
// ScanFolderId of the *most recent* scan that discovered them via a nullable
// FK on Repo.ScanFolderId. Old repos stay NULL until rescanned — no backfill.
//
// VersionProbe stores the result of the hybrid HEAD-then-clone version probe
// for each repo. Empty in Phase 2.1; populated starting in Phase 2.3.

const (
	TableScanFolder   = "ScanFolder"
	TableVersionProbe = "VersionProbe"
)

// SQL: create ScanFolder table.
const SQLCreateScanFolder = `CREATE TABLE IF NOT EXISTS ScanFolder (
	ScanFolderId  INTEGER PRIMARY KEY AUTOINCREMENT,
	AbsolutePath  TEXT NOT NULL,
	Label         TEXT DEFAULT '',
	Notes         TEXT DEFAULT '',
	LastScannedAt TEXT DEFAULT CURRENT_TIMESTAMP,
	CreatedAt     TEXT DEFAULT CURRENT_TIMESTAMP
)`

// SQL: unique index on AbsolutePath so EnsureScanFolder is idempotent.
const SQLCreateScanFolderPathIndex = "CREATE UNIQUE INDEX IF NOT EXISTS IdxScanFolder_AbsolutePath ON ScanFolder(AbsolutePath)"

// SQL: create VersionProbe table (populated starting Phase 2.3).
const SQLCreateVersionProbe = `CREATE TABLE IF NOT EXISTS VersionProbe (
	VersionProbeId  INTEGER PRIMARY KEY AUTOINCREMENT,
	RepoId          INTEGER NOT NULL REFERENCES Repo(RepoId) ON DELETE CASCADE,
	ProbedAt        TEXT DEFAULT CURRENT_TIMESTAMP,
	NextVersionTag  TEXT DEFAULT '',
	NextVersionNum  INTEGER DEFAULT 0,
	Method          TEXT DEFAULT '',
	IsAvailable     INTEGER DEFAULT 0,
	Error           TEXT DEFAULT ''
)`

// SQL: index for fast latest-probe lookups per repo.
const SQLCreateVersionProbeRepoIndex = "CREATE INDEX IF NOT EXISTS IdxVersionProbe_RepoId ON VersionProbe(RepoId, ProbedAt DESC)"

// SQL: ALTER Repo with nullable ScanFolderId FK. Idempotent via
// addColumnIfNotExists. SQLite cannot add a REFERENCES clause via ALTER
// without a table rebuild, so the column stores the FK value without a
// declared FOREIGN KEY constraint — application code enforces validity.
const SQLAddRepoScanFolderId = "ALTER TABLE Repo ADD COLUMN ScanFolderId INTEGER DEFAULT NULL"

// SQL: ScanFolder operations.
const (
	SQLUpsertScanFolder = `INSERT INTO ScanFolder (AbsolutePath, Label, Notes)
		VALUES (?, ?, ?)
		ON CONFLICT(AbsolutePath) DO UPDATE SET
			LastScannedAt=CURRENT_TIMESTAMP,
			Label=CASE WHEN excluded.Label = '' THEN ScanFolder.Label ELSE excluded.Label END,
			Notes=CASE WHEN excluded.Notes = '' THEN ScanFolder.Notes ELSE excluded.Notes END`

	SQLSelectAllScanFolders = `SELECT ScanFolderId, AbsolutePath, Label, Notes, LastScannedAt, CreatedAt
		FROM ScanFolder ORDER BY LastScannedAt DESC, AbsolutePath ASC`

	SQLSelectScanFolderByPath = `SELECT ScanFolderId, AbsolutePath, Label, Notes, LastScannedAt, CreatedAt
		FROM ScanFolder WHERE AbsolutePath = ?`

	SQLSelectScanFolderByID = `SELECT ScanFolderId, AbsolutePath, Label, Notes, LastScannedAt, CreatedAt
		FROM ScanFolder WHERE ScanFolderId = ?`

	SQLCountReposInScanFolder = `SELECT COUNT(*) FROM Repo WHERE ScanFolderId = ?`

	SQLDeleteScanFolderByID   = `DELETE FROM ScanFolder WHERE ScanFolderId = ?`
	SQLDeleteScanFolderByPath = `DELETE FROM ScanFolder WHERE AbsolutePath = ?`

	// Detach: set Repo.ScanFolderId = NULL for any repos pointing at the
	// scan folder being removed. Run BEFORE DELETE to avoid orphan FK ids.
	SQLDetachReposFromScanFolder = `UPDATE Repo SET ScanFolderId = NULL WHERE ScanFolderId = ?`
)

// SQL: drop statements for Reset() ordering.
const (
	SQLDropScanFolder   = "DROP TABLE IF EXISTS ScanFolder"
	SQLDropVersionProbe = "DROP TABLE IF EXISTS VersionProbe"
)

// ScanFolder error messages (Code Red zero-swallow policy).
const (
	ErrSFEnsure      = "failed to ensure scan folder %q: %v"
	ErrSFList        = "failed to list scan folders: %v"
	ErrSFFindByPath  = "no scan folder registered for path: %s"
	ErrSFFindByID    = "no scan folder with id: %d"
	ErrSFRemove      = "failed to remove scan folder: %v"
	ErrSFDetachRepos = "failed to detach repos from scan folder: %v"
	ErrSFAbsResolve  = "failed to resolve absolute path for %q: %v"
	ErrSFInvalidID   = "invalid scan folder id %q: %v"
	ErrSFMissingArg  = "missing required argument: %s"
)

// ScanFolder user-facing CLI strings.
const (
	MsgSFAddedFmt       = "✓ Registered scan folder: %s (id=%d)\n"
	MsgSFAddedExistsFmt = "✓ Scan folder already registered: %s (id=%d, last scanned %s)\n"
	MsgSFRemovedFmt     = "✓ Removed scan folder: %s (id=%d, %d repos detached)\n"
	MsgSFListEmpty      = "No scan folders registered. Run `gitmap scan <dir>` or `gitmap sf add <dir>`.\n"
	MsgSFListHeaderFmt  = "Scan folders (%d):\n"
	MsgSFListRowFmt     = "  [%d] %s\n      label: %s | repos: %d | last scanned: %s\n"
	MsgSFUsageHeader    = "Usage: gitmap sf <add|list|rm> [args]"
	MsgSFUsageAdd       = "  gitmap sf add <absolute-path> [--label <text>] [--notes <text>]"
	MsgSFUsageList      = "  gitmap sf list"
	MsgSFUsageRm        = "  gitmap sf rm <absolute-path|id>"
)

// ScanFolder CLI tokens (avoid magic strings).
const (
	SFSubAdd       = "add"
	SFSubList      = "list"
	SFSubListAlias = "ls"
	SFSubRm        = "rm"
	SFSubRmAlias   = "remove"
	SFFlagLabel    = "--label"
	SFFlagNotes    = "--notes"
)
