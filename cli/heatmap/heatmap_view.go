package heatmap

import (
	"database/sql"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	_ "modernc.org/sqlite"
)

// TableInfo represents metadata for an introspected database table or view.
type TableInfo struct {
	Name        string `json:"name" yaml:"name"`
	Type        string `json:"type" yaml:"type"`
	ColumnCount int    `json:"columnCount" yaml:"columnCount"`
	RowCount    int    `json:"rowCount" yaml:"rowCount"`
}

// RunHeatmapTables lists tables, views, column counts, and row counts.
func RunHeatmapTables(isRc bool, flags map[string]string, out io.Writer) error {
	dbPath, err := resolveHeatmapDB(isRc, flags, "")
	if err != nil {
		return err
	}
	db, err := openViewDB(dbPath, isRc)
	if err != nil {
		return err
	}
	defer db.Close()
	rows, err := db.Query("SELECT name, type FROM sqlite_master WHERE type IN ('table', 'view') AND name NOT LIKE 'sqlite_%' ORDER BY type, name")
	if err != nil {
		return apperror.WrapSimple(err, "heatmap_view.list_tables")
	}
	var rawTables []struct{ name, typ string }
	for rows.Next() {
		var name, typ string
		if rows.Scan(&name, &typ) == nil {
			rawTables = append(rawTables, struct{ name, typ string }{name: name, typ: typ})
		}
	}
	_ = rows.Close()

	var tables []TableInfo
	for _, t := range rawTables {
		var colCount, rowCount int
		if err := db.QueryRow(fmt.Sprintf("SELECT count(*) FROM pragma_table_info('%s')", t.name)).Scan(&colCount); err != nil {
			colCount = 0
		}
		if err := db.QueryRow(fmt.Sprintf("SELECT count(*) FROM \"%s\"", t.name)).Scan(&rowCount); err != nil {
			rowCount = 0
		}
		tables = append(tables, TableInfo{Name: t.name, Type: t.typ, ColumnCount: colCount, RowCount: rowCount})
	}
	return formatTablesOutput(out, dbPath, tables, resolveFormat(flags))
}

// RunHeatmapView inspects records in a specified table with filters and limits.
func RunHeatmapView(isRc bool, flags map[string]string, pos []string, out io.Writer) error {
	table := flags["table"]
	if len(pos) > 0 && strings.TrimSpace(pos[0]) != "" {
		table = strings.TrimSpace(pos[0])
	}
	if table == "" {
		return apperror.NewValidation("heatmap_view", "E1001", "missing table name argument for view command")
	}
	dbPath, err := resolveHeatmapDB(isRc, flags, table)
	if err != nil {
		return err
	}
	db, err := openViewDB(dbPath, isRc)
	if err != nil {
		return err
	}
	defer db.Close()
	actualTable, err := findActualTable(db, table, dbPath)
	if err != nil {
		return err
	}
	cols, pkCol := inspectColumns(db, actualTable)
	clauses, args := buildWhereClauses(flags, cols, pkCol)
	var total int
	countQ := fmt.Sprintf("SELECT count(*) FROM \"%s\"", actualTable)
	if len(clauses) > 0 {
		countQ += " WHERE " + strings.Join(clauses, " AND ")
	}
	if err := db.QueryRow(countQ, args...).Scan(&total); err != nil {
		total = 0
	}
	limit := 20
	if val, err := strconv.Atoi(flags["limit"]); err == nil && val > 0 {
		limit = min(val, 1000)
	}
	records, err := executeViewQuery(db, actualTable, pkCol, clauses, args, limit)
	if err != nil {
		return err
	}
	return formatViewOutput(out, dbPath, actualTable, total, limit, cols, records, resolveFormat(flags))
}

func resolveHeatmapDB(isRc bool, flags map[string]string, table string) (string, error) {
	if isRc {
		return resolveRcHeatmapDB(flags, table)
	}
	return resolveLocalHeatmapDB(flags)
}

func resolveLocalHeatmapDB(flags map[string]string) (string, error) {
	dbFlag := strings.ToLower(flags["db"])
	if dbFlag == "repo" || (flags["repo"] != "" && dbFlag == "") {
		return store.ResolveHeatmapRepoSplitDbPath(resolveTargetRepo(flags["repo"]), ""), nil
	}
	if dbFlag == "root" || dbFlag == "local" || dbFlag == "" {
		return store.ResolveHeatmapRootDbPath(""), nil
	}
	return "", apperror.NewValidation("heatmap_view", "E1004", fmt.Sprintf("invalid db target '%s': expected repo, root, or local", dbFlag))
}

func resolveRcHeatmapDB(flags map[string]string, table string) (string, error) {
	dbFlag := strings.ToLower(flags["db"])
	ai, repo := resolveTargetAi(flags["ai"]), resolveTargetRepo(flags["repo"])
	isRoot := strings.EqualFold(table, "HeatmapAiSession") || strings.EqualFold(table, "HeatmapRepoRegistry") ||
		strings.EqualFold(table, "HeatmapCommitLog") || strings.EqualFold(table, "HeatmapParentTask")
	switch {
	case dbFlag == "root" || (dbFlag == "" && isRoot):
		return store.ResolveRcHeatmapRootDbPath(ai), nil
	case dbFlag == "local" || (dbFlag == "" && strings.EqualFold(table, "HeatmapLocalTracking")):
		return store.ResolveRcHeatmapLocalTrackingDbPath(), nil
	case dbFlag == "repo" || dbFlag == "":
		return store.ResolveRcHeatmapRepoDbPath(ai, repo), nil
	default:
		return "", apperror.NewValidation("heatmap_view", "E1004", fmt.Sprintf("invalid db target '%s': expected repo, root, or local", dbFlag))
	}
}

func openViewDB(dbPath string, isRc bool) (*sql.DB, error) {
	openRes := store.OpenSQLiteDB(dbPath)
	if openRes.IsFailure() {
		return nil, openRes.AppError()
	}
	conn := openRes.Value
	if !isRc {
		_ = store.EnsureHeatmapTables(conn)
	} else {
		_ = store.EnsureRcHeatmapRepoTables(conn)
		_ = store.EnsureRcHeatmapRootTables(conn)
		_ = store.EnsureRcHeatmapLocalTrackingTables(conn)
	}
	return conn, nil
}

func findActualTable(db *sql.DB, target, dbPath string) (string, error) {
	var actual string
	err := db.QueryRow("SELECT name FROM sqlite_master WHERE type IN ('table', 'view') AND name NOT LIKE 'sqlite_%' AND LOWER(name) = LOWER(?)", target).Scan(&actual)
	if err != nil {
		return "", apperror.NewValidation("heatmap_view", "E1002", fmt.Sprintf("table '%s' not found in database '%s'", target, dbPath))
	}
	return actual, nil
}

func inspectColumns(db *sql.DB, table string) ([]string, string) {
	rows, err := db.Query(fmt.Sprintf("PRAGMA table_info(\"%s\")", table))
	if err != nil {
		return nil, ""
	}
	defer rows.Close()
	var (
		cols             []string
		pkCol            string
		cid, notnull, pk int
		name, typ        string
		dflt             *string
	)
	for rows.Next() {
		scanErr := rows.Scan(&cid, &name, &typ, &notnull, &dflt, &pk)
		if scanErr != nil {
			continue
		}
		cols = append(cols, name)
		if pk > 0 && pkCol == "" {
			pkCol = name
		}
	}
	for _, c := range cols {
		if pkCol == "" && (strings.EqualFold(c, table+"Id") || strings.EqualFold(c, "id")) {
			pkCol = c
		}
	}
	return cols, pkCol
}

func buildWhereClauses(flags map[string]string, cols []string, pkCol string) ([]string, []any) {
	var clauses []string
	var args []any
	colStr := strings.ToLower(strings.Join(cols, ","))
	if task := flags["task"]; task != "" && strings.Contains(colStr, "taskslug") {
		clauses, args = append(clauses, "TaskSlug = ?"), append(args, task)
	}
	if repo := flags["repo"]; repo != "" && strings.Contains(colStr, "repo") {
		col := resolveRepoColumn(colStr)
		clauses, args = append(clauses, col+" = ?"), append(args, repo)
	}
	if id := flags["id"]; id != "" && pkCol != "" {
		clauses, args = append(clauses, fmt.Sprintf("\"%s\" = ?", pkCol)), append(args, id)
	}
	return clauses, args
}

func resolveRepoColumn(colStr string) string {
	if strings.Contains(colStr, "reponame") {
		return "RepoName"
	}
	return "RepoSlug"
}

func executeViewQuery(db *sql.DB, table, pk string, clauses []string, args []any, limit int) ([]map[string]any, error) {
	q := fmt.Sprintf("SELECT * FROM \"%s\"", table)
	if len(clauses) > 0 {
		q += " WHERE " + strings.Join(clauses, " AND ")
	}
	if pk != "" {
		q += fmt.Sprintf(" ORDER BY \"%s\" DESC", pk)
	}
	rows, err := db.Query(q+" LIMIT ?", append(args, limit)...)
	if err != nil {
		return nil, apperror.WrapSimple(err, "heatmap_view.execute_query")
	}
	defer rows.Close()

	cols, _ := rows.Columns()
	var records []map[string]any
	for rows.Next() {
		vals, ptrs := make([]any, len(cols)), make([]any, len(cols))
		for i := range vals {
			ptrs[i] = &vals[i]
		}
		scanErr := rows.Scan(ptrs...)
		if scanErr != nil {
			continue
		}
		records = append(records, scanRowValues(cols, vals))
	}
	return records, nil
}

func scanRowValues(cols []string, vals []any) map[string]any {
	rec := make(map[string]any, len(cols))
	for i, name := range cols {
		b, ok := vals[i].([]byte)
		if ok {
			rec[name] = string(b)
			continue
		}
		rec[name] = vals[i]
	}
	return rec
}
