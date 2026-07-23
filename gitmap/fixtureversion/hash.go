package fixtureversion

// Body-hash support for fixture stamps. The hash lets a stamped
// fixture detect *content drift* (someone hand-edited the body
// without bumping the generation) in addition to the existing
// generation-vs-MinGeneration freshness check.
//
// Design rules:
//   - Hash is computed over the body with the `// fixture-stamp:`
//     marker line REMOVED. Including the marker in its own input
//     would create a chicken-and-egg loop where every refresh
//     changes the hash that the marker now needs to record.
//   - SHA-256, hex-encoded, lowercase. We only ever surface the
//     first 12 hex chars in the marker (HashShortLen) — collisions
//     in a small fixture corpus are statistically irrelevant and
//     the short form keeps the marker line readable.
//   - Empty stamp.SHA means "do not enforce" — opt-in per fixture.

import (
	"crypto/sha256"
	"encoding/hex"
	"regexp"
)

// HashShortLen is the number of hex chars we keep in the marker.
// 12 hex = 48 bits — a fixture corpus of <10k entries has ~10^-9
// collision probability, plenty for this use case.
const HashShortLen = 12

// stampLineStripRe matches the entire fixture-stamp line plus its
// trailing newline (if any). Stripping the line — not just blanking
// it — keeps the hash stable when the marker is later widened or
// reformatted, since byte length stops mattering.
var stampLineStripRe = regexp.MustCompile(`(?m)^// fixture-stamp:[^\n]*\n?`)

// BodyHashExcludingMarker returns the lowercase hex SHA-256 (full
// 64 chars) of body with every fixture-stamp line removed. This is
// the canonical input format every consumer must use.
func BodyHashExcludingMarker(body string) string {
	stripped := stampLineStripRe.ReplaceAllString(body, "")
	sum := sha256.Sum256([]byte(stripped))

	return hex.EncodeToString(sum[:])
}

// ShortHash returns the first HashShortLen hex chars of fullHash.
// Returns the input unchanged if it is already shorter (defensive
// — callers should always pass a full hash).
func ShortHash(fullHash string) string {
	if len(fullHash) <= HashShortLen {
		return fullHash
	}

	return fullHash[:HashShortLen]
}

// HashMatches reports whether the body's current short hash matches
// the recorded stamp.SHA. An empty stamp.SHA is treated as "opt-out"
// and returns true unconditionally so old unstamped-hash fixtures
// keep working.
func HashMatches(body, recordedShortHash string) bool {
	if recordedShortHash == "" {
		return true
	}
	got := ShortHash(BodyHashExcludingMarker(body))

	return got == recordedShortHash
}

// shaFieldRe matches `sha=<hex>` up to the next whitespace. Lives
// alongside its semantic peers (the other hash helpers) rather than
// in bump.go so the rewrite logic stays cohesive with the field
// definition.
var shaFieldRe = regexp.MustCompile(`(sha=)\S+`)

// markerLineEndRe matches the end of the fixture-stamp line so we
// can append a new `sha=` field when the marker does not already
// have one. Captures the trailing newline (if any) for re-insertion.
var markerLineEndRe = regexp.MustCompile(`(?m)(^// fixture-stamp:[^\n]*?)(\r?\n|$)`)

// RewriteOrAppendSHA replaces the existing sha= field inside head,
// or appends `sha=<newSHA>` to the marker line when the field is
// missing. The append path preserves any trailing newline so the
// rest of the body remains byte-stable. Empty newSHA is a no-op.
func RewriteOrAppendSHA(head, newSHA string) string {
	if newSHA == "" {
		return head
	}
	if shaFieldRe.MatchString(head) {
		return shaFieldRe.ReplaceAllString(head, "${1}"+newSHA)
	}

	return markerLineEndRe.ReplaceAllString(head, "${1} sha="+newSHA+"${2}")
}
