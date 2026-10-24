// Package release — githubapi.go provides GitHub API types and helpers.
package release

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// GitHubRelease represents the response from creating a GitHub release.
type GitHubRelease struct {
	ID        int    `json:"id"`
	UploadURL string `json:"upload_url"`
}

// CreateGitHubRelease creates a release via the GitHub API and returns the release ID.
func CreateGitHubRelease(
	owner,
	repo,
	tag,
	name,
	body,
	token string,
	draft bool,
) (*GitHubRelease, error) {
	u := buildGitHubReleasesURL(owner, repo)

	if len(body) > 120000 {
		body = body[:120000] + "\n\n... (truncated for GitHub release limit; see changelog.md in repository)"
	}

	payload := map[string]interface{}{
		"tag_name": tag,
		"name":     name,
		"body":     body,
		"draft":    draft,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, apperror.WrapSimple(err, "marshal release payload")
	}

	req := newGitHubRequest(http.MethodPost, u, io.NopCloser(bytes.NewReader(jsonData)), int64(len(jsonData)))
	setGitHubHeaders(req, token)

	resp, err := doGitHubRequest(req)
	if err != nil {
		return nil, apperror.WrapSimple(err, "create release")
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)

		return nil, apperror.NewExecutionError(fmt.Sprintf("GitHub API error %d: %s", resp.StatusCode, string(respBody)))
	}

	var release GitHubRelease
	terr := json.NewDecoder(resp.Body).Decode(&release)

	return &release, terr
}

// setGitHubHeaders sets common headers for GitHub API requests.
func setGitHubHeaders(req *http.Request, token string) {
	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set(constants.HeaderContentType, constants.MediaTypeJSON)
}
