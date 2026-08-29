// Package fsutil — recursive_top_level.go scans directories and prunes traversal upon finding .git.
package fsutil

import (
	"io/fs"
	"os"
	"path/filepath"
)

// DiscoverTopLevelGitRepos walks rootDir recursively and collects top-level Git repositories.
// Crucial: When a directory contains `.git`, it is added to the result and NOT descended into further.
func DiscoverTopLevelGitRepos(rootDir string) ([]string, error) {
	var repos []string

	info, err := os.Stat(rootDir)
	if err != nil || !info.IsDir() {
		return nil, err
	}

	errWalk := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return nil // ignore inaccessible directories gracefully
		}

		if !d.IsDir() {
			return nil
		}

		// Skip hidden tool directories (e.g. .git itself, .lovable, node_modules)
		name := d.Name()
		if path != rootDir && (name == ".git" || name == "node_modules" || name == ".cache" || name == "vendor") {
			return filepath.SkipDir
		}

		// Check if this directory is a Git repository
		gitPath := filepath.Join(path, ".git")
		if _, errGit := os.Stat(gitPath); errGit != nil {
			return nil
		}

		repos = append(repos, path)
		if path != rootDir {
			return filepath.SkipDir
		}

		return nil
	})

	return repos, errWalk
}
