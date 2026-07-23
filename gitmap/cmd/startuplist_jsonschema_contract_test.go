package cmd

// Schema contract for `gitmap startup-list --json`. Pairs the
// runtime encoder (encodeStartupListJSON) with the published schema
// at spec/08-json-schemas/startup-list.schema.json so a drift in
// either side fails the build.
//
// Generic helpers (findSchemaFile, loadSchemaFile, propertyOrder
// extraction) live in jsonschema_helpers_test.go so future
// `*_jsonschema_contract_test.go` files (see
// spec/08-json-schemas/_TODO.md) can reuse them.
//
// Why a hand-rolled mini-validator instead of pulling in a real
// JSON-Schema library?
//
//   1. The project's go.mod is intentionally lean (one direct dep
//      besides stdlib + charmbracelet/sqlite/archives/fuzzy/sys).
//      Adding a 30k-LOC schema validator for one test would be a
//      disproportionate dependency tax.
//   2. The contract surface here is tiny: a top-level array of
//      objects with three required string keys and a fixed key
//      order. A 60-line bespoke check covers the contract precisely
//      and makes the assertions readable in the failure message.
//   3. If/when the schema set in spec/08-json-schemas/ grows past
//      ~5 commands, swapping this for github.com/santhosh-tekuri/
//      jsonschema becomes worthwhile and is a one-file change.

import (
	"bytes"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/startup"
)

// startupListSchemaFilename is the on-disk name under
// spec/08-json-schemas/. Centralized so a future rename touches one
// line.
const startupListSchemaFilename = "startup-list.schema.json"

// itemSchema descends into the array's `items` subschema where the
// per-entry object contract lives. Centralizes the navigation so
// the individual assertions below stay flat.
func itemSchema(t *testing.T, root map[string]any) map[string]any {
	t.Helper()
	items, ok := root["items"].(map[string]any)
	if !ok {
		t.Fatalf("schema has no items object")
	}

	return items
}

// TestStartupListSchema_TopLevelIsArray pins the most fundamental
// shape decision: empty output is `[]`, NOT `null` or `{}`. A
// future encoder bug that emitted `null` for empty would break
// every downstream `jq length` consumer.
func TestStartupListSchema_TopLevelIsArray(t *testing.T) {
	root := loadSchemaFile(t, startupListSchemaFilename)
	if root["type"] != "array" {
		t.Fatalf("schema top-level type = %v, want array", root["type"])
	}
}

// TestStartupListSchema_RequiredKeysMatchEncoder asserts the schema
// requires exactly the three keys the encoder emits. If a future
// PR adds a key to startup.Entry but forgets the schema, this test
// fails with a clear diff.
func TestStartupListSchema_RequiredKeysMatchEncoder(t *testing.T) {
	root := loadSchemaFile(t, startupListSchemaFilename)
	required := stringSliceFromAny(itemSchema(t, root)["required"])
	want := []string{"name", "path", "exec"}
	if !equalStringSlices(required, want) {
		t.Fatalf("schema required = %v, want %v", required, want)
	}
}

// TestStartupListSchema_PropertyOrderMatchesEncoder is the headline
// contract test: encode a real entry, parse the resulting JSON
// preserving key order, and assert the order matches the schema's
// propertyOrder array. This is the ONLY guard that catches a
// reordering of the stablejson.Field slice in startuplistrender.go
// — Go's encoding/json sorts map keys alphabetically so a generic
// json.Unmarshal would mask the bug.
func TestStartupListSchema_PropertyOrderMatchesEncoder(t *testing.T) {
	root := loadSchemaFile(t, startupListSchemaFilename)
	want := stringSliceFromAny(itemSchema(t, root)["propertyOrder"])
	if len(want) == 0 {
		t.Fatalf("schema item has no propertyOrder array")
	}
	entries := []startup.Entry{{Name: "n", Path: "p", Exec: "e"}}
	var buf bytes.Buffer
	if err := encodeStartupListJSON(&buf, entries); err != nil {
		t.Fatalf("encode: %v", err)
	}
	got := extractFirstObjectKeyOrder(t, buf.Bytes())
	if !equalStringSlices(got, want) {
		t.Fatalf("emitted key order = %v, schema propertyOrder = %v", got, want)
	}
}

// TestStartupListSchema_EmptyEncodesAsArray pins the empty-input
// behavior end-to-end. Belt-and-suspenders alongside the byte-exact
// contract test: that one fails if the bytes drift, this one fails
// if the SHAPE drifts (e.g., `null` vs `[]`).
func TestStartupListSchema_EmptyEncodesAsArray(t *testing.T) {
	var buf bytes.Buffer
	if err := encodeStartupListJSON(&buf, nil); err != nil {
		t.Fatalf("encode: %v", err)
	}
	trimmed := bytes.TrimSpace(buf.Bytes())
	if !bytes.Equal(trimmed, []byte("[]")) {
		t.Fatalf("empty encoded as %q, want []", trimmed)
	}
}
