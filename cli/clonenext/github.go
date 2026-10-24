// Package clonenext — github.go checks and creates GitHub repositories.
package clonenext

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/secrets"
)

// RepoExists checks whether a GitHub repository exists via the API.
// Returns true if the repo responds with 200, false on 404, error otherwise.
func RepoExists(owner, repo string) (bool, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return false, apperror.WrapSimple(err, "create request")
	}

	token, _, _ := secrets.Resolve()
	if len(token) > 0 {
		req.Header.Set("Authorization", "token "+token)
	}

	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, apperror.WrapSimple(err, "check repo existence")
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	return false, apperror.NewExecutionError(fmt.Sprintf("GitHub API returned %d checking %s/%s", resp.StatusCode, owner, repo))
}

// RepoCreateParams encapsulates parameters for creating a GitHub repository.
type RepoCreateParams struct {
	ApiURL    string
	RepoName  string
	IsPrivate bool
	Token     string
}

// CreateRepo creates a new GitHub repository under the given owner.
// It detects whether the owner is a user or organization and calls the
// appropriate endpoint. The repo is created as private by default.
func CreateRepo(owner, repoName string, isPrivate bool) error {
	token, err := resolveRepoToken()
	if err != nil {
		return err
	}
	return attemptRepoCreation(owner, repoName, isPrivate, token)
}

func resolveRepoToken() (string, error) {
	token, _, err := secrets.Resolve()
	if err != nil {
		return "", apperror.WrapSimple(err, "resolve GitHub token failed")
	}
	if len(token) == 0 {
		return "", apperror.NewExecutionError("no GitHub token available — set GITHUB_TOKEN or run `gh auth login`")
	}
	return token, nil
}

func attemptRepoCreation(owner, repoName string, isPrivate bool, token string) error {
	orgParams := RepoCreateParams{
		ApiURL:    fmt.Sprintf("https://api.github.com/orgs/%s/repos", owner),
		RepoName:  repoName,
		IsPrivate: isPrivate,
		Token:     token,
	}
	if err := createOrgRepo(orgParams); err == nil {
		return nil
	}
	userParams := RepoCreateParams{
		ApiURL:    "https://api.github.com/user/repos",
		RepoName:  repoName,
		IsPrivate: isPrivate,
		Token:     token,
	}
	return createUserRepo(userParams)
}

// createOrgRepo creates a repo under an organization.
func createOrgRepo(params RepoCreateParams) error {
	return postCreateRepo(params)
}

// createUserRepo creates a repo under the authenticated user.
func createUserRepo(params RepoCreateParams) error {
	return postCreateRepo(params)
}

// postCreateRepo sends the POST request to create a repository.
func postCreateRepo(params RepoCreateParams) error {
	payload := map[string]interface{}{
		"name":    params.RepoName,
		"private": params.IsPrivate,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return apperror.WrapSimple(err, "marshal repo payload")
	}

	req, err := http.NewRequest(http.MethodPost, params.ApiURL, strings.NewReader(string(jsonData)))
	if err != nil {
		return apperror.WrapSimple(err, "create request")
	}

	req.Header.Set("Authorization", "token "+params.Token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set(constants.HeaderContentType, constants.MediaTypeJSON)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return apperror.WrapSimple(err, "create repo")
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		return nil
	}

	body, _ := io.ReadAll(resp.Body)

	return apperror.NewExecutionError(fmt.Sprintf("GitHub API error %d: %s", resp.StatusCode, string(body)))
}

// ParseOwnerRepo extracts owner and repo name from a GitHub remote URL.
// Supports both HTTPS and SSH formats.
func ParseOwnerRepo(remoteURL string) (string, string, error) {
	url := strings.TrimSuffix(remoteURL, ".git")

	// HTTPS: https://github.com/owner/repo
	if strings.HasPrefix(url, "https://") {
		return parseHttpsOwnerRepo(url, remoteURL)
	}

	// SSH: git@github.com:owner/repo
	if !strings.Contains(url, "@") {
		return "", "", apperror.NewValidationError(fmt.Sprintf("unrecognized remote URL format: %s", remoteURL))
	}

	colonIdx := strings.LastIndex(url, ":")
	if colonIdx < 0 {
		return "", "", apperror.NewValidationError(fmt.Sprintf("invalid SSH remote URL: %s", remoteURL))
	}

	path := url[colonIdx+1:]
	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		return "", "", apperror.NewValidationError(fmt.Sprintf("invalid SSH remote path: %s", remoteURL))
	}

	return parts[0], parts[1], nil
}

func parseHttpsOwnerRepo(url, remoteURL string) (string, string, error) {
	parts := strings.Split(url, "/")
	if len(parts) < 5 {
		return "", "", apperror.NewValidationError(fmt.Sprintf("invalid HTTPS remote URL: %s", remoteURL))
	}

	return parts[len(parts)-2], parts[len(parts)-1], nil
}
