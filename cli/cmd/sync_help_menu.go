package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/termhelp"
)

// RenderSyncHelp displays the styled two-column Sync help menu.
func RenderSyncHelp() {
	termhelp.RenderMenu(buildSyncHelpMenu())
}

func buildSyncHelpMenu() termhelp.HelpMenu {
	return termhelp.HelpMenu{
		Title: "Repository Baseline Sync (gitmap sync)",
		UsageLines: []string{
			"gitmap sync <target> [flags]",
			"gitmap sy <target> [flags]",
		},
		Sections: []termhelp.HelpSection{
			buildSyncTargetsSection(),
		},
		FooterFlags: buildSyncFooterFlags(),
		Tips: []string{
			"Line-based targets append missing lines only; existing entries are preserved.",
			"Run 'gitmap sync all --dry-run' to preview planned updates across configs.",
		},
	}
}

func buildSyncFooterFlags() []termhelp.CommandEntry {
	return []termhelp.CommandEntry{
		{Command: "-n, --dry-run", Description: "Print planned additions without touching disk"},
		{Command: "-f, --force", Description: "Overwrite conflicting JSON keys in .prettierrc"},
		{Command: "-h, --help", Description: "Show this sync help menu"},
	}
}

func buildSyncTargetsSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Configuration Sync Targets",
		Entries: []termhelp.CommandEntry{
			{Command: "all", Description: "Run every target below in sequence (full baseline)"},
			{Command: "ignore", Description: "Union-merge curated defaults into .gitignore"},
			{Command: "attributes", Description: "Union-merge curated defaults into .gitattributes"},
			{Command: "lfs-install", Description: "Run git-lfs setup and configure large file filters"},
			{Command: "prettier-ignore", Description: "Union-merge curated defaults into .prettierignore"},
			{Command: "prettier-rc", Description: "Merge curated JSON formatting defaults into .prettierrc"},
		},
	}
}
