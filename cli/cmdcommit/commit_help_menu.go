package cmdcommit

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

// RenderCommitHelp displays the styled two-column Commit & Transfer help menu.
func RenderCommitHelp() {
	termout.RenderMenu(buildCommitHelpMenu())
}

func buildCommitHelpMenu() termout.HelpMenu {
	return termout.HelpMenu{
		Title: "Semantic Commit & Transfer Suite (gitmap commit)",
		UsageLines: []string{
			"gitmap commit \"<message>\" (flat commit with auto-add)",
			"gitmap cm \"<message>\" (short alias for commit)",
			"gitmap commit-both <left> <right> [flags]",
			"gitmap commit-left <left> <right> [flags]",
			"gitmap commit-right <left> <right> [flags]",
			"gitmap cpf \"<message>\" (commit-push-feature)",
			"gitmap cpb \"<message>\" (commit-push-bug)",
			"gitmap cpc \"<message>\" (commit-push-chore)",
			"gitmap cpr \"<message>\" (commit-push-release)",
			"gitmap pcp \"<message>\" (pull-commit-push)",
		},
		Sections: []termout.HelpSection{
			buildCommitTransferSection(),
			buildCommitAIWorkflowSection(),
			buildCommitPRSection(),
		},
		FooterFlags: buildCommitFooterFlags(),
		Tips: []string{
			"Use 'gitmap commit \"message\"' (or 'gitmap cm') for 1-step auto-stage and flat commit.",
			"Use 'gitmap commit-both' to synchronize commits bidirectionally between repositories.",
			"Use 'gitmap cpf \"feat: description\"' for 1-step stage, commit, and push.",
			"Use 'gitmap cpc \"chore: description\"' for 1-step chore stage, commit, and push.",
		},
	}
}

func buildCommitFooterFlags() []termout.CommandEntry {
	return []termout.CommandEntry{
		{Command: "--interleave", Description: "Merge commits by author date order for commit-both"},
		{Command: "--dry-run", Description: "Preview transfer or commit actions without modifying disk"},
		{Command: "--no-push", Description: "Create local commits without pushing to remote"},
		{Command: "--limit <n>", Description: "Limit number of historical commits transferred"},
		{Command: "--since <date>", Description: "Only transfer commits newer than specified date"},
		{Command: "--mirror", Description: "Exact mirror transfer matching tree SHAs"},
		{Command: "-h, --help", Description: "Show this commit help menu"},
	}
}

func buildCommitTransferSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Multi-Repo Commit Transfers",
		Entries: []termout.CommandEntry{
			{Command: "commit-both (cmb)", Description: "Bidirectional commit replay between two repositories"},
			{Command: "commit-left (cml)", Description: "Replay commits from right repo onto left repo"},
			{Command: "commit-right (cmr)", Description: "Replay commits from left repo onto right repo"},
			{Command: "commit-pull (cpull)", Description: "Pull and replay remote commit chain onto local repo"},
		},
	}
}

func buildCommitAIWorkflowSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Commit & AI Workflow Shortcuts",
		Entries: []termout.CommandEntry{
			{Command: "commit, cm <msg>", Description: "commit: flat git commit with auto-stage (git add -A)"},
			{Command: "cpf <msg>", Description: "commit-push-feature: stage all, commit feat, and push"},
			{Command: "cpb <msg>", Description: "commit-push-bug: stage all, commit bugfix, and push"},
			{Command: "cpc <msg>", Description: "commit-push-chore: stage all, commit chore/maintenance, and push"},
			{Command: "cpr <msg>", Description: "commit-push-release: stage all, commit release chore, and push"},
			{Command: "pcp <msg>", Description: "pull-commit-push: pull latest rebase, stage, commit, and push"},
			{Command: "pas", Description: "pull-all-ssh: pull all repositories using SSH transport"},
		},
	}
}

func buildCommitPRSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Pull Request & Split DB Branches",
		Entries: []termout.CommandEntry{
			{Command: "pr list", Description: "Inspect active and merged PR branches in Split-DB"},
			{Command: "pr clean", Description: "Delete merged PR branches locally and remotely"},
			{Command: "pr in", Description: "Generate feature-per-commit pull request branches"},
		},
	}
}
