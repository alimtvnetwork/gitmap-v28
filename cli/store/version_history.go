package store

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

// InsertVersionHistory records a version transition for a repo.
func (db *DB) InsertVersionHistory(r model.RepoVersionHistoryRecord) (int64, error) {
	result, err := ExecWrapper(db.conn, constants.SQLInsertVersionHistory,
		r.RepoID, r.FromVersionTag, r.FromVersionNum,
		r.ToVersionTag, r.ToVersionNum, r.FlattenedPath).Destruct()
	if err != nil {
		return 0, apperror.WrapSimple(err, "failed to query version history")
	}

	return result.LastInsertId()
}

// UpdateRepoVersion updates the current version columns on a Repos row.
func (db *DB) UpdateRepoVersion(repoID int64, versionTag string, versionNum int) error {
	_, err := ExecWrapper(db.conn, constants.SQLUpdateRepoVersion, versionTag, versionNum, repoID).Destruct()
	if err != nil {
		return apperror.WrapSimple(err, "failed to query version history")
	}

	return nil
}

// ListVersionHistory returns all version transitions for a repo.
func (db *DB) ListVersionHistory(repoID int64) ([]model.RepoVersionHistoryRecord, error) {
	rows, err := QueryWrapper(db.conn, constants.SQLSelectVersionHistory, repoID).Destruct()
	if err != nil {
		return nil, apperror.WrapSimple(err, "failed to query version history")
	}

	defer rows.Close()

	var results []model.RepoVersionHistoryRecord

	for rows.Next() {
		var r model.RepoVersionHistoryRecord
		scanErr := rows.Scan(&r.ID, &r.RepoID, &r.FromVersionTag, &r.FromVersionNum,
			&r.ToVersionTag, &r.ToVersionNum, &r.FlattenedPath, &r.CreatedAt)
		if scanErr != nil {
			return nil, apperror.WrapSimple(scanErr, "failed to query version history")
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
