package goldenguard

// Determinism pre-check for golden writers.
//
// Why this exists: the existing two-key gate (AllowUpdate) protects
// against ACCIDENTALLY rewriting goldens, but it does not protect
// against rewriting them with NON-DETERMINISTIC bytes. A writer that
// shuffles map keys, embeds wall-clock timestamps, or formats floats
// locale-dependently will happily produce different bytes on each
// run. If a contributor regenerates such a fixture, pass-1 succeeds
// (the gate is satisfied), pass-2 (verify) fails — but the now-bad
// fixture is already on disk and may be committed by mistake.
//
// AssertWriterDeterministic runs the writer N times in-process and
// fails the test if any two runs disagree byte-for-byte. Call it
// BEFORE writing to disk so a non-deterministic writer is rejected
// at the source instead of poisoning testdata/.
//
// AllowUpdateAfterDeterminism bundles the check with the existing
// gate so callers can switch from `AllowUpdate(t, trigger)` to
// `AllowUpdateAfterDeterminism(t, trigger, label, writer)` and get
// both protections in one line.

import (
	"bytes"
	"fmt"
	"testing"
)

// determinismRunCount is the number of times a writer must produce
// identical bytes to be considered deterministic. Three runs catch
// the common non-determinism sources (map iteration, time.Now,
// random IDs) without making the check expensive. Do NOT lower this
// to 2 — some non-determinism only surfaces on the 3rd+ iteration
// (e.g. Go's randomized map seed is per-process, but per-call code
// paths can still re-roll on subsequent invocations).
const determinismRunCount = 3

// determinismMaxDiffBytes caps how many bytes of the divergence
// snippet are included in the failure message. Keeps the test log
// readable when the writer emits multi-MB blobs.
const determinismMaxDiffBytes = 200

// WriterFn is the contract for a golden writer: it receives no args
// (the caller closes over inputs) and returns the bytes that would
// be written to the fixture on disk. Returning an error fails the
// determinism check immediately — a writer that errors out cannot
// also be deterministic.
type WriterFn func() ([]byte, error)

// AssertWriterDeterministic runs writer determinismRunCount times
// and t.Fatalf's if any pair of runs produces different bytes.
// label is included in failure messages so a test asserting on
// multiple writers can pinpoint which one drifted.
func AssertWriterDeterministic(t *testing.T, label string, writer WriterFn) {
	t.Helper()
	assertWriterDeterministicOn(t, label, writer)
}

// fataler is the slice of *testing.T used by the determinism check.
// Pulled out as an interface so unit tests can supply a fake that
// records Fatalf calls without aborting the test runner. Real users
// always pass a *testing.T.
type fataler interface {
	Helper()
	Fatalf(format string, args ...interface{})
}

// assertWriterDeterministicOn is the interface-based core. The
// public AssertWriterDeterministic is a thin shim so callers don't
// have to know about the fataler abstraction.
func assertWriterDeterministicOn(t fataler, label string, writer WriterFn) {
	t.Helper()
	runs, err := collectWriterRuns(writer)
	if err != nil {
		t.Fatalf("goldenguard: writer %q failed during determinism check: %v",
			label, err)
		return
	}
	assertAllRunsEqualOn(t, label, runs)
}

// collectWriterRuns invokes writer determinismRunCount times and
// returns the resulting byte slices. Any error short-circuits the
// loop because a writer must succeed deterministically before its
// output can be compared.
func collectWriterRuns(writer WriterFn) ([][]byte, error) {
	out := make([][]byte, 0, determinismRunCount)
	for i := 0; i < determinismRunCount; i++ {
		b, err := writer()
		if err != nil {
			return nil, fmt.Errorf("run %d/%d: %w", i+1, determinismRunCount, err)
		}
		out = append(out, b)
	}
	return out, nil
}

// nonDeterministicWriterMsgFmt is the failure template used when a
// writer's bytes diverge across runs. Centralized so the message
// stays one source of truth and assertAllRunsEqualOn fits the
// 15-line function budget.
const nonDeterministicWriterMsgFmt = "goldenguard: writer %q is non-deterministic — " +
	"run 1 vs run %d differ (%d vs %d bytes).\n" +
	"  run 1 head: %s\n  run %d head: %s\n" +
	"Fix the writer (likely culprits: map iteration, time.Now, " +
	"randomness, locale-dependent formatting) BEFORE regenerating fixtures."

// assertAllRunsEqualOn compares every run after the first against
// run[0]. The first divergence triggers Fatalf via the centralized
// failure template above. Takes a fataler so unit tests can inject
// a fake.
func assertAllRunsEqualOn(t fataler, label string, runs [][]byte) {
	t.Helper()
	for i := 1; i < len(runs); i++ {
		if bytes.Equal(runs[0], runs[i]) {
			continue
		}
		t.Fatalf(nonDeterministicWriterMsgFmt,
			label, i+1, len(runs[0]), len(runs[i]),
			snippet(runs[0]), i+1, snippet(runs[i]))
		return
	}
}

// snippet returns up to determinismMaxDiffBytes of b, %q-quoted so
// invisible characters (newlines, tabs, NULs) are visible in the
// failure message.
func snippet(b []byte) string {
	if len(b) > determinismMaxDiffBytes {
		return fmt.Sprintf("%q… (truncated)", b[:determinismMaxDiffBytes])
	}
	return fmt.Sprintf("%q", b)
}

// AllowUpdateAfterDeterminism bundles AssertWriterDeterministic
// with AllowUpdate. When trigger is false this is a fast no-op
// (returns false without invoking the writer at all) so non-update
// runs incur zero cost. When trigger is true, the writer must
// survive determinism BEFORE the gate is consulted — a flaky
// writer cannot regenerate fixtures, period.
func AllowUpdateAfterDeterminism(t *testing.T, trigger bool, label string, writer WriterFn) bool {
	t.Helper()
	if !trigger {
		return false
	}
	AssertWriterDeterministic(t, label, writer)
	return AllowUpdate(t, trigger)
}
