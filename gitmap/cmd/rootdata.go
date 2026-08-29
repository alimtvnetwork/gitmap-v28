package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// dispatchData routes data management, history, profiles, and TUI commands.
func dispatchData(command string) (bool, error) {
	return runDispatchTable(command, dataDispatchEntries())
}

// dataDispatchEntries returns the routing table for data commands.
func dataDispatchEntries() []dispatchEntry {
	entries := make([]dispatchEntry, 0, 25)
	entries = append(entries, dataListingEntries()...)
	entries = append(entries, dataProfileEntries()...)
	entries = append(entries, dataDatabaseEntries()...)
	entries = append(entries, dataExecutionEntries()...)
	return entries
}

func dataListingEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{constants.CmdList, constants.CmdListAlias}, func() error { return runList(argsTail()) }},
		{[]string{constants.CmdGroup, constants.CmdGroupAlias}, func() error { return runGroup(argsTail()) }},
		{[]string{constants.CmdMultiGroup, constants.CmdMultiGroupAlias}, func() error { return runMultiGroup(argsTail()) }},
		{[]string{constants.CmdHistory, constants.CmdHistoryAlias}, func() error { return runHistory(argsTail()) }},
		{[]string{constants.CmdHistoryReset, constants.CmdHistoryResetAlias}, func() error { return runHistoryReset(argsTail()) }},
		{[]string{constants.CmdStats, constants.CmdStatsAlias}, func() error { return runStats(argsTail()) }},
		{[]string{constants.CmdBookmark, constants.CmdBookmarkAlias}, func() error { return runBookmark(argsTail()) }},
	}
}

func dataProfileEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{constants.CmdExport, constants.CmdExportAlias}, func() error { return runExport(argsTail()) }},
		{[]string{constants.CmdImport, constants.CmdImportAlias}, func() error { return runImport(argsTail()) }},
		{[]string{"import-export", "ie"}, func() error { return runImportExport(argsTail()) }},
		{[]string{constants.CmdProfile, constants.CmdProfileAlias}, func() error { return runProfile(argsTail()) }},
		{[]string{constants.CmdDiffProfiles, constants.CmdDiffProfilesAlias}, func() error { return runDiffProfiles(argsTail()) }},
		{[]string{constants.CmdCD, constants.CmdCDAlias}, func() error { return runCD(argsTail()) }},
		{[]string{constants.CmdWatch, constants.CmdWatchAlias}, func() error { return runWatch(argsTail()) }},
		{[]string{constants.CmdInteractive, constants.CmdInteractiveAlias}, runInteractive},
	}
}

func dataDatabaseEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{constants.CmdDBReset}, func() error { return runDBReset(argsTail()) }},
		{[]string{constants.CmdReset}, func() error { return runReset(argsTail()) }},
		{[]string{constants.CmdDBMigrate, constants.CmdDBMigrateAlias}, func() error { return runDBMigrate(argsTail()) }},
		{[]string{constants.CmdAmend, constants.CmdAmendAlias}, func() error { return runAmend(argsTail()) }},
		{[]string{constants.CmdAmendList, constants.CmdAmendListAlias}, func() error { return runAmendList(argsTail()) }},
		{[]string{constants.CmdDashboard, constants.CmdDashboardAlias}, func() error { return runDashboard(argsTail()) }},
		{[]string{constants.CmdVersionHistory, constants.CmdVersionHistoryAlias}, func() error { return runVersionHistory(argsTail()) }},
	}
}

func dataExecutionEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{"execute", "exec"}, func() error { return runExecuteCmd(argsTail()) }},
		{[]string{"macro", "m"}, func() error { return runMacroCmd(argsTail()) }},
		{[]string{"record", "rec"}, func() error { return runMacroCmd(append([]string{"record"}, argsTail()...)) }},
		{[]string{"mv", "move"}, func() error { return runMove(argsTail()) }},
	}
}
