package cmd

// JSON contract tests for `gitmap watch --json`.
//
// watch emits a single top-level object with a nested repos array
// and summary object. The contract covers:
//
//   - Top-level key order: timestamp, repos, summary.
//   - Repo item key order: name, path, branch, status, ahead, behind, stash.
//   - Summary key order: total, dirty, behind, stash.
//
// Regenerate fixtures with:
//
//   GITMAP_UPDATE_GOLDEN=1 go test ./cmd/ -run WatchJSONContract

import (
	"bytes"
	"testing"
)

// TestWatchJSONContract_EmptyRepos asserts the top-level shape when
// no repos are present. The repos array must be `[]` and the summary
// must show all-zero counts.
func TestWatchJSONContract_EmptyRepos(t *testing.T) {
	assertGoldenBytesDeterministic(t, "watch_empty.json", func() ([]byte, error) {
		var buf bytes.Buffer
		err := encodeWatchJSON(&buf, nil, watchSummary{}, "2025-01-01T12:00:00Z")

		return buf.Bytes(), err
	})
}

// TestWatchJSONContract_CanonicalObject_KeyOrders asserts the key
// order of the top-level object, the first repo item, and the
// summary object.
func TestWatchJSONContract_CanonicalObject_KeyOrders(t *testing.T) {
	snapshots := []watchSnapshot{canonicalWatchSnapshot()}
	summary := buildWatchSummary(snapshots)
	var buf bytes.Buffer
	if err := encodeWatchJSON(&buf, snapshots, summary, "2025-01-01T12:00:00Z"); err != nil {
		t.Fatalf("encode: %v", err)
	}
	assertSchemaKeysFirstObject(t, buf.Bytes(), "watch")
}

// canonicalWatchSnapshot builds a deterministic repo snapshot for
// key-order and structural tests.
func canonicalWatchSnapshot() watchSnapshot {
	return watchSnapshot{
		Name:   "widget",
		Path:   "/repos/acme/widget",
		Branch: "main",
		Status: "clean",
		Ahead:  0,
		Behind: 2,
		Stash:  1,
	}
}
