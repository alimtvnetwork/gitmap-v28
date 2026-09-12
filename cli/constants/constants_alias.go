package constants

// gitmap:cmd top-level
// Alias command names.
const (
	CmdAlias        = "alias"
	CmdAliasShort   = "a"
	SubCmdAliasSet  = "set"
	SubCmdAliasRm   = "remove"
	SubCmdAliasList = "list"
	SubCmdAliasShow = "show"
	SubCmdAliasSug  = "suggest"
)

// Alias table name (v15: singular).
const TableAlias = "Alias"

// Legacy plural retained for migration detection.
const LegacyTableAliases = "Aliases"

// SQL: create Alias table (v15: singular + AliasId PK). FK references Repo(RepoId).
const SQLCreateAlias = `CREATE TABLE IF NOT EXISTS Alias (
	AliasId   INTEGER PRIMARY KEY AUTOINCREMENT,
	Alias     TEXT NOT NULL UNIQUE,
	RepoId    INTEGER NOT NULL REFERENCES Repo(RepoId) ON DELETE CASCADE,
	CreatedAt TEXT DEFAULT CURRENT_TIMESTAMP
)`

// SQL: alias operations (v15: Alias singular, AliasId PK).
const (
	SQLInsertAlias = `INSERT INTO Alias (Alias, RepoId) VALUES (?, ?)`

	SQLUpdateAlias = `UPDATE Alias SET RepoId = ? WHERE Alias = ?`

	SQLSelectAllAliases = `SELECT a.AliasId, a.Alias, a.RepoId, a.CreatedAt
		FROM Alias a ORDER BY a.Alias`

	SQLSelectAliasByName = `SELECT a.AliasId, a.Alias, a.RepoId, a.CreatedAt
		FROM Alias a WHERE a.Alias = ?`

	SQLSelectAliasByRepoID = `SELECT a.AliasId, a.Alias, a.RepoId, a.CreatedAt
		FROM Alias a WHERE a.RepoId = ?`

	SQLDeleteAlias = `DELETE FROM Alias WHERE Alias = ?`

	SQLSelectAliasWithRepo = `SELECT a.AliasId, a.Alias, a.RepoId, a.CreatedAt,
		r.AbsolutePath, r.Slug
		FROM Alias a JOIN Repo r ON a.RepoId = r.RepoId
		WHERE a.Alias = ?`

	SQLSelectAllAliasesWithRepo = `SELECT a.AliasId, a.Alias, a.RepoId, a.CreatedAt,
		r.AbsolutePath, r.Slug
		FROM Alias a JOIN Repo r ON a.RepoId = r.RepoId
		ORDER BY a.Alias`

	SQLSelectUnaliasedRepos = `SELECT r.RepoId, r.Slug, r.RepoName
		FROM Repo r LEFT JOIN Alias a ON r.RepoId = a.RepoId
		WHERE a.AliasId IS NULL ORDER BY r.Slug`
)

// SQL: drop Alias table (and legacy plural for safety on Reset).
const (
	SQLDropAlias   = "DROP TABLE IF EXISTS Alias"
	SQLDropAliases = "DROP TABLE IF EXISTS Aliases" // legacy
)

// Alias flag descriptions.
const (
	FlagDescAliasApply = "Auto-accept all alias suggestions"
	FlagDescAliasFlag  = "Target a repository by its alias"
)

// Alias messages.
const (
	MsgAliasCreated     = "  ✓ Alias %q → %s\n"
	MsgAliasUpdated     = "  ✓ Updated alias %q → %s\n"
	MsgAliasRemoved     = "  ✓ Removed alias %q\n"
	MsgAliasResolved    = "  → Resolved alias %q → %s (slug: %s)\n"
	MsgAliasSuggest     = "  %-20s → %-10s Accept? (y/N): "
	MsgAliasSuggestDone = "  ✓ Created %d alias(es).\n"
	MsgAliasSuggestNone = "  All repos already have aliases."
	MsgAliasListHeader  = "\n  Aliases (%d):\n\n"
	MsgAliasListRow     = "  %-15s → %s\n"
	MsgAliasConflict    = "  ⚠ Alias %q already points to %s.\n"
	MsgAliasReassign    = "  → Reassign to %s? (y/N): "
	MsgAliasBothWarn    = "  ⚠ Both alias and slug provided — using alias %q.\n"
)

// Alias error messages.
const (
	ErrAliasNotFound    = "no alias found: %s"
	ErrAliasEmpty       = "alias name cannot be empty"
	ErrAliasInvalid     = "alias must be alphanumeric with hyphens: %s"
	ErrAliasShadow      = "alias cannot shadow command: %s"
	ErrAliasCreate      = "failed to create alias: %v"
	ErrAliasQuery       = "failed to query aliases: %v"
	ErrAliasDelete      = "failed to delete alias: %v"
	ErrAliasRepoMissing = "repo not found for alias target: %s"
)
