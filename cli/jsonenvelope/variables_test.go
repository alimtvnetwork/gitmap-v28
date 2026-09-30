package jsonenvelope

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestExpandVariables_Substitution(t *testing.T) {
	input := []byte(`{"path": "${keyPath}", "dir": "$variables.workDir", "nested": "${variables.name}"}`)
	vars := map[string]any{
		"keyPath": `C:\Users\Administrator\.ssh\id_rsa`,
		"workDir": `D:\work`,
		"name":    "gitmap-node",
	}

	output := ExpandVariables(input, vars)
	outStr := string(output)

	if !strings.Contains(outStr, `C:\\Users\\Administrator\\.ssh\\id_rsa`) {
		t.Errorf("expected escaped backslashes in output, got: %s", outStr)
	}
	if !strings.Contains(outStr, `D:\\work`) {
		t.Errorf("expected workDir expanded, got: %s", outStr)
	}
	if !strings.Contains(outStr, `gitmap-node`) {
		t.Errorf("expected name expanded, got: %s", outStr)
	}
}

func TestEnvelopeAttributes_WorkDirectoryObjectDuality(t *testing.T) {
	// 1. Scalar string
	jsonString := []byte(`{
		"type": "ssh-nodes",
		"workDirectory": "D:\\work",
		"defaultWorkDirectory": "D:\\work"
	}`)
	var attrs1 EnvelopeAttributes
	if err := json.Unmarshal(jsonString, &attrs1); err != nil {
		t.Fatalf("unmarshal string workDirectory failed: %v", err)
	}
	if attrs1.WorkDirectory != `D:\work` {
		t.Errorf("expected WorkDirectory D:\\work, got %s", attrs1.WorkDirectory)
	}

	// 2. Object form
	jsonObject := []byte(`{
		"type": "ssh-nodes",
		"workDirectory": {
			"path": "${workDir}",
			"defaultPath": "D:\\work",
			"isApplied": true,
			"isEnforced": false,
			"variables": {
				"workDir": "D:\\work"
			}
		}
	}`)
	var attrs2 EnvelopeAttributes
	if err := json.Unmarshal(jsonObject, &attrs2); err != nil {
		t.Fatalf("unmarshal object workDirectory failed: %v", err)
	}
	if attrs2.WorkDirectoryConfig == nil {
		t.Fatalf("expected WorkDirectoryConfig to be non-nil")
	}
	if attrs2.WorkDirectory != "${workDir}" {
		t.Errorf("expected WorkDirectory ${workDir}, got %s", attrs2.WorkDirectory)
	}
	if !attrs2.IsWorkDirectoryApplied {
		t.Errorf("expected IsWorkDirectoryApplied to be true")
	}
}
