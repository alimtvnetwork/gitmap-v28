package cmd

import (
	"bufio"
	"fmt"
	"path/filepath"
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
	scanLogLinesIntoMap(rawLogs, jobMap, &order)

	return assembleJobItems(jobMap, order, rawLogs)
}

func scanLogLinesIntoMap(rawLogs string, jobMap map[string]*FailedJobItem, order *[]string) {
	scanner := bufio.NewScanner(strings.NewReader(rawLogs))
	contextRemaining := 0
	var lastKey string
	for scanner.Scan() {
		processLogLine(scanner.Text(), jobMap, order, &contextRemaining, &lastKey)
	}
}

func processLogLine(raw string, jobMap map[string]*FailedJobItem, order *[]string, ctxRem *int, lastKey *string) {
	job, step, text, isError := parseLogLine(raw)
	if text == "" || isIgnoredLogLine(text) {
		return
	}

	key := job + "|||" + step
	if isError {
		recordErrorLine(jobMap, order, key, job, step, text, ctxRem, lastKey)

		return
	}

	appendContextLine(jobMap, key, text, ctxRem, lastKey)
}

func recordErrorLine(jobMap map[string]*FailedJobItem, order *[]string, key, job, step, text string, ctxRem *int, lastKey *string) {
	item := getOrCreateJobItem(jobMap, order, key, job, step)
	item.ErrorLines = append(item.ErrorLines, text)
	updateJobSummary(item, text)
	*ctxRem = 25
	*lastKey = key
}

func appendContextLine(jobMap map[string]*FailedJobItem, key, text string, ctxRem *int, lastKey *string) {
	if *ctxRem > 0 && *lastKey == key {
		item := jobMap[key]
		item.ErrorLines = append(item.ErrorLines, "    "+text)
		*ctxRem--
	}
}

func parseLogLine(raw string) (string, string, string, bool) {
	clean := ansiRegex.ReplaceAllString(raw, "")
	parts := strings.Split(clean, "\t")
	isError := strings.Contains(raw, "##[error]") || hasFailureMarker(raw)
	if len(parts) >= 3 {
		return extractThreePartLogLine(parts, isError)
	}

	if len(parts) == 2 {
		return strings.TrimSpace(parts[0]), "", cleanLogText(parts[1]), isError
	}

	return "", "", cleanLogText(clean), isError
}

func extractThreePartLogLine(parts []string, isError bool) (string, string, string, bool) {
	job := strings.TrimSpace(parts[0])
	step := strings.TrimSpace(parts[1])
	text := stripTimestamp(strings.Join(parts[2:], "\t"))

	return job, step, cleanLogText(text), isError || hasFailureMarker(text)
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
	if strings.Contains(lower, "post job cleanup") || strings.Contains(lower, "safe.directory") {
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

	if strings.Contains(candidate, "Expected ") || strings.Contains(candidate, "gofmt") {
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

// CorrelateFailedJobs merges parsed log items with gh run view jobs so zero failing steps are omitted.
func CorrelateFailedJobs(rawLogs string, ghJobs []ghJobItem) []FailedJobItem {
	parsedJobs := ParseFailedLogLines(rawLogs)
	if len(ghJobs) == 0 {
		return parsedJobs
	}

	ghFailed := extractFailingJobsAndSteps(ghJobs)
	if len(ghFailed) == 0 {
		return parsedJobs
	}

	return mergeCorrelatedFailedJobs(parsedJobs, ghFailed, rawLogs)
}

// CorrelateRunFailedJobs fetches jobs via gh run view --json jobs and correlates them with raw logs.
func CorrelateRunFailedJobs(repo string, runId uint64, rawLogs string) []FailedJobItem {
	ghJobs := queryRunJobs(repo, runId)
	if len(ghJobs) == 0 {
		return ParseFailedLogLines(rawLogs)
	}

	return CorrelateFailedJobs(rawLogs, ghJobs)
}

func mergeCorrelatedFailedJobs(parsed, ghFailed []FailedJobItem, rawLogs string) []FailedJobItem {
	var merged []FailedJobItem
	matchedIndices := make(map[int]bool)
	for _, target := range ghFailed {
		item := matchOrBuildTargetFailure(target, parsed, matchedIndices, rawLogs)
		merged = append(merged, item)
	}

	for i, p := range parsed {
		if !matchedIndices[i] {
			merged = append(merged, p)
		}
	}

	return merged
}

func matchOrBuildTargetFailure(target FailedJobItem, parsed []FailedJobItem, matchedIndices map[int]bool, rawLogs string) FailedJobItem {
	for i, p := range parsed {
		if !matchedIndices[i] && isMatchingJobStep(p, target) {
			matchedIndices[i] = true

			return mergeMatchedFailure(target, p)
		}
	}

	return searchRawLogsForTarget(target, rawLogs)
}

func isMatchingJobStep(p, target FailedJobItem) bool {
	hasJobMatch := isMatchingJobName(p.JobName, target.JobName)
	hasStepMatch := isMatchingStepName(p.StepName, target.StepName)

	return hasJobMatch && (hasStepMatch || p.StepName == "Job Execution")
}

func isMatchingJobName(a, b string) bool {
	if a == b || strings.EqualFold(a, b) {
		return true
	}

	lowerA, lowerB := strings.ToLower(a), strings.ToLower(b)

	return strings.Contains(lowerA, lowerB) || strings.Contains(lowerB, lowerA)
}

func isMatchingStepName(a, b string) bool {
	if a == b || strings.EqualFold(a, b) {
		return true
	}

	lowerA, lowerB := strings.ToLower(a), strings.ToLower(b)

	return strings.Contains(lowerA, lowerB) || strings.Contains(lowerB, lowerA)
}

func mergeMatchedFailure(target, p FailedJobItem) FailedJobItem {
	item := FailedJobItem{
		JobName:        target.JobName,
		StepName:       target.StepName,
		FailureSummary: p.FailureSummary,
		ErrorLines:     p.ErrorLines,
	}

	if len(item.FailureSummary) == 0 {
		item.FailureSummary = target.FailureSummary
	}

	if len(item.ErrorLines) == 0 {
		item.ErrorLines = target.ErrorLines
	}

	return item
}

func searchRawLogsForTarget(target FailedJobItem, rawLogs string) FailedJobItem {
	lines := extractMatchingLogLines(rawLogs, target.StepName)
	if len(lines) == 0 {
		lines = extractMatchingLogLines(rawLogs, target.JobName)
	}

	if len(lines) > 0 {
		target.ErrorLines = lines
		target.FailureSummary = lines[0]
	}

	return target
}

func extractMatchingLogLines(rawLogs, query string) []string {
	if len(query) == 0 || len(rawLogs) == 0 {
		return nil
	}

	var matches []string
	for _, line := range strings.Split(rawLogs, "\n") {
		clean := cleanLogText(ansiRegex.ReplaceAllString(line, ""))
		if strings.Contains(strings.ToLower(clean), strings.ToLower(query)) && !isIgnoredLogLine(clean) {
			matches = append(matches, clean)
		}
	}

	return capErrorLines(matches, 20)
}

func capErrorLines(lines []string, limit int) []string {
	if len(lines) > limit {
		return lines[:limit]
	}

	return lines
}

func extractAllSectionFailures(failedRuns []FailedRunItem) []SectionFailure {
	var sections []SectionFailure
	for _, run := range failedRuns {
		sections = append(sections, extractRunSectionFailures(run)...)
	}

	return sections
}

func extractRunSectionFailures(run FailedRunItem) []SectionFailure {
	jobs := resolveRunCorrelatedJobs(run)
	var out []SectionFailure
	for _, job := range jobs {
		out = append(out, buildSectionFailureFromRunJob(run, job))
	}

	return out
}

func resolveRunCorrelatedJobs(run FailedRunItem) []FailedJobItem {
	if run.RunId == 0 {
		return run.FailedJobs
	}

	ghJobs := queryRunJobs("", run.RunId)
	if len(ghJobs) == 0 {
		return run.FailedJobs
	}

	return CorrelateFailedJobs(run.RawErrors, ghJobs)
}

func buildSectionFailureFromRunJob(run FailedRunItem, job FailedJobItem) SectionFailure {
	return SectionFailure{
		WorkflowName:   run.WorkflowName,
		RunId:          run.RunId,
		JobName:        job.JobName,
		StepName:       job.StepName,
		FailureSummary: job.FailureSummary,
		ErrorLines:     job.ErrorLines,
		SavedLogFile:   toRelativeGitPath(run.SavedLogFile),
		CreatedAt:      run.CreatedAt,
	}
}

func toRelativeGitPath(targetPath string) string {
	if len(targetPath) == 0 {
		return ""
	}

	root := resolveRepoRootDir()
	rel, err := filepath.Rel(root, targetPath)
	if err != nil || strings.HasPrefix(rel, "..") {
		return filepath.ToSlash(targetPath)
	}

	return filepath.ToSlash(rel)
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

	if len(sec.FailureSummary) > 0 {
		sb.WriteString(fmt.Sprintf("  │ Summary:   %s\n", sec.FailureSummary))
	}
}

func formatSectionErrors(sb *strings.Builder, sec SectionFailure) {
	for _, line := range sec.ErrorLines {
		sb.WriteString(fmt.Sprintf("  │   %s\n", line))
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
	if len(job.FailureSummary) > 0 {
		sb.WriteString(fmt.Sprintf("      Summary: %s\n", job.FailureSummary))
	}

	for _, l := range job.ErrorLines {
		sb.WriteString(fmt.Sprintf("      %s\n", l))
	}
}
