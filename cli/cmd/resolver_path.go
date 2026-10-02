package cmd

import (
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/fsutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

func resolveByPath(target string, all []model.ScanRecord) *model.ScanRecord {
	cleanTarget := fsutil.NormalizeSlashes(fsutil.TrimTrailingSlashes(target))
	for _, r := range all {
		if fsutil.EqualPaths(r.AbsolutePath, cleanTarget) {
			return &r
		}
	}

	if found := findByCleanTargetAbs(cleanTarget, all); found != nil {
		return found
	}

	for _, r := range all {
		rClean := fsutil.NormalizeSlashes(r.AbsolutePath)
		if strings.EqualFold(filepath.Base(rClean), filepath.Base(cleanTarget)) || strings.EqualFold(r.Slug, filepath.Base(cleanTarget)) {
			return &r
		}
	}

	return nil
}

func findByAbsolutePath(abs string, all []model.ScanRecord) *model.ScanRecord {
	for _, r := range all {
		if fsutil.EqualPaths(r.AbsolutePath, abs) {
			return &r
		}
	}

	return nil
}

func findByCleanTargetAbs(cleanTarget string, all []model.ScanRecord) *model.ScanRecord {
	abs, err := filepath.Abs(cleanTarget)
	if err != nil {
		return nil
	}

	return findByAbsolutePath(abs, all)
}
