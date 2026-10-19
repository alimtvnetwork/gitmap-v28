package cmddb

import (
	"database/sql"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/appfault"
)

const dbFindHelp = `gitmap db-find — Search across SQLite table columns by wildcard or regex.

Usage:
  gitmap db-find <alias-or-path> -table <t1,t2> --find/-f <text> [flags]
  gitmap db-find <alias-or-path> -table <t1,t2> --grep/-r <regex> [flags]

Flags:
  -table, --table, -t <t1,t2>     Comma-separated table names to search (or '*' for all)
  --find, -f <text>               Substring or wildcard search ('text', 'text*', '*text')
  --grep, -r <regex>              Regular expression search across column values
  --columns, -c <c1,c2>           Restrict search to comma-separated column names
  --limit, -l, -n <N>             Maximum matching rows per table (default 100)
  --json, -j                      Output results in JSON format
  --help, -h                      Show this help text

Examples:
  gitmap db-find master -table Repository --find "gitmap*"
  gitmap db-find tasks -table Subtask,Plan --grep "^E[0-9]+" --columns Title,Code
  gitmap db-find pipeline -table Runs -f "*failed" --limit 25 --json
`

type findParams struct {
	dbArg    string
	tables   []string
	findText string
	grepPat  string
	columns  []string
	limit    int
	isJson   bool
}

// TableFindResult captures search results from a specific database table.
type TableFindResult struct {
	Table   string   `json:"table"`
	Columns []string `json:"columns"`
	Rows    [][]any  `json:"rows"`
	Count   int      `json:"count"`
}

// RunDbFind executes search across SQLite table columns by wildcard or regex.
func RunDbFind(args []string) *appfault.AppError {
	if len(args) == 0 || (len(args) == 1 && isDbHelpToken(args[0])) {
		fmt.Print(dbFindHelp)
		return nil
	}

	p, parseErr := parseFindArgs(args)
	if parseErr != nil {
		return parseErr
	}

	conn, dbPath, openErr := OpenResolvedDb(p.dbArg)
	if openErr != nil {
		return apperror.WrapSimple(openErr, "db.find.open")
	}
	defer conn.Close()

	tables, tblErr := resolveSearchTables(conn, p.tables)
	if tblErr != nil {
		return apperror.WrapSimple(tblErr, "db.find.tables")
	}

	results, searchErr := executeSearchAcrossTables(conn, tables, p)
	if searchErr != nil {
		return searchErr
	}

	return outputFindResults(os.Stdout, dbPath, results, p.isJson)
}

func resolveSearchTables(conn *sql.DB, requested []string) ([]string, error) {
	if len(requested) > 0 && !(len(requested) == 1 && requested[0] == "*") {
		return requested, nil
	}
	rows, err := conn.Query("SELECT name FROM sqlite_master WHERE type = 'table' AND name NOT LIKE 'sqlite_%' ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var name string
		if scanErr := rows.Scan(&name); scanErr == nil {
			tables = append(tables, name)
		}
	}
	return tables, rows.Err()
}

func executeSearchAcrossTables(conn *sql.DB, tables []string, p findParams) ([]TableFindResult, *appfault.AppError) {
	var results []TableFindResult
	for _, table := range tables {
		res, err := searchSingleTable(conn, table, p)
		if err != nil {
			return nil, err
		}
		if res.Count > 0 {
			results = append(results, res)
		}
	}
	return results, nil
}

func searchSingleTable(conn *sql.DB, table string, p findParams) (TableFindResult, *appfault.AppError) {
	cols, err := fetchColumnInfos(conn, table)
	if err != nil || len(cols) == 0 {
		return TableFindResult{}, nil
	}

	targetCols := filterColList(cols, p.columns)
	if len(targetCols) == 0 {
		return TableFindResult{}, nil
	}

	if p.grepPat != "" {
		return searchTableRegex(conn, table, targetCols, p.grepPat, p.limit)
	}
	return searchTableLike(conn, table, targetCols, p.findText, p.limit)
}

func filterColList(cols []ColumnInfo, allowed []string) []ColumnInfo {
	if len(allowed) == 0 {
		return cols
	}
	allowMap := make(map[string]bool, len(allowed))
	for _, a := range allowed {
		allowMap[strings.ToLower(a)] = true
	}
	var out []ColumnInfo
	for _, c := range cols {
		if allowMap[strings.ToLower(c.Name)] {
			out = append(out, c)
		}
	}
	return out
}

func searchTableLike(conn *sql.DB, table string, cols []ColumnInfo, findText string, limit int) (TableFindResult, *appfault.AppError) {
	pattern := toSqlLikePattern(findText)
	clauses := make([]string, len(cols))
	args := make([]any, 0, len(cols)+1)
	for i, c := range cols {
		clauses[i] = fmt.Sprintf(`"%s" LIKE ?`, strings.ReplaceAll(c.Name, `"`, `""`))
		args = append(args, pattern)
	}
	args = append(args, limit)

	query := fmt.Sprintf(`SELECT * FROM "%s" WHERE (%s) LIMIT ?`,
		strings.ReplaceAll(table, `"`, `""`),
		strings.Join(clauses, " OR "),
	)

	qRes, qErr := executeDynamicQuery(conn, query, args...)
	if qErr != nil {
		return TableFindResult{}, apperror.WrapSimple(qErr, "db.find.like")
	}

	return TableFindResult{
		Table:   table,
		Columns: qRes.Columns,
		Rows:    qRes.Rows,
		Count:   qRes.Count,
	}, nil
}

func toSqlLikePattern(raw string) string {
	if strings.Contains(raw, "*") {
		return strings.ReplaceAll(raw, "*", "%")
	}
	return "%" + raw + "%"
}

func searchTableRegex(conn *sql.DB, table string, targetCols []ColumnInfo, pattern string, limit int) (TableFindResult, *appfault.AppError) {
	re, regErr := regexp.Compile(pattern)
	if regErr != nil {
		return TableFindResult{}, appfault.NewValidationError(fmt.Sprintf("invalid regex pattern '%s': %v", pattern, regErr))
	}

	query := fmt.Sprintf(`SELECT * FROM "%s"`, strings.ReplaceAll(table, `"`, `""`))
	qRes, qErr := executeDynamicQuery(conn, query)
	if qErr != nil {
		return TableFindResult{}, apperror.WrapSimple(qErr, "db.find.regex")
	}

	matchedRows := filterRowsByRegex(qRes, targetCols, re, limit)
	return TableFindResult{
		Table:   table,
		Columns: qRes.Columns,
		Rows:    matchedRows,
		Count:   len(matchedRows),
	}, nil
}
