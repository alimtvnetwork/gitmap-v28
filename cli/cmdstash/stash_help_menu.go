package cmdstash

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/termhelp"
)

// RenderStashHelp displays the styled two-column Stash / Fix help menu.
func RenderStashHelp() {
	termhelp.RenderMenu(buildStashHelpMenu())
}

func buildStashHelpMenu() termhelp.HelpMenu {
	return termhelp.HelpMenu{
		Title: "Workspace Remediation & Stash Engine (gitmap fix / stash)",
		UsageLines: []string{
			"gitmap fix [target] [flags]",
			"gitmap stash [flags]",
			"gitmap wip [flags]",
			"gitmap discard [flags]",
			"gitmap fix --all",
		},
		Sections: []termhelp.HelpSection{
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

func buildStashFooterFlags() []termhelp.CommandEntry {
	return []termhelp.CommandEntry{
		{Command: "-a, --all", Description: "Apply remediation recipe to all pending repositories"},
		{Command: "-p, --prompt", Description: "Interactive prompt to choose remediation per repository"},
		{Command: "-y, --yes", Description: "Bypass confirmation prompts automatically"},
		{Command: "--dry-run", Description: "Preview remediation actions without modifying disk"},
		{Command: "-h, --help", Description: "Show this remediation/stash help menu"},
	}
}

func buildStashActionsSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Remediation Recipes",
		Entries: []termhelp.CommandEntry{
			{Command: "stash [1]", Description: "Save working modifications to git stash stack"},
			{Command: "wip [2]", Description: "Commit all working changes to a temporary WIP commit"},
			{Command: "discard [3]", Description: "Hard reset and clean untracked working tree changes"},
		},
	}
}

func buildStashBulkSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Management & Diagnostics",
		Entries: []termhelp.CommandEntry{
			{Command: "fix ls", Description: "List all repositories requiring remediation"},
			{Command: "fix agy", Description: "Autonomous AI-assisted workspace remediation"},
		},
	}
}
