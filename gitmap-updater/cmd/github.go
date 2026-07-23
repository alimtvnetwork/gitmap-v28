package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
)

// releaseResponse is a minimal GitHub release API response.
type releaseResponse struct {
	TagName string `json:"tag_name"`
}

// fetchLatestTag queries the GitHub releases API for the latest tag.
func fetchLatestTag() (string, error) {
	req, err := http.NewRequest("GET", GitHubAPILatest, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "gitmap-updater/"+Version)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("network error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))

		return "", fmt.Errorf(
			"GitHub API returned HTTP %d\n"+
				"  URL: %s\n"+
				"  Response: %s\n"+
				"  Possible causes:\n"+
				"    - No published releases in the repository\n"+
				"    - Repository is private (needs authentication)\n"+
				"    - Repository name has changed\n"+
				"  Try: https://github.com/%s/releases",
			resp.StatusCode, GitHubAPILatest, string(body), RepoSlug,
		)
	}

	var release releaseResponse
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", fmt.Errorf("failed to parse release JSON: %w", err)
	}

	return release.TagName, nil
}

// getInstalledVersion runs `gitmap version` and returns the output.
func getInstalledVersion() (string, error) {
	cmd := exec.Command(GitMapBin, "version")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

// normalizeVersion strips the "v" prefix for comparison.
func normalizeVersion(v string) string {
	return strings.TrimPrefix(strings.TrimSpace(v), "v")
}
