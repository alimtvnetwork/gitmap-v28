// Package cmd — chromeprofile_register_test.go: verifies CPC registers
// arbitrary destination dirs in Chrome Local State and gitmap's listing.
package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestRegisterChromeProfileCreatesOrderWhenMissing(t *testing.T) {
	root := setupFakeChromeRoot(t, map[string]string{"Profile 15": "Lovable"})
	dst := filepath.Join(root, "lv2")
	if err := os.MkdirAll(dst, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := registerChromeProfileInLocalState("Profile 15", "lv2", "lv2"); err != nil {
		t.Fatalf("register: %v", err)
	}
	state := readRawChromeLocalState(t, root)
	assertChromeLocalStateProfile(t, state, "lv2")
}

func TestAvailableChromeProfileNamesIncludesRegisteredCustomDir(t *testing.T) {
	root := setupFakeChromeRoot(t, map[string]string{"Profile 15": "Lovable"})
	if err := os.MkdirAll(filepath.Join(root, "lv2"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := registerChromeProfileInLocalState("Profile 15", "lv2", "lv2"); err != nil {
		t.Fatalf("register: %v", err)
	}
	assertStringSliceContains(t, availableChromeProfileNames(), "lv2")
}

func readRawChromeLocalState(t *testing.T, root string) map[string]any {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join(root, "Local State"))
	if err != nil {
		t.Fatal(err)
	}
	var out map[string]any
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatal(err)
	}
	return out
}

func assertChromeLocalStateProfile(t *testing.T, state map[string]any, dir string) {
	t.Helper()
	profile := state["profile"].(map[string]any)
	info := profile["info_cache"].(map[string]any)
	entry := info[dir].(map[string]any)
	if entry["name"] != dir || entry["user_name"] != nil {
		t.Fatalf("bad entry: %+v", entry)
	}
	assertChromeProfileOrder(t, profile, dir)
}

func assertChromeProfileOrder(t *testing.T, profile map[string]any, dir string) {
	t.Helper()
	order, ok := profile["profiles_order"].([]any)
	if !ok {
		t.Fatalf("profiles_order missing: %+v", profile)
	}
	assertAnySliceContains(t, order, dir)
}

func assertStringSliceContains(t *testing.T, values []string, want string) {
	t.Helper()
	for _, value := range values {
		if value == want {
			return
		}
	}
	t.Fatalf("%q not found in %v", want, values)
}

func assertAnySliceContains(t *testing.T, values []any, want string) {
	t.Helper()
	for _, value := range values {
		if value == want {
			return
		}
	}
	t.Fatalf("%q not found in %v", want, values)
}
