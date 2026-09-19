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
	cluster_role TEXT DEFAULT 'worker',
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

func execAlterColumn(db *sql.DB, stmt string) error {
	_, err := db.Exec(stmt)
	if err != nil && !isBenignAlterError(err) {
		return apperror.WrapSimple(err, "ensureHostColumns")
	}

	return nil
}

func ensureHostColumns(db *sql.DB) error {
	cols := []string{
		"ALTER TABLE ssh_hosts ADD COLUMN port INTEGER DEFAULT 22;",
		"ALTER TABLE ssh_hosts ADD COLUMN encrypted_password TEXT;",
		"ALTER TABLE ssh_hosts ADD COLUMN cluster_role TEXT DEFAULT 'worker';",
	}
	for _, col := range cols {
		if err := execAlterColumn(db, col); err != nil {
			return err
		}
	}

	return nil
}

// SQLCreateSSHConnectionTable defines the DDL for legacy SSHConnection table compatibility.
const SQLCreateSSHConnectionTable = `CREATE TABLE IF NOT EXISTS SSHConnection (
	Alias TEXT PRIMARY KEY,
	IPAddress TEXT NOT NULL,
	Username TEXT NOT NULL,
	EncryptedPassword TEXT NOT NULL,
	KeyPath TEXT,
	OS TEXT DEFAULT 'linux',
	CreatedAt TIMESTAMP NOT NULL
);`

// SQLCreateSSHKnownHostsTable defines the DDL for the ssh_known_hosts table.
const SQLCreateSSHKnownHostsTable = `CREATE TABLE IF NOT EXISTS ssh_known_hosts (
	id TEXT PRIMARY KEY,
	host TEXT NOT NULL,
	key_type TEXT NOT NULL,
	public_key TEXT NOT NULL,
	fingerprint TEXT NOT NULL,
	comment TEXT,
	created_at DATETIME NOT NULL,
	updated_at DATETIME NOT NULL
);`

const SQLCreateSSHKnownHostsIndex = `CREATE INDEX IF NOT EXISTS idx_ssh_known_hosts_host ON ssh_known_hosts (host);`

// EnsureSSHTables creates both ssh_hosts, ssh_history, and SSHConnection tables if they do not exist.
func EnsureSSHTables(db *sql.DB) error {
	if db == nil {
		return nil
	}

	if err := executeTableDDL(db, SQLCreateSSHHostsTable, "EnsureSSHTables_Hosts"); err != nil {
		return err
	}
	if err := ensureHostColumns(db); err != nil {
		return err
	}
	if err := executeTableDDL(db, SQLCreateSSHConnectionTable, "EnsureSSHTables_SSHConnection"); err != nil {
		return err
	}
	_ = executeTableDDL(db, SQLCreateSSHKnownHostsTable, "EnsureSSHTables_KnownHosts")
	_ = executeTableDDL(db, SQLCreateSSHKnownHostsIndex, "EnsureSSHTables_KnownHosts_Index")

	return executeTableDDL(db, SQLCreateSSHHistoryTable, "EnsureSSHTables_History")
}

// RegisterSSHHostMigration creates the ssh_hosts and ssh_history tables if they do not exist.
func RegisterSSHHostMigration(db *sql.DB, version int, isForced bool) error {
	if err := EnsureSSHTables(db); err != nil {
		return apperror.WrapSimple(err, "RegisterSSHHostMigration")
	}

	return nil
}
