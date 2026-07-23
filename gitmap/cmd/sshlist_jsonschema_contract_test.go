package cmd

// JSON schema contract for `gitmap ssh list --json`. Pairs the
// runtime encoder (encodeSSHListJSON / buildSSHListJSONItems
// in sshlistrender.go) with the published schema at
// spec/08-json-schemas/ssh-list.schema.json so drift in either
// side fails the build.

import (
	"bytes"
	"sort"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

const sshListSchemaFilename = "ssh-list.schema.json"

// sshListTopLevelRequiredKeys mirrors the schema's items.required array.
var sshListTopLevelRequiredKeys = []string{
	"createdAt",
	"fingerprint",
	"id",
	"name",
	"privatePath",
}

// TestSSHListJSONSchema_TopLevelShape pins the root type (array)
// and the items.required key set against the schema.
func TestSSHListJSONSchema_TopLevelShape(t *testing.T) {
	root := loadSchemaFile(t, sshListSchemaFilename)
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
	if !equalStringSlices(got, sshListTopLevelRequiredKeys) {
		t.Fatalf("items.required = %v, want %v", got, sshListTopLevelRequiredKeys)
	}
}

// TestSSHListJSONSchema_EncoderMatchesSchema runs the real
// stablejson encoder, then asserts every key in the first emitted
// object is declared in the schema's items.properties map.
func TestSSHListJSONSchema_EncoderMatchesSchema(t *testing.T) {
	root := loadSchemaFile(t, sshListSchemaFilename)
	items, _ := root["items"].(map[string]any)
	props, ok := items["properties"].(map[string]any)
	if !ok {
		t.Fatalf("schema missing items.properties object")
	}

	records := []model.SSHKey{canonicalSSHKey()}
	var buf bytes.Buffer
	if err := encodeSSHListJSON(&buf, records); err != nil {
		t.Fatalf("encode: %v", err)
	}
	gotKeys := extractFirstObjectKeyOrder(t, buf.Bytes())
	for _, key := range gotKeys {
		if _, allowed := props[key]; !allowed {
			t.Errorf("encoder emitted %q not declared in schema items.properties", key)
		}
	}
}
