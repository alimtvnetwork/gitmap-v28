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
	if tokens, ok := tryReadTokenFile(trimmed); ok {
		return tokens
	}
	return splitTokens(trimmed)
}

func tryReadTokenFile(path string) ([]string, bool) {
	if !isFilePathToken(path) {
		return nil, false
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, false
	}
	return splitTokens(string(content)), true
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
	pPath := strings.ToLower(filepath.ToSlash(filepath.Clean(p.GetPath())))
	pBase := strings.ToLower(filepath.Base(p.GetPath()))

	for _, t := range tokens {
		tNorm := strings.ToLower(filepath.ToSlash(filepath.Clean(t)))
		if t == pID || t == pName || t == pBase || tNorm == pBase {
			return true
		}
		if strings.HasPrefix(pID, t) || strings.HasPrefix(pName, t) || strings.HasPrefix(pBase, t) {
			return true
		}
		if pPath != "" && (tNorm == pPath || strings.Contains(pPath, tNorm) || strings.Contains(pPath, t)) {
			return true
		}
	}
	return false
}
