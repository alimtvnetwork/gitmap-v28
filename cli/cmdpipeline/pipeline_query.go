package cmdpipeline

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/ghtoken"
	"github.com/alimtvnetwork/gitmap-v28/cli/pipelinedb"
)

func runGHCommandWithCustomTimeout(timeout time.Duration, args ...string) ([]byte, error) {
	if os.Getenv("GITMAP_MOCK_GH") == "1" {
		return []byte("[]"), nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "gh", args...)
	attachGHTokenEnv(cmd)

	return cmd.CombinedOutput()
}

func attachGHTokenEnv(cmd *exec.Cmd) {
	tok, _, err := ghtoken.Resolve()
	if err != nil || len(tok) == 0 {
		return
	}

	cmd.Env = buildSubprocessEnvWithToken(tok)
}

func buildSubprocessEnvWithToken(tok string) []string {
	var env []string
	for _, entry := range os.Environ() {
		if isGHTokenEnvKey(entry) {
			continue
		}

		env = append(env, entry)
	}

	env = append(env, "GH_TOKEN="+tok)
	env = append(env, "GITHUB_TOKEN="+tok)

	return env
}

func isGHTokenEnvKey(entry string) bool {
	return strings.HasPrefix(entry, "GH_TOKEN=") || strings.HasPrefix(entry, "GITHUB_TOKEN=")
}

func runGHCommandWithTimeout(args ...string) ([]byte, error) {
	return runGHCommandWithCustomTimeout(10*time.Second, args...)
}

func queryWorkflowRuns(repo string) []ghRunItem {
	if len(repo) == 0 {
		return nil
	}

	out, err := runGHCommandWithTimeout("run", "list", "--repo", repo, "--limit", "30", "--json",
		"databaseId,name,status,conclusion,createdAt,updatedAt,headBranch,headSha,url,displayTitle,event")
	if err != nil {
		return queryRunsFromDB(repo)
	}

	return parseRunsOrFallback(out, repo)
}

func parseRunsOrFallback(out []byte, repo string) []ghRunItem {
	var runs []ghRunItem
	if err := json.Unmarshal(out, &runs); err != nil {
		return queryRunsFromDB(repo)
	}

	return runs
}

func queryPendingPRs(repo string) int {
	if len(repo) == 0 {
		return 0
	}

	out, err := runGHCommandWithTimeout("pr", "list", "--repo", repo, "--state", "open", "--json", "number")
	if err != nil {
		return 0
	}

	return parseOpenPRsCount(out)
}

func parseOpenPRsCount(out []byte) int {
	var prs []map[string]any
	if err := json.Unmarshal(out, &prs); err != nil {
		return 0
	}

	return len(prs)
}

func queryLatestTagRelease(repo string) string {
	tag := queryGHLatestTag(repo)
	if len(tag) > 0 {
		return tag
	}

	tagOut, err := exec.Command("git", "describe", "--tags", "--abbrev=0").Output()
	if err == nil && len(tagOut) > 0 {
		return strings.TrimSpace(string(tagOut))
	}

	return "v" + constants.Version
}

type ghReleaseTagItem struct {
	TagName string `json:"tagName"`
}

func queryGHLatestTag(repo string) string {
	if len(repo) == 0 {
		return ""
	}

	out, err := runGHCommandWithTimeout("release", "list", "--repo", repo, "--limit", "1", "--json", "tagName")
	if err != nil || len(out) == 0 {
		return ""
	}

	return parseFirstReleaseTag(out)
}

func parseFirstReleaseTag(out []byte) string {
	var releases []ghReleaseTagItem
	if err := json.Unmarshal(out, &releases); err == nil && len(releases) > 0 {
		return releases[0].TagName
	}

	return ""
}

func queryFailedRunLogs(repo string, runId uint64) string {
	if runId == 0 || len(repo) == 0 {
		return ""
	}

	if cached, ok := readCachedPipelineLogForRepo(repo, runId); ok {
		return cached
	}

	return fetchAndCacheRunLogs(repo, runId)
}

func fetchAndCacheRunLogs(repo string, runId uint64) string {
	idStr := strconv.FormatUint(runId, 10)
	out, err := runGHCommandWithCustomTimeout(60*time.Second, "run", "view", idStr, "--repo", repo, "--log-failed")
	if err == nil && len(out) > 0 {
		logStr := string(out)
		_ = writeCachedPipelineLog(runId, logStr, repo)

		return logStr
	}

	return handleFailedRunLogsFallback(repo, runId, err, out)
}

func handleFailedRunLogsFallback(repo string, runId uint64, err error, out []byte) string {
	fallback := buildFallbackRunLogs(repo, runId)
	if len(fallback) > 0 {
		_ = writeCachedPipelineLog(runId, fallback, repo)

		return fallback
	}

	if err != nil {
		return formatGHFailedError(err, out, repo, runId)
	}

	return "Unable to fetch failed logs via gh CLI."
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
	if !isJobFailingOrCancelled(j) {
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

func isJobFailingOrCancelled(j ghJobItem) bool {
	return j.Conclusion == "failure" || j.Conclusion == "cancelled" || j.Conclusion == "timed_out" || j.Conclusion == "startup_failure"
}

func formatConclusionTag(conclusion string) string {
	if conclusion == "cancelled" {
		return "CANCELLED"
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

func queryRunDiagnostic(repo string, runId uint64) string {
	if runId == 0 || len(repo) == 0 {
		return ""
	}

	idStr := strconv.FormatUint(runId, 10)
	out, err := runGHCommandWithTimeout("run", "view", idStr, "--repo", repo)
	if err != nil || len(out) == 0 {
		return ""
	}

	return extractDiagnosticFromRunView(string(out))
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

func queryRunsFromDB(repo string) []ghRunItem {
	db, err := pipelinedb.OpenPipelineSplitDb(repo)
	if err != nil {
		return nil
	}

	defer db.Close()

	runRes := db.QueryRecentRuns(20)
	if runRes.IsFailure() {
		return nil
	}

	return mapDbRunsToGhRuns(runRes.Data)
}

func queryFreshDbRuns(repo string) ([]ghRunItem, bool) {
	dbRuns := queryRunsFromDB(repo)
	hasRuns := len(dbRuns) > 0
	if hasRuns {
		return dbRuns, true
	}

	return nil, false
}

func resolveCachedDbRunsIfFresh(repo string) ([]ghRunItem, bool) {
	dbPath := pipelinedb.PipelineDbPath(repo)
	isCacheHit := checkTtlCacheHit(dbPath)
	if isCacheHit {
		return queryFreshDbRuns(repo)
	}

	return nil, false
}

func resolveCachedRunsOrFetch(repo string) []ghRunItem {
	dbRuns, isFresh := resolveCachedDbRunsIfFresh(repo)
	if isFresh {
		return dbRuns
	}

	runs := queryWorkflowRuns(repo)
	hasRuns := len(runs) > 0
	if hasRuns {
		RecordFetchedRunsToSplitDb(repo, runs)

		return runs
	}

	return queryRunsFromDB(repo)
}

func mapDbRunsToGhRuns(dbRuns []pipelinedb.PipelineRunRecord) []ghRunItem {
	runs := make([]ghRunItem, 0, len(dbRuns))
	for _, r := range dbRuns {
		runs = append(runs, convertDbRunToGhRun(r))
	}

	return runs
}

func convertDbRunToGhRun(r pipelinedb.PipelineRunRecord) ghRunItem {
	return ghRunItem{
		DatabaseId: r.RunId,
		Name:       r.WorkflowName,
		Status:     r.Status,
		Conclusion: r.Conclusion,
		HeadBranch: r.Branch,
		HeadSha:    r.Sha,
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
		Url:        r.RunUrl,
	}
}

func formatRunTimestamp(raw string) string {
	if raw == "" {
		return "unknown time"
	}

	t, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return raw
	}

	return formatParsedRunTime(t)
}

func formatParsedRunTime(t time.Time) string {
	utcStr := t.UTC().Format("2006-01-02 15:04:05 UTC")
	elapsed := time.Since(t)
	if elapsed < 0 {
		return utcStr
	}

	return fmt.Sprintf("%s (%s ago)", utcStr, formatDurationShort(elapsed))
}

func formatDurationShort(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}

	if d < time.Hour {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	}

	if d < 24*time.Hour {
		return fmt.Sprintf("%dh", int(d.Hours()))
	}

	return fmt.Sprintf("%dd", int(d.Hours()/24))
}

func calculateRunDuration(createdAt, updatedAt string) int {
	t1, err1 := time.Parse(time.RFC3339, createdAt)
	t2, err2 := time.Parse(time.RFC3339, updatedAt)
	if err1 != nil || err2 != nil || !t2.After(t1) {
		return 0
	}

	return int(t2.Sub(t1).Seconds())
}

func formatDurationSeconds(sec int) string {
	if sec <= 0 {
		return "<1s"
	}
	if sec < 60 {
		return fmt.Sprintf("%ds", sec)
	}
	if sec < 3600 {
		return formatMinutesAndSeconds(sec/60, sec%60)
	}

	return formatHoursMinutesSeconds(sec/3600, (sec%3600)/60, sec%60)
}

func formatMinutesAndSeconds(m, s int) string {
	if s == 0 {
		return fmt.Sprintf("%dm", m)
	}

	return fmt.Sprintf("%dm %ds", m, s)
}

func formatHoursMinutesSeconds(h, m, s int) string {
	if m == 0 && s == 0 {
		return fmt.Sprintf("%dh", h)
	}
	if s == 0 {
		return fmt.Sprintf("%dh %dm", h, m)
	}

	return fmt.Sprintf("%dh %dm %ds", h, m, s)
}

func formatEtaDisplay(sec int) string {
	if sec <= 0 {
		return "-"
	}
	if sec < 60 {
		return fmt.Sprintf("~%ds", sec)
	}

	return fmt.Sprintf("~%s (%ds)", formatDurationSeconds(sec), sec)
}

func resolveCurrentRepoSlug() string {
	out, err := exec.Command("git", "config", "--get", "remote.origin.url").Output()
	if err == nil && len(out) > 0 {
		return parseSlugFromGitURL(strings.TrimSpace(string(out)))
	}

	return "alimtvnetwork/gitmap-v28"
}

func parseSlugFromGitURL(raw string) string {
	clean := strings.TrimSuffix(raw, ".git")
	if strings.Contains(clean, "github.com/") {
		return extractSlugAfterToken(clean, "github.com/")
	}

	if strings.Contains(clean, "github.com:") {
		return extractSlugAfterToken(clean, "github.com:")
	}

	return clean
}

func extractSlugAfterToken(clean, token string) string {
	parts := strings.Split(clean, token)
	if len(parts) > 1 {
		return parts[1]
	}

	return clean
}
