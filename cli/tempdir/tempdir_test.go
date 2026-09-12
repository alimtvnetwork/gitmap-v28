package tempdir

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRepoTempDir_WithSubdirs_ReturnsScopedPath(t *testing.T) {
	path := RepoTempDir("test", "sub")
	expectedSub := filepath.Join("gitmap", "test", "sub")
	if !strings.Contains(path, expectedSub) {
		t.Errorf("expected path %q to contain %q", path, expectedSub)
	}

	info, err := os.Stat(path)
	if err != nil || !info.IsDir() {
		t.Errorf("expected %q to exist as directory, err: %v", path, err)
	}
}

func TestBuildTempDir_Always_ReturnsBuildSubfolder(t *testing.T) {
	path := BuildTempDir()
	expectedSub := filepath.Join("gitmap", "build")
	if !strings.Contains(path, expectedSub) {
		t.Errorf("expected BuildTempDir %q to contain %q", path, expectedSub)
	}
}

func TestTestTempDir_Always_ReturnsTestSubfolder(t *testing.T) {
	path := TestTempDir()
	expectedSub := filepath.Join("gitmap", "test")
	if !strings.Contains(path, expectedSub) {
		t.Errorf("expected TestTempDir %q to contain %q", path, expectedSub)
	}
}

func TestClearRepoBuildTempDir_WithExistingArtifacts_PurgesArtifacts(t *testing.T) {
	bDir := BuildTempDir()
	dummyFile := filepath.Join(bDir, "test-artifact.tmp")
	_ = os.WriteFile(dummyFile, []byte("stale binary"), 0o644)

	err := ClearRepoBuildTempDir()
	if err != nil {
		t.Fatalf("unexpected error clearing build temp dir: %v", err)
	}

	_, statErr := os.Stat(dummyFile)
	if !os.IsNotExist(statErr) {
		t.Errorf("expected dummyFile to be removed after clear, got err: %v", statErr)
	}
}
