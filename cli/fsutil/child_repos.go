// Package fsutil — child_repos.go discovers child Git repositories in a directory.
package fsutil

import (
	"os"
	"path/filepath"
)

// isDefaultScanExcluded checks whether a directory name matches DefaultScanExcludeDirs.
func isDefaultScanExcluded(name string) bool {
	return IsDefaultScanExcluded(name)
}

// DiscoverChildGitRepos scans immediate subdirectories of parentDir for .git folders.
func DiscoverChildGitRepos(parentDir string) ([]string, error) {
	var repos []string

	entries, err := os.ReadDir(parentDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		subPath := filepath.Join(parentDir, entry.Name())
		if isDefaultScanExcluded(entry.Name()) || IsExcludedBelowScanRoot(parentDir, subPath) {
			continue
		}

		if IsValidGitRepo(subPath) {
			repos = append(repos, subPath)
		}
	}

	return repos, nil
}
