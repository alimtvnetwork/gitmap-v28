// Package store manages the SQLite database for gitmap.
package store

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

// InsertTempRelease records a new temp-release branch in the database.
func (db *DB) InsertTempRelease(
	branch,
	versionPrefix string,
	seq int,
	commit,
	message string,
) error {
	_, err := ExecWrapper(db.conn, constants.SQLInsertTempRelease, branch, versionPrefix, seq, commit, message).Destruct()
	if err != nil {
		return apperror.WrapSimple(err, "failed to create temp-release")
	}

	return nil
}

// ListTempReleases returns all temp-release records ordered by sequence.
func (db *DB) ListTempReleases() ([]model.TempRelease, error) {
	rows, err := QueryWrapper(db.conn, constants.SQLSelectAllTempReleases).Destruct()
	if err != nil {
		return nil, apperror.WrapSimple(err, "failed to query temp-releases")
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
	_, err := ExecWrapper(db.conn, constants.SQLDeleteTempRelease, branch).Destruct()
	if err != nil {
		return apperror.WrapSimple(err, "failed to delete temp-release")
	}

	return nil
}

// DeleteAllTempReleases removes all temp-release records.
func (db *DB) DeleteAllTempReleases() error {
	_, err := ExecWrapper(db.conn, constants.SQLDeleteAllTempReleases).Destruct()
	if err != nil {
		return apperror.WrapSimple(err, "failed to delete temp-release")
	}

	return nil
}

// CountTempReleases returns the total number of temp-release records.
func (db *DB) CountTempReleases() (int, error) {
	var count int

	err := db.conn.QueryRow(constants.SQLCountTempReleases).Scan(&count)

	return count, err
}
