package cmdmerge

import (
	"database/sql"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func scanForColumnName(rows *sql.Rows, targetCol string) bool {
	for rows.Next() {
		var cid, notnull, pk int
		var name, ctype string
		var dflt any
		scanErr := rows.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk)
		isMatch := scanErr == nil && strings.EqualFold(name, targetCol)
		if isMatch {
			return true
		}
	}

	return false
}

func expandFileGlobs(patterns []string) []string {
	var collected []string
	for _, p := range patterns {
		matches, err := filepath.Glob(p)
		if err == nil && len(matches) > 0 {
			collected = append(collected, matches...)
			continue
		}
		collected = append(collected, p)
	}
	return dedupeAndFilterJSONFiles(collected)
}

func scanDirectoryJSONFiles(dir string) ([]string, error) {
	var matches []string
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if strings.EqualFold(filepath.Ext(e.Name()), ".json") {
			matches = append(matches, filepath.Join(dir, e.Name()))
		}
	}
	sort.Strings(matches)
	return matches, nil
}

func dedupeAndFilterJSONFiles(files []string) []string {
	seen := make(map[string]bool)
	var out []string
	for _, f := range files {
		clean := filepath.Clean(f)
		if seen[clean] {
			continue
		}
		seen[clean] = true
		out = append(out, clean)
	}
	return out
}
