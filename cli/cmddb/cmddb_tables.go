package cmddb

import (
	"database/sql"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

// RunDbTables lists tables, views, column counts, and record counts.
func RunDbTables(dbArg string, isJson bool, out io.Writer) error {
	conn, dbPath, err := OpenResolvedDb(dbArg)
	if err != nil {
		return err
	}
	defer conn.Close()

	return executeDbTablesList(conn, dbPath, isJson, out)
}

func executeDbTablesList(conn *sql.DB, dbPath string, isJson bool, out io.Writer) error {
	tables, listErr := ListTablesAndViews(conn, dbPath)
	if listErr != nil {
		return listErr
	}
	if isJson {
		return outputJson(out, tables)
	}

	renderTablesTerminal(out, dbPath, tables)

	return nil
}

// ListTablesAndViews fetches table and view metadata, using in-memory cache when dbPath is known.
func ListTablesAndViews(conn *sql.DB, dbPath string) ([]TableInfo, error) {
	if dbPath != "" {
		return GetCachedTablesAndViews(conn, dbPath)
	}

	return listTablesAndViews(conn)
}

type rawTableEntry struct {
	name      string
	tableType string
}

func listTablesAndViews(conn *sql.DB) ([]TableInfo, error) {
	query := `SELECT name, type FROM sqlite_master WHERE type IN ('table', 'view') AND name NOT LIKE 'sqlite_%' ORDER BY type, name`
	rows, err := conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query sqlite_master: %w", err)
	}
	defer rows.Close()

	rawEntries, err := collectRawTableEntries(rows)
	if err != nil {
		return nil, err
	}

	return populateTableInfoList(conn, rawEntries), nil
}

func populateTableInfoList(conn *sql.DB, rawEntries []rawTableEntry) []TableInfo {
	tables := make([]TableInfo, 0, len(rawEntries))
	for _, entry := range rawEntries {
		cols, count := fetchTableStats(conn, entry.name, entry.tableType)
		tables = append(tables, TableInfo{
			Name:        entry.name,
			Type:        entry.tableType,
			ColumnCount: cols,
			RowCount:    count,
		})
	}

	return tables
}

func collectRawTableEntries(rows *sql.Rows) ([]rawTableEntry, error) {
	var entries []rawTableEntry
	for rows.Next() {
		entry, scanErr := scanRawTableEntry(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		entries = append(entries, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

func scanRawTableEntry(rows *sql.Rows) (rawTableEntry, error) {
	var name, tableType string
	if err := rows.Scan(&name, &tableType); err != nil {
		return rawTableEntry{}, err
	}

	return rawTableEntry{name: name, tableType: tableType}, nil
}

func fetchTableStats(conn *sql.DB, name, tableType string) (int, int64) {
	colCount := countTableColumns(conn, name)
	rowCount := countTableRows(conn, name)

	return colCount, rowCount
}

func countTableColumns(conn *sql.DB, name string) int {
	var count int
	if err := conn.QueryRow("SELECT count(*) FROM pragma_table_info(?)", name).Scan(&count); err != nil {
		return 0
	}

	return count
}

func countTableRows(conn *sql.DB, name string) int64 {
	escaped := strings.ReplaceAll(name, `"`, `""`)
	safeQuery := fmt.Sprintf(`SELECT count(*) FROM "%s"`, escaped)
	var count int64
	if err := conn.QueryRow(safeQuery).Scan(&count); err != nil {
		return 0
	}

	return count
}

func buildTableOutputColumns() []termout.Column {
	return []termout.Column{
		{Title: "NAME", MinWidth: 24, Align: termout.AlignLeft},
		{Title: "TYPE", MinWidth: 10, Align: termout.AlignLeft},
		{Title: "COLUMNS", MinWidth: 10, Align: termout.AlignRight},
		{Title: "ROWS", MinWidth: 10, Align: termout.AlignRight},
	}
}

func buildTableOutputRows(tables []TableInfo) []termout.Row {
	rows := make([]termout.Row, 0, len(tables))
	for _, t := range tables {
		rows = append(rows, termout.Row{
			Cells: []string{
				t.Name,
				strings.ToUpper(t.Type),
				strconv.Itoa(t.ColumnCount),
				strconv.FormatInt(t.RowCount, 10),
			},
		})
	}

	return rows
}

func renderTablesTerminal(out io.Writer, dbPath string, tables []TableInfo) {
	if len(tables) == 0 {
		fmt.Fprintf(out, "Database %s contains no tables or views.\n", dbPath)
		return
	}

	fmt.Fprintf(out, "Database: %s (%d tables/views)\n\n", dbPath, len(tables))
	tableRows := buildTableOutputRows(tables)
	tableStr := termout.RenderTable(termout.TableConfig{
		Columns:    buildTableOutputColumns(),
		Rows:       tableRows,
		HasBorders: true,
	})
	fmt.Fprint(out, tableStr)
}
