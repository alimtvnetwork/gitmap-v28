// Package store — clone_recent.go: stamps and reads the per-repo
// LastClonedAt column added in schema v26. Powers the
// `gitmap release` auto-cd fallback (most-recent-clone rule).
package store

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// RecentClone summarizes the most recently cloned repo.
type RecentClone struct {
	AbsolutePath string
	RepoName     string
	ClonedAt     string
}

// MarkCloned stamps the LastClonedAt column on the row matching
// absPath. No-op (returns nil) when no row matches — the upstream
// upsert is the only thing that creates Repo rows.
func (db *DB) MarkCloned(absPath string) error {
	if _, err := ExecWrapper(db.conn, constants.SQLUpdateRepoLastClonedAt, absPath).Destruct(); err != nil {
		return fmt.Errorf(constants.ErrDBUpsert, err)
	}

	return nil
}

// MostRecentClone returns the most recently cloned repo, or
// (zero, false, nil) when no clone has ever been recorded.
func (db *DB) MostRecentClone() (RecentClone, bool, error) {
	var r RecentClone
	row := QueryRowWrapper(db.conn, constants.SQLSelectMostRecentClone)
	err := row.Scan(&r.AbsolutePath, &r.RepoName, &r.ClonedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return RecentClone{}, false, nil
	}
	if err != nil {
		return RecentClone{}, false, fmt.Errorf(constants.ErrDBQuery, err)
	}

	return r, true, nil
}
