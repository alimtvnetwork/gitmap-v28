package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/termhelp"
)

// RenderBranchHelp displays the styled two-column Branch help menu.
func RenderBranchHelp() {
	termhelp.RenderMenu(buildBranchHelpMenu())
}

func buildBranchHelpMenu() termhelp.HelpMenu {
	return termhelp.HelpMenu{
		Title: "Branch Resolution & Management (gitmap branch)",
		UsageLines: []string{
			"gitmap branch default [flags]",
			"gitmap b def",
			"gitmap latest-branch",
		},
		Sections: []termhelp.HelpSection{
			buildBranchOpsSection(),
		},
		FooterFlags: buildBranchFooterFlags(),
		Tips: []string{
			"Quickly resolve default branch (main/master) for any cloned repo.",
			"Works offline and falls back gracefully to 'main' when no remotes are configured.",
		},
	}
}

func buildBranchFooterFlags() []termhelp.CommandEntry {
	return []termhelp.CommandEntry{
		{Command: "-j, --json", Description: "Output branch metadata in structured JSON format"},
		{Command: "-h, --help", Description: "Show this branch help menu"},
	}
}

func buildBranchOpsSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Branch Operations",
		Entries: []termhelp.CommandEntry{
			{Command: "default (def)", Description: "Resolve repository's default branch using symbolic refs"},
			{Command: "latest-branch (lb)", Description: "Inspect the newest branch with recent commit activity"},
		},
	}
}
