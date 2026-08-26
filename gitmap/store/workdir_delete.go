// Package store — workdir_delete.go removes work directories.
package store

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

// DeleteWorkDir removes a work directory by ID or path.
func (db *DB) DeleteWorkDir(idOrPath string) error {
	if db == nil || db.conn == nil {
		return apperror.New("DeleteWorkDir", "E_NIL_DB", map[string]any{"target": idOrPath})
	}

	_, err := ExecWrapper(db.conn, SQLDeleteWorkDir, idOrPath, idOrPath).Destruct()
	return err
}
