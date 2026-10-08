package cmdpipeline

import (
	"fmt"
	"path/filepath"
	"strings"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func renderSingleFailureTerminal(p PipelineErrorLogsPayload) {
	fmt.Printf("  %s● Latest Pipeline Failure [%s #%d]:%s\n\n",
		constants.ColorRed, p.WorkflowName, p.RunId, constants.ColorReset)
	renderPayloadSourceDetails(p)
	clean := extractCleanErrorLines(p.ErrorLogs)
	printLogsContent(clean, p.ErrorLogs)
	renderSavedLocationsTerminal(p)
	printRerunETA(p.RerunEtaSeconds)
}

func renderPayloadSourceDetails(p PipelineErrorLogsPayload) {
	logger := "github-actions"
	if len(p.WorkflowName) > 0 {
		logger = p.WorkflowName
	}
	logFile := filepath.ToSlash(FormatRelativeDbPath(p.SavedReportFile))
	if len(logFile) == 0 {
		logFile = "in-memory / active stream"
	}

	fmt.Printf("    Source Details:\n")
	fmt.Printf("      Logger:    %s\n", logger)
	fmt.Printf("      Workflow:  %s\n", p.WorkflowName)
	if len(p.Url) > 0 {
		fmt.Printf("      URL:       %s\n", p.Url)
	}
	fmt.Printf("      Log File:  %s\n\n", logFile)
}

func renderFailedRunCard(fr FailedRunItem, idx, total int) {
	fmt.Printf("  %s┌─ [%d/%d] %s #%d ──────────────────────────%s\n",
		constants.ColorRed, idx, total, fr.WorkflowName, fr.RunId, constants.ColorReset)
	renderRunCardMeta(fr)
	for _, job := range fr.FailedJobs {
		renderFailedJobSection(job)
	}

	fmt.Printf("  %s└──────────────────────────────────────────────────────────%s\n\n",
		constants.ColorRed, constants.ColorReset)
}

func renderRunCardMeta(fr FailedRunItem) {
	if len(fr.CreatedAt) > 0 {
		fmt.Printf("  │ When Run:  %s\n", formatRunTimestamp(fr.CreatedAt))
	}
	if fr.DurationSeconds > 0 {
		fmt.Printf("  │ Duration:  %s\n", formatDurationSeconds(fr.DurationSeconds))
	}
	renderRunCardBranchAndLog(fr)
}

func renderRunCardBranchAndLog(fr FailedRunItem) {
	if len(fr.Branch) > 0 {
		fmt.Printf("  │ Branch:    %s | Commit: %s\n", fr.Branch, fr.Sha)
	}
	if len(fr.SavedLogFile) > 0 {
		fmt.Printf("  │ Saved Log: %s\n", filepath.ToSlash(FormatRelativeDbPath(fr.SavedLogFile)))
	}
	if len(fr.Url) > 0 {
		fmt.Printf("  │ URL:       %s\n", fr.Url)
	}
}

func renderSavedLocationsTerminal(p PipelineErrorLogsPayload) {
	fmt.Println("  💾 Pipeline Error Logs & Artifacts:")
	if len(p.SavedReportFile) > 0 {
		fmt.Printf("    • Combined Report: %s\n", filepath.ToSlash(FormatRelativeDbPath(p.SavedReportFile)))
	}
	if len(p.SavedLogFile) > 0 {
		fmt.Printf("    • Latest Run Log:  %s\n", filepath.ToSlash(FormatRelativeDbPath(p.SavedLogFile)))
	}
	renderSavedDbAndUrl(p)
}

func renderSavedDbAndUrl(p PipelineErrorLogsPayload) {
	if len(p.DbPath) > 0 {
		fmt.Printf("    • Pipeline DB:     %s\n", filepath.ToSlash(FormatRelativeDbPath(p.DbPath)))
		fmt.Printf("    • DB Size:         %s\n", ResolveDbFileSize(p.DbPath))
		fmt.Printf("    • Cleanup:         gitmap pipeline clear -y\n")
	}
	if len(p.Url) > 0 {
		fmt.Printf("    • Web Run URL:     %s\n\n", p.Url)
	}
}

func renderFailedJobSection(job FailedJobItem) {
	fmt.Printf("  │ Job:   %s\n", job.JobName)
	fmt.Printf("  │ Step:  %s\n", job.StepName)
	if len(job.FailureSummary) > 0 {
		cleaned := cleanDisplayErrorText(job.FailureSummary, job.JobName, job.StepName)
		fmt.Printf("  │ Error: %s%s%s\n", constants.ColorRed, cleaned, constants.ColorReset)
	}

	renderJobCardWarnings(job.Warnings)
	renderJobCardDetails(job.ErrorLines, job.FailureSummary)
	renderJobStackTrace(job.StackTrace)
}

func renderJobCardWarnings(warnings []string) {
	if len(warnings) == 0 {
		return
	}

	fmt.Printf("  │ Warnings (%d preceding):\n", len(warnings))
	capped := capErrorLines(warnings, 25)
	for _, w := range capped {
		fmt.Printf("  │   %s%s%s\n", constants.ColorYellow, w, constants.ColorReset)
	}

	printRemainingWarningCount(len(warnings), len(capped))
}

func printRemainingWarningCount(total, capped int) {
	if total > capped {
		fmt.Printf("  │   %s... (%d more warnings)%s\n", constants.ColorDim, total-capped, constants.ColorReset)
	}
}

func renderJobCardDetails(lines []string, summary string) {
	dedupLines := filterOutSummaryLine(lines, summary)
	for _, line := range dedupLines {
		fmt.Printf("  │   %s\n", line)
	}
}

func renderJobStackTrace(stack string) {
	if len(stack) == 0 {
		return
	}

	fmt.Printf("  │ Stack Trace:\n")
	for _, line := range strings.Split(strings.TrimSpace(stack), "\n") {
		fmt.Printf("  │   %s%s%s\n", constants.ColorDim, line, constants.ColorReset)
	}
}

func printLogsContent(clean, raw string) {
	if len(clean) > 0 {
		fmt.Println(clean)

		return
	}

	fmt.Println(raw)
}

func printRerunETA(eta int) {
	if eta <= 0 {
		return
	}

	fmt.Printf("\n  %s● Estimated pipeline rerun duration (ETA): %s%s\n",
		constants.ColorYellow, formatEtaDisplay(eta), constants.ColorReset)
	fmt.Println("    (Based on historical successful pipeline runs baseline)")
}

func buildClipboardErrorReport(p PipelineErrorLogsPayload) string {
	var sb strings.Builder
	appendClipboardMetaHeader(&sb, "GITMAP PIPELINE ERROR REPORT", p)
	appendClipboardFailureTree(&sb, p)
	appendClipboardBodyContent(&sb, p)

	return strings.TrimSpace(sb.String())
}

func appendClipboardBodyContent(sb *strings.Builder, p PipelineErrorLogsPayload) {
	if len(p.ErrorLogs) > 0 {
		sb.WriteString(p.ErrorLogs)
		sb.WriteString("\n")

		return
	}
	if len(p.CombinedErrors) > 0 {
		sb.WriteString(p.CombinedErrors)
		sb.WriteString("\n")
	}
}

func buildClipboardCleanReport(p PipelineErrorLogsPayload) string {
	var sb strings.Builder
	appendClipboardMetaHeader(&sb, "GITMAP PIPELINE STATUS: CLEAN (No errors found)", p)
	sb.WriteString("All recent pipeline workflow runs are PASSING (100% green).\n")

	return strings.TrimSpace(sb.String())
}

func appendClipboardMetaHeader(sb *strings.Builder, title string, p PipelineErrorLogsPayload) {
	sb.WriteString("================================================================================\n")
	sb.WriteString(title + "\n")
	sb.WriteString("================================================================================\n")
	sb.WriteString(fmt.Sprintf("Repo:                 %s\n", p.Repo))
	if len(p.RepoUrl) > 0 {
		sb.WriteString(fmt.Sprintf("Repo URL:             %s\n", p.RepoUrl))
	}
	appendClipboardMetaDetails(sb, p)
}

func appendClipboardMetaDetails(sb *strings.Builder, p PipelineErrorLogsPayload) {
	if p.IsFromCache {
		sb.WriteString(fmt.Sprintf("Cache:                Served from local SQLite DB (commit %s)\n", p.LastHash))
	}
	sb.WriteString(fmt.Sprintf("Branch:               %s\n", p.LatestBranch))
	sb.WriteString(fmt.Sprintf("Last Commit:          %s\n", p.LastHash))
	sb.WriteString(fmt.Sprintf("Last Release:         %s\n", p.LastReleaseVersion))
	sb.WriteString(fmt.Sprintf("Open PRs:             %d\n", p.OpenPRsCount))
	appendClipboardRunStatus(sb, p)
}

func appendClipboardRunStatus(sb *strings.Builder, p PipelineErrorLogsPayload) {
	if len(p.Status) > 0 && len(p.Conclusion) > 0 {
		sb.WriteString(fmt.Sprintf("Status:               %s (conclusion: %s)\n", p.Status, p.Conclusion))
	}
	if len(p.Url) > 0 {
		sb.WriteString(fmt.Sprintf("Pipeline Run URL:     %s\n", p.Url))
	}
	sb.WriteString("================================================================================\n\n")
}

func copyReportToClipboard(content string, isFailure bool) {
	if !isClipboardWriteAllowed(false) {
		return
	}
	if len(content) == 0 {
		return
	}

	err := writeClipboard(content)
	if err != nil {
		return
	}

	printClipboardNotice(isFailure)
}

func printClipboardNotice(isFailure bool) {
	if isFailure {
		fmt.Printf("\n  📋 Copied pipeline error logs to clipboard\n")

		return
	}

	fmt.Printf("\n  📋 Copied pipeline status to clipboard\n")
}
