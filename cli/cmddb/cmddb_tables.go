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

	tables, listErr := listTablesAndViews(conn)
	if listErr != nil {
		return listErr
	}

	if isJson {
		return outputJson(out, tables)
	}

	renderTablesTerminal(out, dbPath, tables)

	return nil
}

func listTablesAndViews(conn *sql.DB) ([]TableInfo, error) {
	query := `SELECT name, type FROM sqlite_master WHERE type IN ('table', 'view') AND name NOT LIKE 'sqlite_%' ORDER BY type, name`
	rows, err := conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query sqlite_master: %w", err)
	}
	defer rows.Close()

	var tables []TableInfo
	for rows.Next() {
		var name, tableType string
		if scanErr := rows.Scan(&name, &tableType); scanErr != nil {
			return nil, scanErr
		}
		cols, count := fetchTableStats(conn, name, tableType)
		tables = append(tables, TableInfo{
			Name:        name,
			Type:        tableType,
			ColumnCount: cols,
			RowCount:    count,
		})
	}

	return tables, rows.Err()
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

func renderTablesTerminal(out io.Writer, dbPath string, tables []TableInfo) {
	if len(tables) == 0 {
		fmt.Fprintf(out, "Database %s contains no tables or views.\n", dbPath)
		return
	}

	fmt.Fprintf(out, "Database: %s (%d tables/views)\n\n", dbPath, len(tables))
	cols := []termout.Column{
		{Title: "NAME", MinWidth: 24, Align: termout.AlignLeft},
		{Title: "TYPE", MinWidth: 10, Align: termout.AlignLeft},
		{Title: "COLUMNS", MinWidth: 10, Align: termout.AlignRight},
		{Title: "ROWS", MinWidth: 10, Align: termout.AlignRight},
	}

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

	tableStr := termout.RenderTable(termout.TableConfig{
		Columns:    cols,
		Rows:       rows,
		HasBorders: true,
	})
	fmt.Fprint(out, tableStr)
}
