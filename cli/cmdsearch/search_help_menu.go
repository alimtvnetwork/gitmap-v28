package cmdsearch

import (
	searchconstants "github.com/alimtvnetwork/gitmap-v28/cli/cmdsearch/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

// RenderSearchHelp displays the styled two-column Search help menu.
func RenderSearchHelp() {
	termout.RenderMenu(buildSearchHelpMenu())
}

func buildSearchHelpMenu() termout.HelpMenu {
	return termout.HelpMenu{
		Title: searchconstants.HelpTitleGlobalSearch,
		UsageLines: []string{
			searchconstants.UsageSearchFlags,
			searchconstants.UsageSearchHistory,
			searchconstants.UsageSearchCleanSimple,
			searchconstants.UsageAUMSearch,
		},
		Sections: []termout.HelpSection{
			buildSearchOpsSection(),
			buildSearchCacheSection(),
		},
		FooterFlags: buildSearchFooterFlags(),
		Tips: []string{
			searchconstants.TipHotQueryCache,
			searchconstants.TipHistoryInspect,
		},
	}
}

func buildSearchFooterFlags() []termout.CommandEntry {
	return []termout.CommandEntry{
		{Command: searchconstants.FlagHelpEntryJSON, Description: searchconstants.DescJSONFormat},
		{Command: searchconstants.FlagHelpEntryPath, Description: searchconstants.DescPathFilter},
		{Command: searchconstants.FlagHelpEntryExt, Description: searchconstants.DescExtFilter},
		{Command: searchconstants.FlagHelpEntryRegex, Description: searchconstants.DescRegexMatch},
		{Command: searchconstants.FlagHelpEntryIgnoreCase, Description: searchconstants.DescIgnoreCase},
		{Command: searchconstants.FlagHelpEntryLimit, Description: searchconstants.DescResultLimit},
		{Command: searchconstants.FlagHelpEntryHelp, Description: searchconstants.DescHelpMenu},
	}
}

func buildSearchOpsSection() termout.HelpSection {
	return termout.HelpSection{
		Title: searchconstants.SectionSearchOperations,
		Entries: []termout.CommandEntry{
			{Command: searchconstants.HelpEntryQuery, Description: searchconstants.DescStreamingSearch},
			{Command: searchconstants.HelpEntryHistory, Description: searchconstants.DescSearchHistory},
			{Command: searchconstants.HelpEntryClean, Description: searchconstants.DescSearchClean},
		},
	}
}

func buildSearchCacheSection() termout.HelpSection {
	return termout.HelpSection{
		Title: searchconstants.SectionSQLiteCaching,
		Entries: []termout.CommandEntry{
			{Command: searchconstants.HelpEntryDH2DCache, Description: searchconstants.DescDH2DCache},
			{Command: searchconstants.HelpEntryRAMHotCache, Description: searchconstants.DescRAMHotCache},
		},
	}
}
