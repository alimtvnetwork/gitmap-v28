package cmddb

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"

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
		return nil, dbPath, fmt.Errorf("database file not found: %s", dbPath)
	}

	conn, openErr := store.OpenSQLiteDB(dbPath)
	if openErr != nil {
		return nil, dbPath, fmt.Errorf("open database '%s': %w", dbPath, openErr)
	}

	IncrementAliasSearchCount(dbArg)

	return conn, dbPath, nil
}

func outputJson(out io.Writer, data any) error {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal json output: %w", err)
	}
	_, writeErr := fmt.Fprintln(out, string(bytes))

	return writeErr
}
