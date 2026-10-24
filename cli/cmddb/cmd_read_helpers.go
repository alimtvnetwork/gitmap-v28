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
		return nil, dbPath, apperror.NewNotFoundError(fmt.Sprintf("database file not found: %s", dbPath))
	}

	openRes := store.OpenSQLiteDB(dbPath)
	if openRes.IsFailure() {
		return nil, dbPath, apperror.WrapSimple(openRes.AppError(), "open database '"+dbPath+"'")
	}
	conn := openRes.Value

	return enforceReadOnlyPragma(conn, dbPath)
}

func enforceReadOnlyPragma(conn *sql.DB, dbPath string) (*sql.DB, string, error) {
	initPragma := "PRAGMA query_only = ON; PRAGMA busy_timeout = 5000;"
	if _, pragmaErr := conn.Exec(initPragma); pragmaErr != nil {
		conn.Close()

		return nil, dbPath, apperror.WrapSimple(pragmaErr, "enforce read-only pragma")
	}

	return conn, dbPath, nil
}

func registerDbPathAlias(rawPath, aliasName string) (string, error) {
	cleanName := strings.TrimSpace(aliasName)
	if cleanName == "" {
		return "", apperror.NewValidationError("alias name cannot be empty")
	}

	absPath, err := filepath.Abs(strings.TrimSpace(rawPath))
	if err != nil {
		return "", apperror.WrapSimple(err, "resolve database path")
	}

	return commitDbPathAlias(absPath, cleanName)
}

func commitDbPathAlias(absPath, cleanName string) (string, error) {
	if _, statErr := os.Stat(absPath); statErr != nil {
		return "", apperror.NewNotFoundError(fmt.Sprintf("database file not found: %s", absPath))
	}

	if regErr := RegisterAlias(cleanName, absPath, "Registered via db-read"); regErr != nil {
		return "", apperror.WrapSimple(regErr, "save database alias")
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
		return TableOverview{}, apperror.NewNotFoundError(fmt.Sprintf("table '%s' not found in database %s", tableName, dbPath))
	}

	colCount := countTableColumns(conn, tableName)
	rowCount := countTableRows(conn, tableName)
	previewRes := queryTablePreview(conn, tableName, limit)

	return TableOverview{
		Table:          tableName,
		RowCount:       rowCount,
		ColumnCount:    colCount,
		Schema:         schema,
		PreviewColumns: previewRes.Columns,
		PreviewRows:    previewRes.Rows,
	}, nil
}

func queryTablePreview(conn *sql.DB, tableName string, limit int) QueryResult {
	previewLimit := limit
	if previewLimit <= 0 {
		previewLimit = 5
	}

	escaped := strings.ReplaceAll(tableName, `"`, `""`)
	previewQuery := fmt.Sprintf(`SELECT * FROM "%s" LIMIT %d`, escaped, previewLimit)
	res, _ := executeDynamicQuery(conn, previewQuery)

	return res
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
	columnOverviews, err := collectColumnOverviews(conn, dbPath, tables)
	if err != nil {
		return failRead(os.Stdout, flags.isJson, "db-read:table", dbPath, dbArg, "ERR_TABLE_NOT_FOUND", 22, err.Error(), "Verify table name via 'gitmap db-read ls'.")
	}

	elapsed := time.Since(start)
	if flags.isJson {
		outputReadSuccess(os.Stdout, "db-read:table", dbPath, dbArg, map[string]any{"tables": columnOverviews}, elapsed)

		return nil
	}

	renderColumnsList(os.Stdout, dbPath, columnOverviews, elapsed)

	return nil
}

func collectColumnOverviews(conn *sql.DB, dbPath string, tables []string) ([]TableColumnOverview, error) {
	var overviews []TableColumnOverview
	for _, t := range tables {
		cols, err := fetchColumnInfos(conn, t)
		if err != nil || len(cols) == 0 {
			return nil, apperror.NewNotFoundError(fmt.Sprintf("table '%s' not found or has no columns in %s", t, dbPath))
		}

		overviews = append(overviews, TableColumnOverview{Table: t, Columns: cols})
	}

	return overviews, nil
}

func renderColumnsList(out io.Writer, dbPath string, overviews []TableColumnOverview, elapsed time.Duration) {
	for _, co := range overviews {
		renderColumnsTerminal(out, dbPath, co.Table, co.Columns)
		fmt.Println()
	}

	fmt.Println(formatCompletionTime(elapsed))
}

func handleTableOverviewOutput(flags readFlags, conn *sql.DB, dbPath, dbArg string, tables []string, start time.Time) *appfault.AppError {
	overviews, err := collectTableOverviews(conn, dbPath, tables, flags.limit)
	if err != nil {
		return failRead(os.Stdout, flags.isJson, "db-read:table", dbPath, dbArg, "ERR_TABLE_NOT_FOUND", 22, err.Error(), "Verify table name via 'gitmap db-read ls'.")
	}

	elapsed := time.Since(start)
	if flags.isJson {
		outputReadSuccess(os.Stdout, "db-read:table", dbPath, dbArg, map[string]any{"tables": overviews}, elapsed)

		return nil
	}

	renderTableOverviewsList(os.Stdout, overviews, elapsed)

	return nil
}

func collectTableOverviews(conn *sql.DB, dbPath string, tables []string, limit int) ([]TableOverview, error) {
	var overviews []TableOverview
	for _, t := range tables {
		overview, err := fetchTableOverview(conn, t, dbPath, limit)
		if err != nil {
			return nil, err
		}

		overviews = append(overviews, overview)
	}

	return overviews, nil
}

func renderTableOverviewsList(out io.Writer, overviews []TableOverview, elapsed time.Duration) {
	for _, ov := range overviews {
		renderTableOverviewTerminal(out, ov)
		fmt.Println()
	}

	fmt.Println(formatCompletionTime(elapsed))
}

func handleReadInspect(flags readFlags, dbArg string) *appfault.AppError {
	start := time.Now()
	conn, dbPath, err := openReadOnlyResolvedDb(dbArg)
	if err != nil {
		return failRead(os.Stdout, flags.isJson, "db-read:inspect", dbArg, dbArg, "ERR_DB_NOT_FOUND", 20, err.Error(), "Verify database alias or run 'gitmap scan' to initialize.")
	}
	defer conn.Close()

	tables, colMap, inspectErr := fetchOrCacheMetadata(conn, dbPath)
	elapsed := time.Since(start)
	if inspectErr != nil {
		return failRead(os.Stdout, flags.isJson, "db-read:inspect", dbPath, dbArg, "ERR_EXECUTION", 23, inspectErr.Error(), "Check database schema integrity.")
	}

	return outputInspectResult(flags, dbPath, dbArg, tables, colMap, elapsed)
}

func outputInspectResult(flags readFlags, dbPath, dbArg string, tables []TableInfo, colMap map[string][]ColumnInfo, elapsed time.Duration) *appfault.AppError {
	if flags.isJson {
		data := map[string]any{"tables": tables, "columns": colMap, "totalTables": len(tables)}
		outputReadSuccess(os.Stdout, "db-read:inspect", dbPath, dbArg, data, elapsed)

		return nil
	}

	renderInspectTerminal(os.Stdout, dbPath, tables, colMap, flags.isShowColumns)
	fmt.Println(formatCompletionTime(elapsed))

	return nil
}

func fetchOrCacheMetadata(conn *sql.DB, dbPath string) ([]TableInfo, map[string][]ColumnInfo, error) {
	cache := GetGlobalTableMetadataCache()
	tables, colMap, hasCached := cache.Get(dbPath)
	if hasCached {
		return tables, colMap, nil
	}

	tables, err := listTablesAndViews(conn)
	if err != nil {
		return nil, nil, err
	}

	colMap = populateColumnsMap(conn, tables)
	cache.Set(dbPath, tables, colMap)

	return tables, colMap, nil
}

func populateColumnsMap(conn *sql.DB, tables []TableInfo) map[string][]ColumnInfo {
	colMap := make(map[string][]ColumnInfo, len(tables))
	for _, t := range tables {
		cols, err := fetchColumnInfos(conn, t.Name)
		if err == nil {
			colMap[t.Name] = cols
		}
	}

	return colMap
}

func renderInspectTerminal(out io.Writer, dbPath string, tables []TableInfo, colMap map[string][]ColumnInfo, isShowColumns bool) {
	renderTablesTerminal(out, dbPath, tables)
	if !isShowColumns {
		return
	}

	for _, t := range tables {
		cols := colMap[t.Name]
		if len(cols) > 0 {
			fmt.Println()
			renderColumnsTerminal(out, dbPath, t.Name, cols)
		}
	}
}
