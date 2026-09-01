package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// helpRow is a single command/flag line tagged with its sub-group
// header so filtered output keeps its context.
type helpRow struct {
	Group string
	Line  string
}

// allHelpRows returns every command + flag row rendered by the
// full `gitmap help` screen. Source of truth for `--filter` search.
func allHelpRows() []helpRow {
	rows := make([]helpRow, 0, 128)
	addGroup(&rows, constants.HelpGroupScanning,
		constants.HelpScan, constants.HelpRescan, constants.HelpList)
	addGroup(&rows, constants.HelpGroupCloning,
		constants.HelpClone, constants.HelpCloneSync, constants.HelpCloneNext,
		constants.HelpDesktopSync, constants.HelpGitHubDesktop)
	addGroup(&rows, constants.HelpGroupGitOps,
		constants.HelpPull, constants.HelpExec, constants.HelpStatus,
		constants.HelpWatch, constants.HelpHasAnyUpdates, constants.HelpLatestBr)
	addGroup(&rows, constants.HelpGroupNavigation,
		constants.HelpCD, constants.HelpGroup, constants.HelpMultiGroup,
		constants.HelpSf, constants.HelpAlias, constants.HelpDiffProfiles)
	addGroup(&rows, constants.HelpGroupRelease,
		constants.HelpRelease, constants.HelpReleasePull,
		constants.HelpReleaseSelf, constants.HelpReleaseBr, constants.HelpTempRelease)
	addGroup(&rows, constants.HelpGroupReleaseInfo,
		constants.HelpChangelog, constants.HelpChangelogGen,
		constants.HelpListVersions, constants.HelpListReleases,
		constants.HelpReleasePend, constants.HelpRevert,
		constants.HelpClearReleaseJSON, constants.HelpPrune)
	addGroup(&rows, constants.HelpGroupData,
		constants.HelpExport, constants.HelpImport, constants.HelpProfile,
		constants.HelpBookmark, constants.HelpMV, constants.HelpRm, constants.HelpDBReset)
	addGroup(&rows, constants.HelpGroupImportExport,
		constants.HelpImportExport, constants.HelpExportSummary, constants.HelpImportSummary)
	addGroup(&rows, constants.HelpGroupHistory,
		constants.HelpHistory, constants.HelpHistoryReset,
		constants.HelpVersionHistory, constants.HelpStats)
	addGroup(&rows, constants.HelpGroupAmendGroup,
		constants.HelpAmend, constants.HelpAmendList)
	addGroup(&rows, constants.HelpGroupProject,
		constants.HelpGoRepos, constants.HelpNodeRepos,
		constants.HelpReactRepos, constants.HelpCppRepos, constants.HelpCsharpRepos)
	addGroup(&rows, constants.HelpGroupSSH, constants.HelpSSH)
	addGroup(&rows, constants.HelpGroupZip, constants.HelpZipGroup)
	addGroup(&rows, constants.HelpGroupEnvTools,
		constants.HelpEnv, constants.HelpInstall, constants.HelpUninstall)
	addGroup(&rows, constants.HelpGroupTasks,
		constants.HelpTask, constants.HelpPending, constants.HelpDoPending,
		constants.HelpMacro, constants.HelpExecute)
	addGroup(&rows, constants.HelpGroupVisualize, constants.HelpDashboard)
	addGroup(&rows, constants.HelpGroupCommitXfer,
		constants.HelpCommitRight, constants.HelpCommitLeft, constants.HelpCommitBoth)
	addGroup(&rows, constants.HelpGroupChromeProf,
		constants.HelpChrome,
		constants.HelpChromeProfileCopy, constants.HelpChromeProfileExport,
		constants.HelpChromeProfileImport, constants.HelpChromeProfileList,
		constants.HelpChromeProfileDelete, constants.HelpChromeProfileMerge)
	addGroup(&rows, constants.HelpGroupTemplates,
		constants.HelpAddIgnore, constants.HelpAddAttributes,
		constants.HelpAddLFSInstall, constants.HelpTemplatesInit,
		constants.HelpTemplatesList, constants.HelpTemplatesShow,
		constants.HelpTemplatesDiff, constants.HelpSync, constants.HelpCommons)
	addGroup(&rows, constants.HelpGroupSearchFind,
		constants.HelpFindFiles, constants.HelpFindFilesAny, constants.HelpFindFilesStart,
		constants.HelpFindFilesEnd, constants.HelpFind, constants.HelpListFiles)
	addGroup(&rows, constants.HelpGroupInstallers,
		constants.HelpInstall, constants.HelpUninstall, constants.HelpInstaller, constants.HelpMacro, constants.HelpSetup)
	addGroup(&rows, constants.HelpGroupIntegrations,
		constants.HelpVSCode, constants.HelpAgy, constants.HelpSchedule, constants.HelpPipeline, constants.HelpUI)
	addGroup(&rows, constants.HelpGroupCluster,
		constants.HelpServersClients, constants.HelpClients, constants.HelpCluster, constants.HelpServe)
	addGroup(&rows, constants.HelpGroupUtilities,
		constants.HelpDoctor, constants.HelpUpdate,
		constants.HelpUpdateCleanup, constants.HelpVersion, constants.HelpCompletion,
		constants.HelpInteractive, constants.HelpDocs, constants.HelpHelpDash,
		constants.HelpGoMod, constants.HelpSEOWrite, constants.HelpLLMDocs,
		constants.HelpFixRepo, constants.HelpMakePublic, constants.HelpMakePrivate,
		constants.HelpCloneFixRepo, constants.HelpCloneFixRepoPub,
		constants.HelpCmdOpen, constants.HelpHelp)

	return rows
}

func addGroup(rows *[]helpRow, group string, lines ...string) {
	for _, ln := range lines {
		*rows = append(*rows, helpRow{Group: group, Line: ln})
	}
}
