package cmd

// JSON schema contract tests for `gitmap startup-list --format=json`.
//
// Two strictness levels per the project policy on output stability:
//
//   1. Byte-exact golden comparison for an EMPTY list and a SINGLE
//      hand-authored canonical entry. These are the fixtures
//      downstream consumers most commonly bake into integration
//      tests, so the bytes themselves (including indentation, key
//      order, and trailing newline) are pinned.
//
//   2. Structural key-order check for a multi-entry list with mixed
//      content. This catches schema drift (renamed/added/reordered
//      fields) without being brittle on payload-shape variations.
//
// To intentionally regenerate the golden fixtures after a deliberate
// schema change, run:
//
//   GITMAP_UPDATE_GOLDEN=1 go test ./cmd/ -run StartupListJSONContract
//
// Then commit the updated files under cmd/testdata/ and bump the
// consumer-facing changelog.

import (
	"bytes"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/startup"
)

// TestStartupListJSONContract_EmptyIsArrayNotNull is the headline
// jq-compat guarantee: zero managed entries must encode as `[]\n`,
// never `null`, so `<output> | length` works without conditionals.
func TestStartupListJSONContract_EmptyIsArrayNotNull(t *testing.T) {
	assertGoldenBytesDeterministic(t, "startup_list_empty.json", func() ([]byte, error) {
		var buf bytes.Buffer
		err := encodeStartupListJSON(&buf, nil)

		return buf.Bytes(), err
	})
}

// TestStartupListJSONContract_EmptyNonNilSliceAlsoIsArray covers
// the second way "empty" reaches the encoder: a non-nil slice with
// length 0 (e.g., a filter that matched no rows but still allocated
// the result). Go's encoding/json renders a nil slice as `null` by
// default; only the explicit `[]Entry{}` case proves the encoder
// also normalizes the non-nil-but-empty case to `[]\n`. Without
// this test, a future refactor that drops the nil-normalization
// step would still pass EmptyIsArrayNotNull but silently break
// callers that pass a pre-allocated empty slice.
func TestStartupListJSONContract_EmptyNonNilSliceAlsoIsArray(t *testing.T) {
	assertGoldenBytesDeterministic(t, "startup_list_empty.json", func() ([]byte, error) {
		var buf bytes.Buffer
		err := encodeStartupListJSON(&buf, []startup.Entry{})

		return buf.Bytes(), err
	})
}

// TestStartupListJSONContract_CanonicalEntry pins the exact bytes
// for a known single entry. Any change to indentation, key order,
// tag names, or trailing-newline behavior breaks this test.
func TestStartupListJSONContract_CanonicalEntry(t *testing.T) {
	entries := []startup.Entry{
		{
			Name: "gitmap-sync-watcher",
			Path: "/home/user/.config/autostart/gitmap-sync-watcher.desktop",
			Exec: "/usr/local/bin/gitmap watch ~/projects",
		},
	}
	assertGoldenBytesDeterministic(t, "startup_list_single.json", func() ([]byte, error) {
		var buf bytes.Buffer
		err := encodeStartupListJSON(&buf, entries)

		return buf.Bytes(), err
	})
}

// TestStartupListJSONContract_KeyOrderStable validates that each
// element in a multi-entry list has its keys in the declared order
// (name, path, exec). Uses the structural check so the test stays
// robust against value-shape changes that do NOT affect the schema.
func TestStartupListJSONContract_KeyOrderStable(t *testing.T) {
	entries := []startup.Entry{
		{Name: "gitmap-a", Path: "/p/a.desktop", Exec: "/bin/a"},
		{Name: "gitmap-b", Path: "/p/b.desktop", Exec: "/bin/b --flag"},
		{Name: "gitmap-c", Path: "/p/c.desktop", Exec: ""},
	}
	var buf bytes.Buffer
	if err := encodeStartupListJSON(&buf, entries); err != nil {
		t.Fatalf("encode: %v", err)
	}
	assertSchemaKeysFirstObject(t, buf.Bytes(), "startup-list")
}
