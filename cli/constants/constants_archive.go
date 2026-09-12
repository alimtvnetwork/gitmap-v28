package constants

// gitmap:cmd top-level
// Archive CLI commands and shorthand aliases (Slice — uzc + zip family).
const (
	CmdUnzipCompact      = "unzip-compact"
	CmdUnzipCompactAlias = "uzc"
	CmdZip               = "zip"
	// Note: 'z' is already taken by zip-group. Archive 'zip' is the only
	// surface here — users wanting a shorthand can `gitmap alias add`.
)

// Archive flag names (long + short).
const (
	FlagArchiveBest     = "best"
	FlagArchiveFast     = "fast"
	FlagArchiveStandard = "standard"
	FlagArchiveStdShort = "s"
	FlagArchiveOut      = "out"
	FlagArchiveOutShort = "o"
	FlagArchiveList     = "list"
	FlagArchiveListShrt = "l"
	FlagArchiveInclude  = "include"
	FlagArchiveExclude  = "exclude"
	FlagArchiveYes      = "yes"
)

// Archive flag descriptions.
const (
	FlagDescArchiveBest     = "Use the highest compression level (slowest, smallest)"
	FlagDescArchiveFast     = "Use the fastest compression level (largest)"
	FlagDescArchiveStandard = "Use the default/standard compression level"
	FlagDescArchiveOut      = "Output file path (extension drives format selection)"
	FlagDescArchiveList     = "List archive contents without extracting"
	FlagDescArchiveInclude  = "Comma-separated globs; only matching paths are processed"
	FlagDescArchiveExclude  = "Comma-separated globs; matching paths are skipped"
	FlagDescArchiveYes      = "Skip confirmation prompts"
)

// Compression mode tags persisted in ArchiveHistory.CompressionMode.
const (
	CompressionBest     = "Best"
	CompressionFast     = "Fast"
	CompressionStandard = "Standard"
)

// Archive command names persisted in ArchiveHistory.CommandName.
const (
	ArchiveCmdUnzipCompact = "unzip-compact"
	ArchiveCmdZip          = "zip"
)

// ArchiveHistory.Status values.
const (
	ArchiveStatusSuccess = "Success"
	ArchiveStatusFailed  = "Failed"
)

// SQL: create ArchiveHistory table (PascalCase, INTEGER PK AUTOINCREMENT,
// matches v15 conventions). Stored under .gitmap/db/<profile>.db like
// every other table — no separate per-command schema.
const SQLCreateArchiveHistory = `CREATE TABLE IF NOT EXISTS ArchiveHistory (
	ArchiveHistoryId        INTEGER PRIMARY KEY AUTOINCREMENT,
	CommandName             TEXT    NOT NULL,
	InputSources            TEXT    NOT NULL DEFAULT '[]',
	OutputPath              TEXT    NOT NULL DEFAULT '',
	ArchiveFormat           TEXT    NOT NULL DEFAULT '',
	CompressionMode         TEXT    NOT NULL DEFAULT '',
	UsedTemporaryDirectory  INTEGER NOT NULL DEFAULT 0,
	Status                  TEXT    NOT NULL DEFAULT '',
	ErrorMessage            TEXT    NOT NULL DEFAULT '',
	StartedAt               TEXT    NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CompletedAt             TEXT    NOT NULL DEFAULT ''
)`

// SQL: ArchiveHistory operations.
const (
	SQLInsertArchiveHistory = `INSERT INTO ArchiveHistory
		(CommandName, InputSources, OutputPath, ArchiveFormat, CompressionMode,
		 UsedTemporaryDirectory, Status, ErrorMessage, StartedAt, CompletedAt)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	SQLUpdateArchiveHistoryFinish = `UPDATE ArchiveHistory
		SET OutputPath = ?, ArchiveFormat = ?, UsedTemporaryDirectory = ?,
		    Status = ?, ErrorMessage = ?, CompletedAt = ?
		WHERE ArchiveHistoryId = ?`

	SQLSelectArchiveHistoryRecent = `SELECT ArchiveHistoryId, CommandName, InputSources,
		OutputPath, ArchiveFormat, CompressionMode, UsedTemporaryDirectory, Status,
		ErrorMessage, StartedAt, CompletedAt
		FROM ArchiveHistory ORDER BY ArchiveHistoryId DESC LIMIT ?`

	SQLDropArchiveHistory = `DROP TABLE IF EXISTS ArchiveHistory`
)

// User-facing archive messages.
const (
	MsgArchiveBanner          = "▶ gitmap %s v%s"
	MsgArchiveResolving       = "  ▸ Resolving %d input source(s)…"
	MsgArchiveResolved        = "  ✓ %s → %s"
	MsgArchiveExtractStart    = "  ▸ Extracting %s → %s"
	MsgArchiveExtractDone     = "  ✓ Extracted to %s (%d entries, format=%s)"
	MsgArchiveCompactFlatten  = "  ✓ Flattened %d duplicate-name layer(s)"
	MsgArchiveListHeader      = "  Archive: %s (format=%s, %d entries)"
	MsgArchiveListEntry       = "    %s  %d bytes"
	MsgArchiveCreateStart     = "  ▸ Creating %s (compression=%s, %d source(s))"
	MsgArchiveCreateDone      = "  ✓ Created %s"
	MsgArchiveTempCleanup     = "  ◦ Cleaning temporary workspace %s"
	MsgArchiveAutoPicked      = "  ✓ Auto-detected single archive in current folder: %s"
	MsgArchiveHistoryRecorded = "  ✓ Recorded ArchiveHistory[%d] (%s)"
	WarnArchiveHistoryWrite   = "  ⚠ Could not record ArchiveHistory: %v"
	ErrArchiveNoSource        = "%s: no input source(s) supplied and current folder has no single archive"
	ErrArchiveMultiInCwd      = "%s: current folder has %d archives — pass one explicitly"
	ErrArchiveUnknownFormat   = "%s: could not determine archive format for %q"
	ErrArchiveCreateNeedsOut  = "zip: --out <path> is required (extension drives format)"
	ErrArchiveBadCompression  = "%s: choose at most one of --best / --fast / --standard"
)

// HelpUnzipCompact / HelpZip strings shown by the help generator.
const (
	HelpUnzipCompact = "  unzip-compact (uzc) [src] [dest]  Extract local/URL archive into a single normalized folder"
	HelpZip          = "  zip (z) <src...> --out <file>     Create archive from folders / URLs / git repos"
)

// MaxCompactFlattenLayers caps how many duplicate same-name folder layers
// the compact-extract algorithm collapses. The spec calls for 2..4 layers
// of tolerance; we set the upper bound at 4 to match.
const MaxCompactFlattenLayers = 4
