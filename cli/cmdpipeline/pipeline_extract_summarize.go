package cmdpipeline

import (
	"strings"
)

func updateJobSummary(item *FailedJobItem, text string) {
	if len(item.FailureSummary) == 0 || isStrongerSummary(text, item.FailureSummary) {
		item.FailureSummary = text
	}
}

func isStrongerSummary(candidate, current string) bool {
	if len(current) == 0 {
		return true
	}

	if isTestFailureSummary(candidate) && (isGenericExitCode(current) || isStepFailureNotice(current) || !isTestFailureSummary(current)) {
		return true
	}

	if isTestStatusBanner(current) && (isLocationSummary(candidate) || strings.Contains(candidate, "Expected ") || strings.Contains(candidate, "AssertionError")) {
		return true
	}

	if isTestFailureSummary(current) && !isTestFailureSummary(candidate) {
		return false
	}

	if isGenericExitCode(candidate) && !isGenericExitCode(current) {
		return false
	}

	if isStepFailureNotice(candidate) && !isStepFailureNotice(current) {
		return false
	}

	if isLocationSummary(candidate) && !isLocationSummary(current) && !isTestFailureSummary(current) {
		return true
	}

	if strings.Contains(candidate, "Expected ") || strings.Contains(candidate, "gofmt") {
		return true
	}

	if strings.Contains(candidate, "fatal error:") {
		return true
	}

	if isBundlerOrBuildFailure(candidate) && !isBundlerOrBuildFailure(current) {
		return true
	}

	return isGenericExitCode(current) || isStepFailureNotice(current)
}

func isTestStatusBanner(s string) bool {
	trimmed := strings.TrimSpace(s)
	if strings.HasPrefix(trimmed, "FAIL: ") && strings.Contains(trimmed, "(") {
		return false
	}

	return strings.HasPrefix(trimmed, "--- FAIL: ") || strings.HasPrefix(trimmed, "FAIL\t") || strings.HasPrefix(trimmed, "FAIL: ")
}

func isTestFailureSummary(s string) bool {
	trimmed := strings.TrimSpace(s)
	if hasFailPrefixOrContains(trimmed) {
		return true
	}
	if strings.HasPrefix(trimmed, "ERROR: ") || strings.Contains(trimmed, "ERROR: ") {
		return true
	}
	if isPytestFailureLine(trimmed) {
		return true
	}

	return strings.HasPrefix(trimmed, "--- FAIL: ") || strings.Contains(trimmed, "--- FAIL: ")
}

func isPytestFailureLine(trimmed string) bool {
	if !strings.HasPrefix(trimmed, "FAILED ") && !strings.Contains(trimmed, "FAILED ") {
		return false
	}

	return !strings.Contains(trimmed, "FAILED (")
}

func hasFailPrefixOrContains(trimmed string) bool {
	if strings.HasPrefix(trimmed, "FAIL: ") || strings.HasPrefix(trimmed, "FAIL:\t") {
		return true
	}

	return strings.Contains(trimmed, "FAIL: ")
}

func isStepFailureNotice(s string) bool {
	if strings.Contains(s, "Step '") && strings.Contains(s, "step #") {
		return true
	}
	if strings.Contains(s, "Job Execution\t") || strings.Contains(s, "Workflow Execution") {
		return true
	}

	return strings.Contains(s, "ended with conclusion")
}

func isBundlerOrBuildFailure(s string) bool {
	lower := strings.ToLower(s)
	if strings.Contains(lower, "failed to bundle") || strings.Contains(lower, "failed to copy") {
		return true
	}
	if strings.Contains(lower, "failed to build") || strings.Contains(lower, "does not exist") {
		return true
	}

	return strings.Contains(lower, "error[e")
}

func isLocationSummary(s string) bool {
	return strings.Contains(s, ".go:") || strings.Contains(s, ".py:") || strings.Contains(s, ".sh:")
}

func isGenericExitCode(s string) bool {
	return strings.Contains(s, "Process completed with exit code") || strings.Contains(s, "exit status 1")
}

func attachStackToItems(items []FailedJobItem, stack string) {
	if len(items) > 0 && len(stack) > 0 {
		items[0].StackTrace = stack
	}
}

func assembleJobItems(jobMap map[string]*FailedJobItem, order []string, rawLogs string) []FailedJobItem {
	if len(order) == 0 {
		items := buildFallbackJobItems(rawLogs)
		attachStackToItems(items, extractStackTraceFromLog(rawLogs))

		return items
	}

	var results []FailedJobItem
	for _, k := range order {
		item := *jobMap[k]
		if len(item.StackTrace) == 0 {
			item.StackTrace = extractJobStackTrace(rawLogs, item.JobName)
		}
		results = append(results, item)
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
