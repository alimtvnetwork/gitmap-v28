// Package cmd — dir_scope_matcher.go matches working directory against registered work paths.
package cmd

import (
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// IsCurrentDirWorkDir returns true if cwd matches any registered work directory in SQLite.
func IsCurrentDirWorkDir(cwd string) bool {
	db, errDB := store.OpenDefault()
	if errDB != nil {
		return false
	}
	defer db.Close()

	dirs, errList := db.ListWorkDirs()
	if errList != nil {
		return false
	}

	cleanCWD := filepath.Clean(cwd)
	for _, d := range dirs {
		if stringsEqualAbs(cleanCWD, d.AbsolutePath) {
			return true
		}
	}
	return false
}

func stringsEqualAbs(p1, p2 string) bool {
	return strings.EqualFold(filepath.Clean(p1), filepath.Clean(p2))
}
