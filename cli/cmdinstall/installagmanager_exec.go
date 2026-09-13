package cmdinstall

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func isAgManagerInstalled() (string, bool) {
	for _, name := range []string{"ag-manager", "Antigravity.Tools", "Antigravity-Manager"} {
		if bin := resolveToolBinaryPath(name); bin != "" {
			return "installed", true
		}
	}

	if isAgManagerWindowsInstalled() {
		return "installed", true
	}

	return "", false
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
