package cmdhygiene

import (
	"os"
	"path/filepath"
)

// IsGitRepoFn is wired by cmd/di_hooks.go to the canonical implementation.
var IsGitRepoFn func(path string) bool

// IsGitRepo reports whether path contains a .git entry (file or
// directory). New callers should use IsGitRepo.
func IsGitRepo(path string) bool {
	if IsGitRepoFn != nil {
		return IsGitRepoFn(path)
	}

	fi, err := os.Stat(filepath.Join(path, ".git"))
	return err == nil && fi != nil
}
