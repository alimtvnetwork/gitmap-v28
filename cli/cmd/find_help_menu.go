package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/termhelp"
)

// RenderFindHelp displays the styled two-column Find/Find-Files help menu.
func RenderFindHelp() {
	termhelp.RenderMenu(buildFindHelpMenu())
}

func buildFindHelpMenu() termhelp.HelpMenu {
	return termhelp.HelpMenu{
		Title: "Fast File Locator (gitmap find / find-files)",
		UsageLines: []string{
			"gitmap find <pattern> [flags]",
			"gitmap ff <exact-name> [-ext <ext>]",
			"gitmap find-files-any <substr> [-ext <ext>]",
			"gitmap find-files-startswith <prefix>",
			"gitmap find-files-endswith <suffix>",
		},
		Sections: []termhelp.HelpSection{
			buildFindOpsSection(),
			buildFindFiltersSection(),
		},
		FooterFlags: buildFindFooterFlags(),
		Tips: []string{
			"Use 'gitmap ff filename.go' for ultra-fast single file lookup.",
			"Wildcards like 'gitmap find \"*config*.json\"' scan 2,900+ files in milliseconds.",
		},
	}
}

func buildFindFooterFlags() []termhelp.CommandEntry {
	return []termhelp.CommandEntry{
		{Command: "-e, -ext <ext>", Description: "Filter results by file extension (e.g. go, ts, json)"},
		{Command: "-l, --limit <n>", Description: "Limit maximum number of returned file paths"},
		{Command: "-j, --json", Description: "Emit matched paths as structured JSON array"},
		{Command: "--full-path", Description: "Output absolute filesystem paths"},
		{Command: "-h, --help", Description: "Show this find help menu"},
	}
}

func buildFindOpsSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Locator Subcommands & Shortcuts",
		Entries: []termhelp.CommandEntry{
			{Command: "find / f", Description: "Universal wildcard and glob file search"},
			{Command: "find-files / ff", Description: "Exact filename match with fast index lookup"},
			{Command: "find-files-any / ffa", Description: "Substring matching within filenames"},
			{Command: "find-files-startswith / ffs", Description: "Prefix matching on filenames"},
			{Command: "find-files-endswith / ffe", Description: "Suffix matching on filenames (e.g. _test.go)"},
		},
	}
}

func buildFindFiltersSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Traversal & Performance",
		Entries: []termhelp.CommandEntry{
			{Command: "Zero-Alloc Scanner", Description: "Iterative depth traversal avoiding heap allocations"},
			{Command: "Auto Exclusions", Description: "Automatically ignores node_modules, .git, vendor, bin"},
		},
	}
}
