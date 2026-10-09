package cmdsearch

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

// RenderSearchHelp displays the styled two-column Search help menu.
func RenderSearchHelp() {
	termout.RenderMenu(buildSearchHelpMenu())
}

func buildSearchHelpMenu() termout.HelpMenu {
	return termout.HelpMenu{
		Title: "Global Indexed Search (gitmap search)",
		UsageLines: []string{
			"gitmap search <query> [flags]",
			"gitmap search history",
			"gitmap search clean",
			"gitmap aum search <query> [dir] [flags]",
		},
		Sections: []termout.HelpSection{
			buildSearchOpsSection(),
			buildSearchCacheSection(),
		},
		FooterFlags: buildSearchFooterFlags(),
		Tips: []string{
			"Hot queries are cached in SQLite with instant DH2D lookup (< 0.05ms).",
			"Use 'gitmap search history' to inspect query frequency and hit counts.",
		},
	}
}

func buildSearchFooterFlags() []termout.CommandEntry {
	return []termout.CommandEntry{
		{Command: "-j, --json", Description: "Emit results in structured JSON format"},
		{Command: "-p, --path <dir>", Description: "Restrict search to target subdirectory"},
		{Command: "-e, --ext <exts>", Description: "Comma-separated extensions filter (.go, .ts)"},
		{Command: "-r, --regex", Description: "Evaluate search query as regular expression"},
		{Command: "-i, --ignore-case", Description: "Case-insensitive matching fast-path"},
		{Command: "-l, --limit <n>", Description: "Maximum number of search results to return"},
		{Command: "-h, --help", Description: "Show this search help menu"},
	}
}

func buildSearchOpsSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Search Operations",
		Entries: []termout.CommandEntry{
			{Command: "<query>", Description: "Multi-core streaming keyword search across repository"},
			{Command: "history", Description: "Display recent search history and hot query rankings"},
			{Command: "clean", Description: "Purge search history cache and vacuum search database"},
		},
	}
}

func buildSearchCacheSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "High-Speed SQLite Caching",
		Entries: []termout.CommandEntry{
			{Command: "DH2D SQLite Cache", Description: "Deterministic SQL ID lookup caching verified matches"},
			{Command: "RAM Hot-Cache", Description: "Auto-promotes frequent queries for microsecond retrieval"},
		},
	}
}
