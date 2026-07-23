package cmd

// JSON schema contract for `gitmap <type>-repos --json`. Pairs the
// runtime encoder (encodeProjectReposJSON / buildProjectReposJSONItems
// in projectreposrender.go) with the published schema at
// spec/08-json-schemas/project-repos.schema.json so drift fails the build.

import (
	"bytes"
	"sort"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

const projectReposSchemaFilename = "project-repos.schema.json"

// projectReposTopLevelRequiredKeys mirrors the schema's items.required array.
var projectReposTopLevelRequiredKeys = []string{
	"absolutePath",
	"id",
	"projectName",
	"projectType",
	"repoId",
}

// TestProjectReposJSONSchema_TopLevelShape pins root type + items.required.
func TestProjectReposJSONSchema_TopLevelShape(t *testing.T) {
	root := loadSchemaFile(t, projectReposSchemaFilename)
	if root["type"] != "array" {
		t.Fatalf("top-level type = %v, want array", root["type"])
	}
	items, ok := root["items"].(map[string]any)
	if !ok {
		t.Fatalf("schema missing items object")
	}
	if items["type"] != "object" {
		t.Fatalf("items.type = %v, want object", items["type"])
	}
	got := stringSliceFromAny(items["required"])
	sort.Strings(got)
	if !equalStringSlices(got, projectReposTopLevelRequiredKeys) {
		t.Fatalf("items.required = %v, want %v", got, projectReposTopLevelRequiredKeys)
	}
}

// TestProjectReposJSONSchema_EncoderMatchesSchema runs the real
// stablejson encoder, then asserts every key in the first emitted
// object is declared in the schema's items.properties map.
func TestProjectReposJSONSchema_EncoderMatchesSchema(t *testing.T) {
	root := loadSchemaFile(t, projectReposSchemaFilename)
	items, _ := root["items"].(map[string]any)
	props, ok := items["properties"].(map[string]any)
	if !ok {
		t.Fatalf("schema missing items.properties object")
	}

	projects := []model.DetectedProject{canonicalDetectedProject()}
	var buf bytes.Buffer
	if err := encodeProjectReposJSON(&buf, projects); err != nil {
		t.Fatalf("encode: %v", err)
	}
	gotKeys := extractFirstObjectKeyOrder(t, buf.Bytes())
	for _, key := range gotKeys {
		if _, allowed := props[key]; !allowed {
			t.Errorf("encoder emitted %q not declared in schema items.properties", key)
		}
	}
}
