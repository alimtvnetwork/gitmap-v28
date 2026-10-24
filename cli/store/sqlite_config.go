package store

import (
	"database/sql"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
	storeconstants "github.com/alimtvnetwork/gitmap-v28/cli/store/constants"
)

// OpenSQLiteDB opens a SQLite database and applies connection limits and pragmas.
func OpenSQLiteDB(path string) SQLiteDBWrap {
	conn, err := sql.Open(storeconstants.DriverSQLite, path)
	if err != nil {
		return result.Fail[*sql.DB](apperror.Wrap(err, "OpenSQLiteDB.open", map[string]any{
			"path":   path,
			"driver": storeconstants.DriverSQLite,
		}))
	}

	if cfgErr := ConfigureSQLiteConn(conn); cfgErr != nil {
		_ = conn.Close()

		return result.Fail[*sql.DB](cfgErr.WithContext("path", path).WithContext("driver", storeconstants.DriverSQLite))
	}

	return result.Ok(conn)
}

// ConfigureSQLiteConn sets connection limits and configures standard SQLite PRAGMAs.
func ConfigureSQLiteConn(conn *sql.DB) *apperror.AppError {
	conn.SetMaxOpenConns(1)

	return storeconstants.ApplyDefaultPragmas(conn)
}
