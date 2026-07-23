package clonefrom

// Tests for the JSON-Schema emit surface backing
// `gitmap clone-from --emit-schema=<kind>`. Three contracts:
//
//  1. Both kinds emit valid, parseable JSON.
//  2. Both schemas declare the draft-2020-12 dialect via `$schema`
//     and a stable `$id`.
//  3. The report schema's `schemaVersion` const tracks the live
//     constants.CloneFromReportSchemaVersion — so a bump there is
//     guaranteed to reach downstream validators.
//
// Unknown-kind handling is also pinned: it must surface the
// user-facing error format from constants so the CLI message stays
// stable for shell-script consumers.

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/clonenow"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

func TestEmitSchema_ReportShape(t *testing.T) {
	body, err := EmitSchema(constants.EmitSchemaKindReport)
	if err != nil {
		t.Fatalf("EmitSchema(report) returned error: %v", err)
	}
	root := decodeSchema(t, body)
	assertString(t, root, "$schema", constants.JSONSchemaDialect2020_12)
	assertString(t, root, "$id", constants.CloneFromSchemaIDReport)
	props, ok := root["properties"].(map[string]any)
	if !ok {
		t.Fatalf("report schema missing properties object: %T", root["properties"])
	}
	for _, key := range []string{"schemaVersion", "transport", "rows"} {
		if _, hasKey := props[key]; !hasKey {
			t.Errorf("report schema missing required property %q", key)
		}
	}
	verifySchemaVersionConst(t, props)
}

func TestEmitSchema_InputShape(t *testing.T) {
	body, err := EmitSchema(constants.EmitSchemaKindInput)
	if err != nil {
		t.Fatalf("EmitSchema(input) returned error: %v", err)
	}
	root := decodeSchema(t, body)
	assertString(t, root, "$schema", constants.JSONSchemaDialect2020_12)
	assertString(t, root, "$id", constants.CloneFromSchemaIDInput)
	assertString(t, root, "type", "array")
	item, ok := root["items"].(map[string]any)
	if !ok {
		t.Fatalf("input schema items must be an object, got %T", root["items"])
	}
	itemProps, ok := item["properties"].(map[string]any)
	if !ok {
		t.Fatalf("input schema items.properties must be an object, got %T", item["properties"])
	}
	for _, name := range clonenow.KnownScanFields() {
		if _, hasKey := itemProps[name]; !hasKey {
			t.Errorf("input schema missing accepted field %q", name)
		}
	}
}

func TestEmitSchema_UnknownKindUsesConstantMessage(t *testing.T) {
	_, err := EmitSchema("nope")
	if err == nil {
		t.Fatal("expected error for unknown kind, got nil")
	}
	// Verify the error string matches the user-facing constant
	// format (substring match — fmt.Errorf adds the %q-quoted value).
	if !strings.Contains(err.Error(), "nope") {
		t.Errorf("error %q should mention the bad kind", err.Error())
	}
	if !strings.Contains(err.Error(), "report") || !strings.Contains(err.Error(), "input") {
		t.Errorf("error %q should list both accepted kinds", err.Error())
	}
}

// decodeSchema parses the emitted bytes as generic JSON, failing
// the test on any parse error. Returns the root object.
func decodeSchema(t *testing.T, body []byte) map[string]any {
	t.Helper()
	var root map[string]any
	if err := json.Unmarshal(body, &root); err != nil {
		t.Fatalf("emitted schema is not valid JSON: %v\n---\n%s", err, body)
	}

	return root
}

// assertString fails the test if obj[key] is not the expected
// string. Centralized so call sites stay one-liners.
func assertString(t *testing.T, obj map[string]any, key, want string) {
	t.Helper()
	got, ok := obj[key].(string)
	if !ok {
		t.Errorf("expected %q to be string, got %T", key, obj[key])

		return
	}
	if got != want {
		t.Errorf("%q = %q; want %q", key, got, want)
	}
}

// verifySchemaVersionConst checks that the report schema's
// schemaVersion property is a `const` integer equal to the live
// constants.CloneFromReportSchemaVersion. Encoding/json decodes
// JSON numbers as float64 — convert before comparing.
func verifySchemaVersionConst(t *testing.T, props map[string]any) {
	t.Helper()
	sv, ok := props["schemaVersion"].(map[string]any)
	if !ok {
		t.Fatalf("schemaVersion must be a sub-schema object, got %T", props["schemaVersion"])
	}
	constVal, hasConst := sv["const"]
	if !hasConst {
		t.Fatal("schemaVersion sub-schema must declare a const value")
	}
	asFloat, isNumber := constVal.(float64)
	if !isNumber {
		t.Fatalf("schemaVersion const must be numeric, got %T", constVal)
	}
	if int(asFloat) != constants.CloneFromReportSchemaVersion {
		t.Errorf("schemaVersion const = %v; want %d (live constant)",
			asFloat, constants.CloneFromReportSchemaVersion)
	}
}
