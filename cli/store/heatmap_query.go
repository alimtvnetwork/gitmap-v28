package store

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	_ "modernc.org/sqlite"
)

// OpenReadOnlyHeatmapDB opens a SQLite database connection in read-only mode with query safety pragmas.
func OpenReadOnlyHeatmapDB(dbPath string) (*sql.DB, error) {
	effectivePath := dbPath
	if _, err := os.Stat(effectivePath); err != nil {
		if fallback := findGlobalFallbackDB(filepath.Base(dbPath), dbPath); fallback != "" {
			effectivePath = fallback
		} else {
			return nil, apperror.NewNotFound("store.OpenReadOnlyHeatmapDB", "E1001", fmt.Sprintf("database file not found: %s", dbPath))
		}
	}

	dsn := fmt.Sprintf("file:%s?mode=ro&cache=shared", filepath.ToSlash(effectivePath))
	conn, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, apperror.WrapSimple(err, "store.OpenReadOnlyHeatmapDB.open")
	}

	conn.SetMaxOpenConns(1)
	const pragmas = "PRAGMA query_only = ON; PRAGMA busy_timeout = 5000;"
	if _, err := conn.Exec(pragmas); err != nil {
		_ = conn.Close()
		return nil, apperror.WrapSimple(err, "store.OpenReadOnlyHeatmapDB.pragmas")
	}

	return conn, nil
}

// ExecuteReadOnlyQuery opens a read-only SQLite database connection, executes the query, and extracts dynamic results.
func ExecuteReadOnlyQuery(dbPath string, query string) ([]string, [][]any, error) {
	conn, err := OpenReadOnlyHeatmapDB(dbPath)
	if err != nil {
		return nil, nil, err
	}
	defer conn.Close()

	return ExecuteReadOnlyQueryOnDB(conn, query)
}

// ExecuteReadOnlyQueryOnDB executes a read-only query on an existing connection and dynamically extracts values.
func ExecuteReadOnlyQueryOnDB(conn *sql.DB, query string) ([]string, [][]any, error) {
	rows, err := conn.Query(query)
	if err != nil {
		return nil, nil, apperror.New("store.ExecuteReadOnlyQueryOnDB", "E1002", map[string]any{"err": err.Error()})
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, nil, apperror.WrapSimple(err, "store.ExecuteReadOnlyQueryOnDB.columns")
	}

	var allRows [][]any
	for rows.Next() {
		row, scanErr := scanDynamicRow(rows, len(cols))
		if scanErr != nil {
			return nil, nil, scanErr
		}
		allRows = append(allRows, row)
	}

	if err := rows.Err(); err != nil {
		return nil, nil, apperror.WrapSimple(err, "store.ExecuteReadOnlyQueryOnDB.rowsErr")
	}

	return cols, allRows, nil
}

func scanDynamicRow(rows *sql.Rows, colCount int) ([]any, error) {
	values := make([]any, colCount)
	valuePtrs := make([]any, colCount)
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	if err := rows.Scan(valuePtrs...); err != nil {
		return nil, apperror.WrapSimple(err, "store.scanDynamicRow")
	}

	row := make([]any, colCount)
	for i, val := range values {
		row[i] = normalizeCellValue(val)
	}

	return row, nil
}

func normalizeCellValue(raw any) any {
	if raw == nil {
		return nil
	}
	switch v := raw.(type) {
	case []byte:
		return string(v)
	default:
		return v
	}
}
