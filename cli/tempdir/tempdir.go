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
