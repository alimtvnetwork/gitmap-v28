package constants

// Subcommand names and aliases for cmdhealth.
const (
	CmdHealth      = "health"
	CmdHealthAlias = "hlth"
)

// Environment variable names for configuration.
const (
	EnvGitmapWorkdir = "GITMAP_WORKDIR"
)

// Flag names used in FlagSet parsing.
const (
	FlagWorkdir = "workdir"
	FlagHelp    = "help"
)

// CLI flag argument strings.
const (
	FlagArgWorkdir   = "--workdir"
	FlagArgHelp      = "--help"
	FlagArgHelpShort = "-h"
)

// Flag descriptions for help menu rendering.
const (
	FlagDescWorkdir = "repo workdir"
)

// Default directory names and fallback paths.
const (
	DefaultWorkspaceDir = "workspace"
	DefaultReposDir     = "repos"
	FallbackWorkdir     = "."
)

// Wrap tokens for error context wrapping.
const (
	WrapHealthFlags = "health.flags"
)

// Thresholds and limits for health checks and rollups.
const (
	MaxDirtyReposDisplay   = 5
	MaxRedPipelinesDisplay = 5
	MaxPipelineProbedRepos = 50
	StaleDbDaysCutoff      = 30
	HoursPerDay            = 24
)

// GitHub CLI command names, arguments, and JSON tokens.
const (
	CmdGh                 = "gh"
	ArgGhRun              = "run"
	ArgGhList             = "list"
	FlagGhRepo            = "-R"
	FlagGhLimit           = "--limit"
	ValGhLimitOne         = "1"
	FlagGhJSON            = "--json"
	ValGhConclusion       = "conclusion"
	JsonConclusionFailure = "\"failure\""
	PathSeparator         = "/"
)

// Diagnostic notes and advisory strings.
const (
	NoteGhNotInstalled          = "gh not installed; pipeline check skipped"
	NoteWorkdirUnreadable       = "workdir unreadable; pipeline check skipped"
	NoteNodeRegistryUnreachable = "node registry unreachable"
	NoteFleetSkewHint           = "run `gitmap nodes version` for skew"
)

// UI layout strings, section titles, and report formatting.
const (
	HeaderWorkstationHealth     = "workstation health:"
	FmtReposTotal               = "  repos:     %d (%d dirty"
	FmtDirtyItem                = ": %s"
	ParenClose                  = ")"
	FmtPipelinesNote            = "  pipelines: %s\n"
	FmtPipelinesFailing         = "  pipelines: %d failing"
	MsgFleetRegistryUnreachable = "  fleet:     registry unreachable"
	FmtFleetRegistered          = "  fleet:     %d node(s) registered (%s)\n"
	FmtDbsStale                 = "  dbs:       %d stale"
	FmtWorstLine                = "  worst:     %s\n"
	ItemSeparator               = ", "
)

// Worst item recommendations and follow-up guidance.
const (
	WorstTriageFmt            = "pe triage -> %s (%d failing)"
	WorstDbDoctorRepair       = "db doctor --repair (stale databases)"
	WorstDirtyWorkingTreesFmt = "dirty working trees (%d) — review before pull/push"
	WorstNothingUrgent        = "nothing urgent — all checks green"
)

// HelpTextHealth provides usage instructions and flags for the health command.
const HelpTextHealth = `health — daily workstation health aggregation.

Description:
  Read-only rollup of the morning checks: dirty repos, failing
  pipelines, fleet node count, and stale databases. Points at the
  worst item first with the command that digs deeper.

Usage:
  gitmap health [--workdir <dir>]
  gitmap hlth [--workdir <dir>]

Flags:
  --workdir <dir>   repo workdir (default: $GITMAP_WORKDIR or ~/workspace/repos)
  --help, -h        show this help

Example:
  gitmap health
`
