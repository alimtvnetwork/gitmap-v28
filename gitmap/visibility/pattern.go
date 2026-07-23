// Package visibility — wildcard pattern engine for the bulk visibility
// commands (make-all-public, make-all-private, MAPUB, MAPRI).
//
// Pattern grammar (per spec/01-app/116-bulk-visibility-mapub-mapri.md §3):
//
//	macro       → exact match
//	macro*      → prefix
//	*macro      → suffix
//	*macro*     → contains
//	lotus*v1    → prefix + suffix
//	a*b*c       → ordered substrings (greedy)
//	*           → ERROR (bare wildcard is a footgun)
//
// Only `*` is special. No regex, no `?`, no character classes.
// Matching is case-sensitive (GitHub/GitLab slugs are case-sensitive on
// the API side).
package visibility

import (
	"fmt"
	"strings"
)

// Pattern is the compiled form of a single user-supplied token. The
// engine slices the input on `*` and the matcher walks the parts in
// order, anchoring the first and last where appropriate.
type Pattern struct {
	Raw     string   // original user input (for "matched by" display)
	parts   []string // non-empty literal segments between `*` chars
	anchorL bool     // matched substring must start at position 0
	anchorR bool     // matched substring must end at end of input
}

// ParsePattern compiles one token. Returns a typed error for the bare
// `*` case so callers can surface the spec-mandated footgun message.
func ParsePattern(raw string) (Pattern, error) {
	trimmed := strings.TrimSpace(raw)
	if len(trimmed) == 0 {
		return Pattern{}, fmt.Errorf("Error: empty pattern (operation: parse-pattern, reason: blank token)")
	}
	if trimmed == "*" {
		return Pattern{}, fmt.Errorf("Error: bare '*' pattern is refused at %q (operation: parse-pattern, reason: would match every repo under the owner)", raw)
	}

	parts, anchorL, anchorR := splitPatternParts(trimmed)
	if len(parts) == 0 {
		return Pattern{}, fmt.Errorf("Error: pattern %q has no literal segments (operation: parse-pattern, reason: only wildcards)", raw)
	}

	return Pattern{Raw: trimmed, parts: parts, anchorL: anchorL, anchorR: anchorR}, nil
}

// splitPatternParts returns the literal segments and the anchor flags.
// `anchorL` is true when the pattern does NOT start with `*`; `anchorR`
// is true when the pattern does NOT end with `*`.
func splitPatternParts(raw string) ([]string, bool, bool) {
	anchorL := strings.HasPrefix(raw, "*") == false
	anchorR := strings.HasSuffix(raw, "*") == false

	rawParts := strings.Split(raw, "*")
	parts := make([]string, 0, len(rawParts))
	for _, p := range rawParts {
		if len(p) > 0 {
			parts = append(parts, p)
		}
	}

	return parts, anchorL, anchorR
}

// Matches returns true if `name` satisfies the pattern under the rules
// above. Empty `name` never matches (no real repo has an empty name).
func (p Pattern) Matches(name string) bool {
	if len(name) == 0 {
		return false
	}
	if len(p.parts) == 0 {
		return false
	}

	cursor := 0
	for i, part := range p.parts {
		idx := strings.Index(name[cursor:], part)
		if idx < 0 {
			return false
		}
		if i == 0 && p.anchorL && idx > 0 {
			return false
		}
		cursor += idx + len(part)
	}

	if p.anchorR && cursor != len(name) {
		return false
	}

	return true
}

// ParsePatternList splits a comma-separated raw list into compiled
// Patterns, trimming spaces, deduping by raw form, and rejecting any
// empty token. Returns the first parse error with the offending token
// index (1-based) for clear user feedback.
func ParsePatternList(raw string) ([]Pattern, error) {
	if len(strings.TrimSpace(raw)) == 0 {
		return nil, fmt.Errorf("Error: empty pattern list (operation: parse-pattern-list, reason: arg is blank)")
	}

	tokens := strings.Split(raw, ",")
	seen := make(map[string]bool, len(tokens))
	out := make([]Pattern, 0, len(tokens))
	for i, tok := range tokens {
		trimmed := strings.TrimSpace(tok)
		if len(trimmed) == 0 {
			return nil, fmt.Errorf("Error: empty pattern at token %d (operation: parse-pattern-list, reason: blank between commas)", i+1)
		}
		if seen[trimmed] {
			continue
		}
		seen[trimmed] = true

		pat, err := ParsePattern(trimmed)
		if err != nil {
			return nil, fmt.Errorf("Error: token %d %q: %w", i+1, trimmed, err)
		}
		out = append(out, pat)
	}

	return out, nil
}
