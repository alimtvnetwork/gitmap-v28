package cmdpipeline

import (
	"fmt"
	"strings"
)

func formatCombinedSectionFailures(sections []SectionFailure) string {
	if len(sections) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("● Combined Pipeline Section Failures [%d failed section(s)]:\n", len(sections)))
	for i, sec := range sections {
		formatSingleSectionFailure(&sb, sec, i+1, len(sections))
	}

	return strings.TrimRight(sb.String(), "\n")
}

func formatSingleSectionFailure(sb *strings.Builder, sec SectionFailure, idx, total int) {
	sb.WriteString(fmt.Sprintf("\n  ┌─ Section [%d/%d]: %s #%d ➔ Job: %s | Step: %s\n",
		idx, total, sec.WorkflowName, sec.RunId, sec.JobName, sec.StepName))
	formatSectionMetadata(sb, sec)
	formatSectionErrors(sb, sec)
	sb.WriteString("  └──────────────────────────────────────────────────────────\n")
}

func formatSectionMetadata(sb *strings.Builder, sec SectionFailure) {
	if len(sec.CreatedAt) > 0 {
		sb.WriteString(fmt.Sprintf("  │ When Run:  %s\n", formatRunTimestamp(sec.CreatedAt)))
	}

	if len(sec.SavedLogFile) > 0 {
		sb.WriteString(fmt.Sprintf("  │ Saved Log: %s\n", toRelativeGitPath(sec.SavedLogFile)))
	}

	formatExtractedScriptAndLocation(sb, sec.ErrorLines)

	if len(sec.FailureSummary) > 0 {
		sb.WriteString(fmt.Sprintf("  │ Summary:   %s\n", sec.FailureSummary))
	}
}

func formatExtractedScriptAndLocation(sb *strings.Builder, lines []string) {
	script := findPrefixedLine(lines, "Script:")
	if len(script) > 0 {
		sb.WriteString(fmt.Sprintf("  │ %s\n", script))
	}

	file := findPrefixedLine(lines, "File:")
	if len(file) > 0 {
		sb.WriteString(fmt.Sprintf("  │ %s\n", file))
	}
}

func findPrefixedLine(lines []string, prefix string) string {
	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		if strings.HasPrefix(trimmed, "│ "+prefix) {
			return strings.TrimPrefix(trimmed, "│ ")
		}
		if strings.HasPrefix(trimmed, prefix) {
			return trimmed
		}
	}

	return ""
}

func formatSectionErrors(sb *strings.Builder, sec SectionFailure) {
	for _, line := range sec.ErrorLines {
		sb.WriteString(fmt.Sprintf("  │   %s\n", line))
	}
	formatSectionStackLines(sb, sec.StackTrace)
}

func formatSectionStackLines(sb *strings.Builder, stack string) {
	if len(stack) == 0 {
		return
	}

	sb.WriteString("  │   Stack Trace:\n")
	for _, line := range strings.Split(strings.TrimSpace(stack), "\n") {
		sb.WriteString(fmt.Sprintf("  │     %s\n", line))
	}
}

// formatAggregatedErrorLogs generates a comprehensive human-readable summary of all failed runs.
func formatAggregatedErrorLogs(failedRuns []FailedRunItem) string {
	if len(failedRuns) == 0 {
		return ""
	}

	sections := extractAllSectionFailures(failedRuns)
	combinedText := formatCombinedSectionFailures(sections)
	detailedText := formatAllRunsDetailed(failedRuns)

	return combinedText + "\n\n" + detailedText
}

func formatAllRunsDetailed(failedRuns []FailedRunItem) string {
	var sb strings.Builder
	sb.WriteString("================================================================================\n")
	sb.WriteString("DETAILED PIPELINE RUN BREAKDOWNS\n")
	sb.WriteString("================================================================================\n\n")
	for i, run := range failedRuns {
		if i > 0 {
			sb.WriteString("\n\n")
		}

		formatSingleRunDetailed(&sb, run)
	}

	return sb.String()
}

func formatSingleRunDetailed(sb *strings.Builder, run FailedRunItem) {
	sb.WriteString(fmt.Sprintf("==> Failed Run: %s (#%d)\n", run.WorkflowName, run.RunId))
	formatRunTimingMeta(sb, run)
	formatRunSourceMeta(sb, run)
	for _, job := range run.FailedJobs {
		formatJobLines(sb, job)
	}
}

func formatRunTimingMeta(sb *strings.Builder, run FailedRunItem) {
	if len(run.CreatedAt) > 0 {
		sb.WriteString(fmt.Sprintf("    When Run:  %s\n", formatRunTimestamp(run.CreatedAt)))
	}

	if run.DurationSeconds > 0 {
		sb.WriteString(fmt.Sprintf("    Duration:  %s\n", formatDurationSeconds(run.DurationSeconds)))
	}
}

func formatRunSourceMeta(sb *strings.Builder, run FailedRunItem) {
	if len(run.Branch) > 0 {
		sb.WriteString(fmt.Sprintf("    Branch:    %s | Commit: %s\n", run.Branch, run.Sha))
	}

	if len(run.SavedLogFile) > 0 {
		sb.WriteString(fmt.Sprintf("    Saved Log: %s\n", toRelativeGitPath(run.SavedLogFile)))
	}

	if len(run.Url) > 0 {
		sb.WriteString(fmt.Sprintf("    URL:       %s\n", run.Url))
	}
}

func formatJobLines(sb *strings.Builder, job FailedJobItem) {
	sb.WriteString(fmt.Sprintf("    ● Job: %s | Step: %s\n", job.JobName, job.StepName))
	if len(job.Warnings) > 0 {
		sb.WriteString(fmt.Sprintf("      Warnings: (%d detected)\n", len(job.Warnings)))
		for _, w := range capErrorLines(job.Warnings, 5) {
			sb.WriteString(fmt.Sprintf("        %s\n", w))
		}
	}
	if len(job.FailureSummary) > 0 {
		sb.WriteString(fmt.Sprintf("      Summary: %s\n", job.FailureSummary))
	}

	for _, l := range filterOutSummaryLine(job.ErrorLines, job.FailureSummary) {
		sb.WriteString(fmt.Sprintf("      %s\n", l))
	}
	formatJobStackLines(sb, job.StackTrace)
}

func formatJobStackLines(sb *strings.Builder, stack string) {
	if len(stack) == 0 {
		return
	}

	sb.WriteString("      Stack Trace:\n")
	for _, line := range strings.Split(strings.TrimSpace(stack), "\n") {
		sb.WriteString(fmt.Sprintf("        %s\n", line))
	}
}
