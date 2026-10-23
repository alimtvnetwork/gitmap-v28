package cmdpipeline

import (
	"strconv"
	"strings"
	"sync"
	"time"
)

func queryFailedRunLogs(repo string, runId uint64) string {
	if runId == 0 || len(repo) == 0 {
		return ""
	}

	cached, ok := readCachedPipelineLogForRepo(repo, runId)
	if ok && !isOnlyFallbackLogs(cached) {
		return cached
	}

	return fetchAndCacheRunLogs(repo, runId)
}

func isOnlyFallbackLogs(content string) bool {
	trimmed := strings.TrimSpace(content)
	if len(trimmed) == 0 {
		return true
	}

	lines := strings.Split(trimmed, "\n")
	hasFallbackLine := false
	for _, rawLine := range lines {
		line := strings.TrimSpace(rawLine)
		if len(line) == 0 {
			continue
		}
		if isFallbackLine(line) {
			hasFallbackLine = true

			continue
		}

		return false
	}

	return hasFallbackLine
}

func isFallbackLine(line string) bool {
	if strings.Contains(line, " (step #") || (strings.Contains(line, "Step '") && strings.Contains(line, "step #")) {
		return true
	}
	if strings.Contains(line, "Job Execution\tFAIL:") || strings.Contains(line, "Job Execution\tCANCELED:") || strings.Contains(line, "Job Execution\tTIMED_OUT:") {
		return true
	}
	if strings.Contains(line, "\tJob Execution\t") || strings.HasPrefix(line, "Workflow Execution\t") {
		return true
	}
	if strings.Contains(line, "ended with conclusion '") {
		return true
	}
	if strings.HasPrefix(line, "gh command notice") || strings.HasPrefix(line, "Unable to fetch failed logs") || strings.HasPrefix(line, "gh command failed") {
		return true
	}

	return false
}

func fetchAndCacheRunLogs(repo string, runId uint64) string {
	idStr := strconv.FormatUint(runId, 10)
	out, err := runGHCommandWithCustomTimeout(15*time.Second, "run", "view", idStr, "--repo", repo, "--log-failed")
	if err == nil && len(out) > 0 {
		logStr := string(out)
		_ = writeCachedPipelineLog(runId, logStr, repo)

		return logStr
	}

	return handleFailedRunLogsFallback(repo, runId, err, out)
}

func handleFailedRunLogsFallback(repo string, runId uint64, err error, out []byte) string {
	if realLogs, hasRealLogs := resolveJobLevelFailedLogs(repo, runId); hasRealLogs {
		cacheRunFallbackIfFinished(repo, runId, realLogs)

		return realLogs
	}

	fallback := buildFallbackRunLogs(repo, runId)
	if len(fallback) > 0 {
		cacheRunFallbackIfFinished(repo, runId, fallback)

		return fallback
	}

	if err != nil {
		return formatGHFailedError(err, out, repo, runId)
	}

	conclusion := queryRunConclusion(repo, runId)
	if strings.HasPrefix(conclusion, "cancel") {
		return "Workflow run was cancelled before completion (conclusion: cancelled). No step failure logs were generated."
	}

	return "Unable to fetch failed logs via gh CLI."
}

func resolveJobLevelFailedLogs(repo string, runId uint64) (string, bool) {
	jobs := queryRunJobs(repo, runId)
	if len(jobs) == 0 {
		return "", false
	}

	return fetchJobLevelFailedLogs(repo, runId, jobs)
}

func fetchFailedJobLogs(repo string, jobId uint64) (string, error) {
	args := []string{"run", "view", "--job", strconv.FormatUint(jobId, 10), "--log-failed"}
	if len(repo) > 0 {
		args = append(args, "--repo", repo)
	}

	out, err := runGHCommandWithCustomTimeout(15*time.Second, args...)
	if err != nil {
		return "", err
	}

	return string(out), nil
}

func fetchJobLevelFailedLogs(repo string, runId uint64, jobs []ghJobItem) (string, bool) {
	results := executeParallelJobLogFetches(repo, jobs)

	return assembleJobLogsOrFallback(repo, runId, jobs, results)
}

func resolveJobLogConcurrency(total int) int {
	if total <= 8 {
		if total < 1 {
			return 1
		}
		return total
	}

	return 8
}

func executeParallelJobLogFetches(repo string, jobs []ghJobItem) []string {
	results := make([]string, len(jobs))
	sem := make(chan struct{}, resolveJobLogConcurrency(len(jobs)))
	var wg sync.WaitGroup

	for i, j := range jobs {
		if !isJobFailingOrCanceled(j) || !isEligibleForFailedLogFetch(j) {
			continue
		}

		wg.Add(1)
		go dispatchFetchJobLogWorker(&wg, sem, results, repo, j, i)
	}
	wg.Wait()

	return results
}

func dispatchFetchJobLogWorker(wg *sync.WaitGroup, sem chan struct{}, results []string, repo string, j ghJobItem, idx int) {
	defer wg.Done()
	sem <- struct{}{}
	defer func() { <-sem }()

	results[idx] = fetchAndNormalizeJobLog(repo, j)
}

func fetchAndNormalizeJobLog(repo string, j ghJobItem) string {
	jobLog, err := fetchFailedJobLogs(repo, j.DatabaseId)
	if err != nil || len(strings.TrimSpace(jobLog)) == 0 {
		return ""
	}

	return normalizeJobLogOutput(j.Name, jobLog)
}

func assembleJobLogsOrFallback(repo string, runId uint64, jobs []ghJobItem, results []string) (string, bool) {
	var sb strings.Builder
	diag := ""

	for i, j := range jobs {
		if !isJobFailingOrCanceled(j) {
			continue
		}

		if len(results[i]) > 0 {
			sb.WriteString(results[i])

			continue
		}

		diag = ensureRunDiagCached(repo, runId, diag)
		appendJobFallbackLogs(&sb, j, diag)
	}

	if sb.Len() > 0 {
		return sb.String(), true
	}

	return "", false
}

func isEligibleForFailedLogFetch(j ghJobItem) bool {
	if j.DatabaseId == 0 {
		return false
	}
	if strings.HasPrefix(j.Conclusion, "cancel") && !hasJobAnyFailedStep(j.Steps) {
		return false
	}

	return true
}

func hasJobAnyFailedStep(steps []ghStepItem) bool {
	for _, s := range steps {
		if s.Conclusion == "failure" || s.Conclusion == "timed_out" {
			return true
		}
	}

	return false
}

func normalizeJobLogOutput(jobName, rawLog string) string {
	if len(rawLog) == 0 {
		return ""
	}

	prefix := jobName + "\t"
	var sb strings.Builder
	for _, rawLine := range strings.Split(rawLog, "\n") {
		line := strings.TrimRight(rawLine, "\r")
		if len(strings.TrimSpace(line)) > 0 {
			appendNormalizedJobLogLine(&sb, prefix, line)
		}
	}

	return sb.String()
}

func appendNormalizedJobLogLine(sb *strings.Builder, prefix, line string) {
	if !strings.HasPrefix(line, prefix) {
		sb.WriteString(prefix)
	}
	sb.WriteString(line)
	sb.WriteByte('\n')
}

func queryAllRunLogs(repo string, runId uint64) string {
	if runId == 0 || len(repo) == 0 {
		return ""
	}

	idStr := strconv.FormatUint(runId, 10)
	out, err := runGHCommandWithCustomTimeout(60*time.Second, "run", "view", idStr, "--repo", repo, "--log")
	if err == nil && len(out) > 0 {
		return string(out)
	}

	return ""
}
