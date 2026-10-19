package cmddb

import (
	"encoding/csv"
	"strings"
)

func parseTokensFromWords(words []string) []string {
	if len(words) == 0 {
		return nil
	}
	raw := strings.Join(words, " ")
	if strings.Contains(raw, ",") {
		return parseCsvTokens(raw)
	}
	out := make([]string, 0, len(words))
	for _, w := range words {
		cleaned := strings.Trim(strings.TrimSpace(w), `,"'`)
		if cleaned != "" {
			out = append(out, cleaned)
		}
	}
	return out
}

func parseCsvTokens(raw string) []string {
	reader := csv.NewReader(strings.NewReader(raw))
	reader.TrimLeadingSpace = true
	records, err := reader.ReadAll()
	if err == nil && len(records) > 0 {
		return records[0]
	}
	return splitRawCsvFallback(raw)
}

func splitRawCsvFallback(raw string) []string {
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		out = append(out, strings.Trim(strings.TrimSpace(p), `"'`))
	}
	return out
}

func checkHasKeyVal(tokens []string) bool {
	for _, t := range tokens {
		if strings.Contains(t, "=") {
			return true
		}
	}
	return false
}

func splitKeyValueToken(t string) (string, string, bool) {
	parts := strings.SplitN(t, "=", 2)
	if len(parts) != 2 {
		return "", "", false
	}
	colName := strings.TrimSpace(parts[0])
	val := strings.Trim(strings.TrimSpace(parts[1]), `"'`)
	return colName, val, colName != ""
}

func indexColumns(cols []ColumnInfo) map[string]bool {
	m := make(map[string]bool, len(cols))
	for _, c := range cols {
		m[strings.ToLower(c.Name)] = true
	}
	return m
}

func findPrimaryKeyCol(cols []ColumnInfo) *ColumnInfo {
	for _, c := range cols {
		if c.IsPrimaryKey {
			return &c
		}
	}
	return nil
}

func filterNonPkCols(cols []ColumnInfo, pkName string) []ColumnInfo {
	var out []ColumnInfo
	for _, c := range cols {
		if !strings.EqualFold(c.Name, pkName) {
			out = append(out, c)
		}
	}
	return out
}
