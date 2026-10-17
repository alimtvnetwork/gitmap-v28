package cmddb

import (
	"database/sql"
	"fmt"
	"io"
	"strings"
)

// RunDbSchema dumps formatted CREATE TABLE and index DDL for tables.
func RunDbSchema(tableArg, dbArg string, isJson bool, out io.Writer) error {
	conn, dbPath, err := OpenResolvedDb(dbArg)
	if err != nil {
		return err
	}
	defer conn.Close()

	schemas, fetchErr := fetchSchemas(conn, tableArg, dbPath)
	if fetchErr != nil {
		return fetchErr
	}

	if isJson {
		return outputJson(out, schemas)
	}

	renderSchemasTerminal(out, dbPath, schemas)

	return nil
}

func fetchSchemas(conn *sql.DB, tableArg, dbPath string) ([]SchemaInfo, error) {
	if strings.TrimSpace(tableArg) != "" {
		return fetchSingleTableSchema(conn, tableArg, dbPath)
	}

	return fetchAllTableSchemas(conn)
}

func fetchSingleTableSchema(conn *sql.DB, tableName, dbPath string) ([]SchemaInfo, error) {
	query := `SELECT type, name, sql FROM sqlite_master WHERE (name = ? OR name = ? COLLATE NOCASE) AND sql IS NOT NULL LIMIT 1`
	var schemaType, name, ddl string
	err := conn.QueryRow(query, tableName, tableName).Scan(&schemaType, &name, &ddl)
	if err != nil {
		return nil, fmt.Errorf("table or view '%s' not found in database %s", tableName, dbPath)
	}

	indexes := fetchTableIndexes(conn, name)
	info := SchemaInfo{
		TableName: name,
		Type:      schemaType,
		SQL:       ensureSemicolon(ddl),
		Indexes:   indexes,
	}

	return []SchemaInfo{info}, nil
}

func fetchAllTableSchemas(conn *sql.DB) ([]SchemaInfo, error) {
	query := `SELECT type, name, sql FROM sqlite_master WHERE type IN ('table', 'view') AND name NOT LIKE 'sqlite_%' AND sql IS NOT NULL ORDER BY type, name`
	rows, err := conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query all schemas: %w", err)
	}
	defer rows.Close()

	var schemas []SchemaInfo
	for rows.Next() {
		var schemaType, name, ddl string
		if scanErr := rows.Scan(&schemaType, &name, &ddl); scanErr != nil {
			return nil, scanErr
		}
		indexes := fetchTableIndexes(conn, name)
		schemas = append(schemas, SchemaInfo{
			TableName: name,
			Type:      schemaType,
			SQL:       ensureSemicolon(ddl),
			Indexes:   indexes,
		})
	}

	return schemas, rows.Err()
}

func fetchTableIndexes(conn *sql.DB, tableName string) []string {
	query := `SELECT sql FROM sqlite_master WHERE type = 'index' AND tbl_name = ? AND sql IS NOT NULL ORDER BY name`
	rows, err := conn.Query(query, tableName)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var indexes []string
	for rows.Next() {
		var ddl string
		if rows.Scan(&ddl) == nil && ddl != "" {
			indexes = append(indexes, ensureSemicolon(ddl))
		}
	}

	return indexes
}

func ensureSemicolon(s string) string {
	trimmed := strings.TrimSpace(s)
	if strings.HasSuffix(trimmed, ";") {
		return trimmed
	}

	return trimmed + ";"
}

func renderSchemasTerminal(out io.Writer, dbPath string, schemas []SchemaInfo) {
	if len(schemas) == 0 {
		fmt.Fprintf(out, "No schemas found in database %s.\n", dbPath)
		return
	}

	fmt.Fprintf(out, "-- Database: %s\n\n", dbPath)
	for _, s := range schemas {
		fmt.Fprintf(out, "-- %s: %s\n", strings.ToUpper(s.Type), s.TableName)
		fmt.Fprintln(out, s.SQL)
		for _, idx := range s.Indexes {
			fmt.Fprintln(out, idx)
		}
		fmt.Fprintln(out)
	}
}
