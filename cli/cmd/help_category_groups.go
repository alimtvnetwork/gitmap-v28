package cmd

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var helpCategoryGroups = map[string]func(){
	"scanning":       printGroupScanning,
	"cloning":        printGroupCloning,
	"gitops":         printGroupGitOps,
	"navigation":     printGroupNavigation,
	"release":        printGroupRelease,
	"release-info":   printGroupReleaseInfo,
	"releaseinfo":    printGroupReleaseInfo,
	"data":           printGroupData,
	"import-export":  printGroupImportExport,
	"importexport":   printGroupImportExport,
	"history":        printGroupHistory,
	"amend":          printGroupAmend,
	"project":        printGroupProject,
	"ssh":            printGroupSSH,
	"zip":            printGroupZip,
	"env-tools":      printGroupEnvTools,
	"envtools":       printGroupEnvTools,
	"env":            printGroupEnvTools,
	"tasks":          printGroupTasks,
	"visualize":      printGroupVisualize,
	"commit-xfer":    printGroupCommitXfer,
	"commitxfer":     printGroupCommitXfer,
	"pr":             printGroupPR,
	"chrome-profile": printGroupChromeProfile,
	"chrome":         printGroupChromeProfile,
	"templates":      printGroupTemplates,
	"cluster":        printGroupCluster,
	"user":           printGroupUser,
	"installers":     printGroupInstallers,
	"integrations":   printGroupIntegrations,
	"search-find":    printGroupSearchFind,
	"search":         printGroupSearchFind,
	"find":           printGroupSearchFind,
	"utilities":      printGroupUtilities,
	"utility":        printGroupUtilities,
	"advanced":       printGroupUtilities,
}

func isHelpCategoryGroup(topic string) bool {
	low := strings.ToLower(strings.TrimSpace(topic))
	_, exists := helpCategoryGroups[low]
	return exists
}

func printUsageCategoryGroup(topic string) {
	low := strings.ToLower(strings.TrimSpace(topic))
	fn, exists := helpCategoryGroups[low]
	if !exists {
		return
	}

	measuringHelp = true
	maxHelpCmdLen = 0
	fn()

	measuringHelp = false
	fmt.Printf(constants.UsageHeaderFmt, constants.Version)
	fmt.Println()
	fn()

	printCategorySpecificTips(low)
	fmt.Println()
	printUsageFooterShort()
}

func printCategorySpecificTips(group string) {
	if group == "search" || group == "search-find" || group == "find" {
		printSearchTaxonomyTip()
	}
	if group == "integrations" || group == "agm" {
		printAGMTip()
	}
}

func printSearchTaxonomyTip() {
	fmt.Println()
	fmt.Printf("  %s💡 SEARCH TAXONOMY & COMMAND DISTINCTION:%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("    • %sgitmap search <text>%s      Multi-repository code content/grep search walk\n", constants.ColorYellow, constants.ColorReset)
	fmt.Printf("    • %sgitmap aum search <text>%s  Automation Unit Manager macro & script index search\n", constants.ColorYellow, constants.ColorReset)
	fmt.Printf("    • %sgitmap find / sf <name>%s   Fast filename and path index search\n", constants.ColorYellow, constants.ColorReset)
}

func printAGMTip() {
	fmt.Println()
	fmt.Printf("  %s💡 ANTIGRAVITY MANAGER (AGM) WORKFLOWS:%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("    • %sgitmap agm install%s        Install Antigravity Manager GUI desktop application\n", constants.ColorYellow, constants.ColorReset)
	fmt.Printf("    • %sgitmap agm update [--ssh]%s Update AGM desktop app locally or across fleet nodes\n", constants.ColorYellow, constants.ColorReset)
	fmt.Printf("    • %sgitmap agy%s                Antigravity prompt queues, reruns, and SUG engine\n", constants.ColorYellow, constants.ColorReset)
}
