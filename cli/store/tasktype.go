// Package store — tasktype.go manages the TaskType reference table.
package store

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// SeedTaskTypes inserts default task types if not present.
func (db *DB) SeedTaskTypes() error {
	_, err := ExecWrapper(db.conn, constants.SQLSeedTaskTypes).Destruct()
	if err != nil {
		return apperror.WrapSimple(err, "failed to insert pending task (operation: insert)")
	}

	return nil
}

// GetTaskTypeID returns the ID for a named task type.
func (db *DB) GetTaskTypeID(name string) (int64, error) {
	row := QueryRowWrapper(db.conn, constants.SQLSelectTaskTypeByName, name)

	var id int64

	err := row.Scan(&id)
	if err != nil {
		return 0, apperror.WrapNotFound(err, fmt.Sprintf("task type not found: %s", name))
	}

	return id, nil
}

// EnsureTaskTypeID returns the ID for a named task type, creating it if absent.
func (db *DB) EnsureTaskTypeID(name string) (int64, error) {
	id, err := db.GetTaskTypeID(name)
	if err == nil {
		return id, nil
	}

	res, insErr := ExecWrapper(db.conn, "INSERT OR IGNORE INTO TaskType (Name) VALUES (?)", name).Destruct()
	if insErr != nil {
		return 0, apperror.WrapSimple(insErr, "failed to insert pending task (operation: insert)")
	}

	lastID, _ := res.LastInsertId()
	if lastID > 0 {
		return lastID, nil
	}

	return db.GetTaskTypeID(name)
}
