package heatmap

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
	"gopkg.in/yaml.v3"
)

var allowedPrefixes = map[string]bool{"SELECT": true, "PRAGMA": true, "EXPLAIN": true, "WITH": true}
var blockedKeywords = map[string]bool{
	"INSERT": true, "UPDATE": true, "DELETE": true, "DROP": true, "ALTER": true, "CREATE": true, "REPLACE": true,
	"ATTACH": true, "DETACH": true, "VACUUM": true, "REINDEX": true, "TRANSACTION": true, "COMMIT": true, "ROLLBACK": true,
}

// ExecuteHeatmapQuery validates, resolves, executes, and formats read-only heatmap SQL queries.
func ExecuteHeatmapQuery(isRc bool, flags map[string]string, pos []string, out io.Writer) error {
	rawQuery, err := extractQueryString(flags, pos)
	if err != nil {
		return err
	}
	cleanSql, valErr := ValidateSafeReadOnlyQuery(rawQuery)
	if valErr != nil {
		return valErr
	}
	limitVal := 100
	if l, parseErr := strconv.Atoi(resolveLimitFlag(flags)); parseErr == nil && l > 0 {
		limitVal = l
	}
	execSql := injectDefaultLimit(cleanSql, limitVal)
	dbPath := resolveQueryDbPath(isRc, flags, execSql)
	start := time.Now()
	cols, rows, qErr := store.ExecuteReadOnlyQuery(dbPath, execSql)
	if qErr != nil {
		return qErr
	}
	durationMs := float64(time.Since(start).Microseconds()) / 1000.0
	return renderQueryResult(execSql, cols, rows, durationMs, resolveFormat(flags), out)
}

// RunHeatmapQuery provides compatibility with alternative dispatcher argument signatures.
func RunHeatmapQuery(flags map[string]string, pos []string, isRc bool, out io.Writer) error {
	return ExecuteHeatmapQuery(isRc, flags, pos, out)
}

func extractQueryString(flags map[string]string, pos []string) (string, error) {
	for _, q := range []string{strings.Join(pos, " "), flags["sql"], flags["query"]} {
		if s := strings.TrimSpace(q); s != "" {
			return s, nil
		}
	}
	return "", apperror.NewValidation("heatmap_query", "E1000", "missing SQL query argument")
}

// ValidateSafeReadOnlyQuery ensures the SQL query is strictly read-only, single-statement, and free of mutating tokens.
func ValidateSafeReadOnlyQuery(rawQuery string) (string, error) {
	clean, hasMulti := stripComments(rawQuery)
	if hasMulti {
		return "", apperror.NewValidation("heatmap_query", "E1000", "multi-statement queries are prohibited in read-only execution")
	}
	if clean == "" {
		return "", apperror.NewValidation("heatmap_query", "E1000", "SQL query cannot be empty")
	}
	tokens := tokenizeSql(clean)
	if len(tokens) == 0 || !allowedPrefixes[strings.ToUpper(tokens[0])] {
		first := ""
		if len(tokens) > 0 {
			first = tokens[0]
		}
		return "", apperror.NewValidation("heatmap_query", "E1000", fmt.Sprintf("only read-only queries (SELECT, PRAGMA, EXPLAIN, WITH) are permitted; query starts with %q", first))
	}
	for _, t := range tokens {
		if blockedKeywords[strings.ToUpper(t)] {
			return "", apperror.NewValidation("heatmap_query", "E1000", fmt.Sprintf("mutating keyword %q is forbidden in read-only queries", t))
		}
	}
	return clean, nil
}

func stripComments(raw string) (string, bool) {
	var sb strings.Builder
	runes, n, inQuote, hasSemi := []rune(raw), len(raw), rune(0), false
	for i := 0; i < n; i++ {
		ch := runes[i]
		if inQuote != 0 {
			sb.WriteRune(ch)
			if ch == inQuote {
				inQuote = 0
			}
			continue
		}
		if ch == '\'' || ch == '"' || ch == '`' {
			inQuote = ch
			sb.WriteRune(ch)
			continue
		}
		if ch == '-' && i+1 < n && runes[i+1] == '-' {
			for i < n && runes[i] != '\n' {
				i++
			}
			sb.WriteRune(' ')
			continue
		}
		if ch == '/' && i+1 < n && runes[i+1] == '*' {
			for i+1 < n && !(runes[i] == '*' && runes[i+1] == '/') {
				i++
			}
			i += 2
			sb.WriteRune(' ')
			continue
		}
		if ch == ';' {
			hasSemi = true
			continue
		}
		if hasSemi && !unicode.IsSpace(ch) {
			return "", true
		}
		sb.WriteRune(ch)
	}
	return strings.TrimSpace(sb.String()), false
}

func tokenizeSql(s string) []string {
	var tokens []string
	var cur strings.Builder
	inQuote := rune(0)
	for _, ch := range s {
		if inQuote != 0 {
			if ch == inQuote {
				inQuote = 0
			}
			continue
		}
		if ch == '\'' || ch == '"' || ch == '`' {
			inQuote = ch
			continue
		}
		if unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '_' {
			cur.WriteRune(ch)
			continue
		}
		if cur.Len() > 0 {
			tokens, cur = append(tokens, cur.String()), strings.Builder{}
		}
	}
	if cur.Len() > 0 {
		tokens = append(tokens, cur.String())
	}
	return tokens
}

func injectDefaultLimit(query string, limitVal int) string {
	for _, t := range tokenizeSql(query) {
		if strings.EqualFold(t, "LIMIT") {
			return query
		}
	}
	if limitVal <= 0 {
		limitVal = 100
	}
	return fmt.Sprintf("%s LIMIT %d", query, limitVal)
}

func resolveQueryDbPath(isRc bool, flags map[string]string, query string) string {
	if !isRc {
		db := strings.ToLower(strings.TrimSpace(flags["db"]))
		if db == "root" || db == "local" || (db == "" && flags["repo"] == "") {
			return store.HeatmapRootDbPath()
		}
		if db == "repo" || (db == "" && flags["repo"] != "") {
			return store.ResolveHeatmapRepoSplitDbPath(resolveTargetRepo(flags["repo"]), "")
		}
		return flags["db"]
	}
	aiName, repoSlug := resolveTargetAi(flags["ai"]), resolveTargetRepo(flags["repo"])
	tier, db := strings.ToLower(strings.TrimSpace(flags["tier"])), strings.ToLower(strings.TrimSpace(flags["db"]))
	switch {
	case db == "root" || tier == "root":
		return store.ResolveRcHeatmapRootDbPath(aiName)
	case db == "local" || tier == "local":
		return store.ResolveRcHeatmapLocalTrackingDbPath()
	case db == "repo" || tier == "repo":
		return store.ResolveRcHeatmapRepoDbPath(aiName, repoSlug)
	case db != "":
		return flags["db"]
	}
	upper := strings.ToUpper(query)
	if strings.Contains(upper, "HEATMAPAISESSION") || strings.Contains(upper, "HEATMAPREPOREGISTRY") || strings.Contains(upper, "HEATMAPCOMMITLOG") {
		return store.ResolveRcHeatmapRootDbPath(aiName)
	}
	if strings.Contains(upper, "HEATMAPLOCALTRACKING") {
		return store.ResolveRcHeatmapLocalTrackingDbPath()
	}
	return store.ResolveRcHeatmapRepoDbPath(aiName, repoSlug)
}

func resolveLimitFlag(flags map[string]string) string {
	for _, k := range []string{"limit", "l", "n"} {
		if v := flags[k]; v != "" {
			return v
		}
	}
	return ""
}

func renderQueryResult(query string, cols []string, rows [][]any, durationMs float64, format string, out io.Writer) error {
	switch format {
	case "table":
		return renderQueryTable(cols, rows, durationMs, out)
	case "json":
		return renderQueryJSON(query, cols, rows, durationMs, out)
	case "yaml", "yml":
		return renderQueryYAML(query, cols, rows, durationMs, out)
	default:
		return apperror.NewValidation("heatmap_query", "E1003", fmt.Sprintf("unsupported format %q, must be table, json, or yaml", format))
	}
}

func renderQueryTable(cols []string, rows [][]any, durationMs float64, out io.Writer) error {
	if len(rows) == 0 {
		fmt.Fprintf(out, "No records found matching query.\n(0 rows returned in %.2fms)\n", durationMs)
		return nil
	}
	tableCols := make([]termout.Column, len(cols))
	for i, c := range cols {
		tableCols[i] = termout.Column{Title: strings.ToUpper(c), Align: termout.AlignLeft, MaxWidth: 35}
	}
	tableRows := make([]termout.Row, len(rows))
	for i, r := range rows {
		cells := make([]string, len(r))
		for j, val := range r {
			cells[j] = formatCellValue(val)
		}
		tableRows[i] = termout.Row{Cells: cells}
	}
	fmt.Fprint(out, termout.RenderTable(termout.TableConfig{Columns: tableCols, Rows: tableRows}))
	fmt.Fprintf(out, "(%d rows returned in %.2fms)\n", len(rows), durationMs)
	return nil
}

func renderQueryJSON(query string, cols []string, rows [][]any, durationMs float64, out io.Writer) error {
	payload := map[string]any{"query": query, "count": len(rows), "durationMs": durationMs, "columns": cols, "rows": buildRowMaps(cols, rows)}
	enc := json.NewEncoder(out)
	enc.SetIndent("", "  ")
	return enc.Encode(payload)
}

func renderQueryYAML(query string, cols []string, rows [][]any, durationMs float64, out io.Writer) error {
	payload := map[string]any{"query": query, "count": len(rows), "durationMs": durationMs, "rows": buildRowMaps(cols, rows)}
	data, err := yaml.Marshal(payload)
	if err != nil {
		return apperror.WrapSimple(err, "heatmap_query.renderYAML")
	}
	_, writeErr := out.Write(data)
	return writeErr
}

func buildRowMaps(cols []string, rows [][]any) []map[string]any {
	rowMaps := make([]map[string]any, len(rows))
	for i, r := range rows {
		m := make(map[string]any, len(cols))
		for j, c := range cols {
			m[c] = r[j]
		}
		rowMaps[i] = m
	}
	return rowMaps
}

func formatCellValue(val any) string {
	switch v := val.(type) {
	case nil:
		return "NULL"
	case string:
		return v
	case bool:
		return strconv.FormatBool(v)
	case []byte:
		return string(v)
	default:
		return fmt.Sprintf("%v", v)
	}
}
