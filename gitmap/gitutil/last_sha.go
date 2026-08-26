// Package gitutil — last_sha.go resolves short commit SHA of HEAD.
package gitutil

import (
	"os/exec"
	"strings"
)

// GetLastCommitSHA returns the 7-character short commit hash for HEAD in repoPath.
func GetLastCommitSHA(repoPath string) string {
	cmd := exec.Command("git", "-C", repoPath, "rev-parse", "--short=7", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return "-"
	}
	return strings.TrimSpace(string(out))
}
