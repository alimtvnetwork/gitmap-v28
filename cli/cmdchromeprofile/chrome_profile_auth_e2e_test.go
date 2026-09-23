//go:build e2e

package cmdchromeprofile

import (
	"archive/zip"
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestChromeProfileAuthSessionRoundtrip(t *testing.T) {
	srcDir := resolveSourceTestProfile(t)
	tmpDir := t.TempDir()
	snapFile := filepath.Join(tmpDir, "profile_export.json")

	exp := executeAndVerifyExport(t, srcDir, snapFile)
	destDir := filepath.Join(tmpDir, "imported_profile")
	if err := applyChromeExport(&exp, destDir); err != nil {
		t.Fatalf("applyChromeExport failed: %v", err)
	}

	verifyImportedPreferencesAuth(t, destDir, exp.Email)
	verifyImportedTokenService(t, srcDir, destDir)
	verifyImportedCookiesDatabase(t, srcDir, destDir)
}

func TestChromeProfileAuthZipRoundtrip(t *testing.T) {
	srcDir := resolveSourceTestProfile(t)
	tmpDir := t.TempDir()
	zipFile := filepath.Join(tmpDir, "profile_export.zip")

	name := filepath.Base(srcDir)
	if _, err := writeChromeExportZIP(srcDir, name, zipFile); err != nil {
		t.Fatalf("writeChromeExportZIP failed: %v", err)
	}

	r, err := zip.OpenReader(zipFile)
	if err != nil {
		t.Fatalf("open zip failed: %v", err)
	}
	defer r.Close()

	destDir := filepath.Join(tmpDir, "imported_zip_profile")
	if err := extractSingleProfileZip(r, destDir); err != nil {
		t.Fatalf("extractSingleProfileZip failed: %v", err)
	}

	verifyImportedPreferencesAuth(t, destDir, "")
	verifyImportedTokenService(t, srcDir, destDir)
	verifyImportedCookiesDatabase(t, srcDir, destDir)
}

func resolveSourceTestProfile(t *testing.T) string {
	t.Helper()
	root := chromeUserDataDir()
	candidate := filepath.Join(root, "Profile 1")
	if !hasChromeProfileAuth(candidate) {
		t.Skip("Live Profile 1 with auth tokens not found on this VM; skipping e2e test")
	}

	return candidate
}

func hasChromeProfileAuth(profileDir string) bool {
	webData := filepath.Join(profileDir, "Web Data")
	cookies := filepath.Join(profileDir, "Network", "Cookies")
	hasWebData := fileExistsAndNonEmpty(webData)
	hasCookies := fileExistsAndNonEmpty(cookies)

	return hasWebData && hasCookies
}

func fileExistsAndNonEmpty(path string) bool {
	info, err := os.Stat(path)

	return err == nil && info.Size() > 0
}

func executeAndVerifyExport(t *testing.T, srcDir, snapFile string) chromeExport {
	t.Helper()
	name := filepath.Base(srcDir)
	if _, err := writeChromeExport(srcDir, name, snapFile); err != nil {
		t.Fatalf("writeChromeExport failed: %v", err)
	}

	raw, err := os.ReadFile(snapFile)
	if err != nil {
		t.Fatalf("read snapshot failed: %v", err)
	}

	var exp chromeExport
	if err := json.Unmarshal(raw, &exp); err != nil {
		t.Fatalf("parse snapshot failed: %v", err)
	}

	validateExportAuthFields(t, exp)

	return exp
}

func validateExportAuthFields(t *testing.T, exp chromeExport) {
	t.Helper()
	if exp.GaiaID == "" {
		t.Errorf("expected GaiaID to be populated in export")
	}
	if exp.TokenVault == nil || exp.TokenVault.Count == 0 {
		t.Errorf("expected TokenVault with at least 1 token")
	}
	if exp.CookiesRawBase64 == "" {
		t.Errorf("expected CookiesRawBase64 to be populated")
	}
	if exp.WebDataRawBase64 == "" {
		t.Errorf("expected WebDataRawBase64 to be populated")
	}
}

func verifyImportedPreferencesAuth(t *testing.T, destDir, expectedEmail string) {
	t.Helper()
	prefPath := filepath.Join(destDir, "Preferences")
	raw, err := os.ReadFile(prefPath)
	if err != nil {
		t.Fatalf("read Preferences: %v", err)
	}

	var root map[string]any
	if err := json.Unmarshal(raw, &root); err != nil {
		t.Fatalf("parse Preferences: %v", err)
	}

	signin, _ := root["signin"].(map[string]any)
	if signin == nil || signin["allowed"] != true {
		t.Errorf("expected signin.allowed == true in imported Preferences")
	}

	accs, _ := root["account_info"].([]any)
	if len(accs) == 0 {
		t.Errorf("expected account_info to be preserved in imported Preferences")
	}
}

func verifyImportedTokenService(t *testing.T, srcDir, destDir string) {
	t.Helper()
	srcTokens := readSQLiteTokens(t, filepath.Join(srcDir, "Web Data"))
	destTokens := readSQLiteTokens(t, filepath.Join(destDir, "Web Data"))

	if len(destTokens) == 0 {
		t.Fatalf("no tokens found in imported Web Data token_service")
	}

	for k, srcVal := range srcTokens {
		destVal, ok := destTokens[k]
		if !ok || string(srcVal) != string(destVal) {
			t.Errorf("token mismatch for %s: expected %q, got %q", k, srcVal, destVal)
		}
	}
}

func readSQLiteTokens(t *testing.T, dbPath string) map[string][]byte {
	t.Helper()
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("open %s: %v", dbPath, err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT service, encrypted_token FROM token_service")
	if err != nil {
		return map[string][]byte{}
	}
	defer rows.Close()

	tokens := make(map[string][]byte)
	for rows.Next() {
		var s string
		var b []byte
		if err := rows.Scan(&s, &b); err == nil {
			tokens[s] = b
		}
	}

	return tokens
}

func verifyImportedCookiesDatabase(t *testing.T, srcDir, destDir string) {
	t.Helper()
	destCookiePath := filepath.Join(destDir, "Network", "Cookies")
	info, err := os.Stat(destCookiePath)
	if err != nil || info.Size() == 0 {
		t.Fatalf("restored Network/Cookies is missing or empty")
	}

	db, err := sql.Open("sqlite", destCookiePath)
	if err != nil {
		t.Fatalf("open restored cookies db: %v", err)
	}
	defer db.Close()

	var count int
	_ = db.QueryRow("SELECT COUNT(*) FROM cookies").Scan(&count)
	if count == 0 {
		t.Errorf("expected restored cookies table to have cookies, got 0")
	}
}
