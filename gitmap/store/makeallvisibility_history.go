// Package store — makeallvisibility_history.go: SELECT helper for
// `gitmap visibility-history`. Pure read path — no mutations.
//
// Spec: spec/01-app/116-bulk-visibility-mapub-mapri.md §history.
package store

import (
	"database/sql"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// SelectRecentMakeAllVisibilityRuns returns the newest `limit` rows.
func (db *DB) SelectRecentMakeAllVisibilityRuns(limit int) ([]model.MakeAllVisibilityRunRecord, error) {
	rows, err := db.conn.Query(constants.SQLSelectRecentRuns, limit)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrHistorySelectFmt, err, err.Error())
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
			return nil, fmt.Errorf(constants.ErrHistorySelectFmt, err, err.Error())
		}
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(constants.ErrHistorySelectFmt, err, err.Error())
	}

	return out, nil
}
