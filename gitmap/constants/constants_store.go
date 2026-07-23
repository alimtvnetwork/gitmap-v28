package constants

// Database location.
const (
	DBDir  = "data"
	DBFile = "gitmap.db"
)

// Lock file.
const (
	LockFileName       = "gitmap.lock"
	LockFilePermission = 0o644
	ErrLockHeld        = "another gitmap process is running (PID %d).\n  If incorrect, delete: %s"
)

// Table names (v15: PascalCase + singular + {Table}Id PK).
const (
	TableRepo      = "Repo"
	TableGroup     = "Group"
	TableGroupRepo = "GroupRepo"
	TableRelease   = "Release"
)

// Legacy table names retained only for migration detection (do not use in new SQL).
const (
	LegacyTableRepos      = "Repos"
	LegacyTableGroups     = "Groups"
	LegacyTableGroupRepos = "GroupRepos"
	LegacyTableReleases   = "Releases"
)

// SQL: create Repo table (v15: singular + RepoId PK).
const SQLCreateRepo = `CREATE TABLE IF NOT EXISTS Repo (
	RepoId           INTEGER PRIMARY KEY AUTOINCREMENT,
	Slug             TEXT NOT NULL,
	RepoName         TEXT NOT NULL,
	HttpsUrl         TEXT NOT NULL,
	SshUrl           TEXT NOT NULL,
	Branch           TEXT NOT NULL,
	RelativePath     TEXT NOT NULL,
	AbsolutePath     TEXT NOT NULL,
	CloneInstruction TEXT NOT NULL,
	Notes            TEXT DEFAULT '',
	LastInjectedDesktopAt TEXT DEFAULT '',
	LastInjectedVSCodeAt  TEXT DEFAULT '',
	LastClonedAt     TEXT DEFAULT '',
	IdentifiedTransport TEXT NOT NULL DEFAULT '',
	CreatedAt        TEXT DEFAULT CURRENT_TIMESTAMP,
	UpdatedAt        TEXT DEFAULT CURRENT_TIMESTAMP
)`

// SQL: migration 008 — backfill IdentifiedTransport column on
// existing Repo tables. Idempotent via PRAGMA table_info detect in
// migrateRepoIdentifiedTransport. Empty string is the sentinel for
// "unknown — backfill from origin URL on next scan" so we can ship
// the column without forcing a one-shot scan on upgrade.
// See spec/04-generic-cli/04-reclone-honors-stored-transport.md.
const SQLAddRepoIdentifiedTransport = `ALTER TABLE Repo ADD COLUMN IdentifiedTransport TEXT NOT NULL DEFAULT ''`

// SQL: create Group table (v15 singular). "Group" is a SQL reserved word so
// it MUST be double-quoted everywhere it appears in DDL/DML.
const SQLCreateGroup = `CREATE TABLE IF NOT EXISTS "Group" (
	GroupId     INTEGER PRIMARY KEY AUTOINCREMENT,
	Name        TEXT NOT NULL UNIQUE,
	Description TEXT DEFAULT '',
	Color       TEXT DEFAULT '',
	CreatedAt   TEXT DEFAULT CURRENT_TIMESTAMP
)`

// SQL: create GroupRepo join table (v15: singular). FKs reference v15 PKs.
const SQLCreateGroupRepo = `CREATE TABLE IF NOT EXISTS GroupRepo (
	GroupId INTEGER NOT NULL REFERENCES "Group"(GroupId) ON DELETE CASCADE,
	RepoId  INTEGER NOT NULL REFERENCES Repo(RepoId) ON DELETE CASCADE,
	PRIMARY KEY (GroupId, RepoId)
)`

// SQL: create Release table (v17: + RepoId FK to Repo, composite UNIQUE(RepoId, Tag)).
// See spec/04-generic-cli/24-release-repo-relationship.md for the rationale.
const SQLCreateRelease = `CREATE TABLE IF NOT EXISTS Release (
	ReleaseId    INTEGER PRIMARY KEY AUTOINCREMENT,
	RepoId       INTEGER NOT NULL REFERENCES Repo(RepoId) ON DELETE CASCADE,
	Version      TEXT NOT NULL,
	Tag          TEXT NOT NULL,
	Branch       TEXT NOT NULL,
	SourceBranch TEXT NOT NULL,
	CommitSha    TEXT NOT NULL,
	Changelog    TEXT DEFAULT '',
	Notes        TEXT DEFAULT '',
	IsDraft      INTEGER DEFAULT 0,
	IsPreRelease INTEGER DEFAULT 0,
	IsLatest     INTEGER DEFAULT 0,
	Source       TEXT DEFAULT 'release',
	CreatedAt    TEXT DEFAULT CURRENT_TIMESTAMP,
	UNIQUE (RepoId, Tag)
)`

// SQL: index on Release.RepoId for fast per-repo filtering (future multi-repo mode).
const SQLCreateReleaseRepoIdIndex = "CREATE INDEX IF NOT EXISTS IdxRelease_RepoId ON Release(RepoId)"

// SQL: add Source column — v15: now targets singular Release table.
const SQLAddSourceColumn = "ALTER TABLE Release ADD COLUMN Source TEXT DEFAULT 'release'"

// SQL: enable foreign keys.
const SQLEnableFK = "PRAGMA foreign_keys = ON"

// SQL: WAL journal + relaxed sync for fast migrations on Windows.
// Default `synchronous=FULL` calls FlushFileBuffers on every commit,
// which serializes parallel tests behind kernel fsync and triggers
// 10-min CI timeouts. WAL + NORMAL is the SQLite-recommended default
// for application-level durability and is safe across crashes (only
// risk is losing the last in-flight commit on power loss — fine for
// gitmap's local cache role).
const (
	SQLPragmaJournalWAL     = "PRAGMA journal_mode = WAL"
	SQLPragmaSynchronousNor = "PRAGMA synchronous = NORMAL"
	SQLPragmaBusyTimeout5s  = "PRAGMA busy_timeout = 5000"
)

// SQL: repo operations (v15: Repo table, RepoId PK).
const (
	SQLUpsertRepo = `INSERT INTO Repo (Slug, RepoName, HttpsUrl, SshUrl, Branch, RelativePath, AbsolutePath, CloneInstruction, Notes, IdentifiedTransport)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(AbsolutePath) DO UPDATE SET
			Slug=excluded.Slug, RepoName=excluded.RepoName, HttpsUrl=excluded.HttpsUrl,
			SshUrl=excluded.SshUrl, Branch=excluded.Branch, RelativePath=excluded.RelativePath,
			CloneInstruction=excluded.CloneInstruction, Notes=excluded.Notes,
			IdentifiedTransport=excluded.IdentifiedTransport, UpdatedAt=CURRENT_TIMESTAMP`

	SQLSelectAllRepos = "SELECT RepoId, Slug, RepoName, HttpsUrl, SshUrl, Branch, RelativePath, AbsolutePath, CloneInstruction, Notes, IdentifiedTransport FROM Repo ORDER BY Slug"

	SQLSelectRepoBySlug = "SELECT RepoId, Slug, RepoName, HttpsUrl, SshUrl, Branch, RelativePath, AbsolutePath, CloneInstruction, Notes, IdentifiedTransport FROM Repo WHERE Slug = ?"

	SQLSelectRepoByPath = "SELECT RepoId, Slug, RepoName, HttpsUrl, SshUrl, Branch, RelativePath, AbsolutePath, CloneInstruction, Notes, IdentifiedTransport FROM Repo WHERE AbsolutePath = ?"

	SQLDeleteRepoByPath = "DELETE FROM Repo WHERE AbsolutePath = ?"
	SQLDeleteRepoBySlug = "DELETE FROM Repo WHERE Slug = ?"
)

// SQL: upsert by AbsolutePath (spec requirement).
const SQLUpsertRepoByPath = `INSERT INTO Repo (Slug, RepoName, HttpsUrl, SshUrl, Branch, RelativePath, AbsolutePath, CloneInstruction, Notes, IdentifiedTransport)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(AbsolutePath) DO UPDATE SET
		Slug=excluded.Slug, RepoName=excluded.RepoName, HttpsUrl=excluded.HttpsUrl,
		SshUrl=excluded.SshUrl, Branch=excluded.Branch, RelativePath=excluded.RelativePath,
		CloneInstruction=excluded.CloneInstruction, Notes=excluded.Notes,
		IdentifiedTransport=excluded.IdentifiedTransport, UpdatedAt=CURRENT_TIMESTAMP`

// SQL: create unique index on AbsolutePath for upsert-by-path (v15: IdxRepo_AbsolutePath).
const SQLCreateAbsPathIndex = "CREATE UNIQUE INDEX IF NOT EXISTS IdxRepo_AbsolutePath ON Repo(AbsolutePath)"

// SQL: drop the legacy index name from pre-v15 installs.
const SQLDropLegacyAbsPathIndex = "DROP INDEX IF EXISTS idx_Repos_AbsolutePath"

// SQL: group operations (v15: "Group" singular, GroupId PK).
const (
	SQLInsertGroup = `INSERT INTO "Group" (Name, Description, Color) VALUES (?, ?, ?)`

	SQLSelectAllGroups = `SELECT GroupId, Name, Description, Color, CreatedAt FROM "Group" ORDER BY Name`

	SQLSelectGroupByName = `SELECT GroupId, Name, Description, Color, CreatedAt FROM "Group" WHERE Name = ?`

	SQLDeleteGroup = `DELETE FROM "Group" WHERE Name = ?`

	SQLInsertGroupRepo = "INSERT OR IGNORE INTO GroupRepo (GroupId, RepoId) VALUES (?, ?)"

	SQLDeleteGroupRepo = "DELETE FROM GroupRepo WHERE GroupId = ? AND RepoId = ?"

	SQLSelectGroupRepos = `SELECT r.RepoId, r.Slug, r.RepoName, r.HttpsUrl, r.SshUrl, r.Branch,
		r.RelativePath, r.AbsolutePath, r.CloneInstruction, r.Notes
		FROM Repo r JOIN GroupRepo gr ON r.RepoId = gr.RepoId WHERE gr.GroupId = ? ORDER BY r.Slug`

	SQLCountGroupRepos = "SELECT COUNT(*) FROM GroupRepo WHERE GroupId = ?"
)

// SQL: import-side group insert (used by store/import.go to insert without conflict).
const SQLImportInsertGroup = `INSERT OR IGNORE INTO "Group" (Name, Description, Color) VALUES (?, ?, ?)`

// SQL: release operations (v17: composite UNIQUE(RepoId, Tag), RepoId FK).
const (
	SQLUpsertRelease = `INSERT INTO Release (RepoId, Version, Tag, Branch, SourceBranch, CommitSha, Changelog, Notes, IsDraft, IsPreRelease, IsLatest, Source, CreatedAt)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(RepoId, Tag) DO UPDATE SET
			Version=excluded.Version, Branch=excluded.Branch, SourceBranch=excluded.SourceBranch,
			CommitSha=excluded.CommitSha, Changelog=excluded.Changelog, Notes=excluded.Notes,
			IsDraft=excluded.IsDraft, IsPreRelease=excluded.IsPreRelease,
			IsLatest=excluded.IsLatest, Source=excluded.Source`

	SQLSelectAllReleases = `SELECT ReleaseId, RepoId, Version, Tag, Branch, SourceBranch, CommitSha, Changelog, Notes, IsDraft, IsPreRelease, IsLatest, Source, CreatedAt
		FROM Release ORDER BY CreatedAt DESC`

	// SQLSelectAllReleasesAcrossRepos powers `gitmap releases --all-repos`.
	// The JOIN on Repo.RepoId = Release.RepoId is what exercises (and
	// pre-pays) the IdxRelease_RepoId secondary index added in v17.
	SQLSelectAllReleasesAcrossRepos = `SELECT R.ReleaseId, R.RepoId, P.Slug, P.AbsolutePath,
			R.Version, R.Tag, R.Branch, R.CommitSha, R.Source,
			R.IsDraft, R.IsLatest, R.IsPreRelease, R.CreatedAt
		FROM Release R
		JOIN Repo P ON P.RepoId = R.RepoId
		ORDER BY R.CreatedAt DESC, P.Slug ASC`

	SQLSelectReleaseByTag = `SELECT ReleaseId, RepoId, Version, Tag, Branch, SourceBranch, CommitSha, Changelog, Notes, IsDraft, IsPreRelease, IsLatest, Source, CreatedAt
		FROM Release WHERE Tag = ?`

	SQLClearLatestRelease = "UPDATE Release SET IsLatest = 0 WHERE IsLatest = 1 AND RepoId = ?"

	SQLAddNotesColumn = "ALTER TABLE Release ADD COLUMN Notes TEXT DEFAULT ''"
)

// SQL: reset operations (v15 names + legacy plurals kept for safe drop on upgraded DBs).
const (
	SQLDropGroupRepo  = "DROP TABLE IF EXISTS GroupRepo"
	SQLDropGroupRepos = "DROP TABLE IF EXISTS GroupRepos" // legacy
	SQLDropGroup      = `DROP TABLE IF EXISTS "Group"`
	SQLDropGroups     = "DROP TABLE IF EXISTS Groups" // legacy
	SQLDropRepo       = "DROP TABLE IF EXISTS Repo"
	SQLDropRepos      = "DROP TABLE IF EXISTS Repos" // legacy, kept for migrateLegacyIDs
	SQLDropRelease    = "DROP TABLE IF EXISTS Release"
	SQLDropReleases   = "DROP TABLE IF EXISTS Releases" // legacy
)

// Store error messages.
const (
	ErrDBOpen          = "failed to open database at %s: %v (operation: open)"
	ErrDBMigrate       = "failed to initialize tables: %v"
	ErrDBUpsert        = "failed to upsert repo: %v"
	ErrDBQuery         = "failed to query repos: %v"
	ErrDBNoMatch       = "no repo matches slug: %s\n"
	ErrDBCreateDir     = "failed to create database directory at %s: %v (operation: mkdir)"
	ErrDBGroupCreate   = "failed to create group: %v"
	ErrDBGroupQuery    = "failed to query groups: %v"
	ErrDBGroupAdd      = "failed to add repo to group: %v"
	ErrDBGroupRemove   = "failed to remove repo from group: %v"
	ErrDBGroupDelete   = "failed to delete group: %v"
	ErrDBGroupNone     = "no group found: %s"
	ErrDBGroupExists   = "group already exists: %s"
	ErrDBReleaseUpsert = "failed to upsert release: %v"
	ErrDBReleaseQuery  = "failed to query releases: %v"
)

// Phase 1 v15 migration messages.
const (
	MsgV15RepoMigrationStart = "→ Migrating database to v15 schema (Repos → Repo)..."
	MsgV15RepoMigrationDone  = "✓ Migrated Repos → Repo (RepoId PK). Existing data preserved."
	ErrV15RepoMigration      = "v15 Repo migration failed: %v"
	ErrV15RepoCountMismatch  = "v15 Repo migration count mismatch: old=%d new=%d"
	ErrV15Phase2Migration    = "v15 Phase 1.2 migration failed: %v"
	ErrV15Phase3Migration    = "v15 Phase 1.3 migration failed: %v"
	ErrV15Phase4Migration    = "v15 Phase 1.4 migration failed: %v"
	ErrV15Phase5Migration    = "v15 Phase 1.5 migration failed: %v"
	ErrV15Phase6Migration    = "v15 Phase 1.6 migration failed: %v"
)

// Phase 1.6 (v17) Release.RepoId FK migration messages.
const (
	MsgV15Phase6Start = "→ Migrating Release table: adding RepoId FK to Repo (v3.17.0 schema)..."
	MsgV15Phase6Wipe  = "  Existing Release rows will be wiped and re-imported from .gitmap/release/v*.json on next list-releases."
	MsgV15Phase6Done  = "✓ Release table dropped — will be recreated with RepoId FK on next CREATE pass."
	ErrReleaseNoRepo  = "cannot persist release: no repo registered for path %q. Run `gitmap scan` first to register the repo, then retry."
)
