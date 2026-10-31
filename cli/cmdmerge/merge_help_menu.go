// Package cmdmerge — merge_help_menu.go renders the rich terminal help menu for gitmap merge.
package cmdmerge

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

// RenderMergeHelp displays the styled two-column Merge & Transfer help menu.
func RenderMergeHelp() {
	termout.RenderMenu(buildMergeHelpMenu())
}

func buildMergeHelpMenu() termout.HelpMenu {
	return termout.HelpMenu{
		Title:       "Merge & Transfer Suite (gitmap merge / thelp merge)",
		UsageLines:  buildMergeUsageLines(),
		Sections:    buildMergeHelpSections(),
		FooterFlags: buildMergeFooterFlags(),
		Tips:        buildMergeTips(),
	}
}

func buildMergeUsageLines() []string {
	return []string{
		"gitmap merge-right <left> <right> [flags]  (or gitmap merge right / mr)",
		"gitmap merge-left <left> <right> [flags]   (or gitmap merge left / ml)",
		"gitmap merge-both <left> <right> [flags]   (or gitmap merge both / mb)",
		"gitmap merge-in <target> <inputs...> [flags] (or gitmap merge in / cin)",
		"gitmap commit-in <target> <inputs...> [flags]",
		"gitmap commit-right <left> <right> [flags] (or gitmap cmr)",
		"gitmap commit-left <left> <right> [flags]  (or gitmap cml)",
		"gitmap commit-both <left> <right> [flags]  (or gitmap cmb)",
	}
}

func buildMergeHelpSections() []termout.HelpSection {
	return []termout.HelpSection{
		buildFileMergeSection(),
		buildHistoryMergeSection(),
		buildPRWorkflowSection(),
	}
}

func buildFileMergeSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "One-Way & Two-Way File-Level Merges",
		Entries: []termout.CommandEntry{
			{Command: "merge-right (mr, mrg)", Description: "One-way file merge writing into RIGHT (from LEFT)"},
			{Command: "merge-left (ml)", Description: "One-way file merge writing into LEFT (from RIGHT)"},
			{Command: "merge-both (mb)", Description: "Two-way bidirectional file merge between endpoints"},
			{Command: "merge right <left> <right>", Description: "Subcommand routing for merge-right"},
			{Command: "merge left <left> <right>", Description: "Subcommand routing for merge-left"},
			{Command: "merge both <left> <right>", Description: "Subcommand routing for merge-both"},
		},
	}
}

func buildHistoryMergeSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Chronological Multi-Source History Merges",
		Entries: []termout.CommandEntry{
			{Command: "merge-in (cin)", Description: "Replay & stitch multi-repo commit history into target"},
			{Command: "merge in <target> <inputs>", Description: "Subcommand routing for commit-in"},
			{Command: "commit-in <target> <inputs>", Description: "Append source commits in author-date order"},
			{Command: "commit-right (cmr)", Description: "Replay commits from left repo onto right repo"},
			{Command: "commit-left (cml)", Description: "Replay commits from right repo onto left repo"},
			{Command: "commit-both (cmb)", Description: "Bidirectional commit replay between repositories"},
		},
	}
}

func buildPRWorkflowSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Pull Request & GitHub Integration",
		Entries: []termout.CommandEntry{
			{Command: "merge write (mr)", Description: "Open a GitHub pull request from CLI (--title, --body)"},
			{Command: "pr in", Description: "Multi-source commit-in replay with PR tracking"},
			{Command: "pr list", Description: "Inspect active and merged PR branches in Split-DB"},
		},
	}
}

func buildMergeFooterFlags() []termout.CommandEntry {
	return []termout.CommandEntry{
		{Command: "-y, --yes", Description: "Bypass prompts and auto-accept default conflict resolution"},
		{Command: "--prefer-left", Description: "LEFT always wins on conflicts (default for merge-right)"},
		{Command: "--prefer-right", Description: "RIGHT always wins on conflicts (default for merge-left)"},
		{Command: "--prefer-newer", Description: "Newer mtime wins on conflicts"},
		{Command: "--no-push", Description: "Create local commits without pushing to remote"},
		{Command: "--dry-run", Description: "Preview merge actions without modifying disk"},
		{Command: "-h, --help", Description: "Show this merge help menu"},
	}
}

func buildMergeTips() []string {
	return []string{
		"Use 'gitmap merge-right ./src https://github.com/org/repo -y' to publish changes safely without modifying local.",
		"Use 'gitmap merge-in ./canonical all --default' to stitch versioned history across repositories chronologically.",
		"Use 'gitmap thelp merge' (or 'gitmap thelp merge-right') for detailed flags and examples.",
		"Use 'gitmap cmb' or 'gitmap commit-both' for commit-level bidirectional synchronization.",
	}
}
