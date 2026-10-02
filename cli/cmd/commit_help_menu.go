package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/termhelp"
)

// RenderCommitHelp displays the styled two-column Commit & Transfer help menu.
func RenderCommitHelp() {
	termhelp.RenderMenu(buildCommitHelpMenu())
}

func buildCommitHelpMenu() termhelp.HelpMenu {
	return termhelp.HelpMenu{
		Title: "Semantic Commit & Transfer Suite (gitmap commit)",
		UsageLines: []string{
			"gitmap commit \"<message>\" (flat commit with auto-add)",
			"gitmap cm \"<message>\" (short alias for commit)",
			"gitmap commit-both <left> <right> [flags]",
			"gitmap commit-left <left> <right> [flags]",
			"gitmap commit-right <left> <right> [flags]",
			"gitmap cpf \"<message>\" (commit-push-feature)",
			"gitmap cpb \"<message>\" (commit-push-bug)",
		},
		Sections: []termhelp.HelpSection{
			buildCommitTransferSection(),
			buildCommitAIWorkflowSection(),
			buildCommitPRSection(),
		},
		FooterFlags: buildCommitFooterFlags(),
		Tips: []string{
			"Use 'gitmap commit \"message\"' (or 'gitmap cm') for 1-step auto-stage and flat commit.",
			"Use 'gitmap commit-both' to synchronize commits bidirectionally between repositories.",
			"Use 'gitmap cpf \"feat: description\"' for 1-step stage, commit, and push.",
		},
	}
}

func buildCommitFooterFlags() []termhelp.CommandEntry {
	return []termhelp.CommandEntry{
		{Command: "--interleave", Description: "Merge commits by author date order for commit-both"},
		{Command: "--dry-run", Description: "Preview transfer or commit actions without modifying disk"},
		{Command: "--no-push", Description: "Create local commits without pushing to remote"},
		{Command: "--limit <n>", Description: "Limit number of historical commits transferred"},
		{Command: "--since <date>", Description: "Only transfer commits newer than specified date"},
		{Command: "--mirror", Description: "Exact mirror transfer matching tree SHAs"},
		{Command: "-h, --help", Description: "Show this commit help menu"},
	}
}

func buildCommitTransferSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Multi-Repo Commit Transfers",
		Entries: []termhelp.CommandEntry{
			{Command: "commit-both (cmb)", Description: "Bidirectional commit replay between two repositories"},
			{Command: "commit-left (cml)", Description: "Replay commits from right repo onto left repo"},
			{Command: "commit-right (cmr)", Description: "Replay commits from left repo onto right repo"},
			{Command: "commit-pull (cpull)", Description: "Pull and replay remote commit chain onto local repo"},
		},
	}
}

func buildCommitAIWorkflowSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Commit & AI Workflow Shortcuts",
		Entries: []termhelp.CommandEntry{
			{Command: "commit, cm <msg>", Description: "Flat git commit with auto-stage (git add -A)"},
			{Command: "cpf <msg>", Description: "commit-push-feature: stage all, commit feat, and push"},
			{Command: "cpb <msg>", Description: "commit-push-bug: stage all, commit bugfix, and push"},
			{Command: "cpr <msg>", Description: "commit-push-release: stage, commit release chore, and push"},
			{Command: "pcp <msg>", Description: "pull-commit-push: pull latest, stage, commit, and push"},
		},
	}
}

func buildCommitPRSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Pull Request & Split DB Branches",
		Entries: []termhelp.CommandEntry{
			{Command: "pr list", Description: "Inspect active and merged PR branches in Split-DB"},
			{Command: "pr clean", Description: "Delete merged PR branches locally and remotely"},
			{Command: "pr in", Description: "Generate feature-per-commit pull request branches"},
		},
	}
}
