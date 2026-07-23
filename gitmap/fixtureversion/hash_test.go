package fixtureversion

import (
	"strings"
	"testing"
)

const hashSampleBody = "// fixture-stamp: name=foo generation=1 min-current=12 for=v9->v12\nbody bytes follow\n"

func TestBodyHashExcludingMarkerStripsMarker(t *testing.T) {
	got := BodyHashExcludingMarker(hashSampleBody)
	if len(got) != 64 {
		t.Fatalf("expected 64 hex chars, got %d (%q)", len(got), got)
	}
	// Same body without the marker line must hash identically.
	noMarker := "body bytes follow\n"
	if got != BodyHashExcludingMarker(noMarker) {
		t.Fatalf("hash differs after marker strip; marker not removed for hashing")
	}
}

func TestBodyHashIsStableAcrossMarkerEdits(t *testing.T) {
	first := BodyHashExcludingMarker(hashSampleBody)
	bumped, ok := BumpStampInBody(hashSampleBody, BumpRequest{NewGeneration: 99, NewSHA: "deadbeefcafe"})
	if !ok {
		t.Fatalf("BumpStampInBody returned ok=false")
	}
	second := BodyHashExcludingMarker(bumped)
	if first != second {
		t.Fatalf("hash changed after marker-only edit:\n  first=%s\n second=%s", first, second)
	}
}

func TestShortHashTruncatesTo12(t *testing.T) {
	full := strings.Repeat("a", 64)
	got := ShortHash(full)
	if len(got) != HashShortLen {
		t.Errorf("ShortHash len = %d, want %d", len(got), HashShortLen)
	}
}

func TestHashMatchesEmptyOptOut(t *testing.T) {
	if !HashMatches("anything", "") {
		t.Errorf("empty recorded hash must always match (opt-out)")
	}
}

func TestHashMatchesDetectsDrift(t *testing.T) {
	correct := ShortHash(BodyHashExcludingMarker(hashSampleBody))
	if !HashMatches(hashSampleBody, correct) {
		t.Errorf("HashMatches with correct sha returned false")
	}
	if HashMatches(hashSampleBody+"trailing edit\n", correct) {
		t.Errorf("HashMatches did not detect trailing-edit drift")
	}
}

func TestRewriteOrAppendSHAAppendsWhenMissing(t *testing.T) {
	out := RewriteOrAppendSHA(hashSampleBody, "abc123def456")
	if !strings.Contains(out, "sha=abc123def456\n") {
		t.Fatalf("sha= field not appended at end of marker:\n%s", out)
	}
}

func TestRewriteOrAppendSHAReplacesExisting(t *testing.T) {
	body := "// fixture-stamp: name=foo generation=2 sha=oldhashvalue\nrest\n"
	out := RewriteOrAppendSHA(body, "newhashvalue")
	if !strings.Contains(out, "sha=newhashvalue") {
		t.Fatalf("sha= field not replaced:\n%s", out)
	}
	if strings.Contains(out, "oldhashvalue") {
		t.Fatalf("old sha= value still present:\n%s", out)
	}
}

func TestValidateBodyDetectsDrift(t *testing.T) {
	stamp := Stamp{Name: "foo", Generation: 1, MinCurrent: 12, SHA: "0000deadbeef"}
	want := Expectation{MinGeneration: 1, CurrentVersion: 12, RegenerateRecipe: "regen"}
	err := ValidateBody(hashSampleBody, stamp, want)
	if err == nil {
		t.Fatalf("expected hash-mismatch error, got nil")
	}
	if !strings.Contains(err.Error(), "body hash mismatch") {
		t.Fatalf("error message missing 'body hash mismatch':\n%v", err)
	}
}

func TestValidateBodyPassesWhenHashMatches(t *testing.T) {
	correct := ShortHash(BodyHashExcludingMarker(hashSampleBody))
	stamp := Stamp{Name: "foo", Generation: 1, MinCurrent: 12, SHA: correct}
	want := Expectation{MinGeneration: 1, CurrentVersion: 12, RegenerateRecipe: "regen"}
	if err := ValidateBody(hashSampleBody, stamp, want); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateBodyEmptySHAIsOptOut(t *testing.T) {
	stamp := Stamp{Name: "foo", Generation: 1, MinCurrent: 12} // SHA empty
	want := Expectation{MinGeneration: 1, CurrentVersion: 12, RegenerateRecipe: "regen"}
	if err := ValidateBody("anything totally different\n", stamp, want); err != nil {
		t.Fatalf("empty stamp.SHA must opt-out of hash check, got: %v", err)
	}
}

func TestStampMarkerRoundTripsSHA(t *testing.T) {
	in := Stamp{Name: "foo", Generation: 3, MinCurrent: 12, CreatedFor: "x", SHA: "abc123def456"}
	body := Marker(in) + "\nbody\n"
	out, ok := ParseMarker(body)
	if !ok {
		t.Fatalf("ParseMarker returned ok=false")
	}
	if out.SHA != in.SHA {
		t.Errorf("SHA not round-tripped: got %q, want %q", out.SHA, in.SHA)
	}
}
