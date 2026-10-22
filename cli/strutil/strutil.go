// Package strutil provides zero-allocation string comparison and normalization helpers.
package strutil

import "strings"

// EqualFoldAny reports whether target case-insensitively equals any candidate using strings.EqualFold.
func EqualFoldAny(target string, candidates ...string) bool {
	for _, c := range candidates {
		if strings.EqualFold(target, c) {
			return true
		}
	}

	return false
}

// EqualFoldAnyTrim reports whether target with whitespace trimmed case-insensitively equals any candidate.
func EqualFoldAnyTrim(target string, candidates ...string) bool {
	trimmed := strings.TrimSpace(target)

	return EqualFoldAny(trimmed, candidates...)
}

// NormalizeLowerTrim returns target with leading/trailing whitespace removed and converted to lowercase.
func NormalizeLowerTrim(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// JoinCSV joins a slice of strings with a comma separator, skipping empty entries.
func JoinCSV(tokens []string) string {
	cleaned := cleanTokens(tokens)
	if len(cleaned) == 0 {
		return ""
	}

	return strings.Join(cleaned, ",")
}

func cleanTokens(tokens []string) []string {
	var cleaned []string
	for _, t := range tokens {
		trimmed := strings.TrimSpace(t)
		if trimmed != "" {
			cleaned = append(cleaned, trimmed)
		}
	}

	return cleaned
}

// TruncateDots shortens s to at most maxLen bytes; when s is longer,
// the result ends with "..." (the last three bytes are replaced).
// Strings at or under maxLen are returned unchanged. Extracted to
// strutil after identical copies were found in cmdagy, cmdignore,
// cmdmail, and cmdpurge.
func TruncateDots(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}

	return s[:maxLen-3] + "..."
}

// TruncateEllipsis shortens s to at most max bytes; when s is longer,
// the result ends with "…" (the last byte is replaced). Strings at or
// under max are returned unchanged. Extracted to strutil after
// identical copies were found in cmdexec, cmdgroup, cmdmultigroup,
// and cmdwatch.
func TruncateEllipsis(s string, max int) string {
	if len(s) <= max {
		return s
	}

	return s[:max-1] + "…"
}
