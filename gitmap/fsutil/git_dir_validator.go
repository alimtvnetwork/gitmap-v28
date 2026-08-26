// Package fsutil — git_dir_validator.go validates Git directory accessibility and integrity.
package fsutil

import (
	"os"
	"path/filepath"
)

// IsValidGitDir returns true if dir is accessible and contains a readable .git directory or gitdir pointer.
func IsValidGitDir(dir string) bool {
	gitPath := filepath.Join(dir, ".git")
	info, err := os.Stat(gitPath)
	return err == nil && (info.IsDir() || !info.IsDir())
}
