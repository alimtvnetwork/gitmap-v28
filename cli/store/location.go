package store

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var binaryDataDirOverride string

// SetBinaryDataDirForTesting overrides the binary data directory during tests.
func SetBinaryDataDirForTesting(dir string) {
	binaryDataDirOverride = dir
}

// GlobalUserDataDir returns the canonical user-level data directory.
func GlobalUserDataDir() string {
	if localAppData := os.Getenv("LOCALAPPDATA"); runtime.GOOS == "windows" && localAppData != "" {
		return filepath.Join(localAppData, "gitmap-cli", constants.DBDir)
	}
	home, err := os.UserHomeDir()
	if err == nil && home != "" {
		return filepath.Join(home, ".gitmap", constants.DBDir)
	}
	return ""
}

func findGlobalFallbackDB(dbFile, currentDBPath string) string {
	if binaryDataDirOverride != "" {
		return ""
	}
	if _, err := os.Stat(currentDBPath); !os.IsNotExist(err) {
		return ""
	}
	globalDir := GlobalUserDataDir()
	if globalDir == "" {
		return ""
	}
	globalDBPath := filepath.Join(globalDir, dbFile)
	if _, err := os.Stat(globalDBPath); err == nil {
		return globalDBPath
	}
	return ""
}

// BinaryDataDir returns the data directory relative to the running
// executable's physical location. This ensures the SQLite database
// is always co-located with the binary, regardless of the working
// directory from which gitmap is invoked.
func BinaryDataDir() string {
	if binaryDataDirOverride != "" {
		return binaryDataDirOverride
	}

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

	if fallback := findGlobalFallbackDB(dbFile, dbPath); fallback != "" {
		return openDBAt(fallback)
	}

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

// OpenGlobalDefault opens the database from the canonical global user data directory.
func OpenGlobalDefault() (*DB, error) {
	globalDir := GlobalUserDataDir()
	if globalDir == "" {
		return nil, os.ErrNotExist
	}
	dbPath := filepath.Join(globalDir, constants.DBFile)

	return openDBAt(dbPath)
}

// DefaultDBPath returns the resolved database path for diagnostics.
func DefaultDBPath() string {
	dir := BinaryDataDir()
	baseDir := filepath.Dir(dir)
	dbFile := ActiveProfileDBFile(baseDir)

	return filepath.Join(dir, dbFile)
}
