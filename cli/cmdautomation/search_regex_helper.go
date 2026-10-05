package cmdautomation

import (
	"regexp"
	"strings"
)

// normalizeRegexPattern converts escaped alternation pipes (\|) to unescaped pipes (|).
func normalizeRegexPattern(pattern string) string {
	if !strings.Contains(pattern, `\|`) {
		return pattern
	}
	var sb strings.Builder
	runes := []rune(pattern)
	n := len(runes)
	for i := 0; i < n; i++ {
		isEscapedPipe := runes[i] == '\\' && i+1 < n && runes[i+1] == '|' && !isPrecededByEscapingBackslash(runes, i)
		if isEscapedPipe {
			sb.WriteRune('|')
			i++
			continue
		}
		sb.WriteRune(runes[i])
	}

	return sb.String()
}

func isPrecededByEscapingBackslash(runes []rune, idx int) bool {
	backslashes := 0
	for j := idx - 1; j >= 0 && runes[j] == '\\'; j-- {
		backslashes++
	}
	return backslashes%2 != 0
}

// hasRegexPatternSyntax tests whether a pattern contains common regex metacharacters.
func hasRegexPatternSyntax(pattern string) bool {
	if strings.Contains(pattern, `\|`) || strings.Contains(pattern, ".*") || strings.Contains(pattern, ".+") {
		return true
	}
	if strings.Contains(pattern, "|") || strings.Contains(pattern, "(?") {
		return true
	}
	if strings.Contains(pattern, `\d`) || strings.Contains(pattern, `\w`) || strings.Contains(pattern, `\s`) || strings.Contains(pattern, `\b`) {
		return true
	}
	if (strings.HasPrefix(pattern, "^") && len(pattern) > 1) || (strings.HasSuffix(pattern, "$") && len(pattern) > 1) {
		return true
	}
	return strings.Contains(pattern, "[") && strings.Contains(pattern, "]")
}

// autoPromoteRegex adjusts search options to regex mode if regex syntax is detected.
func autoPromoteRegex(opts *SearchOptions) {
	if opts.IsRegex {
		opts.Pattern = normalizeRegexPattern(opts.Pattern)
		return
	}
	if !hasRegexPatternSyntax(opts.Pattern) {
		return
	}
	candidate := normalizeRegexPattern(opts.Pattern)
	if _, err := regexp.Compile(candidate); err == nil {
		opts.IsRegex = true
		opts.Pattern = candidate
	}
}
