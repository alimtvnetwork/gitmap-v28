package cmdstatus

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

// RenderStatusHelp displays the styled two-column Status help menu.
func RenderStatusHelp() {
	termout.RenderMenu(buildStatusHelpMenu())
}

func buildStatusHelpMenu() termout.HelpMenu {
	return termout.HelpMenu{
		Title: "Workspace Repository Status (gitmap status)",
		UsageLines: []string{
			"gitmap status [flags]",
			"gitmap st [flags]",
			"gitmap status --dirty",
		},
		Sections: []termout.HelpSection{
			buildStatusScopeSection(),
			buildStatusFormattingSection(),
		},
		FooterFlags: buildStatusFooterFlags(),
		Tips: []string{
			"Quickly audit modified repositories across your entire multi-repo workspace.",
			"Pass '--json' to pipe dirty repo inventory directly to AI fspath.",
		},
	}
}

func buildStatusFooterFlags() []termout.CommandEntry {
	return []termout.CommandEntry{
		{Command: "-d, --dirty", Description: "Filter to only repositories with uncommitted working changes"},
		{Command: "--ahead", Description: "Filter to repositories with commits ahead of remote tracking"},
		{Command: "--behind", Description: "Filter to repositories with remote commits waiting to pull"},
		{Command: "-j, --json", Description: "Output status records in structured JSON array"},
		{Command: "--table", Description: "Render rich ANSI status table with colored columns"},
		{Command: "--compact", Description: "Single-line per repository compact output"},
		{Command: "-h, --help", Description: "Show this status help menu"},
	}
}

func buildStatusScopeSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Filter Modes",
		Entries: []termout.CommandEntry{
			{Command: "all repos", Description: "Scan status across all repositories in current workspace"},
			{Command: "--dirty", Description: "Show only modified, untracked, or staged repositories"},
			{Command: "--ahead / --behind", Description: "Isolate branches out of sync with upstream remotes"},
		},
	}
}

func buildStatusFormattingSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Presentation & Output",
		Entries: []termout.CommandEntry{
			{Command: "Table Mode", Description: "Interactive branch, dirty count, ahead/behind matrix"},
			{Command: "JSON Mode", Description: "Machine-readable status telemetry for automation scripts"},
		},
	}
}
