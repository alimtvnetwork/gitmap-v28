package cmdchromeprofile

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestChromeImportCheckFileExportJSON(t *testing.T) {
	tempUserData := t.TempDir()
	t.Setenv("GITMAP_CHROME_USER_DATA", tempUserData)

	workDir := t.TempDir()
	createTestSnapshotJSON(t, workDir, "Default.json", "Default", "Personal", "test@test.com")
	createTestSnapshotJSON(t, workDir, "Profile 1.json", "Profile 1", "Work", "work@test.com")

	outDir := t.TempDir()
	outFile := filepath.Join(outDir, "sub", "preview.json")

	args := []string{workDir, "--json", "--file", outFile}
	if err := runChromeProfileImportCheck(args); err != nil {
		t.Fatalf("runChromeProfileImportCheck with --file failed: %v", err)
	}

	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("failed reading exported preview file: %v", err)
	}

	var candidates []DiscoveredProfileCandidate
	if err := json.Unmarshal(data, &candidates); err != nil {
		t.Fatalf("unmarshal exported JSON failed: %v\nContent: %s", err, string(data))
	}

	if len(candidates) != 2 {
		t.Errorf("expected 2 candidates in exported JSON, got %d", len(candidates))
	}
}

func TestChromeImportCheckFnfFlag(t *testing.T) {
	tempUserData := t.TempDir()
	t.Setenv("GITMAP_CHROME_USER_DATA", tempUserData)

	workDir := t.TempDir()
	createTestSnapshotJSON(t, workDir, "Default.json", "Default", "Personal", "test@test.com")

	outDir := t.TempDir()
	fnfOutFile := filepath.Join(outDir, "fnf_out.json")

	args := []string{workDir, "--json", "--fnf", fnfOutFile}
	if err := runChromeProfileImportCheck(args); err != nil {
		t.Fatalf("runChromeProfileImportCheck with --fnf failed: %v", err)
	}

	data, err := os.ReadFile(fnfOutFile)
	if err != nil {
		t.Fatalf("failed reading fnf exported file: %v", err)
	}

	var candidates []DiscoveredProfileCandidate
	if err := json.Unmarshal(data, &candidates); err != nil {
		t.Fatalf("unmarshal fnf JSON failed: %v", err)
	}

	if len(candidates) != 1 {
		t.Errorf("expected 1 candidate, got %d", len(candidates))
	}
}

func TestChromeImportCheckFnfFailureOnEmpty(t *testing.T) {
	emptyDir := t.TempDir()
	args := []string{emptyDir, "--fnf"}

	err := runChromeProfileImportCheck(args)
	if err == nil {
		t.Fatalf("expected error when --fnf asserted on empty directory, got nil")
	}

	if !strings.Contains(err.Error(), "--fnf asserted") {
		t.Errorf("expected error message to mention '--fnf asserted', got: %v", err)
	}
}

func TestChromeImportCheckTempFile(t *testing.T) {
	tempUserData := t.TempDir()
	t.Setenv("GITMAP_CHROME_USER_DATA", tempUserData)

	workDir := t.TempDir()
	createTestSnapshotJSON(t, workDir, "Default.json", "Default", "Personal", "test@test.com")

	tempFileName := "test-preview-candidate.json"
	args := []string{workDir, "--json", "--tempfile", tempFileName}

	if err := runChromeProfileImportCheck(args); err != nil {
		t.Fatalf("runChromeProfileImportCheck with --tempfile failed: %v", err)
	}

	expectedPath := filepath.Join(resolveTempDir(), tempFileName)
	defer os.Remove(expectedPath)

	data, err := os.ReadFile(expectedPath)
	if err != nil {
		t.Fatalf("expected temp file at %s: %v", expectedPath, err)
	}

	var candidates []DiscoveredProfileCandidate
	if err := json.Unmarshal(data, &candidates); err != nil {
		t.Fatalf("unmarshal tempfile JSON failed: %v", err)
	}

	if len(candidates) != 1 {
		t.Errorf("expected 1 candidate in tempfile, got %d", len(candidates))
	}
}
