package cmdpending

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func renderTableBanner(border, reset string) {
	bold := constants.ColorBold
	fmt.Println()
	fmt.Printf("  %s┌──────────────────────────────────────────────────────────────────────────────┐%s\n", border, reset)
	fmt.Printf("  %s│%s%s                    GITMAP PENDING COMMITS SUMMARY                            %s%s│%s\n", border, reset, bold, reset, border, reset)
	fmt.Printf("  %s├──────────────────────────────────────────────────────────────────────────────┤%s\n", border, reset)
}

func renderTableSummary(payload PendingCommitsPayload, border, reset string) {
	summaryText := fmt.Sprintf(" Scanned: %d repos │ Dirty: %d repos │ Uncommitted: %d files │ Unpushed: %d",
		payload.TotalReposScanned, payload.TotalDirtyRepos, payload.TotalUncommittedFiles, payload.TotalUnpushedCommits)
	fmt.Printf("  %s│%s%s%s│%s\n", border, reset, padRight(summaryText, 78), border, reset)
	fmt.Printf("  %s├──────────────────────────────────────────────────────────────────────────────┤%s\n", border, reset)
}

func renderTableHeader(border, reset string) {
	headerText := fmt.Sprintf(" %-22s  %-16s  %12s  %10s  %-9s",
		"REPOSITORY", "VER/BRANCH", "UNCOMMITTED", "UNPUSHED", "STATUS")
	fmt.Printf("  %s│%s%s%s│%s\n", border, reset, padRight(headerText, 78), border, reset)
	fmt.Printf("  %s├──────────────────────────────────────────────────────────────────────────────┤%s\n", border, reset)
}

func renderTableRows(repos []RepoPendingCommitRecord, border, reset string) {
	for _, rec := range repos {
		renderPendingRepoRow(rec, border, reset)
		if rec.IsDirty {
			renderTreeRemediationHints(rec, border, reset)
		}
	}
}

func renderPendingRepoRow(rec RepoPendingCommitRecord, border, reset string) {
	repoDisplay := truncateStringWithEllipsis(rec.RepoName, 22)
	branchDisplay := rec.ShortVersionBranch
	if branchDisplay == "" {
		branchDisplay = rec.CurrentBranch
	}
	branchDisplay = truncateStringWithEllipsis(branchDisplay, 16)
	statusText := resolveStatusLabel(rec, reset)
	fmt.Printf("  %s│%s %-22s  %-16s  %12d  %10d  %s%s│%s\n",
		border, reset,
		padRight(repoDisplay, 22), padRight(branchDisplay, 16),
		rec.TotalUncommitted, rec.UnpushedCommitsCount,
		statusText, border, reset)
}

func resolveStatusLabel(rec RepoPendingCommitRecord, reset string) string {
	if rec.IsDirty || rec.HasUnpushed {
		return constants.ColorYellow + "● PEND   " + reset
	}
	return constants.ColorGreen + "○ CLEAN  " + reset
}

func renderTreeRemediationHints(rec RepoPendingCommitRecord, border, reset string) {
	totalOpts := len(rec.RemediationOptions)
	for idx, opt := range rec.RemediationOptions {
		prefix := "├──"
		if idx == totalOpts-1 {
			prefix = "└──"
		}
		rawLine := fmt.Sprintf("   %s Option %d: %s", prefix, opt.OptionNumber, opt.Command)
		formattedLine := formatTreeLine(rawLine)
		fmt.Printf("  %s│%s%s%s│%s\n", border, reset, formattedLine, border, reset)
	}
}

func formatTreeLine(raw string) string {
	runes := []rune(raw)
	const boxWidth = 78
	if len(runes) > boxWidth {
		return string(runes[:boxWidth-1]) + "…"
	}
	return raw + strings.Repeat(" ", boxWidth-len(runes))
}

func renderCleanSummaryIfNeeded(payload PendingCommitsPayload, totalScanned int, allInspected []RepoPendingCommitRecord, border, reset string) {
	cleanCount := countCleanRepositories(allInspected)
	if payload.DetailMode != "summary" || cleanCount <= 0 || len(payload.Repositories) >= totalScanned {
		return
	}
	renderCleanSummaryRow(cleanCount, border, reset)
}

func countCleanRepositories(records []RepoPendingCommitRecord) int {
	count := 0
	for _, rec := range records {
		if rec.IsClean {
			count++
		}
	}
	return count
}

func renderCleanSummaryRow(cleanCount int, border, reset string) {
	cleanLabel := fmt.Sprintf("%d clean repos", cleanCount)
	cleanStatus := constants.ColorGreen + "○ CLEAN  " + reset
	fmt.Printf("  %s│%s %-22s  %-16s  %12d  %10d  %s%s│%s\n",
		border, reset,
		padRight(cleanLabel, 22), padRight("-", 16),
		0, 0,
		cleanStatus, border, reset)
}

func renderTableFooter(border, reset string) {
	fmt.Printf("  %s├──────────────────────────────────────────────────────────────────────────────┤%s\n", border, reset)
	footerMsg := ` Fleet Remediation: gitmap cpar "wip: save changes"`
	fmt.Printf("  %s│%s%s%s│%s\n", border, reset, padRight(footerMsg, 78), border, reset)
	fmt.Printf("  %s└──────────────────────────────────────────────────────────────────────────────┘%s\n", border, reset)
}

func padRight(s string, width int) string {
	rc := utf8.RuneCountInString(s)
	if rc >= width {
		return s
	}
	return s + strings.Repeat(" ", width-rc)
}
