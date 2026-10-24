package constants

// Command names and aliases for cmdbranch.
const (
	CmdBranch          = "branch"
	CmdBranchAlias     = "b"
	CmdSubDefault      = "default"
	CmdSubDefaultAlias = "def"
	CmdSubLatest       = "latest-branch"
	CmdSubLatestAlias  = "lb"
	CmdSubList         = "list"
	CmdSubCurrent      = "current"
	CmdSubHelp         = "help"
)

// Flag names used for CLI parsing.
const (
	FlagJSON   = "json"
	FlagHelp   = "help"
	FlagSort   = "sort"
	FlagFilter = "filter"
	FlagFormat = "format"
)

// Flag CLI argument strings and formatted menu entries.
const (
	FlagArgJSON       = "--json"
	FlagArgJSONShort  = "-j"
	FlagArgHelp       = "--help"
	FlagArgHelpShort  = "-h"
	FlagHelpEntryJSON = "-j, --json"
	FlagHelpEntryHelp = "-h, --help"
)

// Flag descriptions for help text and menu definitions.
const (
	FlagDescJSON = "Output branch metadata in structured JSON format"
	FlagDescHelp = "Show this branch help menu"
)

// Help menu layout strings, section titles, and tips.
const (
	HelpTitleBranch         = "Branch Resolution & Management (gitmap branch)"
	HelpUsageDefault        = "gitmap branch default [flags]"
	HelpUsageDefAlias       = "gitmap b def"
	HelpUsageLatestBranch   = "gitmap latest-branch"
	HelpSectionOperations   = "Branch Operations"
	HelpEntryDefault        = "default (def)"
	HelpDescDefault         = "Resolve repository's default branch using symbolic refs"
	HelpEntryLatest         = "latest-branch (lb)"
	HelpDescLatest          = "Inspect the newest branch with recent commit activity"
	HelpTipResolveDefault   = "Quickly resolve default branch (main/master) for any cloned repo."
	HelpTipOfflineFallback  = "Works offline and falls back gracefully to 'main' when no remotes are configured."
)

// Error codes, error messages, and logging format strings.
const (
	ErrCodeE9000                  = "E9000"
	ErrFatalMessage               = "fatal error"
	ErrBranchUnknownSubcommand    = "  ✗ Unknown branch subcommand: %q (expected: default | def)\n"
	ErrBranchUnknownSubcommandFmt = "  ✗ Unknown branch subcommand: %q (expected: default | def)\n"
	ErrBranchMissingSubcommand    = "  ✗ Missing subcommand for `gitmap branch` (expected: default | def)\n"
	ErrBranchNotRepo              = "  ✗ Not inside a git work tree\n"
	ErrBranchDefaultFailed        = "  ✗ Failed to checkout default branch %s: %v\n"
	ErrBranchDefaultFailedFmt     = "  ✗ Failed to checkout default branch %s: %v\n"
	MsgBranchDefaultSwitching     = "  ▶ Switching to default branch %s...\n"
	MsgBranchDefaultSwitchingFmt  = "  ▶ Switching to default branch %s...\n"
)

// Sentinels and default values.
const (
	SentinelCurrentDir            = "."
	SentinelDefaultBranchFallback = "main"
)
