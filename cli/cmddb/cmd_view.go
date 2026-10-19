package cmddb

import (
	"database/sql"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/appfault"
)

const dbViewHelp = `gitmap db-view — Query tables with automatic join inference, optional views, and filtering.

Usage:
  gitmap db-view <alias-or-path> -t "<tables>" [flags]

Flags:
  -t, --table, -table <spec>    Table declarations with aliases (e.g. "table1 as t1, table2 as t2")
  --viewname, -v, --name <name> Create a VIEW (or TEMP VIEW) in the database
  --find, -f <condition>        Filter condition (e.g. "t1.field = text and t2.field like 'x*'")
  --limit, -l <N>               Row limit (default 100)
  --json, -j                    Output result as JSON with execution timing
  --help, -h                    Show this help text

Examples:
  gitmap db-view master -t "Repository as r, ScanHistory as s" --find "r.IsActive = 1"
  gitmap db-view tasks -t "Subtask as s, Plan as p" --viewname v_tasks --find "s.Status = 'DONE'"
  gitmap db-view tasks -t "Subtask as s, Plan as p" --find "s.Status = DONE and p.Title like 'v6*'" --json
`

type viewParams struct {
	dbArg     string
	tableSpec string
	viewName  string
	findCond  string
	limit     int
	isJson    bool
}

type tableAlias struct {
	Table string
	Alias string
}

// ViewOutput holds JSON output for a view query including timing.
type ViewOutput struct {
	ViewName  string   `json:"viewName,omitempty"`
	SQL       string   `json:"sql"`
	Columns   []string `json:"columns"`
	Rows      [][]any  `json:"rows"`
	Count     int      `json:"count"`
	ElapsedMs int64    `json:"elapsedMs"`
}

// RunDbView executes table joins, optional view creation, and condition filtering.
func RunDbView(args []string) *appfault.AppError {
	if len(args) == 0 || (len(args) == 1 && isDbHelpToken(args[0])) {
		fmt.Print(dbViewHelp)
		return nil
	}

	p, parseErr := parseViewArgs(args)
	if parseErr != nil {
		return parseErr
	}

	conn, dbPath, openErr := OpenResolvedDb(p.dbArg)
	if openErr != nil {
		return apperror.WrapSimple(openErr, "db.view.open")
	}
	defer conn.Close()

	tables := parseTableAliases(p.tableSpec)
	if len(tables) == 0 {
		return appfault.NewValidationError("no tables specified in -t flag")
	}

	fromClause := buildFromClause(conn, tables)
	query, viewErr := prepareViewQuery(conn, p, fromClause)
	if viewErr != nil {
		return viewErr
	}

	start := time.Now()
	qRes, qErr := executeDynamicQuery(conn, query)
	elapsed := time.Since(start)
	if qErr != nil {
		return apperror.WrapSimple(qErr, "db.view.query")
	}

	return outputViewResults(os.Stdout, dbPath, p, query, qRes, elapsed)
}

func parseViewArgs(args []string) (viewParams, *appfault.AppError) {
	p := viewParams{limit: 100}
	i := 0
	for i < len(args) {
		consumed, isFlag := parseViewFlagToken(args, i, &p)
		if isFlag {
			i += consumed
			continue
		}
		assignViewPositional(args[i], &p)
		i++
	}

	if p.dbArg == "" || p.tableSpec == "" {
		return p, appfault.NewValidationError("usage: gitmap db-view <alias-or-path> -t \"<tables>\" [--viewname <name>] [--find <condition>] [--json]")
	}

	return p, nil
}

func parseViewFlagToken(args []string, idx int, p *viewParams) (int, bool) {
	arg := args[idx]
	if (arg == "-t" || arg == "--table" || arg == "-table") && idx+1 < len(args) {
		p.tableSpec = args[idx+1]
		return 2, true
	}
	if (arg == "--viewname" || arg == "-v" || arg == "-viewname" || arg == "--name") && idx+1 < len(args) {
		p.viewName = args[idx+1]
		return 2, true
	}
	if (arg == "--find" || arg == "-f" || arg == "-find") && idx+1 < len(args) {
		p.findCond = args[idx+1]
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
	return parseViewFlagInline(arg, p)
}

func parseViewFlagInline(arg string, p *viewParams) (int, bool) {
	if strings.HasPrefix(arg, "-t=") || strings.HasPrefix(arg, "--table=") {
		p.tableSpec = strings.TrimPrefix(strings.TrimPrefix(arg, "--table="), "-t=")
		return 1, true
	}
	if strings.HasPrefix(arg, "--viewname=") || strings.HasPrefix(arg, "-v=") {
		p.viewName = strings.TrimPrefix(strings.TrimPrefix(arg, "--viewname="), "-v=")
		return 1, true
	}
	if strings.HasPrefix(arg, "--find=") || strings.HasPrefix(arg, "-f=") {
		p.findCond = strings.TrimPrefix(strings.TrimPrefix(arg, "--find="), "-f=")
		return 1, true
	}
	return 0, false
}

func assignViewPositional(arg string, p *viewParams) {
	if strings.HasPrefix(arg, "-") {
		return
	}
	if p.dbArg == "" {
		p.dbArg = arg
		return
	}
	if p.tableSpec == "" {
		p.tableSpec = arg
	}
}

func parseTableAliases(spec string) []tableAlias {
	parts := strings.Split(spec, ",")
	out := make([]tableAlias, 0, len(parts))
	for _, p := range parts {
		ta := parseSingleTableAlias(strings.TrimSpace(p))
		if ta.Table != "" {
			out = append(out, ta)
		}
	}
	return out
}

func parseSingleTableAlias(token string) tableAlias {
	fields := strings.Fields(token)
	if len(fields) == 3 && strings.EqualFold(fields[1], "as") {
		return tableAlias{Table: fields[0], Alias: fields[2]}
	}
	if len(fields) == 2 {
		return tableAlias{Table: fields[0], Alias: fields[1]}
	}
	if len(fields) == 1 {
		return tableAlias{Table: fields[0], Alias: fields[0]}
	}
	return tableAlias{}
}

func buildFromClause(conn *sql.DB, tables []tableAlias) string {
	if len(tables) == 1 {
		return fmt.Sprintf(`FROM "%s" AS %s`, escapeSqlIdent(tables[0].Table), tables[0].Alias)
	}

	joinCond := inferAutoJoin(conn, tables[0], tables[1])
	if joinCond != "" {
		return fmt.Sprintf(`FROM "%s" AS %s JOIN "%s" AS %s ON %s`,
			escapeSqlIdent(tables[0].Table), tables[0].Alias,
			escapeSqlIdent(tables[1].Table), tables[1].Alias,
			joinCond)
	}

	return fmt.Sprintf(`FROM "%s" AS %s, "%s" AS %s`,
		escapeSqlIdent(tables[0].Table), tables[0].Alias,
		escapeSqlIdent(tables[1].Table), tables[1].Alias)
}

func normalizeFindCondition(cond string) string {
	if cond == "" {
		return ""
	}
	reLike := regexp.MustCompile(`(?i)\blike\s+'([^']*)'`)
	s := reLike.ReplaceAllStringFunc(cond, func(m string) string {
		return strings.ReplaceAll(m, "*", "%")
	})

	reEqBare := regexp.MustCompile(`=\s*([a-zA-Z_][a-zA-Z0-9_]*)\b`)
	return reEqBare.ReplaceAllStringFunc(s, func(m string) string {
		val := strings.TrimSpace(strings.TrimPrefix(m, "="))
		lower := strings.ToLower(val)
		if lower == "true" || lower == "false" || lower == "null" {
			return m
		}
		return fmt.Sprintf("= '%s'", val)
	})
}

func prepareViewQuery(conn *sql.DB, p viewParams, fromClause string) (string, *appfault.AppError) {
	normCond := normalizeFindCondition(p.findCond)
	whereClause := ""
	if normCond != "" {
		whereClause = " WHERE " + normCond
	}

	if p.viewName != "" {
		createViewSql := fmt.Sprintf(`CREATE TEMP VIEW IF NOT EXISTS "%s" AS SELECT * %s`, escapeSqlIdent(p.viewName), fromClause)
		if _, err := conn.Exec(createViewSql); err != nil {
			return "", apperror.WrapSimple(err, "db.view.create")
		}
		return fmt.Sprintf(`SELECT * FROM "%s"%s LIMIT %d`, escapeSqlIdent(p.viewName), whereClause, p.limit), nil
	}

	return fmt.Sprintf(`SELECT * %s%s LIMIT %d`, fromClause, whereClause, p.limit), nil
}

func escapeSqlIdent(s string) string {
	return strings.ReplaceAll(s, `"`, `""`)
}
