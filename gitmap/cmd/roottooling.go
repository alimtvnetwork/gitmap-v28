package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// dispatchTooling routes dev tooling and maintenance commands.
func dispatchTooling(command string) (bool, error) {
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
		{[]string{constants.CmdDesktopSync, constants.CmdDesktopSyncAlias}, func() error { checkHelp("desktop-sync", argsTail()); return runDesktopSync() }},
		{[]string{constants.CmdGitHubDesktop, constants.CmdGitHubDesktopAlias}, func() error { return runGitHubDesktop(argsTail()) }},
		{[]string{constants.CmdRescan, constants.CmdRescanAlias}, func() error { checkHelp("rescan", argsTail()); return runRescan() }},
		{[]string{constants.CmdRescanSubtree, constants.CmdRescanSubtreeAlias}, func() error { return runRescanSubtree(argsTail()) }},
		{[]string{constants.CmdSetup}, func() error { return runSetup(argsTail()) }},
		{[]string{constants.CmdDoctor}, func() error { checkHelp("doctor", argsTail()); return runDoctor(argsTail()) }},
		{[]string{constants.CmdLatestBranch, constants.CmdLatestBranchAlias}, func() error { return runLatestBranch(argsTail()) }},
		{[]string{constants.CmdBranch, constants.CmdBranchAlias}, func() error { return runBranch(argsTail()) }},
	}
}

func toolingDevEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{constants.CmdListVersions, constants.CmdListVersionsAlias}, func() error { return runListVersions(argsTail()) }},
		{[]string{constants.CmdListReleases, constants.CmdListReleasesAlias, constants.CmdReleases}, func() error { return runListReleases(argsTail()) }},
		{[]string{constants.CmdSEOWrite, constants.CmdSEOWriteAlias}, func() error { return runSEOWrite(argsTail()) }},
		{[]string{constants.CmdGoMod, constants.CmdGoModAlias}, func() error { return runGoMod(argsTail()) }},
		{[]string{constants.CmdCompletion, constants.CmdCompletionAlias}, func() error { return runCompletion(argsTail()) }},
		{[]string{constants.CmdZipGroup, constants.CmdZipGroupShort}, func() error { return runZipGroup(argsTail()) }},
		{[]string{constants.CmdAlias, constants.CmdAliasShort}, func() error { return runAlias(argsTail()) }},
		{[]string{constants.CmdSSH}, func() error { return runSSH(argsTail()) }},
		{[]string{constants.CmdBackup}, func() error { return runBackup(argsTail()) }},
	}
}

func toolingAuditEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{constants.CmdStale, constants.CmdStaleAlias}, func() error { return runStale(argsTail()) }},
		{[]string{constants.CmdOrphans}, func() error { return runOrphans(argsTail()) }},
		{[]string{constants.CmdDedupe}, func() error { return runDedupe(argsTail()) }},
		{[]string{constants.CmdSize}, func() error { return runSize(argsTail()) }},
		{[]string{constants.CmdReleaseNotes}, func() error { return runReleaseNotes(argsTail()) }},
		{[]string{constants.CmdReleaseDry}, func() error { return runReleaseDry(argsTail()) }},
		{[]string{constants.CmdTagRename}, func() error { return runTagRename(argsTail()) }},
		{[]string{constants.CmdRecent, constants.CmdRecentAlias}, func() error { return runRecent(argsTail()) }},
		{[]string{constants.CmdTodo}, func() error { return runTodo(argsTail()) }},
	}
}

func toolingOpsEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{constants.CmdOpen, constants.CmdOpenAlias}, func() error { return runOpen(argsTail()) }},
		{[]string{constants.CmdPR, constants.CmdPRAlias}, func() error { return runPR(argsTail()) }},
		{[]string{constants.CmdBlameStats}, func() error { return runBlameStats(argsTail()) }},
		{[]string{constants.CmdSnapshot}, func() error { return runSnapshot(argsTail()) }},
		{[]string{constants.CmdRollback}, func() error { return runRollback(argsTail()) }},
		{[]string{constants.CmdGuard}, func() error { return runGuard(argsTail()) }},
		{[]string{constants.CmdPrune, constants.CmdPruneAlias}, func() error { return runPrune(argsTail()) }},
		{[]string{constants.CmdTempRelease, constants.CmdTempReleaseShort}, func() error { return runTempRelease(argsTail()) }},
		{[]string{constants.CmdTask, constants.CmdTaskAlias}, func() error { return runTask(argsTail()) }},
		{[]string{constants.CmdEnv, constants.CmdEnvAlias}, func() error { return runEnv(argsTail()) }},
	}
}

func toolingInstallEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{"installer", "in"}, func() error { RunInstallerCLI(argsTail()); return nil }},
		{[]string{"cg", "coding-guide", "coding-guidelines", "ct"}, func() error { return runCG(argsTail()) }},
		{[]string{"install-version-json", "init-version"}, func() error { return runCG(append([]string{"install-version-json"}, argsTail()...)) }},
		{[]string{"install-prompts", "install-prompt"}, func() error { return runCG(append([]string{"install-prompts"}, argsTail()...)) }},
		{[]string{"prompts-status"}, func() error { return runCG(append([]string{"prompts-status"}, argsTail()...)) }},
		{[]string{"prompts-version"}, func() error { return runCG(append([]string{"prompts-version"}, argsTail()...)) }},
		{[]string{"workdir", "work-dir", "wd"}, func() error { return runWorkDir(argsTail()) }},
		{[]string{"os", "os-update"}, func() error { RunOSCLI(argsTail()); return nil }},
		{[]string{"sj", "ssh-joiner", "ssh-join", "ssh-joined"}, func() error { return runSJ(argsTail()) }},
		{[]string{"se", "ssh-exe", "ssh-exec", "ssh-execute"}, func() error { return runSSHExec(argsTail()) }},
		{[]string{constants.CmdInstall, constants.CmdInstallAlias}, func() error { return runInstall(argsTail()) }},
		{[]string{constants.CmdUninstall, constants.CmdUninstallAlias}, func() error { return runUninstall(argsTail()) }},
		{[]string{constants.CmdStartupAdd, constants.CmdStartupAddAlias}, func() error { return runStartupAdd(argsTail()) }},
		{[]string{constants.CmdStartupList, constants.CmdStartupListAlias}, func() error { return runStartupList(argsTail()) }},
		{[]string{constants.CmdStartupRemove, constants.CmdStartupRemoveAlias}, func() error { return runStartupRemove(argsTail()) }},
		{[]string{constants.CmdSelfInstall}, func() error { return runSelfInstall(argsTail()) }},
		{[]string{constants.CmdSelfUninstall}, func() error { return runSelfUninstall(argsTail()) }},
		{[]string{constants.CmdSelfUninstallRunner}, runSelfUninstallRunner},
		{[]string{constants.CmdPending}, runPending},
		{[]string{constants.CmdDoPending, constants.CmdDoPendingAlias}, func() error { return runDoPending(argsTail()) }},
	}
}

func toolingUtilEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{"schedule", "sc"}, func() error { return runSchedule(argsTail()) }},
		{[]string{constants.CmdDownloaderConfig, constants.CmdDownloaderConfigAlias}, func() error { return runDownloaderConfig(argsTail()) }},
		{[]string{constants.CmdUnzipCompact, constants.CmdUnzipCompactAlias}, func() error { return runUnzipCompact(argsTail()) }},
		{[]string{constants.CmdFolder, "tree"}, func() error { return runFolder(argsTail()) }},
		{[]string{constants.CmdLowercase}, func() error { return runLowercase(argsTail()) }},
		{[]string{constants.CmdFixSeqFiles, constants.CmdFixSeqFilesAlias}, func() error { return runFixSeqFiles(argsTail()) }},
		{[]string{constants.CmdSequence, constants.CmdSequenceAlias}, func() error { return runSequence(argsTail()) }},
		{[]string{constants.CmdGitRm}, func() error { return runGitRm(argsTail()) }},
		{[]string{constants.CmdCommitPush, constants.CmdCommitPushAlias}, func() error { return runCommitPush(argsTail()) }},
		{[]string{constants.CmdPullCommitPush, constants.CmdPullCommitPushAlias}, func() error { return runPullCommitPush(argsTail()) }},
		{[]string{constants.CmdCommitPushBug, constants.CmdCommitPushBugAlias}, func() error { return runCommitPushBug(argsTail()) }},
		{[]string{constants.CmdCommitPushFeature, constants.CmdCommitPushFeatureAlias}, func() error { return runCommitPushFeature(argsTail()) }},
		{[]string{constants.CmdCommitPushRelease, constants.CmdCommitPushReleaseAlias}, func() error { return runCommitPushRelease(argsTail()) }},
		{[]string{constants.CmdRmGit, constants.CmdRmGitAlias}, func() error { return runRmGit(argsTail()) }},
		{[]string{constants.CmdGitReset, constants.CmdGitResetAlias}, func() error { return runGitReset(argsTail()) }},
		{[]string{constants.CmdIgnore}, func() error { return runIgnore(argsTail()) }},
		{[]string{constants.CmdIgnoreRm}, func() error { return runIgnoreRm(argsTail()) }},
		{[]string{constants.CmdAdd}, func() error { return runAdd(argsTail()) }},
		{[]string{constants.CmdLlm}, func() error { return runLlm(argsTail()) }},
		{[]string{constants.CmdFind, "f"}, func() error { return runFind(argsTail()) }},
		{[]string{constants.CmdFindFiles, "ff"}, func() error { return runFindFiles(argsTail()) }},
		{[]string{constants.CmdFindFilesAny, "ffa"}, func() error { return runFindFilesAny(argsTail()) }},
		{[]string{constants.CmdFindFilesStartsWith, "ffs"}, func() error { return runFindFilesStartsWith(argsTail()) }},
		{[]string{constants.CmdFindFilesEndsWith, "ffe"}, func() error { return runFindFilesEndsWith(argsTail()) }},
		{[]string{constants.CmdListFiles, "lf"}, func() error { return runListFiles(argsTail()) }},
		{[]string{constants.CmdFindRegex}, func() error { return runFindRegex(argsTail()) }},
		{[]string{constants.CmdFindRead}, func() error { return runFindRead(argsTail()) }},
		{[]string{constants.CmdFindReadJson}, func() error { return runFindReadJson(argsTail()) }},
		{[]string{constants.CmdFindRegexRead}, func() error { return runFindRegexRead(argsTail()) }},
		{[]string{constants.CmdFindRegexReadJson}, func() error { return runFindRegexReadJson(argsTail()) }},
		{[]string{constants.CmdFindHelp}, func() error { return runFindHelp(argsTail()) }},
		{[]string{constants.CmdSearchHelp}, func() error { return runSearchHelp(argsTail()) }},
		{[]string{constants.CmdRegexHelp}, func() error { return runRegexHelp(argsTail()) }},
		{[]string{constants.CmdSearch}, func() error { return runSearch(argsTail()) }},
		{[]string{constants.CmdReplace}, func() error { return runReplace(argsTail()) }},
		{[]string{constants.CmdReplaceRegex}, func() error { return runReplaceRegex(argsTail()) }},
		{[]string{constants.CmdRepoSearch, constants.CmdRepoSearchAlias}, func() error { return runRepoSearch(argsTail()) }},
		{[]string{"_index"}, func() error { return runIndex(argsTail()) }},
		{[]string{constants.CmdRepoRegex, constants.CmdRepoRegexAlias}, func() error { return runRepoRegex(argsTail()) }},
		{[]string{constants.CmdRepoSearchJson, constants.CmdRepoSearchJsonAlias}, func() error { return runRepoSearchJson(argsTail()) }},
		{[]string{constants.CmdRepoSearchRegexJson}, func() error { return runRepoSearchRegexJson(argsTail()) }},
		{[]string{constants.CmdSearchReplaceAll}, func() error { return runSearchReplaceAll(argsTail()) }},
		{[]string{constants.CmdZip}, func() error { return runZip(argsTail()) }},
		{[]string{"mkdir"}, func() error { return runMkdir(argsTail()) }},
		{[]string{"cat"}, func() error { return runCat(argsTail()) }},
		{[]string{constants.CmdAppend}, func() error { return runAppend(argsTail()) }},
		{[]string{constants.CmdWrite}, func() error { return runWrite(argsTail()) }},
		{[]string{constants.CmdHead}, func() error { return runHead(argsTail()) }},
		{[]string{constants.CmdTail}, func() error { return runTail(argsTail()) }},
		{[]string{constants.CmdFileSearch}, func() error { return runFileSearch(argsTail()) }},
		{[]string{"reset-and-rescan"}, func() error { return runReset([]string{"--confirm", "--rescan"}) }},
		{[]string{constants.CmdReplace, constants.CmdReplaceAlias}, func() error { return runReplace(argsTail()) }},
		{[]string{constants.CmdRegoldens, constants.CmdRegoldensAlias}, func() error { return runRegoldens(argsTail()) }},
		{[]string{constants.CmdAuditLegacy, constants.CmdAuditLegacyAlias, constants.CmdAuditLegacyAlias2}, func() error { return runAuditLegacy(argsTail()) }},
		{[]string{constants.CmdFixRepo, constants.CmdFixRepoAlias}, func() error { return runFixRepo(argsTail()) }},
		{[]string{constants.CmdUndo, constants.CmdUndoAlias}, func() error { return runUndo(argsTail()) }},
		{[]string{constants.CmdHistoryPurge, constants.CmdHistoryPurgeAlias}, func() error { return runHistoryPurge(argsTail()) }},
		{[]string{constants.CmdHistoryPin, constants.CmdHistoryPinAlias}, func() error { return runHistoryPin(argsTail()) }},
	}
}

func toolingChromeEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{constants.CmdChromeProfileCopy, constants.CmdChromeProfileCopyAlias}, func() error { return runChromeProfileCopy(argsTail()) }},
		{[]string{constants.CmdChromeProfileExport, constants.CmdChromeProfileExportAlias}, func() error { return runChromeProfileExport(argsTail()) }},
		{[]string{constants.CmdChromeProfileImport, constants.CmdChromeProfileImportAlias}, func() error { return runChromeProfileImport(argsTail()) }},
		{[]string{constants.CmdChromeProfileList, constants.CmdChromeProfileListAlias, constants.CmdChromeProfileListAlias2}, func() error { return runChromeProfileList(argsTail()) }},
		{[]string{constants.CmdChromeProfileDelete, constants.CmdChromeProfileDeleteAlias}, func() error { return runChromeProfileDelete(argsTail()) }},
		{[]string{constants.CmdChromeProfileMerge, constants.CmdChromeProfileMergeAlias}, func() error { return runChromeProfileMerge(argsTail()) }},
		{[]string{constants.CmdChrome}, func() error { return runChrome(argsTail()) }},
	}
}

func toolingNetworkEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{constants.CmdServe, constants.CmdServeAlias}, func() error { return runServe(argsTail()) }},
		{[]string{constants.CmdJoin, constants.CmdJoinAlias}, func() error { return runJoin(argsTail()) }},
	}
}
