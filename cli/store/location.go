package store

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var binaryDataDirOverride string

// SetBinaryDataDirForTesting overrides the binary data directory during tests.
func SetBinaryDataDirForTesting(dir string) {
	binaryDataDirOverride = dir
}

// GlobalUserDataDir returns the canonical user-level data directory.
func GlobalUserDataDir() string {
	if runtime.GOOS == "windows" {
		return resolveWindowsUserDataDir()
	}
	home, err := os.UserHomeDir()
	if err == nil && home != "" {
		return filepath.Join(home, ".gitmap", constants.DBDir)
	}
	return ""
}

func resolveWindowsUserDataDir() string {
	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData == "" {
		return ""
	}
	if hostDir := resolveAntigravityHostUserDataDir(localAppData); hostDir != "" {
		return hostDir
	}
	return filepath.Join(localAppData, "gitmap-cli", constants.DBDir)
}

func resolveAntigravityHostUserDataDir(localAppData string) string {
	if !strings.Contains(localAppData, ".antigravity_tools") {
		return ""
	}
	drive := os.Getenv("SystemDrive")
	if drive == "" {
		drive = "C:"
	}
	user := os.Getenv("USERNAME")
	if user == "" {
		return ""
	}
	realUserDir := filepath.Join(drive, "\\Users", user, "AppData", "Local", "gitmap-cli", constants.DBDir)
	targetDB := filepath.Join(realUserDir, constants.DBFile)
	if _, err := os.Stat(targetDB); err != nil {
		return ""
	}
	return realUserDir
}

func isDirWritable(dir string) bool {
	clean := filepath.Clean(dir)
	info, err := os.Stat(clean)
	if os.IsNotExist(err) {
		return isParentDirWritable(clean)
	}
	if err != nil || !info.IsDir() {
		return false
	}

	return isTestFileWritable(clean)
}

func isParentDirWritable(clean string) bool {
	parent := filepath.Dir(clean)
	if parent == clean {
		return false
	}
	return isDirWritable(parent)
}

func isTestFileWritable(dir string) bool {
	testFile := filepath.Join(dir, fmt.Sprintf(".gitmap_write_test_%d_%d", os.Getpid(), time.Now().UnixNano()))
	f, err := os.OpenFile(testFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return false
	}
	_ = f.Close()
	_ = os.Remove(testFile)

	return true
}

func findGlobalFallbackDB(dbFile, currentDBPath string) string {
	if binaryDataDirOverride != "" {
		return ""
	}
	currentDir := filepath.Dir(currentDBPath)
	if !isDirWritable(currentDir) {
		return resolveGlobalDirFallback(dbFile)
	}
	if _, err := os.Stat(currentDBPath); !os.IsNotExist(err) {
		return ""
	}
	return checkGlobalDBExists(dbFile)
}

func resolveGlobalDirFallback(dbFile string) string {
	globalDir := GlobalUserDataDir()
	if globalDir == "" {
		return ""
	}
	return filepath.Join(globalDir, dbFile)
}

func checkGlobalDBExists(dbFile string) string {
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

// EffectiveDataDir returns BinaryDataDir if it is writable. If BinaryDataDir
// is read-only or unwritable (e.g. system install in /usr/local/bin/gitmap-cli/data),
// it automatically falls back to GlobalUserDataDir() (~/.gitmap/data).
func EffectiveDataDir() string {
	dir := BinaryDataDir()
	if isDirWritable(dir) {
		return dir
	}
	if globalDir := GlobalUserDataDir(); globalDir != "" {
		_ = os.MkdirAll(globalDir, 0755)
		return globalDir
	}
	return dir
}

// OpenDefault opens the database from the binary's data directory.
func OpenDefault() (*DB, error) {
	dir := EffectiveDataDir()
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
	dir := EffectiveDataDir()
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
	dir := EffectiveDataDir()
	baseDir := filepath.Dir(dir)
	dbFile := ActiveProfileDBFile(baseDir)

	return filepath.Join(dir, dbFile)
}
