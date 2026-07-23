// Package release — remoteorigin.go parses git remote origin URLs.
package release

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// ParseRemoteOrigin extracts owner/repo from the git remote origin URL.
func ParseRemoteOrigin() (string, string, error) {
	url := getRemoteURL()
	if len(url) == 0 {
		return "", "", fmt.Errorf("no remote origin URL found")
	}

	return parseGitURL(url)
}

// getRemoteURL reads the origin remote URL via git config.
func getRemoteURL() string {
	out, err := gitOutput("config", "--get", "remote.origin.url")
	if err != nil {
		return ""
	}

	return strings.TrimSpace(out)
}

// parseGitURL extracts owner/repo from HTTPS or SSH git URLs.
// ParseGitURLExported is the exported alias for testing.
func ParseGitURLExported(url string) (string, string, error) {
	return parseGitURL(url)
}

// parseGitURL extracts owner/repo from HTTPS or SSH git URLs.
func parseGitURL(url string) (string, string, error) {
	// HTTPS: https://github.com/owner/repo.git
	if strings.HasPrefix(url, "https://") {
		return parseHTTPSURL(url)
	}

	// SSH: git@github.com:owner/repo.git
	if strings.Contains(url, "@") {
		return parseSSHURL(url)
	}

	return "", "", fmt.Errorf("unrecognized remote URL format: %s", url)
}

// parseHTTPSURL parses https://github.com/owner/repo.git
func parseHTTPSURL(url string) (string, string, error) {
	url = strings.TrimSuffix(url, ".git")
	parts := strings.Split(url, "/")

	if len(parts) < 5 {
		return "", "", fmt.Errorf("invalid HTTPS remote: %s", url)
	}

	owner := parts[len(parts)-2]
	repo := parts[len(parts)-1]

	return owner, repo, nil
}

// parseSSHURL parses git@github.com:owner/repo.git
func parseSSHURL(url string) (string, string, error) {
	url = strings.TrimSuffix(url, ".git")
	colonIdx := strings.LastIndex(url, ":")

	if colonIdx < 0 {
		return "", "", fmt.Errorf("invalid SSH remote: %s", url)
	}

	path := url[colonIdx+1:]
	parts := strings.Split(path, "/")

	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid SSH remote path: %s", url)
	}

	return parts[0], parts[1], nil
}

// gitOutput runs a git command and returns stdout.
func gitOutput(args ...string) (string, error) {
	cmd := exec.Command(constants.GitBin, args...)
	out, err := cmd.Output()

	return string(out), err
}
