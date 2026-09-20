package cmdpipeline

import (
	"fmt"
	"strings"
)

// AttachLivePipelineErrors inspects active running jobs and attaches failure diagnostics if any error occurred.
func AttachLivePipelineErrors(payload *PipelineStatusPayload) {
	if payload == nil || payload.LastRunId == 0 {
		return
	}

	jobs := queryRunJobs(payload.Repo, payload.LastRunId)
	if !isRunOrJobFailed(payload, jobs) {
		return
	}

	populateLiveErrorDiagnostics(payload, jobs)
}

func isRunOrJobFailed(payload *PipelineStatusPayload, jobs []ghJobItem) bool {
	if payload.LastConclusion == "failure" {
		return true
	}

	return hasAnyFailingJobOrStep(jobs)
}

func hasAnyFailingJobOrStep(jobs []ghJobItem) bool {
	for _, j := range jobs {
		if isJobFailing(j) || hasAnyFailingStep(j.Steps) {
			return true
		}
	}

	return false
}

func hasAnyFailingStep(steps []ghStepItem) bool {
	for _, s := range steps {
		if isStepFailing(s) {
			return true
		}
	}

	return false
}

func populateLiveErrorDiagnostics(payload *PipelineStatusPayload, jobs []ghJobItem) {
	rawLogs := queryFailedRunLogs(payload.Repo, payload.LastRunId)
	cleanErrors := extractCleanErrorLines(rawLogs)
	failedJobs := resolveFailedJobItems(payload.Repo, payload.LastRunId, rawLogs, jobs)

	payload.HasErrors = true
	payload.IsStopWaiting = true
	payload.RecommendedAction = "fix_errors"
	payload.FailedJobCount = len(failedJobs)
	payload.ErrorSummary = buildLiveErrorSummary(failedJobs)
	payload.ErrorLogs = resolveLiveErrorLogs(cleanErrors, rawLogs, failedJobs)
	payload.ActionableErrorSnippet = buildActionableErrorSnippet(cleanErrors, failedJobs)
	payload.FailedJobs = failedJobs
}

func resolveFailedJobItems(repo string, runId uint64, rawLogs string, jobs []ghJobItem) []FailedJobItem {
	if len(rawLogs) > 0 {
		items := CorrelateFailedJobs(rawLogs, jobs)
		if len(items) > 0 {
			return items
		}
	}

	return extractLiveFailingJobsAndSteps(jobs)
}

func extractLiveFailingJobsAndSteps(jobs []ghJobItem) []FailedJobItem {
	var results []FailedJobItem
	for _, j := range jobs {
		if isJobFailing(j) || hasAnyFailingStep(j.Steps) {
			results = append(results, extractFailingStepsFromJob(j)...)
		}
	}

	return results
}

func buildLiveErrorSummary(failedJobs []FailedJobItem) string {
	if len(failedJobs) == 0 {
		return "Pipeline run reported failure"
	}

	if len(failedJobs) == 1 {
		return failedJobs[0].FailureSummary
	}

	return fmt.Sprintf("%d failing steps detected in pipeline run", len(failedJobs))
}

func resolveLiveErrorLogs(cleanErrors, rawLogs string, failedJobs []FailedJobItem) string {
	if len(strings.TrimSpace(cleanErrors)) > 0 {
		return cleanErrors
	}

	if len(strings.TrimSpace(rawLogs)) > 0 {
		return rawLogs
	}

	return formatFailedJobsAsLogs(failedJobs)
}

func formatFailedJobsAsLogs(failedJobs []FailedJobItem) string {
	var lines []string
	for _, j := range failedJobs {
		lines = append(lines, fmt.Sprintf("FAILED: %s - %s", j.JobName, j.StepName))
		lines = append(lines, j.ErrorLines...)
	}

	return strings.Join(lines, "\n")
}

func buildActionableErrorSnippet(cleanErrors string, failedJobs []FailedJobItem) string {
	if len(strings.TrimSpace(cleanErrors)) > 0 {
		return truncateErrorLines(cleanErrors, 15)
	}

	fallbackText := formatFailedJobsAsLogs(failedJobs)

	return truncateErrorLines(fallbackText, 15)
}

func truncateErrorLines(text string, maxLines int) string {
	lines := strings.Split(text, "\n")
	if len(lines) <= maxLines {
		return text
	}

	return strings.Join(lines[:maxLines], "\n")
}
