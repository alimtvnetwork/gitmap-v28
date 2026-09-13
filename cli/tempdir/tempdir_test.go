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

func TestClearRepoTestTempDir_WithExistingArtifacts_PurgesArtifacts(t *testing.T) {
	tDir := TestTempDir()
	dummyFile := filepath.Join(tDir, "test-output.tmp")
	_ = os.WriteFile(dummyFile, []byte("stale test output"), 0o644)

	err := ClearRepoTestTempDir()
	if err != nil {
		t.Fatalf("unexpected error clearing test temp dir: %v", err)
	}

	_, statErr := os.Stat(dummyFile)
	if !os.IsNotExist(statErr) {
		t.Errorf("expected dummyFile to be removed after clear, got err: %v", statErr)
	}
}

func TestSandboxTempDir_Always_ReturnsSandboxSubfolder(t *testing.T) {
	path := SandboxTempDir()
	expectedSub := filepath.Join("gitmap", "sandbox")
	if !strings.Contains(path, expectedSub) {
		t.Errorf("expected SandboxTempDir %q to contain %q", path, expectedSub)
	}
}

func TestClearRepoSandboxTempDir_WithExistingArtifacts_PurgesArtifacts(t *testing.T) {
	sDir := SandboxTempDir()
	dummyFile := filepath.Join(sDir, "test-sandbox.tmp")
	_ = os.WriteFile(dummyFile, []byte("stale sandbox file"), 0o644)

	err := ClearRepoSandboxTempDir()
	if err != nil {
		t.Fatalf("unexpected error clearing sandbox temp dir: %v", err)
	}

	_, statErr := os.Stat(dummyFile)
	if !os.IsNotExist(statErr) {
		t.Errorf("expected dummyFile to be removed after clear, got err: %v", statErr)
	}
}

func TestClearAllRepoTempDirs_WithArtifactsInMultipleDirs_PurgesAll(t *testing.T) {
	bFile := filepath.Join(BuildTempDir(), "dummy-build.tmp")
	tFile := filepath.Join(TestTempDir(), "dummy-test.tmp")
	sFile := filepath.Join(SandboxTempDir(), "dummy-sandbox.tmp")

	_ = os.WriteFile(bFile, []byte("b"), 0o644)
	_ = os.WriteFile(tFile, []byte("t"), 0o644)
	_ = os.WriteFile(sFile, []byte("s"), 0o644)

	err := ClearAllRepoTempDirs()
	if err != nil {
		t.Fatalf("unexpected error running ClearAllRepoTempDirs: %v", err)
	}

	for _, f := range []string{bFile, tFile, sFile} {
		if _, statErr := os.Stat(f); !os.IsNotExist(statErr) {
			t.Errorf("expected file %q to be removed, got err: %v", f, statErr)
		}
	}
}
