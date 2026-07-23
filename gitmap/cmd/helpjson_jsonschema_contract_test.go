package cmd

// Schema contract for `gitmap help --json`. Pairs the runtime
// encoder (printUsageJSON / helpJSONDoc) with the published schema at
// spec/08-json-schemas/help-json.schema.json so drift in either side
// fails the build. Uses the same generic helpers as
// listreleases_jsonschema_contract_test.go (findSchemaFile,
// loadSchemaFile, stringSliceFromAny) defined in
// jsonschema_helpers_test.go.

import (
	"encoding/json"
	"sort"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

const helpJSONSchemaFilename = "help-json.schema.json"

// TestHelpJSONSchema_TopLevelShape pins the root type and its
// required keys against the schema. A new top-level field added to
// helpJSONDoc without a matching schema update fails here.
func TestHelpJSONSchema_TopLevelShape(t *testing.T) {
	root := loadSchemaFile(t, helpJSONSchemaFilename)
	if root["type"] != "object" {
		t.Fatalf("top-level type = %v, want object", root["type"])
	}
	got := stringSliceFromAny(root["required"])
	sort.Strings(got)
	want := []string{"count", "groups", "version"}
	if !equalStringSlices(got, want) {
		t.Fatalf("required keys = %v, want %v", got, want)
	}
}

// TestHelpJSONSchema_EncoderMatchesSchema builds a real helpJSONDoc
// the same way printUsageJSON does, marshals it, and asserts the
// resulting payload only contains keys advertised in the schema's
// `properties` map. Catches an unsynced field on either side.
func TestHelpJSONSchema_EncoderMatchesSchema(t *testing.T) {
	root := loadSchemaFile(t, helpJSONSchemaFilename)
	props, ok := root["properties"].(map[string]any)
	if !ok {
		t.Fatalf("schema missing properties object")
	}

	rows := allHelpRows()
	doc := helpJSONDoc{Version: constants.Version, Count: len(rows)}
	doc.Groups = append(doc.Groups, helpJSONGroup{Group: "G", Lines: []string{"line"}})

	b, err := json.Marshal(doc)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var generic map[string]any
	if err := json.Unmarshal(b, &generic); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	for key := range generic {
		if _, allowed := props[key]; !allowed {
			t.Errorf("encoder emitted %q not declared in schema properties", key)
		}
	}
	if !strings.Contains(string(b), `"version":`) {
		t.Errorf("payload missing version key: %s", b)
	}
}
