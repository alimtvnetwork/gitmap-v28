package cmdbranch

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdbranch/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

// RenderBranchHelp displays the styled two-column Branch help menu.
func RenderBranchHelp() {
	termout.RenderMenu(buildBranchHelpMenu())
}

func buildBranchHelpMenu() termout.HelpMenu {
	return termout.HelpMenu{
		Title:       constants.HelpTitleBranch,
		UsageLines:  buildBranchUsageLines(),
		Sections:    []termout.HelpSection{buildBranchOpsSection()},
		FooterFlags: buildBranchFooterFlags(),
		Tips:        buildBranchTips(),
	}
}

func buildBranchUsageLines() []string {
	return []string{
		constants.HelpUsageDefault,
		constants.HelpUsageDefAlias,
		constants.HelpUsageLatestBranch,
	}
}

func buildBranchTips() []string {
	return []string{
		constants.HelpTipResolveDefault,
		constants.HelpTipOfflineFallback,
	}
}

func buildBranchFooterFlags() []termout.CommandEntry {
	return []termout.CommandEntry{
		{Command: constants.FlagHelpEntryJSON, Description: constants.FlagDescJSON},
		{Command: constants.FlagHelpEntryHelp, Description: constants.FlagDescHelp},
	}
}

func buildBranchOpsSection() termout.HelpSection {
	return termout.HelpSection{
		Title: constants.HelpSectionOperations,
		Entries: []termout.CommandEntry{
			{Command: constants.HelpEntryDefault, Description: constants.HelpDescDefault},
			{Command: constants.HelpEntryLatest, Description: constants.HelpDescLatest},
		},
	}
}
