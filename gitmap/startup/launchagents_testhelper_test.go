package startup

// Shared test helper: redirect $HOME so AutostartDir / darwinLaunchAgentsDir
// resolve under t.TempDir() and pre-create the Library/LaunchAgents folder.
//
// Previously duplicated in add_darwin_test.go and plist_test.go (identical
// signatures, no build tags) which is a Go redeclaration footgun in the
// flat startup package. Centralized here so future macOS startup tests
// share one definition.

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// withFakeLaunchAgentsDir sets $HOME to a fresh t.TempDir() and returns
// the pre-created <temp>/Library/LaunchAgents path. Skips on non-darwin
// so Linux/Windows CI runs see the suite as skipped, not failed.
func withFakeLaunchAgentsDir(t *testing.T) string {
	t.Helper()
	if runtime.GOOS != "darwin" {
		t.Skip("plist tests are macOS-only; add_test.go covers Linux")
	}
	root := t.TempDir()
	t.Setenv("HOME", root)
	dir := filepath.Join(root, "Library", "LaunchAgents")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	return dir
}
