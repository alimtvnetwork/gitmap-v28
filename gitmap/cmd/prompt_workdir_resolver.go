// Package cmd — prompt_workdir_resolver.go extracts targets from registered work directories.
package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

func ResolveAllWorkDirPromptTargets() ([]string, error) {
	db, errDB := store.OpenDefault()
	if errDB != nil {
		return nil, errDB
	}
	defer db.Close()

	dirs, errList := db.ListWorkDirs()
	if errList != nil {
		return nil, errList
	}

	var allTargets []string
	for _, d := range dirs {
		if childRepos, err := DiscoverPromptChildRepos(d.AbsolutePath); err == nil && len(childRepos) > 0 {
			allTargets = append(allTargets, childRepos...)
		} else {
			allTargets = append(allTargets, d.AbsolutePath)
		}
	}
	return allTargets, nil
}
