package cmd

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdamend"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdbookmark"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcd"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmddashboard"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmddb"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmddiffprofiles"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdexport"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdfind"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdgroup"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdhistory"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdimport"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinteractive"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdlist"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdmacro"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdmail"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdmerge"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdmultigroup"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdprofiles"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdschedule"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdstartup"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdstats"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdstorage"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdsupabase"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdversionhistory"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdwatch"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
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
		{[]string{constants.CmdList, constants.CmdListAlias}, func() error { return cmdlist.RunList(argsTail()) }},
		{[]string{constants.CmdGroup, constants.CmdGroupAlias}, func() error { return cmdgroup.RunGroup(argsTail()) }},
		{[]string{constants.CmdMultiGroup, constants.CmdMultiGroupAlias}, func() error { return cmdmultigroup.RunMultiGroup(argsTail()) }},
		{[]string{constants.CmdHistory, constants.CmdHistoryAlias}, func() error { return cmdhistory.RunHistory(argsTail()) }},
		{[]string{constants.CmdHistoryReset, constants.CmdHistoryResetAlias}, func() error { return cmdhistory.RunHistoryReset(argsTail()) }},
		{[]string{constants.CmdStats, constants.CmdStatsAlias}, func() error { return cmdstats.RunStats(argsTail()) }},
		{[]string{constants.CmdBookmark, constants.CmdBookmarkAlias}, func() error { return cmdbookmark.RunBookmark(argsTail()) }},
	}
}

func dataProfileEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{constants.CmdExport, constants.CmdExportAlias}, func() error { return cmdexport.RunExport(argsTail()) }},
		{[]string{constants.CmdImport, constants.CmdImportAlias}, func() error { return cmdimport.RunImport(argsTail()) }},
		{[]string{"export-all"}, func() error { return cmdexport.RunExportAll(argsTail()) }},
		{[]string{"import-all"}, func() error { return cmdimport.RunImportAllJSONCLI(argsTail()) }},
		{[]string{"export-only"}, func() error { return cmdexport.RunExportOnly(argsTail()) }},
		{[]string{"import-export", "ie"}, func() error { return cmdexport.RunImportExport(argsTail()) }},
		{[]string{constants.CmdProfile, constants.CmdProfileAlias}, func() error { return cmdprofiles.RunProfile(argsTail()) }},
		{[]string{constants.CmdProfiles, constants.CmdProfilesAlias, "git-profiles"}, func() error { return cmdprofiles.RunProfiles(argsTail()) }},
		{[]string{constants.CmdDiffProfiles, constants.CmdDiffProfilesAlias}, func() error { return cmddiffprofiles.RunDiffProfiles(argsTail()) }},
		{[]string{constants.CmdCD, constants.CmdCDAlias}, func() error { return cmdcd.RunCD(argsTail()) }},
		{[]string{constants.CmdWatch, constants.CmdWatchAlias}, func() error { return cmdwatch.RunWatch(argsTail()) }},
		{[]string{constants.CmdInteractive, constants.CmdInteractiveAlias}, cmdinteractive.RunInteractive},
	}
}

func dataDatabaseEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{"db-read", "dbr"}, func() error {
			if appErr := cmddb.RunDbRead(argsTail()); appErr != nil {
				return appErr
			}
			return nil
		}},
		{[]string{"db-add", "dba", "dbadd"}, func() error {
			if appErr := cmddb.RunDbAdd(argsTail()); appErr != nil {
				return appErr
			}
			return nil
		}},
		{[]string{"db-edit", "dbe", "dbedit"}, func() error {
			if appErr := cmddb.RunDbEdit(argsTail()); appErr != nil {
				return appErr
			}
			return nil
		}},
		{[]string{"db-find", "dbf", "dbfind"}, func() error {
			if appErr := cmddb.RunDbFind(argsTail()); appErr != nil {
				return appErr
			}
			return nil
		}},
		{[]string{"db-listing", "dbl", "dbls", "db-list"}, func() error {
			if appErr := cmddb.RunDbListing(argsTail()); appErr != nil {
				return appErr
			}
			return nil
		}},
		{[]string{"db-view", "dbv", "dbview"}, func() error {
			if appErr := cmddb.RunDbView(argsTail()); appErr != nil {
				return appErr
			}
			return nil
		}},
		{[]string{"db-rest", "dbrest"}, func() error {
			if appErr := cmddb.RunDbRest(argsTail()); appErr != nil {
				return appErr
			}
			return nil
		}},
		{[]string{"db-record-del", "dbrd", "dbrecorddel"}, func() error {
			if appErr := cmddb.RunDbRecordDel(argsTail()); appErr != nil {
				return appErr
			}
			return nil
		}},
		{[]string{"db-rm", "dbrm", "db-delete", "dbdelete"}, func() error {
			if appErr := cmddb.RunDbRm(argsTail()); appErr != nil {
				return appErr
			}
			return nil
		}},
		{[]string{constants.CmdDBReset}, func() error { return cmddb.RunDBReset(argsTail()) }},
		{[]string{constants.CmdDB}, func() error { return cmddb.RunDB(argsTail()) }},
		{[]string{constants.CmdStartFresh}, func() error { return cmddb.RunStartFresh(argsTail()) }},
		{[]string{constants.CmdFindDuplicates, constants.CmdFindDuplicatesAlias, constants.CmdFindDuplicatesAlias2}, func() error { return cmdfind.RunFindDuplicates("", argsTail()) }},
		{[]string{constants.CmdReset}, func() error { return runReset(argsTail()) }},
		{[]string{constants.CmdDBMigrate, constants.CmdDBMigrateAlias}, func() error { return cmddb.RunDBMigrate(argsTail()) }},
		{[]string{constants.CmdAmend, constants.CmdAmendAlias}, func() error { return cmdamend.RunAmend(argsTail()) }},
		{[]string{constants.CmdAmendList, constants.CmdAmendListAlias}, func() error { return cmdamend.RunAmendList(argsTail()) }},
		{[]string{constants.CmdDashboard, constants.CmdDashboardAlias}, func() error { return cmddashboard.RunDashboard(argsTail()) }},
		{[]string{constants.CmdVersionHistory, constants.CmdVersionHistoryAlias}, func() error { return cmdversionhistory.RunVersionHistory(argsTail()) }},
		{[]string{"supabase", "sb"}, func() error { return cmdsupabase.RunSupabaseCLI(argsTail()) }},
	}
}

func dataExecutionEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{"execute", "exec"}, func() error { return cmdmacro.RunExecuteCmd(argsTail()) }},
		{[]string{"macro", "m"}, func() error { return cmdmacro.RunMacroCmd(argsTail()) }},
		{[]string{"macro-run", "macro-exec"}, func() error { return cmdmacro.RunExecuteCmd(argsTail()) }},
		{[]string{"macro-add", "macro-create"}, func() error { return cmdmacro.HandleMacroAdd(argsTail()) }},
		{[]string{"macro-edit", "macro-modify"}, func() error { return cmdmacro.HandleMacroEdit(argsTail()) }},
		{[]string{"macro-list", "macro-ls"}, func() error { return cmdmacro.HandleMacroList(argsTail()) }},
		{[]string{"macro-record", "macro-rec"}, func() error { return cmdmacro.HandleMacroRecord(argsTail()) }},
		{[]string{"macro-show"}, func() error { return cmdmacro.HandleMacroShow(argsTail()) }},
		{[]string{"macro-rm", "macro-del"}, func() error { return cmdmacro.HandleMacroDelete(argsTail()) }},
		{[]string{"macro-export", "macro-exp"}, func() error { return cmdmacro.RunMacroExport(argsTail()) }},
		{[]string{"macro-export-all"}, func() error { return cmdmacro.RunMacroExport(append([]string{"--all"}, argsTail()...)) }},
		{[]string{"macro-export-single"}, func() error { return cmdmacro.RunMacroExport(append([]string{"--single"}, argsTail()...)) }},
		{[]string{"macro-import", "macro-imp"}, func() error { return cmdmacro.RunMacroImport(argsTail()) }},
		{[]string{"macro-import-all"}, func() error { return cmdmacro.RunMacroImport(append([]string{"--all"}, argsTail()...)) }},
		{[]string{"macro-import-single"}, func() error { return cmdmacro.RunMacroImport(append([]string{"--single"}, argsTail()...)) }},
		{[]string{"record", "rec"}, func() error { return cmdmacro.RunMacroCmd(append([]string{"record"}, argsTail()...)) }},
		{[]string{"retry", "loop", "until-success"}, func() error { return cmdmacro.RunMacroUntilSuccess(argsTail()) }},
		{[]string{"mv", "move"}, func() error { return runMoveDisambiguated(argsTail()) }},
		{[]string{constants.CmdStartup, constants.CmdStartupAlias}, func() error { return cmdstartup.RunStartupCmd(argsTail()) }},
		{[]string{constants.CmdAsync, constants.CmdAsyncAlias}, func() error { return RunAsyncCmd(argsTail()) }},
		{[]string{constants.CmdStorage, constants.CmdStorageAlias}, func() error { return cmdstorage.RunStorageCmd(argsTail()) }},
		{[]string{"schedule", "sc", "crontab", "cron"}, func() error { return cmdschedule.RunSchedule(argsTail()) }},
		{[]string{"mail", "mailbox"}, func() error {
			if appErr := cmdmail.RunMail(argsTail()); appErr != nil {
				return appErr
			}
			return nil
		}},
	}
}

func runMoveDisambiguated(args []string) error {
	if isRemoteMoveArgs(args) {
		return cmdssh.RunSSHMvExtended(args)
	}

	return cmdmerge.RunMove(args)
}

func isRemoteMoveArgs(args []string) bool {
	for _, arg := range args {
		if isRemoteMoveFlag(arg) || hasRemoteMoveColon(arg) {
			return true
		}
	}

	return false
}

func isRemoteMoveFlag(arg string) bool {
	if arg == "-t" || arg == "--target" || arg == "--all" {
		return true
	}

	return strings.HasPrefix(arg, "-t=") || strings.HasPrefix(arg, "--target=")
}

func hasRemoteMoveColon(arg string) bool {
	if isDrivePath(arg) {
		return false
	}

	return strings.Contains(arg, ":")
}

func isDrivePath(arg string) bool {
	if len(arg) < 2 {
		return false
	}

	return arg[1] == ':' && isLetterByte(arg[0])
}

func isLetterByte(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}
