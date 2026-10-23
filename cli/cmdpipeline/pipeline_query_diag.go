package cmdpipeline

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

func cacheRunFallbackIfFinished(repo string, runId uint64, fallback string) {
	if !isRunInProgress(repo, runId) {
		_ = writeCachedPipelineLog(runId, fallback, repo)
	}
}

func isRunInProgress(repo string, runId uint64) bool {
	if runId == 0 {
		return false
	}

	idStr := strconv.FormatUint(runId, 10)
	args := []string{"run", "view", idStr, "--json", "status,conclusion"}
	if len(repo) > 0 {
		args = append(args, "--repo", repo)
	}
	out, err := runGHCommandWithTimeout(args...)
	if err != nil || len(out) == 0 {
		return checkJobsActive(repo, runId)
	}

	var resp struct {
		Status     string `json:"status"`
		Conclusion string `json:"conclusion"`
	}
	if json.Unmarshal(out, &resp) == nil && isGHStatusActive(resp.Status) {
		return true
	}

	return checkJobsActive(repo, runId)
}

func checkJobsActive(repo string, runId uint64) bool {
	jobs := queryRunJobs(repo, runId)
	if len(jobs) > 0 {
		return !isAllJobsCompleted(jobs)
	}

	return false
}

func queryRunConclusion(repo string, runId uint64) string {
	if runId == 0 {
		return ""
	}
	idStr := strconv.FormatUint(runId, 10)
	args := []string{"run", "view", idStr, "--json", "conclusion"}
	if len(repo) > 0 {
		args = append(args, "--repo", repo)
	}
	out, err := runGHCommandWithTimeout(args...)
	if err != nil || len(out) == 0 {
		return ""
	}
	var resp struct {
		Conclusion string `json:"conclusion"`
	}
	if json.Unmarshal(out, &resp) == nil {
		return strings.ToLower(strings.TrimSpace(resp.Conclusion))
	}

	return ""
}

func isGHStatusActive(status string) bool {
	return status == "in_progress" || status == "queued" || status == "waiting" || status == "pending"
}

func buildFallbackRunLogs(repo string, runId uint64) string {
	jobs := queryRunJobs(repo, runId)
	diag := queryRunDiagnostic(repo, runId)
	hasJobs := len(jobs) > 0
	if !hasJobs && len(diag) == 0 {
		return ""
	}

	var sb strings.Builder
	for _, j := range jobs {
		appendJobFallbackLogs(&sb, j, diag)
	}
	if sb.Len() == 0 && len(diag) > 0 {
		sb.WriteString(fmt.Sprintf("Workflow Execution\tJob Execution\tFAIL: %s\n", diag))
	}

	return sb.String()
}

func appendJobFallbackLogs(sb *strings.Builder, j ghJobItem, diag string) {
	if !isJobFailingOrCanceled(j) {
		return
	}

	hasFailedStep := appendStepFallbackLogs(sb, j)
	if !hasFailedStep {
		reason := resolveJobFailureReason(j, diag)
		sb.WriteString(fmt.Sprintf("%s\tJob Execution\t%s: %s\n",
			j.Name, formatConclusionTag(j.Conclusion), reason))
	}
}

func appendStepFallbackLogs(sb *strings.Builder, j ghJobItem) bool {
	hasFailedStep := false
	for _, s := range j.Steps {
		if isStepFailing(s) {
			hasFailedStep = true
			sb.WriteString(fmt.Sprintf("%s\t%s\t%s: Step '%s' (step #%d) %s\n",
				j.Name, s.Name, formatConclusionTag(s.Conclusion), s.Name, s.Number, s.Conclusion))
		}
	}

	return hasFailedStep
}

func isJobFailingOrCanceled(j ghJobItem) bool {
	return j.Conclusion == "failure" || strings.HasPrefix(j.Conclusion, "cancel") || j.Conclusion == "timed_out" || j.Conclusion == "startup_failure"
}

func formatConclusionTag(conclusion string) string {
	if strings.HasPrefix(conclusion, "cancel") {
		return "CANCELED"
	}
	if conclusion == "timed_out" {
		return "TIMED_OUT"
	}

	return "FAIL"
}

func resolveJobFailureReason(j ghJobItem, diag string) string {
	if len(diag) > 0 {
		return fmt.Sprintf("Job '%s' %s: %s", j.Name, j.Conclusion, diag)
	}

	return fmt.Sprintf("Job '%s' ended with conclusion '%s'", j.Name, j.Conclusion)
}

func ensureRunDiagCached(repo string, runId uint64, diag string) string {
	if len(diag) > 0 {
		return diag
	}

	return queryRunDiagnostic(repo, runId)
}

var (
	runDiagCacheLock sync.RWMutex
	runDiagCache     = make(map[uint64]string)
)

func queryRunDiagnostic(repo string, runId uint64) string {
	if runId == 0 || len(repo) == 0 {
		return ""
	}

	runDiagCacheLock.RLock()
	cached, isFound := runDiagCache[runId]
	runDiagCacheLock.RUnlock()
	if isFound {
		return cached
	}

	idStr := strconv.FormatUint(runId, 10)
	out, err := runGHCommandWithTimeout("run", "view", idStr, "--repo", repo)
	result := ""
	if err == nil && len(out) > 0 {
		result = extractDiagnosticFromRunView(string(out))
	}

	runDiagCacheLock.Lock()
	runDiagCache[runId] = result
	runDiagCacheLock.Unlock()

	return result
}

func extractDiagnosticFromRunView(text string) string {
	for _, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(line)
		if isDiagnosticLine(trimmed) {
			return strings.TrimPrefix(trimmed, "X ")
		}
	}

	return ""
}

func isDiagnosticLine(trimmed string) bool {
	if strings.Contains(trimmed, "workflow file issue") || strings.Contains(trimmed, "likely failed") {
		return true
	}
	if strings.Contains(trimmed, "exceeded the maximum execution time") || strings.Contains(trimmed, "timed out") {
		return true
	}
	if strings.HasPrefix(trimmed, "X ") && !strings.Contains(trimmed, " · ") && !strings.Contains(trimmed, " in ") {
		return true
	}

	return false
}

func formatGHFailedError(err error, out []byte, repo string, runId uint64) string {
	msg := strings.TrimSpace(string(out))
	diag := queryRunDiagnostic(repo, runId)

	return buildGHFailedErrorMessage(err, msg, diag)
}

func buildGHFailedErrorMessage(err error, rawMsg, diag string) string {
	var sb strings.Builder
	errMsg := "command execution failed"
	if err != nil {
		errMsg = extractCleanErrorMessage(err)
	}
	sb.WriteString(fmt.Sprintf("gh command notice (%s):\n", errMsg))
	hasRawMsg := len(rawMsg) > 0
	if hasRawMsg {
		sb.WriteString(fmt.Sprintf("  %s\n", rawMsg))
	}

	hasDiag := len(diag) > 0
	if hasDiag {
		sb.WriteString(fmt.Sprintf("  Diagnostic: %s\n", diag))
	}

	return sb.String()
}

func extractCleanErrorMessage(err error) string {
	msg := err.Error()
	if idx := strings.Index(msg, "\n"); idx != -1 {
		msg = msg[:idx]
	}

	return strings.TrimSpace(msg)
}
