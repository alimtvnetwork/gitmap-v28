package constants

// gitmap:cmd top-level
// Move/merge command IDs and aliases.
//
// Spec: spec/01-app/97-move-and-merge.md
const (
	CmdMv         = "mv"
	CmdMove       = "move"
	CmdMergeBoth  = "merge-both"
	CmdMergeBothA = "mb"
	CmdMergeLeft  = "merge-left"
	CmdMergeLeftA = "ml"
	CmdMergeRight = "merge-right"
	CmdMergeRgtA  = "mr"
)

// Move/merge flag names.
const (
	FlagMMYes        = "yes"
	FlagMMYesShort   = "y"
	FlagMMAccept     = "accept-all"
	FlagMMAcceptShrt = "a"
	FlagMMPreferL    = "prefer-left"
	FlagMMPreferR    = "prefer-right"
	FlagMMPreferNew  = "prefer-newer"
	FlagMMPreferSkip = "prefer-skip"
	FlagMMNoPush     = "no-push"
	FlagMMNoCommit   = "no-commit"
	FlagMMForceFold  = "force-folder"
	FlagMMPullFold   = "pull"
	FlagMMInit       = "init"
	FlagMMDryRun     = "dry-run"
	FlagMMIncludeVCS = "include-vcs"
	FlagMMIncludeNM  = "include-node-modules"
)

// Log prefixes per command.
const (
	LogPrefixMv         = "[mv]"
	LogPrefixMergeBoth  = "[merge-both]"
	LogPrefixMergeLeft  = "[merge-left]"
	LogPrefixMergeRight = "[merge-right]"
)

// Commit message templates.
const (
	CommitMsgMv         = "gitmap mv from %s"
	CommitMsgMergeBoth  = "gitmap merge-both with %s"
	CommitMsgMergeLeft  = "gitmap merge-left from %s"
	CommitMsgMergeRight = "gitmap merge-right from %s"
)

// Git argument tokens used by movemerge finalize.
const (
	GitAddCmd     = "add"
	GitAddAllArg  = "-A"
	GitCommitCmd  = "commit"
	GitMessageArg = "-m"
)

// Conflict prompt and error messages.
const (
	ConflictPromptLine = "  [L]eft  [R]ight  [S]kip  [A]ll-left  [B]all-right  [Q]uit"
	ConflictPromptCue  = "  > "
	ErrMMUsageFmt      = "Usage: gitmap %s LEFT RIGHT [flags]\n"
	ErrMMSameFolderFmt = "error: LEFT and RIGHT resolve to the same folder: %s"
	ErrMMNestedFmt     = "error: RIGHT is nested inside LEFT (or vice versa): LEFT=%s RIGHT=%s"
	ErrMMOriginFmt     = "error: folder '%s' exists but its remote is '%s', not '%s'. Pass --force-folder to overwrite, or rename it."
	ErrMMSrcMissingFmt = "error: source '%s' does not exist"
	ErrMMQuit          = "user pressed Q (quit)"
	ErrMMPushFailFmt   = "Push failed. Local commit is preserved at %s. Resolve manually or re-run with --no-push to skip."
)
