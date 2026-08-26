// Package cmd — prompt_target_resolver.go resolves directory targets for prompt installation.
package cmd

import (
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/fsutil"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// ResolvePromptTarget resolves a path, alias, ID, or discovers child repos.
func ResolvePromptTarget(target string) ([]string, error) {
	if target != "" {
		abs, err := filepath.Abs(target)
		if err == nil {
			if info, errStat := os.Stat(abs); errStat == nil && info.IsDir() {
				// If directory is a git repo, return it directly
				if fsutil.IsValidGitDir(abs) {
					return []string{abs}, nil
				}
				// If directory has child git repos, discover them
				if childRepos, errDisc := fsutil.DiscoverTopLevelGitRepos(abs); errDisc == nil && len(childRepos) > 0 {
					return childRepos, nil
				}
				return []string{abs}, nil
			}
		}

		// Try DB lookup by slug or ID
		db, errDB := store.OpenDefault()
		if errDB == nil {
			defer db.Close()
			if repos, errList := db.ListRepos(); errList == nil {
				for _, r := range repos {
					if r.Slug == target || r.RepoName == target {
						return []string{r.AbsolutePath}, nil
					}
				}
			}
		}
	}

	// Fallback to CWD
	cwd, _ := os.Getwd()
	if fsutil.IsValidGitDir(cwd) {
		return []string{cwd}, nil
	}
	if childRepos, errDisc := fsutil.DiscoverTopLevelGitRepos(cwd); errDisc == nil && len(childRepos) > 0 {
		return childRepos, nil
	}
	return []string{cwd}, nil
}
