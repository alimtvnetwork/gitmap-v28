// Package fsutil — workdir_match.go checks if a directory matches a registered work path.
package fsutil

import (
	"path/filepath"
	"strings"
)

// IsInsideWorkDir checks if targetPath is equal to or inside workDirPath.
func IsInsideWorkDir(targetPath, workDirPath string) bool {
	cleanTarget := filepath.Clean(targetPath)
	cleanWork := filepath.Clean(workDirPath)

	if strings.EqualFold(cleanTarget, cleanWork) {
		return true
	}

	rel, err := filepath.Rel(cleanWork, cleanTarget)
	if err != nil {
		return false
	}
	return !strings.HasPrefix(rel, "..") && rel != "."
}
