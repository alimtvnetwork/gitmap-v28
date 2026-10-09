package jsonx

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestNewEnvelope_Serialization(t *testing.T) {
	type SampleData struct {
		Message string `json:"message"`
		Count   int    `json:"count"`
	}

	payload := SampleData{Message: "hello world", Count: 42}
	env := NewEnvelope(TypeSSHNodes, "test-source", "gitmap test", "1.0", payload)

	bytes, err := json.Marshal(env)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	if !IsEnvelope(bytes) {
		t.Errorf("expected IsEnvelope to be true for marshaled envelope")
	}

	extracted, attrs, err := ExtractPayload(bytes)
	if err != nil {
		t.Fatalf("ExtractPayload failed: %v", err)
	}

	if attrs.Type != TypeSSHNodes {
		t.Errorf("expected type %s, got %s", TypeSSHNodes, attrs.Type)
	}

	var recovered SampleData
	err = json.Unmarshal(extracted, &recovered)
	if err != nil {
		t.Fatalf("Unmarshal extracted payload failed: %v", err)
	}
	if recovered.Message != "hello world" || recovered.Count != 42 {
		t.Errorf("payload mismatch: %+v", recovered)
	}
}

func TestDetectFormat_Envelope(t *testing.T) {
	raw := []byte(`{
		"attributes": {
			"type": "macro",
			"source": "cli/cmdmacro",
			"version": "1.0"
		},
		"data": {
			"name": "build-all",
			"steps": []
		}
	}`)

	desc, attrs, isMatched := DetectFormat(raw)
	if !isMatched {
		t.Fatalf("expected format to be matched")
	}
	if desc.Type != TypeMacro {
		t.Errorf("expected type %s, got %s", TypeMacro, desc.Type)
	}
	if attrs.Type != TypeMacro {
		t.Errorf("expected attrs.Type %s, got %s", TypeMacro, attrs.Type)
	}
}

func TestDetectFormat_Legacy(t *testing.T) {
	legacySSH := []byte(`{
		"schema_version": "1.0",
		"total_nodes": 2,
		"nodes": [
			{"alias": "worker-1", "ip_address": "192.168.1.10"}
		]
	}`)

	desc, attrs, isMatched := DetectFormat(legacySSH)
	if !isMatched {
		t.Fatalf("expected legacy SSH nodes to be matched")
	}
	if desc.Type != TypeSSHNodes {
		t.Errorf("expected type %s, got %s", TypeSSHNodes, desc.Type)
	}
	if attrs.Version != "legacy" {
		t.Errorf("expected attrs.Version to be 'legacy', got %s", attrs.Version)
	}
}

func TestDetectFormat_Unmatched(t *testing.T) {
	unmatchedJSON := []byte(`{
		"arbitrary_field": 12345,
		"unrelated_data": "some value"
	}`)

	_, _, isMatched := DetectFormat(unmatchedJSON)
	if isMatched {
		t.Errorf("expected arbitrary JSON to be unmatched")
	}
}

func TestDetectFormat_Fixtures(t *testing.T) {
	fixtures := []struct {
		filename     string
		expectedType string
		shouldMatch  bool
	}{
		{"ssh-nodes.json", TypeSSHNodes, true},
		{"macro.json", TypeMacro, true},
		{"commit-pull-config.json", TypeCommitPullConfig, true},
		{"ui-settings.json", TypeUISettings, true},
		{"unmatched.json", "", false},
	}

	for _, tc := range fixtures {
		path := filepath.Join("fixtures", tc.filename)
		bytes, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("failed to read fixture %s: %v", tc.filename, err)
		}

		desc, _, isMatched := DetectFormat(bytes)
		if isMatched != tc.shouldMatch {
			t.Errorf("fixture %s: expected match %v, got %v", tc.filename, tc.shouldMatch, isMatched)
		}
		if tc.shouldMatch && desc.Type != tc.expectedType {
			t.Errorf("fixture %s: expected type %s, got %s", tc.filename, tc.expectedType, desc.Type)
		}
	}
}
