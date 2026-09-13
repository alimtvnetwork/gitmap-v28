package cmdinstall

import (
	"os"
	"path/filepath"
	"runtime"
)

func getAgyCandidatePaths() []string {
	home, _ := os.UserHomeDir()
	list := []string{
		filepath.Join(home, ".local", "bin", "agy"),
		filepath.Join(home, ".local", "bin", "antigravity"),
		filepath.Join(home, ".local", "agy", "bin", "agy"),
		filepath.Join(home, ".antigravity", "bin", "agy"),
		filepath.Join(home, ".agy", "bin", "agy"),
		"/usr/local/bin/agy",
	}
	if runtime.GOOS == "windows" {
		localApp := os.Getenv("LOCALAPPDATA")
		list = append(list,
			filepath.Join(localApp, "agy", "bin", "agy.exe"),
			filepath.Join(home, ".antigravity", "bin", "agy.exe"),
		)
	}
	return list
}

func uninstallAgyFiles(purge bool) {
	for _, p := range getAgyCandidatePaths() {
		removeFileIfExists(p)
	}
	if !purge {
		return
	}
	home, _ := os.UserHomeDir()
	removeDirIfExists(filepath.Join(home, ".cache", "antigravity"))
	removeDirIfExists(filepath.Join(home, ".gemini", "antigravity-cli"))
	if runtime.GOOS == "windows" {
		removeDirIfExists(filepath.Join(os.Getenv("LOCALAPPDATA"), "antigravity"))
		removeDirIfExists(filepath.Join(os.Getenv("LOCALAPPDATA"), "agy"))
	}
}

func getAntigravityAppPaths() []string {
	home, _ := os.UserHomeDir()
	list := []string{
		"/opt/antigravity",
		filepath.Join(home, ".local", "share", "antigravity"),
		filepath.Join(home, ".local", "bin", "antigravity"),
		filepath.Join(home, ".local", "share", "applications", "antigravity.desktop"),
	}
	if runtime.GOOS == "windows" {
		localApp := os.Getenv("LOCALAPPDATA")
		progFiles := os.Getenv("ProgramFiles")
		list = append(list,
			filepath.Join(localApp, "Programs", "Antigravity"),
			filepath.Join(localApp, "Programs", "antigravity"),
			filepath.Join(progFiles, "Antigravity"),
		)
	}
	return list
}

func uninstallAntigravityFiles(purge bool) {
	for _, p := range getAntigravityAppPaths() {
		removeDirIfExists(p)
		removeFileIfExists(p)
	}
	if runtime.GOOS != "windows" {
		home, _ := os.UserHomeDir()
		desktopDir := filepath.Join(home, ".local", "share", "applications")
		updateDesktopDatabase(desktopDir)
	}
}

func uninstallScriptsDirectory() {
	removeDirIfExists(`D:\gitmap-scripts`)
	home, _ := os.UserHomeDir()
	removeDirIfExists(filepath.Join(home, "Desktop", "gitmap-scripts"))
	removeDirIfExists(filepath.Join(home, "gitmap-scripts"))
}

func uninstallAgManagerFiles() {
	home, _ := os.UserHomeDir()
	removeFileIfExists(filepath.Join(home, ".local", "bin", "ag-manager"))
	removeFileIfExists("/usr/local/bin/ag-manager")
	if runtime.GOOS == "windows" {
		removeDirIfExists(filepath.Join(os.Getenv("LOCALAPPDATA"), "Programs", "antigravity-manager"))
	}
}

func uninstallScriptToolFiles(tool string) {
	home, _ := os.UserHomeDir()
	bin := tool
	removeFileIfExists(filepath.Join(home, ".local", "bin", bin))
	removeFileIfExists(filepath.Join("/usr/local/bin", bin))
	if runtime.GOOS == "windows" {
		localApp := os.Getenv("LOCALAPPDATA")
		removeFileIfExists(filepath.Join(localApp, tool, tool+".exe"))
		removeFileIfExists(filepath.Join(home, "bin", tool+".exe"))
	}
}

func uninstallArchiveAppFiles(tool string) {
	home, _ := os.UserHomeDir()
	removeFileIfExists(filepath.Join(home, ".local", "bin", tool))
	removeFileIfExists(filepath.Join("/usr/local/bin", tool))
	removeDirIfExists(filepath.Join(home, ".local", "share", tool))
	removeDirIfExists(filepath.Join("/opt", tool))
	desktopFile := filepath.Join(home, ".local", "share", "applications", tool+".desktop")
	removeFileIfExists(desktopFile)
	if runtime.GOOS != "windows" {
		updateDesktopDatabase(filepath.Dir(desktopFile))
	}
}

