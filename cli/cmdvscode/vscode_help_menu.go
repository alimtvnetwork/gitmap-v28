package cmdvscode

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

// RenderVSCodeHelp displays the styled two-column VS Code help menu.
func RenderVSCodeHelp() {
	termout.RenderMenu(buildVSCodeHelpMenu())
}

func buildVSCodeHelpMenu() termout.HelpMenu {
	return termout.HelpMenu{
		Title: "VS Code Project Manager Suite (gitmap vscode)",
		UsageLines: []string{
			"gitmap vscode [command] [flags]",
			"gitmap code [alias] [path] [extras...]",
			"gitmap vscode add <path>",
		},
		Sections: []termout.HelpSection{
			buildVSCodeProjectSection(),
			buildVSCodeDiagnosticsSection(),
		},
		FooterFlags: buildVSCodeFooterFlags(),
		Tips: []string{
			"Use 'gitmap code <alias> <path>' to register and open directly in VS Code.",
			"Run 'gitmap vscode optimize-projects' to purge duplicate paths from projects.json.",
			"Run 'gitmap vscode repair' to fix startup crashes (ICU error) and sync projects.json.",
		},
	}
}

func buildVSCodeFooterFlags() []termout.CommandEntry {
	return []termout.CommandEntry{
		{Command: "-v, --verbose", Description: "Show verbose details during sync and optimization"},
		{Command: "-h, --help", Description: "Show this VS Code help menu"},
	}
}

func buildVSCodeProjectSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Project Management",
		Entries: []termout.CommandEntry{
			{Command: "list (ls)", Description: "List all registered VS Code Project Manager entries"},
			{Command: "add <path>", Description: "Register directory path into projects.json"},
			{Command: "rm <target>", Description: "Remove project entry by name or directory path"},
			{Command: "group <sub>", Description: "Group and categorize projects by folder or tag", HasSubcommands: true},
		},
	}
}

func buildVSCodeDiagnosticsSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Diagnostics & Optimization",
		Entries: []termout.CommandEntry{
			{Command: "profiles", Description: "Inspect and list VS Code configuration profiles"},
			{Command: "optimize-projects", Description: "Clean up and deduplicate identical workspace paths"},
			{Command: "find-duplicates", Description: "Detect multiple aliases pointing to the same folder"},
			{Command: "clear (clean)", Description: "Purge broken and non-existent project references"},
			{Command: "repair (fix)", Description: "Diagnose & repair VS Code executable and projects.json"},
		},
	}
}
