// Package store — installer_create.go provides CreateInstaller for installer scripts.
package store

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

// SQLInsertInstallerScript inserts a new installer script record into SQLite.
const SQLInsertInstallerScript = `INSERT INTO installer_scripts (name, slug, description, target_os, version, instructions)
VALUES (?, ?, ?, ?, ?, ?);`

// CreateInstaller inserts a new installer script record into the database.
func (db *DB) CreateInstaller(script *model.InstallerScript) error {
	if script == nil {
		appErr := apperror.New("CreateInstaller", "E_INSTALLER_INVALID_INPUT", map[string]any{
			"error": "script cannot be nil",
		})

		return appErr
	}

	res, errExec := ExecWrapper(db.conn, SQLInsertInstallerScript,
		script.Name,
		script.Slug,
		script.Description,
		script.TargetOS,
		script.Version,
		script.Instructions,
	).Destruct()
	if errExec != nil {
		appErr := apperror.Wrap(errExec, "CreateInstaller", map[string]any{
			"name": script.Name,
			"slug": script.Slug,
		})
		appErr.Code = "E_INSTALLER_CREATE_FAILED"

		return appErr
	}

	if id, errID := res.LastInsertId(); errID == nil && script.ID == 0 {
		script.ID = id
	}

	return nil
}

// SQLUpdateInstallerScript updates an existing installer script record in SQLite.
const SQLUpdateInstallerScript = `UPDATE installer_scripts
SET name = ?, description = ?, target_os = ?, version = ?, instructions = ?, updated_at = CURRENT_TIMESTAMP
WHERE slug = ?;`

// UpdateInstaller updates an existing installer script record in the database.
func (db *DB) UpdateInstaller(script *model.InstallerScript) error {
	if script == nil {
		appErr := apperror.New("UpdateInstaller", "E_INSTALLER_INVALID_INPUT", map[string]any{
			"error": "script cannot be nil",
		})

		return appErr
	}

	return db.execUpdateInstaller(script)
}

func (db *DB) execUpdateInstaller(script *model.InstallerScript) error {
	_, errExec := ExecWrapper(db.conn, SQLUpdateInstallerScript,
		script.Name, script.Description, script.TargetOS,
		script.Version, script.Instructions, script.Slug,
	).Destruct()
	if errExec != nil {
		appErr := apperror.Wrap(errExec, "UpdateInstaller", map[string]any{"slug": script.Slug})
		appErr.Code = "E_INSTALLER_UPDATE_FAILED"

		return appErr
	}

	return nil
}
