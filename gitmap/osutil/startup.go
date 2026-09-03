package osutil

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

// AddToStartup registers the binary for OS startup.
func AddToStartup(binPath string) error {
	isWindows := runtime.GOOS == "windows"
	isLinux := runtime.GOOS == "linux"

	if isWindows {
		return addWindowsStartup(binPath)
	}

	if isLinux {
		return addLinuxStartup(binPath)
	}

	return apperror.NewSimple("AddToStartup", "UNSUPPORTED_OS")
}

// addWindowsStartup adds to Windows Registry Run keys.
func addWindowsStartup(binPath string) error {
	cmd := exec.Command("reg", "add", `HKCU\Software\Microsoft\Windows\CurrentVersion\Run`, "/v", "GitmapMacro", "/t", "REG_SZ", "/d", binPath, "/f")
	if err := cmd.Run(); err != nil {
		return apperror.WrapSimple(err, "addWindowsStartup")
	}

	return nil
}

// addLinuxStartup adds to ~/.config/autostart/.
func addLinuxStartup(binPath string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return apperror.WrapSimple(err, "addLinuxStartup")
	}

	return writeDesktopFile(homeDir, binPath)
}

// writeDesktopFile handles the creation of the .desktop file.
func writeDesktopFile(homeDir, binPath string) error {
	autostartDir := filepath.Join(homeDir, ".config", "autostart")
	if err := os.MkdirAll(autostartDir, 0755); err != nil {
		return apperror.WrapSimple(err, "mkdirAutostart")
	}

	desktopFile := filepath.Join(autostartDir, "gitmap.desktop")
	content := fmt.Sprintf("[Desktop Entry]\nType=Application\nExec=%s\nHidden=false\nName=Gitmap\n", binPath)

	if err := os.WriteFile(desktopFile, []byte(content), 0644); err != nil {
		return apperror.WrapSimple(err, "writeDesktopFile")
	}

	return nil
}
