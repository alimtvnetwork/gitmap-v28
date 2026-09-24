// Package cmd — version_tags.go handles querying, parsing, and rendering GitHub release tags for GitMap and AGM.
package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// GitHubReleaseTagInfo represents a GitHub release tag with metadata and active state.
type GitHubReleaseTagInfo struct {
	TagName      string `json:"tag_name"`
	Name         string `json:"name"`
	PublishedAt  string `json:"published_at"`
	IsPrerelease bool   `json:"prerelease"`
	IsDraft      bool   `json:"draft"`
	IsActive     bool   `json:"is_active"`
	IsLatest     bool   `json:"is_latest"`
	Status       string `json:"status"`
}

// FetchGitHubReleaseTags retrieves release tags from GitHub API with an HTTP timeout.
func FetchGitHubReleaseTags(owner, repo string, limit int) ([]GitHubReleaseTagInfo, error) {
	if limit <= 0 {
		limit = 30
	}
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases?per_page=%d", owner, repo, limit)
	client := &http.Client{Timeout: 8 * time.Second}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gitmap-version-checker")
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		if resp != nil {
			_ = resp.Body.Close()
		}
		return fetchFallbackTags(owner, repo, limit)
	}
	defer resp.Body.Close()

	var rawReleases []struct {
		TagName      string `json:"tag_name"`
		Name         string `json:"name"`
		PublishedAt  string `json:"published_at"`
		IsPrerelease bool   `json:"prerelease"`
		IsDraft      bool   `json:"draft"`
	}
	if decodeErr := json.NewDecoder(resp.Body).Decode(&rawReleases); decodeErr != nil {
		return nil, decodeErr
	}
	var out []GitHubReleaseTagInfo
	for i, r := range rawReleases {
		pubDate := r.PublishedAt
		if len(pubDate) >= 10 {
			pubDate = pubDate[:10]
		}
		out = append(out, GitHubReleaseTagInfo{
			TagName:      r.TagName,
			Name:         r.Name,
			PublishedAt:  pubDate,
			IsPrerelease: r.IsPrerelease,
			IsDraft:      r.IsDraft,
			IsLatest:     i == 0,
		})
	}
	if len(out) == 0 {
		return fetchFallbackTags(owner, repo, limit)
	}
	return out, nil
}

func fetchFallbackTags(owner, repo string, limit int) ([]GitHubReleaseTagInfo, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/tags?per_page=%d", owner, repo, limit)
	client := &http.Client{Timeout: 8 * time.Second}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gitmap-version-checker")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github api status: %d", resp.StatusCode)
	}
	var rawTags []struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&rawTags); err != nil {
		return nil, err
	}
	var out []GitHubReleaseTagInfo
	for i, t := range rawTags {
		out = append(out, GitHubReleaseTagInfo{
			TagName:  t.Name,
			IsLatest: i == 0,
		})
	}
	return out, nil
}

// CompareSemver compares two semantic versions (e.g. "v6.326.0" vs "v6.324.0").
// Returns 1 if v1 > v2, -1 if v1 < v2, and 0 if equal.
func CompareSemver(v1, v2 string) int {
	clean1 := strings.TrimPrefix(strings.TrimSpace(v1), "v")
	clean2 := strings.TrimPrefix(strings.TrimSpace(v2), "v")
	if clean1 == clean2 {
		return 0
	}
	p1 := strings.Split(clean1, ".")
	p2 := strings.Split(clean2, ".")
	maxLen := len(p1)
	if len(p2) > maxLen {
		maxLen = len(p2)
	}
	for i := 0; i < maxLen; i++ {
		var n1, n2 int
		if i < len(p1) {
			_, _ = fmt.Sscanf(p1[i], "%d", &n1)
		}
		if i < len(p2) {
			_, _ = fmt.Sscanf(p2[i], "%d", &n2)
		}
		if n1 > n2 {
			return 1
		}
		if n1 < n2 {
			return -1
		}
	}
	return 0
}

// RenderReleaseTagsTable prints a formatted table of release tags with status relative to current version.
func RenderReleaseTagsTable(title, currentVersion string, tags []GitHubReleaseTagInfo) {
	fmt.Printf("\n%s=== %s Releases & Version Tags ===%s\n", constants.ColorCyan, title, constants.ColorReset)
	fmt.Printf("Current Installed/Running Version: %s%s%s\n\n", constants.ColorGreen, currentVersion, constants.ColorReset)
	fmt.Println("┌──────────────────┬──────────────┬──────────────────────────────┐")
	fmt.Println("│ RELEASE TAG      │ DATE         │ STATUS                       │")
	fmt.Println("├──────────────────┼──────────────┼──────────────────────────────┤")
	if len(tags) == 0 {
		fmt.Println("│ (no tags found)  │ -            │ -                            │")
	}
	for _, t := range tags {
		cmp := CompareSemver(t.TagName, currentVersion)
		status := "▼ older version"
		if cmp == 0 {
			status = "● CURRENTLY ACTIVE"
		} else if cmp > 0 {
			status = "▲ NEWER RELEASE"
		}
		if t.IsLatest && cmp == 0 {
			status = "● CURRENTLY ACTIVE (LATEST)"
		} else if t.IsLatest {
			status = "▲ LATEST RELEASE"
		}
		pub := t.PublishedAt
		fmt.Printf("│ %-16s │ %-12s │ %-28s │\n", t.TagName, pub, status)
	}
	fmt.Println("└──────────────────┴──────────────┴──────────────────────────────┘")
	fmt.Println()
}

// RunGitMapVersionTagsLS queries and renders all GitHub release tags for GitMap.
func RunGitMapVersionTagsLS() error {
	fmt.Printf("Fetching available GitMap release tags from GitHub (%s/%s)...\n", constants.UpdateRepoOwner, constants.UpdateCurrentRepoSlug)
	tags, err := FetchGitHubReleaseTags(constants.UpdateRepoOwner, constants.UpdateCurrentRepoSlug, 30)
	if err != nil {
		fmt.Printf("Note: could not query GitHub API directly (%v); showing current running version.\n", err)
	}
	RenderReleaseTagsTable("GitMap CLI (gitmap-v28)", "v"+constants.Version, tags)
	fmt.Println("To install/switch to an older or specific GitMap release:")
	fmt.Println("  gitmap update <version>         (e.g. gitmap update v6.324.0)")
	fmt.Println("  gitmap update --version <ver>   (e.g. gitmap update --version v6.325.0)")
	return nil
}

// RunAGMVersionTagsLS queries and renders all GitHub release tags for Antigravity Manager.
func RunAGMVersionTagsLS() error {
	const owner = "alimtvnetwork"
	const repo = "Antigravity-Manager"
	fmt.Printf("Fetching available Antigravity Manager release tags from GitHub (%s/%s)...\n", owner, repo)
	tags, err := FetchGitHubReleaseTags(owner, repo, 30)
	if err != nil {
		fmt.Printf("Note: could not query GitHub API directly (%v).\n", err)
	}
	RenderReleaseTagsTable("Antigravity Manager (AGM)", "v4.7.1", tags)
	fmt.Println("To install/update to a specific Antigravity Manager release:")
	fmt.Println("  gitmap agm install <version>    (e.g. gitmap agm install 4.7.1)")
	fmt.Println("  gitmap agm update <version>     (e.g. gitmap agm update 4.7.0)")
	return nil
}
