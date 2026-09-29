package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/termhelp"
)

// RenderDiffHelp displays the styled two-column Diff help menu.
func RenderDiffHelp() {
	termhelp.RenderMenu(buildDiffHelpMenu())
}

func buildDiffHelpMenu() termhelp.HelpMenu {
	return termhelp.HelpMenu{
		Title: "Diff & Inspection Suite (gitmap diff)",
		UsageLines: []string{
			"gitmap diff [target] [flags]",
			"gitmap diff-profiles <profile1> <profile2>",
		},
		Sections: []termhelp.HelpSection{
			buildDiffOpsSection(),
			buildDiffProfilesSection(),
		},
		FooterFlags: buildDiffFooterFlags(),
		Tips: []string{
			"Inspect uncommitted modifications across polyglot repositories.",
			"Use 'gitmap diff-profiles' to detect drift between workstation setups.",
		},
	}
}

func buildDiffFooterFlags() []termhelp.CommandEntry {
	return []termhelp.CommandEntry{
		{Command: "--stat", Description: "Show compact change statistics (insertions/deletions)"},
		{Command: "--cached, --staged", Description: "Diff staged changes against repository HEAD"},
		{Command: "--name-only", Description: "Show only names of modified files"},
		{Command: "-j, --json", Description: "Emit diff metadata in structured JSON format"},
		{Command: "-h, --help", Description: "Show this diff help menu"},
	}
}

func buildDiffOpsSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Repository Diff Actions",
		Entries: []termhelp.CommandEntry{
			{Command: "diff [path]", Description: "Inspect unstaged file modifications with colored diffs"},
			{Command: "diff --cached", Description: "Inspect staged changes prepared for commit"},
			{Command: "diff --stat", Description: "Compact overview of affected files and line deltas"},
		},
	}
}

func buildDiffProfilesSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Profile Comparison",
		Entries: []termhelp.CommandEntry{
			{Command: "diff-profiles <p1> <p2>", Description: "Compare software packages, aliases, and settings between profiles"},
		},
	}
}
