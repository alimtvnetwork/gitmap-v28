package gitutil

import (
	"os/exec"
	"strings"
)

// GetLatestTag queries the most recent tag in the repository.
func GetLatestTag(repoPath string) string {
	hasPath := len(repoPath) > 0
	if hasPath == false {
		return ""
	}

	tag := queryDescribeTag(repoPath)
	if len(tag) > 0 {
		return tag
	}

	return querySortTag(repoPath)
}

func queryDescribeTag(repoPath string) string {
	cmd := exec.Command("git", "-C", repoPath, "describe", "--tags", "--abbrev=0")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}

	tag := strings.TrimSpace(string(out))
	hasTag := len(tag) > 0
	if hasTag {
		return tag
	}

	return ""
}

func querySortTag(repoPath string) string {
	cmd := exec.Command("git", "-C", repoPath, "tag", "-l", "v*", "--sort=-v:refname")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		hasLine := len(trimmed) > 0
		if hasLine {
			return trimmed
		}
	}

	return ""
}
