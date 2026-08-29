// Package store — installer_get.go provides GetInstallerBySlug to retrieve installer scripts.
package store

import (
	"database/sql"
	"errors"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// SQLSelectInstallerBySlug queries an installer script record by its slug.
const SQLSelectInstallerBySlug = `
SELECT id, name, slug, description, target_os, version, instructions, created_at, updated_at
FROM installer_scripts
WHERE slug = ?
`

// GetInstallerBySlug retrieves an installer script by its unique slug.
func (db *DB) GetInstallerBySlug(slug string) (*model.InstallerScript, error) {
	if db == nil || db.conn == nil {
		appErr := apperror.New("GetInstallerBySlug", "E_INSTALLER_NIL_DB", map[string]any{"slug": slug})
		return nil, appErr
	}

	var script model.InstallerScript
	err := db.conn.QueryRow(SQLSelectInstallerBySlug, slug).Scan(
		&script.ID,
		&script.Name,
		&script.Slug,
		&script.Description,
		&script.TargetOS,
		&script.Version,
		&script.Instructions,
		&script.CreatedAt,
		&script.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		appErr := apperror.Wrap(apperror.ErrNotFound, "GetInstallerBySlug", map[string]any{"slug": slug})
		appErr.Code = "E_INSTALLER_NOT_FOUND"
		return nil, appErr
	}
	if err != nil {
		appErr := apperror.Wrap(err, "GetInstallerBySlug", map[string]any{"slug": slug})
		appErr.Code = "E_INSTALLER_GET_FAILED"
		return nil, appErr
	}

	return &script, nil
}
