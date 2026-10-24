package constants

// Subcommand identifiers for cmdplans.
const (
	CmdPlans      = "plans"
	CmdPlansAlias = "plan"
)

// Flag literals for cmdplans execution options.
const (
	FlagLimitShort     = "-n"
	FlagLimit          = "--limit"
	FlagAllShort       = "-a"
	FlagAll            = "--all"
	FlagPendingShort   = "-p"
	FlagPending        = "--pending"
	FlagCompletedShort = "-c"
	FlagCompleted      = "--completed"
	FlagJSONShort      = "-j"
	FlagJSON           = "--json"
	FlagYAML           = "--yaml"
	FlagYML            = "--yml"
	FlagHelpShort      = "-h"
	FlagHelp           = "--help"
	FlagStatus         = "--status"
	FlagFormat         = "--format"
	FlagDir            = "--dir"
)

// Directory and path tokens for repository plans discovery.
const (
	DirAiMemory       = ".ai-memory"
	DirPlans          = "plans"
	DirPending        = "pending"
	DirCompleted      = "completed"
	DirSubtasks       = "subtasks"
	PathAiMemoryPlans = ".ai-memory/plans"
	PathPending       = "pending/"
	PathCompleted     = "completed/"
	PathSubtasks      = "subtasks/"
	ExtMarkdown       = ".md"
)

// Default execution parameters.
const (
	DefaultLimit      = 10
	DefaultRepoPath   = "."
	DefaultDateLayout = "2006-01-02"
	DefaultJSONIndent = "  "
)

// Table column headers.
const (
	ColID       = "ID"
	ColTitle    = "Title"
	ColStatus   = "Status"
	ColFilePath = "File Path"
)

// UI styling, colors, and badge labels.
const (
	ColorHeader         = "#7aa2f7"
	ColorBorder         = "#3b4261"
	ColorStripe         = "#1f2335"
	TitleCompletedPlans = "Completed Plans"
	TitlePendingPlans   = "Pending / Queued Plans"
	HeaderGitMapPlans   = "GitMap Plans — Repository Plan & Task Inventory"
	TagActive           = "ACTIVE"
	TagInProgress       = "IN PROGRESS"
	TagNotice           = "NOTICE"
	TagOverflow         = "OVERFLOW"
	TagInfo             = "INFO"
	TagCompleted        = "COMPLETED"
	PillCompleted       = "Completed"
	PillPending         = "Pending"
	PillTotal           = "Total"
	DashID              = "-"
)

// UsagePlans contains the CLI manual help text for `gitmap plans`.
const UsagePlans = `Usage: gitmap plans [N] [flags]
Aliases: gitmap plan

Discovers and displays completed and pending architectural plans in .ai-memory/plans/.

Flags:
  -n, --limit <N>     Limit number of plans shown per section (default: 10)
  -a, --all           Show all completed and pending plans without truncation
  -p, --pending       Show only pending / active plans
  -c, --completed     Show only completed plans
  -j, --json          Emit machine-readable JSON to stdout
      --yaml, --yml   Emit machine-readable YAML to stdout
  -h, --help          Print this usage help

Examples:
  gitmap plans
  gitmap plans 20
  gitmap plans --all
  gitmap plans --json
  gitmap plans --yaml
`
