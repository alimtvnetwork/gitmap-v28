package cmdupdate

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

// RenderUpdateRichHelp displays the styled two-column Update help menu.
func RenderUpdateRichHelp() {
	termout.RenderMenu(buildUpdateHelpMenu())
}

func buildUpdateHelpMenu() termout.HelpMenu {
	return termout.HelpMenu{
		Title: "GitMap Self & Fleet Update (gitmap update / ua)",
		UsageLines: []string{
			"gitmap update [target|version] [flags]",
			"gitmap ua (alias: update all / update-all)",
			"gitmap update all [--target <node>] [--except <nodes>]",
			"gitmap update <package> [node] [flags]",
		},
		Sections: []termout.HelpSection{
			buildUpdateFleetSection(),
			buildUpdateLocalSection(),
		},
		FooterFlags: buildUpdateFooterFlags(),
		Tips: []string{
			"Use 'gitmap ua' for ultra-fast parallel fleet updates across all SSH nodes.",
			"External machines communicate via JSON so warnings and installer noise are suppressed.",
			"Run 'gitmap update ls' to inspect versions across your cluster fleet.",
		},
	}
}

func buildUpdateFleetSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Fleet & Cluster Node Updates (SSH)",
		Entries: []termout.CommandEntry{
			{Command: "ua (update all)", Description: "Update GitMap CLI across all cluster nodes in parallel"},
			{Command: "update-all (updateall)", Description: "Hyphenated alias for updating all fleet nodes"},
			{Command: "update -all", Description: "Flag alias for updating all fleet nodes"},
			{Command: "update <node>", Description: "Update GitMap on specific node by alias, IP, or seq"},
			{Command: "update ls", Description: "Query software inventory and version across all fleet nodes"},
			{Command: "update agm", Description: "Update Antigravity-Manager across cluster nodes"},
		},
	}
}

func buildUpdateLocalSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Local Self Update",
		Entries: []termout.CommandEntry{
			{Command: "update", Description: "Check for updates and download canonical installer"},
			{Command: "update --force (-f)", Description: "Force reinstallation of latest release binary"},
			{Command: "update <version>", Description: "Install a specific GitMap version (e.g. v6.360.0)"},
			{Command: "update cleanup", Description: "Clean temporary files and stale handoff copies"},
			{Command: "update --dry-run", Description: "Preview update plan without downloading or writing"},
		},
	}
}

func buildUpdateFooterFlags() []termout.CommandEntry {
	return []termout.CommandEntry{
		{Command: "--all, -a, -all", Description: "Target all registered cluster nodes in parallel"},
		{Command: "--target, -t <node>", Description: "Target specific node by alias, IP, sequence, or ID"},
		{Command: "--except, -e <nodes>", Description: "Exclude specific nodes (comma-separated list)"},
		{Command: "--force, -f", Description: "Force reinstall even if already at latest version"},
		{Command: "--dry-run", Description: "Preview update actions without executing changes"},
		{Command: "-h, --help", Description: "Show this update help menu"},
	}
}
