// Package clonenext — github.go checks and creates GitHub repositories.
package clonenext

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/ghtoken"
)

// RepoExists checks whether a GitHub repository exists via the API.
// Returns true if the repo responds with 200, false on 404, error otherwise.
func RepoExists(owner, repo string) (bool, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return false, fmt.Errorf("create request: %w", err)
	}

	token, _, _ := ghtoken.Resolve()
	if len(token) > 0 {
		req.Header.Set("Authorization", "token "+token)
	}
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("check repo existence: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, nil
	}
	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	return false, fmt.Errorf("GitHub API returned %d checking %s/%s", resp.StatusCode, owner, repo)
}

// CreateRepo creates a new GitHub repository under the given owner.
// It detects whether the owner is a user or organization and calls the
// appropriate endpoint. The repo is created as private by default.
func CreateRepo(owner, repoName string, private bool) error {
	token, _, err := ghtoken.Resolve()
	if err != nil || len(token) == 0 {
		return fmt.Errorf("no GitHub token available — set GITHUB_TOKEN or run `gh auth login`")
	}

	// Try org endpoint first; if 404, fall back to user endpoint.
	err = createOrgRepo(owner, repoName, private, token)
	if err == nil {
		return nil
	}

	// Fallback: create under authenticated user.
	return createUserRepo(repoName, private, token)
}

// createOrgRepo creates a repo under an organization.
func createOrgRepo(org, repoName string, private bool, token string) error {
	url := fmt.Sprintf("https://api.github.com/orgs/%s/repos", org)

	return postCreateRepo(url, repoName, private, token)
}

// createUserRepo creates a repo under the authenticated user.
func createUserRepo(repoName string, private bool, token string) error {
	url := "https://api.github.com/user/repos"

	return postCreateRepo(url, repoName, private, token)
}

// postCreateRepo sends the POST request to create a repository.
func postCreateRepo(apiURL, repoName string, private bool, token string) error {
	payload := map[string]interface{}{
		"name":    repoName,
		"private": private,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal repo payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, apiURL, strings.NewReader(string(jsonData)))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("create repo: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		return nil
	}

	body, _ := io.ReadAll(resp.Body)

	return fmt.Errorf("GitHub API error %d: %s", resp.StatusCode, string(body))
}

// ParseOwnerRepo extracts owner and repo name from a GitHub remote URL.
// Supports both HTTPS and SSH formats.
func ParseOwnerRepo(remoteURL string) (string, string, error) {
	url := strings.TrimSuffix(remoteURL, ".git")

	// HTTPS: https://github.com/owner/repo
	if strings.HasPrefix(url, "https://") {
		parts := strings.Split(url, "/")
		if len(parts) < 5 {
			return "", "", fmt.Errorf("invalid HTTPS remote URL: %s", remoteURL)
		}

		return parts[len(parts)-2], parts[len(parts)-1], nil
	}

	// SSH: git@github.com:owner/repo
	if strings.Contains(url, "@") {
		colonIdx := strings.LastIndex(url, ":")
		if colonIdx < 0 {
			return "", "", fmt.Errorf("invalid SSH remote URL: %s", remoteURL)
		}

		path := url[colonIdx+1:]
		parts := strings.Split(path, "/")
		if len(parts) < 2 {
			return "", "", fmt.Errorf("invalid SSH remote path: %s", remoteURL)
		}

		return parts[0], parts[1], nil
	}

	return "", "", fmt.Errorf("unrecognized remote URL format: %s", remoteURL)
}
