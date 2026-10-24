package cmddb

import (
	"database/sql"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/appfault"
	"github.com/alimtvnetwork/gitmap-v28/cli/lazyregex"
	"github.com/alimtvnetwork/gitmap-v28/cli/strutil"
)

var (
	findLikeRe  = lazyregex.New(`(?i)\blike\s+'([^']*)'`)
	findEqAndRe = lazyregex.New(`=\s*([a-zA-Z0-9_]+(?:\s+[a-zA-Z0-9_]+)*)\s+((?i:and|or))\b`)
	findEqEndRe = lazyregex.New(`=\s*([a-zA-Z0-9_]+(?:\s+[a-zA-Z0-9_]+)*)\s*$`)
)

func inferAutoJoin(conn *sql.DB, t1, t2 tableAlias) string {
	if join, ok := checkForeignKeyJoin(conn, t1, t2); ok {
		return join
	}
	if join, ok := checkForeignKeyJoin(conn, t2, t1); ok {
		return join
	}
	if join, ok := checkIdConventionJoin(conn, t1, t2); ok {
		return join
	}
	if join, ok := checkCommonColumnJoin(conn, t1, t2); ok {
		return join
	}
	return ""
}

func checkForeignKeyJoin(conn *sql.DB, src, dest tableAlias) (string, bool) {
	query := fmt.Sprintf(`PRAGMA foreign_key_list("%s")`, escapeSqlIdent(src.Table))
	rows, err := conn.Query(query)
	if err != nil {
		return "", false
	}
	defer rows.Close()

	for rows.Next() {
		var id, seq int
		var foreignTable, fromCol, toCol, onUpdate, onDelete, match string
		if scanErr := rows.Scan(&id, &seq, &foreignTable, &fromCol, &toCol, &onUpdate, &onDelete, &match); scanErr == nil {
			if strings.EqualFold(foreignTable, dest.Table) {
				return fmt.Sprintf("%s.%s = %s.%s", src.Alias, fromCol, dest.Alias, toCol), true
			}
		}
	}
	return "", false
}

func checkIdConventionJoin(conn *sql.DB, t1, t2 tableAlias) (string, bool) {
	cols1, _ := fetchColumnInfos(conn, t1.Table)
	cols2, _ := fetchColumnInfos(conn, t2.Table)
	pk2 := resolveTablePkName(cols2)
	pk1 := resolveTablePkName(cols1)

	for _, c := range cols1 {
		if strings.EqualFold(c.Name, t2.Table+"Id") || strings.EqualFold(c.Name, t2.Table+"_id") {
			return fmt.Sprintf("%s.%s = %s.%s", t1.Alias, c.Name, t2.Alias, pk2), true
		}
	}

	for _, c := range cols2 {
		if strings.EqualFold(c.Name, t1.Table+"Id") || strings.EqualFold(c.Name, t1.Table+"_id") {
			return fmt.Sprintf("%s.%s = %s.%s", t2.Alias, c.Name, t1.Alias, pk1), true
		}
	}

	return "", false
}

func checkCommonColumnJoin(conn *sql.DB, t1, t2 tableAlias) (string, bool) {
	cols1, _ := fetchColumnInfos(conn, t1.Table)
	cols2, _ := fetchColumnInfos(conn, t2.Table)
	for _, c1 := range cols1 {
		for _, c2 := range cols2 {
			if strings.EqualFold(c1.Name, c2.Name) && (strings.HasSuffix(strings.ToLower(c1.Name), "id") || c1.IsPrimaryKey) {
				return fmt.Sprintf("%s.%s = %s.%s", t1.Alias, c1.Name, t2.Alias, c2.Name), true
			}
		}
	}
	return "", false
}

func outputViewResults(out io.Writer, dbPath string, p viewParams, query string, qRes QueryResult, elapsed time.Duration) *appfault.AppError {
	if p.isJson {
		viewOut := ViewOutput{
			ViewName:  p.viewName,
			SQL:       query,
			Columns:   qRes.Columns,
			Rows:      qRes.Rows,
			Count:     qRes.Count,
			ElapsedMs: elapsed.Milliseconds(),
		}
		if err := outputJson(out, viewOut); err != nil {
			return apperror.WrapSimple(err, "db.view.json")
		}
		return nil
	}

	if p.viewName != "" {
		fmt.Fprintf(out, "Created temporary view '%s' in %s\n\n", p.viewName, dbPath)
	}
	renderQueryTerminal(out, qRes)
	fmt.Fprintf(out, "\n%d row(s) returned in %v.\n", qRes.Count, elapsed)
	return nil
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
	s := findLikeRe.Compiled().ReplaceAllStringFunc(cond, func(m string) string {
		return strings.ReplaceAll(m, "*", "%")
	})

	s = findEqAndRe.Compiled().ReplaceAllStringFunc(s, func(m string) string {
		parts := findEqAndRe.FindStringSubmatch(m)
		if len(parts) >= 3 {
			val := strings.TrimSpace(parts[1])
			op := parts[2]
			return formatNormalizedEqVal(val) + " " + op
		}
		return m
	})

	s = findEqEndRe.Compiled().ReplaceAllStringFunc(s, func(m string) string {
		parts := findEqEndRe.FindStringSubmatch(m)
		if len(parts) >= 2 {
			val := strings.TrimSpace(parts[1])
			return formatNormalizedEqVal(val)
		}
		return m
	})

	return s
}

func formatNormalizedEqVal(val string) string {
	if strutil.EqualFoldAny(val, "true", "false", "null") || strings.Contains(val, ".") {
		return "= " + val
	}
	if _, err := strconv.Atoi(val); err == nil {
		return "= " + val
	}
	return fmt.Sprintf("= '%s'", val)
}

func prepareViewQuery(conn *sql.DB, p viewParams, fromClause string, tables []tableAlias) (string, *appfault.AppError) {
	normCond := normalizeFindCondition(p.findCond)
	whereClause := ""
	if normCond != "" {
		whereClause = " WHERE " + normCond
	}

	if p.viewName != "" {
		createViewSql := fmt.Sprintf(`CREATE VIEW IF NOT EXISTS "%s" AS SELECT * %s`, escapeSqlIdent(p.viewName), fromClause)
		if _, err := conn.Exec(createViewSql); err != nil {
			// Fallback to TEMP VIEW if database is in read-only mode
			tempSql := fmt.Sprintf(`CREATE TEMP VIEW IF NOT EXISTS "%s" AS SELECT * %s`, escapeSqlIdent(p.viewName), fromClause)
			if _, tempErr := conn.Exec(tempSql); tempErr != nil {
				return "", apperror.WrapSimple(err, "db.view.create")
			}
		}

		viewCond := stripTableAliases(normCond, tables)
		viewWhere := ""
		if viewCond != "" {
			viewWhere = " WHERE " + viewCond
		}
		return fmt.Sprintf(`SELECT * FROM "%s"%s LIMIT %d`, escapeSqlIdent(p.viewName), viewWhere, p.limit), nil
	}

	return fmt.Sprintf(`SELECT * %s%s LIMIT %d`, fromClause, whereClause, p.limit), nil
}

func stripTableAliases(cond string, tables []tableAlias) string {
	s := cond
	for _, ta := range tables {
		if ta.Alias != "" {
			re := regexp.MustCompile(`\b` + regexp.QuoteMeta(ta.Alias) + `\.`)
			s = re.ReplaceAllString(s, "")
		}
	}
	return s
}

func escapeSqlIdent(s string) string {
	return strings.ReplaceAll(s, `"`, `""`)
}
