package constants

// Table name for version history (v15: singular preserved; PK renamed).
const TableRepoVersionHistory = "RepoVersionHistory"

// SQL: create RepoVersionHistory table (v15: RepoVersionHistoryId PK).
// FK references v15 Repo(RepoId).
const SQLCreateRepoVersionHistory = `CREATE TABLE IF NOT EXISTS RepoVersionHistory (
	RepoVersionHistoryId INTEGER PRIMARY KEY AUTOINCREMENT,
	RepoId               INTEGER NOT NULL REFERENCES Repo(RepoId) ON DELETE CASCADE,
	FromVersionTag       TEXT NOT NULL,
	FromVersionNum       INTEGER NOT NULL,
	ToVersionTag         TEXT NOT NULL,
	ToVersionNum         INTEGER NOT NULL,
	FlattenedPath        TEXT DEFAULT '',
	CreatedAt            TEXT DEFAULT CURRENT_TIMESTAMP
)`

// SQL: add version columns to Repo (v15 table name).
const (
	SQLAddCurrentVersionTag = "ALTER TABLE Repo ADD COLUMN CurrentVersionTag TEXT DEFAULT ''"
	SQLAddCurrentVersionNum = "ALTER TABLE Repo ADD COLUMN CurrentVersionNum INTEGER DEFAULT 0"
)

// SQL: version history operations (v15: RepoVersionHistoryId PK).
const (
	SQLInsertVersionHistory = `INSERT INTO RepoVersionHistory
		(RepoId, FromVersionTag, FromVersionNum, ToVersionTag, ToVersionNum, FlattenedPath)
		VALUES (?, ?, ?, ?, ?, ?)`

	SQLSelectVersionHistory = `SELECT RepoVersionHistoryId, RepoId, FromVersionTag, FromVersionNum,
		ToVersionTag, ToVersionNum, FlattenedPath, CreatedAt
		FROM RepoVersionHistory WHERE RepoId = ? ORDER BY CreatedAt DESC`

	SQLUpdateRepoVersion = `UPDATE Repo SET CurrentVersionTag = ?, CurrentVersionNum = ?,
		UpdatedAt = CURRENT_TIMESTAMP WHERE RepoId = ?`

	SQLSelectRepoIDByPath = "SELECT RepoId FROM Repo WHERE AbsolutePath = ?"

	SQLDropRepoVersionHistory = "DROP TABLE IF EXISTS RepoVersionHistory"
)

// Version history error messages.
const ErrDBVersionHistory = "failed to query version history: %v"

// Flatten messages.
const (
	MsgFlattenRemoving  = "Removing existing %s for fresh clone...\n"
	MsgFlattenCloning   = "Cloning %s into %s (flattened)...\n"
	MsgFlattenDone      = "✓ Cloned %s into %s\n"
	MsgFlattenVersionDB = "✓ Recorded version transition v%d -> v%d\n"
)
