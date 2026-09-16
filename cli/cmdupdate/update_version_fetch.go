package cmdupdate

import (
	"fmt"
	"net/http"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// fetchRemoteTargetVersion fetches the latest version from GitHub releases or version.json.
func fetchRemoteTargetVersion(slug string) string {
	ghVer := fetchGitHubLatestReleaseVersion(slug)
	hasValidVer := len(ghVer) > 0 && ghVer != constants.VersionUnknown

	if hasValidVer {
		return ghVer
	}

	return fetchVersionJSON(slug)
}

func fetchGitHubLatestReleaseVersion(slug string) string {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", constants.UpdateRepoOwner, slug)
	client := &http.Client{Timeout: 6 * time.Second}
	resp, err := executeGitHubReleaseRequest(client, url)
	isFailed := err != nil || resp.StatusCode != http.StatusOK

	if isFailed {
		closeResponse(resp)

		return ""
	}

	defer resp.Body.Close()

	return decodeReleaseTagName(resp.Body)
}

func executeGitHubReleaseRequest(client *http.Client, url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", constants.UpdaterBin)

	return client.Do(req)
}

func fetchVersionJSON(slug string) string {
	url := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/main/version.json", constants.UpdateRepoOwner, slug)
	client := &http.Client{Timeout: 6 * time.Second}
	resp, err := client.Get(url)
	isFailed := err != nil || resp.StatusCode != http.StatusOK

	if isFailed {
		closeResponse(resp)

		return constants.VersionUnknown
	}

	defer resp.Body.Close()

	return decodeVersionFromMap(resp.Body)
}

func closeResponse(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		_ = resp.Body.Close()
	}
}
