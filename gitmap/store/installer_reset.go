// Package store — installer_reset.go provides ResetInstallers to reset installer scripts and versions.
package store

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

// SQLDeleteAllInstallerVersions deletes all records from installer_versions.
const SQLDeleteAllInstallerVersions = `DELETE FROM installer_versions;`

// SQLDeleteAllInstallerScripts deletes all records from installer_scripts.
const SQLDeleteAllInstallerScripts = `DELETE FROM installer_scripts;`

// SQLDeleteInstallerVersionsBySlug deletes records from installer_versions for a given slug.
const SQLDeleteInstallerVersionsBySlug = `DELETE FROM installer_versions WHERE slug = ?;`

// SQLDeleteInstallerScriptBySlug deletes a record from installer_scripts for a given slug.
const SQLDeleteInstallerScriptBySlug = `DELETE FROM installer_scripts WHERE slug = ?;`

// ResetInstallers deletes installer records and version history for a specific slug or all installers.
func (db *DB) ResetInstallers(slug string, all bool) error {
	if db == nil || db.conn == nil {
		return apperror.New("ResetInstallers", "E_INSTALLER_NIL_DB", map[string]any{
			"slug": slug,
			"all":  all,
		})
	}

	if !all && slug == "" {
		return apperror.New("ResetInstallers", "E_INSTALLER_INVALID_INPUT", map[string]any{
			"error": "slug cannot be empty when all is false",
		})
	}

	if all {
		if _, errExec := ExecWrapper(db.conn, SQLDeleteAllInstallerVersions).Destruct(); errExec != nil {
			appErr := apperror.Wrap(errExec, "ResetInstallers", map[string]any{"all": all, "slug": slug})
			appErr.Code = "E_INSTALLER_RESET_FAILED"

			return appErr
		}

		if _, errExec := ExecWrapper(db.conn, SQLDeleteAllInstallerScripts).Destruct(); errExec != nil {
			appErr := apperror.Wrap(errExec, "ResetInstallers", map[string]any{"all": all, "slug": slug})
			appErr.Code = "E_INSTALLER_RESET_FAILED"

			return appErr
		}

		return nil
	}

	if _, errExec := ExecWrapper(db.conn, SQLDeleteInstallerVersionsBySlug, slug).Destruct(); errExec != nil {
		appErr := apperror.Wrap(errExec, "ResetInstallers", map[string]any{"all": all, "slug": slug})
		appErr.Code = "E_INSTALLER_RESET_FAILED"

		return appErr
	}

	if _, errExec := ExecWrapper(db.conn, SQLDeleteInstallerScriptBySlug, slug).Destruct(); errExec != nil {
		appErr := apperror.Wrap(errExec, "ResetInstallers", map[string]any{"all": all, "slug": slug})
		appErr.Code = "E_INSTALLER_RESET_FAILED"

		return appErr
	}

	return nil
}
