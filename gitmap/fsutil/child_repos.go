// Package fsutil — child_repos.go discovers child Git repositories in a directory.
package fsutil

import (
	"os"
	"path/filepath"
)

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
		gitDir := filepath.Join(subPath, ".git")
		if info, errStat := os.Stat(gitDir); errStat == nil && (info.IsDir() || !info.IsDir()) { // works for worktrees/submodules too
			repos = append(repos, subPath)
		}
	}

	return repos, nil
}
