// Package cmdfixreleasetags provides tag discovery, GitHub release inspection,
// asset integrity checking, and CI/CD workflow auditing.
package cmdfixreleasetags

import (
	"encoding/json"
	"os"
	"strings"
	"time"
)

type ghReleaseItem struct {
	TagName   string        `json:"tagName"`
	IsDraft   bool          `json:"isDraft"`
	CreatedAt string        `json:"createdAt"`
	Assets    []ghAssetItem `json:"assets"`
}

type ghAssetItem struct {
	Name  string `json:"name"`
	Size  int64  `json:"size"`
	State string `json:"state"`
}

type ghRunItem struct {
	DatabaseId uint64 `json:"databaseId"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	Conclusion string `json:"conclusion"`
	CreatedAt  string `json:"createdAt"`
}

// CollectLocalTags reads local git tags via git tag -l.
func CollectLocalTags(repoPath string, executor CommandExecutor) ([]string, error) {
	out, err := executor.Run(repoPath, "git", "tag", "-l")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(out), "\n")
	var tags []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			tags = append(tags, trimmed)
		}
	}

	return tags, nil
}

// CollectRemoteTags reads remote git tags via git ls-remote --tags origin.
func CollectRemoteTags(repoPath string, executor CommandExecutor) (map[string]bool, error) {
	out, err := executor.Run(repoPath, "git", "ls-remote", "--tags", "origin")
	if err != nil {
		return nil, err
	}

	tags := make(map[string]bool)
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		processRemoteTagLine(line, tags)
	}

	return tags, nil
}

func processRemoteTagLine(line string, tags map[string]bool) {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return
	}

	parts := strings.Fields(trimmed)
	if len(parts) < 2 || strings.HasSuffix(parts[1], "^{}") {
		return
	}

	const prefix = "refs/tags/"
	if strings.HasPrefix(parts[1], prefix) {
		tags[strings.TrimPrefix(parts[1], prefix)] = true
	}
}

// FetchGitHubReleases retrieves repository releases via gh release list.
func FetchGitHubReleases(repoPath string, executor CommandExecutor) ([]ghReleaseItem, error) {
	if os.Getenv("GITMAP_MOCK_GH") == "1" {
		return nil, nil
	}

	out, err := executor.Run(repoPath, "gh", "release", "list", "--limit", "100", "--json", "tagName,isDraft,createdAt,assets")
	if err != nil {
		return nil, err
	}

	trimmed := strings.TrimSpace(string(out))
	if trimmed == "" || trimmed == "[]" {
		return nil, nil
	}

	var items []ghReleaseItem
	err = json.Unmarshal([]byte(trimmed), &items)
	if err != nil {
		return nil, err
	}

	return items, nil
}

// FetchTagWorkflows retrieves CI/CD workflow runs for a specific tag or commit.
func FetchTagWorkflows(repoPath, tag, commitSha string, executor CommandExecutor) ([]CIWorkflowRunInfo, error) {
	if os.Getenv("GITMAP_MOCK_GH") == "1" {
		return nil, nil
	}

	args := buildWorkflowQueryArgs(tag, commitSha)
	out, err := executor.Run(repoPath, "gh", args...)
	if err != nil {
		return nil, err
	}

	return parseWorkflowRunItems(out)
}

func buildWorkflowQueryArgs(tag, commitSha string) []string {
	args := []string{"run", "list", "--limit", "10", "--json", "databaseId,name,status,conclusion,createdAt"}
	if commitSha != "" {
		return append(args, "--commit", commitSha)
	}

	if tag != "" {
		return append(args, "--branch", tag)
	}

	return args
}

func parseWorkflowRunItems(data []byte) ([]CIWorkflowRunInfo, error) {
	trimmed := strings.TrimSpace(string(data))
	if trimmed == "" || trimmed == "[]" {
		return nil, nil
	}

	var rawItems []ghRunItem
	err := json.Unmarshal([]byte(trimmed), &rawItems)
	if err != nil {
		return nil, err
	}

	var runs []CIWorkflowRunInfo
	for _, r := range rawItems {
		parsedTime, _ := time.Parse(time.RFC3339, r.CreatedAt)
		isSucc := strings.EqualFold(r.Conclusion, "success")
		isInProg := isRunInProgress(CIWorkflowRunInfo{Status: r.Status})
		runs = append(runs, CIWorkflowRunInfo{
			RunId:        r.DatabaseId,
			WorkflowName: r.Name,
			Status:       r.Status,
			Conclusion:   r.Conclusion,
			CreatedAt:    parsedTime,
			IsSuccessful: isSucc,
			IsInProgress: isInProg,
		})
	}

	return runs, nil
}
