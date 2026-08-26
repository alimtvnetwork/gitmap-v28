// Package store — installer_delete.go provides DeleteInstaller for installer scripts.
package store

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

// SQLDeleteInstallerScript deletes an installer script record by its slug.
const SQLDeleteInstallerScript = `DELETE FROM installer_scripts WHERE slug = ?;`

// SQLDeleteInstallerVersions deletes all installer version records for a given slug.
const SQLDeleteInstallerVersions = `DELETE FROM installer_versions WHERE slug = ?;`

// SQLDeleteInstallerExactVersion deletes a specific version record for a given slug.
const SQLDeleteInstallerExactVersion = `DELETE FROM installer_versions WHERE slug = ? AND version = ?;`

// DeleteInstaller removes an installer script and associated version records by slug.
func (db *DB) DeleteInstaller(slug string) error {
	if db == nil || db.conn == nil {
		return apperror.New("DeleteInstaller", "E_INSTALLER_NIL_DB", map[string]any{"slug": slug})
	}
	if slug == "" {
		return apperror.New("DeleteInstaller", "E_INSTALLER_INVALID_INPUT", map[string]any{"error": "slug cannot be empty"})
	}

	res, errExec := ExecWrapper(db.conn, SQLDeleteInstallerScript, slug).Destruct()
	if errExec != nil {
		appErr := apperror.Wrap(errExec, "DeleteInstaller", map[string]any{"slug": slug})
		appErr.Code = "E_INSTALLER_DELETE_FAILED"
		return appErr
	}

	affected, errAffected := res.RowsAffected()
	if errAffected != nil {
		appErr := apperror.Wrap(errAffected, "DeleteInstaller", map[string]any{"slug": slug})
		appErr.Code = "E_INSTALLER_DELETE_FAILED"
		return appErr
	}
	if affected == 0 {
		appErr := apperror.Wrap(apperror.ErrNotFound, "DeleteInstaller", map[string]any{"slug": slug})
		appErr.Code = "E_INSTALLER_NOT_FOUND"
		return appErr
	}

	if _, errVer := ExecWrapper(db.conn, SQLDeleteInstallerVersions, slug).Destruct(); errVer != nil {
		appErr := apperror.Wrap(errVer, "DeleteInstaller", map[string]any{"slug": slug})
		appErr.Code = "E_INSTALLER_DELETE_FAILED"
		return appErr
	}

	return nil
}

// DeleteInstallerVersion removes a specific version record for an installer.
func (db *DB) DeleteInstallerVersion(slug, version string) error {
	if db == nil || db.conn == nil {
		return apperror.New("DeleteInstallerVersion", "E_INSTALLER_NIL_DB", map[string]any{"slug": slug})
	}
	_, err := ExecWrapper(db.conn, SQLDeleteInstallerExactVersion, slug, version).Destruct()
	return err
}
