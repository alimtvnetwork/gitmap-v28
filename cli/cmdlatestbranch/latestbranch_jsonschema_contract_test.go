package cmdlatestbranch

// JSON schema contract for `gitmap latest-branch --json`. Pairs the
// runtime encoder (encodeLatestBranchJSON / renderLatestBranchTopRaw
// in latestbranchrender.go) with the published schema at
// 02-spec/08-json-schemas/latest-branch.schema.json so drift in either
// side fails the build.
//
// NOTE: this test lives in cmdlatestbranch (not cmd) because the
// encoder (encodeLatestBranchJSON) and the latestBranchResult struct
// are unexported here. Schema helpers are local, minimal equivalents
// of the cmd package's schema-test infrastructure.

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
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
	root := lbLoadSchemaFile(t, latestBranchSchemaFilename)
	if root["type"] != "object" {
		t.Fatalf("top-level type = %v, want object", root["type"])
	}

	got := lbStringSliceFromAny(root["required"])
	sort.Strings(got)
	if !lbEqualStringSlices(got, latestBranchTopLevelRequiredKeys) {
		t.Fatalf("required = %v, want %v", got, latestBranchTopLevelRequiredKeys)
	}
}

// TestLatestBranchJSONSchema_EncoderMatchesSchema runs the real
// stablejson encoder, then asserts every key in the emitted object
// is declared in the schema's properties map.
func TestLatestBranchJSONSchema_EncoderMatchesSchema(t *testing.T) {
	root := lbLoadSchemaFile(t, latestBranchSchemaFilename)
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

	gotKeys := lbFirstObjectKeys(t, buf.Bytes())
	for _, key := range gotKeys {
		if _, allowed := props[key]; !allowed {
			t.Errorf("encoder emitted %q not declared in schema properties", key)
		}
	}
}

// --- local schema-test helpers (minimal equivalents of cmd's infra) ---
// lbPackageDir and lbFirstObjectKeys are defined in
// latestbranchjson_contract_test.go.

// lbFindPublishedSchema walks up from the package dir to locate
// 02-spec/08-json-schemas/<filename>.
func lbFindPublishedSchema(t *testing.T, filename string) string {
	t.Helper()
	dir := filepath.Dir(lbPackageDir())
	for i := 0; i < 8; i++ {
		candidate := filepath.Join(dir, "02-spec", "08-json-schemas", filename)
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}

		dir = parent
	}

	t.Fatalf("published schema %s not found walking up from %s", filename, lbPackageDir())

	return ""
}

func lbLoadSchemaFile(t *testing.T, filename string) map[string]any {
	t.Helper()
	raw, err := os.ReadFile(lbFindPublishedSchema(t, filename))
	if err != nil {
		t.Fatalf("read schema: %v", err)
	}

	var s map[string]any
	if err := json.Unmarshal(raw, &s); err != nil {
		t.Fatalf("parse schema: %v", err)
	}

	return s
}

func lbStringSliceFromAny(v any) []string {
	arr, ok := v.([]any)
	if !ok {
		return nil
	}

	out := make([]string, 0, len(arr))
	for _, e := range arr {
		s, ok := e.(string)
		if !ok {
			return nil
		}

		out = append(out, s)
	}

	return out
}

func lbEqualStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
