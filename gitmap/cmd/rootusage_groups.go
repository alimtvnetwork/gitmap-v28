package cmd

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func printGroupScanning() {
	fmt.Println()
	fmt.Println(colorGroupHeader(constants.HelpGroupScanning))
	fmt.Println(constants.HelpScan)
	fmt.Println(constants.HelpRescan)
	fmt.Println(constants.HelpList)
}

func printGroupCloning() {
	fmt.Println()
	fmt.Println(colorGroupHeader(constants.HelpGroupCloning))
	fmt.Println(constants.HelpClone)
	fmt.Println(constants.HelpCloneNext)
	fmt.Println(constants.HelpDesktopSync)
	fmt.Println(constants.HelpGitHubDesktop)
}

func printGroupGitOps() {
	fmt.Println()
	fmt.Println(colorGroupHeader(constants.HelpGroupGitOps))
	fmt.Println(constants.HelpPull)
	fmt.Println(constants.HelpPullAll)
	fmt.Println(constants.HelpExec)
	fmt.Println(constants.HelpStatus)
	fmt.Println(constants.HelpWatch)
	fmt.Println(constants.HelpHasAnyUpdates)
	fmt.Println(constants.HelpLatestBr)
	fmt.Println(constants.MsgHelpLFSCommon)
}

func printGroupNavigation() {
	fmt.Println()
	fmt.Println(colorGroupHeader(constants.HelpGroupNavigation))
	fmt.Println(constants.HelpCD)
	fmt.Println(constants.HelpGroup)
	fmt.Println(constants.HelpMultiGroup)
	fmt.Println(constants.HelpSf)
	fmt.Println(constants.HelpAlias)
	fmt.Println(constants.HelpDiffProfiles)
}

func printGroupRelease() {
	fmt.Println()
	fmt.Println(colorGroupHeader(constants.HelpGroupRelease))
	fmt.Println(constants.HelpRelease)
	fmt.Println(constants.HelpReleasePull)
	fmt.Println(constants.HelpReleaseSelf)
	fmt.Println(constants.HelpReleaseBr)
	fmt.Println(constants.HelpTempRelease)
}

func printGroupReleaseInfo() {
	fmt.Println()
	fmt.Println(colorGroupHeader(constants.HelpGroupReleaseInfo))
	fmt.Println(constants.HelpChangelog)
	fmt.Println(constants.HelpChangelogGen)
	fmt.Println(constants.HelpListVersions)
	fmt.Println(constants.HelpListReleases)
	fmt.Println(constants.HelpReleasePend)
	fmt.Println(constants.HelpRevert)
	fmt.Println(constants.HelpClearReleaseJSON)
	fmt.Println(constants.HelpPrune)
}

func printGroupData() {
	fmt.Println()
	fmt.Println(colorGroupHeader(constants.HelpGroupData))
	fmt.Println(constants.HelpExport)
	fmt.Println(constants.HelpImport)
	fmt.Println(constants.HelpProfile)
	fmt.Println(constants.HelpBookmark)
	fmt.Println(constants.HelpMV)
	fmt.Println(constants.HelpRm)
	fmt.Println(constants.HelpDBReset)
}

func printGroupHistory() {
	fmt.Println()
	fmt.Println(colorGroupHeader(constants.HelpGroupHistory))
	fmt.Println(constants.HelpHistory)
	fmt.Println(constants.HelpHistoryReset)
	fmt.Println(constants.HelpVersionHistory)
	fmt.Println(constants.HelpStats)
}

func printGroupAmend() {
	fmt.Println()
	fmt.Println(colorGroupHeader(constants.HelpGroupAmendGroup))
	fmt.Println(constants.HelpAmend)
	fmt.Println(constants.HelpAmendList)
}

func printGroupProject() {
	fmt.Println()
	fmt.Println(colorGroupHeader(constants.HelpGroupProject))
	fmt.Println(constants.HelpGoRepos)
	fmt.Println(constants.HelpNodeRepos)
	fmt.Println(constants.HelpReactRepos)
	fmt.Println(constants.HelpCppRepos)
	fmt.Println(constants.HelpCsharpRepos)
}

func printGroupSSH() {
	fmt.Println()
	fmt.Println(colorGroupHeader(constants.HelpGroupSSH))
	fmt.Println(constants.HelpSSH)
}

func printGroupZip() {
	fmt.Println()
	fmt.Println(colorGroupHeader(constants.HelpGroupZip))
	fmt.Println(constants.HelpZipGroup)
}

func printGroupEnvTools() {
	fmt.Println()
	fmt.Println(colorGroupHeader(constants.HelpGroupEnvTools))
	fmt.Println(constants.HelpEnv)
	fmt.Println(constants.HelpInstall)
	fmt.Println(constants.HelpUninstall)
}

func printGroupTasks() {
	fmt.Println()
	fmt.Println(colorGroupHeader(constants.HelpGroupTasks))
	fmt.Println(constants.HelpTask)
	fmt.Println(constants.HelpPending)
	fmt.Println(constants.HelpDoPending)
	fmt.Println(constants.HelpMacro)
	fmt.Println(constants.HelpExecute)
}

func printGroupVisualize() {
	fmt.Println()
	fmt.Println(colorGroupHeader(constants.HelpGroupVisualize))
	fmt.Println(constants.HelpDashboard)
}

func printGroupCommitXfer() {
	fmt.Println()
	fmt.Println(colorGroupHeader(constants.HelpGroupCommitXfer))
	fmt.Println(constants.HelpCommitRight)
	fmt.Println(constants.HelpCommitLeft)
	fmt.Println(constants.HelpCommitBoth)
}

func printGroupChromeProfile() {
	fmt.Println()
	fmt.Println(colorGroupHeader(constants.HelpGroupChromeProf))
	fmt.Println(constants.HelpChromeProfileCopy)
	fmt.Println(constants.HelpChromeProfileExport)
	fmt.Println(constants.HelpChromeProfileImport)
	fmt.Println(constants.HelpChromeProfileList)
	fmt.Println(constants.HelpChromeProfileDelete)
}

func printGroupTemplates() {
	fmt.Println()
	fmt.Println(colorGroupHeader(constants.HelpGroupTemplates))
	fmt.Println(constants.HelpAddIgnore)
	fmt.Println(constants.HelpAddAttributes)
	fmt.Println(constants.HelpAddLFSInstall)
	fmt.Println(constants.HelpTemplatesInit)
	fmt.Println(constants.HelpTemplatesList)
	fmt.Println(constants.HelpTemplatesShow)
	fmt.Println(constants.HelpTemplatesDiff)
	fmt.Println(constants.HelpSync)
	fmt.Println(constants.HelpCommons)
}

func printGroupCluster() {
	fmt.Println()
	fmt.Println(colorGroupHeader(constants.HelpGroupCluster))
	fmt.Println(constants.HelpServersClients)
	fmt.Println(constants.HelpClients)
	fmt.Println(constants.HelpCluster)
	fmt.Println(constants.HelpServe)
}
func printGroupUser() {
	fmt.Println()
	fmt.Println(colorGroupHeader(constants.HelpGroupUser))
	fmt.Println(constants.HelpUser)
}
