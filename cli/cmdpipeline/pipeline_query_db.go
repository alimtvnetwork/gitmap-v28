package cmdpipeline

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/pipelinedb"
)

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
	if sec == 0 {
		return "-"
	}
	if sec < 0 {
		return fmt.Sprintf("overtime +%s", formatDurationSeconds(-sec))
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
