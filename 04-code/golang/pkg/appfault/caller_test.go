package appfault_test

import (
	"encoding/json"
	"strings"
	"testing"

	"coding-guidelines/common/pkg/appfault"
)

func TestCallerInfo_JSONMarshalingAndUnmarshaling(t *testing.T) {
	// 1. Standard struct marshaling
	caller := appfault.CallerInfo{
		Function: "HandleRegistration",
		File:     "auth/service.go",
		Line:     55,
	}

	jsonBytes, err := caller.ToJson()
	if err != nil {
		t.Fatalf("caller.ToJson failed: %v", err)
	}

	if !strings.Contains(string(jsonBytes), `"Function": "HandleRegistration"`) {
		t.Fatalf("expected Function in JSON: %s", string(jsonBytes))
	}

	// 2. Unmarshal from JSON object
	var restored appfault.CallerInfo
	if err := json.Unmarshal(jsonBytes, &restored); err != nil {
		t.Fatalf("unmarshal from object failed: %v", err)
	}

	if restored.Function != "HandleRegistration" || restored.File != "auth/service.go" || restored.Line != 55 {
		t.Fatalf("restored mismatch: %+v", restored)
	}

	// 3. Unmarshal from string format: "auth/service.go:55 (HandleRegistration)"
	strJSON := []byte(`"auth/service.go:55 (HandleRegistration)"`)
	var fromString appfault.CallerInfo
	if err := json.Unmarshal(strJSON, &fromString); err != nil {
		t.Fatalf("unmarshal from string failed: %v", err)
	}

	if fromString.Function != "HandleRegistration" || fromString.File != "auth/service.go" || fromString.Line != 55 {
		t.Fatalf("fromString mismatch: %+v", fromString)
	}

	// 4. Unmarshal from compact "file:line" string: "services/user/validator.go:42"
	compactJSON := []byte(`"services/user/validator.go:42"`)
	var fromCompact appfault.CallerInfo
	if err := json.Unmarshal(compactJSON, &fromCompact); err != nil {
		t.Fatalf("unmarshal from compact failed: %v", err)
	}

	if fromCompact.File != "services/user/validator.go" || fromCompact.Line != 42 {
		t.Fatalf("fromCompact mismatch: %+v", fromCompact)
	}

	// 5. Empty caller marshals to null
	emptyCaller := appfault.CallerInfo{}
	emptyBytes, err := emptyCaller.MarshalJSON()
	if err != nil {
		t.Fatalf("emptyCaller.MarshalJSON failed: %v", err)
	}

	if string(emptyBytes) != "null" {
		t.Fatalf("expected null, got %s", string(emptyBytes))
	}
}
