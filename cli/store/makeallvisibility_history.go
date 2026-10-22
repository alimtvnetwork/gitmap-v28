// Package store — makeallvisibility_history.go: SELECT helper for
// `gitmap visibility-history`. Pure read path — no mutations.
//
// Spec: 02-spec/01-app/116-bulk-visibility-mapub-mapri.md §history.
package store

import (
	"database/sql"

t"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

// SelectRecentMakeAllVisibilityRuns returns the newest `limit` rows.
func (db *DB) SelectRecentMakeAllVisibilityRuns(limit int) ([]model.MakeAllVisibilityRunRecord, error) {
	rows, err := QueryWrapper(db.conn, constants.SQLSelectRecentRuns, limit).Destruct()
	if err != nil {
		return nil, apperror.WrapSimple(err, "select recent runs failed (operation: SQLSelectRecentRuns)")
	}

	defer rows.Close()

	return scanRecentRuns(rows)
}

// scanRecentRuns walks the result set into typed records.
func scanRecentRuns(rows *sql.Rows) ([]model.MakeAllVisibilityRunRecord, error) {
	out := make([]model.MakeAllVisibilityRunRecord, 0, 16)
	for rows.Next() {
		var r model.MakeAllVisibilityRunRecord
		if err := rows.Scan(&r.ID, &r.CommandKind, &r.TargetVisibility, &r.Provider,
			&r.Owner, &r.MatchedCount, &r.OkCount, &r.SkippedCount, &r.FailedCount,
			&r.ExcludedCount, &r.ExitCode, &r.StartedAt, &r.FinishedAt); err != nil {
			return nil, apperror.WrapSimple(err, "select recent runs failed (operation: SQLSelectRecentRuns)")
		}

		out = append(out, r)
	}

	if err := rows.Err(); err != nil {
		return nil, apperror.WrapSimple(err, "select recent runs failed (operation: SQLSelectRecentRuns)")
	}

	return out, nil
}
