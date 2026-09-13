package store

import (
	"database/sql"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// SQLCreateSSHHistoryTable defines the DDL statement to initialize the ssh_history table.
const SQLCreateSSHHistoryTable = `CREATE TABLE IF NOT EXISTS ssh_history (id TEXT PRIMARY KEY, host_ip TEXT, joined_at DATETIME, user TEXT);`

// RegisterSSHHistoryMigration applies the migration for the ssh_history table.
func RegisterSSHHistoryMigration(db *sql.DB, v int, force bool) error {
	_, err := db.Exec(SQLCreateSSHHistoryTable)
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
