package cmddb

import (
	"fmt"
	"io"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/appfault"
)

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
		} else if len(p.tables) == 0 && !strings.HasPrefix(args[i], "-") {
			p.tables = splitCommaItems(args[i])
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
	if isTableFlag(arg) && idx+1 < len(args) {
		p.tables = append(p.tables, splitCommaItems(args[idx+1])...)
		consumed := 2
		for idx+consumed < len(args) && !strings.HasPrefix(args[idx+consumed], "-") {
			p.tables = append(p.tables, splitCommaItems(args[idx+consumed])...)
			consumed++
		}
		return consumed, true
	}
	if isFindFlag(arg) && idx+1 < len(args) {
		p.findText = strings.Trim(args[idx+1], `"'`)
		return 2, true
	}
	if isGrepFlag(arg) && idx+1 < len(args) {
		p.grepPat = strings.Trim(args[idx+1], `"'`)
		return 2, true
	}
	if isColumnsFlag(arg) && idx+1 < len(args) {
		p.columns = append(p.columns, splitCommaItems(args[idx+1])...)
		consumed := 2
		for idx+consumed < len(args) && !strings.HasPrefix(args[idx+consumed], "-") {
			p.columns = append(p.columns, splitCommaItems(args[idx+consumed])...)
			consumed++
		}
		return consumed, true
	}
	if isLimitFlag(arg) && idx+1 < len(args) {
		lim, _ := strconv.Atoi(args[idx+1])
		if lim > 0 {
			p.limit = lim
		}
		return 2, true
	}
	if slices.Contains([]string{"--json", "-j"}, arg) {
		p.isJson = true
		return 1, true
	}
	return parseFindFlagInline(arg, p)
}

func isTableFlag(arg string) bool {
	return slices.Contains([]string{"-table", "--table", "-t"}, arg)
}

func isFindFlag(arg string) bool {
	return slices.Contains([]string{"--find", "-f", "-find", "--find/f", "-f/--find"}, arg)
}

func isGrepFlag(arg string) bool {
	return slices.Contains([]string{"--grep", "-r", "-grep", "--grep/r", "-r/--grep"}, arg)
}

func isColumnsFlag(arg string) bool {
	return slices.Contains([]string{"--columns", "-c", "-columns", "--columns/c", "-c/--columns"}, arg)
}

func isLimitFlag(arg string) bool {
	return slices.Contains([]string{"--limit", "-l", "-n", "--limit/l", "-l/--limit"}, arg)
}

func parseFindFlagInline(arg string, p *findParams) (int, bool) {
	if strings.HasPrefix(arg, "-table") && len(arg) > len("-table") {
		val := strings.TrimPrefix(arg, "-table")
		val = strings.TrimPrefix(val, "=")
		p.tables = append(p.tables, splitCommaItems(strings.Trim(val, `"'`))...)
		return 1, true
	}
	if strings.HasPrefix(arg, "--table") && len(arg) > len("--table") {
		val := strings.TrimPrefix(arg, "--table")
		val = strings.TrimPrefix(val, "=")
		p.tables = append(p.tables, splitCommaItems(strings.Trim(val, `"'`))...)
		return 1, true
	}
	if strings.HasPrefix(arg, "--find=") || strings.HasPrefix(arg, "-f=") {
		p.findText = strings.Trim(strings.TrimPrefix(strings.TrimPrefix(arg, "--find="), "-f="), `"'`)
		return 1, true
	}
	if strings.HasPrefix(arg, "--grep=") || strings.HasPrefix(arg, "-r=") {
		p.grepPat = strings.Trim(strings.TrimPrefix(strings.TrimPrefix(arg, "--grep="), "-r="), `"'`)
		return 1, true
	}
	if strings.HasPrefix(arg, "--columns=") || strings.HasPrefix(arg, "-c=") {
		val := strings.TrimPrefix(strings.TrimPrefix(arg, "--columns="), "-c=")
		p.columns = append(p.columns, splitCommaItems(strings.Trim(val, `"'`))...)
		return 1, true
	}
	if strings.HasPrefix(arg, "--limit=") || strings.HasPrefix(arg, "-l=") {
		lim, _ := strconv.Atoi(strings.TrimPrefix(strings.TrimPrefix(arg, "--limit="), "-l="))
		if lim > 0 {
			p.limit = lim
		}
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
