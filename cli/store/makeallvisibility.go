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
// Spec: 02-spec/01-app/116-bulk-visibility-mapub-mapri.md §plan steps 17-18.
package store

import (
	"context"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/dbengine"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// InsertMakeAllVisibilityRun writes the pre-prompt run row and returns
// its autogen ID. Pure single-statement insert — no tx needed.
func (db *DB) InsertMakeAllVisibilityRun(r model.MakeAllVisibilityRunRecord) Int64Wrap {
	res, err := ExecWrapper(db.conn, constants.SQLInsertMakeAllVisibilityRun,
		r.CommandKind, r.TargetVisibility, r.Provider, r.Owner, r.TargetRaw,
		r.PatternList, boolToInt(r.IsYesFlag), boolToInt(r.IsVerboseFlag),
		r.OwnerRepoTotal, r.MatchedCount, r.StartedAt).Destruct()
	if err != nil {
		return result.Fail[int64](apperror.WrapSimple(err, "insert MakeAllVisibilityRun failed (operation: SQLInsertMakeAllVisibilityRun)"))
	}

	id, idErr := res.LastInsertId()
	if idErr != nil {
		return result.Fail[int64](apperror.WrapSimple(idErr, "insert MakeAllVisibilityRun last insert id"))
	}

	return result.Ok(id)
}

// InsertMakeAllVisibilityPendingResults writes one 'Pending' row per
// matched repo inside a single transaction. Returns the assigned IDs
// in input order so callers can later UPDATE by primary key without
// re-querying.
func (db *DB) InsertMakeAllVisibilityPendingResults(runID int64, rows []model.MakeAllVisibilityResultRecord) Int64SliceWrap {
	wrap, appErr := dbengine.WrapDB(db.conn, dbengine.DBSQLite)
	if appErr != nil {
		return result.FailSlice[int64](apperror.WrapSimple(appErr, "insert MakeAllVisibilityResult failed (operation: SQLInsertMakeAllVisibilityResult)"))
	}

	var ids []int64
	ctx := context.Background()
	appErr = wrap.WithTransaction(ctx, func(tx *dbengine.TxWrapper) *apperror.AppError {
		idsRes := insertPendingResultsInTx(ctx, tx, runID, rows)
		ids = idsRes.Items()

		return idsRes.AppError()
	})
	if appErr != nil {
		return result.FailSlice[int64](apperror.WrapSimple(appErr, "insert MakeAllVisibilityResult failed (operation: SQLInsertMakeAllVisibilityResult)"))
	}

	return result.OkSlice(ids)
}

func insertPendingRow(
	ctx context.Context,
	tx *dbengine.TxWrapper,
	runID int64,
	r model.MakeAllVisibilityResultRecord,
) Int64Wrap {
	execRes := tx.Exec(ctx, constants.SQLInsertMakeAllVisibilityResult,
		runID, r.RepoName, r.MatchedPattern,
		constants.ResultStatusPending, r.StartedAt)
	if execRes.IsFailure() {
		return result.Fail[int64](execRes.AppError())
	}

	id, idErr := execRes.Value.LastInsertId()
	if idErr != nil {
		return result.Fail[int64](apperror.WrapSimple(idErr, "last insert id"))
	}

	return result.Ok(id)
}

// insertPendingResultsInTx is the per-row insert loop executed within a transaction.
func insertPendingResultsInTx(
	ctx context.Context,
	tx *dbengine.TxWrapper,
	runID int64,
	rows []model.MakeAllVisibilityResultRecord,
) Int64SliceWrap {
	ids := make([]int64, 0, len(rows))
	for _, r := range rows {
		idRes := insertPendingRow(ctx, tx, runID, r)
		if idRes.IsFailure() {
			return result.FailSlice[int64](idRes.AppError())
		}

		ids = append(ids, idRes.Value)
	}

	return result.OkSlice(ids)
}

// MarkMakeAllVisibilityResultsExcluded flips the given result rows to
// Status='Excluded' with FinishedAt = now. Single tx for atomicity.
func (db *DB) MarkMakeAllVisibilityResultsExcluded(ids []int64, finishedAt string) error {
	isMissing := len(ids) == 0
	if isMissing {
		return nil
	}

	wrap, appErr := dbengine.WrapDB(db.conn, dbengine.DBSQLite)
	if appErr != nil {
		return apperror.WrapSimple(appErr, "exclude MakeAllVisibilityResult rows failed (operation: SQLUpdateMakeAllVisibilityResultExcluded)")
	}

	ctx := context.Background()
	appErr = wrap.WithTransaction(ctx, func(tx *dbengine.TxWrapper) *apperror.AppError {
		return markIdsExcludedInTx(ctx, tx, ids, finishedAt)
	})
	if appErr != nil {
		return apperror.WrapSimple(appErr, "exclude MakeAllVisibilityResult rows failed (operation: SQLUpdateMakeAllVisibilityResultExcluded)")
	}

	return nil
}

func markIdsExcludedInTx(ctx context.Context, tx *dbengine.TxWrapper, ids []int64, finishedAt string) *apperror.AppError {
	for _, id := range ids {
		if execRes := tx.Exec(ctx, constants.SQLUpdateMakeAllVisibilityResultExcluded, finishedAt, id); execRes.IsFailure() {
			return execRes.AppError()
		}
	}

	return nil
}

// UpdateMakeAllVisibilityResult writes the terminal status for one
// per-repo result row after the apply+verify pipeline finishes.
func (db *DB) UpdateMakeAllVisibilityResult(r model.MakeAllVisibilityResultRecord) error {
	_, err := ExecWrapper(db.conn, constants.SQLUpdateMakeAllVisibilityResult,
		r.Status, r.PrevVisibility, r.NewVisibility, r.FailureMessage,
		r.FinishedAt, r.DurationMs, r.ID).Destruct()
	if err != nil {
		return apperror.WrapSimple(err, "update MakeAllVisibilityResult failed (operation: SQLUpdateMakeAllVisibilityResult)")
	}

	return nil
}

// FinalizeMakeAllVisibilityRun flushes the tallied counts + exit code
// + FinishedAt back to the run row.
func (db *DB) FinalizeMakeAllVisibilityRun(r model.MakeAllVisibilityRunRecord) error {
	_, err := ExecWrapper(db.conn, constants.SQLUpdateMakeAllVisibilityRunCounts,
		r.ExcludedCount, r.OkCount, r.SkippedCount, r.FailedCount,
		r.ExitCode, r.FinishedAt, r.ID).Destruct()
	if err != nil {
		return apperror.WrapSimple(err, "finalize MakeAllVisibilityRun failed (operation: SQLUpdateMakeAllVisibilityRunCounts)")
	}

	return nil
}
