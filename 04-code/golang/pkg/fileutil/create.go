package fileutil

import (
	"os"

	"coding-guidelines/common/pkg/result"
)

// CreateDir is a direct alias to EnsureDir that explicitly implies creating a folder.
func CreateDir(path string, perm FilePermType) result.Wrap[bool] {
	return EnsureDir(path, perm)
}

// CreateFile creates a file at the given path with write-only and truncate flags, creating it if it doesn't exist.
func CreateFile(path string, perm FilePermType) result.Wrap[*os.File] {
	return OpenFile(path, FileOpenCreateTruncate, perm)
}
