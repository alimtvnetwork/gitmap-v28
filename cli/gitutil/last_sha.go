// Package gitutil — last_sha.go resolves short commit SHA of HEAD.
package gitutil

import (
	"os"
	"os/exec"
	"strings"
)

// GetLastCommitSHA returns the 7-character short commit hash for HEAD in repoPath.
// Returns "-" immediately if repoPath is empty, unresolvable, or not a git directory.
func GetLastCommitSHA(repoPath string) string {
	cleanPath := strings.TrimSpace(repoPath)
	if cleanPath == "" {
		return "-"
	}

	fi, err := os.Stat(cleanPath)
	isDir := err == nil && fi.IsDir()
	if !isDir {
		return "-"
	}

	cmd := exec.Command("git", "-C", cleanPath, "rev-parse", "--short=7", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return "-"
	}

	return strings.TrimSpace(string(out))
}
