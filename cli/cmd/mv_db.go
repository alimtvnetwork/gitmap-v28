package cmd

import (
	"context"
	"database/sql"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/dbengine"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func updateRepoInDB(db *store.DB, repoID int64, newPath, newName string) error {
	ctx := context.Background()
	wrapper, appErr := dbengine.WrapDB(db.Conn(), dbengine.DBSQLite)
	if appErr != nil {
		return appErr
	}

	cleanPath := store.NormalizeStoragePath(newPath)
	if txErr := runUpdateMoveTx(ctx, wrapper, repoID, cleanPath, newName); txErr != nil {
		return txErr
	}

	return nil
}

func runUpdateMoveTx(ctx context.Context, wrapper *dbengine.DBWrapper, repoID int64, newPath, newName string) *apperror.AppError {
	return wrapper.WithImmediateTransaction(ctx, func(tx *dbengine.TxWrapper) *apperror.AppError {
		return executeMoveDBTx(ctx, tx, repoID, newPath, newName)
	})
}

func executeMoveDBTx(ctx context.Context, tx *dbengine.TxWrapper, repoID int64, newPath, newName string) *apperror.AppError {
	if err := updateRepoRow(ctx, tx, repoID, newPath, newName); err != nil {
		return err
	}

	return updateAliasRowIfPresent(ctx, tx, repoID, newPath)
}

func updateRepoRow(ctx context.Context, tx *dbengine.TxWrapper, repoID int64, newPath, newName string) *apperror.AppError {
	query := "UPDATE Repo SET AbsolutePath = ?, RepoName = ?, UpdatedAt = CURRENT_TIMESTAMP WHERE RepoId = ?"
	if execRes := tx.Exec(ctx, query, newPath, newName, repoID); execRes.IsSuccess() {
		return nil
	}

	fallbackQuery := "UPDATE Repo SET AbsolutePath = ?, RepoName = ? WHERE Id = ?"
	if fallbackRes := tx.Exec(ctx, fallbackQuery, newPath, newName, repoID); fallbackRes.IsFailure() {
		return apperror.WrapSimple(fallbackRes.AppError(), "mv: update repo row")
	}

	return nil
}

func updateAliasRowIfPresent(ctx context.Context, tx *dbengine.TxWrapper, repoID int64, newPath string) *apperror.AppError {
	hasColumn := hasAliasPathColumn(ctx, tx)
	if !hasColumn {
		return nil
	}

	query := "UPDATE Alias SET AbsolutePath = ? WHERE RepoId = ?"
	if execRes := tx.Exec(ctx, query, newPath, repoID); execRes.IsFailure() {
		return apperror.WrapSimple(execRes.AppError(), "mv: update alias row")
	}

	return nil
}

func hasAliasPathColumn(ctx context.Context, tx *dbengine.TxWrapper) bool {
	rowsRes := tx.Query(ctx, "PRAGMA table_info(Alias)")
	if rowsRes.IsFailure() {
		return false
	}

	rows := rowsRes.Value
	defer rows.Close()

	return scanForColumnName(rows, "AbsolutePath")
}

func scanForColumnName(rows *sql.Rows, targetCol string) bool {
	for rows.Next() {
		var cid, notnull, pk int
		var name, ctype string
		var dflt any
		scanErr := rows.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk)
		isMatch := scanErr == nil && strings.EqualFold(name, targetCol)
		if isMatch {
			return true
		}
	}

	return false
}
