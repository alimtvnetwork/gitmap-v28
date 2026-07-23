package clonenext

// Repo-detection helpers shared between the cn batch path and the
// dispatcher's implicit-batch trigger. Split out of batch.go so each
// file stays focused (batch.go = CSV/walk pipeline, repodetect.go =
// "is this directory a repo / a scan root?" predicates) and to keep
// per-file size budgets healthy.

import (
	"os"
	"path/filepath"
)

// IsGitRepo reports whether path contains a .git entry (file or
// directory — `.git` files exist for git worktrees). Exported so the
// cmd-package dispatcher can decide between single-repo and batch
// mode without importing internal helpers.
func IsGitRepo(path string) bool {
	_, err := os.Stat(filepath.Join(path, ".git"))

	return err == nil
}

// isGitRepo is the unexported alias kept for back-compat with existing
// call sites inside this package. New callers should use IsGitRepo.
func isGitRepo(path string) bool {
	return IsGitRepo(path)
}

// HasGitSubdir reports whether `root` contains at least one immediate
// child directory that is itself a git repo. Designed for the cn
// dispatcher's implicit-batch trigger: it short-circuits on the first
// hit, so the cost is bounded to one ReadDir + at most one Stat per
// candidate up to the first match. Returns false on any I/O error so
// the dispatcher fails closed (single-repo path → clean "no remote"
// error) rather than guessing.
func HasGitSubdir(root string) bool {
	entries, err := os.ReadDir(root)
	if err != nil {
		return false
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		if IsGitRepo(filepath.Join(root, entry.Name())) {
			return true
		}
	}

	return false
}
