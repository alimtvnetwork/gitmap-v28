// Package store — installer_version.go provides SaveVersion for installer versions.
package store

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// SQLInsertInstallerVersion inserts a new installer version record into SQLite.
const SQLInsertInstallerVersion = `INSERT INTO installer_versions (script_id, slug, version, target_os, instructions)
VALUES (?, ?, ?, ?, ?);`

// SaveVersion inserts a new installer version snapshot into the database.
func (db *DB) SaveVersion(version *model.InstallerVersion) error {
	if db == nil || db.conn == nil {
		appErr := apperror.New("SaveVersion", "E_INSTALLER_NIL_DB", map[string]any{})
		return appErr
	}

	if version == nil {
		appErr := apperror.New("SaveVersion", "E_INSTALLER_INVALID_INPUT", map[string]any{
			"error": "version cannot be nil",
		})

		return appErr
	}

	res, errExec := ExecWrapper(db.conn, SQLInsertInstallerVersion,
		version.ScriptID,
		version.Slug,
		version.Version,
		version.TargetOS,
		version.Instructions,
	).Destruct()
	if errExec != nil {
		appErr := apperror.Wrap(errExec, "SaveVersion", map[string]any{
			"scriptId": version.ScriptID,
			"slug":     version.Slug,
			"version":  version.Version,
		})
		appErr.Code = "E_INSTALLER_SAVE_VERSION_FAILED"

		return appErr
	}

	if id, errID := res.LastInsertId(); errID == nil && version.ID == 0 {
		version.ID = id
	}

	return nil
}
