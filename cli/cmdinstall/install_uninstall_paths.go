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
		"/usr/local/bin/antigravity",
	}
	return appendPlatformAgyPaths(list, home)
}

func getWindowsAgyPaths(home, binDir string) []string {
	return []string{
		filepath.Join(binDir, "agy.exe"),
		filepath.Join(binDir, "agy.cmd"),
		filepath.Join(binDir, "agy.ps1"),
		filepath.Join(binDir, "antigravity.cmd"),
		filepath.Join(binDir, "antigravity.ps1"),
		filepath.Join(home, ".antigravity", "bin", "agy.exe"),
		binDir,
	}
}

func appendPlatformAgyPaths(list []string, home string) []string {
	if runtime.GOOS != "windows" {
		return list
	}
	binDir := filepath.Join(os.Getenv("LOCALAPPDATA"), "agy", "bin")
	return append(list, getWindowsAgyPaths(home, binDir)...)
}

func uninstallAgyFiles(hasPurge bool) {
	for _, p := range getAgyCandidatePaths() {
		removeFileIfExists(p)
	}
	purgeAgyDirectories(hasPurge)
	_ = PurgeAntigravityDualDatabase()
}

func purgeAgyDirectories(hasPurge bool) {
	if !hasPurge {
		return
	}
	home, _ := os.UserHomeDir()
	removeDirIfExists(filepath.Join(home, ".cache", "antigravity"))
	removeDirIfExists(filepath.Join(home, ".gemini", "antigravity-cli"))
	purgeWindowsAgyDirs(os.Getenv("LOCALAPPDATA"))
}

func purgeWindowsAgyDirs(localApp string) {
	if runtime.GOOS != "windows" {
		return
	}
	removeDirIfExists(filepath.Join(localApp, "antigravity"))
	removeDirIfExists(filepath.Join(localApp, "agy"))
}

func getAntigravityAppPaths() []string {
	home, _ := os.UserHomeDir()
	if runtime.GOOS == "windows" {
		return getWindowsAntigravityAppPaths(os.Getenv("LOCALAPPDATA"), os.Getenv("ProgramFiles"))
	}
	if runtime.GOOS == "darwin" {
		return getDarwinAntigravityAppPaths(home)
	}
	return getLinuxAntigravityAppPaths(home)
}

func getWindowsAntigravityAppPaths(localApp, progFiles string) []string {
	binDir := filepath.Join(localApp, "agy", "bin")
	return []string{
		filepath.Join(localApp, "Programs", "Antigravity"),
		filepath.Join(localApp, "Programs", "antigravity"),
		filepath.Join(progFiles, "Antigravity"),
		filepath.Join(binDir, "antigravity.cmd"),
		filepath.Join(binDir, "antigravity.ps1"),
		filepath.Join(binDir, "antigravity.exe"),
	}
}

func getLinuxAntigravityAppPaths(home string) []string {
	return []string{
		"/opt/antigravity",
		filepath.Join(home, ".local", "share", "antigravity"),
		filepath.Join(home, ".local", "bin", "antigravity"),
		"/usr/local/bin/antigravity",
		filepath.Join(home, ".local", "share", "applications", "antigravity.desktop"),
		"/usr/share/applications/antigravity.desktop",
	}
}

func getDarwinAntigravityAppPaths(home string) []string {
	return []string{
		"/Applications/Antigravity.app",
		filepath.Join(home, "Applications", "Antigravity.app"),
		"/usr/local/bin/antigravity",
		filepath.Join(home, ".local", "bin", "antigravity"),
	}
}

func uninstallAntigravityFiles(hasPurge bool) {
	for _, p := range getAntigravityAppPaths() {
		removeDirIfExists(p)
		removeFileIfExists(p)
	}
	refreshLinuxDesktopDatabases()
	_ = PurgeAntigravityDualDatabase()
}

func refreshLinuxDesktopDatabases() {
	if runtime.GOOS != "linux" {
		return
	}
	home, _ := os.UserHomeDir()
	updateDesktopDatabase(filepath.Join(home, ".local", "share", "applications"))
	updateDesktopDatabase("/usr/share/applications")
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
