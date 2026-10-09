package cmdfixgit

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

// RenderFixGitHelp displays the styled two-column FixGit help menu.
func RenderFixGitHelp() {
	termout.RenderMenu(buildFixGitHelpMenu())
}

func buildFixGitHelpMenu() termout.HelpMenu {
	return termout.HelpMenu{
		Title: "Git Diagnostic & Self-Healing Engine (gitmap fix-git)",
		UsageLines: []string{
			"gitmap fix-git [path] [flags]",
			"gitmap fg [path] [flags]",
			"gitmap --fix-git [path]",
		},
		Sections: []termout.HelpSection{
			buildFixGitTargetsSection(),
			buildFixGitModesSection(),
		},
		FooterFlags: buildFixGitFooterFlags(),
		Tips: []string{
			"Run 'gitmap fix-git --dry-run' to inspect issues without modifying disk.",
			"Use 'gitmap fg' inside any repository with stale lockfiles or bad permissions.",
		},
	}
}

func buildFixGitFooterFlags() []termout.CommandEntry {
	return []termout.CommandEntry{
		{Command: "-n, --dry-run", Description: "Preview detected issues without applying repairs"},
		{Command: "-v, --verbose", Description: "Emit verbose diagnostic traces and ACL status"},
		{Command: "--json", Description: "Output remediation summary as JSON"},
		{Command: "-h, --help", Description: "Show this fix-git help menu"},
	}
}

func buildFixGitTargetsSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Diagnostic & Remediation Targets",
		Entries: []termout.CommandEntry{
			{Command: "--permissions (--perms)", Description: "Fix Windows NTFS ACLs and remove read-only flags"},
			{Command: "--locks (--locks-only)", Description: "Remove stale .git/index.lock and HEAD.lock files"},
			{Command: "--index (--index-only)", Description: "Safely reconstruct corrupted or 0-byte Git index"},
			{Command: "--safe-dir", Description: "Register repo in Git global safe.directory config"},
		},
	}
}

func buildFixGitModesSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Execution Modes",
		Entries: []termout.CommandEntry{
			{Command: "[path]", Description: "Target repository directory (default: current directory)"},
			{Command: "--dry-run (-n)", Description: "Scan and report issues without altering files"},
			{Command: "--json", Description: "Format diagnostic output as structured JSON"},
		},
	}
}
