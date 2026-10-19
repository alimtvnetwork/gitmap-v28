package cmddb

import (
	"database/sql"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/appfault"
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
