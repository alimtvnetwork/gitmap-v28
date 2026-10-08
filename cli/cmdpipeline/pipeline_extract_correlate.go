package cmdpipeline

import (
	"path/filepath"
	"strings"
)

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
	return CorrelateFailedJobsWithRunJobs(rawLogs, queryRunJobs(repo, runId))
}

// CorrelateFailedJobsWithRunJobs correlates pre-queried run jobs with raw logs.
func CorrelateFailedJobsWithRunJobs(rawLogs string, ghJobs []ghJobItem) []FailedJobItem {
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
	if len(a) == 0 || len(b) == 0 {
		return false
	}
	if a == b || strings.EqualFold(a, b) {
		return true
	}

	lowerA, lowerB := strings.ToLower(a), strings.ToLower(b)

	return strings.Contains(lowerA, lowerB) || strings.Contains(lowerB, lowerA)
}

func isMatchingStepName(a, b string) bool {
	if len(a) == 0 || len(b) == 0 {
		return false
	}
	if a == b || strings.EqualFold(a, b) {
		return true
	}

	lowerA, lowerB := strings.ToLower(a), strings.ToLower(b)

	return strings.Contains(lowerA, lowerB) || strings.Contains(lowerB, lowerA)
}

func mergeMatchedFailure(target, p FailedJobItem) FailedJobItem {
	stack := target.StackTrace
	if len(stack) == 0 {
		stack = p.StackTrace
	}

	item := FailedJobItem{
		JobName:        target.JobName,
		StepName:       target.StepName,
		FailureSummary: p.FailureSummary,
		ErrorLines:     p.ErrorLines,
		Warnings:       p.Warnings,
		StackTrace:     stack,
	}

	if len(item.FailureSummary) == 0 {
		item.FailureSummary = target.FailureSummary
	}
	if len(item.ErrorLines) == 0 {
		item.ErrorLines = target.ErrorLines
	}
	if len(item.Warnings) == 0 {
		item.Warnings = target.Warnings
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
	if len(target.StackTrace) == 0 {
		target.StackTrace = extractJobStackTrace(rawLogs, target.JobName)
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

func resolveRunRawLogs(run FailedRunItem) string {
	if len(run.RawErrors) > 0 {
		return run.RawErrors
	}

	cached, hasCached := readCachedPipelineLog(run.RunId)
	if hasCached {
		return cached
	}

	return ""
}

func resolveRunCorrelatedJobs(run FailedRunItem) []FailedJobItem {
	if run.RunId == 0 {
		return run.FailedJobs
	}

	rawLogs := resolveRunRawLogs(run)
	ghJobs := queryRunJobs("", run.RunId)
	if len(ghJobs) == 0 {
		return resolveJobFallback(run.FailedJobs, rawLogs)
	}

	return CorrelateFailedJobs(rawLogs, ghJobs)
}

func resolveJobFallback(failedJobs []FailedJobItem, rawLogs string) []FailedJobItem {
	if len(failedJobs) > 0 {
		return failedJobs
	}

	return ParseFailedLogLines(rawLogs)
}

func buildSectionFailureFromRunJob(run FailedRunItem, job FailedJobItem) SectionFailure {
	stack := job.StackTrace
	if len(stack) == 0 && len(run.FailedJobs) <= 1 {
		stack = run.StackTrace
	}

	return SectionFailure{
		WorkflowName:   run.WorkflowName,
		RunId:          run.RunId,
		JobName:        job.JobName,
		StepName:       job.StepName,
		FailureSummary: job.FailureSummary,
		ErrorLines:     job.ErrorLines,
		Warnings:       job.Warnings,
		SavedLogFile:   toRelativeGitPath(run.SavedLogFile),
		CreatedAt:      run.CreatedAt,
		StackTrace:     stack,
	}
}

func toRelativeGitPath(targetPath string) string {
	if len(targetPath) == 0 {
		return ""
	}

	root := resolveRepoRootDir()
	rel, err := filepath.Rel(root, targetPath)
	if err != nil || strings.HasPrefix(rel, "..") {
		return filepath.ToSlash(FormatRelativeDbPath(targetPath))
	}

	return filepath.ToSlash(rel)
}
