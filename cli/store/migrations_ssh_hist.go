package store

import (
	"database/sql"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// SQLCreateSSHHistoryTable defines the DDL statement to initialize the ssh_history table.
const SQLCreateSSHHistoryTable = `CREATE TABLE IF NOT EXISTS ssh_history (id TEXT PRIMARY KEY, host_ip TEXT, joined_at DATETIME, user TEXT);`

// SQLCreateSSHHistoryView defines the DDL statement to create the SSHHistory view.
const SQLCreateSSHHistoryView = `
CREATE VIEW IF NOT EXISTS SSHHistory AS
SELECT
    id AS SSHHistoryId,
    host_ip AS HostIP,
    joined_at AS JoinedAt,
    user AS User
FROM ssh_history;`

// RegisterSSHHistoryMigration applies the migration for the ssh_history table and view.
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

	_, err = db.Exec(SQLCreateSSHHistoryView)
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
