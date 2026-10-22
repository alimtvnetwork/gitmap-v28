package cmddb

import (
	"database/sql"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

// RunDbColumns formats and displays PRAGMA table_info for a specified table.
func RunDbColumns(tableArg, dbArg string, isJson bool, out io.Writer) error {
	cleanTable := strings.TrimSpace(tableArg)
	if cleanTable == "" {
		return apperror.NewValidationError("table name is required (e.g. gitmap db columns <table-name>)")
	}

	conn, dbPath, err := OpenResolvedDb(dbArg)
	if err != nil {
		return err
	}
	defer conn.Close()

	cols, pragmaErr := fetchColumnInfos(conn, cleanTable)
	if pragmaErr != nil {
		return pragmaErr
	}

	if len(cols) == 0 {
		return apperror.NewNotFoundError(fmt.Sprintf("table '%s' not found or has no columns in %s", cleanTable, dbPath))
	}

	if isJson {
		return outputJson(out, cols)
	}

	renderColumnsTerminal(out, dbPath, cleanTable, cols)

	return nil
}

func fetchColumnInfos(conn *sql.DB, tableName string) ([]ColumnInfo, error) {
	escaped := strings.ReplaceAll(tableName, `"`, `""`)
	query := fmt.Sprintf(`PRAGMA table_info("%s")`, escaped)
	rows, err := conn.Query(query)
	if err != nil {
		return nil, apperror.WrapSimple(err, "query table info")
	}
	defer rows.Close()

	var cols []ColumnInfo
	for rows.Next() {
		col, scanErr := scanColumnInfo(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		cols = append(cols, col)
	}

	return cols, rows.Err()
}

func scanColumnInfo(rows *sql.Rows) (ColumnInfo, error) {
	var (
		cid       int
		name      string
		colType   string
		notNull   int
		dfltValue sql.NullString
		pk        int
	)

	err := rows.Scan(&cid, &name, &colType, &notNull, &dfltValue, &pk)
	if err != nil {
		return ColumnInfo{}, err
	}

	var dfltPtr *string
	if dfltValue.Valid {
		dfltPtr = &dfltValue.String
	}

	return ColumnInfo{
		Cid:          cid,
		Name:         name,
		Type:         colType,
		NotNull:      notNull == 1,
		DefaultValue: dfltPtr,
		IsPrimaryKey: pk > 0,
	}, nil
}

func renderColumnsTerminal(out io.Writer, dbPath, tableName string, cols []ColumnInfo) {
	fmt.Fprintf(out, "Table: %s in %s (%d columns)\n\n", tableName, dbPath, len(cols))
	tableCols := []termout.Column{
		{Title: "CID", MinWidth: 5, Align: termout.AlignRight},
		{Title: "COLUMN NAME", MinWidth: 20, Align: termout.AlignLeft},
		{Title: "DATA TYPE", MinWidth: 14, Align: termout.AlignLeft},
		{Title: "NOT NULL", MinWidth: 10, Align: termout.AlignCenter},
		{Title: "DEFAULT", MinWidth: 12, Align: termout.AlignLeft},
		{Title: "PK", MinWidth: 6, Align: termout.AlignCenter},
	}

	rows := make([]termout.Row, 0, len(cols))
	for _, c := range cols {
		rows = append(rows, termout.Row{
			Cells: []string{
				strconv.Itoa(c.Cid),
				c.Name,
				c.Type,
				formatBoolBadge(c.NotNull),
				formatDefaultValue(c.DefaultValue),
				formatPkBadge(c.IsPrimaryKey),
			},
		})
	}

	tableStr := termout.RenderTable(termout.TableConfig{
		Columns:    tableCols,
		Rows:       rows,
		HasBorders: true,
	})
	fmt.Fprint(out, tableStr)
}

func formatBoolBadge(val bool) string {
	if val {
		return "YES"
	}
	return "NO"
}

func formatPkBadge(isPk bool) string {
	if isPk {
		return "YES"
	}
	return "-"
}

func formatDefaultValue(val *string) string {
	if val == nil {
		return "NULL"
	}
	return *val
}
