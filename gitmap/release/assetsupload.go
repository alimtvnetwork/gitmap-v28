// Package release — assetsupload.go uploads release assets via GitHub API.
package release

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/verbose"
)

// UploadAsset uploads a single file to a GitHub release.
func UploadAsset(owner, repo string, releaseID int, filePath, token string) error {
	filename := filepath.Base(filePath)
	u := buildGitHubUploadURL(owner, repo, releaseID, filename)

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open asset: %w", err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("stat asset: %w", err)
	}

	req := newGitHubRequest(http.MethodPost, u, file, info.Size())
	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Content-Type", "application/octet-stream")

	resp, err := doGitHubRequest(req)
	if err != nil {
		return fmt.Errorf("upload asset: %w", err)
	}
	defer resp.Body.Close()

	if verbose.IsEnabled() {
		verbose.Get().Log("upload: %s → HTTP %d", filename, resp.StatusCode)
	}

	if resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)

		return &uploadError{
			statusCode: resp.StatusCode,
			message:    string(respBody),
		}
	}

	return nil
}

// uploadError captures HTTP status for retry decisions.
type uploadError struct {
	statusCode int
	message    string
}

// Error implements the error interface.
func (e *uploadError) Error() string {
	return fmt.Sprintf("upload error %d: %s", e.statusCode, e.message)
}

// UploadAllAssets uploads all assets to a GitHub release with exponential backoff retry.
func UploadAllAssets(owner, repo string, releaseID int, assets []string, token string) {
	for _, asset := range assets {
		uploadSingleAsset(owner, repo, releaseID, asset, token)
	}
}

// uploadSingleAsset uploads one asset with retry logic.
func uploadSingleAsset(owner, repo string, releaseID int, asset, token string) {
	filename := filepath.Base(asset)

	if verbose.IsEnabled() {
		info, statErr := os.Stat(asset)
		if statErr == nil {
			verbose.Get().Log("upload-start: %s (%d bytes)", filename, info.Size())
		}
	}

	err := withRetry(filename, constants.RetryMaxAttempts, func() error {
		return UploadAsset(owner, repo, releaseID, asset, token)
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrAssetUploadFinal, filename, err)

		return
	}

	fmt.Printf(constants.MsgAssetUploaded, filename)
}
