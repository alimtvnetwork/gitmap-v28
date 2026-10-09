package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/jsonx"
)

func TestWhichFormat_SingleAndMultiFileInspection(t *testing.T) {
	tempDir := t.TempDir()

	// 1. Valid SSH Nodes JSON (Envelope)
	sshEnv := jsonx.NewEnvelope(jsonx.TypeSSHNodes, "repo-secrets/test.json", "gitmap ssh export", "1.0", map[string]any{
		"nodes": []map[string]any{
			{"alias": "node-1", "ip_address": "192.168.1.10"},
		},
	})
	sshFile := filepath.Join(tempDir, "nodes.json")
	sshBytes, _ := json.MarshalIndent(sshEnv, "", "  ")
	_ = os.WriteFile(sshFile, sshBytes, 0644)

	// 2. Valid Macro JSON (Envelope)
	macroEnv := jsonx.NewEnvelope(jsonx.TypeMacro, "test/macro.json", "gitmap macro export", "1.0", map[string]any{
		"name": "build-flow",
	})
	macroFile := filepath.Join(tempDir, "macro.json")
	macroBytes, _ := json.MarshalIndent(macroEnv, "", "  ")
	_ = os.WriteFile(macroFile, macroBytes, 0644)

	// 3. Unmatched Arbitrary JSON
	unmatchedFile := filepath.Join(tempDir, "random.json")
	_ = os.WriteFile(unmatchedFile, []byte(`{"random_data": true}`), 0644)

	results := inspectJSONFiles([]string{sshFile, macroFile, unmatchedFile})
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}

	// First file: SSH Nodes
	if !results[0].IsMatched || results[0].Type != jsonx.TypeSSHNodes {
		t.Errorf("file 0 expected matched %s, got matched=%v, type=%s",
			jsonx.TypeSSHNodes, results[0].IsMatched, results[0].Type)
	}
	if !strings.Contains(results[0].SuggestedImportCmd, "gitmap sj import") {
		t.Errorf("expected suggested import command with 'gitmap sj import', got: %s", results[0].SuggestedImportCmd)
	}

	// Second file: Macro
	if !results[1].IsMatched || results[1].Type != jsonx.TypeMacro {
		t.Errorf("file 1 expected matched %s, got matched=%v, type=%s",
			jsonx.TypeMacro, results[1].IsMatched, results[1].Type)
	}
	if !strings.Contains(results[1].SuggestedImportCmd, "gitmap macro import") {
		t.Errorf("expected suggested import command with 'gitmap macro import', got: %s", results[1].SuggestedImportCmd)
	}

	// Third file: Unmatched
	if results[2].IsMatched {
		t.Errorf("expected file 2 to be unmatched, but IsMatched is true")
	}
	if !strings.Contains(results[2].UnmatchedReason, "Could not match") {
		t.Errorf("expected unmatched reason, got: %s", results[2].UnmatchedReason)
	}
}

func TestWhichFormat_ScanDirectory(t *testing.T) {
	tempDir := t.TempDir()

	uiFile := filepath.Join(tempDir, "ui.json")
	uiContent := []byte(`{"graphicsMode": "high", "commitInLayout": "split"}`)
	_ = os.WriteFile(uiFile, uiContent, 0644)

	files, err := scanDirectoryJSONFiles(tempDir)
	if err != nil {
		t.Fatalf("scanDirectoryJSONFiles failed: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("expected 1 file found, got %d", len(files))
	}

	results := inspectJSONFiles(files)
	if len(results) != 1 || !results[0].IsMatched || results[0].Type != jsonx.TypeUISettings {
		t.Errorf("expected ui-settings matched, got %+v", results)
	}
}
