// Package cmd — chrome_observe_test.go: unit tests for Chrome process & tab observation.
package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func TestChromeObserve(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("GITMAP_CHROME_USER_DATA", tempDir)
	profDir := filepath.Join(tempDir, "Default")
	_ = os.MkdirAll(profDir, constants.DirPermission)

	// Test observe with table output
	if err := runChromeObserve([]string{"--profile=Default"}); err != nil {
		t.Fatalf("observe table failed: %v", err)
	}

	// Test observe with json output
	if err := runChromeObserve([]string{"--json", "--profile=Default"}); err != nil {
		t.Fatalf("observe json failed: %v", err)
	}

	// Test observe with yaml output
	if err := runChromeObserve([]string{"--yaml", "--profile=Default"}); err != nil {
		t.Fatalf("observe yaml failed: %v", err)
	}
}
