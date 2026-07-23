package cmd

// JSON contract tests for `gitmap probe --json`.
//
// probe emits an array of probeJSONEntry. The contract covers:
//
//   - Top-level array shape (empty must be `[]\n`).
//   - Key order: repoId, slug, absolutePath, nextVersionTag,
//     nextVersionNum, method, isAvailable, error.
//
// Regenerate fixtures with:
//
//   GITMAP_UPDATE_GOLDEN=1 go test ./cmd/ -run ProbeJSONContract

import (
	"bytes"
	"testing"
)

// TestProbeJSONContract_EmptyIsArrayNotNull is the jq-compat
// guarantee: zero rows must encode as `[]\n` even when the input
// slice is nil.
func TestProbeJSONContract_EmptyIsArrayNotNull(t *testing.T) {
	assertGoldenBytesDeterministic(t, "probe_report_empty.json", func() ([]byte, error) {
		var buf bytes.Buffer
		err := encodeProbeJSON(&buf, nil)

		return buf.Bytes(), err
	})
}

// canonicalProbeEntry builds a deterministic single entry whose
// every field is a fixed value, so the golden file's bytes are
// stable across machines and time.
func canonicalProbeEntry() probeJSONEntry {
	return probeJSONEntry{
		RepoID:         42,
		Slug:           "acme/widget",
		AbsolutePath:   "/repos/acme/widget",
		NextVersionTag: "v1.2.3",
		NextVersionNum: 123,
		Method:         "tag-probe",
		IsAvailable:    true,
		Error:          "",
	}
}

// TestProbeJSONContract_CanonicalRow_KeyOrders asserts the key order
// of the emitted object matches the schema declaration.
// Structural-only (no byte-exact golden for the populated row) so
// the test stays robust against future numeric formatting changes
// in encoding/json or value-shape tweaks.
func TestProbeJSONContract_CanonicalRow_KeyOrders(t *testing.T) {
	entries := []probeJSONEntry{canonicalProbeEntry()}
	var buf bytes.Buffer
	if err := encodeProbeJSON(&buf, entries); err != nil {
		t.Fatalf("encode: %v", err)
	}
	assertSchemaKeysFirstObject(t, buf.Bytes(), "probe-report")
}
