package cmddb

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
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

func parseFindArgs(args []string) (findParams, *appfault.AppError) {
	p := findParams{limit: 100}
	i := 0
	for i < len(args) {
		consumed, isFlag := parseFindFlagToken(args, i, &p)
		if isFlag {
			i += consumed
			continue
		}
		if p.dbArg == "" && !strings.HasPrefix(args[i], "-") {
			p.dbArg = args[i]
		}
		i++
	}

	if p.findText == "" && p.grepPat == "" {
		return p, appfault.NewValidationError("specify --find/-f <text> or --grep/-r <regex> to search. Run 'gitmap db-find --help' for details.")
	}

	return p, nil
}

func parseFindFlagToken(args []string, idx int, p *findParams) (int, bool) {
	arg := args[idx]
	if (arg == "-table" || arg == "--table" || arg == "-t") && idx+1 < len(args) {
		p.tables = splitCommaItems(args[idx+1])
		return 2, true
	}
	if (arg == "--find" || arg == "-f" || arg == "-find") && idx+1 < len(args) {
		p.findText = args[idx+1]
		return 2, true
	}
	if (arg == "--grep" || arg == "-r" || arg == "-grep") && idx+1 < len(args) {
		p.grepPat = args[idx+1]
		return 2, true
	}
	if (arg == "--columns" || arg == "-c" || arg == "-columns") && idx+1 < len(args) {
		p.columns = splitCommaItems(args[idx+1])
		return 2, true
	}
	if (arg == "--limit" || arg == "-l" || arg == "-n") && idx+1 < len(args) {
		lim, _ := strconv.Atoi(args[idx+1])
		if lim > 0 {
			p.limit = lim
		}
		return 2, true
	}
	if arg == "--json" || arg == "-j" {
		p.isJson = true
		return 1, true
	}
	return parseFindFlagInline(arg, p)
}

func parseFindFlagInline(arg string, p *findParams) (int, bool) {
	if strings.HasPrefix(arg, "-table=") || strings.HasPrefix(arg, "--table=") {
		val := strings.TrimPrefix(strings.TrimPrefix(arg, "--table="), "-table=")
		p.tables = splitCommaItems(val)
		return 1, true
	}
	if strings.HasPrefix(arg, "--find=") || strings.HasPrefix(arg, "-f=") {
		p.findText = strings.TrimPrefix(strings.TrimPrefix(arg, "--find="), "-f=")
		return 1, true
	}
	if strings.HasPrefix(arg, "--grep=") || strings.HasPrefix(arg, "-r=") {
		p.grepPat = strings.TrimPrefix(strings.TrimPrefix(arg, "--grep="), "-r=")
		return 1, true
	}
	if strings.HasPrefix(arg, "--columns=") || strings.HasPrefix(arg, "-c=") {
		val := strings.TrimPrefix(strings.TrimPrefix(arg, "--columns="), "-c=")
		p.columns = splitCommaItems(val)
		return 1, true
	}
	return 0, false
}

func splitCommaItems(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
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

	qRes, qErr := executeDynamicQuery(conn, query)
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

func filterRowsByRegex(qRes QueryResult, targetCols []ColumnInfo, re *regexp.Regexp, limit int) [][]any {
	colIndex := make(map[int]bool)
	for i, name := range qRes.Columns {
		for _, tc := range targetCols {
			if strings.EqualFold(name, tc.Name) {
				colIndex[i] = true
			}
		}
	}

	var matched [][]any
	for _, row := range qRes.Rows {
		if rowMatchesRegex(row, colIndex, re) {
			matched = append(matched, row)
		}
		if len(matched) >= limit {
			break
		}
	}
	return matched
}

func rowMatchesRegex(row []any, colIndex map[int]bool, re *regexp.Regexp) bool {
	for i, val := range row {
		if colIndex[i] && re.MatchString(fmt.Sprint(val)) {
			return true
		}
	}
	return false
}

func outputFindResults(out io.Writer, dbPath string, results []TableFindResult, isJson bool) *appfault.AppError {
	if isJson {
		if err := outputJson(out, results); err != nil {
			return apperror.WrapSimple(err, "db.find.json")
		}
		return nil
	}

	if len(results) == 0 {
		fmt.Fprintln(out, "No matching records found.")
		return nil
	}

	for _, r := range results {
		fmt.Fprintf(out, "Table: %s in %s (%d matches)\n\n", r.Table, dbPath, r.Count)
		renderQueryTerminal(out, QueryResult{Columns: r.Columns, Rows: r.Rows, Count: r.Count})
		fmt.Fprintln(out)
	}
	return nil
}
