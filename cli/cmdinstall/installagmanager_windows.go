//go:build windows

package cmdinstall

import (
	"os"
	"path/filepath"
)

func findAgManagerWindowsPath() string {
	localAppData := resolveLocalAppDataDir()
	candidates := []string{
		filepath.Join(localAppData, "Programs", "Antigravity.Tools", "Antigravity.Tools.exe"),
		filepath.Join(localAppData, "Programs", "antigravity-tools", "Antigravity.Tools.exe"),
		filepath.Join(localAppData, "Programs", "Antigravity-Manager", "Antigravity-Manager.exe"),
		filepath.Join(localAppData, "Programs", "ag-manager", "ag-manager.exe"),
	}

	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c
		}
	}

	return ""
}

func isAgManagerWindowsInstalled() bool {
	return findAgManagerWindowsPath() != ""
}
