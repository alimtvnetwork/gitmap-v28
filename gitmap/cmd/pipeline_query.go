package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

func queryWorkflowRuns(repo string) []ghRunItem {
	if len(repo) == 0 {
		return nil
	}

	out, err := exec.Command("gh", "run", "list", "--repo", repo, "--limit", "5", "--json",
		"databaseId,name,status,conclusion,createdAt,headBranch,headSha,url").Output()

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

	out, err := exec.Command("gh", "pr", "list", "--repo", repo, "--state", "open", "--json", "number").Output()

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

	out, err := exec.Command("gh", "release", "list", "--repo", repo, "--limit", "1", "--json", "tagName").Output()

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

func queryFailedRunLogs(repo string, runID int64) string {
	if runID <= 0 || len(repo) == 0 {
		return ""
	}

	idStr := strconv.FormatInt(runID, 10)
	out, err := exec.Command("gh", "run", "view", idStr, "--repo", repo, "--log-failed").Output()

	if err == nil && len(out) > 0 {
		return string(out)
	}

	return "Unable to fetch failed logs via gh CLI."
}

func queryAllRunLogs(repo string, runID int64) string {
	if runID <= 0 || len(repo) == 0 {
		return ""
	}

	idStr := strconv.FormatInt(runID, 10)
	out, err := exec.Command("gh", "run", "view", idStr, "--repo", repo, "--log").Output()

	if err == nil && len(out) > 0 {
		return string(out)
	}

	return ""
}

func recordPipelineInDB(p PipelineStatusPayload, runs []ghRunItem) {
	db, err := openDB()

	if err != nil {
		return
	}

	defer db.Close()

	for _, r := range runs {
		_ = db.InsertOrUpdatePipelineRun(store.PipelineRun{
			RunID:        r.DatabaseID,
			Repo:         p.Repo,
			WorkflowName: r.Name,
			Status:       r.Status,
			Conclusion:   r.Conclusion,
			Branch:       r.HeadBranch,
			Sha:          r.HeadSha,
			EtaSeconds:   p.EtaSeconds,
			URL:          r.URL,
		})
	}
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
			DatabaseID: r.RunID,
			Name:       r.WorkflowName,
			Status:     r.Status,
			Conclusion: r.Conclusion,
			HeadBranch: r.Branch,
			HeadSha:    r.Sha,
			URL:        r.URL,
		})
	}

	return runs
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

func resolveTempDir() string {
	db, err := openDB()

	if err != nil {
		return filepath.Join(".", ".lovable", "temp")
	}

	defer db.Close()
	val := db.GetSetting("temp_dir")

	if len(val) > 0 {
		return val
	}

	return filepath.Join(".", ".lovable", "temp")
}

func writeContentToFile(targetPath, content string) error {
	dir := filepath.Dir(targetPath)

	if len(dir) > 0 {
		_ = os.MkdirAll(dir, 0755)
	}

	err := os.WriteFile(targetPath, []byte(content), 0644)

	if err != nil {
		return fmt.Errorf("failed writing to %s: %w", targetPath, err)
	}

	fmt.Printf("  %s✓%s Output written to %s\n", constants.ColorGreen, constants.ColorReset, targetPath)

	return nil
}

func hasArgFlag(args []string, flagName string) bool {
	for _, a := range args {
		if a == flagName || strings.HasPrefix(a, flagName+"=") {
			return true
		}
	}

	return false
}

func extractFlagVal(args []string, flagName string) string {
	for i, arg := range args {
		if arg == flagName && i+1 < len(args) {
			return args[i+1]
		}

		if strings.HasPrefix(arg, flagName+"=") {
			return strings.TrimPrefix(arg, flagName+"=")
		}
	}

	return ""
}

func printJSON(v any) error {
	b, err := json.MarshalIndent(v, "", "  ")

	if err != nil {
		return err
	}

	fmt.Println(string(b))

	return nil
}
