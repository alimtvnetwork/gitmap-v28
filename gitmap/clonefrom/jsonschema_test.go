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

	"github.com/alimtvnetwork/gitmap-v28/gitmap/clonenow"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func TestEmitSchema_ReportShape(t *testing.T) {
	body, err := EmitSchema(constants.EmitSchemaKindReport)
	isSchemaErr := err != nil
	if isSchemaErr == true {
		t.Fatalf("EmitSchema(report) returned error: %v", err)
	}
	root := decodeSchema(t, body)
	assertString(t, root, "$schema", constants.JSONSchemaDialect2020_12)
	assertString(t, root, "$id", constants.CloneFromSchemaIDReport)
	verifyReportProperties(t, root)
}

func verifyReportProperties(t *testing.T, root map[string]any) {
	props, ok := root["properties"].(map[string]any)
	isPropsMissing := !ok
	if isPropsMissing == true {
		t.Fatalf("report schema missing properties object: %T", root["properties"])
	}
	assertReportKeys(t, props)
	verifySchemaVersionConst(t, props)
}

func assertReportKeys(t *testing.T, props map[string]any) {
	for _, key := range []string{"schemaVersion", "transport", "rows"} {
		_, hasKey := props[key]
		isKeyMissing := !hasKey
		if isKeyMissing == true {
			t.Errorf("report schema missing required property %q", key)
		}
	}
}

func TestEmitSchema_InputShape(t *testing.T) {
	body, err := EmitSchema(constants.EmitSchemaKindInput)
	isSchemaErr := err != nil
	if isSchemaErr == true {
		t.Fatalf("EmitSchema(input) returned error: %v", err)
	}
	root := decodeSchema(t, body)
	assertString(t, root, "$schema", constants.JSONSchemaDialect2020_12)
	assertString(t, root, "$id", constants.CloneFromSchemaIDInput)
	assertString(t, root, "type", "array")
	verifyInputItems(t, root)
}

func verifyInputItems(t *testing.T, root map[string]any) {
	item, ok := root["items"].(map[string]any)
	isItemMissing := !ok
	if isItemMissing == true {
		t.Fatalf("input schema items must be an object, got %T", root["items"])
	}
	itemProps, ok := item["properties"].(map[string]any)
	isItemPropsMissing := !ok
	if isItemPropsMissing == true {
		t.Fatalf("input schema items.properties must be an object, got %T", item["properties"])
	}
	assertInputFields(t, itemProps)
}

func assertInputFields(t *testing.T, itemProps map[string]any) {
	for _, name := range clonenow.KnownScanFields() {
		_, hasKey := itemProps[name]
		isFieldMissing := !hasKey
		if isFieldMissing == true {
			t.Errorf("input schema missing accepted field %q", name)
		}
	}
}

func TestEmitSchema_UnknownKindUsesConstantMessage(t *testing.T) {
	_, err := EmitSchema("nope")
	isNilErr := err == nil
	if isNilErr == true {
		t.Fatal("expected error for unknown kind, got nil")
	}
	assertUnknownKindErrorText(t, err.Error())
}

func assertUnknownKindErrorText(t *testing.T, errMsg string) {
	isBadKindMissing := !strings.Contains(errMsg, "nope")
	if isBadKindMissing == true {
		t.Errorf("error %q should mention the bad kind", errMsg)
	}
	isKindsMissing := !strings.Contains(errMsg, "report") || !strings.Contains(errMsg, "input")
	if isKindsMissing == true {
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
	if isUnmarshalFailed == true {
		t.Fatalf("emitted schema is not valid JSON: %v\n---\n%s", err, body)
	}

	return root
}

// assertString fails the test if obj[key] is not the expected
// string. Centralized so call sites stay one-liners.
func assertString(t *testing.T, obj map[string]any, key, want string) {
	t.Helper()
	got, ok := obj[key].(string)
	isTypeMismatch := !ok
	if isTypeMismatch == true {
		t.Errorf("expected %q to be string, got %T", key, obj[key])

		return
	}
	isMismatch := got != want
	if isMismatch == true {
		t.Errorf("%q = %q; want %q", key, got, want)
	}
}

// verifySchemaVersionConst checks that the report schema's
// schemaVersion property is a `const` integer equal to the live
// constants.CloneFromReportSchemaVersion.
func verifySchemaVersionConst(t *testing.T, props map[string]any) {
	t.Helper()
	sv, ok := props["schemaVersion"].(map[string]any)
	isSchemaMismatch := !ok
	if isSchemaMismatch == true {
		t.Fatalf("schemaVersion must be a sub-schema object, got %T", props["schemaVersion"])
	}
	verifySchemaVersionValue(t, sv)
}

func verifySchemaVersionValue(t *testing.T, sv map[string]any) {
	t.Helper()
	constVal, hasConst := sv["const"]
	isConstMissing := !hasConst
	if isConstMissing == true {
		t.Fatal("schemaVersion sub-schema must declare a const value")
	}
	checkSchemaNumericConst(t, constVal)
}

func checkSchemaNumericConst(t *testing.T, constVal any) {
	t.Helper()
	asFloat, isNumber := constVal.(float64)
	isNonNumber := isNumber == false
	if isNonNumber {
		t.Fatalf("schemaVersion const must be numeric, got %T", constVal)
	}
	isMismatch := int(asFloat) != constants.CloneFromReportSchemaVersion
	if isMismatch == true {
		t.Errorf("schemaVersion const = %v; want %d (live constant)",
			asFloat, constants.CloneFromReportSchemaVersion)
	}
}
