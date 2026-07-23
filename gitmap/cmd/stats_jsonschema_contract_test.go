package cmd

// JSON schema contract for `gitmap stats --json`. Pairs the runtime
// encoder (encodeStatsJSON / build*Items in statsrender.go) with the
// published schema at spec/08-json-schemas/stats.schema.json so drift
// in either side fails the build.

import (
	"bytes"
	"sort"
	"testing"
)

const statsSchemaFilename = "stats.schema.json"

// statsTopLevelRequiredKeys mirrors the schema's top-level required array.
var statsTopLevelRequiredKeys = []string{
	"avgDurationMs",
	"commands",
	"overallFailRate",
	"totalCommands",
	"totalFail",
	"totalSuccess",
	"uniqueCommands",
}

// TestStatsJSONSchema_TopLevelShape pins the root type (object) and
// the required-keys set against the schema.
func TestStatsJSONSchema_TopLevelShape(t *testing.T) {
	root := loadSchemaFile(t, statsSchemaFilename)
	if root["type"] != "object" {
		t.Fatalf("top-level type = %v, want object", root["type"])
	}
	got := stringSliceFromAny(root["required"])
	sort.Strings(got)
	if !equalStringSlices(got, statsTopLevelRequiredKeys) {
		t.Fatalf("required = %v, want %v", got, statsTopLevelRequiredKeys)
	}
}

// TestStatsJSONSchema_EncoderMatchesSchema runs the real stablejson
// encoder, then asserts every key in the emitted top-level object is
// declared in the schema's properties map.
func TestStatsJSONSchema_EncoderMatchesSchema(t *testing.T) {
	root := loadSchemaFile(t, statsSchemaFilename)
	props, ok := root["properties"].(map[string]any)
	if !ok {
		t.Fatalf("schema missing properties object")
	}

	overall, commands := canonicalStatsOverall(t)
	var buf bytes.Buffer
	if err := encodeStatsJSON(&buf, overall, commands); err != nil {
		t.Fatalf("encode: %v", err)
	}
	gotKeys := readFirstObjectKeys(t, buf.Bytes())
	for _, key := range gotKeys {
		if _, allowed := props[key]; !allowed {
			t.Errorf("encoder emitted %q not declared in schema properties", key)
		}
	}
}
