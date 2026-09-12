// Package fsutil — copydir.go provides directory copying utilities.
package fsutil

import (
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// CopyDirContents copies all files from src directory to dest directory.
func CopyDirContents(src string, dest string) *apperror.AppError {
	if mkErr := os.MkdirAll(dest, constants.DirPermission); mkErr != nil {
		return apperror.WrapSimple(mkErr, "mkdir dest")
	}

	entries, readErr := os.ReadDir(src)
	if readErr != nil {
		return nil
	}

	return copyEntries(src, dest, entries)
}

func copyEntries(src string, dest string, entries []os.DirEntry) *apperror.AppError {
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		copySingleFile(filepath.Join(src, e.Name()), filepath.Join(dest, e.Name()))
	}
	return nil
}

func copySingleFile(sFile string, dFile string) {
	content, rErr := os.ReadFile(sFile)
	if rErr != nil {
		return
	}
	_ = os.WriteFile(dFile, content, constants.FilePermission)
}
