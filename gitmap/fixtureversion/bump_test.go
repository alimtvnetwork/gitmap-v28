package fixtureversion

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const sampleBody = "// fixture-stamp: name=foo generation=1 min-current=12 for=v9->v12\nbody bytes follow\n"

func TestBumpStampInBodyRewritesGeneration(t *testing.T) {
	out, ok := BumpStampInBody(sampleBody, BumpRequest{NewGeneration: 5})
	if !ok {
		t.Fatalf("BumpStampInBody returned ok=false on a stamped body")
	}
	if !strings.Contains(out, "generation=5") {
		t.Fatalf("expected generation=5 in:\n%s", out)
	}
	if strings.Contains(out, "generation=1 ") {
		t.Fatalf("old generation=1 still present in:\n%s", out)
	}
	if !strings.HasSuffix(out, "body bytes follow\n") {
		t.Fatalf("tail of body mutated; got suffix mismatch:\n%s", out)
	}
}

func TestBumpStampInBodyRewritesAllFields(t *testing.T) {
	out, ok := BumpStampInBody(sampleBody, BumpRequest{
		NewGeneration: 3,
		NewMinCurrent: 14,
		NewCreatedFor: "v12->v14-bump",
	})
	if !ok {
		t.Fatalf("BumpStampInBody returned ok=false")
	}
	for _, want := range []string{"generation=3", "min-current=14", "for=v12->v14-bump"} {
		if !strings.Contains(out, want) {
			t.Fatalf("missing %q in rewritten body:\n%s", want, out)
		}
	}
}

func TestBumpStampInBodyZeroValuesAreNoOp(t *testing.T) {
	out, ok := BumpStampInBody(sampleBody, BumpRequest{})
	if !ok {
		t.Fatalf("BumpStampInBody returned ok=false")
	}
	if out != sampleBody {
		t.Fatalf("zero-value request mutated body:\nwant=%q\n got=%q", sampleBody, out)
	}
}

func TestBumpStampInBodyUnstampedReturnsFalse(t *testing.T) {
	body := "no marker here\nstill nothing\n"
	out, ok := BumpStampInBody(body, BumpRequest{NewGeneration: 9})
	if ok {
		t.Fatalf("expected ok=false on unstamped body")
	}
	if out != body {
		t.Fatalf("unstamped body was mutated; got:\n%s", out)
	}
}

func TestNextGenerationClampsToMin(t *testing.T) {
	cases := []struct {
		stampGen, wantMin, expected int
	}{
		{stampGen: 1, wantMin: 2, expected: 2},
		{stampGen: 1, wantMin: 5, expected: 5}, // clamp jumps multiple
		{stampGen: 4, wantMin: 2, expected: 5}, // already past min
	}
	for _, c := range cases {
		got := NextGeneration(Stamp{Generation: c.stampGen}, Expectation{MinGeneration: c.wantMin})
		if got != c.expected {
			t.Errorf("NextGeneration(stamp=%d, min=%d) = %d, want %d",
				c.stampGen, c.wantMin, got, c.expected)
		}
	}
}

func TestMaybeAutoBumpFileGateOff(t *testing.T) {
	t.Setenv(FixtureAutoBumpEnv, "")
	dir := t.TempDir()
	path := filepath.Join(dir, "fix.txt")
	mustWrite(t, path, sampleBody)
	ok, err := MaybeAutoBumpFile(path, BumpRequest{NewGeneration: 9})
	if err != nil {
		t.Fatalf("unexpected error with gate off: %v", err)
	}
	if ok {
		t.Fatalf("autobump ran with gate off (must be no-op)")
	}
	got := mustRead(t, path)
	if got != sampleBody {
		t.Fatalf("file mutated despite gate off:\n%s", got)
	}
}

func TestMaybeAutoBumpFileGateOnRewritesFile(t *testing.T) {
	t.Setenv(FixtureAutoBumpEnv, "1")
	dir := t.TempDir()
	path := filepath.Join(dir, "fix.txt")
	mustWrite(t, path, sampleBody)
	ok, err := MaybeAutoBumpFile(path, BumpRequest{NewGeneration: 7})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatalf("autobump did not run with gate on")
	}
	got := mustRead(t, path)
	if !strings.Contains(got, "generation=7") {
		t.Fatalf("file was not rewritten:\n%s", got)
	}
}

func TestMaybeAutoBumpFileGateOnUnstampedErrors(t *testing.T) {
	t.Setenv(FixtureAutoBumpEnv, "1")
	dir := t.TempDir()
	path := filepath.Join(dir, "fix.txt")
	mustWrite(t, path, "no marker here\n")
	_, err := MaybeAutoBumpFile(path, BumpRequest{NewGeneration: 7})
	if err == nil {
		t.Fatalf("expected error on unstamped file with gate on")
	}
}

func TestParseGenerationFromBody(t *testing.T) {
	if got := ParseGenerationFromBody(sampleBody); got != 1 {
		t.Errorf("ParseGenerationFromBody stamped = %d, want 1", got)
	}
	if got := ParseGenerationFromBody("nothing here"); got != -1 {
		t.Errorf("ParseGenerationFromBody unstamped = %d, want -1", got)
	}
}

func TestFormatBumpSummary(t *testing.T) {
	got := FormatBumpSummary("/tmp/x.go", 1, 2)
	want := "autobumped /tmp/x.go generation 1 -> 2"
	if got != want {
		t.Errorf("FormatBumpSummary = %q, want %q", got, want)
	}
}

// --- helpers --------------------------------------------------------

func mustWrite(t *testing.T, path, body string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func mustRead(t *testing.T, path string) string {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}

	return string(raw)
}
