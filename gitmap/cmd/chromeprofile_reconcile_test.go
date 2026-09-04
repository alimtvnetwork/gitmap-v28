// Package cmd — chromeprofile_reconcile_test.go: tests for reconciling
// unlinked on-disk profile directories into Chrome Local State.
package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func TestChromeProfileReconcileUnlinkedProfile(t *testing.T) {
	tempUserData := t.TempDir()
	t.Setenv("GITMAP_CHROME_USER_DATA", tempUserData)

	// Local State starts with only "Default"
	localStatePath := filepath.Join(tempUserData, constants.ChromeLocalStateFile)
	localStateData := map[string]any{
		"profile": map[string]any{
			"info_cache": map[string]any{
				"Default": map[string]any{
					"name":      "Main Account",
					"user_name": "main@example.com",
				},
			},
			"profiles_order": []any{"Default"},
		},
	}
	rawLS, _ := json.Marshal(localStateData)
	if err := os.WriteFile(localStatePath, rawLS, 0644); err != nil {
		t.Fatal(err)
	}

	// Create Default folder on disk
	_ = os.MkdirAll(filepath.Join(tempUserData, "Default"), 0755)

	// Create an unlinked "Profile 5" on disk with Preferences
	profile5Dir := filepath.Join(tempUserData, "Profile 5")
	if err := os.MkdirAll(profile5Dir, 0755); err != nil {
		t.Fatal(err)
	}
	prefsContent := `{
		"account_info": [
			{
				"email": "erfan.office.n@gmail.com",
				"full_name": "Erfan Office"
			}
		],
		"profile": {
			"name": "Person 1"
		}
	}`
	if err := os.WriteFile(filepath.Join(profile5Dir, constants.ChromePreferencesFile), []byte(prefsContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Run reconcile
	if err := runChromeProfileReconcile([]string{}); err != nil {
		t.Fatalf("runChromeProfileReconcile failed: %v", err)
	}

	// Verify Local State now contains Profile 5 with full schema
	updatedRaw, err := os.ReadFile(localStatePath)
	if err != nil {
		t.Fatal(err)
	}
	var root map[string]any
	if err := json.Unmarshal(updatedRaw, &root); err != nil {
		t.Fatal(err)
	}

	prof := root["profile"].(map[string]any)
	infoCache := prof["info_cache"].(map[string]any)

	p5Entry, ok := infoCache["Profile 5"].(map[string]any)
	if !ok {
		t.Fatalf("expected Profile 5 to be registered in info_cache")
	}

	if p5Entry["name"] != "Erfan Office" {
		t.Errorf("expected name 'Erfan Office', got %v", p5Entry["name"])
	}
	if p5Entry["user_name"] != "erfan.office.n@gmail.com" {
		t.Errorf("expected user_name 'erfan.office.n@gmail.com', got %v", p5Entry["user_name"])
	}
	if p5Entry["avatar_icon"] != "chrome://theme/IDR_PROFILE_AVATAR_26" {
		t.Errorf("expected default avatar_icon, got %v", p5Entry["avatar_icon"])
	}

	// Verify Profile 5 in profiles_order
	assertChromeProfileOrder(t, prof, "Profile 5")

	// Verify Preferences was patched
	updatedPrefsRaw, _ := os.ReadFile(filepath.Join(profile5Dir, constants.ChromePreferencesFile))
	var updatedPrefs map[string]any
	_ = json.Unmarshal(updatedPrefsRaw, &updatedPrefs)
	updatedProf := updatedPrefs["profile"].(map[string]any)
	if updatedProf["name"] != "Erfan Office" {
		t.Errorf("expected patched preferences profile.name 'Erfan Office', got %v", updatedProf["name"])
	}
}
