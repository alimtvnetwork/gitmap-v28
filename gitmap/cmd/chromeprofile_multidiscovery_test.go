package cmd

import (
	"archive/zip"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func TestChromeImportAllMultiProfileDirectoryWithManifest(t *testing.T) {
	tempUserData := t.TempDir()
	t.Setenv("GITMAP_CHROME_USER_DATA", tempUserData)

	workDir := t.TempDir()
	manifest := chromeProfileManifest{
		GitMapVersion: constants.Version,
		SchemaVersion: 1,
		ExportedAt:    time.Now().UTC().Format(time.RFC3339),
		ProfileCount:  2,
		Profiles: []chromeManifestProfile{
			{Name: "Profile 1", DisplayName: "Roki", Email: "roki@gmail.com", ExtensionCount: 2},
			{Name: "Profile 2", DisplayName: "Alex", Email: "alex@gmail.com", ExtensionCount: 1},
		},
	}
	mRaw, _ := json.MarshalIndent(manifest, "", "  ")
	_ = os.WriteFile(filepath.Join(workDir, "manifest.json"), mRaw, 0644)

	p1Dir := filepath.Join(workDir, "Profile 1")
	_ = os.MkdirAll(p1Dir, 0755)
	createTestSnapshotJSON(t, p1Dir, "Profile 1.json", "Profile 1", "Roki", "roki@gmail.com")

	p2Dir := filepath.Join(workDir, "Profile 2")
	_ = os.MkdirAll(p2Dir, 0755)
	createTestSnapshotJSON(t, p2Dir, "Profile 2.json", "Profile 2", "Alex", "alex@gmail.com")

	if err := runChromeImportAll([]string{workDir}); err != nil {
		t.Fatalf("runChromeImportAll failed: %v", err)
	}

	manifestProfile := filepath.Join(tempUserData, "manifest")
	if chromeProfilePathExists(manifestProfile) {
		t.Fatalf("BUG: directory manifest.json was incorrectly imported as dummy profile 'manifest'")
	}

	if !chromeProfilePathExists(filepath.Join(tempUserData, "Profile 1")) {
		t.Errorf("expected Profile 1 to be imported")
	}
	if !chromeProfilePathExists(filepath.Join(tempUserData, "Profile 2")) {
		t.Errorf("expected Profile 2 to be imported")
	}
}

func TestChromeImportAllMultiProfileZip(t *testing.T) {
	tempUserData := t.TempDir()
	t.Setenv("GITMAP_CHROME_USER_DATA", tempUserData)

	workDir := t.TempDir()
	zipPath := filepath.Join(workDir, "chrome-ext.zip")

	zipFile, err := os.Create(zipPath)
	if err != nil {
		t.Fatalf("create zip failed: %v", err)
	}

	zw := zip.NewWriter(zipFile)
	manifest := chromeProfileManifest{
		GitMapVersion: constants.Version,
		SchemaVersion: 1,
		ProfileCount:  2,
		Profiles: []chromeManifestProfile{
			{Name: "Profile 1", DisplayName: "Roki", Email: "roki@gmail.com"},
			{Name: "Profile 2", DisplayName: "Alex", Email: "alex@gmail.com"},
		},
	}
	mRaw, _ := json.Marshal(manifest)
	mw, _ := zw.Create("manifest.json")
	_, _ = mw.Write(mRaw)

	p1w, _ := zw.Create("Profile 1/Bookmarks")
	_, _ = p1w.Write([]byte(`{"roots":{"bookmark_bar":{"children":[{"type":"url","name":"G","url":"https://google.com"}]}}}`))

	p2w, _ := zw.Create("Profile 2/Bookmarks")
	_, _ = p2w.Write([]byte(`{"roots":{"bookmark_bar":{"children":[{"type":"url","name":"Y","url":"https://yahoo.com"}]}}}`))

	_ = zw.Close()
	_ = zipFile.Close()

	if err := runChromeImportAll([]string{zipPath}); err != nil {
		t.Fatalf("runChromeImportAll with zip failed: %v", err)
	}

	zipDummyProfile := filepath.Join(tempUserData, "chrome-ext")
	if chromeProfilePathExists(zipDummyProfile) {
		t.Fatalf("BUG: zip archive was imported as a single dummy profile named after the zip")
	}

	if !chromeProfilePathExists(filepath.Join(tempUserData, "Profile 1")) {
		t.Errorf("expected Profile 1 to be extracted and imported from zip")
	}
	if !chromeProfilePathExists(filepath.Join(tempUserData, "Profile 2")) {
		t.Errorf("expected Profile 2 to be extracted and imported from zip")
	}
}

func TestChromeImportLsPreflightCommand(t *testing.T) {
	tempUserData := t.TempDir()
	t.Setenv("GITMAP_CHROME_USER_DATA", tempUserData)

	workDir := t.TempDir()
	createTestSnapshotJSON(t, workDir, "Default.json", "Default", "Personal", "test@test.com")

	if err := runChromeProfileImport([]string{"ls", workDir}); err != nil {
		t.Errorf("gitmap chrome import ls failed: %v", err)
	}

	if err := runChromeImportAll([]string{"ls", workDir}); err != nil {
		t.Errorf("gitmap chrome import-all ls failed: %v", err)
	}

	if err := runChromeProfileImportCheck([]string{workDir, "--json"}); err != nil {
		t.Errorf("gitmap chrome import-check --json failed: %v", err)
	}
}

func TestTopLevelImportAllPreflight(t *testing.T) {
	workDir := t.TempDir()
	createTestSnapshotJSON(t, workDir, "Default.json", "Default", "Personal", "test@test.com")

	if err := runImportAll([]string{"scan", workDir}); err != nil {
		t.Errorf("gitmap import-all scan failed: %v", err)
	}
	if err := runImportAll([]string{"profile", workDir}); err != nil {
		t.Errorf("gitmap import-all profile failed: %v", err)
	}
}

func TestTopLevelImportAllExecution(t *testing.T) {
	t.Setenv("GITMAP_CHROME_USER_DATA", t.TempDir())
	workDir := t.TempDir()
	createTestSnapshotJSON(t, workDir, "Default.json", "Default", "Personal", "test@test.com")

	if err := runImportAll([]string{workDir}); err != nil {
		t.Errorf("gitmap import-all <workDir> failed: %v", err)
	}
}
