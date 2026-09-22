package cmdfixgit

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/termhelp"
)

// RenderFixGitHelp displays the styled two-column FixGit help menu.
func RenderFixGitHelp() {
	termhelp.RenderMenu(buildFixGitHelpMenu())
}

func buildFixGitHelpMenu() termhelp.HelpMenu {
	return termhelp.HelpMenu{
		Title: "Git Diagnostic & Self-Healing Engine (gitmap fix-git)",
		UsageLines: []string{
			"gitmap fix-git [path] [flags]",
			"gitmap fg [path] [flags]",
			"gitmap --fix-git [path]",
		},
		Sections: []termhelp.HelpSection{
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

func buildFixGitFooterFlags() []termhelp.CommandEntry {
	return []termhelp.CommandEntry{
		{Command: "-n, --dry-run", Description: "Preview detected issues without applying repairs"},
		{Command: "-v, --verbose", Description: "Emit verbose diagnostic traces and ACL status"},
		{Command: "--json", Description: "Output remediation summary as JSON"},
		{Command: "-h, --help", Description: "Show this fix-git help menu"},
	}
}

func buildFixGitTargetsSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Diagnostic & Remediation Targets",
		Entries: []termhelp.CommandEntry{
			{Command: "--permissions (--perms)", Description: "Fix Windows NTFS ACLs and remove read-only flags"},
			{Command: "--locks (--locks-only)", Description: "Remove stale .git/index.lock and HEAD.lock files"},
			{Command: "--index (--index-only)", Description: "Safely reconstruct corrupted or 0-byte Git index"},
			{Command: "--safe-dir", Description: "Register repo in Git global safe.directory config"},
		},
	}
}

func buildFixGitModesSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Execution Modes",
		Entries: []termhelp.CommandEntry{
			{Command: "[path]", Description: "Target repository directory (default: current directory)"},
			{Command: "--dry-run (-n)", Description: "Scan and report issues without altering files"},
			{Command: "--json", Description: "Format diagnostic output as structured JSON"},
		},
	}
}
