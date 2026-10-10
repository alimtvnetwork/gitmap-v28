// Package cmdfixreleasetags provides formatted two-column help menu structures
// for the release tag management command suite.
package cmdfixreleasetags

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

// BuildFixReleaseTagsHelpMenu constructs the command help menu.
func BuildFixReleaseTagsHelpMenu() termout.HelpMenu {
	return termout.HelpMenu{
		Title: "Release Tag & Broken Release Cleanup Engine (gitmap fix release tags)",
		UsageLines: []string{
			"gitmap fix release tags [flags]",
			"gitmap fix-release-tags [flags]",
			"gitmap frt [flags]",
		},
		Sections: []termout.HelpSection{
			buildAuditAndScopeSection(),
			buildExecutionControlSection(),
			buildOutputFormattingSection(),
		},
		FooterFlags: buildHelpFooterFlags(),
		Tips: []string{
			"Always run with '--dry-run' first to preview orphan tags without deleting.",
			"Use '-y' or '--confirm' in automated GitHub Actions / CI cleanup workflows.",
			"Tags with passing CI/CD or published binary assets are automatically protected.",
		},
	}
}

func buildAuditAndScopeSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Audit Scope & Targets",
		Color: constants.ColorCyan,
		Entries: []termout.CommandEntry{
			{Command: "--repo, -r <path>", Description: "Target repository path or slug (default: current directory)"},
			{Command: "--local-only", Description: "Audit and delete local git tags only (skip remote/GitHub)"},
			{Command: "--remote-only", Description: "Audit and delete remote origin tags and GitHub releases only"},
		},
	}
}

func buildExecutionControlSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Execution & Safety Controls",
		Color: constants.ColorGreen,
		Entries: []termout.CommandEntry{
			{Command: "-n, --dry-run", Description: "Simulate audit and preview actions without modifying repository"},
			{Command: "-y, --yes, --confirm", Description: "Bypass interactive confirmation prompt and execute deletions"},
		},
	}
}

func buildOutputFormattingSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Output & Formatting",
		Color: constants.ColorYellow,
		Entries: []termout.CommandEntry{
			{Command: "--json", Description: "Emit audit report and deletion results as structured JSON"},
			{Command: "-v, --verbose", Description: "Display detailed API HTTP payloads and git command execution"},
		},
	}
}

func buildHelpFooterFlags() []termout.CommandEntry {
	return []termout.CommandEntry{
		{Command: "-h, --help", Description: "Show this help menu and exit"},
	}
}

// RenderFixReleaseTagsHelp renders the help menu to standard output.
func RenderFixReleaseTagsHelp() {
	termout.RenderMenu(BuildFixReleaseTagsHelpMenu())
}
