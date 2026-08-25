package store

import (
	"database/sql"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

// RegisterSSHHistoryMigration applies the migration for the ssh_history table.
func RegisterSSHHistoryMigration(db *sql.DB, v int, force bool) error {
	query := `CREATE TABLE IF NOT EXISTS ssh_history (id TEXT PRIMARY KEY, host_ip TEXT, joined_at DATETIME, user TEXT);`

	_, err := db.Exec(query)
	if err != nil {
		return &apperror.AppError{
			Op:    "RegisterSSHHistoryMigration",
			Code:  "E_INTERNAL_ERROR",
			Cause: err,
			Ctx:   map[string]any{"v": v, "force": force},
		}
	}

	return nil
}
