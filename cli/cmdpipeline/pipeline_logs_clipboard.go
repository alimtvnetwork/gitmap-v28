package cmdpipeline

import (
	"fmt"
	"strings"
)

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
