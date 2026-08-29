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
	if target == "" {
		return resolveFallbackPromptTarget(), nil
	}
	if repos, ok := resolveDirPromptTarget(target); ok {
		return repos, nil
	}
	if repos, ok := resolveDBPromptTarget(target); ok {
		return repos, nil
	}
	return resolveFallbackPromptTarget(), nil
}

func resolveDirPromptTarget(target string) ([]string, bool) {
	abs, err := filepath.Abs(target)
	if err != nil {
		return nil, false
	}
	info, errStat := os.Stat(abs)
	if errStat != nil || !info.IsDir() {
		return nil, false
	}
	if fsutil.IsValidGitDir(abs) {
		return []string{abs}, true
	}
	if childRepos, errDisc := fsutil.DiscoverTopLevelGitRepos(abs); errDisc == nil && len(childRepos) > 0 {
		return childRepos, true
	}
	return []string{abs}, true
}

func resolveDBPromptTarget(target string) ([]string, bool) {
	db, errDB := store.OpenDefault()
	if errDB != nil {
		return nil, false
	}
	defer db.Close()

	repos, errList := db.ListRepos()
	if errList != nil {
		return nil, false
	}
	for _, r := range repos {
		if r.Slug == target || r.RepoName == target {
			return []string{r.AbsolutePath}, true
		}
	}
	return nil, false
}

func resolveFallbackPromptTarget() []string {
	cwd, _ := os.Getwd()
	if fsutil.IsValidGitDir(cwd) {
		return []string{cwd}
	}
	if childRepos, errDisc := fsutil.DiscoverTopLevelGitRepos(cwd); errDisc == nil && len(childRepos) > 0 {
		return childRepos
	}
	return []string{cwd}
}
