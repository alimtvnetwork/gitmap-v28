package cmd

import (
	"os"
	"path/filepath"
	"strings"
)

func parseAgyExceptTokens(raw string) []string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil
	}
	if isFilePathToken(trimmed) {
		if content, err := os.ReadFile(trimmed); err == nil {
			return splitTokens(string(content))
		}
	}
	return splitTokens(trimmed)
}

func isFilePathToken(s string) bool {
	lower := strings.ToLower(s)
	if strings.HasSuffix(lower, ".csv") || strings.HasSuffix(lower, ".txt") {
		return true
	}
	if info, err := os.Stat(s); err == nil && !info.IsDir() {
		return true
	}
	return false
}

func splitTokens(content string) []string {
	content = strings.ReplaceAll(content, "\r\n", "\n")
	content = strings.ReplaceAll(content, "\n", ",")
	rawItems := strings.Split(content, ",")
	var out []string
	for _, it := range rawItems {
		c := strings.TrimSpace(it)
		c = strings.Trim(c, `"'`)
		if c != "" {
			out = append(out, strings.ToLower(c))
		}
	}
	return out
}

func isAgyProjectExceptedWithTokens(p AgyProject, tokens []string) bool {
	if len(tokens) == 0 {
		return false
	}
	pID := strings.ToLower(p.ID)
	pName := strings.ToLower(p.Name)
	pPath := strings.ToLower(filepath.Clean(p.GetPath()))
	pBase := strings.ToLower(filepath.Base(p.GetPath()))

	for _, t := range tokens {
		if t == pID || t == pName || t == pBase {
			return true
		}
		if strings.HasPrefix(pID, t) || strings.HasPrefix(pName, t) || strings.HasPrefix(pBase, t) {
			return true
		}
		if pPath != "" && (t == pPath || strings.Contains(pPath, t)) {
			return true
		}
	}
	return false
}
