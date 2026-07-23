---
Slug: pattern-engine
Status: pending
Created: 2026-06-06
Parent: 01-bulk-visibility-mapub-mapri
---

# Subtask 01 — Wildcard pattern engine

## Goal
Translate user patterns (`exact`, `prefix*`, `*contains*`, `prefix*suffix`,
`a*b*c`) into Go matchers that run against a list of repo names returned by
`gh repo list --json name`.

## Rules (verbatim from user spec)
1. `*` is the ONLY wildcard token.
2. Literal text between `*` markers = substring constraint, IN ORDER.
3. Repo name match is anchored: `^segment0.*segment1.*…segmentN$` where
   leading `*` makes the first anchor optional and trailing `*` makes the
   last anchor optional.

## Implementation
- New file: `gitmap/visibility/pattern.go`
- `type Pattern struct { Raw string; segments []string; anchorStart, anchorEnd bool }`
- `func ParsePattern(raw string) (Pattern, error)` — empty raw → error; reject
  patterns of only `*` (would match everything, surface error per spec safety).
- `func (p Pattern) Match(name string) bool` — iterate segments left-to-right
  with `strings.Index`, advancing offset.
- `func MatchAny(patterns []Pattern, names []string) map[string]Pattern` —
  returns name→first-matching-pattern; preserves input order for display.

## Tests (`gitmap/visibility/pattern_test.go`)
- Exact: `vibe-ext-v5` matches only that.
- `macro*`: matches `macro`, `macromate`, NOT `submacro`.
- `*mate*`: matches `mate`, `xxxmate`, `matexxx`, `xxmatexx`; NOT `mat`.
- `lotus*v1-4`: matches `lotus-v1-4`, `lotusxyzv1-4`; NOT `lotusv1` or `xlotusv1-4`.
- `a*b*c`: ordering — `axbyc` ✓, `cba` ✗.
- Empty / bare `*` → error.

## Verification
`go test ./gitmap/visibility -run Pattern -v` green.
