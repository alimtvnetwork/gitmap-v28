package constants

// pull-release mode flags and messages (canonical command name as of
// v5.6.0 — the legacy `release-pull` alias still routes here). Mode
// flags are mutually exclusive; --ff-only is the safe default. The
// --merge mode passes `--no-rebase` to git so it overrides any
// user-level `pull.rebase=true` setting.
const (
	FlagRPFFOnly  = "ff-only"
	FlagRPRebase  = "rebase"
	FlagRPMerge   = "merge"
	FlagRPDryRun  = "dry-run"
	FlagRPVerbose = "verbose"

	FlagDescRPFFOnly  = "Pull mode: fast-forward only (default). Hard-fails on divergent history."
	FlagDescRPRebase  = "Pull mode: rebase local commits onto upstream. Aborts the rebase on conflict."
	FlagDescRPMerge   = "Pull mode: classic merge (passes --no-rebase). Creates a merge commit on divergence."
	FlagDescRPDryRun  = "Show the git pull command that would run, then skip the pull and forward to release."
	FlagDescRPVerbose = "Echo every git invocation to stderr before running it."

	RPModeFFOnly = "ff-only"
	RPModeRebase = "rebase"
	RPModeMerge  = "merge"

	ErrRPNotInRepo       = "pull-release: must be run inside a git repository\n"
	ErrRPCwdFailedFmt    = "pull-release: cannot resolve cwd: %v\n"
	ErrRPModeConflictFmt = "pull-release: --ff-only, --rebase, and --merge are mutually exclusive (got %s)\n"
	ErrRPPullFailedFmt   = "pull-release: git pull %s failed in %s: %v\n"
	ErrRPRebaseAbortFmt  = "pull-release: rebase failed in %s: %v (attempted git rebase --abort)\n"

	MsgRPPullingFmt = "[pull-release] git pull %s in %s\n"
	MsgRPDryRunFmt  = "[pull-release] dry-run: would run `git pull %s` in %s\n"
)
