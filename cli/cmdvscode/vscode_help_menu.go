package cmdvscode

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/termhelp"
)

// RenderVSCodeHelp displays the styled two-column VS Code help menu.
func RenderVSCodeHelp() {
	termhelp.RenderMenu(buildVSCodeHelpMenu())
}

func buildVSCodeHelpMenu() termhelp.HelpMenu {
	return termhelp.HelpMenu{
		Title: "VS Code Project Manager Suite (gitmap vscode)",
		UsageLines: []string{
			"gitmap vscode [command] [flags]",
			"gitmap code [alias] [path] [extras...]",
			"gitmap vscode add <path>",
		},
		Sections: []termhelp.HelpSection{
			buildVSCodeProjectSection(),
			buildVSCodeDiagnosticsSection(),
		},
		FooterFlags: buildVSCodeFooterFlags(),
		Tips: []string{
			"Use 'gitmap code <alias> <path>' to register and open directly in VS Code.",
			"Run 'gitmap vscode optimize-projects' to purge duplicate paths from projects.json.",
		},
	}
}

func buildVSCodeFooterFlags() []termhelp.CommandEntry {
	return []termhelp.CommandEntry{
		{Command: "-v, --verbose", Description: "Show verbose details during sync and optimization"},
		{Command: "-h, --help", Description: "Show this VS Code help menu"},
	}
}

func buildVSCodeProjectSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Project Management",
		Entries: []termhelp.CommandEntry{
			{Command: "list (ls)", Description: "List all registered VS Code Project Manager entries"},
			{Command: "add <path>", Description: "Register directory path into projects.json"},
			{Command: "rm <target>", Description: "Remove project entry by name or directory path"},
			{Command: "group <sub>", Description: "Group and categorize projects by folder or tag", HasSubcommands: true},
		},
	}
}

func buildVSCodeDiagnosticsSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Diagnostics & Optimization",
		Entries: []termhelp.CommandEntry{
			{Command: "profiles", Description: "Inspect and list VS Code configuration profiles"},
			{Command: "optimize-projects", Description: "Clean up and deduplicate identical workspace paths"},
			{Command: "find-duplicates", Description: "Detect multiple aliases pointing to the same folder"},
			{Command: "clear (clean)", Description: "Purge broken and non-existent project references"},
		},
	}
}
