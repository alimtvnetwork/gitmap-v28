package store

import (
	"context"
	"database/sql"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/dbengine"
)

// SQLDeleteInstallerScript deletes an installer script record by its slug.
const SQLDeleteInstallerScript = `DELETE FROM installer_scripts WHERE slug = ?;`

// SQLDeleteInstallerVersions deletes all installer version records for a given slug.
const SQLDeleteInstallerVersions = `DELETE FROM installer_versions WHERE slug = ?;`

// SQLDeleteInstallerExactVersion deletes a specific version record for a given slug.
const SQLDeleteInstallerExactVersion = `DELETE FROM installer_versions WHERE slug = ? AND version = ?;`

// DeleteInstaller removes an installer script and associated version records by slug.
func (db *DB) DeleteInstaller(slug string) error {
	if appErr := validateInstallerInput(db, slug); appErr != nil {
		return appErr
	}

	wrap, appErr := dbengine.WrapDb(db.conn, dbengine.DbSQLite)
	if appErr != nil {
		return appErr
	}

	if appErr := runDeleteInstallerTx(wrap, slug); appErr != nil {
		return appErr
	}

	return nil
}

func validateInstallerInput(db *DB, slug string) *apperror.AppError {
	if db == nil || db.conn == nil {
		return apperror.New("DeleteInstaller", "E_INSTALLER_NIL_DB", map[string]any{"slug": slug})
	}
	if slug == "" {
		return apperror.New("DeleteInstaller", "E_INSTALLER_INVALID_INPUT", map[string]any{"error": "slug cannot be empty"})
	}

	return nil
}

func runDeleteInstallerTx(wrap *dbengine.DbWrapper, slug string) *apperror.AppError {
	ctx := context.Background()

	appErr := wrap.WithTransaction(ctx, func(tx *dbengine.TxWrapper) *apperror.AppError {
		return executeDeleteInstaller(tx.Tx(), slug)
	})
	if appErr == nil {
		return nil
	}

	if appErr.Code != "E_INSTALLER_NOT_FOUND" {
		appErr.Code = "E_INSTALLER_DELETE_FAILED"
	}

	return appErr
}

func executeDeleteInstaller(tx *sql.Tx, slug string) *apperror.AppError {
	if appErr := deleteInstallerScriptRow(tx, slug); appErr != nil {
		return appErr
	}

	return deleteInstallerVersionRows(tx, slug)
}

func deleteInstallerScriptRow(tx *sql.Tx, slug string) *apperror.AppError {
	res, errExec := ExecWrapper(tx, SQLDeleteInstallerScript, slug).Destruct()
	if errExec != nil {
		return installerError(errExec, "E_INSTALLER_DELETE_FAILED", slug)
	}

	return verifyScriptRowsAffected(res, slug)
}

func verifyScriptRowsAffected(res sql.Result, slug string) *apperror.AppError {
	affected, err := res.RowsAffected()
	if err != nil {
		return installerError(err, "E_INSTALLER_DELETE_FAILED", slug)
	}
	if affected == 0 {
		return installerError(apperror.ErrNotFound, "E_INSTALLER_NOT_FOUND", slug)
	}

	return nil
}

func deleteInstallerVersionRows(tx *sql.Tx, slug string) *apperror.AppError {
	if _, err := ExecWrapper(tx, SQLDeleteInstallerVersions, slug).Destruct(); err != nil {
		return installerError(err, "E_INSTALLER_DELETE_FAILED", slug)
	}

	return nil
}

func installerError(err error, code, slug string) *apperror.AppError {
	appErr := apperror.Wrap(err, "DeleteInstaller", map[string]any{"slug": slug})
	appErr.Code = code

	return appErr
}

// DeleteInstallerVersion removes a specific version record for an installer.
func (db *DB) DeleteInstallerVersion(slug, version string) error {
	if db == nil || db.conn == nil {
		return apperror.New("DeleteInstallerVersion", "E_INSTALLER_NIL_DB", map[string]any{"slug": slug})
	}
	_, err := ExecWrapper(db.conn, SQLDeleteInstallerExactVersion, slug, version).Destruct()

	return err
}
