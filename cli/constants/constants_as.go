package constants

// gitmap:cmd top-level
// `gitmap as` — register a short name (alias) for the current Git repo.
//
// Spec: docs in helptext/as.md and README.
const (
	CmdAs        = "as"
	CmdAsAlias   = "s-alias"
	LogPrefixAs  = "[as]"
	FlagAsForce  = "force"
	FlagAsForceS = "f"
)

// User-facing messages for the `as` command.
const (
	MsgAsRegisteredFmt = "  ✓ Registered repo '%s' as alias '%s' (path: %s)\n"
	MsgAsUpdatedFmt    = "  ✓ Updated alias '%s' -> repo '%s' (path: %s)\n"
	MsgAsDBSyncedFmt   = "  ✓ Database now tracks repo '%s' (slug: %s)\n"
	MsgAsHintNext      = "  Tip: run `gitmap release-alias %s <version>` from anywhere to release this repo.\n"
	ErrAsNotInRepoFmt  = "error: not inside a Git repository (cwd: %s). Run this from inside the repo's working tree."
	ErrAsResolveFmt    = "error: could not resolve repo metadata for %s: %v"
	ErrAsAliasInUseFmt = "error: alias '%s' is already mapped to a different repo (slug '%s'). Pass --force to overwrite."
	ErrAsUsage         = "Usage: gitmap as [alias-name] [--force]\n  (omit alias-name to use the repo folder's basename)"
)
