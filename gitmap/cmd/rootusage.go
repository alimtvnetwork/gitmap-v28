package cmd

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/pterm/pterm"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

const (
	superCategoryLineWidth = 58
	minRuleLength          = 4
)

// printUsage displays grouped help text for all commands and flags.
func printUsage() {
	measuringHelp = true
	maxHelpCmdLen = 0
	printUsageCoreCategories()
	printUsageReleaseAndProjectCategories()
	printUsageAdvancedCategories()

	measuringHelp = false
	printUsageHeader()
	printUsageCoreCategories()
	printUsageReleaseAndProjectCategories()
	printUsageAdvancedCategories()
	printUsageTrailer()
}

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
		printGroupSearchFind()
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
	if measuringHelp {
		body()
		return
	}
	fmt.Println()
	rule := repeatRule(superCategoryLineWidth - len(title))
	fmt.Println("  " + constants.ColorMagenta + "━━ " +
		constants.ColorWhite + title + constants.ColorReset +
		" " + constants.ColorMagenta + rule + constants.ColorReset)
	fmt.Println()
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

var (
	measuringHelp bool
	maxHelpCmdLen int
)

func renderHeader(header string) {
	if measuringHelp {
		return
	}
	fmt.Println()
	if strings.HasPrefix(header, "  ") {
		fmt.Println(colorGroupHeader(header))
	} else {
		fmt.Println(colorGroupHeader("  " + header))
	}
	fmt.Println()
}

func renderFlagHeader(header string) {
	if measuringHelp {
		return
	}
	fmt.Println()
	if strings.HasPrefix(header, "  ") {
		fmt.Println(colorGroupHeader(header))
	} else {
		fmt.Println(colorGroupHeader("  " + header))
	}
	fmt.Println()
}

const maxCmdColumnWidth = 38

func updateMaxHelpCmdLen(cmd string) {
	l := lipgloss.Width(cmd)
	if l > maxCmdColumnWidth {
		l = maxCmdColumnWidth
	}
	if l > maxHelpCmdLen {
		maxHelpCmdLen = l
	}
}

func extractCmdAndDesc(line string) (string, string) {
	line = strings.TrimLeft(line, " ")
	for i := 0; i < len(line)-1; i++ {
		if line[i] == ' ' && line[i+1] == ' ' {
			return strings.TrimRight(line[:i], " "), strings.TrimLeft(line[i:], " ")
		}
	}

	return line, ""
}

// ExpandableMarkerResult contains the extracted command, description, and expandable marker.
type ExpandableMarkerResult struct {
	Cmd    string
	Desc   string
	Marker string
}

func parseExpandableMarker(cmd, desc string) ExpandableMarkerResult {
	if strings.Contains(desc, "(use --help to expand)") {
		cleanDesc := strings.TrimSpace(strings.ReplaceAll(desc, "(use --help to expand)", ""))
		firstWord := strings.Split(cmd, " ")[0]
		return ExpandableMarkerResult{
			Cmd:    cmd,
			Desc:   cleanDesc,
			Marker: fmt.Sprintf("  ▸ subcommands — see `gitmap %s --help`", firstWord),
		}
	}

	if strings.Contains(cmd, "(use --help to expand)") {
		cleanCmd := strings.TrimSpace(strings.ReplaceAll(cmd, "(use --help to expand)", ""))
		firstWord := strings.Split(cmd, " ")[0]
		return ExpandableMarkerResult{
			Cmd:    cleanCmd,
			Desc:   desc,
			Marker: fmt.Sprintf("  ▸ subcommands — see `gitmap %s --help`", firstWord),
		}
	}

	return ExpandableMarkerResult{
		Cmd:    cmd,
		Desc:   desc,
		Marker: "",
	}
}

func renderLongHelpRow(cmd, fullDesc string, termWidth int) {
	fmt.Printf("  %s\n", cmd)
	indent := "    "
	descWidth := termWidth - len(indent)
	if descWidth <= 10 {
		fmt.Printf("%s%s\n", indent, fullDesc)
		return
	}

	printWrappedHelpLines(indent, fullDesc, descWidth)
}

// HelpRowParams encapsulates parameters for rendering a standard 2-column help row.
type HelpRowParams struct {
	Cmd       string
	FullDesc  string
	CmdWidth  int
	TermWidth int
}

func renderStandardHelpRow(params HelpRowParams) {
	pad := maxHelpCmdLen - params.CmdWidth
	if pad < 0 {
		pad = 0
	}

	prefix := fmt.Sprintf("  %s%s  ", params.Cmd, strings.Repeat(" ", pad))
	descWidth := params.TermWidth - lipgloss.Width(prefix)
	if descWidth <= 10 {
		fmt.Printf("%s%s\n", prefix, params.FullDesc)
		return
	}

	printWrappedHelpLines(prefix, params.FullDesc, descWidth)
}

func renderHelpRow(cmd string, fullDesc string, termWidth int) {
	cmdWidth := lipgloss.Width(cmd)
	if cmdWidth > maxHelpCmdLen {
		renderLongHelpRow(cmd, fullDesc, termWidth)
		return
	}

	renderStandardHelpRow(HelpRowParams{
		Cmd:       cmd,
		FullDesc:  fullDesc,
		CmdWidth:  cmdWidth,
		TermWidth: termWidth,
	})
}

func renderLine(line string) {
	if line == "" {
		return
	}

	cmd, desc := extractCmdAndDesc(line)
	markerResult := parseExpandableMarker(cmd, desc)

	if measuringHelp {
		updateMaxHelpCmdLen(markerResult.Cmd)
		return
	}

	termWidth := pterm.GetTerminalWidth()
	if termWidth <= 0 {
		termWidth = 120
	}

	renderHelpRow(markerResult.Cmd, markerResult.Desc+markerResult.Marker, termWidth)
}

func printWrappedHelpLines(prefix, fullDesc string, descWidth int) {
	wrapped := wrapText(fullDesc, descWidth)
	lines := strings.Split(wrapped, "\n")
	indent := strings.Repeat(" ", lipgloss.Width(prefix))
	for i, l := range lines {
		if i == 0 {
			fmt.Printf("%s%s\n", prefix, l)
			continue
		}
		fmt.Printf("%s%s\n", indent, l)
	}
}

func appendWord(out *strings.Builder, w string, wLen int, curLen int, width int) int {
	if curLen+1+wLen > width {
		out.WriteString("\n")
		out.WriteString(w)
		return wLen
	}

	out.WriteString(" ")
	out.WriteString(w)
	return curLen + 1 + wLen
}

func wrapText(text string, width int) string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return ""
	}

	var out strings.Builder
	curLen := 0
	for i, w := range words {
		wLen := lipgloss.Width(w)
		if i == 0 {
			out.WriteString(w)
			curLen = wLen
			continue
		}

		curLen = appendWord(&out, w, wLen, curLen, width)
	}

	return out.String()
}
