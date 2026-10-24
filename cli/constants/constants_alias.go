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

// AliasActionType defines actions for the alias subcommand.
type AliasActionType string

// Backwards-compatible type alias.
type AliasAction = AliasActionType

const (
	AliasActionInvalid AliasActionType = ""
	AliasActionSet     AliasActionType = "set"
	AliasActionRemove  AliasActionType = "remove"
	AliasActionRm      AliasActionType = "rm"
	AliasActionList    AliasActionType = "list"
	AliasActionLs      AliasActionType = "ls"
	AliasActionShow    AliasActionType = "show"
	AliasActionSuggest AliasActionType = "suggest"
)

// String returns the canonical serialized string representation.
func (a AliasActionType) String() string {
	return string(a)
}

// Label returns the display label.
func (a AliasActionType) Label() string {
	return a.String()
}

// IsValid reports whether the alias action is recognized.
func (a AliasActionType) IsValid() bool {
	hasAction := a.IsSet() || a.IsRemove() || a.IsList() || a.IsShow() || a.IsSuggest()

	return hasAction
}

func (a AliasActionType) IsSet() bool {
	return a == AliasActionSet
}

func (a AliasActionType) IsRemove() bool {
	return a == AliasActionRemove || a == AliasActionRm
}

func (a AliasActionType) IsList() bool {
	return a == AliasActionList || a == AliasActionLs
}

func (a AliasActionType) IsShow() bool {
	return a == AliasActionShow
}

func (a AliasActionType) IsSuggest() bool {
	return a == AliasActionSuggest
}

// Alias table name (v15: singular).
const TableAlias = "Alias"

// Legacy plural retained for migration detection.
const LegacyTableAliases = "Aliases"

// SQL: create Alias table (v15: singular + AliasId PK). FK references Repo(RepoId).
const SQLCreateAlias = `CREATE TABLE IF NOT EXISTS Alias (
	AliasId   INTEGER PRIMARY KEY AUTOINCREMENT,
	Alias     TEXT NOT NULL UNIQUE,
	RepoId    INTEGER NOT NULL REFERENCES Repo(RepoId) ON DELETE CASCADE,
	IsPrimary INTEGER DEFAULT 0,
	Source    TEXT DEFAULT 'manual',
	CreatedAt TEXT DEFAULT CURRENT_TIMESTAMP
)`

// SQL: alias operations (v15: Alias singular, AliasId PK).
const (
	SQLInsertAlias     = `INSERT INTO Alias (Alias, RepoId) VALUES (?, ?)`
	SQLInsertAliasFull = `INSERT INTO Alias (Alias, RepoId, IsPrimary, Source) VALUES (?, ?, ?, ?)`

	SQLAddAliasIsPrimary = `ALTER TABLE Alias ADD COLUMN IsPrimary INTEGER DEFAULT 0`
	SQLAddAliasSource    = `ALTER TABLE Alias ADD COLUMN Source TEXT DEFAULT 'manual'`

	SQLUpdateAlias = `UPDATE Alias SET RepoId = ? WHERE Alias = ?`

	SQLSelectAllAliases = `SELECT a.AliasId, a.Alias, a.RepoId, a.CreatedAt
		FROM Alias a ORDER BY a.Alias`

	SQLSelectAliasByName = `SELECT a.AliasId, a.Alias, a.RepoId, a.CreatedAt
		FROM Alias a WHERE a.Alias = ?`

	SQLSelectAliasByRepoID = `SELECT a.AliasId, a.Alias, a.RepoId, a.CreatedAt
		FROM Alias a WHERE a.RepoId = ?`

	SQLSelectAliasesByRepoID = `SELECT a.AliasId, a.Alias, a.RepoId, COALESCE(a.IsPrimary, 0), COALESCE(a.Source, 'manual'), a.CreatedAt
		FROM Alias a WHERE a.RepoId = ? ORDER BY a.IsPrimary DESC, a.AliasId ASC`

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
	ErrAliasCreate      = "failed to create alias: %w"
	ErrAliasQuery       = "failed to query aliases: %w"
	ErrAliasDelete      = "failed to delete alias: %w"
	ErrAliasRepoMissing = "repo not found for alias target: %s"
)
