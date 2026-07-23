package constants

// gitmap:cmd top-level
// Version history CLI commands.
const (
	CmdVersionHistory      = "version-history"
	CmdVersionHistoryAlias = "vh"
)

// Version history help text.
const HelpVersionHistory = "  version-history (vh) Show version transitions for the current repo (--limit N, --json)"

// Version history terminal output.
const (
	MsgVersionHistoryEmpty   = "No version history found for this repo.\n"
	MsgVersionHistoryHeader  = "Version history for %s:\n\n"
	MsgVersionHistoryColumns = "FROM        TO          FOLDER                    TIMESTAMP"
	MsgVersionHistoryRowFmt  = "%-11s %-11s %-25s %s\n"
	MsgVersionHistoryCount   = "\n%d transition(s) recorded.\n"
)

// Version history error messages.
const (
	ErrVersionHistoryCwd = "Error: cannot determine current directory: %v\n"
	ErrVersionHistoryDB  = "Error: failed to query version history: %v\n"
)
