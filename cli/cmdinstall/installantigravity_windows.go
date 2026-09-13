//go:build windows

package cmdinstall

import (
	"os"
	"os/exec"
	"path/filepath"
)

func resolveLocalAppDataDir() string {
	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData != "" {
		return localAppData
	}

	home, _ := os.UserHomeDir()

	return filepath.Join(home, "AppData", "Local")
}

func getAntigravityDesktopWindowsExePath() string {
	localAppData := resolveLocalAppDataDir()

	return filepath.Join(localAppData, "Programs", "Antigravity", "Antigravity.exe")
}

func getWindowsAntigravityCandidatePaths(localAppData string) []string {
	return []string{
		filepath.Join(localAppData, "Programs", "Antigravity", "Antigravity.exe"),
		filepath.Join(localAppData, "Programs", "antigravity", "Antigravity.exe"),
		filepath.Join(localAppData, "Programs", "antigravity", "antigravity.exe"),
		filepath.Join(os.Getenv("ProgramFiles"), "Antigravity", "Antigravity.exe"),
		filepath.Join(os.Getenv("ProgramFiles(x86)"), "Antigravity", "Antigravity.exe"),
	}
}

func checkCandidatePaths(paths []string) (string, bool) {
	for _, candidate := range paths {
		if candidate == "" {
			continue
		}

		if _, err := os.Stat(candidate); err == nil {
			return candidate, true
		}
	}

	return "", false
}

func findAntigravityInPath() (string, bool) {
	for _, name := range []string{"antigravity", "Antigravity"} {
		if p, err := exec.LookPath(name); err == nil {
			return p, true
		}
	}

	return "", false
}

func findInstalledAntigravityDesktopPath() (string, bool) {
	localAppData := resolveLocalAppDataDir()
	candidates := getWindowsAntigravityCandidatePaths(localAppData)
	if path, isFound := checkCandidatePaths(candidates); isFound {
		return path, true
	}

	return findAntigravityInPath()
}

func updateDesktopDatabase(desktopDir string) {
	// No-op on Windows
}

func resolveAppIconPath(installDir string) string {
	return ""
}
