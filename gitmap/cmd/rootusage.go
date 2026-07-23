package cmd

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// printUsage displays grouped help text for all commands and flags.
func printUsage() {
	fmt.Printf(constants.UsageHeaderFmt, constants.Version)
	fmt.Println(constants.HelpUsage)
	fmt.Println()
	printUsageQuickStart()

	printSuperCategory("GET STARTED", func() {
		printGroupScanning()
		printGroupNavigation()
		printGroupEnvTools()
		printGroupTemplates()
	})
	printSuperCategory("WORK WITH REPOS", func() {
		printGroupCloning()
		printGroupGitOps()
		printGroupSSH()
	})
	printSuperCategory("RELEASE & HISTORY", func() {
		printGroupRelease()
		printGroupReleaseInfo()
		printGroupHistory()
		printGroupAmend()
		printGroupCommitXfer()
	})
	printSuperCategory("PROJECTS & DATA", func() {
		printGroupProject()
		printGroupData()
		printGroupChromeProfile()
		printGroupZip()
		printGroupTasks()
		printGroupVisualize()
	})
	printSuperCategory("ADVANCED", func() {
		printGroupUtilities()
	})

	fmt.Println()
	printUsageFlagSections()
	printUsageFooter()
}

// printSuperCategory renders a bold intent-banner above a set of
// related sub-groups, so users can pinpoint the right area without
// scanning 17 sub-headers. The banner uses box-drawing rules that the
// glyph filter downgrades to "==" on legacy PowerShell hosts.
func printSuperCategory(title string, body func()) {
	fmt.Println()
	rule := repeatRule(58 - len(title))
	fmt.Println("  " + constants.ColorMagenta + "━━ " +
		constants.ColorWhite + title + constants.ColorReset +
		" " + constants.ColorMagenta + rule + constants.ColorReset)
	body()
}

func repeatRule(n int) string {
	if n < 4 {
		n = 4
	}
	out := ""
	for i := 0; i < n; i++ {
		out += "━"
	}
	return out
}

// colorGroupHeader wraps a sub-group header line in bold cyan so each
// section stands out from the muted command rows beneath it.
func colorGroupHeader(header string) string {
	return constants.ColorCyan + header + constants.ColorReset
}

// printUsageQuickStart prints examples and the help hint.
func printUsageQuickStart() {
	fmt.Println(colorGroupHeader(constants.HelpGroupExample))
	fmt.Println(constants.HelpExampleScan)
	fmt.Println(constants.HelpExampleList)
	fmt.Println(constants.HelpExamplePull)
	fmt.Println(constants.HelpExampleCD)
	fmt.Println()
	fmt.Println(colorGroupHeader(constants.HelpGroupHint))
	fmt.Println(constants.HelpCompactHint)
}

// printGroupScanning prints the scanning commands.
func printGroupScanning() {
	fmt.Println()
	fmt.Println(colorGroupHeader(constants.HelpGroupScanning))
	fmt.Println(constants.HelpScan)
	fmt.Println(constants.HelpRescan)
	fmt.Println(constants.HelpList)
}

// printGroupCloning prints the cloning commands.
func printGroupCloning() {
	fmt.Println()
	fmt.Println(colorGroupHeader(constants.HelpGroupCloning))
	fmt.Println(constants.HelpClone)
	fmt.Println(constants.HelpCloneNext)
	fmt.Println(constants.HelpDesktopSync)
	fmt.Println(constants.HelpGitHubDesktop)
}

// printGroupGitOps prints the git operations commands.
func printGroupGitOps() {
	fmt.Println()
	fmt.Println(colorGroupHeader(constants.HelpGroupGitOps))
	fmt.Println(constants.HelpPull)
	fmt.Println(constants.HelpExec)
	fmt.Println(constants.HelpStatus)
	fmt.Println(constants.HelpWatch)
	fmt.Println(constants.HelpHasAnyUpdates)
	fmt.Println(constants.HelpLatestBr)
	fmt.Println(constants.MsgHelpLFSCommon)
}

// printGroupNavigation prints the navigation commands.
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

// printGroupRelease prints the release workflow commands.
func printGroupRelease() {
	fmt.Println()
	fmt.Println(colorGroupHeader(constants.HelpGroupRelease))
	fmt.Println(constants.HelpRelease)
	fmt.Println(constants.HelpReleasePull)
	fmt.Println(constants.HelpReleaseSelf)
	fmt.Println(constants.HelpReleaseBr)
	fmt.Println(constants.HelpTempRelease)
}

// printGroupReleaseInfo prints the release info commands.
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

// printGroupData prints the data/profile/bookmark commands.
func printGroupData() {
	fmt.Println()
	fmt.Println(colorGroupHeader(constants.HelpGroupData))
	fmt.Println(constants.HelpExport)
	fmt.Println(constants.HelpImport)
	fmt.Println(constants.HelpProfile)
	fmt.Println(constants.HelpBookmark)
	fmt.Println(constants.HelpRm)
	fmt.Println(constants.HelpDBReset)
}

// printGroupHistory prints the history and stats commands.
func printGroupHistory() {
	fmt.Println()
	fmt.Println(colorGroupHeader(constants.HelpGroupHistory))
	fmt.Println(constants.HelpHistory)
	fmt.Println(constants.HelpHistoryReset)
	fmt.Println(constants.HelpVersionHistory)
	fmt.Println(constants.HelpStats)
}

// printGroupAmend prints the amend commands.
func printGroupAmend() {
	fmt.Println()
	fmt.Println(colorGroupHeader(constants.HelpGroupAmendGroup))
	fmt.Println(constants.HelpAmend)
	fmt.Println(constants.HelpAmendList)
}

// printGroupProject prints the project detection commands.
func printGroupProject() {
	fmt.Println()
	fmt.Println(colorGroupHeader(constants.HelpGroupProject))
	fmt.Println(constants.HelpGoRepos)
	fmt.Println(constants.HelpNodeRepos)
	fmt.Println(constants.HelpReactRepos)
	fmt.Println(constants.HelpCppRepos)
	fmt.Println(constants.HelpCsharpRepos)
}

// printGroupSSH prints the SSH key management commands.
func printGroupSSH() {
	fmt.Println()
	fmt.Println(colorGroupHeader(constants.HelpGroupSSH))
	fmt.Println(constants.HelpSSH)
}

// printGroupZip prints the zip group commands.
func printGroupZip() {
	fmt.Println()
	fmt.Println(colorGroupHeader(constants.HelpGroupZip))
	fmt.Println(constants.HelpZipGroup)
}

// printGroupEnvTools prints the env and install commands.
func printGroupEnvTools() {
	fmt.Println()
	fmt.Println(colorGroupHeader(constants.HelpGroupEnvTools))
	fmt.Println(constants.HelpEnv)
	fmt.Println(constants.HelpInstall)
	fmt.Println(constants.HelpUninstall)
}

// printGroupTasks prints the task commands.
func printGroupTasks() {
	fmt.Println()
	fmt.Println(colorGroupHeader(constants.HelpGroupTasks))
	fmt.Println(constants.HelpTask)
	fmt.Println(constants.HelpPending)
	fmt.Println(constants.HelpDoPending)
}

// printGroupVisualize prints the visualization commands.
func printGroupVisualize() {
	fmt.Println()
	fmt.Println(colorGroupHeader(constants.HelpGroupVisualize))
	fmt.Println(constants.HelpDashboard)
}

// printGroupCommitXfer prints the commit-transfer family. Surfacing
// these in `gitmap help` is the primary discovery path — without it
// users have to read spec/01-app/106 or stumble into `help commit-right`
// to learn the aliases (cml / cmr / cmb) even exist.
func printGroupCommitXfer() {
	fmt.Println()
	fmt.Println(colorGroupHeader(constants.HelpGroupCommitXfer))
	fmt.Println(constants.HelpCommitRight)
	fmt.Println(constants.HelpCommitLeft)
	fmt.Println(constants.HelpCommitBoth)
}

// printGroupChromeProfile prints the chrome-profile-* family so users
// can discover cpc / cpe / cpi / cpl / cpd from a bare `gitmap help`.
// Without this block these commands were invisible to fuzzy search.
func printGroupChromeProfile() {
	fmt.Println()
	fmt.Println(colorGroupHeader(constants.HelpGroupChromeProf))
	fmt.Println(constants.HelpChromeProfileCopy)
	fmt.Println(constants.HelpChromeProfileExport)
	fmt.Println(constants.HelpChromeProfileImport)
	fmt.Println(constants.HelpChromeProfileList)
	fmt.Println(constants.HelpChromeProfileDelete)
}

// printGroupTemplates surfaces `add ignore`, `add attributes`,
// `add lfs-install` and the `templates` family so users can discover
// the scaffolding commands from a bare `gitmap help`. Prior to this
// group they were only reachable via `gitmap help add-ignore` etc.
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
