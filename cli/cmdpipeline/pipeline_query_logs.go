package cmdpipeline

import (
	"strconv"
	"strings"
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
	var sb strings.Builder
	diag := ""

	for _, j := range jobs {
		if !isJobFailingOrCanceled(j) {
			continue
		}

		if appendJobLogIfFetched(&sb, repo, j) {
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

func appendJobLogIfFetched(sb *strings.Builder, repo string, j ghJobItem) bool {
	if !isEligibleForFailedLogFetch(j) {
		return false
	}

	jobLog, err := fetchFailedJobLogs(repo, j.DatabaseId)
	if err != nil || len(strings.TrimSpace(jobLog)) == 0 {
		return false
	}

	sb.WriteString(normalizeJobLogOutput(j.Name, jobLog))

	return true
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
