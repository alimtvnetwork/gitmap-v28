package clonefrom

// summary_scheme.go — URL-scheme classification used by the
// `--output terminal` summary (summary_terminal.go) to render the
// "by mode:" tally. Split into its own file so the renderer stays
// focused and so future per-row previews can reuse ClassifyScheme
// without dragging in the whole summary surface.
//
// Classification logic intentionally mirrors validate.looksLikeGitURL:
// any URL that survived parse-time validation lands in a named
// bucket; truly unrecognized strings (rare — validation already
// rejects most) fall through to "other" so the tally is still
// total-preserving.

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// schemeOrder returns the canonical render order for the per-mode
// tally. https first (most common in practice), ssh next (typical
// interactive-developer alternative), then the less-common modes,
// with "other" always last as the catch-all.
func schemeOrder() []string {
	return []string{
		constants.CloneFromSchemeHTTPS,
		constants.CloneFromSchemeHTTP,
		constants.CloneFromSchemeSSH,
		constants.CloneFromSchemeSCP,
		constants.CloneFromSchemeGit,
		constants.CloneFromSchemeFile,
		constants.CloneFromSchemeOther,
	}
}

// tallySchemes counts each result's URL by scheme. Returns a map
// keyed by the scheme constants so the renderer can look up zero-
// count buckets without panicking on missing keys.
func tallySchemes(results []Result) map[string]int {
	out := make(map[string]int, len(schemeOrder()))
	for _, r := range results {
		out[ClassifyScheme(r.Row.URL)]++
	}

	return out
}

// ClassifyScheme picks one bucket for a URL. Exported so tests and
// any future per-row preview can reuse the rule without re-deriving
// the prefix table.
func ClassifyScheme(url string) string {
	url = strings.TrimSpace(url)
	if hit, ok := matchKnownScheme(url); ok {
		return hit
	}
	if looksLikeSCP(url) {
		return constants.CloneFromSchemeSCP
	}

	return constants.CloneFromSchemeOther
}

// matchKnownScheme walks the prefix table once. Kept tiny so
// ClassifyScheme stays under the function-length budget and the
// table itself is the single source of truth for known schemes.
func matchKnownScheme(url string) (string, bool) {
	prefixes := []struct{ prefix, scheme string }{
		{"https://", constants.CloneFromSchemeHTTPS},
		{"http://", constants.CloneFromSchemeHTTP},
		{"ssh://", constants.CloneFromSchemeSSH},
		{"git://", constants.CloneFromSchemeGit},
		{"file://", constants.CloneFromSchemeFile},
	}
	for _, p := range prefixes {
		if strings.HasPrefix(url, p.prefix) {
			return p.scheme, true
		}
	}

	return "", false
}
