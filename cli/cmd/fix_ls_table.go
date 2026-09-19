// Package cmd — fix_ls_table.go renders aligned tables for git repositories requiring remediation.
package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/termtable"
)

func renderFixLsTable(items []RemediationItem) {
	fmt.Printf("\n%s Repositories Requiring Fix / Remediation (%d):\n\n",
		constants.ColorCyan+"ℹ"+constants.ColorReset, len(items))

	cfg := buildFixLsTableConfig(items)
	termtable.PrintTable(cfg)
	printFixLsFooter()
}

func buildFixLsTableConfig(items []RemediationItem) termtable.TableConfig {
	return termtable.TableConfig{
		HeaderColor:  constants.ColorCyan,
		BorderColor:  constants.ColorDim,
		EllipsisText: "...",
		Columns:      buildFixLsColumns(),
		Rows:         buildFixLsRows(items),
	}
}

func buildFixLsColumns() []termtable.Column {
	return []termtable.Column{
		{Title: "#", MinWidth: 3, Align: termtable.AlignRight},
		{Title: "Repository", MinWidth: 16, MaxWidth: 26, Align: termtable.AlignLeft},
		{Title: "Status / Issues", MinWidth: 22, MaxWidth: 38, Align: termtable.AlignLeft},
		{Title: "Branch", MinWidth: 10, MaxWidth: 16, Align: termtable.AlignLeft},
		{Title: "Suggested Fix", MinWidth: 22, MaxWidth: 36, Align: termtable.AlignLeft},
	}
}

func buildFixLsRows(items []RemediationItem) []termtable.Row {
	rows := make([]termtable.Row, 0, len(items))
	for i, item := range items {
		rows = append(rows, buildSingleFixLsRow(i+1, item))
	}

	return rows
}

func buildSingleFixLsRow(idx int, item RemediationItem) termtable.Row {
	branch, _ := gitutil.CurrentBranch(item.RepoPath)
	if branch == "" {
		branch = "-"
	}

	return termtable.Row{
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

func resolveSuggestedFixCmd(item RemediationItem) string {
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
