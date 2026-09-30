package vscodepm

import (
	"os"
	"path/filepath"
	"strings"
)

// BypassDisallowedPathFilterForTesting allows unit tests using t.TempDir() to test optimize logic.
var BypassDisallowedPathFilterForTesting = false

// IsDisallowedProjectPath reports whether path is empty, temp, or test directory.
func IsDisallowedProjectPath(path string) bool {
	if BypassDisallowedPathFilterForTesting {
		return false
	}
	if strings.TrimSpace(path) == "" {
		return true
	}
	clean := strings.ToLower(filepath.Clean(path))
	if isTempDirectory(clean) || isTestArtifactPath(clean) {
		return true
	}
	return false
}

func isTempDirectory(clean string) bool {
	tempDir := strings.ToLower(filepath.Clean(os.TempDir()))
	if tempDir != "" && strings.HasPrefix(clean, tempDir) {
		return true
	}
	if tempEnv := strings.ToLower(filepath.Clean(os.Getenv("TEMP"))); tempEnv != "" && strings.HasPrefix(clean, tempEnv) {
		return true
	}
	if tmpEnv := strings.ToLower(filepath.Clean(os.Getenv("TMP"))); tmpEnv != "" && strings.HasPrefix(clean, tmpEnv) {
		return true
	}
	return strings.HasPrefix(clean, "/tmp") || strings.HasPrefix(clean, "/var/tmp")
}

func isTestArtifactPath(clean string) bool {
	if strings.Contains(clean, "appdata\\local\\temp") || strings.Contains(clean, "appdata/local/temp") {
		return true
	}
	return strings.Contains(clean, "testprovisiondesttarget") || strings.Contains(clean, "testensureorprovision")
}
