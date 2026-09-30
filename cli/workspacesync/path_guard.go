package workspacesync

import (
	"os"
	"path/filepath"
	"strings"
)

// IsTempOrTestPath reports whether a path is located in an OS temp dir or test directory.
func IsTempOrTestPath(path string) bool {
	if strings.TrimSpace(path) == "" {
		return false
	}
	clean := strings.ToLower(filepath.Clean(path))
	if isTempEnvironment() || matchesTempPrefix(clean) {
		return true
	}
	return matchesTestKeywords(clean)
}

func isTempEnvironment() bool {
	return os.Getenv("GITMAP_TESTING") == "1" || os.Getenv("GO_TEST") == "1"
}

func matchesTempPrefix(clean string) bool {
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

func matchesTestKeywords(clean string) bool {
	if strings.Contains(clean, "appdata\\local\\temp") || strings.Contains(clean, "appdata/local/temp") {
		return true
	}
	return strings.Contains(clean, "testprovisiondesttarget") || strings.Contains(clean, "testensureorprovision")
}

// IsRestrictedPath reports whether a path should not be registered as a workspace.
func IsRestrictedPath(path string) bool {
	if strings.TrimSpace(path) == "" {
		return true
	}
	clean := filepath.Clean(path)
	if abs, err := filepath.Abs(clean); err == nil {
		clean = abs
	}
	if isRootOrDrivePath(clean) {
		return true
	}
	if isUserHomeOrParent(clean) || IsTempOrTestPath(clean) {
		return true
	}
	return isSystemPath(clean)
}

func isRootOrDrivePath(clean string) bool {
	if clean == "/" || clean == "\\" || filepath.Dir(clean) == clean {
		return true
	}
	vol := filepath.VolumeName(clean)
	return vol != "" && (clean == vol || clean == vol+"\\" || clean == vol+"/")
}

func isUserHomeOrParent(clean string) bool {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return false
	}
	cleanHome := filepath.Clean(home)
	if strings.EqualFold(clean, cleanHome) {
		return true
	}
	parent := filepath.Dir(cleanHome)
	return parent != cleanHome && strings.EqualFold(clean, parent)
}

func isSystemPath(clean string) bool {
	candidates := []string{
		os.Getenv("WINDIR"), os.Getenv("SystemRoot"), os.Getenv("ProgramFiles"),
		os.Getenv("ProgramFiles(x86)"), os.Getenv("ProgramData"),
		"/etc", "/usr", "/bin", "/sbin", "/var", "/root", "/System", "/Library",
	}
	for _, sys := range candidates {
		if sys != "" && strings.EqualFold(clean, filepath.Clean(sys)) {
			return true
		}
	}
	return false
}

func isGitRepoPath(path string) bool {
	if strings.TrimSpace(path) == "" {
		return false
	}
	dotGit := filepath.Join(path, ".git")
	info, err := os.Stat(dotGit)
	return err == nil && (info.IsDir() || !info.IsDir())
}
