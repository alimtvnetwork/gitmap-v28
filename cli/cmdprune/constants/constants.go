package constants

// Command names and aliases for cmdprune.
const (
	CmdPrune      = "prune"
	CmdPruneAlias = "prn"
)

// Flag names used for CLI parsing in cmdprune.
const (
	FlagDryRun  = "dry-run"
	FlagConfirm = "confirm"
	FlagRemote  = "remote"
	FlagHelp    = "help"
)

// Flag CLI argument strings and formatted menu entries.
const (
	FlagArgDryRun     = "--dry-run"
	FlagArgConfirm    = "--confirm"
	FlagArgRemote     = "--remote"
	FlagArgHelp       = "--help"
	FlagArgHelpShort  = "-h"
	FlagHelpEntryHelp = "-h, --help"
)

// Flag descriptions for help text and flag parsing.
const (
	FlagDescDryRun  = "Preview changes without applying them"
	FlagDescConfirm = "Confirm destructive operation"
	FlagDescRemote  = "Also delete remote release branches"
	FlagDescHelp    = "Show prune help text"
)

// Git CLI tokens and command names.
const (
	GitBin              = "git"
	GitBranch           = "branch"
	GitBranchListFlag   = "--list"
	GitBranchDeleteFlag = "-D"
	GitPush             = "push"
	GitOrigin           = "origin"
	GitPushDeleteFlag   = "--delete"
	ReleaseBranchPrefix = "release/"
	GitWildcard         = "*"
)

// Interactive prompt responses.
const (
	PromptAnswerY      = "y"
	PromptAnswerYUpper = "Y"
	PromptAnswerN      = "n"
	PromptAnswerNUpper = "N"
)

// Output messages and status banners.
const (
	MsgPruneStaleHeader  = "\n  Stale release branches (%d):\n"
	MsgPruneStaleItem    = "    %s  →  tag %s exists\n"
	MsgPruneDryRunHint   = "\n  Use --confirm to delete, or run without --dry-run for interactive mode.\n"
	MsgPruneDeleting     = "\n  Pruning stale release branches...\n"
	MsgPruneDeleted      = "    ✓ Deleted %s\n"
	MsgPruneRemoteDelete = "    ✓ Deleted remote %s\n"
	MsgPruneRemoteWarn   = "    ⚠ Failed to delete remote %s: %v\n"
	MsgPruneSummary      = "\n  Summary: %d deleted, %d kept.\n"
	MsgPruneNone         = "  No stale release branches found.\n"
	MsgPrunePrompt       = "  Delete %d stale branch(es)? (y/N): "
	MsgPruneAborted      = "  Prune aborted.\n"
)

// Error codes, error messages, and format strings.
const (
	ErrCodePruneFailed   = "EPRUNE01"
	ErrPruneListBranches = "failed to list branches: %v\n"
	ErrPruneDeleteBranch = "    ✗ Failed to delete %s: %v\n"
)

// Help documentation strings and usage specifications.
const (
	HelpPrune = "  prune (prn)         Delete stale release branches that have been tagged"
	HelpUsage = "gitmap prune [--dry-run] [--confirm] [--remote]"
)

// Common sentinels and delimiters.
const (
	Newline     = "\n"
	EmptyString = ""
)
