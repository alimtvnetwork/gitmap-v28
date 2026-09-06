package logleveltype_test

import (
	"encoding/json"
	"testing"

	"coding-guidelines/common/pkg/enum/logleveltype"
)

func TestInvalidZeroValue(t *testing.T) {
	var zero logleveltype.Variant
	if zero != logleveltype.Invalid {
		t.Fatalf("expected zero value to be Invalid, got %v", zero)
	}

	if zero != logleveltype.Unknown {
		t.Fatalf("expected Invalid to equal Unknown")
	}

	if zero.IsValid() {
		t.Fatalf("expected Invalid to not be valid")
	}

	if !zero.IsInvalid() {
		t.Fatalf("expected Invalid to be IsInvalid")
	}

	if zero.String() != "Unknown" {
		t.Fatalf("expected String 'Unknown', got %s", zero.String())
	}
}

func TestVariantsAndLabels(t *testing.T) {
	tests := []struct {
		variant      logleveltype.Variant
		expectedName string
	}{
		{logleveltype.Debug, "Debug"},
		{logleveltype.Info, "Info"},
		{logleveltype.Warn, "Warn"},
		{logleveltype.Error, "Error"},
		{logleveltype.Fatal, "Fatal"},
	}

	for _, tc := range tests {
		if tc.variant.Name() != tc.expectedName {
			t.Errorf("expected name %s, got %s", tc.expectedName, tc.variant.Name())
		}

		if tc.variant.Label() != tc.expectedName {
			t.Errorf("expected label %s, got %s", tc.expectedName, tc.variant.Label())
		}

		if !tc.variant.IsValid() {
			t.Errorf("expected %s to be valid", tc.expectedName)
		}
	}
}

func TestCheckers(t *testing.T) {
	d := logleveltype.Debug
	if !d.IsDebug() {
		t.Fatalf("expected IsDebug true")
	}

	if d.IsInfo() {
		t.Fatalf("expected IsInfo false")
	}

	i := logleveltype.Info
	if !i.IsInfo() {
		t.Fatalf("expected IsInfo true")
	}

	w := logleveltype.Warn
	if !w.IsWarn() {
		t.Fatalf("expected IsWarn true")
	}

	e := logleveltype.Error
	if !e.IsError() {
		t.Fatalf("expected IsError true")
	}

	f := logleveltype.Fatal
	if !f.IsFatal() {
		t.Fatalf("expected IsFatal true")
	}
}

func TestIsEnabled(t *testing.T) {
	if !logleveltype.Info.IsEnabled(logleveltype.Debug) {
		t.Fatalf("expected Info to be enabled when threshold is Debug")
	}

	if !logleveltype.Info.IsEnabled(logleveltype.Info) {
		t.Fatalf("expected Info to be enabled when threshold is Info")
	}

	if logleveltype.Debug.IsEnabled(logleveltype.Info) {
		t.Fatalf("expected Debug to not be enabled when threshold is Info")
	}

	if !logleveltype.Fatal.IsEnabled(logleveltype.Error) {
		t.Fatalf("expected Fatal to be enabled when threshold is Error")
	}
}

func TestParse(t *testing.T) {
	res := logleveltype.Parse("warn")
	if res.IsFailed() {
		t.Fatalf("expected parse to succeed: %v", res.Fault())
	}

	if res.Data() != logleveltype.Warn {
		t.Fatalf("expected Warn, got %v", res.Data())
	}

	resUnknown := logleveltype.Parse("unknown")
	if resUnknown.IsFailed() || resUnknown.Data() != logleveltype.Unknown {
		t.Fatalf("expected Unknown, got %v", resUnknown.Data())
	}

	badRes := logleveltype.Parse("nonexistent-level")
	if badRes.IsSuccess() {
		t.Fatalf("expected failure on bad string")
	}
}

func TestAllAndValues(t *testing.T) {
	all := logleveltype.All()
	if len(all) != 5 {
		t.Fatalf("expected 5 valid variants, got %d", len(all))
	}

	values := logleveltype.Values()
	if len(values) != 5 {
		t.Fatalf("expected 5 string values, got %d", len(values))
	}
}

func TestJSONRoundtrip(t *testing.T) {
	val := logleveltype.Warn

	data, err := json.Marshal(val)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	if string(data) != `"Warn"` {
		t.Fatalf("unexpected json: %s", string(data))
	}

	var parsed logleveltype.Variant
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if parsed != logleveltype.Warn {
		t.Fatalf("expected Warn, got %v", parsed)
	}
}

func TestUnmarshalJSON_InvalidCases(t *testing.T) {
	var v logleveltype.Variant

	// Invalid string must return error
	if err := json.Unmarshal([]byte(`"NonExistentLevel"`), &v); err == nil {
		t.Fatalf("expected error unmarshaling invalid string, got nil")
	}

	// Invalid numeric byte must return error
	if err := json.Unmarshal([]byte(`99`), &v); err == nil {
		t.Fatalf("expected error unmarshaling invalid numeric byte 99, got nil")
	}

	// Null should unmarshal to Invalid without error
	if err := json.Unmarshal([]byte(`null`), &v); err != nil {
		t.Fatalf("expected nil error on null, got %v", err)
	}

	if v != logleveltype.Invalid {
		t.Fatalf("expected Invalid on null, got %v", v)
	}

	// Valid numeric byte
	if err := json.Unmarshal([]byte(`1`), &v); err != nil {
		t.Fatalf("expected nil error on numeric 1, got %v", err)
	}

	if v != logleveltype.Debug {
		t.Fatalf("expected Debug on numeric 1, got %v", v)
	}
}
