package cmdscan

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

// RenderScanHelp displays the styled two-column Scan help menu.
func RenderScanHelp() {
	termout.RenderMenu(buildScanHelpMenu())
}

func buildScanHelpMenu() termout.HelpMenu {
	return termout.HelpMenu{
		Title: "Repository Discovery Scanner (gitmap scan)",
		UsageLines: []string{
			"gitmap scan [dir] [flags]",
			"gitmap s [dir] [flags]",
			"gitmap scan --fix",
			"gitmap scan ~ --force-include .oh-my-zsh",
			"gitmap scan ~ -fi .oh-my-zsh,node_modules",
			"gitmap scan /home/user --force-include all",
		},
		Sections: []termout.HelpSection{
			buildScanOutputSection(),
			buildScanWalkerSection(),
			buildScanIntegrationSection(),
		},
		FooterFlags: buildScanFooterFlags(),
		Tips: []string{
			"Run 'gitmap scan . --output json' to generate a machine-readable repo catalog.",
			"Use 'gitmap scan --fix' to reconcile local disk with GitMap database.",
			"Pass '--max-depth 2' to quickly scan shallow project directories.",
			"Pass '--force-include .oh-my-zsh' (alias: -fi) to force scanning excluded directories or 'all' for everything.",
		},
	}
}

func buildScanFooterFlags() []termout.CommandEntry {
	return []termout.CommandEntry{
		{Command: "--quiet", Description: "Suppress interactive walker spinner and clone hints"},
		{Command: "--open", Description: "Automatically open output directory after scan completes"},
		{Command: "-v, --verbose", Description: "Emit verbose directory traversal and probe logs"},
		{Command: "-h, --help", Description: "Show this scanner help menu"},
	}
}

func buildScanOutputSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Output & Formatting",
		Entries: []termout.CommandEntry{
			{Command: "--output <mode>", Description: "Output format: terminal (default), csv, or json"},
			{Command: "--output-path <dir>", Description: "Directory to write scan artifacts (.gitmap/output)"},
			{Command: "--manifest <dir>", Description: "Alias for --output-path (aligned with reclone)"},
			{Command: "--relative-root <dir>", Description: "Pin base directory for byte-stable relative paths"},
		},
	}
}
