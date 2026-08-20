package cmd

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
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
	printSuperCategory("CLUSTER & NETWORK", func() {
		printGroupCluster()
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
// scanning 17 sub-headers.
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
