package store

// ListReposByTag returns all repositories matching a specific tag stored in metadata or notes.
func (db *DB) ListReposByTag(tag string) ScanRecordSliceWrap {
	return NewRepoRepository(db.conn).ListByTag(tag)
}

// ListReposUnpushed returns repositories flagged with unpushed commits in metadata or local tracking.
func (db *DB) ListReposUnpushed() ScanRecordSliceWrap {
	return NewRepoRepository(db.conn).ListUnpushed()
}
