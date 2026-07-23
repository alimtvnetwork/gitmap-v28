package cmd

// JSON contract tests for `gitmap diff-profiles --json`.
//
// The contract covers key order of the emitted single object:
//   profileA, profileB, onlyInA, onlyInB, different, same.
//
// Regenerate fixtures with:
//
//   GITMAP_UPDATE_GOLDEN=1 go test ./cmd/ -run DiffProfilesJSONContract

import (
	"bytes"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// canonicalDPResult builds a deterministic diff-profiles result.
func canonicalDPResult() dpResult {
	return dpResult{
		onlyInA: []model.ScanRecord{
			{RepoName: "repo-a", AbsolutePath: "/home/user/a"},
		},
		onlyInB: []model.ScanRecord{
			{RepoName: "repo-b", AbsolutePath: "/home/user/b"},
		},
		different: []dpDiff{
			{
				Name:  "repo-diff",
				PathA: "/home/user/diff-a",
				PathB: "/home/user/diff-b",
				ModeA: "https",
				ModeB: "ssh",
			},
		},
		same: []model.ScanRecord{
			{RepoName: "repo-same", AbsolutePath: "/home/user/same"},
		},
	}
}

// TestDiffProfilesJSONContract_CanonicalRecord_KeyOrder asserts
// the key order of the emitted object matches the schema declaration.
func TestDiffProfilesJSONContract_CanonicalRecord_KeyOrder(t *testing.T) {
	encode := func() ([]byte, error) {
		var buf bytes.Buffer
		err := encodeDiffProfilesJSON(&buf, "alpha", "beta", canonicalDPResult())

		return buf.Bytes(), err
	}
	assertGoldenBytesDeterministic(t, "diff_profiles_canonical.json", encode)
	raw, _ := encode()
	assertSchemaKeysFirstObject(t, raw, "diff-profiles")
}
