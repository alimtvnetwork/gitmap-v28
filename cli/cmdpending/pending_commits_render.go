package cmdpending

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
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
	fmt.Printf("  %s│%s%s%s│%s\n", border, reset, termout.PadRightMax(summaryText, 78), border, reset)
	fmt.Printf("  %s├──────────────────────────────────────────────────────────────────────────────┤%s\n", border, reset)
}

func renderTableHeader(border, reset string) {
	headerText := fmt.Sprintf(" %-22s  %-16s  %12s  %10s  %-9s",
		"REPOSITORY", "VER/BRANCH", "UNCOMMITTED", "UNPUSHED", "STATUS")
	fmt.Printf("  %s│%s%s%s│%s\n", border, reset, termout.PadRightMax(headerText, 78), border, reset)
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
		termout.PadRightMax(repoDisplay, 22), termout.PadRightMax(branchDisplay, 16),
		rec.TotalUncommitted, rec.UnpushedCommitsCount,
		statusText, border, reset)
}

func resolveStatusLabel(rec RepoPendingCommitRecord, reset string) string {
	if rec.IsDirty || rec.HasUnpushed {
		return termout.TokyoWarnPill("PEND") + "   "
	}

	return termout.TokyoOKPill("CLEAN") + "  "
}

func renderTreeRemediationHints(rec RepoPendingCommitRecord, border, reset string) {
	totalOpts := len(rec.RemediationOptions)
	for idx, opt := range rec.RemediationOptions {
		isLastOption := idx == totalOpts-1
		renderSingleRemediationSubtree(opt, isLastOption, border, reset)
	}
}

func renderSingleRemediationSubtree(opt RemediationOption, isLastOption bool, border, reset string) {
	renderRemediationTitleBranch(opt, isLastOption, border, reset)
	renderRemediationChildBranch(opt, isLastOption, border, reset)
}

func renderRemediationTitleBranch(opt RemediationOption, isLastOption bool, border, reset string) {
	prefix := "├──"
	if isLastOption {
		prefix = "└──"
	}
	rawTitle := fmt.Sprintf("   %s Option %d: %s", prefix, opt.OptionNumber, opt.Label)
	formattedTitle := formatTreeLine(rawTitle)
	fmt.Printf("  %s│%s%s%s│%s\n", border, reset, formattedTitle, border, reset)
}

func renderRemediationChildBranch(opt RemediationOption, isLastOption bool, border, reset string) {
	childPrefix := "   │   └── "
	if isLastOption {
		childPrefix = "       └── "
	}
	rawChild := childPrefix + opt.Command
	formattedChild := formatTreeLine(rawChild)
	fmt.Printf("  %s│%s%s%s│%s\n", border, reset, formattedChild, border, reset)
}

func formatTreeLine(raw string) string {
	truncated := termout.TruncateRightWithEllipsis(raw, 78)

	return termout.PadRightMax(truncated, 78)
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
	cleanStatus := termout.TokyoOKPill("CLEAN") + "  "
	fmt.Printf("  %s│%s %-22s  %-16s  %12d  %10d  %s%s│%s\n",
		border, reset,
		termout.PadRightMax(cleanLabel, 22), termout.PadRightMax("-", 16),
		0, 0,
		cleanStatus, border, reset)
}

func renderTableFooter(border, reset string) {
	fmt.Printf("  %s├──────────────────────────────────────────────────────────────────────────────┤%s\n", border, reset)
	footerMsg := ` Fleet Remediation: gitmap sends cp all "wip: save changes" (or gitmap cpar)`
	fmt.Printf("  %s│%s%s%s│%s\n", border, reset, termout.PadRightMax(footerMsg, 78), border, reset)
	fmt.Printf("  %s└──────────────────────────────────────────────────────────────────────────────┘%s\n", border, reset)
}
