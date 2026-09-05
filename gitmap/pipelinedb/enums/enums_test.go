package enums

import (
	"encoding/json"
	"testing"
)

func TestEnums_TableConstants(t *testing.T) {
	if PipelineSplitDbTable != "PipelineSplitDb" {
		t.Errorf("expected PipelineSplitDb, got %s", PipelineSplitDbTable)
	}
	if PipelineRunRecordTable != "PipelineRunRecord" {
		t.Errorf("expected PipelineRunRecord, got %s", PipelineRunRecordTable)
	}
	if PipelineRunTable != PipelineRunRecordTable {
		t.Errorf("expected PipelineRunTable == PipelineRunRecordTable")
	}
	if PipelineErrorRecordTable != "PipelineErrorRecord" {
		t.Errorf("expected PipelineErrorRecord, got %s", PipelineErrorRecordTable)
	}
	if PipelineErrorTable != PipelineErrorRecordTable {
		t.Errorf("expected PipelineErrorTable == PipelineErrorRecordTable")
	}
	if PipelineDbStatsTable != "PipelineDbStats" {
		t.Errorf("expected PipelineDbStats, got %s", PipelineDbStatsTable)
	}
}

func TestEnums_FieldTypeReceivers(t *testing.T) {
	field := PipelineRunRecordDb.RunId

	if field.Name() != "RunId" {
		t.Errorf("expected 'RunId', got '%s'", field.Name())
	}
	if field.String() != "RunId" {
		t.Errorf("expected 'RunId', got '%s'", field.String())
	}
	if field.Value() != "RunId" {
		t.Errorf("expected 'RunId', got '%s'", field.Value())
	}
	if !field.IsCompare(PipelineRunRecordDb.RunId) {
		t.Errorf("expected IsCompare true")
	}
	if field.IsCompare(PipelineRunRecordDb.WorkflowName) {
		t.Errorf("expected IsCompare false for different fields")
	}
	if !field.IsEnum() {
		t.Errorf("expected IsEnum true for valid field")
	}

	invalidField := PipelineRunRecordFieldType("UnknownField")
	if invalidField.IsEnum() {
		t.Errorf("expected IsEnum false for unknown field")
	}

	if !field.IsRunId() {
		t.Errorf("expected IsRunId true")
	}
	if field.IsWorkflowName() {
		t.Errorf("expected IsWorkflowName false")
	}
}

func TestEnums_JSONMarshaling(t *testing.T) {
	field := PipelineRunRecordDb.WorkflowName

	// MarshalJSON
	bytes, err := json.Marshal(field)
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}
	if string(bytes) != `"WorkflowName"` {
		t.Errorf("expected %q, got %s", `"WorkflowName"`, string(bytes))
	}

	// UnmarshalJSON valid
	var unmarshaled PipelineRunRecordFieldType
	if err := json.Unmarshal(bytes, &unmarshaled); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}
	if unmarshaled != field {
		t.Errorf("expected '%s', got '%s'", field, unmarshaled)
	}

	// UnmarshalJSON invalid
	badJSON := []byte(`"BogusField"`)
	var badTarget PipelineRunRecordFieldType
	if err := json.Unmarshal(badJSON, &badTarget); err == nil {
		t.Errorf("expected error unmarshaling bogus field enum, got nil")
	}

	// ToJSON
	jsonStr, appErr := field.ToJSON()
	if appErr != nil {
		t.Fatalf("ToJSON failed: %v", appErr)
	}
	if jsonStr != `"WorkflowName"` {
		t.Errorf("expected %q, got %s", `"WorkflowName"`, jsonStr)
	}

	// FromJSON valid
	var fromTarget PipelineRunRecordFieldType
	if fromErr := fromTarget.FromJSON(jsonStr); fromErr != nil {
		t.Fatalf("FromJSON failed: %v", fromErr)
	}
	if fromTarget != field {
		t.Errorf("expected '%s', got '%s'", field, fromTarget)
	}

	// FromJSON invalid
	var fromBad PipelineRunRecordFieldType
	if fromErr := fromBad.FromJSON(`"BogusField"`); fromErr == nil {
		t.Errorf("expected error from FromJSON with invalid field, got nil")
	}
}

func TestEnums_Registry(t *testing.T) {
	reg := PipelineRunRecordDb
	all := reg.All()
	if len(all) == 0 {
		t.Fatalf("expected non-empty All()")
	}

	names := reg.Names()
	if len(names) != len(all) {
		t.Errorf("expected Names() len == All() len (%d vs %d)", len(names), len(all))
	}

	if !reg.IsEnum(reg.RunId) {
		t.Errorf("expected IsEnum true for RunId")
	}
	if reg.IsEnum(PipelineRunRecordFieldType("NoSuchField")) {
		t.Errorf("expected IsEnum false for NoSuchField")
	}

	if !reg.IsRunId(reg.RunId) {
		t.Errorf("expected IsRunId true")
	}
	if reg.IsRunId(reg.WorkflowName) {
		t.Errorf("expected IsRunId false for WorkflowName")
	}

	regJSON, regErr := reg.ToJSON()
	if regErr != nil {
		t.Fatalf("reg.ToJSON failed: %v", regErr)
	}
	if len(regJSON) == 0 {
		t.Errorf("expected non-empty registry JSON")
	}

	// Verify aliases
	if PipelineRunDb.RunId != PipelineRunRecordDb.RunId {
		t.Errorf("expected PipelineRunDb.RunId == PipelineRunRecordDb.RunId")
	}
	if PipelineErrorDb.ErrorText != PipelineErrorRecordDb.ErrorText {
		t.Errorf("expected PipelineErrorDb.ErrorText == PipelineErrorRecordDb.ErrorText")
	}
}
