package helpdisplay

// CommandHelper describes one command's help entry, including nested sub-helpers.
type CommandHelper struct {
	command     string
	description string
	example     string
	url         string
	subHelpers  []CommandHelper
}

// NewCommandHelper builds a CommandHelper; subHelpers nests sub-command help.
func NewCommandHelper(command, description, example, url string, subHelpers ...CommandHelper) CommandHelper {
	return CommandHelper{command: command, description: description, example: example, url: url, subHelpers: subHelpers}
}
