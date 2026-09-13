//go:build !windows

package cmdinstall

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/tempdir"
)

func getAntigravityLinuxCandidatePaths() []string {
	home, _ := os.UserHomeDir()

	return []string{
		"/opt/antigravity/antigravity",
		filepath.Join(home, ".local", "share", "antigravity", "antigravity"),
		filepath.Join(home, ".local", "share", "antigravity", "Antigravity"),
		filepath.Join(home, ".local", "bin", "antigravity"),
		"/usr/local/bin/antigravity",
	}
}

func findAntigravityInPathLinux() (string, bool) {
	if p, err := exec.LookPath("antigravity"); err == nil {
		return p, true
	}

	return "", false
}

func findInstalledAntigravityDesktopPath() (string, bool) {
	for _, p := range getAntigravityLinuxCandidatePaths() {
		if info, err := os.Stat(p); err == nil && !info.IsDir() && info.Size() > 0 {
			return p, true
		}
	}

	return findAntigravityInPathLinux()
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

func isValidBinaryCandidate(p string) bool {
	info, err := os.Stat(p)
	if err != nil || info.IsDir() {
		return false
	}

	return info.Size() > 0
}

func searchDirBinary(dir string) string {
	for _, name := range []string{"antigravity", "Antigravity"} {
		p := filepath.Join(dir, name)
		if isValidBinaryCandidate(p) {
			return p
		}
	}

	return ""
}

func searchSubdirBinary(dir string) string {
	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		if p := searchDirBinary(filepath.Join(dir, e.Name())); p != "" {
			return p
		}
	}

	return ""
}

func locateExtractedLinuxBinary(installDir string) (string, error) {
	if p := searchDirBinary(installDir); p != "" {
		return p, nil
	}
	if p := searchSubdirBinary(installDir); p != "" {
		return p, nil
	}

	return "", apperror.NewSimple("executable Antigravity binary not found in archive", "E9000")
}

func deployAntigravityDesktopLinux(installDir, binDir, desktopDir, archivePath string) error {
	cleanupBrokenLinuxArtifacts(installDir, binDir, desktopDir)
	if err := extractAntigravityTarball(archivePath, installDir); err != nil {
		cleanupBrokenLinuxArtifacts(installDir, binDir, desktopDir)
		return err
	}

	binFile, err := locateExtractedLinuxBinary(installDir)
	if err != nil {
		cleanupBrokenLinuxArtifacts(installDir, binDir, desktopDir)
		return err
	}

	_ = os.Chmod(binFile, 0755)
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return err
	}

	_ = symlinkAntigravityBinary(binFile, filepath.Join(binDir, "antigravity"))
	_ = symlinkAntigravityBinary(binFile, filepath.Join(binDir, "agy"))
	ensureDirInPath(binDir)

	_ = os.MkdirAll(desktopDir, 0755)
	icon := resolveAppIconPath(installDir)
	errDesktop := createAntigravityDesktopEntry(filepath.Join(binDir, "antigravity"), icon, filepath.Join(desktopDir, "antigravity.desktop"))
	updateDesktopDatabase(desktopDir)
	return errDesktop
}

func installAntigravityDesktopPlatform(opts installOptions) error {
	url := getAntigravityDesktopDownloadUrl("linux")
	tempArchive := filepath.Join(tempdir.RepoTempDir("downloads"), "Antigravity.tar.gz")
	defer os.Remove(tempArchive)

	if err := downloadFileToDest(url, tempArchive); err != nil {
		return err
	}

	installDir, binDir, desktopDir := resolveLinuxInstallDirs()
	return deployAntigravityDesktopLinux(installDir, binDir, desktopDir, tempArchive)
}
