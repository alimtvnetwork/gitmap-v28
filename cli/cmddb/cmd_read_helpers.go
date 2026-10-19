package cmddb

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	apperror "github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	appfault "github.com/alimtvnetwork/gitmap-v28/cli/appfault"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// TableOverview holds schema, statistics, and preview rows for a table.
type TableOverview struct {
	Table          string   `json:"table"`
	RowCount       int64    `json:"rowCount"`
	ColumnCount    int      `json:"columnCount"`
	Schema         string   `json:"schema"`
	PreviewColumns []string `json:"previewColumns"`
	PreviewRows    [][]any  `json:"previewRows"`
}

// TableColumnOverview holds column metadata for a single table.
type TableColumnOverview struct {
	Table   string       `json:"table"`
	Columns []ColumnInfo `json:"columns"`
}

func openReadOnlyResolvedDb(dbArg string) (*sql.DB, string, error) {
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

	initPragma := "PRAGMA query_only = ON; PRAGMA busy_timeout = 5000;"
	if _, pragmaErr := conn.Exec(initPragma); pragmaErr != nil {
		conn.Close()
		return nil, dbPath, fmt.Errorf("enforce read-only pragma: %w", pragmaErr)
	}

	return conn, dbPath, nil
}

func registerDbPathAlias(rawPath, aliasName string) (string, error) {
	cleanName := strings.TrimSpace(aliasName)
	if cleanName == "" {
		return "", fmt.Errorf("alias name cannot be empty")
	}

	absPath, err := filepath.Abs(strings.TrimSpace(rawPath))
	if err != nil {
		return "", fmt.Errorf("resolve database path: %w", err)
	}

	if _, statErr := os.Stat(absPath); statErr != nil {
		return "", fmt.Errorf("database file not found: %s", absPath)
	}

	if regErr := RegisterAlias(cleanName, absPath, "Registered via db-read"); regErr != nil {
		return "", fmt.Errorf("save database alias: %w", regErr)
	}

	return absPath, nil
}

func applyQueryLimit(query string, limit int) string {
	clean := strings.TrimSpace(query)
	upper := strings.ToUpper(clean)
	if strings.Contains(upper, " LIMIT ") || strings.HasPrefix(upper, "PRAGMA") || strings.HasPrefix(upper, "EXPLAIN") {
		return clean
	}

	effLimit := limit
	if effLimit <= 0 {
		effLimit = 100
	}

	return fmt.Sprintf("%s LIMIT %d", clean, effLimit)
}

func fetchSingleTableDdl(conn *sql.DB, tableName string) string {
	query := `SELECT sql FROM sqlite_master WHERE (name = ? OR name = ? COLLATE NOCASE) AND sql IS NOT NULL LIMIT 1`
	var ddl string
	if err := conn.QueryRow(query, tableName, tableName).Scan(&ddl); err != nil {
		return ""
	}

	return ensureSemicolon(ddl)
}

func fetchTableOverview(conn *sql.DB, tableName, dbPath string, limit int) (TableOverview, error) {
	schema := fetchSingleTableDdl(conn, tableName)
	if schema == "" {
		return TableOverview{}, fmt.Errorf("table '%s' not found in database %s", tableName, dbPath)
	}

	colCount := countTableColumns(conn, tableName)
	rowCount := countTableRows(conn, tableName)

	previewLimit := limit
	if previewLimit <= 0 {
		previewLimit = 5
	}

	escaped := strings.ReplaceAll(tableName, `"`, `""`)
	previewQuery := fmt.Sprintf(`SELECT * FROM "%s" LIMIT %d`, escaped, previewLimit)
	previewRes, _ := executeDynamicQuery(conn, previewQuery)

	return TableOverview{
		Table:          tableName,
		RowCount:       rowCount,
		ColumnCount:    colCount,
		Schema:         schema,
		PreviewColumns: previewRes.Columns,
		PreviewRows:    previewRes.Rows,
	}, nil
}

func formatCompletionTime(d time.Duration) string {
	ms := float64(d.Microseconds()) / 1000.0
	return fmt.Sprintf("Completed in %.2fms", ms)
}

func outputReadSuccess(out io.Writer, cmd, dbPath, alias string, data map[string]any, elapsed time.Duration) {
	data["executionTimeMs"] = float64(elapsed.Microseconds()) / 1000.0
	envelope := map[string]any{
		"success":   true,
		"command":   cmd,
		"database":  dbPath,
		"alias":     alias,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"data":      data,
		"error":     nil,
	}

	_ = outputJson(out, envelope)
}

func failRead(out io.Writer, isJson bool, cmd, dbPath, alias, code string, exitCode int, msg, remediation string) *appfault.AppError {
	if isJson {
		envelope := map[string]any{
			"success":   false,
			"command":   cmd,
			"database":  dbPath,
			"alias":     alias,
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"data":      nil,
			"error": map[string]any{
				"code":        code,
				"exitCode":    exitCode,
				"message":     msg,
				"remediation": remediation,
			},
		}
		_ = outputJson(out, envelope)
	}

	return apperror.NewWithDetails(cmd, code, msg, "cmddb", apperror.ErrorTypeExecution, apperror.SeverityError, nil)
}

func renderTableOverviewTerminal(out io.Writer, overview TableOverview) {
	fmt.Fprintf(out, "Table: %s (%d columns, %d rows)\n\n", overview.Table, overview.ColumnCount, overview.RowCount)
	fmt.Fprintln(out, "Schema:")
	fmt.Fprintf(out, "  %s\n\n", overview.Schema)

	if len(overview.PreviewColumns) == 0 || len(overview.PreviewRows) == 0 {
		fmt.Fprintln(out, "Preview: (0 rows)")
		return
	}

	fmt.Fprintln(out, "Preview:")
	res := QueryResult{
		Columns: overview.PreviewColumns,
		Rows:    overview.PreviewRows,
		Count:   len(overview.PreviewRows),
	}
	renderQueryTerminal(out, res)
}

func handleTableColumnsOutput(flags readFlags, conn *sql.DB, dbPath, dbArg string, tables []string, start time.Time) *appfault.AppError {
	var columnOverviews []TableColumnOverview
	for _, t := range tables {
		cols, err := fetchColumnInfos(conn, t)
		if err != nil || len(cols) == 0 {
			msg := fmt.Sprintf("table '%s' not found or has no columns in %s", t, dbPath)
			return failRead(os.Stdout, flags.isJson, "db-read:table", dbPath, dbArg, "ERR_TABLE_NOT_FOUND", 22, msg, "Verify table name via 'gitmap db-read ls'.")
		}
		columnOverviews = append(columnOverviews, TableColumnOverview{Table: t, Columns: cols})
	}

	elapsed := time.Since(start)
	if flags.isJson {
		outputReadSuccess(os.Stdout, "db-read:table", dbPath, dbArg, map[string]any{"tables": columnOverviews}, elapsed)
		return nil
	}

	for _, co := range columnOverviews {
		renderColumnsTerminal(os.Stdout, dbPath, co.Table, co.Columns)
		fmt.Println()
	}
	fmt.Println(formatCompletionTime(elapsed))
	return nil
}

func handleTableOverviewOutput(flags readFlags, conn *sql.DB, dbPath, dbArg string, tables []string, start time.Time) *appfault.AppError {
	var overviews []TableOverview
	for _, t := range tables {
		overview, err := fetchTableOverview(conn, t, dbPath, flags.limit)
		if err != nil {
			return failRead(os.Stdout, flags.isJson, "db-read:table", dbPath, dbArg, "ERR_TABLE_NOT_FOUND", 22, err.Error(), "Verify table name via 'gitmap db-read ls'.")
		}
		overviews = append(overviews, overview)
	}

	elapsed := time.Since(start)
	if flags.isJson {
		outputReadSuccess(os.Stdout, "db-read:table", dbPath, dbArg, map[string]any{"tables": overviews}, elapsed)
		return nil
	}

	for _, ov := range overviews {
		renderTableOverviewTerminal(os.Stdout, ov)
		fmt.Println()
	}
	fmt.Println(formatCompletionTime(elapsed))
	return nil
}
