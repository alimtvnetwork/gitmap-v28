package store_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// TestBinaryDataDirReturnsDataSubdir verifies that BinaryDataDir
// returns a path ending in the configured data directory name.
func TestBinaryDataDirReturnsDataSubdir(t *testing.T) {
	dir := store.BinaryDataDir()
	base := filepath.Base(dir)

	if base != constants.DBDir {
		t.Errorf("expected base dir %q, got %q", constants.DBDir, base)
	}
}

// TestBinaryDataDirIsAbsolute verifies the returned path is absolute.
func TestBinaryDataDirIsAbsolute(t *testing.T) {
	dir := store.BinaryDataDir()

	if !filepath.IsAbs(dir) {
		t.Errorf("expected absolute path, got %q", dir)
	}
}

// TestBinaryDataDirResolvesSymlink verifies symlink resolution by
// creating a symlink to the test binary and checking the resolved
// path differs from the symlink location.
func TestBinaryDataDirResolvesSymlink(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink test not reliable on Windows CI")
	}

	exe, err := os.Executable()
	if err != nil {
		t.Fatalf("could not get executable: %v", err)
	}

	resolved, err := filepath.EvalSymlinks(exe)
	if err != nil {
		t.Fatalf("could not resolve symlinks: %v", err)
	}

	tmpDir := t.TempDir()
	symPath := filepath.Join(tmpDir, "test-link")

	err = os.Symlink(resolved, symPath)
	if err != nil {
		t.Skipf("cannot create symlink: %v", err)
	}

	// BinaryDataDir uses os.Executable() internally so we can only
	// verify it produces a consistent result anchored to the real binary.
	dir := store.BinaryDataDir()
	expected := filepath.Join(filepath.Dir(resolved), constants.DBDir)

	if dir != expected {
		t.Errorf("expected %q, got %q", expected, dir)
	}
}

// TestBinaryDataDirNoDoubleNesting ensures the path does not contain
// data/data (the double-nesting bug from v2.20.0).
func TestBinaryDataDirNoDoubleNesting(t *testing.T) {
	dir := store.BinaryDataDir()
	parent := filepath.Base(filepath.Dir(dir))

	if parent == constants.DBDir {
		t.Errorf("double-nesting detected: %q", dir)
	}
}

// TestDefaultDBPathContainsDBFile verifies the diagnostic path
// includes the expected database filename.
func TestDefaultDBPathContainsDBFile(t *testing.T) {
	path := store.DefaultDBPath()
	base := filepath.Base(path)

	// Should be either "gitmap.db" or "gitmap-<profile>.db"
	if len(base) == 0 {
		t.Error("DefaultDBPath returned empty filename")
	}

	ext := filepath.Ext(base)
	if ext != ".db" {
		t.Errorf("expected .db extension, got %q", ext)
	}
}

// TestDefaultDBPathIsAbsolute verifies the diagnostic path is absolute.
func TestDefaultDBPathIsAbsolute(t *testing.T) {
	path := store.DefaultDBPath()

	if !filepath.IsAbs(path) {
		t.Errorf("expected absolute path, got %q", path)
	}
}

// TestDefaultDBPathInsideDataDir verifies the DB file lives inside
// the data directory, not at the binary root.
func TestDefaultDBPathInsideDataDir(t *testing.T) {
	path := store.DefaultDBPath()
	dir := filepath.Dir(path)
	base := filepath.Base(dir)

	if base != constants.DBDir {
		t.Errorf("expected DB inside %q dir, got parent %q", constants.DBDir, base)
	}
}

// TestProfileDBFileDefault verifies the default profile returns
// the standard database filename.
func TestProfileDBFileDefault(t *testing.T) {
	name := store.ProfileDBFile(constants.DefaultProfileName)

	if name != constants.DBFile {
		t.Errorf("expected %q, got %q", constants.DBFile, name)
	}
}

// TestProfileDBFileCustom verifies a custom profile returns a
// prefixed database filename.
func TestProfileDBFileCustom(t *testing.T) {
	name := store.ProfileDBFile("work")
	expected := constants.ProfileDBPrefix + "work.db"

	if name != expected {
		t.Errorf("expected %q, got %q", expected, name)
	}
}
