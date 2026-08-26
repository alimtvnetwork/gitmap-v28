// Package cmd — pull_target_resolver.go resolves repositories for pull with recursive top-level discovery.
package cmd

import (
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/fsutil"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// ResolvePullDirectoryTargets resolves scan records from a work directory or current path.
func ResolvePullDirectoryTargets(dirPath string) []model.ScanRecord {
	discovered, err := fsutil.DiscoverTopLevelGitRepos(dirPath)
	if err != nil || len(discovered) == 0 {
		return nil
	}

	var records []model.ScanRecord
	db, errDB := store.OpenDefault()
	var dbRepos []model.ScanRecord
	if errDB == nil {
		dbRepos, _ = db.ListRepos()
		db.Close()
	}

	for _, d := range discovered {
		name := filepath.Base(d)
		rec := model.ScanRecord{
			RepoName:     name,
			Slug:         name,
			AbsolutePath: d,
		}
		for _, dbR := range dbRepos {
			if stringsEqualAbs(dbR.AbsolutePath, d) {
				rec = dbR
				break
			}
		}
		records = append(records, rec)
	}
	return records
}
