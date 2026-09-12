// Package constants — open command IDs and messages.
//
// `gitmap open` launches the OS-native file explorer (or default
// handler) for the target path.
package constants

// gitmap:cmd top-level
// Open command help-line entry surfaced by `gitmap help`.
const HelpCmdOpen = "  open (o)            Open target (or current repo) in default file explorer"

// Open command messages.
const (
	ErrOpenResolveCwd = "open: ERROR cannot determine current directory: %v\n"
)
