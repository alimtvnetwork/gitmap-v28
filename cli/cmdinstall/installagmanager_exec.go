package cmdinstall

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func isAgManagerInstalled() (string, bool) {
	for _, name := range []string{"ag-manager", "Antigravity.Tools", "Antigravity-Manager"} {
		if bin := resolveToolBinaryPath(name); bin != "" {
			return "installed", true
		}
	}

	if runtime.GOOS == "windows" {
		if path := findAgManagerWindowsPath(); path != "" {
			return "installed", true
		}
	}

	return "", false
}

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

func isAppImageFile(path string) bool {
	low := strings.ToLower(path)

	return strings.HasSuffix(low, ".appimage")
}

func installAgManagerAppImage(srcPath string) error {
	dest, err := installUserBin("ag-manager", srcPath)
	if err != nil {
		return err
	}

	fmt.Printf("  ✓ Installed portable AppImage to %s\n", dest)

	return nil
}

func executeAgManagerInstaller(path string) error {
	if runtime.GOOS == "linux" && isAppImageFile(path) {
		return installAgManagerAppImage(path)
	}

	cmd := buildInstallerCommand(path)
	if cmd != nil {
		cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr

		return cmd.Run()
	}

	return nil
}

func buildInstallerCommand(path string) *exec.Cmd {
	switch runtime.GOOS {
	case "windows":
		return buildWindowsInstallerCmd(path)
	case "darwin":
		return exec.Command("open", path)
	case "linux":
		return buildLinuxInstallerCmd(path)
	}

	return nil
}

func buildWindowsInstallerCmd(path string) *exec.Cmd {
	if strings.HasSuffix(strings.ToLower(path), ".msi") {
		return exec.Command("msiexec", "/i", path, "/qn")
	}

	return exec.Command(path, "/S")
}

func buildLinuxInstallerCmd(path string) *exec.Cmd {
	if strings.HasSuffix(strings.ToLower(path), ".deb") {
		return exec.Command("sudo", "dpkg", "-i", path)
	}

	return exec.Command("sudo", "apt", "install", "-y", path)
}
