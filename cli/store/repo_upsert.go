package store

import (
	"context"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

// NormalizeRepoRecord ensures a record's AbsolutePath is canonicalized before storage.
func NormalizeRepoRecord(r *model.ScanRecord) {
	if r != nil {
		r.AbsolutePath = NormalizeStoragePath(r.AbsolutePath)
	}
}

// IsPathDuplicateCaseInsensitive checks if two paths are equivalent under case-insensitive comparison.
func IsPathDuplicateCaseInsensitive(pathA, pathB string) bool {
	return strings.EqualFold(NormalizeStoragePath(pathA), NormalizeStoragePath(pathB))
}

// UpsertOneRepoRecord inserts or updates a single repository record ensuring path normalization.
func (db *DB) UpsertOneRepoRecord(ctx context.Context, r model.ScanRecord) error {
	NormalizeRepoRecord(&r)

	return db.UpsertRepos([]model.ScanRecord{r})
}
