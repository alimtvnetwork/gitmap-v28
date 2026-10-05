// Package fsutil — child_repos.go discovers child Git repositories in a directory.
package fsutil

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// isDefaultScanExcluded checks whether a directory name matches DefaultScanExcludeDirs.
func isDefaultScanExcluded(name string) bool {
	lower := strings.ToLower(name)
	for _, excl := range constants.DefaultScanExcludeDirs {
		if strings.EqualFold(lower, excl) {
			return true
		}
	}

	return false
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

		if isDefaultScanExcluded(entry.Name()) {
			continue
		}

		subPath := filepath.Join(parentDir, entry.Name())
		gitDir := filepath.Join(subPath, ".git")
		if _, errStat := os.Stat(gitDir); errStat == nil { // works for worktrees/submodules too
			repos = append(repos, subPath)
		}
	}

	return repos, nil
}
