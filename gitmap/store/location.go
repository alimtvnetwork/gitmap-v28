package store

import (
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// BinaryDataDir returns the data directory relative to the running
// executable's physical location. This ensures the SQLite database
// is always co-located with the binary, regardless of the working
// directory from which gitmap is invoked.
func BinaryDataDir() string {
	exe, err := os.Executable()
	if err != nil {
		return filepath.Join(".", constants.DBDir)
	}

	resolved, err := filepath.EvalSymlinks(exe)
	if err != nil {
		resolved = exe
	}

	return filepath.Join(filepath.Dir(resolved), constants.DBDir)
}

// OpenDefault opens the database from the binary's data directory.
func OpenDefault() (*DB, error) {
	dir := BinaryDataDir()
	baseDir := filepath.Dir(dir) // binary dir without /data
	dbFile := ActiveProfileDBFile(baseDir)
	dbPath := filepath.Join(dir, dbFile)

	return openDBAt(dbPath)
}

// OpenDefaultProfile opens a named profile's database from the
// binary's data directory.
func OpenDefaultProfile(profileName string) (*DB, error) {
	dir := BinaryDataDir()
	dbFile := ProfileDBFile(profileName)
	dbPath := filepath.Join(dir, dbFile)

	return openDBAt(dbPath)
}

// DefaultDBPath returns the resolved database path for diagnostics.
func DefaultDBPath() string {
	dir := BinaryDataDir()
	baseDir := filepath.Dir(dir)
	dbFile := ActiveProfileDBFile(baseDir)

	return filepath.Join(dir, dbFile)
}
