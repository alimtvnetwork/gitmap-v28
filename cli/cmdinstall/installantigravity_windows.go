//go:build windows

package cmdinstall

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func getAntigravityDesktopWindowsExePath() string {
	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData == "" {
		home, _ := os.UserHomeDir()
		localAppData = filepath.Join(home, "AppData", "Local")
	}

	return filepath.Join(localAppData, "Programs", "Antigravity", "Antigravity.exe")
}

func findInstalledAntigravityDesktopPath() (string, bool) {
	exePath := getAntigravityDesktopWindowsExePath()
	if _, err := os.Stat(exePath); err == nil {
		return exePath, true
	}

	if p, err := exec.LookPath("antigravity"); err == nil {
		return p, true
	}

	return "", false
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
	tempInstaller := filepath.Join(os.TempDir(), "Antigravity-x64.exe")
	defer os.Remove(tempInstaller)

	if err := downloadFileToDest(url, tempInstaller); err != nil {
		return err
	}

	return executeAndVerifyWindowsInstall(tempInstaller, getAntigravityDesktopWindowsExePath())
}
