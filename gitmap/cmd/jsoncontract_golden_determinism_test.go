package cmd

// Determinism overlay for golden-fixture assertions. Existing
// `assertGoldenBytes(t, name, buf.Bytes())` calls run the encoder
// ONCE and compare the result to disk — that proves the bytes
// match the snapshot but says nothing about whether a SECOND
// encode of the same input would produce the same bytes.
//
// Determinism is the guarantee downstream consumers actually
// depend on (CI diff jobs, content-addressed caches, "did anything
// change?" pipelines). The dedicated startuplistjson_determinism_
// test.go covers it for startup-list at the encoder boundary;
// this helper extends the same coverage to EVERY golden-fixture
// site (find-next, latest-branch, startup-list bytes/contract
// suites) without each test having to handroll the loop.
//
// Usage at a call site:
//
//	assertGoldenBytesDeterministic(t, "startup_list_multi.json",
//	    func() ([]byte, error) {
//	        var buf bytes.Buffer
//	        err := encodeStartupListJSON(&buf, entries)
//	        return buf.Bytes(), err
//	    })
//
// The closure isolates the encoder identity (which encode* call,
// which arguments) so the helper stays generic across schemas.
// On a determinism break, the failure message names which run
// diverged AND prints both byte streams so the offending field is
// obvious without re-running anything.

import (
	"bytes"
	"testing"
)

// determinismRuns is the number of repeat encodings compared
// against run 0. Three is the smallest count that catches both
// "second call diverges" (e.g. mutating package-level state) and
// "third call converges" (e.g. cache-warmup masking a real bug)
// without slowing the test suite measurably — each run is a
// pure in-memory encode of <100 bytes.
const determinismRuns = 3

// assertGoldenBytesDeterministic runs `encode` determinismRuns
// times, asserts every run's bytes match run 0, then delegates to
// assertGoldenBytes for the on-disk snapshot check. A failure on
// the cross-run comparison is reported BEFORE the golden check so
// the developer sees the determinism break first (which is almost
// always the root cause when the golden also drifts).
func assertGoldenBytesDeterministic(t *testing.T, name string, encode func() ([]byte, error)) {
	t.Helper()
	first, err := encode()
	if err != nil {
		t.Fatalf("%s: encode run 0: %v", name, err)
	}
	for i := 1; i < determinismRuns; i++ {
		got, err := encode()
		if err != nil {
			t.Fatalf("%s: encode run %d: %v", name, i, err)
		}
		if !bytes.Equal(got, first) {
			t.Fatalf("%s: determinism broken — run %d differs from run 0\n--- run 0 (%d bytes)\n%s--- run %d (%d bytes)\n%s",
				name, i, len(first), string(first), i, len(got), string(got))
		}
	}
	assertGoldenBytes(t, name, first)
}
