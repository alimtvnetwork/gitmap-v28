// Package visibility — version.go: shared helpers that parse the
// `<base>-v<N>` repo-name convention used across gitmap. Exported so
// both the cmd-layer (make-last-*, except-latest, fuzzy fallback)
// and the store-layer (OwnerRepoNameIndex) share one parser.
package visibility

import (
	"regexp"
	"strconv"
)

// versionSuffixRE captures the trailing `-v<digits>` segment.
var versionSuffixRE = regexp.MustCompile(`(?i)^(.*)-v(\d+)$`)

// trailingDigitsRE captures `<base>-<digits>` (no `v`), used by the
// fuzzy fallback to spot user typos like `macro-ahk-51`.
var trailingDigitsRE = regexp.MustCompile(`^(.*)-(\d+)$`)

// ParseRepoNameMeta returns (base, version, ok) when `name` ends in
// `-v<digits>`, otherwise ("", 0, false). Case-insensitive on the v.
func ParseRepoNameMeta(name string) (string, int, bool) {
	m := versionSuffixRE.FindStringSubmatch(name)
	if m == nil {
		return "", 0, false
	}
	v, err := strconv.Atoi(m[2])
	if err != nil {
		return "", 0, false
	}

	return m[1], v, true
}

// SplitTrailingDigits returns (base, digits, ok) when `name` ends in
// `-<digits>` without a leading `v`. Used by AutoFixVDigitPatterns.
func SplitTrailingDigits(name string) (string, string, bool) {
	m := trailingDigitsRE.FindStringSubmatch(name)
	if m == nil {
		return "", "", false
	}

	return m[1], m[2], true
}

// HighestVersionedMatch walks `names`, keeping the highest -vN entry
// whose base equals `wantBase` (case-sensitive). Returns ("", 0, false)
// if no versioned sibling matches.
func HighestVersionedMatch(names []string, wantBase string) (string, int, bool) {
	bestName := ""
	bestVer := -1
	for _, n := range names {
		base, ver, ok := ParseRepoNameMeta(n)
		if !ok || base != wantBase {
			continue
		}
		if ver > bestVer {
			bestVer = ver
			bestName = n
		}
	}
	if bestVer < 0 {
		return "", 0, false
	}

	return bestName, bestVer, true
}
