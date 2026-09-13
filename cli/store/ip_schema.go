package store

import (
	"database/sql"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

const sqlCreateIPSnapshot = `CREATE TABLE IF NOT EXISTS IPSnapshot (
	IPSnapshotId INTEGER PRIMARY KEY AUTOINCREMENT,
	InterfaceName TEXT NOT NULL,
	IP TEXT NOT NULL,
	Netmask TEXT NOT NULL,
	Gateway TEXT NOT NULL DEFAULT '',
	DNS TEXT NOT NULL DEFAULT '',
	IsDHCP INTEGER NOT NULL DEFAULT 0,
	Timestamp INTEGER NOT NULL,
	Notes TEXT NULL,
	Comments TEXT NULL
);`

// EnsureIPSnapshotTable ensures that the IPSnapshot table exists in the database.
func EnsureIPSnapshotTable(conn *sql.DB) *apperror.AppError {
	if conn == nil {
		return apperror.NewSimple("store.EnsureIPSnapshotTable", "E_NIL_CONN")
	}

	if _, err := conn.Exec(sqlCreateIPSnapshot); err != nil {
		return apperror.WrapSimple(err, "store.EnsureIPSnapshotTable")
	}

	return nil
}

// EnsureIPSnapshotTable ensures table existence on the DB instance.
func (db *DB) EnsureIPSnapshotTable() *apperror.AppError {
	if db == nil || db.conn == nil {
		return apperror.NewSimple("store.DB.EnsureIPSnapshotTable", "E_NIL_CONN")
	}

	return EnsureIPSnapshotTable(db.conn)
}
