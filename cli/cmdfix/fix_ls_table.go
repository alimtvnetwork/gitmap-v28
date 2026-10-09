// Package cmd — fix_ls_table.go renders aligned tables for git repositories requiring remediation.
package cmdfix

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdremediation"
	"fmt"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

func renderFixLsTable(items []cmdremediation.RemediationItem) {
	fmt.Printf("\n%s Repositories Requiring Fix / Remediation (%d):\n\n",
		constants.ColorCyan+"ℹ"+constants.ColorReset, len(items))

	cfg := buildFixLsTableConfig(items)
	termout.PrintTable(cfg)
	printFixLsFooter()
}

func buildFixLsTableConfig(items []cmdremediation.RemediationItem) termout.TableConfig {
	return termout.TableConfig{
		HeaderColor:  constants.ColorCyan,
		BorderColor:  constants.ColorDim,
		EllipsisText: "...",
		Columns:      buildFixLsColumns(),
		Rows:         buildFixLsRows(items),
	}
}

func buildFixLsColumns() []termout.Column {
	return []termout.Column{
		{Title: "#", MinWidth: 3, Align: termout.AlignRight},
		{Title: "Repository", MinWidth: 16, MaxWidth: 26, Align: termout.AlignLeft},
		{Title: "Status / Issues", MinWidth: 22, MaxWidth: 38, Align: termout.AlignLeft},
		{Title: "Branch", MinWidth: 10, MaxWidth: 16, Align: termout.AlignLeft},
		{Title: "Suggested Fix", MinWidth: 22, MaxWidth: 36, Align: termout.AlignLeft},
	}
}

func buildFixLsRows(items []cmdremediation.RemediationItem) []termout.Row {
	rows := make([]termout.Row, 0, len(items))
	for i, item := range items {
		rows = append(rows, buildSingleFixLsRow(i+1, item))
	}

	return rows
}

func buildSingleFixLsRow(idx int, item cmdremediation.RemediationItem) termout.Row {
	branch, _ := gitutil.CurrentBranch(item.RepoPath)
	if branch == "" {
		branch = "-"
	}

	return termout.Row{
		Cells: []string{
			strconv.Itoa(idx),
			item.RepoName,
			item.SummaryReason,
			branch,
			resolveSuggestedFixCmd(item),
		},
		Color: resolveIssueRowColor(item.SummaryReason),
	}
}

func resolveSuggestedFixCmd(item cmdremediation.RemediationItem) string {
	if strings.Contains(item.SummaryReason, "lock") {
		return "gitmap fix-git --locks"
	}

	return fmt.Sprintf("gitmap fix %s 1", item.RepoName)
}

func resolveIssueRowColor(reason string) string {
	low := strings.ToLower(reason)
	if strings.Contains(low, "conflict") || strings.Contains(low, "lock") {
		return constants.ColorRed
	}
	if strings.Contains(low, "dirty") || strings.Contains(low, "uncommitted") {
		return constants.ColorYellow
	}

	return constants.ColorCyan
}

func printFixLsFooter() {
	fmt.Println()
	fmt.Println("  Remediation Commands:")
	fmt.Printf("    %-32s %s\n", constants.ColorCyan+"gitmap fix <repo> [1|2|3]"+constants.ColorReset, "(1=stash, 2=wip, 3=discard)")
	fmt.Printf("    %-32s %s\n", constants.ColorCyan+"gitmap fix all"+constants.ColorReset, "(apply stash to all repositories)")
	fmt.Printf("    %-32s %s\n", constants.ColorCyan+"gitmap fix all [1|2|3]"+constants.ColorReset, "(apply specific strategy to all)")
	fmt.Printf("    %-32s %s\n\n", constants.ColorCyan+"gitmap fix --prompt"+constants.ColorReset, "(step-by-step interactive walkthrough)")
}

func printCleanFixLsOutput(totalScanned int) {
	fmt.Printf("\n%s All repositories are clean! No git issues detected across %d repository(ies).\n\n",
		constants.ColorGreen+"✨"+constants.ColorReset, totalScanned)
}
