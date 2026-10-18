package cmdspace

import (
	"os"
	"path/filepath"
	"strings"
)

// WindowsExecutableExts contains the supported executable file extensions on Windows.
var WindowsExecutableExts = []string{".exe", ".cmd", ".bat"}

// IsWindowsExecutable reports whether path represents an existing executable on Windows.
func IsWindowsExecutable(path string) bool {
	if len(path) == 0 {
		return false
	}
	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		return false
	}
	return HasWindowsExecutableExt(path)
}

// HasWindowsExecutableExt checks if the file path has a valid Windows executable extension.
func HasWindowsExecutableExt(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	if ext == "" {
		return true
	}
	for _, valid := range WindowsExecutableExts {
		if ext == valid {
			return true
		}
	}
	return false
}
