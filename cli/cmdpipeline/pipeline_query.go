package cmdpipeline

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/secrets"
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
	tok, _, _ := secrets.Resolve()
	cmd.Env = buildSubprocessEnvWithToken(tok)
}

func buildSubprocessEnvWithToken(tok string) []string {
	var env []string
	for _, entry := range os.Environ() {
		if (len(tok) > 0 && isGHTokenEnvKey(entry)) || isGHPagerEnvKey(entry) {
			continue
		}

		env = append(env, entry)
	}

	if len(tok) > 0 {
		env = append(env, "GH_TOKEN="+tok)
		env = append(env, "GITHUB_TOKEN="+tok)
	}
	env = append(env, "GH_PAGER=cat")
	env = append(env, "PAGER=cat")

	return env
}

func isGHPagerEnvKey(entry string) bool {
	return strings.HasPrefix(entry, "GH_PAGER=") || strings.HasPrefix(entry, "PAGER=")
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

	targetRepo := repo
	if !strings.Contains(repo, "/") {
		targetRepo = ResolvePipelineTarget(repo)
	}
	if strings.Contains(targetRepo, "/") {
		repo = targetRepo
	}

	out, err := runGHCommandWithTimeout("run", "list", "--repo", repo, "--limit", "30", "--json",
		"databaseId,name,status,conclusion,createdAt,updatedAt,headBranch,headSha,url,displayTitle,event")
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Warning: could not list GitHub workflow runs for %s via gh CLI: %v\n", repo, err)

		return queryRunsFromDB(repo)
	}

	return parseRunsOrFallback(out, repo)
}

func parseRunsOrFallback(out []byte, repo string) []ghRunItem {
	var runs []ghRunItem
	if err := json.Unmarshal(out, &runs); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Warning: failed to parse GitHub runs JSON for %s: %v\n", repo, err)

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

	count, parseErr := parseOpenPRsCount(out)
	if parseErr != nil {
		return 0
	}

	return count
}

func parseOpenPRsCount(out []byte) (int, error) {
	var prs []map[string]any
	if err := json.Unmarshal(out, &prs); err != nil {
		return 0, err
	}

	return len(prs), nil
}

func queryLatestTagRelease(repo string) string {
	tag := queryGHLatestTag(repo)
	if len(tag) > 0 {
		return tag
	}

	if ghTag := queryGHTags(repo); len(ghTag) > 0 {
		return ghTag
	}

	if tag := resolveLocalTargetTag(repo); len(tag) > 0 {
		return tag
	}

	if isGitmapRepo(repo) {
		return "v" + constants.Version
	}

	return ""
}

func resolveLocalTargetTag(repo string) string {
	targetDir := resolveRepoLocalDir(repo)
	if len(targetDir) > 0 {
		return queryLocalDirectoryTag(targetDir)
	}
	if isCurrentRepoMatching(repo) {
		return queryLocalDirectoryTag(".")
	}
	return ""
}

func queryLocalDirectoryTag(dir string) string {
	tagOut, err := exec.Command("git", "-C", dir, "describe", "--tags", "--abbrev=0").Output()
	if err == nil && len(tagOut) > 0 {
		return strings.TrimSpace(string(tagOut))
	}
	return readVersionJSONVersion(dir)
}

func queryGHTags(repo string) string {
	if len(repo) == 0 {
		return ""
	}
	endpoint := fmt.Sprintf("repos/%s/tags", repo)
	out, err := runGHCommandWithTimeout("api", endpoint, "--jq", ".[0].name")
	if err != nil || len(out) == 0 {
		return ""
	}
	tag := strings.TrimSpace(string(out))
	if tag == "" || tag == "null" {
		return ""
	}
	return tag
}

func isGitmapRepo(repo string) bool {
	low := strings.ToLower(repo)
	return low == "gitmap" || strings.Contains(low, "gitmap-v") || strings.HasSuffix(low, "/gitmap")
}

func isCurrentRepoMatching(repo string) bool {
	cwdSlug := resolveCurrentRepoSlug()
	return strings.EqualFold(cwdSlug, repo) || strings.EqualFold(filepath.Base(cwdSlug), filepath.Base(repo))
}

func resolveRepoLocalDir(repo string) string {
	_, absPath, ok := queryRepoSlugAndPathFromDB(repo)
	if !ok || absPath == "" {
		return ""
	}
	info, err := os.Stat(absPath)
	if err != nil || !info.IsDir() {
		return ""
	}
	return absPath
}

type versionFileMinimal struct {
	Version string `json:"version"`
}

func readVersionJSONVersion(dir string) string {
	data, err := os.ReadFile(filepath.Join(dir, "version.json"))
	if err != nil {
		return ""
	}
	var vf versionFileMinimal
	if err := json.Unmarshal(data, &vf); err != nil || vf.Version == "" {
		return ""
	}
	if !strings.HasPrefix(vf.Version, "v") {
		return "v" + vf.Version
	}
	return vf.Version
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
