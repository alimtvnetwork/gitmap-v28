// Package dashboard collects Git repository data for the HTML dashboard.
package dashboard

import (
	"os/exec"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// queryBranches returns raw branch lines from git branch -a --format.
func queryBranches(repoPath string) ([]string, error) {
	out, err := runDashGit(repoPath,
		"branch", "-a", "--format="+constants.GitBranchDashFormat)
	if err != nil {
		return nil, err
	}

	return splitNonEmpty(out), nil
}

// queryTags returns raw tag lines from git tag --sort=-creatordate --format.
func queryTags(repoPath string) ([]string, error) {
	out, err := runDashGit(repoPath,
		"tag", "--sort=-creatordate", "--format="+constants.GitTagDashFormat)
	if err != nil {
		return nil, err
	}

	return splitNonEmpty(out), nil
}

// queryLog returns raw commit lines from git log with --numstat.
// An optional limit caps the number of commits returned.
// An optional since restricts commits to those after a date string.
func queryLog(repoPath string, limit int, since string, noMerges bool) (string, error) {
	args := []string{"log", "--format=" + constants.GitLogDashFormat, "--numstat"}

	if limit > 0 {
		args = append(args, "-n", strconv.Itoa(limit))
	}
	if len(since) > 0 {
		args = append(args, "--since="+since)
	}
	if noMerges {
		args = append(args, "--no-merges")
	}

	out, err := runDashGit(repoPath, args...)
	if err != nil {
		return "", err
	}

	return out, nil
}

// queryTagDistance returns the commit count between two refs (e.g. v1.0..v1.1).
func queryTagDistance(repoPath, fromRef, toRef string) int {
	out, err := runDashGit(repoPath,
		"rev-list", "--count", fromRef+".."+toRef)
	if err != nil {
		return 0
	}

	count, _ := strconv.Atoi(strings.TrimSpace(out))

	return count
}

// queryRepoName extracts the repository name from the remote URL.
// Falls back to the directory base name if no remote is configured.
func queryRepoName(repoPath string) string {
	out, err := runDashGit(repoPath,
		"config", "--get", "remote.origin.url")
	if err != nil {
		return baseName(repoPath)
	}

	url := strings.TrimSpace(out)
	url = strings.TrimSuffix(url, ".git")

	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}

	return baseName(repoPath)
}

// runDashGit executes a git command in the given directory.
func runDashGit(dir string, args ...string) (string, error) {
	cmd := exec.Command(constants.GitBin, args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}

// splitNonEmpty splits output by newline and drops empty lines.
func splitNonEmpty(s string) []string {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	result := make([]string, 0, len(lines))

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if len(trimmed) > 0 {
			result = append(result, trimmed)
		}
	}

	return result
}

// baseName returns the last path segment.
func baseName(path string) string {
	path = strings.TrimRight(path, "/\\")
	idx := strings.LastIndexAny(path, "/\\")
	if idx >= 0 {
		return path[idx+1:]
	}

	return path
}
