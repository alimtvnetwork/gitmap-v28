package constants

// gitmap:cmd top-level
// Multi-group CLI commands.
const (
	CmdMultiGroup      = "multi-group"
	CmdMultiGroupAlias = "mg"
)

// gitmap:cmd top-level
// Multi-group subcommands.
const (
	CmdMGClear  = "clear"  // gitmap:cmd skip
	CmdMGPull   = "pull"   // gitmap:cmd skip
	CmdMGStatus = "status" // gitmap:cmd skip
	CmdMGExec   = "exec"   // gitmap:cmd skip
)

// Multi-group help text is declared in constants_cli.go (HelpMultiGroup).

// Multi-group messages.
const (
	MsgMGActive       = "Active multi-group: %s\n"
	MsgMGSet          = "Multi-group set: %s\n"
	MsgMGCleared      = "Multi-group selection cleared.\n"
	MsgMGNone         = "No multi-group set. Use 'gitmap mg g1,g2' to select groups.\n"
	ErrMGUsage        = "Usage: gitmap multi-group <group1,group2,...|clear|pull|status|exec>\n"
	ErrMGGroupMissing = "Group not found: %s\n"
)
