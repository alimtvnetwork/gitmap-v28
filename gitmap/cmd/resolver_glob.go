package cmd

import (
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/fsutil"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

func resolveByGlob(pat string, all []model.ScanRecord) []model.ScanRecord {
	var out []model.ScanRecord
	for _, r := range all {
		base := filepath.Base(fsutil.NormalizeSlashes(r.AbsolutePath))
		if globHit(pat, r.Slug) || globHit(pat, base) {
			out = append(out, r)
		}
	}
	return out
}

func isGlob(s string) bool {
	return strings.ContainsAny(s, "*?[")
}

func globHit(pat, name string) bool {
	ok, err := filepath.Match(pat, name)
	return err == nil && ok
}
