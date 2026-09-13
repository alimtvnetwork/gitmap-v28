//go:build windows

package cmdinstall

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/tempdir"
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
	for _, name := range []string{"antigravity", "Antigravity", "agy"} {
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

func runAntigravitySilentInstaller(installerPath string) error {
	cmd := exec.Command(installerPath, "/S")

	return cmd.Run()
}

func waitForAntigravityExe(exePath string, maxSeconds int) bool {
	for i := 0; i < maxSeconds; i++ {
		if _, err := os.Stat(exePath); err == nil {
			return true
		}

		time.Sleep(1 * time.Second)
	}

	return false
}

func executeAndVerifyWindowsInstall(installerPath, exePath string) error {
	if err := runAntigravitySilentInstaller(installerPath); err != nil {
		return err
	}

	if isCreated := waitForAntigravityExe(exePath, 30); !isCreated {
		return fmt.Errorf("timeout waiting for Antigravity.exe to install")
	}

	return nil
}

func installAntigravityDesktopPlatform(opts installOptions) error {
	url := getAntigravityDesktopDownloadUrl("windows")
	tempInstaller := filepath.Join(tempdir.RepoTempDir("downloads"), "Antigravity-x64.exe")
	defer os.Remove(tempInstaller)

	if err := downloadFileToDest(url, tempInstaller); err != nil {
		return err
	}

	return executeAndVerifyWindowsInstall(tempInstaller, getAntigravityDesktopWindowsExePath())
}
