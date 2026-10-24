package constants

// Command names and aliases for cmdstale.
const (
	CmdStale      = "stale"
	CmdStaleAlias = "sta"
	SubCmdHelp    = "help"
)

// Flag names used for CLI parsing in cmdstale.
const (
	FlagDays      = "days"
	FlagRoot      = "root"
	FlagArchive   = "archive"
	FlagDryRun    = "dry-run"
	FlagFormat    = "format"
	FlagHelp      = "help"
	FlagHelpShort = "h"
)

// Flag CLI argument strings and formatted entries.
const (
	FlagArgDays       = "--days"
	FlagArgRoot       = "--root"
	FlagArgArchive    = "--archive"
	FlagArgDryRun     = "--dry-run"
	FlagArgFormat     = "--format"
	FlagArgHelp       = "--help"
	FlagArgHelpShort  = "-h"
	FlagHelpEntryHelp = "-h, --help"
)

// Flag descriptions for help text and CLI flag definitions.
const (
	FlagDescDays    = "report repos with no commits in the last N days"
	FlagDescRoot    = "scan root directory"
	FlagDescArchive = "move stale repos into .gitmap/archive/"
	FlagDescDryRun  = "preview archive moves without touching disk"
	FlagDescFormat  = "output format: table|json|csv"
	FlagDescHelp    = "show stale help text"
)

// Default thresholds, boundaries, and sizing values.
const (
	DefaultStaleDays      = 90
	DefaultScanRoot       = "."
	DefaultOutputFormat   = "table"
	MinHygieneWorkers     = 2
	MaxHygieneWorkers     = 8
	HoursPerDay           = 24
	DefaultArchiveDirPerm = 0o755
)

// Directory names and paths.
const (
	DirDotGit  = ".git"
	DirGitmap  = ".gitmap"
	DirArchive = "archive"
)

// Date, time, and serialization format layouts.
const (
	DateFormatStale        = "2006-01-02"
	TimestampFormatArchive = "20060102T150405Z"
	JSONIndent             = "  "
	JSONPrefix             = ""
)

// CSV column header names.
const (
	HeaderPath          = "path"
	HeaderLastCommitUTC = "last_commit_utc"
	HeaderAgeDays       = "age_days"
)

// Git CLI executable name, subcommands, and command-line arguments.
const (
	GitBin               = "git"
	GitFlagDir           = "-C"
	GitSubCmdLog         = "log"
	GitArgLimitOne       = "-1"
	GitFormatCommitEpoch = "--format=%ct"
)

// ANSI terminal color escape codes.
const (
	ColorCyan   = "\033[36m"
	ColorYellow = "\033[33m"
	ColorGreen  = "\033[32m"
	ColorRed    = "\033[31m"
	ColorReset  = "\033[0m"
)

// User-facing output templates and format strings.
const (
	MsgNoStaleReposFmt = "\n  no repos stale beyond %d days\n\n"
	MsgStaleHeaderFmt  = "\n  \033[36m%d stale repo(s)\033[0m (no commits in %d days)\n\n"
	MsgStaleRowFmt     = "  \033[33m%4dd\033[0m  %s  (last: %s)\n"
	MsgWouldArchiveFmt = "  \033[33mwould archive\033[0m %s -> %s\n"
	MsgArchivedFmt     = "  \033[32marchived\033[0m %s -> %s\n"
	MsgFailedFmt       = "  \033[31mfailed\033[0m %s: %v\n"
)

// Error format strings and messages.
const (
	ErrJSONEncodeFmt    = "  json encode: %v\n"
	ErrArchiveMkdirFmt  = "  archive mkdir: %v\n"
	ErrInvalidFormatFmt = "invalid --format %q (want table|json|csv)"
)

// Common format specifiers and string sentinels.
const (
	AgeDaysIntFmt = "%d"
	SecScanFmt    = "%d"
	EmptyString   = ""
	Newline       = "\n"
)
