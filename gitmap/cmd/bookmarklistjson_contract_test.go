package cmd

// JSON contract tests for `gitmap bookmark list --json`.
//
// bookmark-list emits an array of model.BookmarkRecord. The contract
// covers:
//
//   - Top-level array shape (empty must be `[]\n`).
//   - Key order: id, name, command, args, flags, createdAt.
//
// Regenerate fixtures with:
//
//   GITMAP_UPDATE_GOLDEN=1 go test ./cmd/ -run BookmarkListJSONContract

import (
	"bytes"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// TestBookmarkListJSONContract_EmptyIsArrayNotNull is the jq-compat guarantee.
func TestBookmarkListJSONContract_EmptyIsArrayNotNull(t *testing.T) {
	assertGoldenBytesDeterministic(t, "bookmark_list_empty.json", func() ([]byte, error) {
		var buf bytes.Buffer
		err := encodeBookmarkListJSON(&buf, nil)

		return buf.Bytes(), err
	})
}

// canonicalBookmarkRecord builds a deterministic single row.
func canonicalBookmarkRecord() model.BookmarkRecord {
	return model.BookmarkRecord{
		ID:        7,
		Name:      "my-bookmark",
		Command:   "scan",
		Args:      "--depth=3",
		Flags:     "--verbose",
		CreatedAt: "2025-01-01T12:00:00Z",
	}
}

// TestBookmarkListJSONContract_CanonicalRow_KeyOrder asserts the key
// order of the emitted object matches the schema declaration.
func TestBookmarkListJSONContract_CanonicalRow_KeyOrder(t *testing.T) {
	records := []model.BookmarkRecord{canonicalBookmarkRecord()}
	var buf bytes.Buffer
	if err := encodeBookmarkListJSON(&buf, records); err != nil {
		t.Fatalf("encode: %v", err)
	}
	assertSchemaKeysFirstObject(t, buf.Bytes(), "bookmark-list")
}
