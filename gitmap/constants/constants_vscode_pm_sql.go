package constants

// SQL for the VSCodeProject table — DB source of truth for the
// VS Code Project Manager sync (v3.38.0+).
//
// `tags` is not stored in the DB on purpose; it lives only in
// `projects.json` and is preserved verbatim across syncs so user edits in
// the extension UI are never clobbered.
//
// `Paths` (multi-root) IS stored in the DB as a JSON-encoded TEXT column
// added in schema v20 (v3.39.0). The DB-side list is UNIONed with whatever
// the user added in the VS Code UI on every sync — gitmap never removes a
// user-added path. Use `gitmap code paths rm` to drop a gitmap-managed
// extra root explicitly.

const TableVSCodeProject = "VSCodeProject"

const SQLCreateVSCodeProject = `CREATE TABLE IF NOT EXISTS VSCodeProject (
	VSCodeProjectId INTEGER PRIMARY KEY AUTOINCREMENT,
	RootPath        TEXT NOT NULL,
	Name            TEXT NOT NULL,
	Paths           TEXT NOT NULL DEFAULT '[]',
	Enabled         INTEGER NOT NULL DEFAULT 1,
	Profile         TEXT NOT NULL DEFAULT '',
	LastSeenAt      TEXT DEFAULT CURRENT_TIMESTAMP,
	CreatedAt       TEXT DEFAULT CURRENT_TIMESTAMP,
	UpdatedAt       TEXT DEFAULT CURRENT_TIMESTAMP
)`

// COLLATE NOCASE so Windows path matching is case-insensitive while
// staying byte-exact on Unix when the user happens to use the same case.
const SQLCreateVSCodeProjectRootPathIndex = `CREATE UNIQUE INDEX IF NOT EXISTS UX_VSCodeProject_RootPath ON VSCodeProject(RootPath COLLATE NOCASE)`

// Idempotent additive migration — safe to re-run on legacy v18/v19 DBs
// that pre-date the Paths column. SQLite's ALTER TABLE ADD COLUMN errors
// when the column already exists, so callers must IGNORE that single error.
const SQLAddVSCodeProjectPathsColumn = `ALTER TABLE VSCodeProject ADD COLUMN Paths TEXT NOT NULL DEFAULT '[]'`

const SQLDropVSCodeProject = "DROP TABLE IF EXISTS VSCodeProject"

const (
	SQLUpsertVSCodeProject = `INSERT INTO VSCodeProject (RootPath, Name, Paths, Enabled, Profile, LastSeenAt, UpdatedAt)
		VALUES (?, ?, '[]', 1, '', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		ON CONFLICT(RootPath) DO UPDATE SET
			Name=excluded.Name,
			LastSeenAt=CURRENT_TIMESTAMP,
			UpdatedAt=CURRENT_TIMESTAMP`

	SQLSelectAllVSCodeProjects = `SELECT VSCodeProjectId, RootPath, Name, Paths, Enabled, Profile, LastSeenAt, CreatedAt, UpdatedAt
		FROM VSCodeProject ORDER BY UpdatedAt DESC, RootPath ASC`

	SQLSelectVSCodeProjectByPath = `SELECT VSCodeProjectId, RootPath, Name, Paths, Enabled, Profile, LastSeenAt, CreatedAt, UpdatedAt
		FROM VSCodeProject WHERE RootPath = ? COLLATE NOCASE`

	SQLSelectVSCodeProjectByName = `SELECT VSCodeProjectId, RootPath, Name, Paths, Enabled, Profile, LastSeenAt, CreatedAt, UpdatedAt
		FROM VSCodeProject WHERE Name = ? COLLATE NOCASE LIMIT 1`

	SQLRenameVSCodeProject = `UPDATE VSCodeProject
		SET Name = ?, UpdatedAt = CURRENT_TIMESTAMP
		WHERE RootPath = ? COLLATE NOCASE`

	SQLUpdateVSCodeProjectPaths = `UPDATE VSCodeProject
		SET Paths = ?, UpdatedAt = CURRENT_TIMESTAMP
		WHERE RootPath = ? COLLATE NOCASE`

	SQLDeleteVSCodeProjectByPath = `DELETE FROM VSCodeProject WHERE RootPath = ? COLLATE NOCASE`
)

// Error messages.
const (
	ErrVSCodePMUpsert        = "failed to upsert VSCodeProject %q: %v"
	ErrVSCodePMList          = "failed to list VSCodeProject rows: %v"
	ErrVSCodePMRename        = "failed to rename VSCodeProject %q: %v"
	ErrVSCodePMDelete        = "failed to delete VSCodeProject %q: %v"
	ErrVSCodePMUpdatePaths   = "failed to update Paths for VSCodeProject %q: %v"
	ErrVSCodePMPathsEncode   = "failed to encode Paths for VSCodeProject %q: %v"
	ErrVSCodePMPathsDecode   = "failed to decode Paths for VSCodeProject %q: %v"
	ErrVSCodePMAliasNotFound = "no VS Code project registered with alias %q (register one first via `gitmap code %s`)"
)

// User-facing messages for the `code paths` subcommand.
const (
	MsgVSCodePMPathsAdded   = "  ✓ added extra path to %q: %s\n"
	MsgVSCodePMPathsRemoved = "  ✓ removed extra path from %q: %s\n"
	MsgVSCodePMPathsExists  = "  • path already attached to %q: %s\n"
	MsgVSCodePMPathsMissing = "  • path not attached to %q: %s\n"
	MsgVSCodePMPathsList    = "%s (%s)\n  rootPath: %s\n  paths   : %s\n"
	MsgVSCodePMPathsNone    = "  (no extra paths)\n"
)
