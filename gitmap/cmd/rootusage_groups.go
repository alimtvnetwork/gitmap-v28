package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func printGroupScanning() {
	renderHeader(constants.HelpGroupScanning)
	renderLine(constants.HelpScan)
	renderLine(constants.HelpRescan)
	renderLine(constants.HelpList)
}

func printGroupCloning() {
	renderHeader(constants.HelpGroupCloning)
	renderLine(constants.HelpClone)
	renderLine(constants.HelpCloneSync)
	renderLine(constants.HelpCloneNext)
	renderLine(constants.HelpDesktopSync)
	renderLine(constants.HelpGitHubDesktop)
}

func printGroupGitOps() {
	renderHeader(constants.HelpGroupGitOps)
	renderLine(constants.HelpPull)
	renderLine(constants.HelpPullAll)
	renderLine(constants.HelpFix)
	renderLine(constants.HelpExec)
	renderLine(constants.HelpStatus)
	renderLine(constants.HelpWatch)
	renderLine(constants.HelpHasAnyUpdates)
	renderLine(constants.HelpLatestBr)
	renderLine(constants.MsgHelpLFSCommon)
}

func printGroupNavigation() {
	renderHeader(constants.HelpGroupNavigation)
	renderLine(constants.HelpCD)
	renderLine(constants.HelpGroup)
	renderLine(constants.HelpMultiGroup)
	renderLine(constants.HelpSf)
	renderLine(constants.HelpAlias)
	renderLine(constants.HelpDiffProfiles)
}

func printGroupRelease() {
	renderHeader(constants.HelpGroupRelease)
	renderLine(constants.HelpRelease)
	renderLine(constants.HelpReleasePull)
	renderLine(constants.HelpReleaseSelf)
	renderLine(constants.HelpReleaseBr)
	renderLine(constants.HelpTempRelease)
}

func printGroupReleaseInfo() {
	renderHeader(constants.HelpGroupReleaseInfo)
	renderLine(constants.HelpChangelog)
	renderLine(constants.HelpChangelogGen)
	renderLine(constants.HelpListVersions)
	renderLine(constants.HelpListReleases)
	renderLine(constants.HelpReleasePend)
	renderLine(constants.HelpRevert)
	renderLine(constants.HelpClearReleaseJSON)
	renderLine(constants.HelpPrune)
}

func printGroupData() {
	renderHeader(constants.HelpGroupData)
	renderLine(constants.HelpExport)
	renderLine(constants.HelpImport)
	renderLine(constants.HelpProfile)
	renderLine(constants.HelpBookmark)
	renderLine(constants.HelpMV)
	renderLine(constants.HelpRm)
	renderLine(constants.HelpDB)
	renderLine(constants.HelpDBReset)
	renderLine(constants.HelpStartFresh)
	renderLine(constants.HelpFindDuplicates)
}

func printGroupImportExport() {
	renderHeader(constants.HelpGroupImportExport)
	renderLine(constants.HelpImportExport)
	renderLine(constants.HelpExportSummary)
	renderLine(constants.HelpImportSummary)
}

func printGroupHistory() {
	renderHeader(constants.HelpGroupHistory)
	renderLine(constants.HelpHistory)
	renderLine(constants.HelpHistoryReset)
	renderLine(constants.HelpVersionHistory)
	renderLine(constants.HelpStats)
}

func printGroupAmend() {
	renderHeader(constants.HelpGroupAmendGroup)
	renderLine(constants.HelpAmend)
	renderLine(constants.HelpAmendList)
}

func printGroupProject() {
	renderHeader(constants.HelpGroupProject)
	renderLine(constants.HelpGoRepos)
	renderLine(constants.HelpNodeRepos)
	renderLine(constants.HelpReactRepos)
	renderLine(constants.HelpCppRepos)
	renderLine(constants.HelpCsharpRepos)
}

func printGroupSSH() {
	renderHeader(constants.HelpGroupSSH)
	renderLine(constants.HelpSSH)
	renderLine(constants.HelpSSHJoin)
}

func printGroupZip() {
	renderHeader(constants.HelpGroupZip)
	renderLine(constants.HelpZipGroup)
}

func printGroupEnvTools() {
	renderHeader(constants.HelpGroupEnvTools)
	renderLine(constants.HelpEnv)
	renderLine(constants.HelpCodingGuideline)
}

func printGroupTasks() {
	renderHeader(constants.HelpGroupTasks)
	renderLine(constants.HelpTask)
	renderLine(constants.HelpPending)
	renderLine(constants.HelpDoPending)
	renderLine(constants.HelpMacro)
	renderLine(constants.HelpExecute)
}

func printGroupVisualize() {
	renderHeader(constants.HelpGroupVisualize)
	renderLine(constants.HelpDashboard)
}

func printGroupCommitXfer() {
	renderHeader(constants.HelpGroupCommitXfer)
	renderLine(constants.HelpCommitRight)
	renderLine(constants.HelpCommitLeft)
	renderLine(constants.HelpCommitBoth)
}

func printGroupChromeProfile() {
	renderHeader(constants.HelpGroupChromeProf)
	renderLine(constants.HelpChrome)
	renderLine(constants.HelpChromeProfileCopy)
	renderLine(constants.HelpChromeProfileExport)
	renderLine(constants.HelpChromeProfileImport)
	renderLine(constants.HelpChromeProfileList)
	renderLine(constants.HelpChromeProfileDelete)
	renderLine(constants.HelpChromeProfileMerge)
}

func printGroupTemplates() {
	renderHeader(constants.HelpGroupTemplates)
	renderLine(constants.HelpAddIgnore)
	renderLine(constants.HelpAddAttributes)
	renderLine(constants.HelpAddLFSInstall)
	renderLine(constants.HelpTemplatesInit)
	renderLine(constants.HelpTemplatesList)
	renderLine(constants.HelpTemplatesShow)
	renderLine(constants.HelpTemplatesDiff)
	renderLine(constants.HelpSync)
	renderLine(constants.HelpCommons)
}

func printGroupCluster() {
	renderHeader(constants.HelpGroupCluster)
	renderLine(constants.HelpServersClients)
	renderLine(constants.HelpClients)
	renderLine(constants.HelpCluster)
	renderLine(constants.HelpServe)
}

func printGroupUser() {
	renderHeader(constants.HelpGroupUser)
	renderLine(constants.HelpUser)
}

func printGroupInstallers() {
	renderHeader(constants.HelpGroupInstallers)
	renderLine(constants.HelpInstall)
	renderLine(constants.HelpUninstall)
	renderLine(constants.HelpInstaller)
	renderLine(constants.HelpMacro)
	renderLine(constants.HelpSetup)
}

func printGroupIntegrations() {
	renderHeader(constants.HelpGroupIntegrations)
	renderLine(constants.HelpVSCode)
	renderLine(constants.HelpAgy)
	renderLine(constants.HelpSchedule)
	renderLine(constants.HelpPipeline)
	renderLine(constants.HelpUI)
}

func printGroupSearchFind() {
	renderHeader(constants.HelpGroupSearchFind)
	renderLine(constants.HelpFindFiles)
	renderLine(constants.HelpFindFilesAny)
	renderLine(constants.HelpFindFilesStart)
	renderLine(constants.HelpFindFilesEnd)
	renderLine(constants.HelpFind)
	renderLine(constants.HelpListFiles)
}
