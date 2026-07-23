package constants

// gitmap:cmd top-level
// `gitmap release-alias` and `gitmap release-alias-pull` — release a repo by
// its registered alias from anywhere on disk.
//
// Spec: docs in helptext/release-alias.md and README.
const (
	CmdReleaseAlias          = "release-alias"
	CmdReleaseAliasShort     = "ra"
	CmdReleaseAliasPull      = "release-alias-pull"
	CmdReleaseAliasPullShort = "rap"
	LogPrefixReleaseAlias    = "[release-alias]"
	FlagRAPull               = "pull"
	FlagRANoStash            = "no-stash"
	FlagRADryRun             = "dry-run"
)

// Auto-stash bookkeeping.
const (
	RAStashMessageFmt    = "gitmap-release-alias autostash %s"
	MsgRAStashCreatedFmt = "  ✓ Auto-stashed dirty changes (label: %s)\n"
	MsgRAStashPoppedFmt  = "  ✓ Restored auto-stashed changes (label: %s)\n"
	MsgRAPullingFmt      = "  ↻ Pulling latest from origin in %s\n"
	MsgRAReleasingFmt    = "  ▸ Releasing repo '%s' (path: %s) version=%s\n"
	WarnRAStashPopFailed = "  ⚠ Auto-stash pop failed; your changes remain in `git stash`. Resolve manually."
	ErrRAUsage           = "Usage: gitmap release-alias <alias> <version> [--pull] [--no-stash] [--dry-run]"
	ErrRAUnknownAliasFmt = "error: alias '%s' is not registered. Run `gitmap as %s` from the repo first."
	ErrRAChdirFailedFmt  = "error: could not change directory to '%s': %v"
	ErrRAPullFailedFmt   = "error: git pull failed in %s: %v"
	ErrRAStashFailedFmt  = "error: auto-stash failed in %s: %v"
)
