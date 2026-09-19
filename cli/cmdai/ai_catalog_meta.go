package cmdai

import (
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

func buildDynamicMetadata(filename, dir string) ScriptMetadata {
	num := parseScriptNumber(filename)
	slug := parseScriptSlug(filename)
	fullPath := filepath.Join(dir, filename)
	content, _ := os.ReadFile(fullPath)
	desc := parseDocstringDesc(string(content))
	hasFix := strings.Contains(string(content), "--fix")
	cat := detectCategoryFromSlug(slug)

	return newScript(num, filename, slug, []string{slug}, cat, desc, hasFix)
}

func parseScriptNumber(filename string) string {
	var digits []rune
	for _, r := range filename {
		isDigit := unicode.IsDigit(r)
		if !isDigit {
			break
		}
		digits = append(digits, r)
	}

	hasDigits := len(digits) > 0
	if hasDigits {
		return string(digits)
	}

	return "dyn"
}

func parseScriptSlug(filename string) string {
	trimmed := strings.TrimSuffix(filename, ".py")
	parts := strings.SplitN(trimmed, "-", 2)
	hasPrefix := len(parts) == 2 && isAllDigits(parts[0])
	if hasPrefix {
		return parts[1]
	}

	return trimmed
}

func isAllDigits(s string) bool {
	isEmpty := len(s) == 0
	if isEmpty {
		return false
	}
	for _, r := range s {
		isDigit := unicode.IsDigit(r)
		if !isDigit {
			return false
		}
	}

	return true
}

func parseDocstringDesc(content string) string {
	start := strings.Index(content, `"""`)
	hasStart := start >= 0
	if !hasStart {
		return "Dynamic repository automation script"
	}

	rest := content[start+3:]
	end := strings.Index(rest, `"""`)
	hasEnd := end > 0
	if !hasEnd {
		return "Dynamic repository automation script"
	}

	lines := strings.Split(strings.TrimSpace(rest[:end]), "\n")
	first := strings.TrimSpace(lines[0])
	hasFirst := first != ""
	if hasFirst {
		return first
	}

	return "Dynamic repository automation script"
}

func detectCategoryFromSlug(slug string) ScriptCategoryType {
	isFix := strings.Contains(slug, "fix")
	if isFix {
		return CategoryFormatting
	}
	isLinter := strings.Contains(slug, "lint") || strings.Contains(slug, "check")
	if isLinter {
		return CategoryGuideline
	}
	isAudit := strings.Contains(slug, "audit")
	if isAudit {
		return CategoryAudit
	}

	return CategoryMaintenance
}
