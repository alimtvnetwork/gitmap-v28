package cmd

// Regression contract: `gitmap startup-list --format=json` must
// produce byte-identical output across two back-to-back executions
// over the same input. This is the determinism guarantee downstream
// consumers actually depend on:
//
//   - CI diff jobs that re-run a command and `diff` the outputs
//   - Content-addressed caches keyed on the output hash
//   - "Did anything change?" pipelines using sha256 of the JSON
//
// Verified at the encoder boundary (encodeStartupListJSON), NOT by
// shelling out to the built CLI binary. The CLI dispatcher is a
// thin wrapper around this encoder + os.Stdout — testing the
// encoder isolates the determinism question from filesystem state
// (XDG autostart dir contents change between machines / runs and
// would make a binary-shelling test flaky for unrelated reasons).
//
// Pinned across THREE shapes so a future regression caused by e.g.
// random map iteration in a helper would be caught regardless of
// input size:
//
//   1. Empty list           — exercises the `[]\n` short-circuit
//   2. Single entry         — minimal happy path
//   3. Multi-entry with     — exercises the inter-element separator
//      varied content         path AND special-character escaping
//                             (the two places stdlib map iteration
//                             order has historically leaked through)
//
// Also pinned across BOTH indent settings the flag exposes (0 and
// the default 2) so a determinism regression in either the minified
// or pretty path fails this test, not just the one most callers use.

import (
	"bytes"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/startup"
)

// TestStartupListJSON_DeterministicAcrossRuns runs the encoder
// twice for each fixture and asserts byte-equal output. A failure
// here means a downstream content-addressed cache would generate
// spurious cache misses on every gitmap invocation.
func TestStartupListJSON_DeterministicAcrossRuns(t *testing.T) {
	fixtures := []struct {
		name    string
		entries []startup.Entry
	}{
		{
			name:    "empty",
			entries: nil,
		},
		{
			name: "single",
			entries: []startup.Entry{
				{Name: "gitmap-a", Path: "/p/a.desktop", Exec: "/bin/a"},
			},
		},
		{
			name: "multi_with_special_chars",
			entries: []startup.Entry{
				{Name: "gitmap-a", Path: "/p/a.desktop", Exec: "/bin/a"},
				{Name: "gitmap-中文", Path: "/p/\"q\"\\b.desktop", Exec: "a\tb\nc"},
				{Name: "gitmap-c", Path: "/p/c.desktop", Exec: ""},
			},
		},
	}
	for _, fx := range fixtures {
		t.Run(fx.name, func(t *testing.T) {
			assertEncoderDeterministic(t, fx.entries)
		})
	}
}

// assertEncoderDeterministic encodes `entries` twice via the
// default-indent path AND the minified path and asserts both pairs
// are byte-equal. Two separate run pairs (rather than one combined
// run-3-times sweep) keep the failure message specific: the test
// name reports which indent setting drifted.
func assertEncoderDeterministic(t *testing.T, entries []startup.Entry) {
	t.Helper()
	t.Run("indent_default", func(t *testing.T) {
		first := encodeOnceDefault(t, entries)
		second := encodeOnceDefault(t, entries)
		assertBytesEqualOrDiff(t, "default-indent", first, second)
	})
	t.Run("indent_minified", func(t *testing.T) {
		first := encodeOnceIndent(t, entries, 0)
		second := encodeOnceIndent(t, entries, 0)
		assertBytesEqualOrDiff(t, "minified", first, second)
	})
}

// encodeOnceDefault runs the legacy 2-arg wrapper (indent=2). Kept
// separate from encodeOnceIndent so the determinism test verifies
// the wrapper's stability, not just the underlying indented path —
// a future change that adds e.g. a timestamp to the wrapper alone
// would still fail this test.
func encodeOnceDefault(t *testing.T, entries []startup.Entry) []byte {
	t.Helper()
	var buf bytes.Buffer
	if err := encodeStartupListJSON(&buf, entries); err != nil {
		t.Fatalf("encode: %v", err)
	}

	return buf.Bytes()
}

// encodeOnceIndent runs the indent-aware encoder at the requested
// width. Returns a fresh []byte each call so the comparison cannot
// accidentally pass via aliased buffer reuse.
func encodeOnceIndent(t *testing.T, entries []startup.Entry, indent int) []byte {
	t.Helper()
	var buf bytes.Buffer
	if err := encodeStartupListJSONIndent(&buf, entries, indent); err != nil {
		t.Fatalf("encode (indent=%d): %v", indent, err)
	}
	return buf.Bytes()
}

// assertBytesEqualOrDiff prints both runs verbatim on mismatch so
// the developer sees exactly which bytes drifted (length difference,
// whitespace flip, escape sequence change). Plain bytes.Equal would
// just say "false" and force a manual reproduction.
func assertBytesEqualOrDiff(t *testing.T, label string, first, second []byte) {
	t.Helper()
	if bytes.Equal(first, second) {
		return
	}
	t.Fatalf("%s: byte drift between runs\n--- run 1 (%d bytes)\n%s--- run 2 (%d bytes)\n%s",
		label, len(first), first, len(second), second)
}
