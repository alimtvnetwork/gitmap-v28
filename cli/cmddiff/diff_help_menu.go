package cmddiff

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

// RenderDiffHelp displays the styled two-column Diff help menu.
func RenderDiffHelp() {
	termout.RenderMenu(buildDiffHelpMenu())
}

func buildDiffHelpMenu() termout.HelpMenu {
	return termout.HelpMenu{
		Title: "Diff & Inspection Suite (gitmap diff)",
		UsageLines: []string{
			"gitmap diff [target] [flags]",
			"gitmap diff-profiles <profile1> <profile2>",
		},
		Sections: []termout.HelpSection{
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

func buildDiffFooterFlags() []termout.CommandEntry {
	return []termout.CommandEntry{
		{Command: "--stat", Description: "Show compact change statistics (insertions/deletions)"},
		{Command: "--cached, --staged", Description: "Diff staged changes against repository HEAD"},
		{Command: "--name-only", Description: "Show only names of modified files"},
		{Command: "-j, --json", Description: "Emit diff metadata in structured JSON format"},
		{Command: "-h, --help", Description: "Show this diff help menu"},
	}
}

func buildDiffOpsSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Repository Diff Actions",
		Entries: []termout.CommandEntry{
			{Command: "diff [path]", Description: "Inspect unstaged file modifications with colored diffs"},
			{Command: "diff --cached", Description: "Inspect staged changes prepared for commit"},
			{Command: "diff --stat", Description: "Compact overview of affected files and line deltas"},
		},
	}
}

func buildDiffProfilesSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Profile Comparison",
		Entries: []termout.CommandEntry{
			{Command: "diff-profiles <p1> <p2>", Description: "Compare software packages, aliases, and settings between profiles"},
		},
	}
}
