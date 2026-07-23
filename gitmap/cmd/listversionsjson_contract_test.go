package cmd

// JSON contract tests for `gitmap list-versions --json`.
//
// list-versions emits an array of {version, source?, changelog?}.
// The contract covers:
//
//   - Top-level array shape (empty must be `[]\n`).
//   - Key order: version, source, changelog.
//   - source/changelog are omitted when empty (legacy omitempty shape).
//
// Regenerate fixtures with:
//
//   GITMAP_UPDATE_GOLDEN=1 go test ./cmd/ -run ListVersionsJSONContract

import (
	"bytes"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/release"
)

// TestListVersionsJSONContract_EmptyIsArrayNotNull is the jq-compat guarantee.
func TestListVersionsJSONContract_EmptyIsArrayNotNull(t *testing.T) {
	assertGoldenBytesDeterministic(t, "list_versions_empty.json", func() ([]byte, error) {
		var buf bytes.Buffer
		err := encodeListVersionsJSON(&buf, nil)

		return buf.Bytes(), err
	})
}

// canonicalListVersionsEntries builds a deterministic two-row fixture:
// the first row exercises all fields; the second omits both optional
// fields so the omitempty contract is pinned.
func canonicalListVersionsEntries(t *testing.T) []versionEntry {
	t.Helper()
	v1, err := release.Parse("v5.72.0")
	if err != nil {
		t.Fatalf("parse v5.72.0: %v", err)
	}
	v2, err := release.Parse("v5.71.0")
	if err != nil {
		t.Fatalf("parse v5.71.0: %v", err)
	}

	return []versionEntry{
		{Version: v1, Source: "local", Notes: []string{"first note", "second note"}},
		{Version: v2},
	}
}

// TestListVersionsJSONContract_CanonicalRows pins the bytes of a
// canonical two-row sample, locking key order AND the omitempty
// behavior for source/changelog.
func TestListVersionsJSONContract_CanonicalRows(t *testing.T) {
	entries := canonicalListVersionsEntries(t)
	assertGoldenBytesDeterministic(t, "list_versions_canonical.json", func() ([]byte, error) {
		var buf bytes.Buffer
		err := encodeListVersionsJSON(&buf, entries)

		return buf.Bytes(), err
	})
}

// TestListVersionsJSONContract_KeyOrder asserts the first object's
// key order matches the schema registry declaration.
func TestListVersionsJSONContract_KeyOrder(t *testing.T) {
	entries := canonicalListVersionsEntries(t)
	var buf bytes.Buffer
	if err := encodeListVersionsJSON(&buf, entries); err != nil {
		t.Fatalf("encode: %v", err)
	}
	assertSchemaKeysFirstObject(t, buf.Bytes(), "list-versions")
}
