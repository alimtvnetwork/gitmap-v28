// Package store — makeallvisibility_undo.go: read-side helpers for
// `gitmap visibility-undo`. Pure SELECTs — no mutations live here.
//
// Spec: spec/01-app/116-bulk-visibility-mapub-mapri.md §undo-redo.
package store

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// SelectLatestUndoableMakeAllVisibilityRun returns the most recent run
// that has at least one Ok result with a captured PrevVisibility.
// Returns (zero, nil) when no such run exists — callers map that to
// ErrUndoNoRunFound at the CLI boundary.
func (db *DB) SelectLatestUndoableMakeAllVisibilityRun() (model.MakeAllVisibilityRunRecord, error) {
	var r model.MakeAllVisibilityRunRecord
	err := db.conn.QueryRow(constants.SQLSelectLatestUndoableRun).Scan(
		&r.ID, &r.CommandKind, &r.TargetVisibility, &r.Provider,
		&r.Owner, &r.TargetRaw, &r.OkCount,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return model.MakeAllVisibilityRunRecord{}, nil
	}
	if err != nil {
		return model.MakeAllVisibilityRunRecord{}, fmt.Errorf(constants.ErrUndoSelectRunFmt, err, err.Error())
	}

	return r, nil
}

// SelectUndoableResultsForRun returns the Ok rows of `runID` whose
// PrevVisibility differs from NewVisibility (i.e. actually mutated
// state worth reversing).
func (db *DB) SelectUndoableResultsForRun(runID int64) ([]model.MakeAllVisibilityResultRecord, error) {
	rows, err := db.conn.Query(constants.SQLSelectUndoableResultsForRun, runID)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrUndoSelectResultsFmt, err, err.Error())
	}
	defer rows.Close()

	return scanUndoableResults(rows)
}

// scanUndoableResults walks the result set. Extracted so the parent
// stays under the 15-line cap.
func scanUndoableResults(rows *sql.Rows) ([]model.MakeAllVisibilityResultRecord, error) {
	out := make([]model.MakeAllVisibilityResultRecord, 0, 8)
	for rows.Next() {
		var r model.MakeAllVisibilityResultRecord
		if err := rows.Scan(&r.ID, &r.RepoName, &r.MatchedPattern, &r.PrevVisibility, &r.NewVisibility); err != nil {
			return nil, fmt.Errorf(constants.ErrUndoSelectResultsFmt, err, err.Error())
		}
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(constants.ErrUndoSelectResultsFmt, err, err.Error())
	}

	return out, nil
}

// SelectMakeAllVisibilityRunByID resolves an explicit --run <id>.
// Returns (zero, nil) when the id is unknown — callers map to
// ErrUndoRunNotFoundFmt at the CLI boundary.
func (db *DB) SelectMakeAllVisibilityRunByID(id int64) (model.MakeAllVisibilityRunRecord, error) {
	return scanRunRow(db.conn.QueryRow(constants.SQLSelectRunByID, id))
}

// SelectLatestMakeAllVisibilityRunByKind picks the most recent run
// matching the given CommandKind that still has Ok rows.
func (db *DB) SelectLatestMakeAllVisibilityRunByKind(kind string) (model.MakeAllVisibilityRunRecord, error) {
	return scanRunRow(db.conn.QueryRow(constants.SQLSelectLatestRunByKind, kind))
}

// scanRunRow is the shared Scan path for the three run-select queries.
// Uses errors.Is(sql.ErrNoRows) so the absent-row case is a clean
// (zero, nil) return instead of a wrapped error.
func scanRunRow(row *sql.Row) (model.MakeAllVisibilityRunRecord, error) {
	var r model.MakeAllVisibilityRunRecord
	err := row.Scan(&r.ID, &r.CommandKind, &r.TargetVisibility, &r.Provider,
		&r.Owner, &r.TargetRaw, &r.OkCount)
	if errors.Is(err, sql.ErrNoRows) {
		return model.MakeAllVisibilityRunRecord{}, nil
	}
	if err != nil {
		return model.MakeAllVisibilityRunRecord{}, fmt.Errorf(constants.ErrUndoSelectRunFmt, err, err.Error())
	}

	return r, nil
}
