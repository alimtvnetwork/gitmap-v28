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
			"gitmap sug <target> [flags]",
			"gitmap agy shutdown-until <command> [flags]",
			"gitmap agy sug <command> [flags]",
		},
		Sections: []termhelp.HelpSection{
			buildSugCommandsSection(),
			buildSugWorkflowSection(),
			buildSugTargetDefinitionSection(),
		},
		FooterFlags: buildSugFooterFlags(),
		Tips: []string{
			"Run 'gitmap sug agy-running-projects' (or 'gitmap sug arp') to auto-register all active workspaces.",
			"Run 'gitmap sug status' to check if a watch loop process is currently active or idle.",
			"Run 'gitmap sug ui' (or 'gitmap sug web') to open the dark-mode dashboard in your browser.",
			"Direct path: 'gitmap sug d:/work/my-project' registers and monitors the target immediately.",
		},
	}
}

func buildSugCommandsSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Watch List Commands",
		Entries: []termhelp.CommandEntry{
			{Command: "ls / status (st)", Description: "List registered targets and check live watch loop process state"},
			{Command: "add-projects <targets...>", Description: "Register project targets into watch list (alias: add, /add, ap)"},
			{Command: "rm <targets...>", Description: "Remove targets from watch list (alias: /rm, remove, del, delete)"},
			{Command: "agy-running-projects", Description: "Auto-populate with running AGY projects (alias: running-projects, arp, rp)"},
			{Command: "ui (web)", Description: "Launch dark-mode local browser dashboard for visual project monitoring"},
			{Command: "help", Description: "Show this comprehensive interactive help menu"},
		},
	}
}

func buildSugWorkflowSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Execution & Monitoring",
		Entries: []termhelp.CommandEntry{
			{Command: "run / watch (w) [-t <dur>]", Description: "Start continuous watch loop until all monitored pipelines turn green"},
			{Command: "run --dry-run (-n)", Description: "Dry-run watch loop; simulates final OS shutdown when all pipelines pass"},
			{Command: "run --once (-1)", Description: "Evaluate pipeline status once and exit without looping"},
			{Command: "<target> [-t <dur>]", Description: "Directly register and watch target (e.g. 'gitmap sug d:/work/my-repo')"},
		},
	}
}

func buildSugTargetDefinitionSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Target Types & Examples",
		Entries: []termhelp.CommandEntry{
			{Command: "Local Path", Description: "Directory on disk (e.g. 'd:/work/antigravity-manager', '.', '../repo')"},
			{Command: "Repo Alias / Name", Description: "Indexed repository in GitMap DB (e.g. 'gitmap-v28', 'antigravity-manager')"},
			{Command: "Git Remote URL", Description: "GitHub/Git repository URL (e.g. 'https://github.com/alimtvnetwork/gitmap-v28')"},
		},
	}
}

func buildSugFooterFlags() []termhelp.CommandEntry {
	return []termhelp.CommandEntry{
		{Command: "-t, --time <duration>", Description: "Polling check interval (default: 5m, minimum 2m; 10s in dry-run)"},
		{Command: "-n, --dry-run", Description: "Simulate OS shutdown execution without powering off system"},
		{Command: "-1, --once", Description: "Execute a single status evaluation cycle and exit immediately"},
		{Command: "--no-open", Description: "Do not automatically launch web browser when running 'ui'"},
		{Command: "-h, --help", Description: "Display this comprehensive help menu"},
	}
}
