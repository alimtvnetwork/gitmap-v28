package vscodepm

import (
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/fsutil"
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

	clean := fsutil.CanonicalPathKey(path)
	if isTempDirectory(clean) || isTestArtifactPath(clean) {
		return true
	}

	return false
}

func isTempDirectory(clean string) bool {
	tempDir := fsutil.CanonicalPathKey(os.TempDir())
	if tempDir != "" && strings.HasPrefix(clean, tempDir) {
		return true
	}

	tempEnv := fsutil.CanonicalPathKey(os.Getenv("TEMP"))
	if tempEnv != "" && strings.HasPrefix(clean, tempEnv) {
		return true
	}

	tmpEnv := fsutil.CanonicalPathKey(os.Getenv("TMP"))
	if tmpEnv != "" && strings.HasPrefix(clean, tmpEnv) {
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
