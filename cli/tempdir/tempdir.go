package tempdir

import (
	"os"
	"path/filepath"
)

// RepoTempSubdir defines the canonical repository-scoped subdirectory name under OS temp.
const RepoTempSubdir = "gitmap"

// RepoTempDir returns a repository-scoped path under os.TempDir(), e.g. <os.TempDir()>/gitmap/<subdirs...>.
// It automatically creates the directory if it does not exist.
func RepoTempDir(subdirs ...string) string {
	parts := append([]string{os.TempDir(), RepoTempSubdir}, subdirs...)
	target := filepath.Join(parts...)
	_ = os.MkdirAll(target, 0o755)

	return target
}

// BuildTempDir returns the dedicated repository build temporary directory.
func BuildTempDir() string {
	return RepoTempDir("build")
}

// TestTempDir returns the dedicated repository test temporary directory.
func TestTempDir() string {
	return RepoTempDir("test")
}

// ClearRepoBuildTempDir cleans all artifacts in the build directory before a build runs,
// preventing storage bloat and respecting disk storage reuse.
func ClearRepoBuildTempDir() error {
	bDir := BuildTempDir()
	entries, err := os.ReadDir(bDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		_ = os.RemoveAll(filepath.Join(bDir, entry.Name()))
	}

	return nil
}

// ClearRepoTestTempDir cleans all artifacts in the test directory before tests run.
func ClearRepoTestTempDir() error {
	tDir := TestTempDir()
	entries, err := os.ReadDir(tDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		_ = os.RemoveAll(filepath.Join(tDir, entry.Name()))
	}

	return nil
}

// SandboxTempDir returns the dedicated repository sandbox temporary directory.
func SandboxTempDir() string {
	return RepoTempDir("sandbox")
}

// ClearRepoSandboxTempDir cleans all artifacts in the sandbox directory before/after tests run.
func ClearRepoSandboxTempDir() error {
	sDir := SandboxTempDir()
	entries, err := os.ReadDir(sDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		_ = os.RemoveAll(filepath.Join(sDir, entry.Name()))
	}

	return nil
}

// ClearAllRepoTempDirs sweeps all standard repository temporary directories to free disk space.
func ClearAllRepoTempDirs() error {
	_ = ClearRepoBuildTempDir()
	_ = ClearRepoTestTempDir()
	_ = ClearRepoSandboxTempDir()

	for _, sub := range []string{"purge", "downloads", "handoff"} {
		dir := RepoTempDir(sub)
		entries, err := os.ReadDir(dir)
		if err == nil {
			for _, entry := range entries {
				_ = os.RemoveAll(filepath.Join(dir, entry.Name()))
			}
		}
	}

	return nil
}
