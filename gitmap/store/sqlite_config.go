package store

import (
	"database/sql"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

var sqlitePragmas = []string{
	constants.SQLPragmaBusyTimeout5s,
	constants.SQLPragmaJournalWAL,
	constants.SQLPragmaSynchronousNor,
	constants.SQLEnableFK,
}

// OpenSQLiteDB opens a SQLite database and applies connection limits and pragmas.
func OpenSQLiteDB(path string) (*sql.DB, *apperror.AppError) {
	conn, err := sql.Open("sqlite", path)
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

	return applySQLitePragmas(conn)
}

func applySQLitePragmas(conn *sql.DB) error {
	for _, pragma := range sqlitePragmas {
		if _, err := conn.Exec(pragma); err != nil {
			return fmt.Errorf("apply sqlite pragma %q: %w", pragma, err)
		}
	}

	return nil
}
