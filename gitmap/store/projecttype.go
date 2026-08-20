// Package store — projecttype.go manages the ProjectTypes reference table.
package store

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// SeedProjectTypes inserts all supported project types if not present.
func (db *DB) SeedProjectTypes() error {
	_, err := ExecWrapper(db.conn, constants.SQLSeedProjectTypes).Destruct()

	return err
}
