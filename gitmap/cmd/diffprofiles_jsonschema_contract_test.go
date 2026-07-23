package cmd

// JSON schema contract for `gitmap diff-profiles --json`. Pairs
// the runtime encoder (encodeDiffProfilesJSON in diffprofilesrender.go)
// with the published schema at
// spec/08-json-schemas/diff-profiles.schema.json so drift in either
// side fails the build.

import (
	"bytes"
	"sort"
	"testing"
)

const diffProfilesSchemaFilename = "diff-profiles.schema.json"

// diffProfilesTopLevelRequiredKeys mirrors the schema's required array.
var diffProfilesTopLevelRequiredKeys = []string{
	"different",
	"onlyInA",
	"onlyInB",
	"profileA",
	"profileB",
	"same",
}

// TestDiffProfilesJSONSchema_TopLevelShape pins the root type
// (object) and the required key set against the schema.
func TestDiffProfilesJSONSchema_TopLevelShape(t *testing.T) {
	root := loadSchemaFile(t, diffProfilesSchemaFilename)
	if root["type"] != "object" {
		t.Fatalf("top-level type = %v, want object", root["type"])
	}
	got := stringSliceFromAny(root["required"])
	sort.Strings(got)
	if !equalStringSlices(got, diffProfilesTopLevelRequiredKeys) {
		t.Fatalf("required = %v, want %v", got, diffProfilesTopLevelRequiredKeys)
	}
}

// TestDiffProfilesJSONSchema_EncoderMatchesSchema runs the real
// stablejson encoder, then asserts every key in the emitted object
// is declared in the schema's properties map.
func TestDiffProfilesJSONSchema_EncoderMatchesSchema(t *testing.T) {
	root := loadSchemaFile(t, diffProfilesSchemaFilename)
	props, ok := root["properties"].(map[string]any)
	if !ok {
		t.Fatalf("schema missing properties object")
	}

	result := canonicalDPResult()
	var buf bytes.Buffer
	if err := encodeDiffProfilesJSON(&buf, "alpha", "beta", result); err != nil {
		t.Fatalf("encode: %v", err)
	}
	gotKeys := readFirstObjectKeys(t, buf.Bytes())
	for _, key := range gotKeys {
		if _, allowed := props[key]; !allowed {
			t.Errorf("encoder emitted %q not declared in schema properties", key)
		}
	}
}
