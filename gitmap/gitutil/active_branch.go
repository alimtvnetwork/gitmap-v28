// Package gitutil — active_branch.go resolves active branch name or detached state.
package gitutil

import (
	"os/exec"
	"strings"
)

// GetActiveBranch returns the current active branch name or "(detached)" if in detached HEAD.
func GetActiveBranch(repoPath string) string {
	cmd := exec.Command("git", "-C", repoPath, "rev-parse", "--abbrev-ref", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return "-"
	}
	branch := strings.TrimSpace(string(out))
	if branch == "HEAD" {
		return "(detached)"
	}
	return branch
}
