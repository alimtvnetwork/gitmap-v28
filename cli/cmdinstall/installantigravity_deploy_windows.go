//go:build windows

package cmdinstall

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/sys/windows/registry"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/tempdir"
)

func buildNsisInstallerArgs(customDir string) []string {
	args := []string{"/S"}
	if customDir != "" {
		args = append(args, "/D="+customDir)
	}

	return args
}

func runAntigravitySilentInstaller(installerPath, customDir string) error {
	args := buildNsisInstallerArgs(customDir)
	cmd := exec.Command(installerPath, args...)

	return cmd.Run()
}

func runAntigravityInteractiveInstaller(installerPath string) error {
	cmd := exec.Command(installerPath)

	return cmd.Run()
}

func executeWindowsInstallerWithFallback(installerPath, customDir string) error {
	if err := runAntigravitySilentInstaller(installerPath, customDir); err == nil {
		return nil
	}

	if err := runAntigravityInteractiveInstaller(installerPath); err != nil {
		return apperror.WrapSimple(err, "windows.installerInteractive")
	}

	return nil
}

func checkAntigravityExeExists(exePath string) bool {
	info, err := os.Stat(exePath)
	if err != nil {
		return false
	}

	return !info.IsDir() && info.Size() > 0
}

func waitForAntigravityExe(exePath string, maxSeconds int) bool {
	for i := 0; i < maxSeconds; i++ {
		if hasExe := checkAntigravityExeExists(exePath); hasExe {
			return true
		}

		time.Sleep(1 * time.Second)
	}

	return false
}

func resolveExpectedWindowsExe(customDir string) string {
	if customDir != "" {
		return filepath.Join(customDir, "Antigravity.exe")
	}

	return getAntigravityDesktopWindowsExePath()
}

func verifyAntigravityInstallation(expectedExe string) error {
	if isFound := waitForAntigravityExe(expectedExe, 30); isFound {
		return nil
	}

	if _, isInstalled := findInstalledAntigravityDesktopPath(); isInstalled {
		return nil
	}

	return apperror.NewSimple("timeout waiting for Antigravity.exe to install", "E9000")
}

func resolveWindowsAgyBinDir() string {
	localApp := resolveLocalAppDataDir()

	return filepath.Join(localApp, "agy", "bin")
}

func writeCmdShim(shimPath, targetExe string) error {
	content := fmt.Sprintf("@echo off\r\n\"%s\" %%*\r\n", targetExe)
	if err := os.WriteFile(shimPath, []byte(content), 0755); err != nil {
		return apperror.WrapSimple(err, "writeCmdShim")
	}

	return nil
}

func writePs1Shim(shimPath, targetExe string) error {
	content := fmt.Sprintf("& \"%s\" @args\r\n", targetExe)
	if err := os.WriteFile(shimPath, []byte(content), 0755); err != nil {
		return apperror.WrapSimple(err, "writePs1Shim")
	}

	return nil
}

func writeDualShimsForName(binDir, shimName, targetExe string) error {
	cmdPath := filepath.Join(binDir, shimName+".cmd")
	if err := writeCmdShim(cmdPath, targetExe); err != nil {
		return err
	}

	ps1Path := filepath.Join(binDir, shimName+".ps1")

	return writePs1Shim(ps1Path, targetExe)
}

func createDualWindowsShims(targetExe string) error {
	binDir := resolveWindowsAgyBinDir()
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return apperror.WrapSimple(err, "mkdir.windowsAgyBin")
	}

	if err := writeDualShimsForName(binDir, "antigravity", targetExe); err != nil {
		return err
	}

	return writeDualShimsForName(binDir, "agy", targetExe)
}

func hasDirInPathString(pathStr, dir string) bool {
	parts := strings.Split(pathStr, ";")
	targetLow := strings.ToLower(filepath.Clean(dir))
	for _, part := range parts {
		if strings.ToLower(filepath.Clean(strings.TrimSpace(part))) == targetLow {
			return true
		}
	}

	return false
}

func joinWindowsPaths(curPath, dir string) string {
	trimmed := strings.TrimSpace(curPath)
	if trimmed == "" {
		return dir
	}

	if strings.HasSuffix(trimmed, ";") {
		return trimmed + dir
	}

	return trimmed + ";" + dir
}

func appendDirToUserPath(dir string) error {
	ensureDirInPath(dir)
	key, err := registry.OpenKey(registry.CURRENT_USER, "Environment", registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		return apperror.WrapSimple(err, "openUserEnvRegKey")
	}
	defer key.Close()

	curPath, _, _ := key.GetStringValue("Path")
	if hasDirInPathString(curPath, dir) {
		return nil
	}

	return key.SetStringValue("Path", joinWindowsPaths(curPath, dir))
}

func finalizeWindowsShimsAndPath(activeExe string) error {
	if err := createDualWindowsShims(activeExe); err != nil {
		return err
	}

	binDir := resolveWindowsAgyBinDir()
	_ = appendDirToUserPath(binDir)

	return nil
}

func resolveActiveWindowsExe(expectedExe string) string {
	activeExe, isFound := findInstalledAntigravityDesktopPath()
	if !isFound {
		return expectedExe
	}

	return activeExe
}

func deployAntigravityDesktopWindows(installerPath, customDir string) error {
	expectedExe := resolveExpectedWindowsExe(customDir)
	if err := executeWindowsInstallerWithFallback(installerPath, customDir); err != nil {
		return err
	}

	if err := verifyAntigravityInstallation(expectedExe); err != nil {
		return err
	}

	activeExe := resolveActiveWindowsExe(expectedExe)

	return finalizeWindowsShimsAndPath(activeExe)
}

func installAntigravityDesktopPlatform(opts installOptions) error {
	url := getAntigravityDesktopDownloadUrl("windows")
	tempInstaller := filepath.Join(tempdir.RepoTempDir("downloads"), "Antigravity-x64.exe")
	defer os.Remove(tempInstaller)

	if err := downloadFileToDest(url, tempInstaller); err != nil {
		return apperror.WrapSimple(err, "download.windowsInstaller")
	}

	return deployAntigravityDesktopWindows(tempInstaller, opts.Prefix)
}
