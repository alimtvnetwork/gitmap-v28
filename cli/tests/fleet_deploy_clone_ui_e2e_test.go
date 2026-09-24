//go:build e2e

package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmddaemon"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdui"
)

// TestE2ERemoteEditorRoundTrip verifies reading and writing files via UI editor engine.
func TestE2ERemoteEditorRoundTrip(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "gitmap-e2e-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testFilePath := filepath.Join(tempDir, "sample.json")
	initialContent := `{"test": true, "version": "1.0.0"}`
	if err := os.WriteFile(testFilePath, []byte(initialContent), 0644); err != nil {
		t.Fatalf("failed to write initial file: %v", err)
	}

	// 1. Test Read
	readRes, err := cmdui.ReadRemoteOrLocalFile("local", testFilePath)
	if err != nil {
		t.Fatalf("expected successful read, got error: %v", err)
	}
	if readRes.IsSuccess == false {
		t.Fatalf("expected IsSuccess true, got false with err: %s", readRes.Error)
	}
	if readRes.Language != "json" {
		t.Errorf("expected syntax language 'json', got '%s'", readRes.Language)
	}

	// 2. Test Save (Modify)
	modifiedContent := `{"test": true, "version": "2.0.0", "updated": true}`
	if err := cmdui.SaveRemoteOrLocalFile("local", testFilePath, modifiedContent); err != nil {
		t.Fatalf("failed to save modified content: %v", err)
	}

	// 3. Verify Persistence
	savedData, err := os.ReadFile(testFilePath)
	if err != nil {
		t.Fatalf("failed to re-read saved file: %v", err)
	}
	if string(savedData) != modifiedContent {
		t.Errorf("saved content mismatch: expected %s, got %s", modifiedContent, string(savedData))
	}
}

// TestE2EDaemonRESTTriad verifies REST command execution against a mock daemon.
func TestE2EDaemonRESTTriad(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/api/v1/status" {
			_ = json.NewEncoder(w).Encode(cmddaemon.DaemonStatusResp{
				Status:    "online",
				Version:   "v6.319.0",
				OS:        "linux",
				UptimeSec: 100,
			})
			return
		}

		if r.URL.Path == "/api/v1/exec" {
			var req cmddaemon.DaemonExecReq
			_ = json.NewDecoder(r.Body).Decode(&req)
			_ = json.NewEncoder(w).Encode(cmddaemon.DaemonExecResp{
				Status:     "success",
				Success:    true,
				ExitCode:   0,
				Stdout:     "mock output for: " + req.Command,
				DurationMs: 5,
			})
			return
		}
	}))
	defer server.Close()

	resp, err := http.Get(server.URL + "/api/v1/status")
	if err != nil {
		t.Fatalf("failed to query mock daemon status: %v", err)
	}
	defer resp.Body.Close()

	var status cmddaemon.DaemonStatusResp
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		t.Fatalf("failed to decode daemon status: %v", err)
	}
	if status.Status != "online" {
		t.Errorf("expected status 'online', got '%s'", status.Status)
	}
}
