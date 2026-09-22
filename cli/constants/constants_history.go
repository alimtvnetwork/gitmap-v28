package constants

// History table name (v15: singular preserved; PK Id → CommandHistoryId).
const TableCommandHistory = "CommandHistory"

// SQL: create CommandHistory table (v15: CommandHistoryId PK).
const SQLCreateCommandHistory = `CREATE TABLE IF NOT EXISTS CommandHistory (
	CommandHistoryId INTEGER PRIMARY KEY AUTOINCREMENT,
	Command          TEXT NOT NULL,
	Alias            TEXT DEFAULT '',
	Args             TEXT DEFAULT '',
	Flags            TEXT DEFAULT '',
	StartedAt        TEXT NOT NULL,
	FinishedAt       TEXT DEFAULT '',
	DurationMs       INTEGER DEFAULT 0,
	ExitCode         INTEGER DEFAULT 0,
	Summary          TEXT DEFAULT '',
	RepoCount        INTEGER DEFAULT 0,
	CreatedAt        TEXT DEFAULT CURRENT_TIMESTAMP
)`

// SQL: command history operations (v15: CommandHistoryId PK).
const (
	SQLInsertHistory = `INSERT INTO CommandHistory
		(Command, Alias, Args, Flags, StartedAt, FinishedAt, DurationMs, ExitCode, Summary, RepoCount)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	SQLUpdateHistory = `UPDATE CommandHistory
		SET FinishedAt = ?, DurationMs = ?, ExitCode = ?, Summary = ?, RepoCount = ?
		WHERE CommandHistoryId = ?`

	SQLSelectAllHistory = `SELECT CommandHistoryId, Command, Alias, Args, Flags, StartedAt, FinishedAt,
		DurationMs, ExitCode, Summary, RepoCount, CreatedAt
		FROM CommandHistory ORDER BY CreatedAt DESC`

	SQLSelectHistoryByCommand = `SELECT CommandHistoryId, Command, Alias, Args, Flags, StartedAt, FinishedAt,
		DurationMs, ExitCode, Summary, RepoCount, CreatedAt
		FROM CommandHistory WHERE Command = ? ORDER BY CreatedAt DESC`

	SQLDeleteAllHistory = "DELETE FROM CommandHistory"

	SQLDropCommandHistory = "DROP TABLE IF EXISTS CommandHistory"
)

// gitmap:cmd top-level
// History CLI commands.
const (
	CmdHistory           = "history"
	CmdHistoryAlias      = "hi"
	CmdHistoryReset      = "history-reset"
	CmdHistoryResetAlias = "hr"
)

// History help text.
const (
	HelpHistory      = "  history (hi)        Show command execution audit log (--limit N, --json, --command, --detail)"
	HelpHistoryReset = "  history-reset (hr)  Clear command history (--confirm required)"
)

// History flag descriptions.
const (
	FlagDescDetail  = "Detail level: basic, standard, or detailed (default: standard)"
	FlagDescCommand = "Filter by command name"
)

// History detail levels.
const (
	DetailBasic    = "basic"
	DetailStandard = "standard"
	DetailDetailed = "detailed"
)

// History terminal columns.
const (
	MsgHistoryColumnsBasic    = "COMMAND         TIMESTAMP                STATUS"
	MsgHistoryColumnsStandard = "COMMAND         TIMESTAMP                FLAGS                    STATUS  DURATION"
	MsgHistoryColumnsDetailed = "COMMAND         TIMESTAMP                ARGS             FLAGS                    STATUS  DURATION  REPOS  SUMMARY"
	MsgHistoryRowBasicFmt     = "%-15s %-24s %s\n"
	MsgHistoryRowStdFmt       = "%-15s %-24s %-24s %-7s %s\n"
	MsgHistoryRowDetailFmt    = "%-15s %-24s %-16s %-24s %-7s %-9s %-6s %s\n"
)

// History messages.
const (
	MsgHistoryEmpty          = "No command history found.\n"
	MsgHistoryResetDone      = "Command history cleared.\n"
	ErrHistoryResetFailed    = "failed to reset command history: %v\n"
	ErrHistoryResetNoConfirm = "history-reset requires --confirm flag\n"
	ErrHistoryQuery          = "failed to query command history: %v"
	MsgHistoryStatusOK       = "OK"
	MsgHistoryStatusFail     = "FAIL"
)
