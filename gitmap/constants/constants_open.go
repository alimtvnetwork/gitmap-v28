// Package constants — open command IDs and messages.
//
// `gitmap open` launches GitHub Desktop AND VS Code on the current
// repo (cwd or git toplevel). Re-uses the inject pipeline's Desktop
// + VS Code helpers so newly-moved or freshly-cloned repos always
// land in both tools without a separate `inject` call.
package constants

// gitmap:cmd top-level
// Open command help-line entry surfaced by `gitmap help`.
const HelpCmdOpen = "  open (op)                              Open current repo in GitHub Desktop AND VS Code"

// Open command messages.
const (
	MsgOpenStart       = "Opening %q (%s) in GitHub Desktop and VS Code...\n"
	MsgOpenDone        = "  ✓ open: %q ready in both tools\n"
	WarnOpenNotGitRepo = "open: WARN current folder is not inside a git repo with origin remote — proceeding with Desktop + VS Code only\n"
	ErrOpenResolveCwd  = "open: ERROR cannot determine current directory: %v\n"
)
