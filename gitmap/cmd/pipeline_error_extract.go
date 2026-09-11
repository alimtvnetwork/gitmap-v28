package cmd

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
)

var ansiRegex = regexp.MustCompile(`(?:\x1b|\^\[)\[[0-9;]*[a-zA-Z]`)

var failureMarkers = []string{
	"##[error]",
	"--- FAIL:",
	"FAIL\t",
	"FAIL:",
	"FAILED",
	"FAIL ",
	"Expected",
	"fatal error:",
	"syntax error:",
	"exit status",
	"exit code",
	"exited ",
	"panic:",
	"Stack Trace:",
	"Unexpected token",
	"Not Found - ",
	"Error:",
	"error:",
	"[gosec-",
	"gofmt",
	"FAILED TESTS",
	"diff (baseline-diff",
	"Suite:",
	"got \"",
	"want \"",
	"lint error",
	"Process completed with exit code",
}

// extractCleanErrorLines filters noisy logs to isolate only failure and error lines.
func extractCleanErrorLines(rawLogs string) string {
	jobs := ParseFailedLogLines(rawLogs)
	if len(jobs) == 0 {
		return ""
	}

	var parts []string
	for _, j := range jobs {
		lines := append([]string{j.FailureSummary}, j.ErrorLines...)
		cleanLines := deduplicateLines(lines)
		parts = append(parts, strings.Join(cleanLines, "\n"))
	}

	return strings.Join(parts, "\n\n")
}

func deduplicateLines(lines []string) []string {
	var out []string
	seen := make(map[string]bool)
	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		if len(trimmed) > 0 && !seen[trimmed] {
			seen[trimmed] = true
			out = append(out, l)
		}
	}

	return out
}

// ParseFailedLogLines parses raw GitHub Actions log text into structured failed jobs.
func ParseFailedLogLines(rawLogs string) []FailedJobItem {
	if strings.TrimSpace(rawLogs) == "" {
		return nil
	}

	jobMap := make(map[string]*FailedJobItem)
	var order []string
	scanner := bufio.NewScanner(strings.NewReader(rawLogs))
	contextRemaining := 0
	var lastKey string

	for scanner.Scan() {
		raw := scanner.Text()
		job, step, text, isError := parseLogLine(raw)
		if text == "" || isIgnoredLogLine(text) {
			continue
		}
		key := job + "|||" + step
		if isError {
			item := getOrCreateJobItem(jobMap, &order, key, job, step)
			item.ErrorLines = append(item.ErrorLines, text)
			updateJobSummary(item, text)
			contextRemaining = 25
			lastKey = key
			continue
		}
		if contextRemaining > 0 && lastKey == key {
			item := jobMap[key]
			item.ErrorLines = append(item.ErrorLines, "    "+text)
			contextRemaining--
		}
	}

	return assembleJobItems(jobMap, order, rawLogs)
}

func parseLogLine(raw string) (string, string, string, bool) {
	clean := ansiRegex.ReplaceAllString(raw, "")
	parts := strings.Split(clean, "\t")
	isError := strings.Contains(raw, "##[error]") || hasFailureMarker(raw)

	if len(parts) >= 3 {
		job := strings.TrimSpace(parts[0])
		step := strings.TrimSpace(parts[1])
		text := stripTimestamp(strings.Join(parts[2:], "\t"))
		return job, step, cleanLogText(text), isError || hasFailureMarker(text)
	}
	if len(parts) == 2 {
		return strings.TrimSpace(parts[0]), "", cleanLogText(parts[1]), isError
	}

	return "", "", cleanLogText(clean), isError
}

func stripTimestamp(text string) string {
	if idx := strings.Index(text, "Z "); idx != -1 {
		return text[idx+2:]
	}
	if idx := strings.Index(text, "Z\t"); idx != -1 {
		return text[idx+2:]
	}

	return text
}

func cleanLogText(text string) string {
	t := strings.TrimPrefix(text, "##[error]")
	t = strings.TrimPrefix(t, "##[group]")
	t = strings.TrimPrefix(t, "##[endgroup]")

	return strings.TrimSpace(t)
}

func isIgnoredLogLine(text string) bool {
	lower := strings.ToLower(text)
	if strings.Contains(lower, "post job cleanup") {
		return true
	}
	if strings.Contains(lower, "safe.directory") {
		return true
	}
	if strings.Contains(lower, "removing ssh command") {
		return true
	}

	return strings.Contains(lower, "terminate orphan process")
}

func hasFailureMarker(line string) bool {
	for _, m := range failureMarkers {
		if strings.Contains(line, m) {
			return true
		}
	}

	return false
}

func getOrCreateJobItem(jobMap map[string]*FailedJobItem, order *[]string, key, job, step string) *FailedJobItem {
	if item, exists := jobMap[key]; exists {
		return item
	}
	if step == "UNKNOWN STEP" || step == "" {
		step = "Job Execution"
	}
	item := &FailedJobItem{
		JobName:  job,
		StepName: step,
	}
	jobMap[key] = item
	*order = append(*order, key)

	return item
}

func updateJobSummary(item *FailedJobItem, text string) {
	if len(item.FailureSummary) == 0 || isStrongerSummary(text, item.FailureSummary) {
		item.FailureSummary = text
	}
}

func isStrongerSummary(candidate, current string) bool {
	if len(current) == 0 {
		return true
	}
	if isGenericExitCode(candidate) && !isGenericExitCode(current) {
		return false
	}
	if strings.Contains(candidate, "Expected ") {
		return true
	}
	if strings.Contains(candidate, "gofmt") {
		return true
	}
	if strings.Contains(candidate, "fatal error:") {
		return true
	}

	return isGenericExitCode(current)
}

func isGenericExitCode(s string) bool {
	return strings.Contains(s, "Process completed with exit code") || strings.Contains(s, "exit status 1")
}

func assembleJobItems(jobMap map[string]*FailedJobItem, order []string, rawLogs string) []FailedJobItem {
	if len(order) == 0 {
		return buildFallbackJobItems(rawLogs)
	}

	var results []FailedJobItem
	for _, k := range order {
		results = append(results, *jobMap[k])
	}

	return results
}

func buildFallbackJobItems(rawLogs string) []FailedJobItem {
	tail := extractTailLines(rawLogs, 50)
	if tail == "" {
		return nil
	}

	return []FailedJobItem{
		{
			JobName:        "Pipeline Execution",
			StepName:       "Failed Steps",
			FailureSummary: "Step execution failed (see error lines)",
			ErrorLines:     strings.Split(tail, "\n"),
		},
	}
}

func extractTailLines(rawLogs string, n int) string {
	lines := strings.Split(strings.TrimSpace(rawLogs), "\n")
	if len(lines) <= n {
		return strings.Join(lines, "\n")
	}

	return strings.Join(lines[len(lines)-n:], "\n")
}

func extractAllSectionFailures(failedRuns []FailedRunItem) []SectionFailure {
	var sections []SectionFailure
	for _, run := range failedRuns {
		sections = append(sections, extractRunSectionFailures(run)...)
	}

	return sections
}

func extractRunSectionFailures(run FailedRunItem) []SectionFailure {
	var out []SectionFailure
	for _, job := range run.FailedJobs {
		sf := SectionFailure{
			WorkflowName: run.WorkflowName, RunId: run.RunId,
			JobName: job.JobName, StepName: job.StepName,
			FailureSummary: job.FailureSummary, ErrorLines: job.ErrorLines,
			SavedLogFile: run.SavedLogFile, CreatedAt: run.CreatedAt,
		}
		out = append(out, sf)
	}

	return out
}

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
	if len(sec.CreatedAt) > 0 {
		sb.WriteString(fmt.Sprintf("  │ When Run:  %s\n", formatRunTimestamp(sec.CreatedAt)))
	}
	if len(sec.SavedLogFile) > 0 {
		sb.WriteString(fmt.Sprintf("  │ Saved Log: %s\n", sec.SavedLogFile))
	}
	if len(sec.FailureSummary) > 0 {
		sb.WriteString(fmt.Sprintf("  │ Summary:   %s\n", sec.FailureSummary))
	}
	for _, line := range sec.ErrorLines {
		sb.WriteString(fmt.Sprintf("  │   %s\n", line))
	}
	sb.WriteString("  └──────────────────────────────────────────────────────────\n")
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
	formatRunMetaLines(sb, run)
	for _, job := range run.FailedJobs {
		formatJobLines(sb, job)
	}
}

func formatRunMetaLines(sb *strings.Builder, run FailedRunItem) {
	if len(run.CreatedAt) > 0 {
		sb.WriteString(fmt.Sprintf("    When Run:  %s\n", formatRunTimestamp(run.CreatedAt)))
	}
	if run.DurationSeconds > 0 {
		sb.WriteString(fmt.Sprintf("    Duration:  %s\n", formatDurationSeconds(run.DurationSeconds)))
	}
	if len(run.Branch) > 0 {
		sb.WriteString(fmt.Sprintf("    Branch:    %s | Commit: %s\n", run.Branch, run.Sha))
	}
	if len(run.SavedLogFile) > 0 {
		sb.WriteString(fmt.Sprintf("    Saved Log: %s\n", run.SavedLogFile))
	}
	if len(run.Url) > 0 {
		sb.WriteString(fmt.Sprintf("    URL:       %s\n", run.Url))
	}
}

func formatJobLines(sb *strings.Builder, job FailedJobItem) {
	sb.WriteString(fmt.Sprintf("    ● Job: %s | Step: %s\n", job.JobName, job.StepName))
	if len(job.FailureSummary) > 0 {
		sb.WriteString(fmt.Sprintf("      Summary: %s\n", job.FailureSummary))
	}
	for _, l := range job.ErrorLines {
		sb.WriteString(fmt.Sprintf("      %s\n", l))
	}
}
