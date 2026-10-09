package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

// RenderStashHelp displays the styled two-column Stash / Fix help menu.
func RenderStashHelp() {
	termout.RenderMenu(buildStashHelpMenu())
}

func buildStashHelpMenu() termout.HelpMenu {
	return termout.HelpMenu{
		Title: "Workspace Remediation & Stash Engine (gitmap fix / stash)",
		UsageLines: []string{
			"gitmap fix [target] [flags]",
			"gitmap stash [flags]",
			"gitmap wip [flags]",
			"gitmap discard [flags]",
			"gitmap fix --all",
		},
		Sections: []termout.HelpSection{
			buildStashActionsSection(),
			buildStashBulkSection(),
		},
		FooterFlags: buildStashFooterFlags(),
		Tips: []string{
			"Run 'gitmap fix --all' to apply remediation across all pending repos.",
			"Use 'gitmap stash' to safely save uncommitted work before updates or branch switches.",
		},
	}
}

func buildStashFooterFlags() []termout.CommandEntry {
	return []termout.CommandEntry{
		{Command: "-a, --all", Description: "Apply remediation recipe to all pending repositories"},
		{Command: "-p, --prompt", Description: "Interactive prompt to choose remediation per repository"},
		{Command: "-y, --yes", Description: "Bypass confirmation prompts automatically"},
		{Command: "--dry-run", Description: "Preview remediation actions without modifying disk"},
		{Command: "-h, --help", Description: "Show this remediation/stash help menu"},
	}
}

func buildStashActionsSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Remediation Recipes",
		Entries: []termout.CommandEntry{
			{Command: "stash [1]", Description: "Save working modifications to git stash stack"},
			{Command: "wip [2]", Description: "Commit all working changes to a temporary WIP commit"},
			{Command: "discard [3]", Description: "Hard reset and clean untracked working tree changes"},
		},
	}
}

func buildStashBulkSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Management & Diagnostics",
		Entries: []termout.CommandEntry{
			{Command: "fix ls", Description: "List all repositories requiring remediation"},
			{Command: "fix agy", Description: "Autonomous AI-assisted workspace remediation"},
		},
	}
}
