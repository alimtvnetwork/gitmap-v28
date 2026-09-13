package store

import (
	"database/sql"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// SQLCreateSSHHostsTable defines the DDL statement to initialize the ssh_hosts table.
const SQLCreateSSHHostsTable = `CREATE TABLE IF NOT EXISTS ssh_hosts (id TEXT PRIMARY KEY, alias TEXT, ip TEXT, username TEXT, created_at DATETIME);`

// RegisterSSHHostMigration creates the ssh_hosts table if it does not exist.
func RegisterSSHHostMigration(db *sql.DB, version int, force bool) error {
	if _, err := db.Exec(SQLCreateSSHHostsTable); err != nil {
		e := apperror.WrapSimple(err, "RegisterSSHHostMigration")
		e.Code = "E_INTERNAL_ERROR"

		return e
	}

	return nil
}
