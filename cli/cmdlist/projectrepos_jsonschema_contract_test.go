package cmdlist

// JSON schema contract for `gitmap <type>-repos --json`. Pairs the
// runtime encoder (encodeProjectReposJSON / buildProjectReposJSONItems
// in projectreposrender.go) with the published schema at
// 02-spec/08-json-schemas/project-repos.schema.json so drift fails the build.
//
// NOTE: this test lives in cmdlist (not cmd) because the encoder
// (encodeProjectReposJSON) is unexported here. Schema helpers are
// shared with projectreposjson_contract_test.go in this package.

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/model"
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
	root := pjLoadSchemaFile(t, projectReposSchemaFilename)
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

	got := pjStringSliceFromAny(items["required"])
	sort.Strings(got)
	if !pjEqualStringSlices(got, projectReposTopLevelRequiredKeys) {
		t.Fatalf("items.required = %v, want %v", got, projectReposTopLevelRequiredKeys)
	}
}

// TestProjectReposJSONSchema_EncoderMatchesSchema runs the real
// stablejson encoder, then asserts every key in the first emitted
// object is declared in the schema's items.properties map.
func TestProjectReposJSONSchema_EncoderMatchesSchema(t *testing.T) {
	root := pjLoadSchemaFile(t, projectReposSchemaFilename)
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

	gotKeys := pjFirstObjectKeys(t, buf.Bytes())
	for _, key := range gotKeys {
		if _, allowed := props[key]; !allowed {
			t.Errorf("encoder emitted %q not declared in schema items.properties", key)
		}
	}
}

// --- local schema-test helpers (minimal equivalents of cmd's infra) ---
// pjPackageDir and pjFirstObjectKeys are defined in
// projectreposjson_contract_test.go.

// pjFindPublishedSchema walks up from the package dir to locate
// 02-spec/08-json-schemas/<filename>.
func pjFindPublishedSchema(t *testing.T, filename string) string {
	t.Helper()
	dir := filepath.Dir(pjPackageDir())
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

	t.Fatalf("published schema %s not found walking up from %s", filename, pjPackageDir())

	return ""
}

func pjLoadSchemaFile(t *testing.T, filename string) map[string]any {
	t.Helper()
	raw, err := os.ReadFile(pjFindPublishedSchema(t, filename))
	if err != nil {
		t.Fatalf("read schema: %v", err)
	}

	var s map[string]any
	if err := json.Unmarshal(raw, &s); err != nil {
		t.Fatalf("parse schema: %v", err)
	}

	return s
}

func pjStringSliceFromAny(v any) []string {
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

func pjEqualStringSlices(a, b []string) bool {
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
