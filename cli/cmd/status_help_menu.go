package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/termhelp"
)

// RenderStatusHelp displays the styled two-column Status help menu.
func RenderStatusHelp() {
	termhelp.RenderMenu(buildStatusHelpMenu())
}

func buildStatusHelpMenu() termhelp.HelpMenu {
	return termhelp.HelpMenu{
		Title: "Workspace Repository Status (gitmap status)",
		UsageLines: []string{
			"gitmap status [flags]",
			"gitmap st [flags]",
			"gitmap status --dirty",
		},
		Sections: []termhelp.HelpSection{
			buildStatusScopeSection(),
			buildStatusFormattingSection(),
		},
		FooterFlags: buildStatusFooterFlags(),
		Tips: []string{
			"Quickly audit modified repositories across your entire multi-repo workspace.",
			"Pass '--json' to pipe dirty repo inventory directly to AI scripts.",
		},
	}
}

func buildStatusFooterFlags() []termhelp.CommandEntry {
	return []termhelp.CommandEntry{
		{Command: "-d, --dirty", Description: "Filter to only repositories with uncommitted working changes"},
		{Command: "--ahead", Description: "Filter to repositories with commits ahead of remote tracking"},
		{Command: "--behind", Description: "Filter to repositories with remote commits waiting to pull"},
		{Command: "-j, --json", Description: "Output status records in structured JSON array"},
		{Command: "--table", Description: "Render rich ANSI status table with colored columns"},
		{Command: "--compact", Description: "Single-line per repository compact output"},
		{Command: "-h, --help", Description: "Show this status help menu"},
	}
}

func buildStatusScopeSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Filter Modes",
		Entries: []termhelp.CommandEntry{
			{Command: "all repos", Description: "Scan status across all repositories in current workspace"},
			{Command: "--dirty", Description: "Show only modified, untracked, or staged repositories"},
			{Command: "--ahead / --behind", Description: "Isolate branches out of sync with upstream remotes"},
		},
	}
}

func buildStatusFormattingSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Presentation & Output",
		Entries: []termhelp.CommandEntry{
			{Command: "Table Mode", Description: "Interactive branch, dirty count, ahead/behind matrix"},
			{Command: "JSON Mode", Description: "Machine-readable status telemetry for automation scripts"},
		},
	}
}
