package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/dbengine"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

// EnsureScanFolder upserts a scan folder by AbsolutePath and returns the
// resulting row. Idempotent: re-invoking with the same path bumps
// LastScannedAt and only overwrites Label/Notes when the new values are
// non-empty (so subsequent scans don't wipe a manually set label).
func (db *DB) EnsureScanFolder(absPath, label, notes string) (model.ScanFolder, error) {
	if _, err := ExecWrapper(db.conn, constants.SQLUpsertScanFolder, absPath, label, notes).Destruct(); err != nil {
		return model.ScanFolder{}, apperror.WrapSimple(err, fmt.Sprintf("failed to ensure scan folder %q", absPath))
	}

	return db.findScanFolderByPath(absPath)
}

// ListScanFolders returns every registered scan folder, newest-scanned first.
func (db *DB) ListScanFolders() ([]model.ScanFolder, error) {
	rows, err := QueryWrapper(db.conn, constants.SQLSelectAllScanFolders).Destruct()
	if err != nil {
		return nil, apperror.WrapSimple(err, "failed to list scan folders")
	}

	defer rows.Close()

	return scanScanFolderRows(rows)
}

// CountReposInScanFolder returns how many Repo rows currently FK to the
// given scan folder. Used by `gitmap sf rm` to report what gets detached.
func (db *DB) CountReposInScanFolder(scanFolderID int64) (int, error) {
	var count int
	err := db.conn.QueryRow(constants.SQLCountReposInScanFolder, scanFolderID).Scan(&count)
	if err != nil {
		return 0, apperror.WrapSimple(err, "failed to list scan folders")
	}

	return count, nil
}

// RemoveScanFolderByPath detaches every repo pointing at the folder
// (sets Repo.ScanFolderId = NULL) and then deletes the ScanFolder row.
// Returns (removedFolder, detachedRepoCount).
func (db *DB) RemoveScanFolderByPath(absPath string) (model.ScanFolder, int, error) {
	folder, err := db.findScanFolderByPath(absPath)
	if err != nil {
		return model.ScanFolder{}, 0, err
	}

	return db.removeScanFolderRow(folder)
}

// RemoveScanFolderByID is the same as RemoveScanFolderByPath but keyed by ID.
func (db *DB) RemoveScanFolderByID(id int64) (model.ScanFolder, int, error) {
	folder, err := db.findScanFolderByID(id)
	if err != nil {
		return model.ScanFolder{}, 0, err
	}

	return db.removeScanFolderRow(folder)
}

// removeScanFolderRow performs detach-then-delete in a transactional pair.
func (db *DB) removeScanFolderRow(folder model.ScanFolder) (model.ScanFolder, int, error) {
	count, err := db.CountReposInScanFolder(folder.ID)
	if err != nil {
		return folder, 0, err
	}

	wrap, appErr := dbengine.WrapDB(db.conn, dbengine.DBSQLite)
	if appErr != nil {
		return folder, 0, appErr
	}

	if appErr := runDetachAndDeleteTx(wrap, folder.ID); appErr != nil {
		return folder, 0, appErr
	}

	return folder, count, nil
}

func runDetachAndDeleteTx(wrap *dbengine.DBWrapper, folderID int64) *apperror.AppError {
	ctx := context.Background()

	return wrap.WithTransaction(ctx, func(tx *dbengine.TxWrapper) *apperror.AppError {
		return executeDetachAndDeleteSF(tx.Tx(), folderID)
	})
}

func executeDetachAndDeleteSF(tx *sql.Tx, folderID int64) *apperror.AppError {
	if _, err := ExecWrapper(tx, constants.SQLDetachReposFromScanFolder, folderID).Destruct(); err != nil {
		return apperror.WrapSimple(err, "detach repos from scan folder")
	}

	if _, err := ExecWrapper(tx, constants.SQLDeleteScanFolderByID, folderID).Destruct(); err != nil {
		return apperror.WrapSimple(err, "delete scan folder by id")
	}

	return nil
}

// findScanFolderByPath returns the row matching AbsolutePath.
func (db *DB) findScanFolderByPath(absPath string) (model.ScanFolder, error) {
	row := QueryRowWrapper(db.conn, constants.SQLSelectScanFolderByPath, absPath)
	folder, err := scanOneScanFolder(row)
	if errors.Is(err, sql.ErrNoRows) {
		return model.ScanFolder{}, apperror.NewNotFoundError(fmt.Sprintf("no scan folder registered for path: %s", absPath))
	}

	if err != nil {
		return model.ScanFolder{}, apperror.WrapSimple(err, "failed to list scan folders")
	}

	return folder, nil
}

// findScanFolderByID returns the row matching ScanFolderId.
func (db *DB) findScanFolderByID(id int64) (model.ScanFolder, error) {
	row := QueryRowWrapper(db.conn, constants.SQLSelectScanFolderByID, id)
	folder, err := scanOneScanFolder(row)
	if errors.Is(err, sql.ErrNoRows) {
		return model.ScanFolder{}, apperror.NewNotFoundError(fmt.Sprintf("no scan folder with id: %d", id))
	}

	if err != nil {
		return model.ScanFolder{}, apperror.WrapSimple(err, "failed to list scan folders")
	}

	return folder, nil
}

// scanOneScanFolder reads a single ScanFolder row.
func scanOneScanFolder(row interface{ Scan(dest ...any) error }) (model.ScanFolder, error) {
	var f model.ScanFolder
	err := row.Scan(&f.ID, &f.AbsolutePath, &f.Label, &f.Notes, &f.LastScannedAt, &f.CreatedAt)

	return f, err
}

// scanScanFolderRows reads ScanFolder values from query result rows.
func scanScanFolderRows(rows interface {
	Next() bool
	Scan(dest ...any) error
}) ([]model.ScanFolder, error) {
	var results []model.ScanFolder
	for rows.Next() {
		f, err := scanOneScanFolder(rows)
		if err != nil {
			return nil, apperror.WrapSimple(err, "failed to list scan folders")
		}

		results = append(results, f)
	}

	return results, nil
}
