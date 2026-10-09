package cmdsync

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

// RenderSyncHelp displays the styled two-column Sync help menu.
func RenderSyncHelp() {
	termout.RenderMenu(buildSyncHelpMenu())
}

func buildSyncHelpMenu() termout.HelpMenu {
	return termout.HelpMenu{
		Title: "Repository Baseline Sync (gitmap sync)",
		UsageLines: []string{
			"gitmap sync <target> [flags]",
			"gitmap sy <target> [flags]",
		},
		Sections: []termout.HelpSection{
			buildSyncTargetsSection(),
		},
		FooterFlags: buildSyncFooterFlags(),
		Tips: []string{
			"Line-based targets append missing lines only; existing entries are preserved.",
			"Run 'gitmap sync all --dry-run' to preview planned updates across configs.",
		},
	}
}

func buildSyncFooterFlags() []termout.CommandEntry {
	return []termout.CommandEntry{
		{Command: "-n, --dry-run", Description: "Print planned additions without touching disk"},
		{Command: "-f, --force", Description: "Overwrite conflicting JSON keys in .prettierrc"},
		{Command: "-h, --help", Description: "Show this sync help menu"},
	}
}

func buildSyncTargetsSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Configuration Sync Targets",
		Entries: []termout.CommandEntry{
			{Command: "all", Description: "Run every target below in sequence (full baseline)"},
			{Command: "ignore", Description: "Union-merge curated defaults into .gitignore"},
			{Command: "attributes", Description: "Union-merge curated defaults into .gitattributes"},
			{Command: "lfs-install", Description: "Run git-lfs setup and configure large file filters"},
			{Command: "prettier-ignore", Description: "Union-merge curated defaults into .prettierignore"},
			{Command: "prettier-rc", Description: "Merge curated JSON formatting defaults into .prettierrc"},
		},
	}
}
