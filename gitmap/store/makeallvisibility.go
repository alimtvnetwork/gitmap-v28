// Package store — makeallvisibility.go: persistence layer for the
// bulk wildcard visibility audit trail. All public helpers wrap their
// multi-statement writes in a sql.Tx so a crash mid-run cannot leave
// half-written rows (zero-swallow + transactional integrity).
//
// Call sequence from runMakeAllVisibility:
//
//  1. InsertMakeAllVisibilityRun(run)                   → runID
//  2. InsertMakeAllVisibilityPendingResults(runID, ms)  → []resultID
//  3. (optional) MarkMakeAllVisibilityResultsExcluded(ids)
//  4. UpdateMakeAllVisibilityResult(id, status, ...)    × N
//  5. FinalizeMakeAllVisibilityRun(runID, counts, exitCode)
//
// Spec: spec/01-app/116-bulk-visibility-mapub-mapri.md §plan steps 17-18.
package store

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// InsertMakeAllVisibilityRun writes the pre-prompt run row and returns
// its autogen ID. Pure single-statement insert — no tx needed.
func (db *DB) InsertMakeAllVisibilityRun(r model.MakeAllVisibilityRunRecord) (int64, error) {
	res, err := db.conn.Exec(constants.SQLInsertMakeAllVisibilityRun,
		r.CommandKind, r.TargetVisibility, r.Provider, r.Owner, r.TargetRaw,
		r.PatternList, boolToInt(r.YesFlag), boolToInt(r.VerboseFlag),
		r.OwnerRepoTotal, r.MatchedCount, r.StartedAt)
	if err != nil {
		return 0, fmt.Errorf(constants.ErrMakeAllRunInsertFmt, err, err.Error())
	}

	return res.LastInsertId()
}

// InsertMakeAllVisibilityPendingResults writes one 'Pending' row per
// matched repo inside a single transaction. Returns the assigned IDs
// in input order so callers can later UPDATE by primary key without
// re-querying.
func (db *DB) InsertMakeAllVisibilityPendingResults(runID int64, rows []model.MakeAllVisibilityResultRecord) ([]int64, error) {
	tx, err := db.conn.Begin()
	if err != nil {
		return nil, fmt.Errorf(constants.ErrMakeAllResultInsertFmt, err, err.Error())
	}

	ids, err := insertPendingResultsInTx(tx, runID, rows)
	if err != nil {
		_ = tx.Rollback()

		return nil, err
	}

	if commitErr := tx.Commit(); commitErr != nil {
		return nil, fmt.Errorf(constants.ErrMakeAllResultInsertFmt, commitErr, commitErr.Error())
	}

	return ids, nil
}

// insertPendingResultsInTx is the per-row insert loop. Extracted so the
// outer function stays under the 15-line cap and the loop is testable
// against an injected tx in unit tests.
func insertPendingResultsInTx(tx txExecer, runID int64, rows []model.MakeAllVisibilityResultRecord) ([]int64, error) {
	ids := make([]int64, 0, len(rows))
	for _, r := range rows {
		res, err := tx.Exec(constants.SQLInsertMakeAllVisibilityResult,
			runID, r.RepoName, r.MatchedPattern,
			constants.ResultStatusPending, r.StartedAt)
		if err != nil {
			return nil, fmt.Errorf(constants.ErrMakeAllResultInsertFmt, err, err.Error())
		}
		id, idErr := res.LastInsertId()
		if idErr != nil {
			return nil, fmt.Errorf(constants.ErrMakeAllResultInsertFmt, idErr, idErr.Error())
		}
		ids = append(ids, id)
	}

	return ids, nil
}

// MarkMakeAllVisibilityResultsExcluded flips the given result rows to
// Status='Excluded' with FinishedAt = now. Single tx for atomicity.
func (db *DB) MarkMakeAllVisibilityResultsExcluded(ids []int64, finishedAt string) error {
	if len(ids) == 0 {
		return nil
	}
	tx, err := db.conn.Begin()
	if err != nil {
		return fmt.Errorf(constants.ErrMakeAllResultExcludeFmt, err, err.Error())
	}
	for _, id := range ids {
		if _, execErr := tx.Exec(constants.SQLUpdateMakeAllVisibilityResultExcluded, finishedAt, id); execErr != nil {
			_ = tx.Rollback()

			return fmt.Errorf(constants.ErrMakeAllResultExcludeFmt, execErr, execErr.Error())
		}
	}

	return commitOrWrap(tx, constants.ErrMakeAllResultExcludeFmt)
}

// UpdateMakeAllVisibilityResult writes the terminal status for one
// per-repo result row after the apply+verify pipeline finishes.
func (db *DB) UpdateMakeAllVisibilityResult(r model.MakeAllVisibilityResultRecord) error {
	_, err := db.conn.Exec(constants.SQLUpdateMakeAllVisibilityResult,
		r.Status, r.PrevVisibility, r.NewVisibility, r.FailureMessage,
		r.FinishedAt, r.DurationMs, r.ID)
	if err != nil {
		return fmt.Errorf(constants.ErrMakeAllResultUpdateFmt, err, err.Error())
	}

	return nil
}

// FinalizeMakeAllVisibilityRun flushes the tallied counts + exit code
// + FinishedAt back to the run row.
func (db *DB) FinalizeMakeAllVisibilityRun(r model.MakeAllVisibilityRunRecord) error {
	_, err := db.conn.Exec(constants.SQLUpdateMakeAllVisibilityRunCounts,
		r.ExcludedCount, r.OkCount, r.SkippedCount, r.FailedCount,
		r.ExitCode, r.FinishedAt, r.ID)
	if err != nil {
		return fmt.Errorf(constants.ErrMakeAllRunFinalizeFmt, err, err.Error())
	}

	return nil
}
