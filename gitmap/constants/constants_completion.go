package constants

// Completion shells.

// Completion shells.
const (
	ShellPowerShell = "powershell"
	ShellBash       = "bash"
	ShellZsh        = "zsh"
)

// Completion list flags.
const (
	CompListRepos      = "--list-repos"
	CompListGroups     = "--list-groups"
	CompListCommands   = "--list-commands"
	CompListAliases    = "--list-aliases"
	CompListZipGroups  = "--list-zip-groups"
	CompListHelpGroups = "--list-help-groups"
)

// Completion file names.
const (
	CompFilePS   = "completions.ps1"
	CompFileBash = "completions.bash"
	CompFileZsh  = "completions.zsh"
	CompDirName  = "gitmap"
)

// Completion help text.
const HelpCompletionLong = "  completion (cmp)    Generate or install shell tab-completion scripts"

// Completion messages.
const (
	MsgCompInstalled    = "Shell completion installed for %s\n"
	MsgCompAlreadyDone  = "Shell completion already configured for %s\n"
	MsgCompProfileWrite = "Added source line to %s\n"
	ErrCompUsage        = "usage: gitmap completion <powershell|bash|zsh> [--list-repos|--list-groups|--list-commands|--list-aliases|--list-zip-groups|--list-help-groups]\n"
	ErrCompUnknownShell = "unknown shell: %s (use powershell, bash, or zsh)\n"
	ErrCompProfileWrite = "failed to update profile at %s: %v (operation: write)\n"
)

// Completion flag descriptions.
const (
	FlagDescCompListRepos      = "Print repo slugs one per line"
	FlagDescCompListGroups     = "Print group names one per line"
	FlagDescCompListCommands   = "Print all command names one per line"
	FlagDescCompListAliases    = "Print alias names one per line"
	FlagDescCompListZipGroups  = "Print zip group names one per line"
	FlagDescCompListHelpGroups = "Print help group names one per line"
)
