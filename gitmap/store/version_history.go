package store

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// InsertVersionHistory records a version transition for a repo.
func (db *DB) InsertVersionHistory(r model.RepoVersionHistoryRecord) (int64, error) {
	result, err := db.conn.Exec(constants.SQLInsertVersionHistory,
		r.RepoID, r.FromVersionTag, r.FromVersionNum,
		r.ToVersionTag, r.ToVersionNum, r.FlattenedPath)
	if err != nil {
		return 0, fmt.Errorf(constants.ErrDBVersionHistory, err)
	}

	return result.LastInsertId()
}

// UpdateRepoVersion updates the current version columns on a Repos row.
func (db *DB) UpdateRepoVersion(repoID int64, versionTag string, versionNum int) error {
	_, err := db.conn.Exec(constants.SQLUpdateRepoVersion, versionTag, versionNum, repoID)
	if err != nil {
		return fmt.Errorf(constants.ErrDBVersionHistory, err)
	}

	return nil
}

// ListVersionHistory returns all version transitions for a repo.
func (db *DB) ListVersionHistory(repoID int64) ([]model.RepoVersionHistoryRecord, error) {
	rows, err := db.conn.Query(constants.SQLSelectVersionHistory, repoID)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrDBVersionHistory, err)
	}
	defer rows.Close()

	var results []model.RepoVersionHistoryRecord

	for rows.Next() {
		var r model.RepoVersionHistoryRecord
		scanErr := rows.Scan(&r.ID, &r.RepoID, &r.FromVersionTag, &r.FromVersionNum,
			&r.ToVersionTag, &r.ToVersionNum, &r.FlattenedPath, &r.CreatedAt)
		if scanErr != nil {
			return nil, fmt.Errorf(constants.ErrDBVersionHistory, scanErr)
		}
		results = append(results, r)
	}

	return results, nil
}

// GetRepoIDByPath returns the Repos.Id for a given absolute path, or 0 if not found.
func (db *DB) GetRepoIDByPath(absPath string) (int64, error) {
	var id int64
	err := db.conn.QueryRow(constants.SQLSelectRepoIDByPath, absPath).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
