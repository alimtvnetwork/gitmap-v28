// Package store — installer_list.go provides ListInstallers to retrieve all installer scripts.
package store

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// SQLSelectAllInstallers queries all installer script records ordered by id.
const SQLSelectAllInstallers = `SELECT id, name, slug, description, target_os, version, instructions, created_at, updated_at FROM installer_scripts ORDER BY id ASC;`

// ListInstallers retrieves all installer scripts from the database.
func (db *DB) ListInstallers() ([]model.InstallerScript, error) {
	if db == nil || db.conn == nil {
		return nil, apperror.New("ListInstallers", "E_INSTALLER_NIL_DB", nil)
	}

	rows, errQuery := QueryWrapper(db.conn, SQLSelectAllInstallers).Destruct()
	if errQuery != nil {
		appErr := apperror.Wrap(errQuery, "ListInstallers", nil)
		appErr.Code = "E_INSTALLER_LIST_FAILED"
		return nil, appErr
	}
	defer rows.Close()

	var scripts []model.InstallerScript
	for rows.Next() {
		var script model.InstallerScript
		if errScan := rows.Scan(
			&script.ID, &script.Name, &script.Slug, &script.Description,
			&script.TargetOS, &script.Version, &script.Instructions,
			&script.CreatedAt, &script.UpdatedAt,
		); errScan != nil {
			appErr := apperror.Wrap(errScan, "ListInstallers", nil)
			appErr.Code = "E_INSTALLER_LIST_FAILED"
			return nil, appErr
		}
		scripts = append(scripts, script)
	}

	if errRows := rows.Err(); errRows != nil {
		appErr := apperror.Wrap(errRows, "ListInstallers", nil)
		appErr.Code = "E_INSTALLER_LIST_FAILED"
		return nil, appErr
	}

	return scripts, nil
}

// ListInstallHistory retrieves all installer script records for history listing.
func (db *DB) ListInstallHistory() ([]model.InstallerScript, error) {
	return db.ListInstallers()
}
