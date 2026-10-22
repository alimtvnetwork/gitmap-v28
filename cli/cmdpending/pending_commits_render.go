package cmdpending

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

func renderTableBanner(border, reset string) {
	fmt.Println()
	topBorder := strings.Repeat("─", 78)
	fmt.Printf("  %s╭%s╮%s\n", border, topBorder, reset)
	styledTitle := termout.SafeTokyoStyle(termout.TokyoPrimary).Bold(true).Render("GITMAP PENDING COMMITS SUMMARY")
	paddedTitle := centerVisual(styledTitle, 78)
	fmt.Printf("  %s│%s%s%s│%s\n", border, reset, paddedTitle, border, reset)
	fmt.Printf("  %s├%s┤%s\n", border, topBorder, reset)
}

func centerVisual(s string, target int) string {
	vis := termout.VisualStringWidth(s)
	if vis >= target {
		return s
	}
	pad := (target - vis) / 2
	rightPad := target - vis - pad

	return termout.SpacePool(pad) + s + termout.SpacePool(rightPad)
}

func renderTableSummary(payload PendingCommitsPayload, border, reset string) {
	summaryText := formatSummaryMetrics(payload)
	fmt.Printf("  %s│%s%s%s│%s\n", border, reset, padRightVisual(" "+summaryText, 78), border, reset)
	midBorder := strings.Repeat("─", 78)
	fmt.Printf("  %s├%s┤%s\n", border, midBorder, reset)
}

func formatSummaryMetrics(payload PendingCommitsPayload) string {
	pills := []string{
		termout.TokyoMetricPill("Scanned", payload.TotalReposScanned),
		termout.TokyoMetricPill("Dirty", payload.TotalDirtyRepos),
		termout.TokyoMetricPill("Uncommitted", payload.TotalUncommittedFiles),
		termout.TokyoMetricPill("Unpushed", payload.TotalUnpushedCommits),
	}

	return strings.Join(pills, "  ")
}

func padRightVisual(s string, target int) string {
	vis := termout.VisualStringWidth(s)
	if vis >= target {
		return s
	}

	return s + termout.SpacePool(target-vis)
}

func renderTableHeader(border, reset string) {
	headerText := fmt.Sprintf(" %-22s  %-16s  %12s  %10s  %-9s",
		"REPOSITORY", "VER/BRANCH", "UNCOMMITTED", "UNPUSHED", "STATUS")
	midBorder := strings.Repeat("─", 78)
	fmt.Printf("  %s│%s%s%s│%s\n", border, reset, padRightVisual(headerText, 78), border, reset)
	fmt.Printf("  %s├%s┤%s\n", border, midBorder, reset)
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
	branchGlyph := resolveTreeBranchGlyph(isLastOption)
	optLabel := fmt.Sprintf("Option %d: %s", opt.OptionNumber, opt.Label)
	titleText := formatRemediationTitle(branchGlyph, optLabel)
	formattedTitle := padRightVisual("   "+titleText, 78)
	fmt.Printf("  %s│%s%s%s│%s\n", border, reset, formattedTitle, border, reset)
}

func resolveTreeBranchGlyph(isLastOption bool) string {
	branch := "├──"
	if isLastOption {
		branch = "└──"
	}

	return termout.SafeTokyoStyle(termout.TokyoMuted).Render(branch)
}

func formatRemediationTitle(branch, label string) string {
	styledLabel := termout.SafeTokyoStyle(termout.TokyoPrimary).Render(label)

	return fmt.Sprintf("%s %s", branch, styledLabel)
}

func renderRemediationChildBranch(opt RemediationOption, isLastOption bool, border, reset string) {
	childGlyph := resolveChildBranchGlyph(isLastOption)
	styledCmd := termout.SafeTokyoStyle(termout.TokyoAccent).Render(opt.Command)
	formattedChild := padRightVisual(childGlyph+styledCmd, 78)
	fmt.Printf("  %s│%s%s%s│%s\n", border, reset, formattedChild, border, reset)
}

func resolveChildBranchGlyph(isLastOption bool) string {
	if isLastOption {
		branch := "       └── "
		return termout.SafeTokyoStyle(termout.TokyoMuted).Render(branch)
	}
	branch := "   │   └── "

	return termout.SafeTokyoStyle(termout.TokyoMuted).Render(branch)
}

func formatTreeLine(raw string) string {
	truncated := termout.TruncateRightWithEllipsis(raw, 78)

	return padRightVisual(truncated, 78)
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
	midBorder := strings.Repeat("─", 78)
	botBorder := strings.Repeat("─", 78)
	fmt.Printf("  %s├%s┤%s\n", border, midBorder, reset)
	footerMsg := formatFooterMessage()
	fmt.Printf("  %s│%s%s%s│%s\n", border, reset, padRightVisual(footerMsg, 78), border, reset)
	fmt.Printf("  %s╰%s╯%s\n", border, botBorder, reset)
}

func formatFooterMessage() string {
	prefix := termout.SafeTokyoStyle(termout.TokyoMuted).Render(" Fleet Remediation: ")
	cmd1 := termout.SafeTokyoStyle(termout.TokyoAccent).Render(`gitmap sends cp all "wip: save changes"`)
	cmd2 := termout.SafeTokyoStyle(termout.TokyoPrimary).Render(`gitmap cpar`)

	return fmt.Sprintf("%s%s (or %s)", prefix, cmd1, cmd2)
}
