package fixtureversion

// Bump helpers: pure string rewrites + an env-gated file-level
// auto-bumper that tests can call from MustValidateBody to make
// "fixture is stale" errors self-healing during a regenerate run.
//
// Design rules (mirrors the rest of fixtureversion):
//   - Pure functions take/return strings; no I/O, no testing.T.
//   - Side-effecting helpers live behind an explicit env gate
//     (FixtureAutoBumpEnv) so a normal `go test` run can NEVER
//     mutate source on disk by accident.
//   - Every helper is independently unit-testable.

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// FixtureAutoBumpEnv is the env-var name that opts MaybeAutoBumpFile
// in. Centralized so producers and CI configs cannot drift.
const FixtureAutoBumpEnv = "GITMAP_FIXTURE_AUTOBUMP"

// readWritePerm is the file mode used when rewriting fixture source
// files. 0o644 matches the project's source-file convention.
const readWritePerm = 0o644

// generationFieldRe matches `generation=<int>` inside a marker line,
// capturing the prefix for ReplaceAllString-friendly substitution.
var generationFieldRe = regexp.MustCompile(`(generation=)(\d+)`)

// minCurrentFieldRe matches `min-current=<int>` similarly.
var minCurrentFieldRe = regexp.MustCompile(`(min-current=)(\d+)`)

// createdForFieldRe matches `for=<token>` up to the next whitespace.
var createdForFieldRe = regexp.MustCompile(`(for=)\S+`)

// BumpRequest describes the in-place rewrite a caller wants. Zero
// values are no-ops: pass NewMinCurrent=0 / NewCreatedFor="" to
// leave them unchanged.
type BumpRequest struct {
	NewGeneration int
	NewMinCurrent int
	NewCreatedFor string
	// NewSHA, when non-empty, is written to the marker as `sha=<v>`.
	// If the marker already has a sha= field it is replaced; if not,
	// the field is appended to the marker line. Pass "" to leave any
	// existing sha= field untouched.
	NewSHA string
}

// BumpStampInBody returns body with the marker line's generation=
// (and, optionally, min-current= / for=) replaced. The marker is
// matched in the first 512 bytes only — same window ParseMarker
// uses, so the two cannot disagree on which line is "the stamp".
//
// Returns (newBody, true) on a successful rewrite; (body, false)
// when no marker is present so callers can decide whether that's
// an error.
func BumpStampInBody(body string, req BumpRequest) (string, bool) {
	headLen := minInt(len(body), 512)
	head, tail := body[:headLen], body[headLen:]
	if !markerLineRe.MatchString(head) {
		return body, false
	}
	head = rewriteGeneration(head, req.NewGeneration)
	head = rewriteMinCurrent(head, req.NewMinCurrent)
	head = rewriteCreatedFor(head, req.NewCreatedFor)
	head = RewriteOrAppendSHA(head, req.NewSHA)

	return head + tail, true
}

// rewriteGeneration replaces the captured digits with newGen when
// newGen > 0. Returns head unchanged otherwise.
func rewriteGeneration(head string, newGen int) string {
	if newGen <= 0 {
		return head
	}

	return generationFieldRe.ReplaceAllString(head, fmt.Sprintf("${1}%d", newGen))
}

// rewriteMinCurrent mirrors rewriteGeneration for the min-current
// field. Same zero-means-skip contract.
func rewriteMinCurrent(head string, newMin int) string {
	if newMin <= 0 {
		return head
	}

	return minCurrentFieldRe.ReplaceAllString(head, fmt.Sprintf("${1}%d", newMin))
}

// rewriteCreatedFor swaps the `for=...` value when a non-empty
// replacement is supplied. Whitespace-terminated to avoid eating
// the rest of the line.
func rewriteCreatedFor(head, newFor string) string {
	if newFor == "" {
		return head
	}

	return createdForFieldRe.ReplaceAllString(head, "${1}"+newFor)
}

// NextGeneration returns stamp.Generation+1, clamped to at least
// want.MinGeneration so a single bump always satisfies the failing
// expectation.
func NextGeneration(stamp Stamp, want Expectation) int {
	candidate := stamp.Generation + 1
	if candidate < want.MinGeneration {
		candidate = want.MinGeneration
	}

	return candidate
}

// MaybeAutoBumpFile is the env-gated escape hatch tests can call
// when MustValidateBody is about to fail. When FixtureAutoBumpEnv
// is set to "1", it reads path, rewrites the marker via
// BumpStampInBody, and writes the result back. Returns (true, nil)
// on a successful in-place rewrite; (false, nil) when the gate is
// off (the caller should then proceed to t.Fatal as usual).
//
// The two-step contract — pure rewrite + opt-in I/O — is what keeps
// `make fixtures-bump` honest: nothing on disk changes without the
// explicit env gate, so accidental autobumps in CI are impossible.
func MaybeAutoBumpFile(path string, req BumpRequest) (bool, error) {
	if !isAutoBumpEnabled() {
		return false, nil
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return false, fmt.Errorf("autobump: read %s: %w", path, err)
	}
	out, ok := BumpStampInBody(string(raw), req)
	if !ok {
		return false, fmt.Errorf("autobump: no stamp marker in %s", path)
	}
	if err := os.WriteFile(path, []byte(out), readWritePerm); err != nil {
		return false, fmt.Errorf("autobump: write %s: %w", path, err)
	}

	return true, nil
}

// isAutoBumpEnabled centralizes the env-gate check so future
// extensions (e.g. a `=dry-run` value) can grow in one place.
func isAutoBumpEnabled() bool {
	return os.Getenv(FixtureAutoBumpEnv) == "1"
}

// FormatBumpSummary renders a human-readable one-liner suitable for
// printing from a `make` target after MaybeAutoBumpFile succeeds.
func FormatBumpSummary(path string, oldGen, newGen int) string {
	return strings.Join([]string{
		"autobumped",
		path,
		fmt.Sprintf("generation %d -> %d", oldGen, newGen),
	}, " ")
}

// ParseGenerationFromBody is a convenience for shell tooling that
// just needs the current generation number (e.g. to print before/
// after in a make target). Returns -1 when no marker is present.
func ParseGenerationFromBody(body string) int {
	stamp, ok := ParseMarker(body)
	if !ok {
		return -1
	}

	return stamp.Generation
}

// minInt is a tiny stdlib gap-filler kept local so this file stays
// dependency-free.
func minInt(a, b int) int {
	if a < b {
		return a
	}

	return b
}
