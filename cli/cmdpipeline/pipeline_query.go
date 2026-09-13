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
)

func runGHCommandWithCustomTimeout(timeout time.Duration, args ...string) ([]byte, error) {
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
		"databaseId,name,status,conclusion,createdAt,updatedAt,headBranch,headSha,url")
	if err != nil {
		return queryRunsFromDB(repo)
	}

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

func queryGHLatestTag(repo string) string {
	if len(repo) == 0 {
		return ""
	}

	out, err := runGHCommandWithTimeout("release", "list", "--repo", repo, "--limit", "1", "--json", "tagName")
	if err != nil || len(out) == 0 {
		return ""
	}

	var releases []struct {
		TagName string `json:"tagName"`
	}

	if err := json.Unmarshal(out, &releases); err == nil && len(releases) > 0 {
		return releases[0].TagName
	}

	return ""
}

func queryFailedRunLogs(repo string, runId uint64) string {
	if runId == 0 || len(repo) == 0 {
		return ""
	}

	if cached, ok := readCachedPipelineLog(runId); ok {
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
		return formatGHFailedError(err, out)
	}

	return "Unable to fetch failed logs via gh CLI."
}

func buildFallbackRunLogs(repo string, runId uint64) string {
	jobs := queryRunJobs(repo, runId)
	if len(jobs) == 0 {
		return ""
	}

	var sb strings.Builder
	for _, j := range jobs {
		appendJobFallbackLogs(&sb, j)
	}

	return sb.String()
}

func appendJobFallbackLogs(sb *strings.Builder, j ghJobItem) {
	if j.Conclusion != "failure" {
		return
	}

	for _, s := range j.Steps {
		if s.Conclusion == "failure" {
			sb.WriteString(fmt.Sprintf("%s\t%s\tFAIL: Step '%s' (step #%d) failed\n",
				j.Name, s.Name, s.Name, s.Number))
		}
	}
}

func formatGHFailedError(err error, out []byte) string {
	msg := strings.TrimSpace(string(out))
	if len(msg) > 0 {
		return fmt.Sprintf("gh command failed (%v):\n%s", err, msg)
	}

	return fmt.Sprintf("gh command failed: %v", err)
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
	db, err := openDB()
	if err != nil {
		return nil
	}

	defer db.Close()

	dbRuns, err := db.ListRecentPipelineRuns(repo, 5)
	if err != nil {
		return nil
	}

	var runs []ghRunItem
	for _, r := range dbRuns {
		runs = append(runs, ghRunItem{
			DatabaseId: safeInt64ToUint64(r.RunID),
			Name:       r.WorkflowName,
			Status:     r.Status,
			Conclusion: r.Conclusion,
			HeadBranch: r.Branch,
			HeadSha:    r.Sha,
			CreatedAt:  r.CreatedAt,
			UpdatedAt:  r.UpdatedAt,
			Url:        r.URL,
		})
	}

	return runs
}

func formatRunTimestamp(raw string) string {
	if raw == "" {
		return "unknown time"
	}

	t, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return raw
	}

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

	m := sec / 60
	s := sec % 60
	if s == 0 {
		return fmt.Sprintf("%dm", m)
	}

	return fmt.Sprintf("%dm %ds", m, s)
}

func safeInt64ToUint64(val int64) uint64 {
	if val < 0 {
		return 0
	}

	return uint64(val)
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
