package gitutil

import (
	"strconv"
	"strings"
)

// GetRepoUnpushedCommitCount returns the number of unpushed commits ahead of upstream tracking.
func GetRepoUnpushedCommitCount(repoPath string) (int, error) {
	out, err := RunGitInDir(repoPath, "rev-list", "--count", "@{u}..HEAD")
	if err != nil {
		return 0, err
	}

	clean := strings.TrimSpace(string(out))
	count, parseErr := strconv.Atoi(clean)
	if parseErr != nil {
		return 0, parseErr
	}

	return count, nil
}

// GetRepoMergeBase finds the best common ancestor commit between two references.
func GetRepoMergeBase(repoPath, ref1, ref2 string) (string, error) {
	out, err := RunGitInDir(repoPath, "merge-base", ref1, ref2)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}
