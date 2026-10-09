package cmdbranch

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

// RenderBranchHelp displays the styled two-column Branch help menu.
func RenderBranchHelp() {
	termout.RenderMenu(buildBranchHelpMenu())
}

func buildBranchHelpMenu() termout.HelpMenu {
	return termout.HelpMenu{
		Title: "Branch Resolution & Management (gitmap branch)",
		UsageLines: []string{
			"gitmap branch default [flags]",
			"gitmap b def",
			"gitmap latest-branch",
		},
		Sections: []termout.HelpSection{
			buildBranchOpsSection(),
		},
		FooterFlags: buildBranchFooterFlags(),
		Tips: []string{
			"Quickly resolve default branch (main/master) for any cloned repo.",
			"Works offline and falls back gracefully to 'main' when no remotes are configured.",
		},
	}
}

func buildBranchFooterFlags() []termout.CommandEntry {
	return []termout.CommandEntry{
		{Command: "-j, --json", Description: "Output branch metadata in structured JSON format"},
		{Command: "-h, --help", Description: "Show this branch help menu"},
	}
}

func buildBranchOpsSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Branch Operations",
		Entries: []termout.CommandEntry{
			{Command: "default (def)", Description: "Resolve repository's default branch using symbolic refs"},
			{Command: "latest-branch (lb)", Description: "Inspect the newest branch with recent commit activity"},
		},
	}
}
