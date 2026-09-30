package workspacesync

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsTempOrTestPath(t *testing.T) {
	tempDir := t.TempDir()
	if !IsTempOrTestPath(tempDir) {
		t.Errorf("expected tempDir %q to be identified as temp or test path", tempDir)
	}

	testTarget := filepath.Join(tempDir, "auto-dest-repo")
	if !IsTempOrTestPath(testTarget) {
		t.Errorf("expected %q to be identified as temp or test path", testTarget)
	}

	validWorkPath := filepath.Join("D:", "work", "my-app")
	if os.Getenv("GITMAP_TESTING") != "1" && IsTempOrTestPath(validWorkPath) {
		t.Errorf("expected %q not to be temp path in non-testing env", validWorkPath)
	}
}

func TestSyncAll_RestrictedTempPath(t *testing.T) {
	tempDir := t.TempDir()
	// SyncAll must safely no-op for temp paths and not panic.
	SyncAll(tempDir, "temp-repo")
	SyncWithoutDesktop(tempDir, "temp-repo")
	SyncAgyOnly(tempDir, "temp-repo")
}
