package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// dispatchTooling routes dev tooling and maintenance commands.
func dispatchTooling(command string) bool {
	return runDispatchTable(command, toolingDispatchEntries())
}

// toolingDispatchEntries returns the routing table for tooling commands.
func toolingDispatchEntries() []dispatchEntry {
	entries := make([]dispatchEntry, 0, 65)
	entries = append(entries, toolingWorkspaceEntries()...)
	entries = append(entries, toolingDevEntries()...)
	entries = append(entries, toolingAuditEntries()...)
	entries = append(entries, toolingOpsEntries()...)
	entries = append(entries, toolingInstallEntries()...)
	entries = append(entries, toolingUtilEntries()...)
	entries = append(entries, toolingChromeEntries()...)
	entries = append(entries, toolingNetworkEntries()...)
	return entries
}

func toolingWorkspaceEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{constants.CmdDesktopSync, constants.CmdDesktopSyncAlias}, func() { checkHelp("desktop-sync", argsTail()); runDesktopSync() }},
		{[]string{constants.CmdGitHubDesktop, constants.CmdGitHubDesktopAlias}, func() { runGitHubDesktop(argsTail()) }},
		{[]string{constants.CmdRescan, constants.CmdRescanAlias}, func() { checkHelp("rescan", argsTail()); runRescan() }},
		{[]string{constants.CmdRescanSubtree, constants.CmdRescanSubtreeAlias}, func() { runRescanSubtree(argsTail()) }},
		{[]string{constants.CmdSetup}, func() { runSetup(argsTail()) }},
		{[]string{constants.CmdDoctor}, func() { checkHelp("doctor", argsTail()); runDoctor(argsTail()) }},
		{[]string{constants.CmdLatestBranch, constants.CmdLatestBranchAlias}, func() { runLatestBranch(argsTail()) }},
		{[]string{constants.CmdBranch, constants.CmdBranchAlias}, func() { runBranch(argsTail()) }},
	}
}

func toolingDevEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{constants.CmdListVersions, constants.CmdListVersionsAlias}, func() { runListVersions(argsTail()) }},
		{[]string{constants.CmdListReleases, constants.CmdListReleasesAlias, constants.CmdReleases}, func() { runListReleases(argsTail()) }},
		{[]string{constants.CmdSEOWrite, constants.CmdSEOWriteAlias}, func() { runSEOWrite(argsTail()) }},
		{[]string{constants.CmdGoMod, constants.CmdGoModAlias}, func() { runGoMod(argsTail()) }},
		{[]string{constants.CmdCompletion, constants.CmdCompletionAlias}, func() { runCompletion(argsTail()) }},
		{[]string{constants.CmdZipGroup, constants.CmdZipGroupShort}, func() { runZipGroup(argsTail()) }},
		{[]string{constants.CmdAlias, constants.CmdAliasShort}, func() { runAlias(argsTail()) }},
		{[]string{constants.CmdSSH}, func() { runSSH(argsTail()) }},
		{[]string{constants.CmdBackup}, func() { runBackup(argsTail()) }},
	}
}

func toolingAuditEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{constants.CmdStale, constants.CmdStaleAlias}, func() { runStale(argsTail()) }},
		{[]string{constants.CmdOrphans}, func() { runOrphans(argsTail()) }},
		{[]string{constants.CmdDedupe}, func() { runDedupe(argsTail()) }},
		{[]string{constants.CmdSize}, func() { runSize(argsTail()) }},
		{[]string{constants.CmdReleaseNotes}, func() { runReleaseNotes(argsTail()) }},
		{[]string{constants.CmdReleaseDry}, func() { runReleaseDry(argsTail()) }},
		{[]string{constants.CmdTagRename}, func() { runTagRename(argsTail()) }},
		{[]string{constants.CmdRecent, constants.CmdRecentAlias}, func() { runRecent(argsTail()) }},
		{[]string{constants.CmdTodo}, func() { runTodo(argsTail()) }},
	}
}

func toolingOpsEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{constants.CmdOpen, constants.CmdOpenAlias}, func() { runOpen(argsTail()) }},
		{[]string{constants.CmdPR, constants.CmdPRAlias}, func() { runPR(argsTail()) }},
		{[]string{constants.CmdBlameStats}, func() { runBlameStats(argsTail()) }},
		{[]string{constants.CmdSnapshot}, func() { runSnapshot(argsTail()) }},
		{[]string{constants.CmdRollback}, func() { runRollback(argsTail()) }},
		{[]string{constants.CmdGuard}, func() { runGuard(argsTail()) }},
		{[]string{constants.CmdPrune, constants.CmdPruneAlias}, func() { runPrune(argsTail()) }},
		{[]string{constants.CmdTempRelease, constants.CmdTempReleaseShort}, func() { runTempRelease(argsTail()) }},
		{[]string{constants.CmdTask, constants.CmdTaskAlias}, func() { runTask(argsTail()) }},
		{[]string{constants.CmdEnv, constants.CmdEnvAlias}, func() { runEnv(argsTail()) }},
	}
}

func toolingInstallEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{"installer", "in"}, func() { RunInstallerCLI(argsTail()) }},
		{[]string{"cg", "coding-guide", "coding-guidelines", "ct"}, func() { runCG(argsTail()) }},
		{[]string{"install-version-json", "init-version"}, func() { runCG(append([]string{"install-version-json"}, argsTail()...)) }},
		{[]string{"install-prompts", "install-prompt"}, func() { runCG(append([]string{"install-prompts"}, argsTail()...)) }},
		{[]string{"prompts-status"}, func() { runCG(append([]string{"prompts-status"}, argsTail()...)) }},
		{[]string{"prompts-version"}, func() { runCG(append([]string{"prompts-version"}, argsTail()...)) }},
		{[]string{"workdir", "work-dir", "wd"}, func() { runWorkDir(argsTail()) }},
		{[]string{"os", "os-update"}, func() { RunOSCLI(argsTail()) }},
		{[]string{"sj", "ssh-joiner", "ssh-join", "ssh-joined"}, func() { runSJ(argsTail()) }},
		{[]string{"se", "ssh-exe", "ssh-exec", "ssh-execute"}, func() { runSSHExec(argsTail()) }},
		{[]string{constants.CmdInstall, constants.CmdInstallAlias}, func() { runInstall(argsTail()) }},
		{[]string{constants.CmdUninstall, constants.CmdUninstallAlias}, func() { runUninstall(argsTail()) }},
		{[]string{constants.CmdStartupAdd, constants.CmdStartupAddAlias}, func() { runStartupAdd(argsTail()) }},
		{[]string{constants.CmdStartupList, constants.CmdStartupListAlias}, func() { runStartupList(argsTail()) }},
		{[]string{constants.CmdStartupRemove, constants.CmdStartupRemoveAlias}, func() { runStartupRemove(argsTail()) }},
		{[]string{constants.CmdSelfInstall}, func() { runSelfInstall(argsTail()) }},
		{[]string{constants.CmdSelfUninstall}, func() { runSelfUninstall(argsTail()) }},
		{[]string{constants.CmdSelfUninstallRunner}, func() { runSelfUninstallRunner() }},
		{[]string{constants.CmdPending}, func() { runPending() }},
		{[]string{constants.CmdDoPending, constants.CmdDoPendingAlias}, func() { runDoPending(argsTail()) }},
	}
}

func toolingUtilEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{"schedule"}, func() { runSchedule(argsTail()) }},
		{[]string{constants.CmdDownloaderConfig, constants.CmdDownloaderConfigAlias}, func() { runDownloaderConfig(argsTail()) }},
		{[]string{constants.CmdUnzipCompact, constants.CmdUnzipCompactAlias}, func() { runUnzipCompact(argsTail()) }},
		{[]string{constants.CmdFolder}, func() { runFolder(argsTail()) }},
		{[]string{constants.CmdGitRm}, func() { runGitRm(argsTail()) }},
		{[]string{constants.CmdCommitPush, constants.CmdCommitPushAlias}, func() { runCommitPush(argsTail()) }},
		{[]string{constants.CmdCommitPushPull, constants.CmdCommitPushPullAlias}, func() { runCommitPushPull(argsTail()) }},
		{[]string{constants.CmdCommitPushBug, constants.CmdCommitPushBugAlias}, func() { runCommitPushBug(argsTail()) }},
		{[]string{constants.CmdCommitPushFeature, constants.CmdCommitPushFeatureAlias}, func() { runCommitPushFeature(argsTail()) }},
		{[]string{constants.CmdCommitPushRelease, constants.CmdCommitPushReleaseAlias}, func() { runCommitPushRelease(argsTail()) }},
		{[]string{constants.CmdRmGit, constants.CmdRmGitAlias}, func() { runRmGit(argsTail()) }},
		{[]string{constants.CmdIgnore}, func() { runIgnore(argsTail()) }},
		{[]string{constants.CmdIgnoreRm}, func() { runIgnoreRm(argsTail()) }},
		{[]string{constants.CmdAdd}, func() { runAdd(argsTail()) }},
		{[]string{constants.CmdAg, constants.CmdAntigravity}, func() { runAg(argsTail()) }},
				{[]string{constants.CmdLlm}, func() { runLlm(argsTail()) }},
		{[]string{constants.CmdFind}, func() { runFind(argsTail()) }},
		{[]string{constants.CmdFindRegex}, func() { runFindRegex(argsTail()) }},
		{[]string{constants.CmdFindRead}, func() { runFindRead(argsTail()) }},
		{[]string{constants.CmdFindReadJson}, func() { runFindReadJson(argsTail()) }},
		{[]string{constants.CmdFindRegexRead}, func() { runFindRegexRead(argsTail()) }},
		{[]string{constants.CmdFindRegexReadJson}, func() { runFindRegexReadJson(argsTail()) }},
		{[]string{constants.CmdFindHelp}, func() { runFindHelp(argsTail()) }},
		{[]string{constants.CmdSearchHelp}, func() { runSearchHelp(argsTail()) }},
		{[]string{constants.CmdRegexHelp}, func() { runRegexHelp(argsTail()) }},
		{[]string{constants.CmdSearch}, func() { runSearch(argsTail()) }},
		{[]string{constants.CmdReplace}, func() { runReplace(argsTail()) }},
		{[]string{constants.CmdReplaceRegex}, func() { runReplaceRegex(argsTail()) }},
		{[]string{constants.CmdRepoSearch, constants.CmdRepoSearchAlias}, func() { runRepoSearch(argsTail()) }},
		{[]string{"_index"}, func() { runIndex(argsTail()) }},
		{[]string{constants.CmdRepoRegex, constants.CmdRepoRegexAlias}, func() { runRepoRegex(argsTail()) }},
		{[]string{constants.CmdRepoSearchJson, constants.CmdRepoSearchJsonAlias}, func() { runRepoSearchJson(argsTail()) }},
		{[]string{constants.CmdRepoSearchRegexJson}, func() { runRepoSearchRegexJson(argsTail()) }},
		{[]string{constants.CmdSearchReplaceAll}, func() { runSearchReplaceAll(argsTail()) }},
		{[]string{constants.CmdZip}, func() { runZip(argsTail()) }},
		{[]string{"mkdir"}, func() { runMkdir(argsTail()) }},
		{[]string{"cat"}, func() { runCat(argsTail()) }},
		{[]string{constants.CmdAppend}, func() { runAppend(argsTail()) }},
		{[]string{constants.CmdWrite}, func() { runWrite(argsTail()) }},
		{[]string{"reset-and-rescan"}, func() { runReset([]string{"--confirm", "--rescan"}) }},
		{[]string{constants.CmdReplace, constants.CmdReplaceAlias}, func() { runReplace(argsTail()) }},
		{[]string{constants.CmdRegoldens, constants.CmdRegoldensAlias}, func() { runRegoldens(argsTail()) }},
		{[]string{constants.CmdAuditLegacy, constants.CmdAuditLegacyAlias, constants.CmdAuditLegacyAlias2}, func() { runAuditLegacy(argsTail()) }},
		{[]string{constants.CmdFixRepo, constants.CmdFixRepoAlias}, func() { runFixRepo(argsTail()) }},
		{[]string{constants.CmdUndo, constants.CmdUndoAlias}, func() { runUndo(argsTail()) }},
		{[]string{constants.CmdHistoryPurge, constants.CmdHistoryPurgeAlias}, func() { runHistoryPurge(argsTail()) }},
		{[]string{constants.CmdHistoryPin, constants.CmdHistoryPinAlias}, func() { runHistoryPin(argsTail()) }},
	}
}

func toolingChromeEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{constants.CmdChromeProfileCopy, constants.CmdChromeProfileCopyAlias}, func() { runChromeProfileCopy(argsTail()) }},
		{[]string{constants.CmdChromeProfileExport, constants.CmdChromeProfileExportAlias}, func() { runChromeProfileExport(argsTail()) }},
		{[]string{constants.CmdChromeProfileImport, constants.CmdChromeProfileImportAlias}, func() { runChromeProfileImport(argsTail()) }},
		{[]string{constants.CmdChromeProfileList, constants.CmdChromeProfileListAlias, constants.CmdChromeProfileListAlias2}, func() { runChromeProfileList(argsTail()) }},
		{[]string{constants.CmdChromeProfileDelete, constants.CmdChromeProfileDeleteAlias}, func() { runChromeProfileDelete(argsTail()) }},
		{[]string{constants.CmdChromeProfileMerge, constants.CmdChromeProfileMergeAlias}, func() { runChromeProfileMerge(argsTail()) }},
		{[]string{constants.CmdChrome}, func() { runChrome(argsTail()) }},
	}
}

func toolingNetworkEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{constants.CmdServe, constants.CmdServeAlias}, func() { runServe(argsTail()) }},
		{[]string{constants.CmdJoin, constants.CmdJoinAlias}, func() { runJoin(argsTail()) }},
	}
}
