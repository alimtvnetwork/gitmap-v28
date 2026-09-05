package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func createTestSnapshotJSON(t *testing.T, dir, filename, name, displayName, email string) string {
	t.Helper()
	exp := chromeExport{
		SchemaVersion: chromeExportSchemaVersion,
		GitMapVersion: constants.Version,
		Name:          name,
		DisplayName:   displayName,
		Email:         email,
		ExportedAt:    time.Now().UTC().Format(time.RFC3339),
		Bookmarks:     json.RawMessage(`{"roots":{"bookmark_bar":{"children":[{"type":"url","name":"Site1","url":"https://example.com"}]}}}`),
		Preferences:   json.RawMessage(`{"account_info":[{"email":"` + email + `"}]}`),
		ExtensionIDs:  []string{"extension-id-1", "extension-id-2"},
	}
	raw, err := json.MarshalIndent(exp, "", "  ")
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	filePath := filepath.Join(dir, filename)
	if err := os.WriteFile(filePath, raw, 0644); err != nil {
		t.Fatalf("write failed: %v", err)
	}
	return filePath
}

func TestChromeProfileImportCheck(t *testing.T) {
	tempUserData := t.TempDir()
	t.Setenv("GITMAP_CHROME_USER_DATA", tempUserData)

	workDir := t.TempDir()
	createTestSnapshotJSON(t, workDir, "Default.json", "Default", "Personal", "personal@gmail.com")
	createTestSnapshotJSON(t, workDir, "Work.json", "Work", "Work Profile", "work@company.com")

	if err := runChromeProfileImportCheck([]string{workDir}); err != nil {
		t.Fatalf("runChromeProfileImportCheck failed: %v", err)
	}
}

func TestChromeProfileImportProtectsExistingProfileAndAllocatesNew(t *testing.T) {
	tempUserData := t.TempDir()
	t.Setenv("GITMAP_CHROME_USER_DATA", tempUserData)

	// Pre-populate Chrome User Data with an existing Default profile
	defaultProfileDir := filepath.Join(tempUserData, "Default")
	if err := os.MkdirAll(defaultProfileDir, 0755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}
	_ = os.WriteFile(filepath.Join(defaultProfileDir, "Bookmarks"), []byte(`dummy-bookmarks-content-over-50-bytes-1234567890-abcdefghij-klmnopqrst`), 0644)

	// Pre-register Default in Local State with personal@gmail.com
	localStatePath := filepath.Join(tempUserData, constants.ChromeLocalStateFile)
	localStateData := map[string]any{
		"profile": map[string]any{
			"info_cache": map[string]any{
				"Default": map[string]any{
					"name":      "Personal Default",
					"user_name": "personal@gmail.com",
				},
			},
			"profiles_order": []any{"Default"},
		},
	}
	rawLS, _ := json.Marshal(localStateData)
	_ = os.WriteFile(localStatePath, rawLS, 0644)

	// Now create a backup snapshot for a completely different email (work@company.com) whose original Name was "Default"
	workDir := t.TempDir()
	snapPath := createTestSnapshotJSON(t, workDir, "Default.json", "Default", "Work Profile", "work@company.com")

	// Import snapshot into Chrome
	if err := runChromeProfileImport([]string{snapPath}); err != nil {
		t.Fatalf("runChromeProfileImport failed: %v", err)
	}

	// Verify that existing Default was NOT overwritten and still has personal@gmail.com
	defaultBM, _ := os.ReadFile(filepath.Join(defaultProfileDir, "Bookmarks"))
	if !strings.Contains(string(defaultBM), "dummy-bookmarks-content") {
		t.Errorf("expected existing Default profile bookmarks to be preserved, got: %s", string(defaultBM))
	}

	// Verify that a new profile directory ("Profile 1") was created for work@company.com
	profile1Dir := filepath.Join(tempUserData, "Profile 1")
	if !chromeProfilePathExists(profile1Dir) {
		t.Fatalf("expected new profile directory Profile 1 to be created")
	}

	// Verify Local State registers Profile 1 with work@company.com
	rawUpdatedLS, err := os.ReadFile(localStatePath)
	if err != nil {
		t.Fatalf("read Local State failed: %v", err)
	}
	var root map[string]any
	_ = json.Unmarshal(rawUpdatedLS, &root)
	prof := root["profile"].(map[string]any)
	infoCache := prof["info_cache"].(map[string]any)
	p1Entry, ok := infoCache["Profile 1"].(map[string]any)
	if !ok {
		t.Fatalf("Profile 1 not registered in info_cache")
	}
	if p1Entry["user_name"] != "work@company.com" {
		t.Errorf("expected user_name work@company.com, got %v", p1Entry["user_name"])
	}
}

func TestChromeProfileImportByEmail(t *testing.T) {
	tempUserData := t.TempDir()
	t.Setenv("GITMAP_CHROME_USER_DATA", tempUserData)

	workDir := t.TempDir()
	oldWd, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(oldWd) }()

	createTestSnapshotJSON(t, workDir, "Personal.json", "Personal", "Personal", "personal@gmail.com")
	createTestSnapshotJSON(t, workDir, "Work.json", "Work", "Work", "target.work@corp.com")

	// Run import targeting single email
	if err := runChromeProfileImport([]string{"target.work@corp.com"}); err != nil {
		t.Fatalf("import by email failed: %v", err)
	}

	// Verify only target.work@corp.com was imported
	names := availableChromeProfileNames()
	if len(names) != 1 {
		t.Fatalf("expected exactly 1 imported profile, got %d: %v", len(names), names)
	}
	_, email := resolveProfileNameAndEmail(names[0], nil)
	if email != "target.work@corp.com" {
		t.Errorf("expected email target.work@corp.com, got %s", email)
	}
}

func TestChromeProfileImportWithExceptAndLimit(t *testing.T) {
	tempUserData := t.TempDir()
	t.Setenv("GITMAP_CHROME_USER_DATA", tempUserData)

	workDir := t.TempDir()
	createTestSnapshotJSON(t, workDir, "ProfA.json", "ProfA", "ProfA", "a@test.com")
	createTestSnapshotJSON(t, workDir, "ProfB.json", "ProfB", "ProfB", "b@test.com")
	createTestSnapshotJSON(t, workDir, "ProfC.json", "ProfC", "ProfC", "c@test.com")

	// Test import directory with --except "ProfA,b@test.com" and --limit 1
	err := runChromeProfileImport([]string{workDir, "--except=ProfA,b@test.com", "--limit=1"})
	if err != nil {
		t.Fatalf("import with except and limit failed: %v", err)
	}

	// Only ProfC should have been imported
	names := availableChromeProfileNames()
	if len(names) != 1 {
		t.Fatalf("expected 1 profile imported, got %d: %v", len(names), names)
	}
	_, email := resolveProfileNameAndEmail(names[0], nil)
	if email != "c@test.com" {
		t.Errorf("expected email c@test.com, got %s", email)
	}
}

func TestChromeProfileListDiscoveredSnapshots(t *testing.T) {
	tempUserData := t.TempDir()
	t.Setenv("GITMAP_CHROME_USER_DATA", tempUserData)

	workDir := t.TempDir()
	oldWd, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(oldWd) }()

	createTestSnapshotJSON(t, workDir, "Sample.json", "Sample", "Sample Display", "sample@domain.com")

	if err := runChromeProfileList([]string{}); err != nil {
		t.Fatalf("runChromeProfileList failed: %v", err)
	}
}

func TestChromeImportAllSmartDelegation(t *testing.T) {
	tempUserData := t.TempDir()
	t.Setenv("GITMAP_CHROME_USER_DATA", tempUserData)

	// Pre-populate existing Default profile
	defaultProfileDir := filepath.Join(tempUserData, "Default")
	_ = os.MkdirAll(defaultProfileDir, 0755)
	_ = os.WriteFile(filepath.Join(defaultProfileDir, "Bookmarks"), []byte(`dummy-bookmarks-content-over-50-bytes-1234567890-abcdefghij-klmnopqrst`), 0644)

	workDir := t.TempDir()
	createTestSnapshotJSON(t, workDir, "Default.json", "Default", "Exported Default", "newaccount@corp.com")

	if err := runChromeImportAll([]string{workDir}); err != nil {
		t.Fatalf("runChromeImportAll failed: %v", err)
	}

	// Default bookmarks must not be overwritten
	defaultBM, _ := os.ReadFile(filepath.Join(defaultProfileDir, "Bookmarks"))
	if !strings.Contains(string(defaultBM), "dummy-bookmarks-content") {
		t.Errorf("expected existing Default profile bookmarks to be preserved")
	}

	// Profile 1 must be created
	profile1Dir := filepath.Join(tempUserData, "Profile 1")
	if !chromeProfilePathExists(profile1Dir) {
		t.Fatalf("expected new profile directory Profile 1 to be created by import-all")
	}
}

func TestChromeImportCurrentDirectoryDot(t *testing.T) {
	tempUserData := t.TempDir()
	t.Setenv("GITMAP_CHROME_USER_DATA", tempUserData)

	workDir := t.TempDir()
	oldWd, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(oldWd) }()

	createTestSnapshotJSON(t, workDir, "Default.json", "Default", "My Profile", "dotuser@example.com")

	// Test importing "."
	if err := runChromeProfileImport([]string{"."}); err != nil {
		t.Fatalf("runChromeProfileImport with . failed: %v", err)
	}

	// Test import-check "ls"
	if err := runChromeProfileImportCheck([]string{"ls"}); err != nil {
		t.Fatalf("runChromeProfileImportCheck with ls failed: %v", err)
	}
}
