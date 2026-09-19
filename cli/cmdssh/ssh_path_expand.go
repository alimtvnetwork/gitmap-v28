package cmdssh

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var winEnvRegex = regexp.MustCompile(`%([A-Za-z0-9_ -]+)%`)

// ExpandUniversalPath resolves ~, Windows macros (%win%, %win-drive%, %temp%, %appdata%),
// and OS environment variables.
func ExpandUniversalPath(rawPath, osType string) string {
	if rawPath == "" {
		return rawPath
	}

	isWin := strings.EqualFold(osType, "windows") || isWindowsDrivePath(rawPath)
	expanded := expandHomeTilde(rawPath, isWin)
	expanded = expandKnownWinMacros(expanded, isWin)
	expanded = expandAllEnvVariables(expanded)

	return normalizePathSeparators(expanded, isWin)
}

func isWindowsDrivePath(p string) bool {
	return len(p) >= 2 && p[1] == ':'
}

func expandHomeTilde(p string, isWin bool) string {
	if p == "~" {
		return resolveHomeDir(isWin)
	}
	if strings.HasPrefix(p, "~/") || strings.HasPrefix(p, "~\\") {
		return filepath.Join(resolveHomeDir(isWin), p[2:])
	}

	return p
}

func resolveHomeDir(isWin bool) string {
	home, err := os.UserHomeDir()
	if err == nil && home != "" {
		return home
	}
	if isWin {
		return "C:\\Users\\Default"
	}

	return "/root"
}

func expandKnownWinMacros(p string, isWin bool) string {
	winDir := "C:\\Windows"
	winDrive := "C:"
	tempDir := os.TempDir()
	appData := resolveAppDataDir()

	replacer := strings.NewReplacer(
		"%win%", winDir, "%WIN%", winDir,
		"%win-drive%", winDrive, "%WIN-DRIVE%", winDrive,
		"%temp%", tempDir, "%TEMP%", tempDir,
		"%appdata%", appData, "%APPDATA%", appData,
	)

	return replacer.Replace(p)
}

func resolveAppDataDir() string {
	if appData := os.Getenv("APPDATA"); appData != "" {
		return appData
	}
	home, _ := os.UserHomeDir()

	return filepath.Join(home, "AppData", "Roaming")
}

func expandAllEnvVariables(p string) string {
	expanded := os.ExpandEnv(p)

	return winEnvRegex.ReplaceAllStringFunc(expanded, func(m string) string {
		varName := m[1 : len(m)-1]
		if val := os.Getenv(varName); val != "" {
			return val
		}

		return m
	})
}

func normalizePathSeparators(p string, isWin bool) string {
	if isWin {
		return strings.ReplaceAll(p, "/", "\\")
	}

	return strings.ReplaceAll(p, "\\", "/")
}
