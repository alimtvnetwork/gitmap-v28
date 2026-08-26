package cmd

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

const (
	superCategoryLineWidth = 58
	minRuleLength          = 4
)

// printUsage displays grouped help text for all commands and flags.
func printUsage() {
	printUsageHeader()
	printUsageCoreCategories()
	printUsageReleaseAndProjectCategories()
	printUsageAdvancedCategories()
	printUsageTrailer()
}

// printUsageHeader prints the header banner, usage line, and quick-start section.
func printUsageHeader() {
	fmt.Printf(constants.UsageHeaderFmt, constants.Version)
	fmt.Println(constants.HelpUsage)
	fmt.Println()
	printUsageQuickStart()
}

// printUsageCoreCategories prints the Get Started and Work With Repos categories.
func printUsageCoreCategories() {
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
}

// printUsageReleaseAndProjectCategories prints the Release and Projects categories.
func printUsageReleaseAndProjectCategories() {
	printSuperCategory("RELEASE & HISTORY", printUsageReleaseGroups)
	printSuperCategory("PROJECTS & DATA", printUsageProjectsGroups)
}

func printUsageReleaseGroups() {
	printGroupRelease()
	printGroupReleaseInfo()
	printGroupHistory()
	printGroupAmend()
	printGroupCommitXfer()
}

func printUsageProjectsGroups() {
	printGroupProject()
	printGroupData()
	printGroupImportExport()
	printGroupChromeProfile()
	printGroupZip()
	printGroupTasks()
	printGroupVisualize()
}

// printUsageAdvancedCategories prints the Cluster and Advanced categories.
func printUsageAdvancedCategories() {
	printSuperCategory("INSTALL & AUTOMATION", func() {
		printGroupInstallers()
	})
	printSuperCategory("INTEGRATIONS", func() {
		printGroupIntegrations()
	})
	printSuperCategory("CLUSTER & NETWORK", func() {

		printGroupCluster()
		printGroupUser()
	})
	printSuperCategory("ADVANCED", func() {
		printGroupUtilities()
	})
}

// printUsageTrailer prints the flag sections and usage footer.
func printUsageTrailer() {
	fmt.Println()
	printUsageFlagSections()
	printUsageFooter()
}

// printSuperCategory renders a bold intent-banner above a set of
// related sub-groups, so users can pinpoint the right area without
// scanning 17 sub-headers.
func printSuperCategory(title string, body func()) {
	fmt.Println()
	rule := repeatRule(superCategoryLineWidth - len(title))
	fmt.Println("  " + constants.ColorMagenta + "━━ " +
		constants.ColorWhite + title + constants.ColorReset +
		" " + constants.ColorMagenta + rule + constants.ColorReset)
	body()
}

func repeatRule(count int) string {
	if count < minRuleLength {
		count = minRuleLength
	}
	out := ""
	for idx := 0; idx < count; idx++ {
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
	fmt.Println()
	fmt.Println(constants.HelpExampleScan)
	fmt.Println(constants.HelpExampleList)
	fmt.Println(constants.HelpExamplePull)
	fmt.Println(constants.HelpExampleCD)
	fmt.Println()
	fmt.Println(colorGroupHeader(constants.HelpGroupHint))
	fmt.Println()
	fmt.Println(constants.HelpCompactHint)
}
