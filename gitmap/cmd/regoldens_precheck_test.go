package cmd

// Tests for the regoldens determinism pre-check helpers. The
// goroutine-spawning runPrecheckGoTest is exercised by CI; here we
// pin down the pure marker-scan logic so a future refactor of the
// goldenguard message cannot silently neutralize the pre-check.

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func TestPrecheckFoundNonDeterminism_DetectsMarker(t *testing.T) {
	out := []byte(`--- FAIL: TestFoo (0.01s)
    foo_test.go:42: goldenguard: writer "FooReport" is non-deterministic — run 1 vs run 2 differ
FAIL`)
	if !precheckFoundNonDeterminism(out) {
		t.Fatalf("expected marker hit, got miss for output: %s", out)
	}
}

func TestPrecheckFoundNonDeterminism_IgnoresGateOnlyFailure(t *testing.T) {
	out := []byte(`--- FAIL: TestFoo (0.01s)
    goldenguard.go:72: golden update requested but GITMAP_ALLOW_GOLDEN_UPDATE is not set
FAIL`)
	if precheckFoundNonDeterminism(out) {
		t.Fatalf("expected miss (gate-only failure), got hit for output: %s", out)
	}
}

func TestPrecheckFoundNonDeterminism_EmptyOutput_Miss(t *testing.T) {
	if precheckFoundNonDeterminism(nil) {
		t.Fatalf("nil output must not register as non-determinism")
	}
}

func TestPrecheckMarkerConstant_NonEmpty(t *testing.T) {
	if constants.RegoldensNonDetMarker == "" {
		t.Fatalf("RegoldensNonDetMarker must be a non-empty substring " +
			"or the pre-check will accept ALL output as deterministic")
	}
}
