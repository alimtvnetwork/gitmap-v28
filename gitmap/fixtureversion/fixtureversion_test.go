package fixtureversion

// Unit tests for the fixture-stamp marker round-trip and the
// Validate failure messages. Locks the contract that downstream
// fixtures (cmd/fixrepo_rewrite_v9tov12_test.go and friends) rely
// on: round-trippable marker, name-required, generation+min-current
// gates produce actionable error text.

import (
	"strings"
	"testing"
)

func TestMarkerRoundTrip(t *testing.T) {
	in := Stamp{
		Name: "v9-to-v12", Generation: 2,
		MinCurrent: 12, CreatedFor: "v9->v12-width-cross",
	}
	body := Marker(in) + "\nrest of fixture\n"
	got, ok := ParseMarker(body)
	if !ok {
		t.Fatalf("ParseMarker failed on body: %q", body)
	}
	if got != in {
		t.Errorf("round-trip mismatch:\n  got  = %+v\n  want = %+v", got, in)
	}
}

func TestParseMarkerMissing(t *testing.T) {
	if _, ok := ParseMarker("no marker here\n"); ok {
		t.Error("expected ParseMarker to return false on unstamped body")
	}
}

func TestParseMarkerMissingName(t *testing.T) {
	body := "// fixture-stamp: generation=1 min-current=12 for=oops\n"
	if _, ok := ParseMarker(body); ok {
		t.Error("expected ParseMarker to reject a marker without name=")
	}
}

func TestValidateStaleGeneration(t *testing.T) {
	stamp := Stamp{Name: "demo", Generation: 1, MinCurrent: 12}
	err := Validate(stamp, Expectation{
		MinGeneration: 2, CurrentVersion: 12,
		RegenerateRecipe: "go test -run TestRegen -update",
	})
	if err == nil {
		t.Fatal("expected stale-generation error, got nil")
	}
	if !strings.Contains(err.Error(), "generation 1") ||
		!strings.Contains(err.Error(), "expects >=2") ||
		!strings.Contains(err.Error(), "TestRegen -update") {
		t.Errorf("error missing actionable details: %v", err)
	}
}

func TestValidateMinCurrentTooHigh(t *testing.T) {
	stamp := Stamp{Name: "demo", Generation: 2, MinCurrent: 14}
	err := Validate(stamp, Expectation{
		MinGeneration: 1, CurrentVersion: 12,
		RegenerateRecipe: "regen-cmd",
	})
	if err == nil {
		t.Fatal("expected min-current error, got nil")
	}
	if !strings.Contains(err.Error(), "min-current=14") ||
		!strings.Contains(err.Error(), "current=12") {
		t.Errorf("error missing version details: %v", err)
	}
}

func TestValidateOK(t *testing.T) {
	stamp := Stamp{Name: "demo", Generation: 2, MinCurrent: 12}
	err := Validate(stamp, Expectation{
		MinGeneration: 2, CurrentVersion: 12,
		RegenerateRecipe: "regen-cmd",
	})
	if err != nil {
		t.Errorf("unexpected validation error: %v", err)
	}
}
