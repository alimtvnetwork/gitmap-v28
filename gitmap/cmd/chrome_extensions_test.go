// Package cmd — chrome_extensions_test.go: unit tests for Chrome extension
// scanning, enable/disable pattern filtering, and installation.
package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func TestChromeExtensionLifecycle(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("GITMAP_CHROME_USER_DATA", tempDir)
	profDir := filepath.Join(tempDir, "Default")
	_ = os.MkdirAll(profDir, constants.DirPermission)

	// Seed Preferences with an extension
	prefData := map[string]any{
		"extensions": map[string]any{
			"settings": map[string]any{
				"mock_ext_123": map[string]any{
					"state": float64(1),
				},
			},
		},
	}
	prefRaw, _ := json.Marshal(prefData)
	_ = os.WriteFile(filepath.Join(profDir, "Preferences"), prefRaw, constants.FilePermission)

	// Seed Extensions folder
	extVerDir := filepath.Join(profDir, "Extensions", "mock_ext_123", "1.0.0")
	_ = os.MkdirAll(extVerDir, constants.DirPermission)
	manifestData := `{"name": "Mock Test Extension", "version": "1.0.0", "description": "For testing"}`
	_ = os.WriteFile(filepath.Join(extVerDir, "manifest.json"), []byte(manifestData), constants.FilePermission)

	// 1. Test scanning & listing
	exts, err := scanExtensionsForProfile("Default")
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}
	if len(exts) != 1 {
		t.Fatalf("expected 1 extension, got %d", len(exts))
	}
	if exts[0].Name != "Mock Test Extension" || !exts[0].IsEnabled {
		t.Errorf("unexpected extension info: %+v", exts[0])
	}

	// 2. Test disable by pattern
	if err := runChromeExtensionDisable([]string{"Mock", "--profile=Default"}); err != nil {
		t.Fatalf("disable failed: %v", err)
	}
	extsAfterDisable, _ := scanExtensionsForProfile("Default")
	if extsAfterDisable[0].IsEnabled {
		t.Errorf("expected extension to be disabled")
	}

	// 3. Test enable
	if err := runChromeExtensionEnable([]string{"mock_ext_123", "--profile=Default"}); err != nil {
		t.Fatalf("enable failed: %v", err)
	}
	extsAfterEnable, _ := scanExtensionsForProfile("Default")
	if !extsAfterEnable[0].IsEnabled {
		t.Errorf("expected extension to be enabled")
	}

	// 4. Test disable all
	if err := runChromeExtensionDisableAll([]string{"--profile=Default"}); err != nil {
		t.Fatalf("disable all failed: %v", err)
	}
	extsAfterDisableAll, _ := scanExtensionsForProfile("Default")
	if extsAfterDisableAll[0].IsEnabled {
		t.Errorf("expected extension to be disabled")
	}
}

func TestChromeExtensionInstall(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("GITMAP_CHROME_USER_DATA", tempDir)
	profDir := filepath.Join(tempDir, "Default")
	_ = os.MkdirAll(profDir, constants.DirPermission)

	// Create a dummy unpacked extension folder
	unpackedDir := filepath.Join(tempDir, "my-extension")
	_ = os.MkdirAll(unpackedDir, constants.DirPermission)
	_ = os.WriteFile(filepath.Join(unpackedDir, "manifest.json"), []byte(`{"name": "Injected Ext", "version": "2.0.0"}`), constants.FilePermission)

	if err := runChromeExtensionInstall([]string{unpackedDir, "--profile=Default"}); err != nil {
		t.Fatalf("install failed: %v", err)
	}
	exts, _ := scanExtensionsForProfile("Default")
	if len(exts) == 0 {
		t.Fatalf("expected installed extension to be discovered")
	}
}
