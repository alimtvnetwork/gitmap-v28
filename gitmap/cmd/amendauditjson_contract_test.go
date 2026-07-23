package cmd

// JSON contract tests for `gitmap amend audit` file output.
//
// The contract covers key order of the emitted single object:
//   id, timestamp, branch, fromCommit, toCommit, totalCommits,
//   previousAuthor, newAuthor, mode, forcePushed, commits.
//
// Regenerate fixtures with:
//
//   GITMAP_UPDATE_GOLDEN=1 go test ./cmd/ -run AmendAuditJSONContract

import (
	"bytes"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// canonicalAmendmentRecord builds a deterministic single audit record.
func canonicalAmendmentRecord() model.AmendmentRecord {
	return model.AmendmentRecord{
		ID:           42,
		Timestamp:    "2025-01-01T12:00:00Z",
		Branch:       "main",
		FromCommit:   "abc123",
		ToCommit:     "def456",
		TotalCommits: 2,
		PreviousAuthor: model.AmendAuthor{
			Name:  "Old Name",
			Email: "old@example.com",
		},
		NewAuthor: model.AmendAuthor{
			Name:  "New Name",
			Email: "new@example.com",
		},
		Mode:        "all",
		ForcePushed: true,
		Commits: []model.CommitEntry{
			{SHA: "abc123", Message: "First"},
			{SHA: "def456", Message: "Second"},
		},
	}
}

// TestAmendAuditJSONContract_CanonicalRecord_KeyOrder asserts the
// key order of the emitted object matches the schema declaration.
func TestAmendAuditJSONContract_CanonicalRecord_KeyOrder(t *testing.T) {
	encode := func() ([]byte, error) {
		var buf bytes.Buffer
		err := encodeAmendAuditJSON(&buf, canonicalAmendmentRecord())

		return buf.Bytes(), err
	}
	assertGoldenBytesDeterministic(t, "amend_audit_canonical.json", encode)
	raw, _ := encode()
	assertSchemaKeysFirstObject(t, raw, "amend-audit")
}
