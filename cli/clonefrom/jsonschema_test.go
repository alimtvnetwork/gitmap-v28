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

	"github.com/alimtvnetwork/gitmap-v28/cli/clonenow"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func TestEmitSchema_ReportShape(t *testing.T) {
	body, err := EmitSchema(constants.EmitSchemaKindReport)
	isSchemaErr := err != nil
	if isSchemaErr {
		t.Fatalf("EmitSchema(report) returned error: %v", err)
	}

	root := decodeSchema(t, body)
	assertString(t, root, "$schema", constants.JSONSchemaDialect2020_12)
	assertString(t, root, "$id", constants.CloneFromSchemaIDReport)
	verifyReportProperties(t, root)
}

func verifyReportProperties(t *testing.T, root map[string]any) {
	props, isPropsMap := root["properties"].(map[string]any)
	if !isPropsMap {
		t.Fatalf("report schema missing properties object: %T", root["properties"])
	}

	assertReportKeys(t, props)
	verifySchemaVersionConst(t, props)
}

func assertReportKeys(t *testing.T, props map[string]any) {
	for _, key := range []string{"schemaVersion", "transport", "rows"} {
		_, hasKey := props[key]
		if !hasKey {
			t.Errorf("report schema missing required property %q", key)
		}
	}
}

func TestEmitSchema_InputShape(t *testing.T) {
	body, err := EmitSchema(constants.EmitSchemaKindInput)
	isSchemaErr := err != nil
	if isSchemaErr {
		t.Fatalf("EmitSchema(input) returned error: %v", err)
	}

	root := decodeSchema(t, body)
	assertString(t, root, "$schema", constants.JSONSchemaDialect2020_12)
	assertString(t, root, "$id", constants.CloneFromSchemaIDInput)
	assertString(t, root, "type", "array")
	verifyInputItems(t, root)
}

func verifyInputItems(t *testing.T, root map[string]any) {
	item, isItemMap := root["items"].(map[string]any)
	if !isItemMap {
		t.Fatalf("input schema items must be an object, got %T", root["items"])
	}

	itemProps, isItemPropsMap := item["properties"].(map[string]any)
	if !isItemPropsMap {
		t.Fatalf("input schema items.properties must be an object, got %T", item["properties"])
	}

	assertInputFields(t, itemProps)
}

func assertInputFields(t *testing.T, itemProps map[string]any) {
	for _, name := range clonenow.KnownScanFields() {
		_, hasKey := itemProps[name]
		if !hasKey {
			t.Errorf("input schema missing accepted field %q", name)
		}
	}
}

func TestEmitSchema_UnknownKindUsesConstantMessage(t *testing.T) {
	_, err := EmitSchema("nope")
	isNilErr := err == nil
	if isNilErr {
		t.Fatal("expected error for unknown kind, got nil")
	}

	assertUnknownKindErrorText(t, err.Error())
}

func assertUnknownKindErrorText(t *testing.T, errMsg string) {
	isBadKindMissing := !strings.Contains(errMsg, "nope")
	if isBadKindMissing {
		t.Errorf("error %q should mention the bad kind", errMsg)
	}

	isKindsMissing := !strings.Contains(errMsg, "report") || !strings.Contains(errMsg, "input")
	if isKindsMissing {
		t.Errorf("error %q should list both accepted kinds", errMsg)
	}
}

// decodeSchema parses the emitted bytes as generic JSON, failing
// the test on any parse error. Returns the root object.
func decodeSchema(t *testing.T, body []byte) map[string]any {
	t.Helper()
	var root map[string]any
	err := json.Unmarshal(body, &root)
	isUnmarshalFailed := err != nil
	if isUnmarshalFailed {
		t.Fatalf("emitted schema is not valid JSON: %v\n---\n%s", err, body)
	}

	return root
}

// assertString fails the test if obj[key] is not the expected
// string. Centralized so call sites stay one-liners.
func assertString(t *testing.T, obj map[string]any, key, want string) {
	t.Helper()
	got, isString := obj[key].(string)
	if !isString {
		t.Errorf("expected %q to be string, got %T", key, obj[key])

		return
	}

	isMismatch := got != want
	if isMismatch {
		t.Errorf("%q = %q; want %q", key, got, want)
	}
}

// verifySchemaVersionConst checks that the report schema's
// schemaVersion property is a `const` integer equal to the live
// constants.CloneFromReportSchemaVersion.
func verifySchemaVersionConst(t *testing.T, props map[string]any) {
	t.Helper()
	sv, isSchemaMap := props["schemaVersion"].(map[string]any)
	if !isSchemaMap {
		t.Fatalf("schemaVersion must be a sub-schema object, got %T", props["schemaVersion"])
	}

	verifySchemaVersionValue(t, sv)
}

func verifySchemaVersionValue(t *testing.T, sv map[string]any) {
	t.Helper()
	constVal, hasConst := sv["const"]
	if !hasConst {
		t.Fatal("schemaVersion sub-schema must declare a const value")
	}

	checkSchemaNumericConst(t, constVal)
}

func checkSchemaNumericConst(t *testing.T, constVal any) {
	t.Helper()
	asFloat, isNumber := constVal.(float64)
	if !isNumber {
		t.Fatalf("schemaVersion const must be numeric, got %T", constVal)
	}

	isMismatch := int(asFloat) != constants.CloneFromReportSchemaVersion
	if isMismatch {
		t.Errorf("schemaVersion const = %v; want %d (live constant)",
			asFloat, constants.CloneFromReportSchemaVersion)
	}
}
