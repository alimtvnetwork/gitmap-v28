package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// docsSiteDownloadSources returns the ordered list of URLs to try when
// auto-fetching docs-site.zip at runtime. We try the current repo's
// versioned tag first, then `latest`, so freshly-built binaries and
// older installs both have a chance.
func docsSiteDownloadSources() []string {
	owner := constants.UpdateRepoOwner
	slug := constants.UpdateCurrentRepoSlug
	asset := constants.DocsSiteArchive
	ver := constants.Version

	return []string{
		fmt.Sprintf("https://github.com/%s/%s/releases/download/v%s/%s", owner, slug, ver, asset),
		fmt.Sprintf("https://github.com/%s/%s/releases/latest/download/%s", owner, slug, asset),
	}
}

// downloadDocsSiteArchive tries each source URL in order and writes the
// first successful response to destPath. Returns the URL that succeeded
// and the number of bytes written, or a wrapped error describing the
// last failure.
func downloadDocsSiteArchive(destPath string) (string, int64, error) {
	sources := docsSiteDownloadSources()

	client := &http.Client{Timeout: time.Duration(constants.DocsSiteDownloadTimeoutSec) * time.Second}

	var lastErr error

	for _, url := range sources {
		fmt.Printf(constants.MsgDocsSiteDownload, url)

		n, err := fetchToFile(client, url, destPath)
		if err == nil {
			return url, n, nil
		}

		lastErr = err

		fmt.Fprintf(os.Stderr, "    skip: %v\n", err)
	}

	return "", 0, fmt.Errorf("tried %d source(s): %w", len(sources), lastErr)
}

// fetchToFile streams an HTTP GET into destPath, capped at maxDocsSiteSize.
func fetchToFile(client *http.Client, url, destPath string) (int64, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, fmt.Errorf("build request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("http get: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("http %d", resp.StatusCode)
	}

	out, err := os.Create(destPath) // #nosec G304 — destPath derived from resolveBinaryDir
	if err != nil {
		return 0, fmt.Errorf("create file: %w", err)
	}
	defer out.Close()

	n, err := io.CopyN(out, resp.Body, maxDocsSiteSize)
	if err != nil && err != io.EOF {
		return n, fmt.Errorf("download body: %w", err)
	}

	if n == 0 {
		return 0, fmt.Errorf("empty response body")
	}

	return n, nil
}
