package cmd

// JSON schema contract for `gitmap amend audit` file output. Pairs
// the runtime encoder (encodeAmendAuditJSON / buildAuditRecord in
// amendaudit.go + amendauditrender.go) with the published schema at
// spec/08-json-schemas/amend-audit.schema.json so drift in either
// side fails the build.

import (
	"bytes"
	"sort"
	"testing"
)

const amendAuditSchemaFilename = "amend-audit.schema.json"

// amendAuditTopLevelRequiredKeys mirrors the schema's required array.
var amendAuditTopLevelRequiredKeys = []string{
	"branch",
	"commits",
	"forcePushed",
	"fromCommit",
	"id",
	"mode",
	"newAuthor",
	"previousAuthor",
	"timestamp",
	"toCommit",
	"totalCommits",
}

// TestAmendAuditJSONSchema_TopLevelShape pins the root type (object)
// and the required key set against the schema.
func TestAmendAuditJSONSchema_TopLevelShape(t *testing.T) {
	root := loadSchemaFile(t, amendAuditSchemaFilename)
	if root["type"] != "object" {
		t.Fatalf("top-level type = %v, want object", root["type"])
	}
	got := stringSliceFromAny(root["required"])
	sort.Strings(got)
	if !equalStringSlices(got, amendAuditTopLevelRequiredKeys) {
		t.Fatalf("required = %v, want %v", got, amendAuditTopLevelRequiredKeys)
	}
}

// TestAmendAuditJSONSchema_EncoderMatchesSchema runs the real
// stablejson encoder, then asserts every key in the emitted object
// is declared in the schema's properties map.
func TestAmendAuditJSONSchema_EncoderMatchesSchema(t *testing.T) {
	root := loadSchemaFile(t, amendAuditSchemaFilename)
	props, ok := root["properties"].(map[string]any)
	if !ok {
		t.Fatalf("schema missing properties object")
	}

	record := canonicalAmendmentRecord()
	var buf bytes.Buffer
	if err := encodeAmendAuditJSON(&buf, record); err != nil {
		t.Fatalf("encode: %v", err)
	}
	gotKeys := readFirstObjectKeys(t, buf.Bytes())
	for _, key := range gotKeys {
		if _, allowed := props[key]; !allowed {
			t.Errorf("encoder emitted %q not declared in schema properties", key)
		}
	}
}
