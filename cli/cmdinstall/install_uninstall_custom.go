package cmdinstall

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// IsCustomStandaloneTool reports whether tool is installed outside package managers.
func IsCustomStandaloneTool(tool string) bool {
	canonical := resolveToolAlias(tool)
	switch canonical {
	case constants.ToolAgy, constants.ToolAntigravity,
		constants.ToolAgManager, constants.ToolScriptsFixer,
		constants.ToolCodingGuidelines, constants.ToolMacroAhk,
		constants.ToolCtx, constants.ToolVSCodeCtx,
		constants.ToolPwshCtx, constants.ToolAgCtx,
		constants.ToolScripts,
		constants.ToolBeyondCompare, constants.ToolBeyondCompare4, constants.ToolBeyondCompare5:
		return true
	default:
		return hasInstalledArchiveFiles(canonical)
	}
}

func hasInstalledArchiveFiles(tool string) bool {
	home, _ := os.UserHomeDir()
	candidates := []string{
		filepath.Join(home, ".local", "share", tool),
		filepath.Join("/opt", tool),
		filepath.Join(home, ".local", "bin", tool),
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return true
		}
	}
	return false
}

func dispatchCustomRemoval(canonical string, purge bool) {
	switch canonical {
	case constants.ToolAgy:
		uninstallAgyFiles(purge)
	case constants.ToolAntigravity:
		uninstallAntigravityFiles(purge)
	case constants.ToolAgManager:
		uninstallAgManagerFiles()
	case constants.ToolCtx, constants.ToolVSCodeCtx, constants.ToolPwshCtx, constants.ToolAgCtx:
		runUninstallCtx()
	case constants.ToolScripts:
		uninstallScriptsDirectory()
	case constants.ToolBeyondCompare, constants.ToolBeyondCompare4, constants.ToolBeyondCompare5:
		uninstallBeyondCompare(canonical, purge)
	default:
		uninstallScriptToolFiles(canonical)
		uninstallArchiveAppFiles(canonical)
	}
}

// RunUninstallCustomTool removes files for custom tools and purges DB entries.
func RunUninstallCustomTool(tool string, purge bool) error {
	canonical := resolveToolAlias(tool)
	fmt.Printf("Removing %s...\n", canonical)

	dispatchCustomRemoval(canonical, purge)

	cleanToolFromDatabases(canonical)
	if tool != canonical {
		cleanToolFromDatabases(tool)
	}

	fmt.Printf(constants.MsgUninstallSuccess, tool)
	return nil
}

func removeFileIfExists(path string) {
	if path == "" {
		return
	}
	_ = os.Remove(path)
}

func removeDirIfExists(path string) {
	if path == "" {
		return
	}
	_ = os.RemoveAll(path)
}

func cleanToolFromDatabases(name string) {
	splitDB, err := store.OpenInstallationSplitDB()
	if err == nil {
		defer splitDB.Close()
		_ = splitDB.RemoveInstalledTool(name)
	}
}

func uninstallBeyondCompare(canonical string, purge bool) {
	if runtime.GOOS == "windows" {
		uninstallBeyondCompareWindows(canonical, purge)
		return
	}

	uninstallBeyondCompareLinux(canonical, purge)
}

func uninstallBeyondCompareWindows(canonical string, purge bool) {
	uninstaller := findBCWindowsUninstaller(canonical)
	if uninstaller != "" {
		cmd := exec.Command(uninstaller, "/VERYSILENT", "/NORESTART", "/SUPPRESSMSGBOXES")
		_ = cmd.Run()
	}

	if purge {
		purgeBeyondCompareAppData()
	}
}

func findBCWindowsUninstaller(canonical string) string {
	paths := getBCWindowsUninstallerPaths(canonical)
	for _, p := range paths {
		if isFileExist(p) {
			return p
		}
	}

	return ""
}

func getBCWindowsUninstallerPaths(canonical string) []string {
	if canonical == constants.ToolBeyondCompare4 {
		return []string{
			`C:\Program Files\Beyond Compare 4\unins000.exe`,
			`C:\Program Files (x86)\Beyond Compare 4\unins000.exe`,
		}
	}

	if canonical == constants.ToolBeyondCompare5 {
		return []string{
			`C:\Program Files\Beyond Compare 5\unins000.exe`,
		}
	}

	return []string{
		`C:\Program Files\Beyond Compare 5\unins000.exe`,
		`C:\Program Files\Beyond Compare 4\unins000.exe`,
		`C:\Program Files (x86)\Beyond Compare 4\unins000.exe`,
	}
}

func purgeBeyondCompareAppData() {
	appData := os.Getenv("APPDATA")
	if appData == "" {
		return
	}

	removeDirIfExists(filepath.Join(appData, "Scooter Software"))
}

func uninstallBeyondCompareLinux(canonical string, purge bool) {
	if isBinaryOnPath("apt-get") {
		_ = exec.Command("sudo", "apt-get", "remove", "-y", "bcompare").Run()
	}

	removeBCLinuxDirs(canonical)
	removeBCLinuxSymlinks()
	if purge {
		purgeBCLinuxConfig()
	}
}

func removeBCLinuxDirs(canonical string) {
	if canonical == constants.ToolBeyondCompare4 {
		removeDirIfExists("/opt/beyondcompare4")
		return
	}

	if canonical == constants.ToolBeyondCompare5 {
		removeDirIfExists("/opt/beyondcompare5")
		return
	}

	removeDirIfExists("/opt/beyondcompare5")
	removeDirIfExists("/opt/beyondcompare4")
	removeDirIfExists("/opt/beyondcompare")
}

func removeBCLinuxSymlinks() {
	symlinks := []string{
		"/usr/local/bin/bcompare",
		"/usr/local/bin/bcompare4",
		"/usr/local/bin/bcompare5",
	}

	for _, link := range symlinks {
		removeFileIfExists(link)
	}
}

func purgeBCLinuxConfig() {
	home, _ := os.UserHomeDir()
	if home == "" {
		return
	}

	removeDirIfExists(filepath.Join(home, ".config", "bcompare"))
}
