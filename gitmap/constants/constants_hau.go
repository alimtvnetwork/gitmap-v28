package constants

// gitmap:cmd top-level
// Has-any-updates command.
const (
	CmdHasAnyUpdates      = "has-any-updates"
	CmdHasAnyUpdatesAlias = "hau"
	CmdHasAnyChanges      = "has-any-changes"
	CmdHasAnyChangesAlias = "hac"
)

// Has-any-updates help text.
const HelpHasAnyUpdates = "  has-any-updates (hau/hac) Check if remote has new commits"

// Has-any-updates messages.
const (
	MsgHAUChecking    = "  Checking for updates...\n"
	MsgHAUYes         = "\n  ✓ Yes, you have %d new update(s) from remote.\n    Run 'git pull' to sync.\n"
	MsgHAUNo          = "\n  ✓ You are up to date. No new changes.\n"
	MsgHAUAhead       = "\n  ✓ You are %d commit(s) ahead of remote. No incoming changes.\n"
	MsgHAUDiverged    = "\n  ⚠ Branch has diverged: %d ahead, %d behind remote.\n    Run 'git pull --rebase' or 'git pull' to reconcile.\n"
	MsgHAUNoUpstream  = "\n  ⚠ No upstream tracking branch configured.\n    Run 'git branch --set-upstream-to=origin/<branch>' first.\n"
	ErrHAUNotRepo     = "Error: not inside a Git repository.\n"
	ErrHAUFetchFailed = "  Warning: fetch failed, using cached refs: %v\n"
)
