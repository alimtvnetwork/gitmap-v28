package store

// ReleaseRepoIntegrity returns:
//   - orphaned: count of Release rows whose RepoId has no matching Repo row
//   - reposNoRel: count of Repo rows that have zero Release rows
//
// Both queries are read-only and safe to run from doctor.
func (db *DB) ReleaseRepoIntegrity() (orphaned, reposNoRel int, err error) {
	if !db.tableExists("Release") || !db.tableExists("Repo") {
		return 0, 0, nil
	}

	if !db.columnExists("Release", "RepoId") {
		return 0, 0, nil // pre-v17 schema, skip
	}

	err = db.conn.QueryRow(
		"SELECT COUNT(*) FROM Release WHERE RepoId NOT IN (SELECT RepoId FROM Repo)",
	).Scan(&orphaned)
	if err != nil {
		return 0, 0, err
	}

	err = db.conn.QueryRow(
		"SELECT COUNT(*) FROM Repo WHERE RepoId NOT IN (SELECT DISTINCT RepoId FROM Release)",
	).Scan(&reposNoRel)
	if err != nil {
		return orphaned, 0, err
	}

	return orphaned, reposNoRel, nil
}
