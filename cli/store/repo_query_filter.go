package store

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

// ListReposByTag returns all repositories matching a specific tag stored in metadata or notes.
func (db *DB) ListReposByTag(tag string) ([]model.ScanRecord, *apperror.AppError) {
	return NewRepoRepository(db.conn).ListByTag(tag)
}

// ListReposUnpushed returns repositories flagged with unpushed commits in metadata or local tracking.
func (db *DB) ListReposUnpushed() ([]model.ScanRecord, *apperror.AppError) {
	return NewRepoRepository(db.conn).ListUnpushed()
}
