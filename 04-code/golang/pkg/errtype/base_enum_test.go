package errtype_test

import (
	"encoding/json"
	"testing"

	"coding-guidelines/common/pkg/errtype"
)

func TestBaseEnum_VariationConforms(t *testing.T) {
	var e errtype.BaseEnumer = errtype.Validation
	if e.Name() != "Validation" {
		t.Fatalf("expected Name() == 'Validation', got %s", e.Name())
	}

	if e.ValueString() != "2" {
		t.Fatalf("expected ValueString() == '2', got %s", e.ValueString())
	}

	if !e.IsEnum() {
		t.Fatal("expected Validation to be registered enum")
	}

	var aliasE errtype.BaseEnum = e
	if aliasE.Name() != "Validation" {
		t.Fatalf("expected alias Name() == 'Validation', got %s", aliasE.Name())
	}

	var ne errtype.NumberEnumer = errtype.NotFound
	if ne.Int() != 3 || ne.Code() != 3 {
		t.Fatalf("unexpected number enum values: int=%d code=%d", ne.Int(), ne.Code())
	}

	var aliasNE errtype.NumberEnum = ne
	if aliasNE.Int() != 3 || aliasNE.Code() != 3 {
		t.Fatalf("unexpected alias number enum values: int=%d code=%d", aliasNE.Int(), aliasNE.Code())
	}
}

func TestProcessStateType_Lifecycle(t *testing.T) {
	state := errtype.ProcessStateRunning

	if state.Name() != "Running" {
		t.Fatalf("expected Name() == 'Running', got %s", state.Name())
	}

	if !state.IsValid() {
		t.Fatal("expected Running to be valid")
	}

	if !state.IsEnum() {
		t.Fatal("expected Running to be registered enum")
	}

	// JSON roundtrip
	data, err := json.Marshal(state)
	if err != nil {
		t.Fatalf("json marshal failed: %v", err)
	}

	var unmarshaled errtype.ProcessStateType
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("json unmarshal failed: %v", err)
	}

	if unmarshaled != errtype.ProcessStateRunning {
		t.Fatalf("expected unmarshaled == Running, got %s", unmarshaled)
	}

	// Case-insensitive parsing
	parsed := errtype.ParseProcessState("running")
	if parsed != errtype.ProcessStateRunning {
		t.Fatalf("expected ParseProcessState('running') to match, got %s", parsed)
	}
}

func TestLogLevelType_Lifecycle(t *testing.T) {
	lvl := errtype.LogLevelWarn

	if lvl.Name() != "Warn" {
		t.Fatalf("expected Name() == 'Warn', got %s", lvl.Name())
	}

	if lvl.Code() != 3 || lvl.Int() != 3 {
		t.Fatalf("unexpected code/int: %d/%d", lvl.Code(), lvl.Int())
	}

	if !lvl.IsValid() || !lvl.IsEnum() {
		t.Fatal("expected LogLevelWarn to be valid and enum")
	}

	// JSON string marshal
	data, err := json.Marshal(lvl)
	if err != nil {
		t.Fatalf("json marshal failed: %v", err)
	}

	var unmarshaled errtype.LogLevelType
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("json unmarshal failed: %v", err)
	}

	if unmarshaled != errtype.LogLevelWarn {
		t.Fatalf("expected unmarshaled == LogLevelWarn, got %v", unmarshaled)
	}

	// Parsing
	parsed := errtype.ParseLogLevel("warn")
	if parsed != errtype.LogLevelWarn {
		t.Fatalf("expected ParseLogLevel('warn') to match, got %v", parsed)
	}
}

func TestToEnum_GenericHelper(t *testing.T) {
	// Search by name
	foundState, ok := errtype.ToEnum("Completed", errtype.AllProcessStates())
	if !ok || foundState != errtype.ProcessStateCompleted {
		t.Fatalf("expected ToEnum to find ProcessStateCompleted, got %v, ok=%v", foundState, ok)
	}

	// Search by value string for number enum
	foundLvl, ok := errtype.ToEnum("4", errtype.AllLogLevels())
	if !ok || foundLvl != errtype.LogLevelError {
		t.Fatalf("expected ToEnum to find LogLevelError by '4', got %v, ok=%v", foundLvl, ok)
	}

	// Unknown enum search
	_, ok = errtype.ToEnum("non-existent", errtype.AllProcessStates())
	if ok {
		t.Fatal("expected non-existent enum to return ok=false")
	}
}

func TestProcessStateType_UnmarshalJSON_Invalid(t *testing.T) {
	var s errtype.ProcessStateType

	if err := json.Unmarshal([]byte(`"invalid_state"`), &s); err == nil {
		t.Fatalf("expected error unmarshaling invalid process state, got nil")
	}

	if err := json.Unmarshal([]byte(`null`), &s); err != nil {
		t.Fatalf("expected nil error on null, got %v", err)
	}

	if s != errtype.ProcessStateUnknown {
		t.Fatalf("expected ProcessStateUnknown on null, got %v", s)
	}
}

func TestLogLevelType_UnmarshalJSON_Invalid(t *testing.T) {
	var l errtype.LogLevelType

	if err := json.Unmarshal([]byte(`"invalid_level"`), &l); err == nil {
		t.Fatalf("expected error unmarshaling invalid log level, got nil")
	}

	if err := json.Unmarshal([]byte(`999`), &l); err == nil {
		t.Fatalf("expected error unmarshaling out-of-range numeric code 999, got nil")
	}

	if err := json.Unmarshal([]byte(`null`), &l); err != nil {
		t.Fatalf("expected nil error on null, got %v", err)
	}

	if err := json.Unmarshal([]byte(`2`), &l); err != nil {
		t.Fatalf("expected nil error on numeric 2, got %v", err)
	}

	if l != errtype.LogLevelInfo {
		t.Fatalf("expected LogLevelInfo on numeric 2, got %v", l)
	}
}
