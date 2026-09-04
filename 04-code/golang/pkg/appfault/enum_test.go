package appfault_test

import (
	"encoding/json"
	"testing"

	"coding-guidelines/common/pkg/appfault"
)

func TestSeverityTypeEnumAndJSON(t *testing.T) {
	sev := appfault.SeverityError
	data, err := json.Marshal(sev)
	if err != nil || string(data) != "\"Error\"" || sev.Name() != "Error" {
		t.Fatalf("expected \"Error\" JSON, got %s", string(data))
	}

	var parsed appfault.SeverityType
	if err := json.Unmarshal([]byte("\"Critical\""), &parsed); err != nil || parsed != appfault.SeverityCritical {
		t.Fatalf("expected SeverityCritical, got %v", parsed)
	}
}

func TestPriorityTypeEnumAndJSON(t *testing.T) {
	pri := appfault.PriorityHigh
	data, err := json.Marshal(pri)
	if err != nil || string(data) != "\"High\"" || pri.Name() != "High" {
		t.Fatalf("expected \"High\" JSON, got %s", string(data))
	}

	var parsed appfault.PriorityType
	if err := json.Unmarshal([]byte("\"Low\""), &parsed); err != nil || parsed != appfault.PriorityLow {
		t.Fatalf("expected PriorityLow, got %v", parsed)
	}
}
