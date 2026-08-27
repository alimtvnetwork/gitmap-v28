package gitutil

import (
	"os/exec"
	"strings"
)

// GetLatestRemoteBranch quickly gets the most recently updated remote branch name.
func GetLatestRemoteBranch(repoPath string) string {
	cmd := exec.Command("git", "-C", repoPath, "for-each-ref", "--sort=-committerdate", "refs/remotes", "--format=%(refname:lstrip=3)", "--count=5")
	out, err := cmd.Output()
	if err != nil {
		return "—"
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l != "" && l != "HEAD" && l != "main" && l != "master" {
			return l
		}
	}
	return "main" // fallback if nothing else found
}
