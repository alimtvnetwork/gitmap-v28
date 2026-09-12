package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestPatchImportedChromeProfilePreferencesScrubAuth(t *testing.T) {
	tempDir := t.TempDir()
	prefPath := filepath.Join(tempDir, "Preferences")

	initialPrefs := map[string]any{
		"account_info": []any{
			map[string]any{
				"email":   "test@gmail.com",
				"gaia_id": "12345",
			},
		},
		"signin": map[string]any{
			"allowed": true,
		},
		"google": map[string]any{
			"services": map[string]any{
				"username": "test@gmail.com",
			},
		},
		"sync": map[string]any{
			"has_setup_completed": true,
		},
		"profile": map[string]any{
			"name":            "Old Name",
			"gaia_name":       "Test User",
			"gaia_given_name": "Test",
		},
	}

	raw, err := json.Marshal(initialPrefs)
	if err != nil {
		t.Fatalf("failed to marshal initial prefs: %v", err)
	}

	if writeErr := os.WriteFile(prefPath, raw, 0644); writeErr != nil {
		t.Fatalf("failed to write initial prefs: %v", writeErr)
	}

	if patchErr := patchImportedChromeProfilePreferences(tempDir, "Clean Profile"); patchErr != nil {
		t.Fatalf("patchImportedChromeProfilePreferences failed: %v", patchErr)
	}

	patchedBytes, readErr := os.ReadFile(prefPath)
	if readErr != nil {
		t.Fatalf("failed to read patched prefs: %v", readErr)
	}

	var patched map[string]any
	if unmarshalErr := json.Unmarshal(patchedBytes, &patched); unmarshalErr != nil {
		t.Fatalf("failed to unmarshal patched prefs: %v", unmarshalErr)
	}

	if _, hasAccountInfo := patched["account_info"]; hasAccountInfo {
		t.Errorf("expected account_info to be scrubbed, but it was found")
	}

	if _, hasGoogle := patched["google"]; hasGoogle {
		t.Errorf("expected google to be scrubbed, but it was found")
	}

	if _, hasSync := patched["sync"]; hasSync {
		t.Errorf("expected sync to be scrubbed, but it was found")
	}

	signinMap, ok := patched["signin"].(map[string]any)
	if !ok || signinMap["allowed"] != false {
		t.Errorf("expected signin.allowed to be false, got: %v", patched["signin"])
	}

	browserMap, ok := patched["browser"].(map[string]any)
	if !ok || browserMap["has_seen_welcome_page"] != true {
		t.Errorf("expected browser.has_seen_welcome_page to be true, got: %v", patched["browser"])
	}

	profMap, ok := patched["profile"].(map[string]any)
	if !ok || profMap["name"] != "Clean Profile" {
		t.Errorf("expected profile.name to be 'Clean Profile', got: %v", patched["profile"])
	}

	if _, hasGaiaName := profMap["gaia_name"]; hasGaiaName {
		t.Errorf("expected gaia_name to be scrubbed from profile")
	}
}

func TestPatchImportedChromeProfilePreferencesKeepSignin(t *testing.T) {
	tempDir := t.TempDir()
	prefPath := filepath.Join(tempDir, "Preferences")

	initialPrefs := map[string]any{
		"account_info": []any{
			map[string]any{
				"email": "keep@gmail.com",
			},
		},
		"profile": map[string]any{
			"name": "Old",
		},
	}

	raw, err := json.Marshal(initialPrefs)
	if err != nil {
		t.Fatalf("failed to marshal initial prefs: %v", err)
	}

	if writeErr := os.WriteFile(prefPath, raw, 0644); writeErr != nil {
		t.Fatalf("failed to write initial prefs: %v", writeErr)
	}

	if patchErr := patchImportedChromeProfilePreferencesWithOptions(tempDir, "Kept Profile", true); patchErr != nil {
		t.Fatalf("patchImportedChromeProfilePreferencesWithOptions failed: %v", patchErr)
	}

	patchedBytes, readErr := os.ReadFile(prefPath)
	if readErr != nil {
		t.Fatalf("failed to read patched prefs: %v", readErr)
	}

	var patched map[string]any
	if unmarshalErr := json.Unmarshal(patchedBytes, &patched); unmarshalErr != nil {
		t.Fatalf("failed to unmarshal patched prefs: %v", unmarshalErr)
	}

	if _, hasAccountInfo := patched["account_info"]; !hasAccountInfo {
		t.Errorf("expected account_info to be kept when keepSignin=true")
	}
}
