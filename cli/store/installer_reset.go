package store

import (
	"context"
	"database/sql"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/dbengine"
)

// SQLDeleteAllInstallerVersions deletes all records from installer_versions.
const SQLDeleteAllInstallerVersions = `DELETE FROM installer_versions;`

// SQLDeleteAllInstallerScripts deletes all records from installer_scripts.
const SQLDeleteAllInstallerScripts = `DELETE FROM installer_scripts;`

// SQLDeleteInstallerVersionsBySlug deletes records from installer_versions for a given slug.
const SQLDeleteInstallerVersionsBySlug = `DELETE FROM installer_versions WHERE slug = ?;`

// SQLDeleteInstallerScriptBySlug deletes a record from installer_scripts for a given slug.
const SQLDeleteInstallerScriptBySlug = `DELETE FROM installer_scripts WHERE slug = ?;`

// ResetInstallers deletes installer records and version history for a specific slug or all installers atomically.
func (db *DB) ResetInstallers(slug string, isAll bool) error {
	if appErr := validateResetInput(db, slug, isAll); appErr != nil {
		return appErr
	}

	wrap, appErr := dbengine.WrapDb(db.conn, dbengine.DbSQLite)
	if appErr != nil {
		return appErr
	}

	if appErr := runResetInstallersTx(wrap, slug, isAll); appErr != nil {
		return appErr
	}

	return nil
}

func validateResetInput(db *DB, slug string, isAll bool) *apperror.AppError {
	if db == nil || db.conn == nil {
		return apperror.New("ResetInstallers", "E_INSTALLER_NIL_DB", map[string]any{
			"slug": slug,
			"all":  isAll,
		})
	}

	if !isAll && slug == "" {
		return apperror.New("ResetInstallers", "E_INSTALLER_INVALID_INPUT", map[string]any{
			"error": "slug cannot be empty when all is false",
		})
	}

	return nil
}

func runResetInstallersTx(wrap *dbengine.DbWrapper, slug string, isAll bool) *apperror.AppError {
	ctx := context.Background()

	appErr := wrap.WithTransaction(ctx, func(tx *dbengine.TxWrapper) *apperror.AppError {
		return dispatchResetTx(tx.Tx(), slug, isAll)
	})
	if appErr != nil {
		return ensureResetErrorCode(appErr)
	}

	return nil
}

func dispatchResetTx(tx *sql.Tx, slug string, isAll bool) *apperror.AppError {
	if isAll {
		return resetAllInstallersTx(tx, slug)
	}

	return resetSingleInstallerTx(tx, slug)
}

func ensureResetErrorCode(appErr *apperror.AppError) *apperror.AppError {
	appErr.Code = "E_INSTALLER_RESET_FAILED"

	return appErr
}

func resetAllInstallersTx(tx *sql.Tx, slug string) *apperror.AppError {
	if _, err := ExecWrapper(tx, SQLDeleteAllInstallerVersions).Destruct(); err != nil {
		return resetError(err, true, slug)
	}

	if _, err := ExecWrapper(tx, SQLDeleteAllInstallerScripts).Destruct(); err != nil {
		return resetError(err, true, slug)
	}

	return nil
}

func resetSingleInstallerTx(tx *sql.Tx, slug string) *apperror.AppError {
	if slug == "" {
		return apperror.New("ResetInstallers", "E_INSTALLER_INVALID_INPUT", map[string]any{
			"error": "slug cannot be empty when all is false",
		})
	}

	return executeSingleInstallerReset(tx, slug)
}

func executeSingleInstallerReset(tx *sql.Tx, slug string) *apperror.AppError {
	if _, err := ExecWrapper(tx, SQLDeleteInstallerVersionsBySlug, slug).Destruct(); err != nil {
		return resetError(err, false, slug)
	}

	if _, err := ExecWrapper(tx, SQLDeleteInstallerScriptBySlug, slug).Destruct(); err != nil {
		return resetError(err, false, slug)
	}

	return nil
}

func resetError(err error, isAll bool, slug string) *apperror.AppError {
	appErr := apperror.Wrap(err, "ResetInstallers", map[string]any{"all": isAll, "slug": slug})
	appErr.Code = "E_INSTALLER_RESET_FAILED"

	return appErr
}
