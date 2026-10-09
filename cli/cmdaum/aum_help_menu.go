package cmdaum

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

// RenderAumHelp displays the styled two-column AUM / Automation help menu.
func RenderAumHelp() {
	termout.RenderMenu(buildAumHelpMenu())
}

func buildAumHelpMenu() termout.HelpMenu {
	return termout.HelpMenu{
		Title: "Polyglot Automation Engine (gitmap aum / automation)",
		UsageLines: []string{
			"gitmap aum <subcommand> [flags]",
			"gitmap automation <subcommand> [flags]",
			"gitmap aum search <query> [dir] [flags]",
			"gitmap aum benchmark [target]",
		},
		Sections: []termout.HelpSection{
			buildAumToolsSection(),
			buildAumFormattingSection(),
		},
		FooterFlags: buildAumFooterFlags(),
		Tips: []string{
			"Multi-core streaming search with lazy regex and zero-allocation fast path.",
			"Run 'gitmap aum benchmark' to compare native Go vs Python performance.",
		},
	}
}

func buildAumFooterFlags() []termout.CommandEntry {
	return []termout.CommandEntry{
		{Command: "-w, --workers <n>", Description: "Number of concurrent worker threads (default: CPU cores)"},
		{Command: "--max-json-kb <n>", Description: "Maximum JSON size in KB before auto-exclusion (default: 500)"},
		{Command: "-e, --ext <exts>", Description: "Filter search by file extensions (e.g. .go, .ts)"},
		{Command: "-r, --regex", Description: "Use regular expression matching"},
		{Command: "-i, --ignore-case", Description: "Perform case-insensitive matching"},
		{Command: "-h, --help", Description: "Show this automation help menu"},
	}
}

func buildAumToolsSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Automation Subcommands",
		Entries: []termout.CommandEntry{
			{Command: "search <pattern> [dir]", Description: "Multi-core streaming search with lazy regex and literal fast path"},
			{Command: "locate <binary>", Description: "Fast binary and toolchain locator (vswhere & vcvars fast path)"},
			{Command: "guard", Description: "Safety probes for oversized JSONs and binary null-byte scans"},
			{Command: "sequence", Description: "Markdown numbering gap detector & H1 title validator"},
			{Command: "benchmark [target]", Description: "Run side-by-side Go vs Python execution benchmarks"},
		},
	}
}

func buildAumFormattingSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Normalization & Hygiene",
		Entries: []termout.CommandEntry{
			{Command: "newlines [dir]", Description: "Polyglot newline normalization across repositories"},
			{Command: "format-go", Description: "Fast Go AST formatting and import grouping"},
			{Command: "clean-artifacts", Description: "Zero-storage cleanup of test outputs, binaries, and caches"},
		},
	}
}
