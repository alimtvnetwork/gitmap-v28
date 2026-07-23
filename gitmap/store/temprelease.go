// Package store manages the SQLite database for gitmap.
package store

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// InsertTempRelease records a new temp-release branch in the database.
func (db *DB) InsertTempRelease(branch, versionPrefix string, seq int, commit, message string) error {
	_, err := db.conn.Exec(constants.SQLInsertTempRelease, branch, versionPrefix, seq, commit, message)
	if err != nil {
		return fmt.Errorf(constants.ErrTRCreate, err)
	}

	return nil
}

// ListTempReleases returns all temp-release records ordered by sequence.
func (db *DB) ListTempReleases() ([]model.TempRelease, error) {
	rows, err := db.conn.Query(constants.SQLSelectAllTempReleases)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrTRQuery, err)
	}
	defer rows.Close()

	var releases []model.TempRelease

	for rows.Next() {
		var r model.TempRelease

		err := rows.Scan(&r.ID, &r.Branch, &r.VersionPrefix, &r.SequenceNumber, &r.CommitSha, &r.CommitMessage, &r.CreatedAt)
		if err != nil {
			continue
		}

		releases = append(releases, r)
	}

	return releases, nil
}

// MaxTempReleaseSeq returns the highest sequence number for a version prefix.
func (db *DB) MaxTempReleaseSeq(versionPrefix string) (int, error) {
	var max int

	err := db.conn.QueryRow(constants.SQLSelectMaxSeqByPrefix, versionPrefix).Scan(&max)
	if err != nil {
		return 0, err
	}

	return max, nil
}

// DeleteTempRelease removes a single temp-release record by branch name.
func (db *DB) DeleteTempRelease(branch string) error {
	_, err := db.conn.Exec(constants.SQLDeleteTempRelease, branch)
	if err != nil {
		return fmt.Errorf(constants.ErrTRDelete, err)
	}

	return nil
}

// DeleteAllTempReleases removes all temp-release records.
func (db *DB) DeleteAllTempReleases() error {
	_, err := db.conn.Exec(constants.SQLDeleteAllTempReleases)
	if err != nil {
		return fmt.Errorf(constants.ErrTRDelete, err)
	}

	return nil
}

// CountTempReleases returns the total number of temp-release records.
func (db *DB) CountTempReleases() (int, error) {
	var count int

	err := db.conn.QueryRow(constants.SQLCountTempReleases).Scan(&count)

	return count, err
}
