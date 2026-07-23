package store

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// TagReposByScanFolder bulk-updates Repo.ScanFolderId for every row whose
// AbsolutePath is in the supplied list. No-op when paths is empty.
func (db *DB) TagReposByScanFolder(scanFolderID int64, paths []string) error {
	if len(paths) == 0 {
		return nil
	}

	placeholders := strings.Repeat("?,", len(paths))
	placeholders = placeholders[:len(placeholders)-1]
	query := fmt.Sprintf(constants.SQLTagReposByScanFolderTpl, placeholders)

	args := make([]any, 0, len(paths)+1)
	args = append(args, scanFolderID)
	for _, p := range paths {
		args = append(args, p)
	}

	if _, err := db.conn.Exec(query, args...); err != nil {
		return fmt.Errorf(constants.ErrProbeTagFail, scanFolderID, err)
	}

	return nil
}

// RecordVersionProbe persists a probe result.
func (db *DB) RecordVersionProbe(p model.VersionProbe) error {
	available := 0
	if p.IsAvailable {
		available = 1
	}

	_, err := db.conn.Exec(constants.SQLInsertVersionProbe,
		p.RepoID, p.NextVersionTag, p.NextVersionNum,
		p.Method, available, p.Error,
	)
	if err != nil {
		return fmt.Errorf(constants.ErrProbeRecord, p.RepoID, err)
	}

	return nil
}

// LatestVersionProbe returns the newest probe row for a repo, or
// sql.ErrNoRows when no probe has run yet.
func (db *DB) LatestVersionProbe(repoID int64) (model.VersionProbe, error) {
	row := db.conn.QueryRow(constants.SQLSelectLatestVersionProbe, repoID)

	var (
		p         model.VersionProbe
		available int
	)
	err := row.Scan(&p.ID, &p.RepoID, &p.ProbedAt,
		&p.NextVersionTag, &p.NextVersionNum,
		&p.Method, &available, &p.Error)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.VersionProbe{}, sql.ErrNoRows
		}

		return model.VersionProbe{}, fmt.Errorf(constants.ErrProbeRecord, repoID, err)
	}
	p.IsAvailable = available == 1

	return p, nil
}
