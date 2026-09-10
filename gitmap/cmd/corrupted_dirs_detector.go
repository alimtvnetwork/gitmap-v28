package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func isProtectedPath(path, homeDir string) bool {
	cleanPath := filepath.ToSlash(filepath.Clean(path))
	cleanHome := filepath.ToSlash(filepath.Clean(homeDir))
	if cleanPath == "/" || cleanPath == "." || cleanPath == "/tmp" {
		return true
	}
	if cleanPath == cleanHome || cleanPath == "/usr" || cleanPath == "/usr/local" || cleanPath == "/usr/local/bin" {
		return true
	}
	return false
}

func isCorruptedDirName(name string) (bool, string) {
	if name == CorruptedTildeDir {
		return true, "literal tilde directory"
	}
	if strings.Contains(name, AnsiEscapePrefix) || strings.Contains(name, AnsiEscapeOctal) {
		return true, "contains ANSI escape sequence"
	}
	if strings.Contains(name, "\n") || strings.Contains(name, "\r") {
		return true, "contains embedded newline or carriage return"
	}
	low := strings.ToLower(name)
	if strings.Contains(low, CorruptedPatternQuickInstaller) || strings.Contains(low, CorruptedPatternInstaller) {
		return true, "contains installer banner prompt text"
	}
	if strings.Contains(name, CorruptedPatternDefault) || strings.Contains(name, CorruptedPatternPrompt) {
		return true, "contains installer prompt text"
	}
	return false, ""
}

func getCandidateRoots() []string {
	var roots []string
	if home, err := os.UserHomeDir(); err == nil && home != "" {
		roots = append(roots, home)
		roots = append(roots, filepath.Join(home, ".local"))
	}
	if cwd, err := os.Getwd(); err == nil && cwd != "" {
		roots = append(roots, cwd)
	}
	if tmp := os.TempDir(); tmp != "" {
		roots = append(roots, tmp)
	}
	return roots
}

func scanCandidateRoot(root, homeDir string) []CorruptedDirInfo {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil
	}
	var results []CorruptedDirInfo
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		fullPath := filepath.Join(root, entry.Name())
		if isProtectedPath(fullPath, homeDir) {
			continue
		}
		isCorrupted, reason := isCorruptedDirName(entry.Name())
		if isCorrupted {
			info := buildCorruptedInfo(fullPath, entry.Name(), reason)
			results = append(results, info)
		}
	}
	return results
}

func buildCorruptedInfo(path, name, reason string) CorruptedDirInfo {
	subEntries, _ := os.ReadDir(path)
	var files []string
	for _, sub := range subEntries {
		files = append(files, sub.Name())
	}
	return CorruptedDirInfo{
		Path:           path,
		Name:           name,
		Reason:         reason,
		IsLiteralTilde: name == CorruptedTildeDir,
		HasFiles:       len(files) > 0,
		Files:          files,
	}
}

// DetectCorruptedDirs scans standard roots for corrupted installation directories.
func DetectCorruptedDirs() ([]CorruptedDirInfo, error) {
	homeDir, _ := os.UserHomeDir()
	roots := getCandidateRoots()
	var detected []CorruptedDirInfo
	seen := make(map[string]bool)
	for _, root := range roots {
		for _, info := range scanCandidateRoot(root, homeDir) {
			if !seen[info.Path] {
				seen[info.Path] = true
				detected = append(detected, info)
			}
		}
	}
	return detected, nil
}

func probeCorruptedInstallDirs() DoctorCheck {
	return DoctorCheck{
		Name:    "install-dirs",
		FixHint: "Run `gitmap clean-corrupted --force` to remove corrupted folders",
		Run: func() (bool, string) {
			detected, err := DetectCorruptedDirs()
			if err != nil {
				return false, "error detecting corrupted dirs: " + err.Error()
			}
			if len(detected) > 0 {
				return false, fmt.Sprintf("%d corrupted folder(s) found", len(detected))
			}
			return true, "clean"
		},
	}
}
