package cmddb

import (
	"database/sql"
	"net/http"
	"strings"
)

// RestSchemaEntry represents a DDL entry for a database object.
type RestSchemaEntry struct {
	Name string `json:"Name"`
	Type string `json:"Type"`
	Sql  string `json:"Sql"`
}

// handleTables returns table and view metadata for an alias.
func handleTables(mgr *MultiDbManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		alias := r.PathValue("alias")
		tables, err := mgr.GetTables(alias)
		if err != nil {
			writeRestEnvelope(w, http.StatusNotFound, nil, nil, err)

			return
		}
		writeRestEnvelope(w, http.StatusOK, tables, nil, nil)
	}
}

// handleSchema returns DDL statements for all or specific tables.
func handleSchema(mgr *MultiDbManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		alias := r.PathValue("alias")
		tableName := r.URL.Query().Get("table")
		conn, _, err := mgr.GetDB(alias)
		if err != nil {
			writeRestEnvelope(w, http.StatusNotFound, nil, nil, err)

			return
		}

		entries, qErr := querySchemaEntries(conn, tableName)
		if qErr != nil {
			writeRestEnvelope(w, http.StatusInternalServerError, nil, nil, qErr)

			return
		}
		writeRestEnvelope(w, http.StatusOK, entries, nil, nil)
	}
}

func querySchemaEntries(conn *sql.DB, tableName string) ([]RestSchemaEntry, error) {
	q, args := buildSchemaQuery(tableName)
	rows, err := conn.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanSchemaRows(rows)
}

func buildSchemaQuery(tableName string) (string, []any) {
	if tableName != "" {
		q := `SELECT name, type, COALESCE(sql, '') FROM sqlite_master WHERE tbl_name = ? AND name NOT LIKE 'sqlite_%' ORDER BY type, name`

		return q, []any{tableName}
	}
	q := `SELECT name, type, COALESCE(sql, '') FROM sqlite_master WHERE sql IS NOT NULL AND name NOT LIKE 'sqlite_%' ORDER BY type, name`

	return q, nil
}

func scanSchemaRows(rows *sql.Rows) ([]RestSchemaEntry, error) {
	var entries []RestSchemaEntry
	for rows.Next() {
		var e RestSchemaEntry
		if err := rows.Scan(&e.Name, &e.Type, &e.Sql); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}

	return entries, rows.Err()
}

// handleFindColumns searches column names across all tables.
func handleFindColumns(mgr *MultiDbManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		alias := r.PathValue("alias")
		nameQuery := strings.ToLower(r.URL.Query().Get("name"))
		matches, err := searchMatchingColumns(mgr, alias, nameQuery)
		if err != nil {
			writeRestEnvelope(w, http.StatusNotFound, nil, nil, err)

			return
		}
		writeRestEnvelope(w, http.StatusOK, matches, nil, nil)
	}
}

func searchMatchingColumns(mgr *MultiDbManager, alias, query string) ([]RestColumnInfo, error) {
	tables, err := mgr.GetTables(alias)
	if err != nil {
		return nil, err
	}

	var results []RestColumnInfo
	for _, tab := range tables {
		cols, _ := mgr.GetColumns(alias, tab.Name)
		for _, col := range cols {
			if query == "" || strings.Contains(strings.ToLower(col.Name), query) {
				results = append(results, col)
			}
		}
	}

	return results, nil
}
