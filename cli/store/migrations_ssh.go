package store

import (
	"database/sql"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// SQLCreateSSHHostsTable defines the DDL statement to initialize the ssh_hosts table.
const SQLCreateSSHHostsTable = `CREATE TABLE IF NOT EXISTS ssh_hosts (
	id TEXT PRIMARY KEY,
	alias TEXT,
	ip TEXT,
	username TEXT,
	port INTEGER DEFAULT 22,
	encrypted_password TEXT,
	created_at DATETIME
);`

func executeTableDDL(db *sql.DB, ddl string, op string) error {
	if _, err := db.Exec(ddl); err != nil {
		appErr := apperror.WrapSimple(err, op)
		appErr.Code = "E_INTERNAL_ERROR"

		return appErr
	}

	return nil
}

func ensureHostColumns(db *sql.DB) {
	_, _ = db.Exec("ALTER TABLE ssh_hosts ADD COLUMN port INTEGER DEFAULT 22;")
	_, _ = db.Exec("ALTER TABLE ssh_hosts ADD COLUMN encrypted_password TEXT;")
}

// EnsureSSHTables creates both ssh_hosts and ssh_history tables if they do not exist.
func EnsureSSHTables(db *sql.DB) error {
	if db == nil {
		return nil
	}

	if err := executeTableDDL(db, SQLCreateSSHHostsTable, "EnsureSSHTables_Hosts"); err != nil {
		return err
	}
	ensureHostColumns(db)

	return executeTableDDL(db, SQLCreateSSHHistoryTable, "EnsureSSHTables_History")
}

// RegisterSSHHostMigration creates the ssh_hosts and ssh_history tables if they do not exist.
func RegisterSSHHostMigration(db *sql.DB, version int, isForced bool) error {
	if err := EnsureSSHTables(db); err != nil {
		return apperror.WrapSimple(err, "RegisterSSHHostMigration")
	}

	return nil
}
