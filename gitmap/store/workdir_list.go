// Package store — workdir_list.go implements work directory querying and insertion.
package store

import (
	"database/sql"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// EnsureWorkDir registers or updates a work directory.
func (db *DB) EnsureWorkDir(absPath, label string, isDefault bool) (*model.WorkDir, error) {
	if db == nil || db.conn == nil {
		return nil, apperror.New("EnsureWorkDir", "E_NIL_DB", map[string]any{"path": absPath})
	}

	defInt := 0
	if isDefault {
		defInt = 1
	}

	_, err := ExecWrapper(db.conn, SQLUpsertWorkDir, absPath, label, defInt, label, label, defInt, defInt).Destruct()
	if err != nil {
		return nil, apperror.Wrap(err, "EnsureWorkDir", map[string]any{"path": absPath})
	}

	return db.GetWorkDirByPath(absPath)
}

// ListWorkDirs returns all registered work directories.
func (db *DB) ListWorkDirs() ([]model.WorkDir, error) {
	if db == nil || db.conn == nil {
		return nil, apperror.NewSimple("ListWorkDirs", "E_NIL_DB")
	}

	rows, err := QueryWrapper(db.conn, SQLSelectAllWorkDirs).Destruct()
	if err != nil {
		return nil, apperror.WrapSimple(err, "ListWorkDirs")
	}
	defer rows.Close()

	var results []model.WorkDir
	for rows.Next() {
		var wd model.WorkDir
		var defInt int
		var label sql.NullString
		if errScan := rows.Scan(&wd.ID, &wd.AbsolutePath, &label, &defInt, &wd.CreatedAt, &wd.UpdatedAt); errScan == nil {
			wd.Label = label.String
			wd.IsDefault = (defInt == 1)
			results = append(results, wd)
		}
	}
	return results, nil
}

// GetWorkDirByPath retrieves a work directory by its absolute path.
func (db *DB) GetWorkDirByPath(absPath string) (*model.WorkDir, error) {
	row := db.conn.QueryRow("SELECT id, absolute_path, label, is_default, created_at, updated_at FROM work_directories WHERE absolute_path = ?", absPath)
	var wd model.WorkDir
	var defInt int
	var label sql.NullString
	if err := row.Scan(&wd.ID, &wd.AbsolutePath, &label, &defInt, &wd.CreatedAt, &wd.UpdatedAt); err != nil {
		return nil, apperror.Wrap(err, "GetWorkDirByPath", map[string]any{"path": absPath})
	}
	wd.Label = label.String
	wd.IsDefault = (defInt == 1)
	return &wd, nil
}
