package cmd

import (
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/fsutil"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

func resolveByPath(target string, all []model.ScanRecord) *model.ScanRecord {
	cleanTarget := fsutil.NormalizeSlashes(fsutil.TrimTrailingSlashes(target))
	for _, r := range all {
		if fsutil.EqualPaths(r.AbsolutePath, cleanTarget) {
			return &r
		}
	}

	abs, err := filepath.Abs(cleanTarget)
	if err == nil {
		for _, r := range all {
			if fsutil.EqualPaths(r.AbsolutePath, abs) {
				return &r
			}
		}
	}

	baseName := strings.ToLower(filepath.Base(cleanTarget))
	for _, r := range all {
		rClean := fsutil.NormalizeSlashes(r.AbsolutePath)
		rBase := strings.ToLower(filepath.Base(rClean))
		if rBase == baseName || strings.ToLower(r.Slug) == baseName {
			return &r
		}
	}
	return nil
}
