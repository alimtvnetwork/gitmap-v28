package cmddb

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/appfault"
)

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
