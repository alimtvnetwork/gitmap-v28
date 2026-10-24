package cmdstatus

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdstatus/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

// RenderStatusHelp displays the styled two-column Status help menu.
func RenderStatusHelp() {
	termout.RenderMenu(buildStatusHelpMenu())
}

func buildStatusHelpMenu() termout.HelpMenu {
	return termout.HelpMenu{
		Title:       constants.StatusHelpTitle,
		UsageLines:  buildStatusUsageLines(),
		Sections:    buildStatusSections(),
		FooterFlags: buildStatusFooterFlags(),
		Tips:        buildStatusTips(),
	}
}

func buildStatusUsageLines() []string {
	return []string{
		constants.UsageStatus,
		constants.UsageStatusAlias,
		constants.UsageStatusDirty,
	}
}

func buildStatusTips() []string {
	return []string{
		constants.StatusHelpTipAudit,
		constants.StatusHelpTipJSON,
	}
}

func buildStatusSections() []termout.HelpSection {
	return []termout.HelpSection{
		buildStatusScopeSection(),
		buildStatusFormattingSection(),
	}
}

func buildStatusFooterFlags() []termout.CommandEntry {
	return []termout.CommandEntry{
		{Command: constants.FlagHelpEntryDirty, Description: constants.FlagDescFilterDirty},
		{Command: constants.FlagArgAhead, Description: constants.FlagDescFilterAhead},
		{Command: constants.FlagArgBehind, Description: constants.FlagDescFilterBehind},
		{Command: constants.FlagHelpEntryJSON, Description: constants.FlagDescJSON},
		{Command: constants.FlagArgTable, Description: constants.FlagDescTable},
		{Command: constants.FlagArgCompact, Description: constants.FlagDescCompact},
		{Command: constants.FlagHelpEntryHelp, Description: constants.FlagDescHelp},
	}
}

func buildStatusScopeSection() termout.HelpSection {
	return termout.HelpSection{
		Title: constants.StatusSectionFilterModes,
		Entries: []termout.CommandEntry{
			{Command: constants.StatusHelpEntryAll, Description: constants.StatusHelpDescAll},
			{Command: constants.FlagArgDirty, Description: constants.StatusHelpDescDirty},
			{Command: constants.StatusHelpEntryAheadBehind, Description: constants.StatusHelpDescAheadBehind},
		},
	}
}

func buildStatusFormattingSection() termout.HelpSection {
	return termout.HelpSection{
		Title: constants.StatusSectionPresentation,
		Entries: []termout.CommandEntry{
			{Command: constants.StatusHelpEntryTableMode, Description: constants.StatusHelpDescTableMode},
			{Command: constants.StatusHelpEntryJSONMode, Description: constants.StatusHelpDescJSONMode},
		},
	}
}
