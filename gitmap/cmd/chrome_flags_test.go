// Package cmd — chrome_flags_test.go: unit tests for inspecting, enabling,
// disabling, and resetting Chrome experimental feature flags.
package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func TestChromeFlagsLifecycle(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("GITMAP_CHROME_USER_DATA", tempDir)

	// Seed Local State
	localStatePath := filepath.Join(tempDir, constants.ChromeLocalStateFile)
	initialState := map[string]any{
		"browser": map[string]any{
			"enabled_labs_experiments": []any{"existing-test-flag@1"},
		},
	}
	raw, _ := json.Marshal(initialState)
	_ = os.WriteFile(localStatePath, raw, constants.FilePermission)

	// 1. Test listing
	flags, _, err := readEnabledChromeFlags()
	if err != nil {
		t.Fatalf("read flags failed: %v", err)
	}
	if len(flags) != 1 || flags[0] != "existing-test-flag@1" {
		t.Errorf("unexpected flags: %v", flags)
	}

	// 2. Test enabling a new flag
	if err := runChromeFlags([]string{"enable", "enable-gpu-rasterization"}); err != nil {
		t.Fatalf("enable flag failed: %v", err)
	}
	flagsAfterEnable, _, _ := readEnabledChromeFlags()
	if len(flagsAfterEnable) != 2 {
		t.Fatalf("expected 2 flags, got %d (%v)", len(flagsAfterEnable), flagsAfterEnable)
	}

	// 3. Test disabling a flag
	if err := runChromeFlags([]string{"disable", "existing-test-flag"}); err != nil {
		t.Fatalf("disable flag failed: %v", err)
	}
	flagsAfterDisable, _, _ := readEnabledChromeFlags()
	if len(flagsAfterDisable) != 1 || flagsAfterDisable[0] != "enable-gpu-rasterization@1" {
		t.Errorf("unexpected flags after disable: %v", flagsAfterDisable)
	}

	// 4. Test resetting flags
	if err := runChromeFlags([]string{"reset"}); err != nil {
		t.Fatalf("reset flags failed: %v", err)
	}
	flagsAfterReset, _, _ := readEnabledChromeFlags()
	if len(flagsAfterReset) != 0 {
		t.Errorf("expected 0 flags after reset, got %v", flagsAfterReset)
	}
}
