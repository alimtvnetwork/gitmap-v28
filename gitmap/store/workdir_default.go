// Package store — workdir_default.go manages the active default work directory.
package store

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// SetDefaultWorkDir marks a specific work directory as default.
func (db *DB) SetDefaultWorkDir(idOrPath string) error {
	if db == nil || db.conn == nil {
		return apperror.New("SetDefaultWorkDir", "E_NIL_DB", map[string]any{"target": idOrPath})
	}

	_, err := ExecWrapper(db.conn, SQLSetDefaultWorkDir, idOrPath, idOrPath).Destruct()
	return err
}

// GetDefaultWorkDir returns the current default work directory.
func (db *DB) GetDefaultWorkDir() (*model.WorkDir, error) {
	dirs, err := db.ListWorkDirs()
	if err != nil {
		return nil, err
	}
	for _, d := range dirs {
		if d.IsDefault {
			return &d, nil
		}
	}
	if len(dirs) > 0 {
		return &dirs[0], nil
	}
	return nil, apperror.NewSimple("GetDefaultWorkDir", "E_NOT_FOUND")
}
