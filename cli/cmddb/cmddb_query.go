package cmddb

import (
	"database/sql"
	"fmt"
	"io"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

// RunDbQuery executes a read-only SQL query and renders tabular or JSON output.
func RunDbQuery(sqlQuery, dbArg string, isJson bool, out io.Writer) error {
	cleanSql := strings.TrimSpace(sqlQuery)
	if cleanSql == "" {
		return apperror.NewValidationError("SQL query string is required (e.g. gitmap db query \"SELECT ...\")")
	}

	if !isQuerySafe(cleanSql) {
		return apperror.NewValidationError("only read-only queries (SELECT, PRAGMA, EXPLAIN, WITH) are permitted")
	}

	conn, _, err := OpenResolvedDb(dbArg)
	if err != nil {
		return err
	}
	defer conn.Close()

	res, queryErr := executeDynamicQuery(conn, cleanSql)
	if queryErr != nil {
		return queryErr
	}

	if isJson {
		return outputJson(out, res)
	}

	renderQueryTerminal(out, res)

	return nil
}

func isQuerySafe(sql string) bool {
	upper := strings.ToUpper(strings.TrimSpace(sql))

	return strings.HasPrefix(upper, "SELECT") ||
		strings.HasPrefix(upper, "PRAGMA") ||
		strings.HasPrefix(upper, "EXPLAIN") ||
		strings.HasPrefix(upper, "WITH")
}

// sqlStatementVerbs are the leading keywords that mark an argument as a SQL
// statement (read or write) rather than a table name.
var sqlStatementVerbs = map[string]bool{
	"SELECT": true, "WITH": true, "PRAGMA": true, "EXPLAIN": true,
	"INSERT": true, "UPDATE": true, "DELETE": true, "REPLACE": true,
	"CREATE": true, "DROP": true, "ALTER": true, "VACUUM": true,
	"REINDEX": true, "ANALYZE": true, "ATTACH": true, "DETACH": true,
	"BEGIN": true, "COMMIT": true, "ROLLBACK": true, "SAVEPOINT": true,
	"RELEASE": true, "TRUNCATE": true,
}

// isSqlStatement reports whether the text begins with a SQL statement verb,
// regardless of whether the statement is read-only.
func isSqlStatement(sql string) bool {
	fields := strings.Fields(strings.ToUpper(strings.TrimSpace(sql)))
	if len(fields) == 0 {
		return false
	}

	return sqlStatementVerbs[fields[0]]
}

func executeDynamicQuery(conn *sql.DB, query string, args ...any) (QueryResult, error) {
	rows, err := conn.Query(query, args...)
	if err != nil {
		return QueryResult{}, apperror.WrapSimple(err, "execute query")
	}
	defer rows.Close()

	colNames, err := rows.Columns()
	if err != nil {
		return QueryResult{}, apperror.WrapSimple(err, "read query columns")
	}

	var data [][]any
	for rows.Next() {
		row, scanErr := scanDynamicRow(rows, len(colNames))
		if scanErr != nil {
			return QueryResult{}, scanErr
		}
		data = append(data, row)
	}

	return QueryResult{
		Columns: colNames,
		Rows:    data,
		Count:   len(data),
	}, rows.Err()
}

func scanDynamicRow(rows *sql.Rows, colCount int) ([]any, error) {
	values := make([]any, colCount)
	valuePtrs := make([]any, colCount)
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	if err := rows.Scan(valuePtrs...); err != nil {
		return nil, apperror.WrapSimple(err, "scan query row")
	}

	cleanRow := make([]any, colCount)
	for i, v := range values {
		cleanRow[i] = normalizeCellValue(v)
	}

	return cleanRow, nil
}

func normalizeCellValue(val any) any {
	if b, ok := val.([]byte); ok {
		return string(b)
	}

	return val
}

func renderQueryTerminal(out io.Writer, res QueryResult) {
	if len(res.Columns) == 0 {
		fmt.Fprintln(out, "Empty query result.")
		return
	}

	cols := make([]termout.Column, len(res.Columns))
	for i, c := range res.Columns {
		cols[i] = termout.Column{
			Title:    c,
			MinWidth: len(c) + 2,
			Align:    termout.AlignLeft,
		}
	}

	rows := make([]termout.Row, 0, len(res.Rows))
	for _, r := range res.Rows {
		cells := make([]string, len(r))
		for i, cell := range r {
			cells[i] = formatCellString(cell)
		}
		rows = append(rows, termout.Row{Cells: cells})
	}

	tableStr := termout.RenderTable(termout.TableConfig{
		Columns:    cols,
		Rows:       rows,
		HasBorders: true,
	})
	fmt.Fprint(out, tableStr)
	fmt.Fprintf(out, "\n%d row(s) returned.\n", res.Count)
}

func formatCellString(cell any) string {
	if cell == nil {
		return "NULL"
	}

	return fmt.Sprint(cell)
}
