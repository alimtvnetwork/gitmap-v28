package cmddb

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// ResolveDbPath resolves a database alias or path to a clean file path.
func ResolveDbPath(dbArg string) (string, error) {
	return ResolveDbPathWithRegistry(dbArg)
}

// OpenResolvedDb opens a SQLite connection to the resolved database path.
func OpenResolvedDb(dbArg string) (*sql.DB, string, error) {
	dbPath, err := ResolveDbPath(dbArg)
	if err != nil {
		return nil, "", err
	}

	if _, statErr := os.Stat(dbPath); statErr != nil {
		return nil, dbPath, apperror.NewNotFoundError(fmt.Sprintf("database file not found: %s", dbPath))
	}

	openRes := store.OpenSQLiteDB(dbPath)
	if openRes.IsFailure() {
		return nil, dbPath, apperror.WrapSimple(openRes.AppError(), "open database '"+dbPath+"'")
	}
	conn := openRes.Value

	IncrementAliasSearchCount(dbArg)

	return conn, dbPath, nil
}

func outputJson(out io.Writer, data any) error {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "marshal json output")
	}
	_, writeErr := fmt.Fprintln(out, string(bytes))

	return writeErr
}
