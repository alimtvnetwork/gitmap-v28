package cmdvscode

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

// RenderGitHubDesktopHelp renders the styled two-column help menu for GitHub Desktop integration.
func RenderGitHubDesktopHelp() {
	termout.RenderMenu(buildGitHubDesktopHelpMenu())
}

func buildGitHubDesktopHelpMenu() termout.HelpMenu {
	return termout.HelpMenu{
		Title: "GitHub Desktop Integration (gitmap gd / github)",
		UsageLines: []string{
			"gitmap gd [folder-path] [flags]",
			"gitmap github [folder-path] [flags]",
			"gitmap github-desktop --all",
			"gitmap gd group <ls|add|rm> [args...]",
		},
		Sections: []termout.HelpSection{
			buildGHDesktopTargetSection(),
			buildGHDesktopGroupSection(),
			buildGHDesktopMaintenanceSection(),
		},
		FooterFlags: buildGHDesktopFooterFlags(),
		Tips: []string{
			"'gd', 'ds', 'github', and 'github-desktop' are interchangeable aliases.",
			"Run 'gitmap gd --install' to automatically install GitHub Desktop if missing.",
			"Run 'gitmap gd' without arguments to register CWD or tracked sub-repositories.",
		},
	}
}

func buildGHDesktopTargetSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Registration Targets",
		Entries: []termout.CommandEntry{
			{Command: "gitmap gd", Description: "Register current repository or all tracked repos under CWD"},
			{Command: "gitmap gd <path>", Description: "Register specific repository folder with GitHub Desktop"},
			{Command: "gitmap gd --all", Description: "Bulk register every tracked repository in the database"},
		},
	}
}

func buildGHDesktopGroupSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Repository Groups",
		Entries: []termout.CommandEntry{
			{Command: "group ls", Description: "List all defined GitHub Desktop repository groups"},
			{Command: "group add <name> <paths...>", Description: "Add one or more repositories to a specified group"},
			{Command: "group rm <name> [path]", Description: "Remove repository from group or delete group entirely"},
		},
	}
}

func buildGHDesktopMaintenanceSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Maintenance & Optimization",
		Entries: []termout.CommandEntry{
			{Command: "optimize-projects, dedupe", Description: "Deduplicate and clean stale GitHub Desktop entries"},
			{Command: "clear, clean", Description: "Remove missing or broken repositories from GitHub Desktop"},
		},
	}
}

func buildGHDesktopFooterFlags() []termout.CommandEntry {
	return []termout.CommandEntry{
		{Command: "-i, --install", Description: "Auto-install GitHub Desktop CLI if missing before registering"},
		{Command: "--all", Description: "Register all repositories in GitMap database"},
		{Command: "-h, --help", Description: "Show this GitHub Desktop help menu"},
	}
}
