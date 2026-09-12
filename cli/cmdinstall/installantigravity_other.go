//go:build !windows

package cmdinstall

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func getAntigravityLinuxCandidatePaths() []string {
	home, _ := os.UserHomeDir()

	return []string{
		"/opt/antigravity/antigravity",
		filepath.Join(home, ".local", "share", "antigravity", "antigravity"),
		"/usr/local/bin/antigravity",
		filepath.Join(home, ".local", "bin", "antigravity"),
	}
}

func findInstalledAntigravityDesktopPath() (string, bool) {
	for _, p := range getAntigravityLinuxCandidatePaths() {
		if _, err := os.Stat(p); err == nil {
			return p, true
		}
	}

	if p, err := exec.LookPath("antigravity"); err == nil {
		return p, true
	}

	return "", false
}

func resolveLinuxInstallDirs() (string, string, string) {
	if os.Geteuid() == 0 {
		return "/opt/antigravity", "/usr/local/bin", "/usr/share/applications"
	}

	home, _ := os.UserHomeDir()

	return filepath.Join(home, ".local", "share", "antigravity"),
		filepath.Join(home, ".local", "bin"),
		filepath.Join(home, ".local", "share", "applications")
}

func extractAntigravityTarball(archivePath, destDir string) error {
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}

	cmd := exec.Command("tar", "-xzf", archivePath, "-C", destDir, "--strip-components=1")
	if err := cmd.Run(); err != nil {
		fallbackCmd := exec.Command("tar", "-xzf", archivePath, "-C", destDir)

		return fallbackCmd.Run()
	}

	return nil
}

func symlinkAntigravityBinary(targetBinary, symlinkPath string) error {
	_ = os.Remove(symlinkPath)

	return os.Symlink(targetBinary, symlinkPath)
}

func createAntigravityDesktopEntry(binPath, iconPath, desktopFile string) error {
	content := fmt.Sprintf("[Desktop Entry]\nName=Google Antigravity\nComment=AI-First Development Platform & Agent Orchestration IDE\nGenericName=Text Editor / IDE\nExec=%s %%U\nIcon=%s\nType=Application\nStartupNotify=true\nStartupWMClass=Antigravity\nCategories=Development;IDE;Utility;\n", binPath, iconPath)

	return os.WriteFile(desktopFile, []byte(content), 0644)
}

func deployAntigravityDesktopLinux(installDir, binDir, desktopDir, archivePath string) error {
	if err := extractAntigravityTarball(archivePath, installDir); err != nil {
		return err
	}

	binFile := filepath.Join(installDir, "antigravity")
	_ = os.Chmod(binFile, 0755)
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return err
	}

	if err := symlinkAntigravityBinary(binFile, filepath.Join(binDir, "antigravity")); err != nil {
		return err
	}

	_ = os.MkdirAll(desktopDir, 0755)
	icon := filepath.Join(installDir, "resources", "app", "resources", "icon.png")

	return createAntigravityDesktopEntry(filepath.Join(binDir, "antigravity"), icon, filepath.Join(desktopDir, "antigravity.desktop"))
}

func installAntigravityDesktopPlatform(opts installOptions) error {
	url := getAntigravityDesktopDownloadUrl("linux")
	tempArchive := filepath.Join(os.TempDir(), "Antigravity.tar.gz")
	defer os.Remove(tempArchive)

	if err := downloadFileToDest(url, tempArchive); err != nil {
		return err
	}

	installDir, binDir, desktopDir := resolveLinuxInstallDirs()

	return deployAntigravityDesktopLinux(installDir, binDir, desktopDir, tempArchive)
}
