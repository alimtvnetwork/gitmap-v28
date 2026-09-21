package cmd

import (
	"strings"
	"unicode"
)

func isSlugChar(r rune) bool {
	isAlpha := unicode.IsLetter(r)
	isDigit := unicode.IsDigit(r)
	isDash := r == '-'

	return isAlpha || isDigit || isDash
}

func sanitizeRuneToSlug(r rune) rune {
	isSpaceOrSep := unicode.IsSpace(r) || r == '_' || r == '.' || r == '/' || r == '\\'
	if isSpaceOrSep {
		return '-'
	}

	lower := unicode.ToLower(r)
	if isSlugChar(lower) {
		return lower
	}

	return -1
}

func appendSlugRune(b *strings.Builder, r rune, hasHyphen bool) bool {
	isHyphen := r == '-'
	isNewHyphen := isHyphen && !hasHyphen
	if isNewHyphen {
		b.WriteRune(r)

		return true
	}
	if isHyphen {
		return true
	}

	b.WriteRune(r)

	return false
}

func collapseHyphens(raw string) string {
	var b strings.Builder
	hasHyphen := false
	for _, r := range raw {
		hasHyphen = appendSlugRune(&b, r, hasHyphen)
	}

	return b.String()
}

// SlugifyRepoName converts a human-readable repository name into a valid GitHub slug.
func SlugifyRepoName(name string) string {
	trimmed := strings.TrimSpace(name)
	mapped := strings.Map(sanitizeRuneToSlug, trimmed)
	collapsed := collapseHyphens(mapped)

	return strings.Trim(collapsed, "-")
}
