package constants

// SQL: create CloneInteractiveSelection table (spec 100, v3.153.0+).
//
// Records every `gitmap clone-pick` invocation so the same selection
// can be re-applied with `--replay <id|name>` -- without forcing the
// user to re-type the path list or remember which sparse-checkout
// flags they used.
//
// Why no FK to Repo(RepoId):
//
//	The repo being picked may not exist in any prior `gitmap scan`
//	output. A user can `clone-pick owner/repo docs` against a fresh
//	clone they've never indexed locally; gating that behind a
//	"must scan first" requirement would gut the command's usefulness.
//	We instead identify the repo by its canonical id (host/owner/repo)
//	stored as text and indexed.
//
// Why no UNIQUE on Name:
//
//	'' is the default Name (auto-saves without --name). SQLite UNIQUE
//	would reject the second nameless save. Uniqueness for non-empty
//	names is enforced in the store layer (a SELECT-then-INSERT race
//	is acceptable here -- worst case the second user sees an error
//	when their `--replay <name>` lookup returns >1 row, which is the
//	MsgClonePickReplayAmbiguous path and prints the candidate IDs).
const SQLCreateCloneInteractiveSelection = `CREATE TABLE IF NOT EXISTS CloneInteractiveSelection (
	SelectionId       INTEGER PRIMARY KEY AUTOINCREMENT,
	Name              TEXT NOT NULL DEFAULT '',
	RepoCanonicalId   TEXT NOT NULL,
	RepoUrl           TEXT NOT NULL,
	Mode              TEXT NOT NULL DEFAULT 'https',
	Branch            TEXT NOT NULL DEFAULT '',
	Depth             INTEGER NOT NULL DEFAULT 1,
	Cone              INTEGER NOT NULL DEFAULT 1,
	KeepGit           INTEGER NOT NULL DEFAULT 1,
	DestDir           TEXT NOT NULL DEFAULT '.',
	PathsCsv          TEXT NOT NULL,
	UsedAsk           INTEGER NOT NULL DEFAULT 0,
	CreatedAt         TEXT DEFAULT CURRENT_TIMESTAMP
)`

// SQLCreateClonePickRepoCanonIndex speeds up "show me everything I've
// ever picked from this repo" lookups (used implicitly when --replay
// resolves a numeric id within the matching repo's history).
const SQLCreateClonePickRepoCanonIndex = `CREATE INDEX IF NOT EXISTS idx_clonepick_repocanon
	ON CloneInteractiveSelection(RepoCanonicalId)`

// SQLCreateClonePickNameIndex is partial: empty Names share no
// uniqueness or lookup pressure, so the index is cheaper.
const SQLCreateClonePickNameIndex = `CREATE INDEX IF NOT EXISTS idx_clonepick_name
	ON CloneInteractiveSelection(Name) WHERE Name <> ''`

// CRUD statements. All column lists are spelled out (no SELECT *) so
// future column adds don't silently break the row scanner.
const (
	SQLInsertClonePickSelection = `INSERT INTO CloneInteractiveSelection
		(Name, RepoCanonicalId, RepoUrl, Mode, Branch, Depth, Cone, KeepGit, DestDir, PathsCsv, UsedAsk)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	SQLSelectClonePickByID = `SELECT SelectionId, Name, RepoCanonicalId, RepoUrl, Mode,
		Branch, Depth, Cone, KeepGit, DestDir, PathsCsv, UsedAsk, CreatedAt
		FROM CloneInteractiveSelection WHERE SelectionId = ?`

	SQLSelectClonePickByName = `SELECT SelectionId, Name, RepoCanonicalId, RepoUrl, Mode,
		Branch, Depth, Cone, KeepGit, DestDir, PathsCsv, UsedAsk, CreatedAt
		FROM CloneInteractiveSelection WHERE Name = ? ORDER BY SelectionId DESC`

	SQLTouchClonePickCreatedAt = `UPDATE CloneInteractiveSelection
		SET CreatedAt = CURRENT_TIMESTAMP WHERE SelectionId = ?`
)
