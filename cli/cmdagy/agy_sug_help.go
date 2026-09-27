package cmdagy

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/termhelp"
)

// RenderAgySugHelp renders the rich two-column boxed terminal UI help menu for shutdown-until-green.
func RenderAgySugHelp() {
	termhelp.RenderMenu(buildAgySugHelpMenu())
}

func buildAgySugHelpMenu() termhelp.HelpMenu {
	return termhelp.HelpMenu{
		Title: "Shutdown Until Green (gitmap shutdown-until / sug)",
		UsageLines: []string{
			"gitmap shutdown-until <command> [flags]",
			"gitmap sug <command> [flags]",
			"gitmap agy shutdown-until <command> [flags]",
			"gitmap agy sug <command> [flags]",
		},
		Sections: []termhelp.HelpSection{
			buildSugCommandsSection(),
			buildSugWorkflowSection(),
		},
		FooterFlags: buildSugFooterFlags(),
		Tips: []string{
			"Run 'gitmap sug agy-running-projects' to auto-register all active workspaces.",
			"Run 'gitmap sug run --dry-run' to test the watch loop safely without OS shutdown.",
			"Use 'gitmap sug ls' to inspect currently registered watch targets.",
		},
	}
}

func buildSugCommandsSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Watch List Commands",
		Entries: []termhelp.CommandEntry{
			{Command: "ls (list)", Description: "List registered project targets in shutdown watch list"},
			{Command: "add-projects <targets...>", Description: "Register one or more projects into watch list"},
			{Command: "rm <targets...>", Description: "Remove designated projects from watch list"},
			{Command: "agy-running-projects", Description: "Auto-populate watch list with all currently running AGY projects"},
			{Command: "help", Description: "Show this two-column interactive help menu"},
		},
	}
}

func buildSugWorkflowSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Execution & Monitoring",
		Entries: []termhelp.CommandEntry{
			{Command: "run [-t <duration>]", Description: "Start continuous watch loop until all monitored pipelines turn green"},
			{Command: "run --dry-run (-n)", Description: "Dry-run watch loop; simulates final OS shutdown when all pipelines pass"},
			{Command: "run --once (-1)", Description: "Evaluate pipeline status once and exit without looping"},
		},
	}
}

func buildSugFooterFlags() []termhelp.CommandEntry {
	return []termhelp.CommandEntry{
		{Command: "-t, --time <duration>", Description: "Polling check interval (default: 5m, minimum 2m; 10s in dry-run)"},
		{Command: "-n, --dry-run", Description: "Simulate OS shutdown execution without powering off system"},
		{Command: "-1, --once", Description: "Execute a single status evaluation cycle and exit immediately"},
		{Command: "-h, --help", Description: "Display this comprehensive help menu"},
	}
}
