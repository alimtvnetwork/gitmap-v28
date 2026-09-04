// Package cmd — fixgit_locks.go: detection and removal of stale Git lockfiles.
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func remediateGitLocks(gitDir string, opts FixGitOptions) ([]FixGitIssue, error) {
	var issues []FixGitIssue

	lockFiles, scanErr := findGitLockFiles(gitDir)
	if scanErr != nil {
		return issues, scanErr
	}

	for _, lockPath := range lockFiles {
		issue := processSingleLockFile(gitDir, lockPath, opts.IsDryRun)
		issues = append(issues, issue)
	}

	return issues, nil
}

func findGitLockFiles(gitDir string) ([]string, error) {
	var locks []string

	err := filepath.Walk(gitDir, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return nil
		}

		if !info.IsDir() && isLockFileName(info.Name()) {
			locks = append(locks, path)
		}

		return nil
	})

	return locks, err
}

func isLockFileName(name string) bool {
	return strings.HasSuffix(name, ".lock")
}

func processSingleLockFile(gitDir, lockPath string, isDryRun bool) FixGitIssue {
	relPath, _ := filepath.Rel(gitDir, lockPath)

	issue := FixGitIssue{
		Category:    "Lockfile",
		Description: fmt.Sprintf("Stale lock file found: .git/%s", relPath),
	}

	if isDryRun {
		issue.Remedy = fmt.Sprintf("Would remove stale lock file .git/%s", relPath)

		return issue
	}

	removeErr := os.Remove(lockPath)
	if removeErr != nil {
		issue.ErrorDetail = removeErr.Error()
		issue.Remedy = fmt.Sprintf("Failed to remove lock file .git/%s", relPath)

		return issue
	}

	issue.IsFixed = true
	issue.Remedy = fmt.Sprintf("Removed stale lock file .git/%s", relPath)

	return issue
}
