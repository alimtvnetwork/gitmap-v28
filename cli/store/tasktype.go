// Package store — tasktype.go manages the TaskType reference table.
package store

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// SeedTaskTypes inserts default task types if not present.
func (db *DB) SeedTaskTypes() error {
	_, err := ExecWrapper(db.conn, constants.SQLSeedTaskTypes).Destruct()
	if err != nil {
		return fmt.Errorf(constants.ErrPendingTaskInsert, err)
	}

	return nil
}

// GetTaskTypeID returns the ID for a named task type.
func (db *DB) GetTaskTypeID(name string) (int64, error) {
	row := QueryRowWrapper(db.conn, constants.SQLSelectTaskTypeByName, name)

	var id int64

	err := row.Scan(&id)
	if err != nil {
		return 0, fmt.Errorf(constants.ErrTaskTypeNotFound, name)
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
		return 0, fmt.Errorf(constants.ErrPendingTaskInsert, insErr)
	}

	lastID, _ := res.LastInsertId()
	if lastID > 0 {
		return lastID, nil
	}

	return db.GetTaskTypeID(name)
}
