package cmddiff

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/cmddiff/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

// RenderDiffHelp displays the styled two-column Diff help menu.
func RenderDiffHelp() {
	termout.RenderMenu(buildDiffHelpMenu())
}

func buildDiffHelpMenu() termout.HelpMenu {
	return termout.HelpMenu{
		Title: constants.HelpTitleDiff,
		UsageLines: []string{
			constants.HelpUsageDiffTarget,
			constants.HelpUsageDiffProfiles,
		},
		Sections: []termout.HelpSection{
			buildDiffOpsSection(),
			buildDiffProfilesSection(),
		},
		FooterFlags: buildDiffFooterFlags(),
		Tips: []string{
			constants.HelpTipUncommitted,
			constants.HelpTipDrift,
		},
	}
}

func buildDiffFooterFlags() []termout.CommandEntry {
	return []termout.CommandEntry{
		{Command: constants.FlagArgStat, Description: constants.FlagDescStat},
		{Command: constants.FlagHelpEntryCachedStaged, Description: constants.FlagDescCachedStaged},
		{Command: constants.FlagArgNameOnly, Description: constants.FlagDescNameOnly},
		{Command: constants.FlagHelpEntryJSON, Description: constants.FlagDescJSONStructured},
		{Command: constants.FlagHelpEntryHelp, Description: constants.FlagDescHelp},
	}
}

func buildDiffOpsSection() termout.HelpSection {
	return termout.HelpSection{
		Title: constants.HelpSectionRepoDiff,
		Entries: []termout.CommandEntry{
			{Command: constants.HelpEntryDiffPath, Description: constants.HelpDescDiffPath},
			{Command: constants.HelpEntryDiffCached, Description: constants.HelpDescDiffCached},
			{Command: constants.HelpEntryDiffStat, Description: constants.HelpDescDiffStat},
		},
	}
}

func buildDiffProfilesSection() termout.HelpSection {
	return termout.HelpSection{
		Title: constants.HelpSectionProfiles,
		Entries: []termout.CommandEntry{
			{Command: constants.HelpEntryDiffProfiles, Description: constants.HelpDescDiffProfiles},
		},
	}
}
