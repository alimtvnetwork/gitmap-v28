package constants

// gitmap:cmd top-level
// Prune command names.
const (
	CmdPrune      = "prune"
	CmdPruneAlias = "prn" // short alias (was "pr" before v5.6.0; "pr" now belongs to pull-release)
)

// Prune flag names.
const (
	PruneFlagDryRun  = "dry-run"
	PruneFlagConfirm = "confirm"
	PruneFlagRemote  = "remote"
)

// Prune messages.
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

// Prune errors.
const (
	ErrPruneListBranches = "failed to list branches: %v\n"
	ErrPruneDeleteBranch = "    ✗ Failed to delete %s: %v\n"
)

// Prune git arguments.
const (
	GitBranchDeleteFlag = "-D"
	GitPushDeleteFlag   = "--delete"
)

// Prune help text.
const HelpPrune = "  prune (prn)         Delete stale release branches that have been tagged"
