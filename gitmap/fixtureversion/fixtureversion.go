package fixtureversion

// Package fixtureversion stamps test fixtures with a machine-readable
// generation marker so stale on-disk fixtures fail loudly instead of
// producing confusing downstream assertion errors (off-by-one
// rewrite counts, mismatched golden bytes, etc).
//
// The historical failure mode this closes: a fixture body baked in
// at generation N is read by a test that has since been updated to
// expect generation N+1. The byte-level rewriter or the golden
// comparer then trips with a low-level diff that does not point at
// the real cause (the fixture itself is out of date). Embedding a
// `// fixture-stamp:` line as the first content line lets every test
// call `MustValidate` once and fail with an explicit, actionable
// message ("fixture <name> is generation 1, test expects >=2 — re-
// generate via <recipe>").
//
// Stamps are intentionally minimal: a Name (so multiple fixtures in
// the same package don't collide), an integer Generation that is
// bumped on every intentional regeneration, an optional MinCurrent
// (the lowest project version this fixture is valid against — keeps
// fixtures from silently outliving their version-window), and a
// free-form CreatedFor note for human readers.

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// MarkerPrefix is the canonical line-prefix every stamped fixture
// embeds. Centralized here so producers (Marker) and consumers
// (ParseMarker) cannot drift.
const MarkerPrefix = "// fixture-stamp:"

// Stamp describes a single fixture's identity + freshness contract.
// All four fields participate in Marker / ParseMarker round-trips.
type Stamp struct {
	Name       string
	Generation int
	MinCurrent int
	CreatedFor string
	// SHA is the optional short (HashShortLen-char) SHA-256 of the
	// fixture body with all fixture-stamp lines stripped. Empty
	// means "do not enforce content drift" so older unstamped-hash
	// fixtures keep validating. See hash.go.
	SHA string
}

// Marker renders s as the one-line comment that should be embedded
// at the top of the fixture file. The format is intentionally
// `key=value` pairs separated by spaces so it stays readable in a
// raw text editor and trivially parseable by ParseMarker.
//
// Example output:
//
//	// fixture-stamp: name=v9-to-v12 generation=2 min-current=12 for=v9->v12 sha=abc123def456
func Marker(s Stamp) string {
	base := fmt.Sprintf("%s name=%s generation=%d min-current=%d for=%s",
		MarkerPrefix, s.Name, s.Generation, s.MinCurrent, s.CreatedFor)
	if s.SHA == "" {
		return base
	}

	return base + " sha=" + s.SHA
}

// markerLineRe matches the marker line anywhere in a body. We only
// inspect the first ~512 bytes (see ParseMarker) so the regex stays
// cheap regardless of fixture size.
var markerLineRe = regexp.MustCompile(
	`(?m)^// fixture-stamp:\s+(.+)$`)

// ParseMarker extracts the first stamp marker from body. Returns
// (Stamp{}, false) if no marker is present or if the marker is
// malformed — callers should treat both as "fixture is unstamped".
func ParseMarker(body string) (Stamp, bool) {
	head := body
	if len(head) > 512 {
		head = head[:512]
	}
	m := markerLineRe.FindStringSubmatch(head)
	if m == nil {
		return Stamp{}, false
	}

	return parseFields(m[1])
}

// parseFields splits the `key=value key=value ...` payload of a
// marker line into a Stamp. Unknown keys are ignored so future
// extensions remain backward-compatible with older test binaries.
func parseFields(payload string) (Stamp, bool) {
	out := Stamp{}
	hasName := false
	for _, raw := range splitFields(payload) {
		key, val, ok := strings.Cut(raw, "=")
		if !ok {
			continue
		}
		switch key {
		case "name":
			out.Name = val
			hasName = true
		case "generation":
			n, err := strconv.Atoi(val)
			if err != nil {
				return Stamp{}, false
			}
			out.Generation = n
		case "min-current":
			n, err := strconv.Atoi(val)
			if err != nil {
				return Stamp{}, false
			}
			out.MinCurrent = n
		case "for":
			out.CreatedFor = val
		case "sha":
			out.SHA = val
		}
	}
	if !hasName {
		return Stamp{}, false
	}

	return out, true
}

// splitFields splits payload on whitespace, but treats the trailing
// `for=...` value as a single token (it is the last key and may
// contain spaces in a future revision). Today every value is a
// single token, so a plain Fields() suffices — kept as a named
// helper to make the contract explicit for future maintainers.
func splitFields(payload string) []string {
	return strings.Fields(payload)
}
