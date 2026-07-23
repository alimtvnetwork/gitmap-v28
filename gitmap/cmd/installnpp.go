package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

const maxNppFileSize = 10 * 1024 * 1024 // 10 MB limit per extracted file

// resolveNppInstallName maps install-npp to the npp binary for install.
func resolveNppInstallName(tool string) string {
	if tool == constants.ToolNppInstall {
		return constants.ToolNpp
	}

	return tool
}

// runNppSettingsOnly syncs Notepad++ settings without installing the binary.
func runNppSettingsOnly() {
	fmt.Print(constants.MsgInstallNppSkipBin)
	runNppSettings()
}

// runNppSettings syncs Notepad++ settings to the AppData directory.
func runNppSettings() {
	fmt.Print(constants.MsgInstallNppSettings)

	target := nppSettingsTarget()
	if target == "" {
		return
	}

	err := os.MkdirAll(target, 0o755)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrNppDirCreate, target, err)

		return
	}

	extractNppSettingsZip(target)
}

// nppSettingsTarget returns the Notepad++ AppData settings path.
func nppSettingsTarget() string {
	if runtime.GOOS != "windows" {
		fmt.Fprintf(os.Stderr, constants.ErrNppWindowsOnly, runtime.GOOS)

		return ""
	}

	appData := os.Getenv("APPDATA")
	if appData == "" {
		fmt.Fprint(os.Stderr, constants.ErrNppNoAppData)

		return ""
	}

	return filepath.Join(appData, "Notepad++")
}

// resolveSettingsPath resolves a settings file path, searching multiple
// locations relative to the binary and CWD. This supports both the legacy
// data/ layout and the current settings/ layout.
func resolveSettingsPath(subpaths ...string) string {
	exe, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ! Could not resolve executable path: %v\n", err)

		return subpaths[0]
	}

	realExe, err := filepath.EvalSymlinks(exe)
	if err != nil {
		realExe = exe
	}

	binDir := filepath.Dir(realExe)

	// Build candidate list: for each subpath, try binary-relative then CWD-relative.
	var candidates []string

	for _, subpath := range subpaths {
		candidates = append(candidates,
			filepath.Join(binDir, subpath),
			subpath,
		)
	}

	for _, candidate := range candidates {
		abs, absErr := filepath.Abs(candidate)
		if absErr != nil {
			abs = candidate
		}
		if _, statErr := os.Stat(candidate); statErr == nil {
			fmt.Printf("  -> Resolved path: %s\n", abs)

			return candidate
		}

		fmt.Printf("  -> Searched: %s (not found)\n", abs)
	}

	return candidates[0]
}

// resolveNppDataPath resolves the npp-settings data path relative to the binary.
// Searches settings/ folder first (current layout), then data/ (legacy).
func resolveNppDataPath(subpath string) string {
	return resolveSettingsPath(
		filepath.Join("settings", "01 - notepad++", subpath),
		filepath.Join("data", "npp-settings", subpath),
		filepath.Join("data", subpath),
	)
}
