package cmd

// Schema-snapshot test for `gitmap startup-list --format=json`.
//
// This test complements (does NOT replace) the byte-exact golden
// fixtures in startuplistjson_contract_test.go. The goldens catch
// ANY change with a single opaque "bytes differ" message. This file
// catches the three SPECIFIC schema regressions that downstream
// consumers care about, with a targeted failure message for each:
//
//   1. Added field   → "unexpected key X in object[i]"
//   2. Removed field → "missing key X in object[i]"
//   3. Renamed field → both of the above fire on the same object
//   4. Reordered field → "key order drift in object[i]: want ... got ..."
//   5. Non-deterministic encoding → "determinism broken: run N differs"
//
// Why both layers? The golden test is the source of truth for the
// EXACT bytes (catches indentation / trailing-newline drift that no
// schema check would notice). This snapshot test is the source of
// truth for the SCHEMA SHAPE (catches what changed and where, in
// language a downstream JSON consumer would understand). When both
// fail, the goldens get regenerated; when only this fails, the
// schema itself drifted and the consumer-facing changelog must be
// bumped before regenerating.

import (
	"bytes"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/startup"
)

// Schema for `gitmap startup-list --format=json|jsonl|csv` rows is
// now sourced from the versioned schema registry (see
// cmd/testdata/schemas/startup-list.vN.json + assertSchemaKeysArray
// in schemaregistry_assert_test.go).
//
// To add a key (e.g. `enabled`): create startup-list.v2.json with
// the new key list, OR run `GITMAP_UPDATE_SCHEMA=startup-list go
// test ./cmd/...` once and acknowledge with `-accept-schema=
// startup-list@v2` thereafter.

// TestStartupListJSONSnapshot_SchemaIsLocked is the headline schema
// guarantee. Encodes a representative multi-entry list, then asserts
// every object has EXACTLY the expected keys in EXACTLY the
// expected order. Added / removed / renamed / reordered fields all
// produce a targeted failure message naming the offending object
// index and key.
func TestStartupListJSONSnapshot_SchemaIsLocked(t *testing.T) {
	entries := []startup.Entry{
		{Name: "gitmap-a", Path: "/p/a.desktop", Exec: "/bin/a"},
		{Name: "gitmap-b", Path: "/p/b.desktop", Exec: "/bin/b --flag"},
		// Empty Exec is intentional — past regressions have come
		// from "helpful" omitempty tags hiding zero-value fields.
		// Including this row guarantees the schema check sees the
		// key whether or not the value is the zero string.
		{Name: "gitmap-c", Path: "/p/c.desktop", Exec: ""},
	}
	var buf bytes.Buffer
	if err := encodeStartupListJSON(&buf, entries); err != nil {
		t.Fatalf("encode: %v", err)
	}
	assertSchemaKeysArray(t, buf.Bytes(), "startup-list")
}

// TestStartupListJSONSnapshot_DeterministicAcrossRuns proves the
// encoder produces byte-identical output when called repeatedly
// with the same input. The most common way determinism breaks is a
// map being introduced into the encoding path (Go map iteration is
// randomized). Running 32 times gives ~1 in 4 billion odds of
// missing a 50/50 ordering bug — high enough confidence for a unit
// test without making the suite slow.
func TestStartupListJSONSnapshot_DeterministicAcrossRuns(t *testing.T) {
	entries := []startup.Entry{
		{Name: "gitmap-x", Path: "/p/x.desktop", Exec: "/bin/x"},
		{Name: "gitmap-y", Path: "/p/y.desktop", Exec: "/bin/y"},
		{Name: "gitmap-z", Path: "/p/z.desktop", Exec: ""},
	}
	const runs = 32
	first := mustEncodeStartupList(t, entries)
	for i := 1; i < runs; i++ {
		got := mustEncodeStartupList(t, entries)
		if !bytes.Equal(first, got) {
			t.Fatalf("determinism broken: run %d differs from run 0\n--- run 0\n%s--- run %d\n%s",
				i, string(first), i, string(got))
		}
	}
}

// TestStartupListJSONSnapshot_EmptyListHasNoObjects guards the edge
// case where the schema check has no objects to inspect — empty
// list MUST still encode as `[]` (verified byte-exactly by the
// golden test), and the schema check's "every object" quantifier
// must not vacuously pass on something malformed.
func TestStartupListJSONSnapshot_EmptyListHasNoObjects(t *testing.T) {
	var buf bytes.Buffer
	if err := encodeStartupListJSON(&buf, nil); err != nil {
		t.Fatalf("encode: %v", err)
	}
	keysPerObject := readEveryObjectKeys(t, buf.Bytes())
	if len(keysPerObject) != 0 {
		t.Fatalf("empty list must produce zero objects, got %d: %v",
			len(keysPerObject), keysPerObject)
	}
}

// All shared JSON parsing + encoding helpers (mustEncodeStartupList,
// readEveryObjectKeys, assertEveryObjectKeysExact, etc.) live in
// jsonsnapshot_helpers_test.go so this file contains only tests.
