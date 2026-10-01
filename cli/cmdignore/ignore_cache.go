package cmdignore

import (
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// FilterReposNeedingCheck filters out repositories whose .gitignore status
// has already been audited within the specified TTL duration.
//
// Fails open (returns all repos) if TTL <= 0 or if split-DB fails to open.
func FilterReposNeedingCheck(repos []model.ScanRecord, ttl time.Duration) ([]model.ScanRecord, error) {
	return store.FilterReposNeedingIgnoreCheck(repos, ttl)
}

// RecordRepoCheckResult persists a repository inspection outcome into the Split-DB.
func RecordRepoCheckResult(repoPath, slug, status string, count int, dur time.Duration) error {
	return store.RecordIgnoreCheckResult(repoPath, slug, status, count, dur)
}

// RecordIgnoreCheckResult wraps RecordRepoCheckResult for caller flexibility.
func RecordIgnoreCheckResult(repoPath, slug, status string, count int, dur time.Duration) error {
	return store.RecordIgnoreCheckResult(repoPath, slug, status, count, dur)
}

// InvalidateRepoCache marks a repository's cache record as inactive.
func InvalidateRepoCache(repoPath string) error {
	db, err := store.OpenGitIgnoreSplitDB()
	if err != nil {
		return err
	}
	defer db.Close()
	return db.InvalidateCache(repoPath)
}

// ClearIgnoreCache invalidates all cached repository verification records.
func ClearIgnoreCache() error {
	db, err := store.OpenGitIgnoreSplitDB()
	if err != nil {
		return err
	}
	defer db.Close()
	return db.InvalidateAll()
}
