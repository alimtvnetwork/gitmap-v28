package cmdclone

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/termhelp"
)

// RenderCloneHelp displays the styled two-column Clone help menu.
func RenderCloneHelp() {
	termhelp.RenderMenu(buildCloneHelpMenu())
}

func buildCloneHelpMenu() termhelp.HelpMenu {
	return termhelp.HelpMenu{
		Title: "Repository Cloner (gitmap clone)",
		UsageLines: []string{
			"gitmap clone <source|json|csv|text|url> [target-dir] [flags]",
			"gitmap c <source|json|csv|text|url> [target-dir] [flags]",
			"gitmap clone next [flags]",
		},
		Sections: []termhelp.HelpSection{
			buildCloneSourceSection(),
			buildCloneSubcmdSection(),
			buildCloneConcurrencySection(),
			buildCloneEnvSection(),
		},
		FooterFlags: buildCloneFooterFlags(),
		Tips: []string{
			"Run 'gitmap clone repos.json --audit' to preview clone paths safely.",
			"Use 'gitmap clone <url> --ssh' to clone using SSH credentials.",
			"Run 'gitmap clone next' to incrementally process a repo manifest.",
		},
	}
}

func buildCloneFooterFlags() []termhelp.CommandEntry {
	return []termhelp.CommandEntry{
		{Command: "-f, --force", Description: "Wipe and re-clone destination directories fresh"},
		{Command: "--dry-run", Description: "Preview clone actions without modifying disk"},
		{Command: "--safe-pull", Description: "Pull existing repos with retry and diagnostics"},
		{Command: "--audit", Description: "Dry-run validation with diff-style command summary"},
		{Command: "--ssh / --https", Description: "Force URL protocol conversion before cloning"},
		{Command: "--no-vscode-sync", Description: "Skip syncing cloned repos to VS Code Project Manager"},
		{Command: "-v, --verbose", Description: "Emit detailed clone debug telemetry"},
		{Command: "-h, --help", Description: "Show this clone help menu"},
	}
}

func buildCloneSourceSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Input Sources & Formats",
		Entries: []termhelp.CommandEntry{
			{Command: "<git-url> [dir]", Description: "Clone single repository directly from Git URL"},
			{Command: "<file.json>", Description: "Batch clone from GitMap scan JSON manifest"},
			{Command: "<file.csv>", Description: "Batch clone from CSV repository inventory"},
			{Command: "<file.txt>", Description: "Batch clone from plain text URL list (one per line)"},
			{Command: "ls / list", Description: "Inspect manifest items without cloning"},
		},
	}
}
