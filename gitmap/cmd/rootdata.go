package cmd

import (
	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// dispatchData routes data management, history, profiles, and TUI commands.
func dispatchData(command string) bool {
	return runDispatchTable(command, dataDispatchEntries())
}

// dataDispatchEntries returns the routing table for data commands.
func dataDispatchEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{constants.CmdList, constants.CmdListAlias}, func() { runList(argsTail()) }},
		{[]string{constants.CmdGroup, constants.CmdGroupAlias}, func() { runGroup(argsTail()) }},
		{[]string{constants.CmdMultiGroup, constants.CmdMultiGroupAlias}, func() { runMultiGroup(argsTail()) }},
		{[]string{constants.CmdHistory, constants.CmdHistoryAlias}, func() { runHistory(argsTail()) }},
		{[]string{constants.CmdHistoryReset, constants.CmdHistoryResetAlias}, func() { runHistoryReset(argsTail()) }},
		{[]string{constants.CmdStats, constants.CmdStatsAlias}, func() { runStats(argsTail()) }},
		{[]string{constants.CmdBookmark, constants.CmdBookmarkAlias}, func() { runBookmark(argsTail()) }},
		{[]string{constants.CmdExport, constants.CmdExportAlias}, func() { runExport(argsTail()) }},
		{[]string{constants.CmdImport, constants.CmdImportAlias}, func() { runImport(argsTail()) }},
		{[]string{constants.CmdProfile, constants.CmdProfileAlias}, func() { runProfile(argsTail()) }},
		{[]string{constants.CmdDiffProfiles, constants.CmdDiffProfilesAlias}, func() { runDiffProfiles(argsTail()) }},
		{[]string{constants.CmdCD, constants.CmdCDAlias}, func() { runCD(argsTail()) }},
		{[]string{constants.CmdWatch, constants.CmdWatchAlias}, func() { runWatch(argsTail()) }},
		{[]string{constants.CmdInteractive, constants.CmdInteractiveAlias}, func() { runInteractive() }},
		{[]string{constants.CmdDBReset}, func() { runDBReset(argsTail()) }},
		{[]string{constants.CmdReset}, func() { runReset(argsTail()) }},
		{[]string{constants.CmdDBMigrate, constants.CmdDBMigrateAlias}, func() { runDBMigrate(argsTail()) }},
		{[]string{constants.CmdAmend, constants.CmdAmendAlias}, func() { runAmend(argsTail()) }},
		{[]string{constants.CmdAmendList, constants.CmdAmendListAlias}, func() { runAmendList(argsTail()) }},
		{[]string{constants.CmdDashboard, constants.CmdDashboardAlias}, func() { runDashboard(argsTail()) }},
		{[]string{constants.CmdVersionHistory, constants.CmdVersionHistoryAlias}, func() { runVersionHistory(argsTail()) }},
	}
}
