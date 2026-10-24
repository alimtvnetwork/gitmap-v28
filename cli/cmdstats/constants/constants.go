package constants

// Command names and aliases for cmdstats.
const (
	CmdStats      = "stats"
	CmdStatsAlias = "ss"
)

// Flag names used in FlagSet parsing.
const (
	FlagCommand      = "command"
	FlagCommandShort = "c"
	FlagJSON         = "json"
	FlagJSONShort    = "j"
	FlagHelp         = "help"
	FlagHelpShort    = "h"
)

// Flag CLI argument strings.
const (
	FlagArgCommand      = "--command"
	FlagArgCommandShort = "-c"
	FlagArgJSON         = "--json"
	FlagArgJSONShort    = "-j"
	FlagArgHelp         = "--help"
	FlagArgHelpShort    = "-h"
)

// Flag descriptions for help text and flag definitions.
const (
	FlagDescStatsCommand = "Show stats for a specific command only"
	FlagDescJSON         = "Output stats in structured JSON format"
	FlagDescHelp         = "Show this stats help menu"
)

// Wire contract keys for top-level stats JSON representation.
const (
	WireKeyTotalCommands   = "totalCommands"
	WireKeyUniqueCommands  = "uniqueCommands"
	WireKeyTotalSuccess    = "totalSuccess"
	WireKeyTotalFail       = "totalFail"
	WireKeyOverallFailRate = "overallFailRate"
	WireKeyAvgDurationMs   = "avgDurationMs"
	WireKeyCommands        = "commands"
)

// Wire contract keys for per-command stats JSON rows.
const (
	WireKeyCmdCommand       = "command"
	WireKeyCmdTotalRuns     = "totalRuns"
	WireKeyCmdSuccessCount  = "successCount"
	WireKeyCmdFailCount     = "failCount"
	WireKeyCmdFailRate      = "failRate"
	WireKeyCmdAvgDurationMs = "avgDurationMs"
	WireKeyCmdMinDurationMs = "minDurationMs"
	WireKeyCmdMaxDurationMs = "maxDurationMs"
	WireKeyCmdLastUsed      = "lastUsed"
)

// SSH forward candidate subcommand keywords.
const (
	SSHForwardFixAuth      = "fix-auth"
	SSHForwardLogin        = "login"
	SSHForwardHost         = "host"
	SSHForwardHosts        = "hosts"
	SSHForwardNodes        = "nodes"
	SSHForwardEnable       = "enable"
	SSHForwardPort         = "port"
	SSHForwardTroubleshoot = "troubleshoot"
	SSHForwardKey          = "key"
	SSHForwardCopyID       = "copy-id"
	SSHForwardAuthKey      = "auth-key"
)

// AllSSHForwardCandidates returns all subcommand names forwarded to SSH.
func AllSSHForwardCandidates() []string {
	return []string{
		SSHForwardFixAuth, SSHForwardLogin, SSHForwardHost,
		SSHForwardHosts, SSHForwardNodes, SSHForwardEnable,
		SSHForwardPort, SSHForwardTroubleshoot, SSHForwardKey,
		SSHForwardCopyID, SSHForwardAuthKey,
	}
}

// Error codes, diagnostic operations, and pattern substrings.
const (
	ErrCodeE1154           = "E1154"
	ErrCodeE1155           = "E1155"
	ErrCodeE1156           = "E1156"
	OpStatsOpenDB          = "cmd.stats.openDB"
	OpStatsLegacyData      = "cmd.stats.legacyData"
	OpStatsQuery           = "cmd.stats.query"
	DomainStats            = "cmd.stats"
	ErrPatternScanError    = "Scan error"
	ErrPatternDriverValue  = "converting driver.Value type string"
	ErrStatsJSONMarshalFmt = "  ✗ Failed to marshal stats to JSON: %v\n"
	SchemaStatsFilename    = "stats.schema.json"
)

// Exit codes for cmdstats execution.
const (
	ExitSuccess      = 0
	ExitGeneralError = 1
)
