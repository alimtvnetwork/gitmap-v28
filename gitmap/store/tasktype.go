// Package store — tasktype.go manages the TaskType reference table.
package store

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// SeedTaskTypes inserts default task types if not present.
func (db *DB) SeedTaskTypes() error {
	_, err := db.conn.Exec(constants.SQLSeedTaskTypes)
	if err != nil {
		return fmt.Errorf(constants.ErrPendingTaskInsert, err)
	}

	return nil
}

// GetTaskTypeID returns the ID for a named task type.
func (db *DB) GetTaskTypeID(name string) (int64, error) {
	row := db.conn.QueryRow(constants.SQLSelectTaskTypeByName, name)

	var id int64

	err := row.Scan(&id)
	if err != nil {
		return 0, fmt.Errorf(constants.ErrTaskTypeNotFound, name)
	}

	return id, nil
}
