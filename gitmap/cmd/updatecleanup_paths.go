package cmd

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

type updateCleanupConfig struct {
	DeployPath  string
	BuildOutput string
}

type updateCleanupContext struct {
	selfPath       string
	tempPatterns   []string
	backupPatterns []string
}

// loadUpdateCleanupContext resolves every directory that update-cleanup must scan.
func loadUpdateCleanupContext() updateCleanupContext {
	selfPath := resolveCleanupSelfPath()
	config := readUpdateCleanupConfig(constants.RepoPath)
	tempDirs := collectTempCleanupDirs(selfPath, constants.RepoPath, config)
	backupDirs := collectBackupCleanupDirs(selfPath, constants.RepoPath, config)

	return updateCleanupContext{
		selfPath:       selfPath,
		tempPatterns:   buildCleanupPatterns(tempDirs, constants.UpdateCopyGlob),
		backupPatterns: buildCleanupPatterns(backupDirs, constants.OldBackupGlob),
	}
}

// resolveCleanupSelfPath returns the active binary path for cleanup scanning.
func resolveCleanupSelfPath() string {
	selfPath, err := os.Executable()
	if err != nil {
		logUpdateCleanupExecutableError(err)

		return ""
	}

	return filepath.Clean(selfPath)
}

// readUpdateCleanupConfig reads the cleanup-relevant values from powershell.json.
func readUpdateCleanupConfig(repoPath string) updateCleanupConfig {
	config := updateCleanupConfig{BuildOutput: constants.DefaultBuildOutput}
	if len(repoPath) == 0 {
		return config
	}

	return readUpdateCleanupConfigFile(repoPath, config)
}

// readUpdateCleanupConfigFile loads cleanup settings from disk.
func readUpdateCleanupConfigFile(repoPath string, config updateCleanupConfig) updateCleanupConfig {
	configPath := filepath.Join(repoPath, constants.GitMapSubdir, constants.PowershellConfigFile)
	data, err := os.ReadFile(configPath)
	if err != nil {
		logUpdateCleanupConfigReadError(configPath, err)

		return config
	}

	config.DeployPath = extractJSONString(data, constants.JSONKeyDeployPath)
	config.BuildOutput = firstNonEmpty(extractJSONString(data, constants.JSONKeyBuildOutput), config.BuildOutput)

	return config
}

// collectTempCleanupDirs returns the directories scanned for gitmap-update-* leftovers.
func collectTempCleanupDirs(selfPath, repoPath string, config updateCleanupConfig) []string {
	dirs := []string{os.TempDir()}
	dirs = appendCleanupDir(dirs, resolveCleanupDir(selfPath))

	return appendResolvedCleanupDirs(dirs, selfPath, repoPath, config)
}

// collectBackupCleanupDirs returns the directories scanned for .old backups.
func collectBackupCleanupDirs(selfPath, repoPath string, config updateCleanupConfig) []string {
	dirs := appendCleanupDir(nil, resolveCleanupDir(selfPath))

	return appendResolvedCleanupDirs(dirs, selfPath, repoPath, config)
}

// appendResolvedCleanupDirs appends derived deploy and build directories.
func appendResolvedCleanupDirs(dirs []string, selfPath, repoPath string, config updateCleanupConfig) []string {
	dirs = appendCleanupDir(dirs, deriveDeployAppDir(selfPath))
	dirs = appendCleanupDir(dirs, resolveConfigDeployAppDir(config.DeployPath))
	dirs = appendCleanupDir(dirs, resolveBuildOutputDir(repoPath, config.BuildOutput))

	return dirs
}

// resolveCleanupDir returns the parent directory for a file path.
func resolveCleanupDir(filePath string) string {
	if len(filePath) == 0 {
		return ""
	}

	return filepath.Dir(filePath)
}

// deriveDeployAppDir mirrors run.ps1 PATH-derived deploy target resolution.
// Recognizes both the new "gitmap-cli" subdir and the legacy "gitmap" subdir
// so cleanup keeps working during the v3.6.0 migration window.
func deriveDeployAppDir(selfPath string) string {
	selfDir := resolveCleanupDir(selfPath)
	if len(selfDir) == 0 {
		return ""
	}
	base := filepath.Base(selfDir)
	if base == constants.GitMapCliSubdir || base == constants.GitMapSubdir {
		return selfDir
	}

	parentDir := filepath.Dir(selfDir)
	if len(parentDir) == 0 || parentDir == selfDir {
		return ""
	}

	return filepath.Join(parentDir, constants.GitMapCliSubdir)
}

// resolveConfigDeployAppDir returns the nested gitmap-cli deploy directory from config.
func resolveConfigDeployAppDir(deployPath string) string {
	if len(deployPath) == 0 {
		return ""
	}

	return filepath.Join(deployPath, constants.GitMapCliSubdir)
}

// resolveBuildOutputDir resolves the build output directory from repo root + config.
func resolveBuildOutputDir(repoPath, buildOutput string) string {
	if len(repoPath) == 0 {
		return ""
	}
	if len(buildOutput) == 0 {
		buildOutput = constants.DefaultBuildOutput
	}
	// Treat both OS-native absolute paths AND POSIX-style leading "/"
	// as absolute. On Windows filepath.IsAbs("/custom/bin") returns
	// false (Windows wants a drive letter), but tests + config files
	// authored on Linux/macOS routinely use "/custom/bin" to mean
	// "absolute, root of current drive". Honor that intent so a
	// shared config doesn't silently get re-rooted under repoPath.
	if isAbsolutePath(buildOutput) {
		return filepath.Clean(buildOutput)
	}

	return filepath.Clean(filepath.Join(repoPath, buildOutput))
}

// buildCleanupPatterns expands scan directories into filepath.Glob patterns.
func buildCleanupPatterns(dirs []string, glob string) []string {
	patterns := make([]string, 0, len(dirs))
	for _, dir := range dirs {
		if len(dir) == 0 {
			continue
		}
		patterns = append(patterns, filepath.Join(dir, glob))
	}

	return patterns
}

// appendCleanupDir appends a cleaned directory once.
func appendCleanupDir(dirs []string, dir string) []string {
	if len(dir) == 0 {
		return dirs
	}

	cleanDir := filepath.Clean(dir)
	if hasCleanupDir(dirs, cleanDir) {
		return dirs
	}

	return append(dirs, cleanDir)
}

// hasCleanupDir reports whether the directory is already tracked.
func hasCleanupDir(dirs []string, dir string) bool {
	normalizedDir := normalizeCleanupPath(dir)
	for _, existingDir := range dirs {
		if normalizeCleanupPath(existingDir) == normalizedDir {
			return true
		}
	}

	return false
}

// firstNonEmpty returns value when set, otherwise fallback.
func firstNonEmpty(value, fallback string) string {
	if len(value) > 0 {
		return value
	}

	return fallback
}

// normalizeCleanupPath normalizes cleanup paths for cross-platform comparisons.
func normalizeCleanupPath(path string) string {
	cleanPath := filepath.Clean(path)
	if runtime.GOOS != "windows" {
		return cleanPath
	}

	return strings.ToLower(cleanPath)
}

// isAbsolutePath treats both OS-native absolute paths and POSIX-style
// "/leading-slash" or "\leading-slash" paths as absolute. On Windows
// filepath.IsAbs("/custom/bin") returns false (Windows wants a drive
// letter), but config files authored on Linux/macOS routinely use
// "/custom/bin" to mean "absolute". Honor that intent so a shared
// config doesn't silently get re-rooted under repoPath.
func isAbsolutePath(path string) bool {
	if filepath.IsAbs(path) {
		return true
	}
	if len(path) == 0 {
		return false
	}
	first := path[0]
	return first == '/' || first == '\\'
}
