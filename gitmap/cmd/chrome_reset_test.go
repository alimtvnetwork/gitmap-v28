// Package cmd — chrome_reset_test.go: unit tests for Chrome cache, cookies,
// and profile reset operations.
package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func TestChromeResetOperations(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("GITMAP_CHROME_USER_DATA", tempDir)
	profDir := filepath.Join(tempDir, "Default")
	_ = os.MkdirAll(profDir, constants.DirPermission)

	// Create cache folder and history file
	cacheDir := filepath.Join(profDir, "Cache")
	_ = os.MkdirAll(cacheDir, constants.DirPermission)
	historyFile := filepath.Join(profDir, "History")
	_ = os.WriteFile(historyFile, []byte("fake history"), constants.FilePermission)

	// 1. Reset cache
	if err := runChromeReset([]string{"--cache", "--profile=Default"}); err != nil {
		t.Fatalf("reset cache failed: %v", err)
	}
	if _, err := os.Stat(cacheDir); err == nil {
		t.Errorf("expected Cache directory to be removed")
	}
	if _, err := os.Stat(historyFile); err != nil {
		t.Errorf("expected History file to still exist")
	}

	// 2. Reset history
	if err := runChromeReset([]string{"--history", "--profile=Default"}); err != nil {
		t.Fatalf("reset history failed: %v", err)
	}
	if _, err := os.Stat(historyFile); err == nil {
		t.Errorf("expected History file to be removed")
	}

	// 3. Full reset
	prefFile := filepath.Join(profDir, "Preferences")
	_ = os.WriteFile(prefFile, []byte(`{"test": 123}`), constants.FilePermission)
	if err := runChromeReset([]string{"--all", "--profile=Default"}); err != nil {
		t.Fatalf("full reset failed: %v", err)
	}
	if _, err := os.Stat(prefFile + ".bak"); err != nil {
		t.Errorf("expected Preferences backup to be created")
	}
}
