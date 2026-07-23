package cmd

// JSON contract tests for `gitmap temp-releaselist --json`.
//
// temp-releaselist emits an array of temporary release branch records.
// The contract covers:
//
//   - Top-level array shape (empty must be `[]\n`).
//   - Key order: id, branch, versionPrefix, sequenceNumber,
//     commit, commitMessage, createdAt.
//
// Regenerate fixtures with:
//
//   GITMAP_UPDATE_GOLDEN=1 go test ./cmd/ -run TempReleaseListJSONContract

import (
	"bytes"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// canonicalTempReleaseList builds a deterministic two-row fixture.
func canonicalTempReleaseList(t *testing.T) []model.TempRelease {
	t.Helper()

	return []model.TempRelease{
		{
			ID:             1,
			Branch:         "release/v5.76.0-temp-3",
			VersionPrefix:  "v5.76",
			SequenceNumber: 3,
			CommitSha:      "abc123def456",
			CommitMessage:  "feat: add version-history stablejson encoder",
			CreatedAt:      "2026-05-26T10:00:00Z",
		},
		{
			ID:             2,
			Branch:         "release/v5.76.0-temp-4",
			VersionPrefix:  "v5.76",
			SequenceNumber: 4,
			CommitSha:      "def789abc012",
			CommitMessage:  "fix: correct edge case in history rewrite",
			CreatedAt:      "2026-05-26T11:30:00Z",
		},
	}
}

// TestTempReleaseListJSONContract_EmptyIsArrayNotNull is the jq-compat guarantee.
func TestTempReleaseListJSONContract_EmptyIsArrayNotNull(t *testing.T) {
	assertGoldenBytesDeterministic(t, "temp_release_list_empty.json", func() ([]byte, error) {
		var buf bytes.Buffer
		err := encodeTempReleaseListJSON(&buf, nil)

		return buf.Bytes(), err
	})
}

// TestTempReleaseListJSONContract_CanonicalRows pins the bytes of a
// canonical two-row sample, locking key order.
func TestTempReleaseListJSONContract_CanonicalRows(t *testing.T) {
	releases := canonicalTempReleaseList(t)
	assertGoldenBytesDeterministic(t, "temp_release_list_canonical.json", func() ([]byte, error) {
		var buf bytes.Buffer
		err := encodeTempReleaseListJSON(&buf, releases)

		return buf.Bytes(), err
	})
}

// TestTempReleaseListJSONContract_KeyOrder asserts the first object's
// key order matches the schema registry declaration.
func TestTempReleaseListJSONContract_KeyOrder(t *testing.T) {
	releases := canonicalTempReleaseList(t)
	var buf bytes.Buffer
	if err := encodeTempReleaseListJSON(&buf, releases); err != nil {
		t.Fatalf("encode: %v", err)
	}
	assertSchemaKeysFirstObject(t, buf.Bytes(), "temp-release-list")
}
