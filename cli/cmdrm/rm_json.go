package cmdrm

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/dbengine"
	"github.com/alimtvnetwork/gitmap-v28/cli/fsutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func removeRepoFromJSON(r model.ScanRecord) {
	paths := []string{
		filepath.Join(constants.DefaultOutputDir, constants.DefaultJSONFile),
		filepath.Join(constants.DefaultOutputFolder, constants.DefaultJSONFile),
	}
	for _, p := range paths {
		filterAndSaveJSON(p, r)
	}
}

func filterAndSaveJSON(path string, r model.ScanRecord) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return
	}

	records, err := model.LoadStatusRecords(path)
	if err != nil {
		return
	}

	filtered := filterOutRecord(records, r)
	data, err := json.MarshalIndent(filtered, "", constants.JSONIndent)
	if err != nil {
		return
	}

	_ = os.WriteFile(path, append(data, '\n'), 0644)
}

func filterOutRecord(records []model.ScanRecord, target model.ScanRecord) []model.ScanRecord {
	var out []model.ScanRecord
	for _, rec := range records {
		if !isMatchingScanRecord(rec, target) {
			out = append(out, rec)
		}
	}

	return out
}

func isMatchingScanRecord(a, b model.ScanRecord) bool {
	if fsutil.EqualPaths(a.AbsolutePath, b.AbsolutePath) {
		return true
	}

	if len(b.Slug) > 0 && strings.EqualFold(a.Slug, b.Slug) {
		return true
	}

	if len(b.RepoName) > 0 && strings.EqualFold(a.RepoName, b.RepoName) {
		return true
	}

	return false
}

func removeRepoDB(db *store.DB, r model.ScanRecord) *apperror.AppError {
	if db == nil {
		return nil
	}

	ctx := context.Background()
	wrapper, appErr := dbengine.WrapDB(db.Conn(), dbengine.DBSQLite)
	if appErr != nil {
		return appErr
	}

	return runRemoveRepoTx(ctx, wrapper, r)
}

func runRemoveRepoTx(ctx context.Context, wrapper *dbengine.DBWrapper, r model.ScanRecord) *apperror.AppError {
	return wrapper.WithImmediateTransaction(ctx, func(tx *dbengine.TxWrapper) *apperror.AppError {
		repoID := resolveRepoIDForRm(ctx, tx, r)

		return executeRemoveRepoTx(ctx, tx, repoID, r.AbsolutePath)
	})
}

func resolveRepoIDForRm(ctx context.Context, tx *dbengine.TxWrapper, r model.ScanRecord) int64 {
	if r.ID > 0 {
		return r.ID
	}

	var id int64
	row, appErr := tx.QueryRow(ctx, "SELECT RepoId FROM Repo WHERE AbsolutePath = ?", r.AbsolutePath)
	if appErr == nil && row.Scan(&id) == nil {
		return id
	}

	return scanFallbackRepoID(ctx, tx, r.AbsolutePath)
}

func scanFallbackRepoID(ctx context.Context, tx *dbengine.TxWrapper, absPath string) int64 {
	var id int64
	row, appErr := tx.QueryRow(ctx, "SELECT Id FROM Repo WHERE AbsolutePath = ?", absPath)
	if appErr != nil {
		return 0
	}

	if err := row.Scan(&id); err != nil {
		return 0
	}

	return id
}

func removeAliasIfPresent(ctx context.Context, tx *dbengine.TxWrapper, repoID int64) *apperror.AppError {
	if repoID <= 0 {
		return nil
	}

	if _, err := tx.Exec(ctx, "DELETE FROM Alias WHERE RepoId = ?", repoID); err != nil {
		return apperror.WrapSimple(err, "rm: delete alias")
	}

	return nil
}

func executeRemoveRepoTx(ctx context.Context, tx *dbengine.TxWrapper, repoID int64, absPath string) *apperror.AppError {
	if err := removeAliasIfPresent(ctx, tx, repoID); err != nil {
		return err
	}

	if repoID > 0 {
		return deleteRepoByIDAndPath(ctx, tx, repoID, absPath)
	}

	return deleteRepoByPath(ctx, tx, absPath)
}

func deleteRepoByIDAndPath(ctx context.Context, tx *dbengine.TxWrapper, repoID int64, absPath string) *apperror.AppError {
	normalizedPath := store.NormalizeStoragePath(absPath)
	query := "DELETE FROM Repo WHERE RepoId = ? OR AbsolutePath = ? OR AbsolutePath = ?"
	if _, err := tx.Exec(ctx, query, repoID, absPath, normalizedPath); err != nil {
		return apperror.WrapSimple(err, "rm: delete repo")
	}

	return recordRmHistory(ctx, tx, absPath)
}

func deleteRepoByPath(ctx context.Context, tx *dbengine.TxWrapper, absPath string) *apperror.AppError {
	normalizedPath := store.NormalizeStoragePath(absPath)
	query := "DELETE FROM Repo WHERE AbsolutePath = ? OR AbsolutePath = ?"
	if _, err := tx.Exec(ctx, query, absPath, normalizedPath); err != nil {
		return apperror.WrapSimple(err, "rm: delete repo")
	}

	return recordRmHistory(ctx, tx, absPath)
}

func recordRmHistory(ctx context.Context, tx *dbengine.TxWrapper, absPath string) *apperror.AppError {
	now := time.Now().UTC().Format(time.RFC3339)
	query := `INSERT INTO CommandHistory (Command, Alias, Args, Flags, StartedAt, FinishedAt, DurationMs, ExitCode, Summary, RepoCount)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := tx.Exec(ctx, query, "rm", "", absPath, "", now, now, 0, 0, "removed repository", 1)
	if err != nil {
		return apperror.WrapSimple(err, "record rm history")
	}

	return nil
}

func removeRepoDisk(absPath string, isDbOnly bool) *apperror.AppError {
	if isDbOnly {
		return nil
	}

	err := safeRemoveWithRetry(absPath)
	if err == nil {
		return nil
	}

	if isProcessLockError(err) {
		fmt.Fprintf(os.Stderr, "rm: warning: some files in %s are locked by another process; untracking repo\n", absPath)
		return nil
	}

	return apperror.WrapSimple(err, "remove dir")
}
