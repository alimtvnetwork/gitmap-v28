package store

import (
	"database/sql"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

// RegisterSSHHostMigration creates the ssh_hosts table if it does not exist.
func RegisterSSHHostMigration(db *sql.DB, version int, force bool) error {
	const query = `CREATE TABLE IF NOT EXISTS ssh_hosts (id TEXT PRIMARY KEY, alias TEXT, ip TEXT, username TEXT, created_at DATETIME);`

	if _, err := db.Exec(query); err != nil {
		e := apperror.Wrap(err, "RegisterSSHHostMigration", nil)
		e.Code = "E_INTERNAL_ERROR"
		return e
	}

	return nil
}
