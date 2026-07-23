package cmd

// Reusable helpers for JSON snapshot tests (extracted from
// startuplistjson_snapshot_test.go to keep both files under the
// 200-line code-style budget). Future `--format=json` snapshot
// tests should land their helpers here too rather than spawning
// per-feature helper files.
//
// Design contract:
//
//   - Token-stream parsing (NOT json.Unmarshal into map[string]any),
//     so on-the-wire key ORDER is observable. A map-based parse
//     would silently lose the very property these tests exist to
//     pin down.
//
//   - Per-object failure messages name the offending object index
//     and the missing/unexpected/reordered key. A schema regression
//     in row 17 of a 100-row output points at row 17, not at the
//     whole blob.
//
//   - No dependency on any specific encoder or schema — these
//     helpers operate on `[]byte` of already-rendered JSON, so they
//     work for stablejson, encoding/json, or any future encoder.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/startup"
)

// mustEncodeStartupList runs encodeStartupListJSON on `entries` and
// returns a COPY of the produced bytes (so a caller comparing two
// runs cannot accidentally alias a reused buffer). Lives in the
// shared helpers file rather than per-test so any future
// startup-list-shaped snapshot can reuse it without duplicating
// the encoder-call + buffer-copy boilerplate.
func mustEncodeStartupList(t *testing.T, entries []startup.Entry) []byte {
	t.Helper()
	var buf bytes.Buffer
	if err := encodeStartupListJSON(&buf, entries); err != nil {
		t.Fatalf("encode: %v", err)
	}

	return append([]byte(nil), buf.Bytes()...)
}

// Per-object key-set assertions (assertEveryObjectKeysExact /
// assertObjectKeysExactAt) used to live here but were dropped once
// no test referenced them — the live snapshot suites pin shape via
// assertGoldenBytes instead. Re-introduce alongside their first
// caller if a future test needs structural-only checks.

// readEveryObjectKeys streams the entire top-level array and returns
// one []string of keys per object. Returns an empty outer slice for
// `[]` so callers can distinguish "no objects" from "malformed".
func readEveryObjectKeys(t *testing.T, raw []byte) [][]string {
	t.Helper()
	dec := json.NewDecoder(bytes.NewReader(raw))
	if err := expectDelim(dec, '['); err != nil {
		t.Fatalf("expected top-level array: %v", err)
	}
	var out [][]string
	for dec.More() {
		if err := expectDelim(dec, '{'); err != nil {
			t.Fatalf("expected object at index %d: %v", len(out), err)
		}
		out = append(out, collectObjectKeys(t, dec))
	}

	return out
}

// expectDelim reads the next token and confirms it's the requested
// JSON delimiter (`[`, `]`, `{`, or `}`). One helper instead of
// per-delimiter wrappers keeps the helpers file under budget while
// keeping failure messages specific (the caller adds context).
func expectDelim(dec *json.Decoder, want byte) error {
	tok, err := dec.Token()
	if err != nil {

		return err
	}
	delim, ok := tok.(json.Delim)
	if !ok || rune(delim) != rune(want) {

		return fmt.Errorf("want delim %q, got %v (%T)", want, tok, tok)
	}

	return nil
}

// collectObjectKeys reads key-value pairs from dec until the
// closing `}` is consumed, returning just the key names in
// wire order. dec must already have consumed the opening `{`.
func collectObjectKeys(t *testing.T, dec *json.Decoder) []string {
	t.Helper()
	var keys []string
	for dec.More() {
		tok, err := dec.Token()
		if err != nil {
			t.Fatalf("reading object key: %v", err)
		}
		key, ok := tok.(string)
		if !ok {
			t.Fatalf("expected string key, got %v (%T)", tok, tok)
		}
		keys = append(keys, key)
		// Skip the value without decoding its type.
		var raw json.RawMessage
		if err := dec.Decode(&raw); err != nil {
			t.Fatalf("skipping value for key %q: %v", key, err)
		}
	}
	// Consume the closing '}'.
	if _, err := dec.Token(); err != nil {
		t.Fatalf("expected closing '}': %v", err)
	}
	return keys
}

// equalStringSlices returns true when a and b have identical length
// and identical element values at every index.
func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// stringSet (former map-based set helper) was removed alongside the
// per-object key-set assertions that used it. Re-introduce it next
// to its first caller if a future helper needs O(1) membership.
