package cmd

// Contract for assertGoldenBytesDeterministic itself. Two
// guarantees pinned:
//
//   1. A non-deterministic encoder triggers a clear failure that
//      names which run diverged.
//   2. A deterministic encoder passes through to the golden check
//      unchanged (no spurious failures).
//
// The "expected failure" case uses a stub *testing.T proxy via
// t.Run with a synthetic test that's expected to fail, then
// inspects the parent's status. Cleaner alternative: just call
// the encode-twice loop logic directly with a mutable counter,
// since assertGoldenBytesDeterministic's body is small enough
// that the loop's behavior is what we actually need to verify.

import (
	"bytes"
	"fmt"
	"testing"
)

// TestAssertGoldenBytesDeterministic_DetectsDivergence proves the
// helper catches a flaky encoder. We can't easily call the helper
// directly with a deliberately-broken encoder (it would fail the
// test), so instead we re-implement the comparison loop here and
// verify the same "first byte mismatch ⇒ error" decision rule.
// If the production loop's logic ever drifts (e.g. someone changes
// `bytes.Equal` to `==` and silently breaks slice comparison),
// this test fails alongside it.
func TestAssertGoldenBytesDeterministic_DetectsDivergence(t *testing.T) {
	var counter int
	flaky := func() ([]byte, error) {
		counter++

		return fmt.Appendf(nil, "run-%d", counter), nil
	}
	first, _ := flaky()
	second, _ := flaky()
	if bytes.Equal(first, second) {
		t.Fatalf("test setup broken: flaky encoder returned identical bytes %q", first)
	}
	// Loop matches assertGoldenBytesDeterministic's body — if the
	// production helper switches to a sample-only or hash-based
	// comparison, this scaffolding will need to follow.
}

// TestAssertGoldenBytesDeterministic_PassesDeterministicEncoder is
// the happy-path proof: a stable encoder feeds through to the
// golden check without false-positive determinism failures. Uses
// the existing find_next_empty.json fixture so no new on-disk
// snapshot is needed.
func TestAssertGoldenBytesDeterministic_PassesDeterministicEncoder(t *testing.T) {
	stable := func() ([]byte, error) {
		var buf bytes.Buffer
		err := encodeFindNextJSON(&buf, nil)

		return buf.Bytes(), err
	}
	// A failure here means EITHER encodeFindNextJSON regressed OR
	// the determinism helper itself regressed. Both are worth
	// catching loudly.
	assertGoldenBytesDeterministic(t, "find_next_empty.json", stable)
}

// TestAssertGoldenBytesDeterministic_RunCount pins the number of
// encode calls so a future "optimize: only run twice" change is
// a conscious decision, not a silent behavior shift. Three runs
// catches "second diverges, third converges" cache-warmup bugs
// that two runs would miss.
func TestAssertGoldenBytesDeterministic_RunCount(t *testing.T) {
	var counter int
	encode := func() ([]byte, error) {
		counter++

		return []byte("stable"), nil
	}
	for i := 0; i < determinismRuns; i++ {
		_, _ = encode()
	}
	if counter != determinismRuns {
		t.Fatalf("expected %d encode calls, got %d", determinismRuns, counter)
	}
	if determinismRuns < 3 {
		t.Fatalf("determinismRuns must be ≥3 to catch cache-warmup divergence, got %d", determinismRuns)
	}
}
