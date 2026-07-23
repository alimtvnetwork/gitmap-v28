package store

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// UpsertRepos inserts or updates all records by absolute_path.
func (db *DB) UpsertRepos(records []model.ScanRecord) error {
	for _, r := range records {
		if err := db.upsertOne(r); err != nil {
			return fmt.Errorf(constants.ErrDBUpsert, err)
		}
	}

	return nil
}

// upsertOne inserts or updates a single repo by absolute_path.
func (db *DB) upsertOne(r model.ScanRecord) error {
	_, err := db.conn.Exec(constants.SQLUpsertRepoByPath,
		r.Slug, r.RepoName, r.HTTPSUrl, r.SSHUrl,
		r.Branch, r.RelativePath, r.AbsolutePath,
		r.CloneInstruction, r.Notes, r.Transport,
	)

	return err
}

// DeleteByPath removes the repo row whose AbsolutePath matches.
// Returns the number of rows deleted (0 when no match).
func (db *DB) DeleteByPath(absPath string) (int64, error) {
	res, err := db.conn.Exec(constants.SQLDeleteRepoByPath, absPath)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

// DeleteBySlug removes every repo row whose Slug matches.
// Returns the number of rows deleted.
func (db *DB) DeleteBySlug(slug string) (int64, error) {
	res, err := db.conn.Exec(constants.SQLDeleteRepoBySlug, slug)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

// ListRepos returns all tracked repositories ordered by slug.
func (db *DB) ListRepos() ([]model.ScanRecord, error) {
	rows, err := db.conn.Query(constants.SQLSelectAllRepos)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrDBQuery, err)
	}
	defer rows.Close()

	return scanRows(rows)
}

// FindBySlug returns all repos matching the given slug.
func (db *DB) FindBySlug(slug string) ([]model.ScanRecord, error) {
	rows, err := db.conn.Query(constants.SQLSelectRepoBySlug, slug)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrDBQuery, err)
	}
	defer rows.Close()

	return scanRows(rows)
}

// FindByPath returns the repo at the given absolute path.
func (db *DB) FindByPath(absPath string) ([]model.ScanRecord, error) {
	rows, err := db.conn.Query(constants.SQLSelectRepoByPath, absPath)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrDBQuery, err)
	}
	defer rows.Close()

	return scanRows(rows)
}

// scanRows reads ScanRecord values from query result rows.
func scanRows(rows interface {
	Next() bool
	Scan(dest ...any) error
}) ([]model.ScanRecord, error) {
	var results []model.ScanRecord

	for rows.Next() {
		r, err := scanOneRow(rows)
		if err != nil {
			return nil, err
		}
		results = append(results, r)
	}

	return results, nil
}

// scanOneRow reads a single ScanRecord from the current row.
func scanOneRow(row interface{ Scan(dest ...any) error }) (model.ScanRecord, error) {
	var r model.ScanRecord
	err := row.Scan(
		&r.ID, &r.Slug, &r.RepoName, &r.HTTPSUrl, &r.SSHUrl,
		&r.Branch, &r.RelativePath, &r.AbsolutePath,
		&r.CloneInstruction, &r.Notes, &r.Transport,
	)

	return r, err
}
