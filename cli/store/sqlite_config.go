package store

import (
	"database/sql"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	storeconstants "github.com/alimtvnetwork/gitmap-v28/cli/store/constants"
)

// OpenSQLiteDB opens a SQLite database and applies connection limits and pragmas.
func OpenSQLiteDB(path string) (*sql.DB, *apperror.AppError) {
	conn, err := sql.Open(storeconstants.DriverSQLite, path)
	if err != nil {
		return nil, apperror.WrapSimple(err, "OpenSQLiteDB.open")
	}

	if err := ConfigureSQLiteConn(conn); err != nil {
		_ = conn.Close()

		return nil, apperror.WrapSimple(err, "OpenSQLiteDB.configure")
	}

	return conn, nil
}

// ConfigureSQLiteConn sets connection limits and configures standard SQLite PRAGMAs.
func ConfigureSQLiteConn(conn *sql.DB) error {
	conn.SetMaxOpenConns(1)

	return storeconstants.ApplyDefaultPragmas(conn)
}
