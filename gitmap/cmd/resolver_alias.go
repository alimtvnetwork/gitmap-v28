package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/fsutil"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

func resolveByAlias(db *store.DB, target string, all []model.ScanRecord) *model.ScanRecord {
	aliasRow, err := db.ResolveAlias(target)
	if err != nil {
		return nil
	}
	for _, r := range all {
		if r.ID == aliasRow.RepoID || fsutil.EqualPaths(r.AbsolutePath, aliasRow.AbsolutePath) {
			return &r
		}
	}
	return nil
}
