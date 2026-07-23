package cmd

// JSON contract tests for `gitmap amend list --json`.
//
// amend-list emits an array of store.AmendmentRow. The contract
// covers:
//
//   - Top-level array shape (empty must be `[]\n`).
//   - Key order: ID, Branch, FromCommit, ToCommit, TotalCommits,
//     PreviousName, PreviousEmail, NewName, NewEmail, Mode,
//     ForcePushed, CreatedAt.
//
// Regenerate fixtures with:
//
//   GITMAP_UPDATE_GOLDEN=1 go test ./cmd/ -run AmendListJSONContract

import (
	"bytes"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/store"
)

// TestAmendListJSONContract_EmptyIsArrayNotNull is the jq-compat
// guarantee.
func TestAmendListJSONContract_EmptyIsArrayNotNull(t *testing.T) {
	assertGoldenBytesDeterministic(t, "amend_list_empty.json", func() ([]byte, error) {
		var buf bytes.Buffer
		err := encodeAmendListJSON(&buf, nil)

		return buf.Bytes(), err
	})
}

// canonicalAmendmentRow builds a deterministic single row.
func canonicalAmendmentRow() store.AmendmentRow {
	return store.AmendmentRow{
		ID:            7,
		Branch:        "main",
		FromCommit:    "a1b2c3d",
		ToCommit:      "e4f5g6h",
		TotalCommits:  3,
		PreviousName:  "Alice Old",
		PreviousEmail: "alice.old@example.com",
		NewName:       "Alice New",
		NewEmail:      "alice.new@example.com",
		Mode:          "rewrite",
		ForcePushed:   1,
		CreatedAt:     "2025-01-01T12:00:00Z",
	}
}

// TestAmendListJSONContract_CanonicalRow_KeyOrders asserts the key
// order of the emitted object matches the schema declaration.
func TestAmendListJSONContract_CanonicalRow_KeyOrders(t *testing.T) {
	rows := []store.AmendmentRow{canonicalAmendmentRow()}
	var buf bytes.Buffer
	if err := encodeAmendListJSON(&buf, rows); err != nil {
		t.Fatalf("encode: %v", err)
	}
	assertSchemaKeysFirstObject(t, buf.Bytes(), "amend-list")
}
