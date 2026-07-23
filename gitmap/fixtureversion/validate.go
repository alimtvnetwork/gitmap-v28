package fixtureversion

// Validation entry points used by tests. Kept in a sibling file so
// fixtureversion.go stays focused on the data model + marker
// rendering, and validate.go owns the failure-message contract.

import (
	"fmt"
	"testing"
)

// Expectation is what a test asserts about a fixture it loads.
// MinGeneration is the lowest acceptable Stamp.Generation; anything
// older means the fixture body has not been re-stamped after the
// test was updated and must be regenerated. CurrentVersion is the
// project version the test is exercising — Stamp.MinCurrent must
// be <= CurrentVersion or the fixture is from a too-old version
// window. RegenerateRecipe is a free-form one-line command (e.g.
// `go test ./gitmap/cmd -run TestRegenFixRepoV9ToV12 -update`)
// surfaced verbatim in the failure message so the human knows
// exactly how to fix it.
type Expectation struct {
	MinGeneration    int
	CurrentVersion   int
	RegenerateRecipe string
}

// Validate checks stamp against want and returns an actionable
// error describing exactly what is stale. Returns nil on success.
// Pure function — no t.* calls — so it can be reused outside
// *testing.T contexts (e.g. a fixture-audit CLI).
//
// Validate does NOT check the body hash — callers that have the
// body bytes should use ValidateBody instead. Splitting them keeps
// the pure-stamp-only check usable in contexts (audit CLIs,
// dashboards) that have the marker but not the full body.
func Validate(stamp Stamp, want Expectation) error {
	if stamp.Name == "" {
		return fmt.Errorf("fixture is unstamped: add %q as the first line",
			Marker(Stamp{Name: "<name>", Generation: 1, MinCurrent: want.CurrentVersion}))
	}
	if stamp.Generation < want.MinGeneration {
		return fmt.Errorf(
			"fixture %q is generation %d, test expects >=%d (stale fixture)\n"+
				"  regenerate via: %s",
			stamp.Name, stamp.Generation, want.MinGeneration, want.RegenerateRecipe)
	}
	if want.CurrentVersion > 0 && stamp.MinCurrent > want.CurrentVersion {
		return fmt.Errorf(
			"fixture %q requires min-current=%d but test runs at current=%d\n"+
				"  regenerate via: %s",
			stamp.Name, stamp.MinCurrent, want.CurrentVersion, want.RegenerateRecipe)
	}

	return nil
}

// ValidateBody runs Validate and additionally checks the body's
// content hash against the stamp's recorded SHA. An empty
// stamp.SHA opts out (validation passes). Use this in tests so a
// hand-edit to the fixture body without a corresponding marker
// refresh fails loudly.
func ValidateBody(body string, stamp Stamp, want Expectation) error {
	if err := Validate(stamp, want); err != nil {
		return err
	}
	if HashMatches(body, stamp.SHA) {
		return nil
	}
	got := ShortHash(BodyHashExcludingMarker(body))

	return fmt.Errorf(
		"fixture %q body hash mismatch: marker records sha=%s but actual body sha=%s\n"+
			"  the fixture was edited without refreshing its // fixture-stamp: marker.\n"+
			"  regenerate via: %s",
		stamp.Name, stamp.SHA, got, want.RegenerateRecipe)
}

// MustValidateBody is the one-call helper tests use: parses the
// marker out of body, validates it against want, and t.Fatals with
// the actionable error on any mismatch. Use this at the top of any
// test that consumes a stamped fixture so a stale fixture trips
// here instead of corrupting downstream assertions.
func MustValidateBody(t *testing.T, body string, want Expectation) Stamp {
	t.Helper()
	stamp, ok := ParseMarker(body)
	if !ok {
		t.Fatalf("fixture is unstamped or marker malformed: add %q as the first line",
			Marker(Stamp{Name: "<name>", Generation: 1, MinCurrent: want.CurrentVersion}))
	}
	if err := ValidateBody(body, stamp, want); err != nil {
		t.Fatal(err)
	}

	return stamp
}

// MustValidateBodyWithAutobump is a drop-in upgrade of
// MustValidateBody that, when the FixtureAutoBumpEnv gate is set,
// rewrites sourcePath in place to satisfy want.MinGeneration before
// failing. Use it from tests whose fixture body lives as a string
// literal in a Go source file (the common case): pass the path of
// that .go file as sourcePath and a stale fixture becomes a
// self-healing condition under `make fixtures-bump`.
//
// Without the env gate this behaves identically to MustValidateBody
// — there is zero risk of accidental source mutation in CI / local
// `go test` runs.
func MustValidateBodyWithAutobump(t *testing.T, body, sourcePath string, want Expectation) Stamp {
	t.Helper()
	stamp, ok := ParseMarker(body)
	if !ok {
		t.Fatalf("fixture is unstamped or marker malformed: add %q as the first line",
			Marker(Stamp{Name: "<name>", Generation: 1, MinCurrent: want.CurrentVersion}))
	}
	validateErr := ValidateBody(body, stamp, want)
	if validateErr == nil {
		return stamp
	}
	if !tryAutobumpAndReport(t, body, stamp, sourcePath, want) {
		t.Fatal(validateErr)
	}

	return stamp
}

// tryAutobumpAndReport runs MaybeAutoBumpFile and logs the outcome.
// Returns true when an autobump was actually applied (caller should
// treat the test as passing for this run since the rewrite landed).
// Refreshes BOTH the generation and the body-hash so a content-only
// drift heals in the same pass as a generation drift.
func tryAutobumpAndReport(t *testing.T, body string, stamp Stamp, sourcePath string, want Expectation) bool {
	t.Helper()
	newGen := NextGeneration(stamp, want)
	newSHA := ShortHash(BodyHashExcludingMarker(body))
	bumped, err := MaybeAutoBumpFile(sourcePath, BumpRequest{
		NewGeneration: newGen,
		NewSHA:        newSHA,
	})
	if err != nil {
		t.Fatalf("autobump attempt failed: %v", err)
	}
	if !bumped {
		return false
	}
	t.Log(FormatBumpSummary(sourcePath, stamp.Generation, newGen))
	t.Logf("autobumped sha %s -> %s", stamp.SHA, newSHA)
	t.Log("re-run the test without GITMAP_FIXTURE_AUTOBUMP=1 to confirm the bumped fixture passes")

	return true
}
