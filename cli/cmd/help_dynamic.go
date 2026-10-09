package cmd

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/helpdoc"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

// RenderDynamicCommandHelp generates and renders a modern framed box help menu
// for any GitMap command, deriving summaries from catalog metadata or generating
// realistic usage patterns and examples.
func RenderDynamicCommandHelp(command string) bool {
	cmdClean := strings.TrimSpace(command)
	if cmdClean == "" {
		return false
	}

	summary := helpdoc.GetTopicDetailedSummary(cmdClean)
	if summary == "" || summary == "Gitmap command line utilities." {
		summary = fmt.Sprintf("Execute, manage, and automate GitMap %s operations.", cmdClean)
	}

	menu := termout.HelpMenu{
		Title: fmt.Sprintf("gitmap %s", cmdClean),
		UsageLines: []string{
			fmt.Sprintf("gitmap %s [arguments] [flags]", cmdClean),
			fmt.Sprintf("gitmap %s help", cmdClean),
		},
		Sections: []termout.HelpSection{
			{
				Title: "Description",
				Entries: []termout.CommandEntry{
					{Command: cmdClean, Description: summary},
				},
			},
			{
				Title: "Commands & Execution",
				Entries: []termout.CommandEntry{
					{Command: cmdClean, Description: "Run command with default parameters"},
					{Command: fmt.Sprintf("%s --dry-run", cmdClean), Description: "Simulate execution without state modifications"},
					{Command: fmt.Sprintf("%s --json", cmdClean), Description: "Output machine-readable JSON telemetry"},
				},
			},
			{
				Title: "Examples",
				Entries: []termout.CommandEntry{
					{Command: fmt.Sprintf("gitmap %s", cmdClean), Description: "Standard invocation"},
					{Command: fmt.Sprintf("gitmap %s --help", cmdClean), Description: "Inspect full options and syntax"},
				},
			},
		},
		FooterFlags: []termout.CommandEntry{
			{Command: "-h, --help", Description: "Display this formatted help screen"},
			{Command: "--json", Description: "Emit results formatted as JSON"},
			{Command: "--dry-run", Description: "Simulate operations without making changes"},
		},
		Tips: []string{
			fmt.Sprintf("Run 'gitmap %s help' anytime for quick reference.", cmdClean),
			"Combine with '--json' for pipeline and AI automation fspath.",
		},
	}

	termout.RenderMenu(menu)
	printUsageFooterShort()

	return true
}
