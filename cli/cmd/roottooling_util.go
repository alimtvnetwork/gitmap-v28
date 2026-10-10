package cmd

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdadd"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdasset"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdauditlegacy"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdauthor"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdchrome"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdchromeprofile"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcommit"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdconfig"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmddownload"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdfind"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdfixgit"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdfixrepo"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdfoldertree"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdgitrm"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdhistory"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdignore"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdimport"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdindex"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdjoin"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdllm"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdmacro"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdmerge"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdports"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpull"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdregoldens"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdreplace"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdschedule"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdsearch"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdsequence"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdservercmd"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdundo"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdzip"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func toolingUtilEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{"schedule", "sc", "crontab", "cron"}, func() error { return cmdschedule.RunSchedule(argsTail()) }},
		{[]string{constants.CmdDownloaderConfig, constants.CmdDownloaderConfigAlias}, func() error { return cmddownload.RunDownloaderConfig(argsTail()) }},
		{[]string{constants.CmdUnzipCompact, constants.CmdUnzipCompactAlias}, func() error { return cmdzip.RunUnzipCompact(argsTail()) }},
		{[]string{constants.CmdFolder, "tree"}, func() error { return cmdfoldertree.RunFolder(argsTail()) }},
		{[]string{"tree-search", "treesearch", "ts"}, func() error { return cmdfoldertree.RunTreeSearchDispatch("tree-search", argsTail()) }},
		{[]string{"tree-search-startsWith", "tree-search-startswith", "tree-search-starts-with", "tree-search-prefix", "tss"}, func() error { return cmdfoldertree.RunTreeSearchDispatch("tree-search-startsWith", argsTail()) }},
		{[]string{"tree-search-contains", "tree-search-contain", "tree-search-substr", "tsc"}, func() error { return cmdfoldertree.RunTreeSearchDispatch("tree-search-contains", argsTail()) }},
		{[]string{"tree-search-endsWith", "tree-search-endswith", "tree-search-ends-with", "tree-search-suffix", "tse"}, func() error { return cmdfoldertree.RunTreeSearchDispatch("tree-search-endsWith", argsTail()) }},
		{[]string{"tree-search-grep", "tree-search-regex", "tsg"}, func() error { return cmdfoldertree.RunTreeSearchDispatch("tree-search-grep", argsTail()) }},
		{[]string{"tree-learn", "treelearn", "tl"}, func() error { return cmdfoldertree.RunTreeLearn(argsTail()) }},
		{[]string{
			constants.CmdLowercase,
			constants.CmdLowerCaseFix,
			constants.CmdLowerCaseFixAlias,
			constants.CmdLowerCaseFixShort,
			constants.CmdLowercaseReadme,
			constants.CmdReadmeLower,
			"lower",
			"lower-case-readme",
			"readme-lowercase",
			"lcr",
			"lc-fix",
		}, func() error {
			checkHelp("lowercase", argsTail())
			return runLowerCaseFixCLI(argsTail())
		}},
		{[]string{constants.CmdFixSeqFiles, constants.CmdFixSeqFilesAlias}, func() error { return runFixSeqFiles(argsTail()) }},
		{[]string{constants.CmdSequence, constants.CmdSequenceAlias}, func() error { return cmdsequence.RunSequence(argsTail()) }},
		{[]string{constants.CmdGitRm}, func() error { return cmdgitrm.RunGitRm(argsTail()) }},
		{[]string{constants.CmdCommitPush, constants.CmdCommitPushAlias}, func() error { return cmdcommit.RunCommitPush(argsTail()) }},
		{[]string{constants.CmdPullCommitPush, constants.CmdPullCommitPushAlias}, func() error { return cmdcommit.RunPullCommitPush(argsTail()) }},
		{[]string{constants.CmdCommitPushBug, constants.CmdCommitPushBugAlias}, func() error { return cmdcommit.RunCommitPushBug(argsTail()) }},
		{[]string{constants.CmdCommitPushFeature, constants.CmdCommitPushFeatureAlias}, func() error { return cmdcommit.RunCommitPushFeature(argsTail()) }},
		{[]string{constants.CmdCommitPushChore}, func() error { return cmdcommit.RunCommitPushChore(argsTail()) }},
		{[]string{"cpc"}, func() error { return runCpcDisambiguated(argsTail()) }},
		{[]string{constants.CmdCommitPushRelease, constants.CmdCommitPushReleaseAlias}, func() error { return cmdcommit.RunCommitPushRelease(argsTail()) }},
		{[]string{constants.CmdRmGit, constants.CmdRmGitAlias}, func() error { return cmdcommit.RunRmGit(argsTail()) }},
		{[]string{constants.CmdGitReset, constants.CmdGitResetAlias}, func() error { return cmdcommit.RunGitReset(argsTail()) }},
		{[]string{constants.CmdIgnore}, func() error { return cmdignore.RunIgnoreCLI(argsTail()) }},
		{[]string{constants.CmdIgnoreRm}, func() error { return cmdignore.RunIgnoreRm(argsTail()) }},
		{[]string{constants.CmdAdd}, func() error { return cmdadd.RunAdd(argsTail()) }},
		{[]string{constants.CmdLlm}, func() error { return cmdllm.RunLlm(argsTail()) }},
		{[]string{"locate", "find-tool", "vcvars"}, func() error { return RunLocateTopLevel(argsTail()) }},
		{[]string{constants.CmdFind, "f"}, func() error { return cmdfind.RunFind(argsTail()) }},
		{[]string{constants.CmdFindFiles, "ff"}, func() error { return cmdfind.RunFindFiles(argsTail()) }},
		{[]string{constants.CmdFindFilesAny, "ffa"}, func() error { return cmdfind.RunFindFilesAny(argsTail()) }},
		{[]string{constants.CmdFindFilesStartsWith, "ffs"}, func() error { return cmdfind.RunFindFilesStartsWith(argsTail()) }},
		{[]string{constants.CmdFindFilesEndsWith, "ffe"}, func() error { return cmdfind.RunFindFilesEndsWith(argsTail()) }},
		{[]string{constants.CmdListFiles, "lf"}, func() error { return cmdfind.RunListFiles(argsTail()) }},
		{[]string{constants.CmdFindRegex}, func() error { return cmdfind.RunFindRegex(argsTail()) }},
		{[]string{constants.CmdFindRead}, func() error { return cmdfind.RunFindRead(argsTail()) }},
		{[]string{constants.CmdFindReadJson}, func() error { return cmdfind.RunFindReadJson(argsTail()) }},
		{[]string{constants.CmdFindRegexRead}, func() error { return cmdfind.RunFindRegexRead(argsTail()) }},
		{[]string{constants.CmdFindRegexReadJson}, func() error { return cmdfind.RunFindRegexReadJson(argsTail()) }},
		{[]string{constants.CmdFindHelp}, func() error { cmdfind.RenderFindHelp(); return nil }},
		{[]string{constants.CmdSearchHelp}, func() error { cmdsearch.RenderSearchHelp(); return nil }},
		{[]string{constants.CmdRegexHelp}, func() error { cmdsearch.RenderSearchHelp(); return nil }},
		{[]string{constants.CmdSearch}, func() error { return cmdsearch.RunSearch(argsTail()) }},
		{[]string{constants.CmdReplace}, func() error { return cmdreplace.RunReplace(argsTail()) }},
		{[]string{constants.CmdReplaceRegex}, func() error { return cmdsearch.RunReplaceRegex(argsTail()) }},
		{[]string{constants.CmdRepoSearch, constants.CmdRepoSearchAlias}, func() error { return cmdsearch.RunRepoSearch(argsTail()) }},
		{[]string{"_index"}, func() error { return cmdindex.RunIndex(argsTail()) }},
		{[]string{constants.CmdRepoRegex, constants.CmdRepoRegexAlias}, func() error { return cmdsearch.RunRepoRegex(argsTail()) }},
		{[]string{constants.CmdRepoSearchJson, constants.CmdRepoSearchJsonAlias}, func() error { return cmdsearch.RunRepoSearchJson(argsTail()) }},
		{[]string{constants.CmdRepoSearchRegexJson}, func() error { return cmdsearch.RunRepoSearchRegexJson(argsTail()) }},
		{[]string{constants.CmdSearchReplaceAll}, func() error { return cmdsearch.RunSearchReplaceAll(argsTail()) }},
		{[]string{"asset", "download-prnt", "prnt", "prnt-download"}, func() error { return cmdasset.Run(argsTail()) }},
		{[]string{constants.CmdZip}, func() error { return cmdzip.RunZip(argsTail()) }},
		{[]string{"mkdir"}, func() error { return RunMkdir(argsTail()) }},
		{[]string{"cat"}, func() error { return cmdmacro.RunCat(argsTail()) }},
		{[]string{constants.CmdAppend}, func() error { return RunAppend(argsTail()) }},
		{[]string{constants.CmdWrite}, func() error { return RunWrite(argsTail()) }},
		{[]string{constants.CmdHead}, func() error { return RunHead(argsTail()) }},
		{[]string{constants.CmdTail}, func() error { return RunTail(argsTail()) }},
		{[]string{constants.CmdFileSearch}, func() error { return runFileSearch(argsTail()) }},
		{[]string{"reset-and-rescan"}, func() error { return runReset([]string{"--confirm", "--rescan"}) }},
		{[]string{constants.CmdReplace, constants.CmdReplaceAlias}, func() error { return cmdreplace.RunReplace(argsTail()) }},
		{[]string{constants.CmdRegoldens, constants.CmdRegoldensAlias}, func() error { return cmdregoldens.RunRegoldens(argsTail()) }},
		{[]string{constants.CmdAuditLegacy, constants.CmdAuditLegacyAlias, constants.CmdAuditLegacyAlias2}, func() error { return cmdauditlegacy.RunAuditLegacy(argsTail()) }},
		{[]string{constants.CmdFixRepo, constants.CmdFixRepoAlias}, func() error { return cmdfixrepo.RunFixRepo(argsTail()) }},
		{[]string{constants.CmdFixGit, constants.CmdFixGitAlias, "--fix-git", "fixgit"}, func() error { return cmdfixgit.RunFixGit(argsTail()) }},
		{[]string{constants.CmdUndo, constants.CmdUndoAlias}, func() error { return cmdundo.RunUndo(argsTail()) }},
		{[]string{constants.CmdHistoryPurge, constants.CmdHistoryPurgeAlias}, func() error { return cmdhistory.RunHistoryPurge(argsTail()) }},
		{[]string{constants.CmdHistoryPin, constants.CmdHistoryPinAlias}, func() error { return cmdhistory.RunHistoryPin(argsTail()) }},
		{[]string{"author"}, func() error { return cmdauthor.RunAuthor(argsTail()) }},
		{[]string{"sponsor"}, func() error { return cmdauthor.RunSponsor(argsTail()) }},
		{[]string{"credits"}, func() error { return cmdauthor.RunCredits(argsTail()) }},
	}
}

func toolingChromeEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{constants.CmdChromeProfileCopy}, func() error { return cmdchromeprofile.RunProfileCopy(argsTail()) }},
		{[]string{constants.CmdChromeProfileExport, constants.CmdChromeProfileExportAlias}, func() error { return cmdchromeprofile.RunProfileExport(argsTail()) }},
		{[]string{constants.CmdChromeProfileImport, constants.CmdChromeProfileImportAlias}, func() error { return cmdchromeprofile.RunProfileImport(argsTail()) }},
		{[]string{constants.CmdChromeProfileList, constants.CmdChromeProfileListAlias, constants.CmdChromeProfileListAlias2}, func() error { return cmdchromeprofile.RunProfileList(argsTail()) }},
		{[]string{constants.CmdChromeProfileDelete, constants.CmdChromeProfileDeleteAlias}, func() error { return cmdchromeprofile.RunProfileDelete(argsTail()) }},
		{[]string{constants.CmdChromeProfileMerge, constants.CmdChromeProfileMergeAlias}, func() error { return cmdchromeprofile.RunProfileMerge(argsTail()) }},
		{[]string{"chrome-profile-copy-all", "cpc-all", "copy-all"}, func() error { return cmdchromeprofile.RunCopyAll(argsTail()) }},
		{[]string{"chrome-profile-export-all", "cpe-all", "export-all"}, func() error { return cmdchromeprofile.RunExportAll(argsTail()) }},
		{[]string{"chrome-profile-import-all", "cpi-all", "import-all"}, func() error { return cmdchromeprofile.RunImportAll(argsTail()) }},
		{[]string{constants.CmdChrome, constants.CmdChromeAlias, constants.CmdChromeAlias2}, func() error { return cmdchrome.RunChrome(argsTail()) }},
	}
}

func runCpcDisambiguated(args []string) error {
	if cmdcommit.IsCommitPushHelpArg(args) {
		return cmdcommit.RunCommitPushChore(args)
	}

	if cmdpull.IsGitRepoCWD() {
		return cmdcommit.RunCommitPushChore(args)
	}

	return cmdchromeprofile.RunProfileCopy(args)
}

func toolingNetworkEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{constants.CmdServe, constants.CmdServeAlias}, func() error { return cmdservercmd.RunServe(argsTail()) }},
		{[]string{constants.CmdJoin, constants.CmdJoinAlias}, func() error {
			if err := cmdjoin.RunJoin(argsTail()); err != nil {
				return err
			}
			return nil
		}},
		{[]string{"ip"}, func() error { return runIP(argsTail()) }},
		{[]string{"ports", "port", "open-ports", "listening-ports"}, func() error { return cmdports.Run(argsTail()) }},
	}
}

func handleWhichSubcommand(args []string) error {
	if len(args) == 0 {
		return cmdconfig.RunWhichFormatCLI(args)
	}
	first := strings.ToLower(args[0])
	if first == "config" || first == "configs" {
		return cmdconfig.RunWhatConfigsCLI(args[1:])
	}
	if first == "format" || first == "json" {
		return cmdconfig.RunWhichFormatCLI(args[1:])
	}
	return cmdconfig.RunWhichFormatCLI(args)
}

func handleFormatSubcommand(args []string) error {
	if len(args) == 0 {
		return cmdconfig.RunWhichFormatCLI(args)
	}
	first := strings.ToLower(args[0])
	if first == "merge" {
		return cmdmerge.RunMergeJSONCLI(args[1:])
	}
	if first == "import" || first == "import-all" || first == "importall" {
		return cmdimport.RunImportAllJSONCLI(args[1:])
	}
	if first == "which" || first == "inspect" || first == "check" {
		return cmdconfig.RunWhichFormatCLI(args[1:])
	}
	return cmdconfig.RunWhichFormatCLI(args)
}

func handleJSONSubcommand(args []string) error {
	if len(args) == 0 {
		return cmdconfig.RunWhichFormatCLI(args)
	}
	first := strings.ToLower(args[0])
	if first == "merge" {
		return cmdmerge.RunMergeJSONCLI(args[1:])
	}
	if first == "import" || first == "import-all" || first == "importall" {
		return cmdimport.RunImportAllJSONCLI(args[1:])
	}
	return cmdconfig.RunWhichFormatCLI(args)
}

func toolingSystemEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{"history clean", "history-clean"}, func() error {
			return cmdhistory.RunHistoryPurgeCLI(argsTail())
		}},
		{[]string{"history-undo", "hu", "history restore", "history-restore", "undo-history"}, func() error {
			return cmdhistory.RunHistoryUndoCLI(argsTail())
		}},
		{[]string{"history"}, func() error {
			return handleHistoryParentDispatcher(argsTail())
		}},
	}
}

func handleHistoryParentDispatcher(args []string) error {
	if len(args) == 0 {
		return cmdhistory.RunHistory(args)
	}

	sub := strings.ToLower(args[0])
	switch sub {
	case "purge", "clean", "p":
		return cmdhistory.RunHistoryPurgeCLI(args[1:])
	case "undo", "restore", "u":
		return cmdhistory.RunHistoryUndoCLI(args[1:])
	default:
		return cmdhistory.RunHistory(args)
	}
}
