package helpdisplay

// CommandHelpGroup bundles related command helpers under one header with hints.
type CommandHelpGroup struct {
	header   string
	commands []CommandHelper
	hints    []string
}

// NewCommandHelpGroup builds a CommandHelpGroup.
func NewCommandHelpGroup(header string, commands []CommandHelper, hints []string) CommandHelpGroup {
	return CommandHelpGroup{header: header, commands: commands, hints: hints}
}
