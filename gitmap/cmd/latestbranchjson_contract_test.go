package cmd

// JSON schema contract tests for `gitmap latest-branch --json`.
//
// Same two-tier strictness as startuplistjson_contract_test.go:
// byte-exact for canonical fixtures, structural for variable data.
// Covers BOTH the top-level object without `top` AND the nested
// top-N array inside `top`, since the key is conditionally emitted.
//
// Regenerate fixtures with:
//
//   GITMAP_UPDATE_GOLDEN=1 go test ./cmd/ -run LatestBranchJSONContract

import (
	"bytes"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/gitutil"
)

// canonicalLatestResult builds a deterministic latestBranchResult
// for byte-exact tests. Every value is a fixed string so the golden
// file's bytes are reproducible across machines, timezones, and
// gitutil format-helper changes (we go around the helpers by
// pre-formatting the strings ourselves at the call site).
func canonicalLatestResult() latestBranchResult {

	return latestBranchResult{
		branchNames:    []string{"main", "develop"},
		selectedRemote: "origin",
		shortSha:       "abc1234",
		commitDate:     "01-Jan-2025 12:00 PM (UTC)",
		latest: gitutil.RemoteBranchInfo{
			RemoteRef:  "refs/remotes/origin/main",
			CommitDate: time.Unix(0, 0).UTC(),
			Sha:        "abc1234567890",
			Subject:    "Initial commit",
		},
	}
}

// TestLatestBranchJSONContract_NoTopOmitsKey pins the bytes when
// `top` is absent (top=0). The `omitempty` tag on Top means the key
// must NOT appear in the output.
func TestLatestBranchJSONContract_NoTopOmitsKey(t *testing.T) {
	encode := func() ([]byte, error) {
		var buf bytes.Buffer
		err := encodeLatestBranchJSON(&buf, canonicalLatestResult(), nil, 0)

		return buf.Bytes(), err
	}
	assertGoldenBytesDeterministic(t, "latest_branch_no_top.json", encode)
	// Schema check uses a fresh encode so a bug in the helper that
	// mutates the returned slice cannot cross-contaminate the two checks.
	raw, _ := encode()
	assertSchemaKeysFirstObject(t, raw, "latest-branch-no-top")
}

// TestLatestBranchJSONContract_WithTopIncludesKey verifies the
// `top` key is present and the nested object's key order matches
// the lbTopKey* constants in latestbranchrender.go. Structural-only
// (no byte-exact golden) because renderLatestBranchTopRaw calls
// gitutil.FormatDisplayDate which uses the LOCAL timezone — bytes
// would drift across CI machines.
func TestLatestBranchJSONContract_WithTopIncludesKey(t *testing.T) {
	items := []gitutil.RemoteBranchInfo{
		{
			RemoteRef:  "refs/remotes/origin/main",
			CommitDate: time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC),
			Sha:        "abc1234567890",
			Subject:    "Initial commit",
		},
	}
	var buf bytes.Buffer
	if err := encodeLatestBranchJSON(&buf, canonicalLatestResult(), items, 1); err != nil {
		t.Fatalf("encode: %v", err)
	}
	assertSchemaKeysFirstObject(t, buf.Bytes(), "latest-branch-with-top")
}
