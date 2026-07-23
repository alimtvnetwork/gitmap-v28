package cmd

// JSON schema contract tests for `gitmap find-next --json`.
//
// find-next emits an array of model.FindNextRow, each containing an
// embedded model.ScanRecord under the `repo` key. The contract
// covers BOTH layers:
//
//   - Top-level array shape (empty must be `[]\n`).
//   - FindNextRow key order: repo, nextVersionTag, nextVersionNum,
//     method, probedAt.
//   - ScanRecord key order (nested under `repo`): id, slug,
//     repoName, httpsUrl, sshUrl, branch, branchSource,
//     relativePath, absolutePath, cloneInstruction, notes, depth.
//
// ScanRecord is large (12 fields) and shared with several other CLI
// outputs, so a rename/reorder there would silently ripple into
// multiple consumers. Pinning it here adds one tripwire that catches
// every such change.
//
// Regenerate fixtures with:
//
//   GITMAP_UPDATE_GOLDEN=1 go test ./cmd/ -run FindNextJSONContract

import (
	"bytes"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// TestFindNextJSONContract_EmptyIsArrayNotNull is the jq-compat
// guarantee: zero rows must encode as `[]\n` even when the input
// slice is nil.
func TestFindNextJSONContract_EmptyIsArrayNotNull(t *testing.T) {
	assertGoldenBytesDeterministic(t, "find_next_empty.json", func() ([]byte, error) {
		var buf bytes.Buffer
		err := encodeFindNextJSON(&buf, nil)

		return buf.Bytes(), err
	})
}

// canonicalFindNextRow builds a deterministic single row whose every
// field is a fixed value, so the golden file's bytes are stable
// across machines and time. Used for the byte-exact + key-order
// tests below.
func canonicalFindNextRow() model.FindNextRow {

	return model.FindNextRow{
		Repo: model.ScanRecord{
			ID:               42,
			Slug:             "acme/widget",
			RepoName:         "widget",
			HTTPSUrl:         "https://github.com/acme/widget.git",
			SSHUrl:           "git@github.com:acme/widget.git",
			Branch:           "main",
			BranchSource:     "remote-head",
			RelativePath:     "acme/widget",
			AbsolutePath:     "/repos/acme/widget",
			CloneInstruction: "git clone https://github.com/acme/widget.git",
			Notes:            "",
		},
		NextVersionTag: "v1.2.3",
		NextVersionNum: 123,
		Method:         "tag-probe",
		ProbedAt:       "2025-01-01T12:00:00Z",
	}
}

// TestFindNextJSONContract_CanonicalRow_KeyOrders asserts the
// FindNextRow-level AND nested ScanRecord-level key orders.
// Structural-only (no byte-exact golden for the populated row) so
// the test stays robust against future numeric formatting changes
// in encoding/json or value-shape tweaks.
func TestFindNextJSONContract_CanonicalRow_KeyOrders(t *testing.T) {
	rows := []model.FindNextRow{canonicalFindNextRow()}
	var buf bytes.Buffer
	if err := encodeFindNextJSON(&buf, rows); err != nil {
		t.Fatalf("encode: %v", err)
	}
	assertSchemaKeysFirstObject(t, buf.Bytes(), "find-next")
}
