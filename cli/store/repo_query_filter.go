package store

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// ListReposByTag returns all repositories matching a specific tag stored in metadata or notes.
func (db *DB) ListReposByTag(tag string) result.Result[[]model.ScanRecord] {
	return NewRepoRepository(db.conn).ListByTag(tag)
}

// ListReposUnpushed returns repositories flagged with unpushed commits in metadata or local tracking.
func (db *DB) ListReposUnpushed() result.Result[[]model.ScanRecord] {
	return NewRepoRepository(db.conn).ListUnpushed()
}
