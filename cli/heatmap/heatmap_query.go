// Package heatmap — heatmap_query.go executes read-only SQL queries against heatmap SQLite databases.
package heatmap

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	_ "modernc.org/sqlite"
)

var allowedPrefixes = map[string]bool{
	"SELECT":  true,
	"PRAGMA":  true,
	"EXPLAIN": true,
	"WITH":    true,
}

var blockedKeywords = map[string]bool{
	"INSERT":   true,
	"UPDATE":   true,
	"DELETE":   true,
	"DROP":     true,
	"ALTER":    true,
	"CREATE":   true,
	"TRUNCATE": true,
	"REPLACE":  true,
	"ATTACH":   true,
	"DETACH":   true,
	"VACUUM":   true,
}

func ExecuteHeatmapQuery(isRc bool, flags map[string]string, pos []string, out io.Writer) error {
	rawQuery, err := extractQueryString(flags, pos)
	if err != nil {
		return err
	}
	safeQuery, err := ValidateSafeReadOnlyQuery(rawQuery)
	if err != nil {
		return err
	}
	dbPath := resolveQueryDbPath(isRc, flags, safeQuery)
	limitVal, _ := strconv.Atoi(resolveLimitFlag(flags))
	finalQuery := injectDefaultLimit(safeQuery, limitVal)
	start := time.Now()
	cols, rows, err := store.ExecuteReadOnlyQuery(dbPath, finalQuery)
	durationMs := float64(time.Since(start).Microseconds()) / 1000.0
	if err != nil {
		return apperror.WrapSimple(err, "heatmap_query.exec")
	}
	return renderQueryResult(dbPath, finalQuery, cols, rows, durationMs, resolveFormat(flags), out)
}

func RunHeatmapQuery(flags map[string]string, pos []string, isRc bool, out io.Writer) error {
	return ExecuteHeatmapQuery(isRc, flags, pos, out)
}

func extractQueryString(flags map[string]string, pos []string) (string, error) {
	if q, ok := flags["sql"]; ok && strings.TrimSpace(q) != "" {
		return strings.TrimSpace(q), nil
	}
	if len(pos) > 0 && strings.TrimSpace(pos[0]) != "" {
		return strings.TrimSpace(strings.Join(pos, " ")), nil
	}
	return "", apperror.NewValidation("heatmap_query", "E1000", "missing SQL query: provide via --sql=\"...\" or as positional argument")
}

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
		first := firstTokenName(tokens)
		return "", apperror.NewValidation("heatmap_query", "E1000", fmt.Sprintf("only read-only queries (SELECT, PRAGMA, EXPLAIN, WITH) are permitted; query starts with %q", first))
	}
	for _, t := range tokens {
		if blockedKeywords[strings.ToUpper(t)] {
			return "", apperror.NewValidation("heatmap_query", "E1000", fmt.Sprintf("mutating keyword %q is forbidden in read-only queries", t))
		}
	}
	return clean, nil
}

func firstTokenName(tokens []string) string {
	if len(tokens) == 0 {
		return ""
	}
	return tokens[0]
}

func stripComments(raw string) (string, bool) {
	var sb strings.Builder
	runes, n, inQuote, hasSemi := []rune(raw), len(raw), rune(0), false
	for i := 0; i < n; i++ {
		ch := runes[i]
		if inQuote != 0 {
			sb.WriteRune(ch)
			inQuote = updateQuoteState(inQuote, ch)
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

func updateQuoteState(current, ch rune) rune {
	if current == ch {
		return 0
	}
	return current
}

func tokenizeSql(s string) []string {
	var tokens []string
	var cur strings.Builder
	inQuote := rune(0)
	for _, ch := range s {
		if inQuote != 0 {
			inQuote = updateQuoteState(inQuote, ch)
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
	if isRc {
		return resolveRcQueryDbPath(flags)
	}
	return resolveLocalQueryDbPath(flags)
}

func resolveLocalQueryDbPath(flags map[string]string) string {
	db := strings.ToLower(strings.TrimSpace(flags["db"]))
	if db == "root" || db == "local" || (db == "" && flags["repo"] == "") {
		return store.HeatmapRootDbPath()
	}
	if db == "repo" || (db == "" && flags["repo"] != "") {
		return store.ResolveHeatmapRepoSplitDbPath(resolveTargetRepo(flags["repo"]), "")
	}
	return flags["db"]
}

func resolveRcQueryDbPath(flags map[string]string) string {
	aiName, repoSlug := resolveTargetAi(flags["ai"]), resolveTargetRepo(flags["repo"])
	tier, db := strings.ToLower(strings.TrimSpace(flags["tier"])), strings.ToLower(strings.TrimSpace(flags["db"]))
	if tier == "local" || db == "local" {
		return store.ResolveRcHeatmapLocalTrackingDbPath()
	}
	if tier == "root" || db == "root" {
		return store.ResolveRcHeatmapRootDbPath(aiName)
	}
	return store.ResolveRcHeatmapRepoDbPath(aiName, repoSlug)
}

func resolveLimitFlag(flags map[string]string) string {
	if lim, ok := flags["limit"]; ok && lim != "" {
		return lim
	}
	return "100"
}
