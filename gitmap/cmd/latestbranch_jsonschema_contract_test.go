package cmd

// JSON schema contract for `gitmap latest-branch --json`. Pairs the
// runtime encoder (encodeLatestBranchJSON / renderLatestBranchTopRaw
// in latestbranchrender.go) with the published schema at
// spec/08-json-schemas/latest-branch.schema.json so drift in either
// side fails the build.

import (
	"bytes"
	"sort"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/gitutil"
)

const latestBranchSchemaFilename = "latest-branch.schema.json"

// latestBranchTopLevelRequiredKeys mirrors the schema's required array.
var latestBranchTopLevelRequiredKeys = []string{
	"branch",
	"commitDate",
	"ref",
	"remote",
	"sha",
	"subject",
}

// TestLatestBranchJSONSchema_TopLevelShape pins the root type (object)
// and the required key set against the schema.
func TestLatestBranchJSONSchema_TopLevelShape(t *testing.T) {
	root := loadSchemaFile(t, latestBranchSchemaFilename)
	if root["type"] != "object" {
		t.Fatalf("top-level type = %v, want object", root["type"])
	}
	got := stringSliceFromAny(root["required"])
	sort.Strings(got)
	if !equalStringSlices(got, latestBranchTopLevelRequiredKeys) {
		t.Fatalf("required = %v, want %v", got, latestBranchTopLevelRequiredKeys)
	}
}

// TestLatestBranchJSONSchema_EncoderMatchesSchema runs the real
// stablejson encoder, then asserts every key in the emitted object
// is declared in the schema's properties map.
func TestLatestBranchJSONSchema_EncoderMatchesSchema(t *testing.T) {
	root := loadSchemaFile(t, latestBranchSchemaFilename)
	props, ok := root["properties"].(map[string]any)
	if !ok {
		t.Fatalf("schema missing properties object")
	}

	result := latestBranchResult{
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

	var buf bytes.Buffer
	if err := encodeLatestBranchJSON(&buf, result, nil, 0); err != nil {
		t.Fatalf("encode: %v", err)
	}
	gotKeys := readFirstObjectKeys(t, buf.Bytes())
	for _, key := range gotKeys {
		if _, allowed := props[key]; !allowed {
			t.Errorf("encoder emitted %q not declared in schema properties", key)
		}
	}
}
