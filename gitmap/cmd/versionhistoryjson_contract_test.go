package cmd

// JSON contract tests for `gitmap version-history --json`.
//
// version-history emits an array of version transition records.
// The contract covers:
//
//   - Top-level array shape (empty must be `[]\n`).
//   - Key order: fromVersionTag, fromVersionNum, toVersionTag,
//     toVersionNum, flattenedPath?, createdAt?, id, repoId.
//   - flattenedPath and createdAt are omitted when empty (legacy
//     omitempty shape).
//
// Regenerate fixtures with:
//
//   GITMAP_UPDATE_GOLDEN=1 go test ./cmd/ -run VersionHistoryJSONContract

import (
	"bytes"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// canonicalVersionHistoryRecords builds a deterministic two-row fixture:
// the first row exercises all fields; the second omits both optional
// fields so the omitempty contract is pinned.
func canonicalVersionHistoryRecords(t *testing.T) []model.RepoVersionHistoryRecord {
	t.Helper()

	return []model.RepoVersionHistoryRecord{
		{
			ID:             1,
			RepoID:         42,
			FromVersionTag: "v5.40.0",
			FromVersionNum: 50400,
			ToVersionTag:   "v5.41.0",
			ToVersionNum:   50410,
			FlattenedPath:  "/home/user/repos/gitmap-v28",
			CreatedAt:      "2026-05-20T14:30:00Z",
		},
		{
			ID:             2,
			RepoID:         42,
			FromVersionTag: "v5.41.0",
			FromVersionNum: 50410,
			ToVersionTag:   "v5.42.0",
			ToVersionNum:   50420,
		},
	}
}

// TestVersionHistoryJSONContract_EmptyIsArrayNotNull is the jq-compat guarantee.
func TestVersionHistoryJSONContract_EmptyIsArrayNotNull(t *testing.T) {
	assertGoldenBytesDeterministic(t, "version_history_empty.json", func() ([]byte, error) {
		var buf bytes.Buffer
		err := encodeVersionHistoryJSON(&buf, nil)

		return buf.Bytes(), err
	})
}

// TestVersionHistoryJSONContract_CanonicalRows pins the bytes of a
// canonical two-row sample, locking key order AND the omitempty
// behavior for flattenedPath/createdAt.
func TestVersionHistoryJSONContract_CanonicalRows(t *testing.T) {
	records := canonicalVersionHistoryRecords(t)
	assertGoldenBytesDeterministic(t, "version_history_canonical.json", func() ([]byte, error) {
		var buf bytes.Buffer
		err := encodeVersionHistoryJSON(&buf, records)

		return buf.Bytes(), err
	})
}

// TestVersionHistoryJSONContract_KeyOrder asserts the first object's
// key order matches the schema registry declaration.
func TestVersionHistoryJSONContract_KeyOrder(t *testing.T) {
	records := canonicalVersionHistoryRecords(t)
	var buf bytes.Buffer
	if err := encodeVersionHistoryJSON(&buf, records); err != nil {
		t.Fatalf("encode: %v", err)
	}
	assertSchemaKeysFirstObject(t, buf.Bytes(), "version-history")
}
