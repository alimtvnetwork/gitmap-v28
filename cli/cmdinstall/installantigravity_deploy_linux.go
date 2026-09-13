//go:build linux

package cmdinstall

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/tempdir"
)

func isPrivilegedLinuxUser() bool {
	return os.Geteuid() == 0
}

func resolveFallbackLinuxDirs() (string, string, string) {
	home, _ := os.UserHomeDir()

	return filepath.Join(home, ".local", "share", "antigravity"),
		filepath.Join(home, ".local", "bin"),
		filepath.Join(home, ".local", "share", "applications")
}

func resolveLinuxInstallDirs() (string, string, string) {
	if isPrivilegedLinuxUser() {
		return "/opt/antigravity", "/usr/local/bin", "/usr/share/applications"
	}

	return resolveFallbackLinuxDirs()
}

func prepareLinuxTargetDirs(reqInstall, reqBin, reqDesktop string) (string, string, string, error) {
	if err := os.MkdirAll(reqInstall, 0755); err == nil {
		return reqInstall, reqBin, reqDesktop, nil
	}

	fallbackInstall, fallbackBin, fallbackDesktop := resolveFallbackLinuxDirs()
	if err := os.MkdirAll(fallbackInstall, 0755); err != nil {
		return "", "", "", apperror.WrapSimple(err, "prepareLinuxTargetDirs")
	}

	return fallbackInstall, fallbackBin, fallbackDesktop, nil
}

func runTarExtractCommand(archivePath, destDir string) error {
	cmd := exec.Command("tar", "-xzf", archivePath, "-C", destDir, "--strip-components=1")
	if err := cmd.Run(); err == nil {
		return nil
	}

	fallbackCmd := exec.Command("tar", "-xzf", archivePath, "-C", destDir)

	return fallbackCmd.Run()
}

func extractAntigravityTarballWithRetry(archivePath, destDir string) error {
	var lastErr error
	for attempt := 1; attempt <= 3; attempt++ {
		lastErr = runTarExtractCommand(archivePath, destDir)
		if lastErr == nil {
			return nil
		}

		time.Sleep(time.Duration(attempt) * 200 * time.Millisecond)
	}

	return apperror.WrapSimple(lastErr, "antigravity.extractTarball")
}

func isValidBinaryCandidate(p string) bool {
	info, err := os.Stat(p)
	if err != nil {
		return false
	}

	return !info.IsDir() && info.Size() > 0
}

func searchDirBinary(dir string) string {
	candidates := []string{"antigravity", "Antigravity"}
	for _, name := range candidates {
		targetPath := filepath.Join(dir, name)
		if isValidBinaryCandidate(targetPath) {
			return targetPath
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

		subPath := filepath.Join(dir, e.Name())
		if p := searchDirBinary(subPath); p != "" {
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

func configureLinuxBinaryPermissions(binFile string) error {
	if err := os.Chmod(binFile, 0755); err != nil {
		return apperror.WrapSimple(err, "chmod.antigravity")
	}

	return nil
}

func configureChromeSandbox(installDir string) {
	if !isPrivilegedLinuxUser() {
		return
	}

	sandboxPath := filepath.Join(installDir, "chrome-sandbox")
	if _, err := os.Stat(sandboxPath); err != nil {
		return
	}

	_ = os.Chmod(sandboxPath, 04755)
}

func symlinkAntigravityBinary(targetBinary, symlinkPath string) error {
	_ = os.Remove(symlinkPath)
	if err := os.Symlink(targetBinary, symlinkPath); err != nil {
		return apperror.WrapSimple(err, "symlink.create")
	}

	return nil
}

func createDualLinuxSymlinks(binFile, binDir string) error {
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return apperror.WrapSimple(err, "mkdir.binDir")
	}

	if err := symlinkAntigravityBinary(binFile, filepath.Join(binDir, "antigravity")); err != nil {
		return err
	}

	return symlinkAntigravityBinary(binFile, filepath.Join(binDir, "agy"))
}

func createDualLinuxSymlinksWithFallback(binFile, preferredBinDir string) (string, error) {
	if err := createDualLinuxSymlinks(binFile, preferredBinDir); err == nil {
		return preferredBinDir, nil
	}

	home, _ := os.UserHomeDir()
	fallbackBinDir := filepath.Join(home, ".local", "bin")
	if err := createDualLinuxSymlinks(binFile, fallbackBinDir); err != nil {
		return "", err
	}

	return fallbackBinDir, nil
}

func createAntigravityDesktopEntry(binPath, iconPath, desktopFile string) error {
	content := fmt.Sprintf("[Desktop Entry]\nName=Google Antigravity\nComment=AI-First Development Platform & Agent Orchestration IDE\nGenericName=Text Editor / IDE\nExec=%s %%U\nIcon=%s\nType=Application\nStartupNotify=true\nStartupWMClass=Antigravity\nCategories=Development;IDE;Utility;\n", binPath, iconPath)

	if err := os.WriteFile(desktopFile, []byte(content), 0644); err != nil {
		return apperror.WrapSimple(err, "write.desktopEntry")
	}

	return nil
}

func installLinuxDesktopLauncher(binPath, installDir, desktopDir string) error {
	if err := os.MkdirAll(desktopDir, 0755); err != nil {
		home, _ := os.UserHomeDir()
		desktopDir = filepath.Join(home, ".local", "share", "applications")
		_ = os.MkdirAll(desktopDir, 0755)
	}

	icon := resolveAppIconPath(installDir)
	desktopFile := filepath.Join(desktopDir, "antigravity.desktop")
	err := createAntigravityDesktopEntry(binPath, icon, desktopFile)
	updateDesktopDatabase(desktopDir)

	return err
}

func linkAndLaunchLinuxDesktop(binFile, installDir, binDir, desktopDir string) error {
	activeBinDir, errSym := createDualLinuxSymlinksWithFallback(binFile, binDir)
	if errSym != nil {
		return errSym
	}

	ensureDirInPath(activeBinDir)

	return installLinuxDesktopLauncher(binFile, installDir, desktopDir)
}

func completeLinuxDeployment(installDir, binDir, desktopDir string) error {
	binFile, err := locateExtractedLinuxBinary(installDir)
	if err != nil {
		cleanupBrokenLinuxArtifacts(installDir, binDir, desktopDir)
		return err
	}

	_ = configureLinuxBinaryPermissions(binFile)
	configureChromeSandbox(installDir)

	return linkAndLaunchLinuxDesktop(binFile, installDir, binDir, desktopDir)
}

func deployAntigravityDesktopLinux(reqInstall, reqBin, reqDesktop, archivePath string) error {
	installDir, binDir, desktopDir, errDir := prepareLinuxTargetDirs(reqInstall, reqBin, reqDesktop)
	if errDir != nil {
		return errDir
	}

	cleanupBrokenLinuxArtifacts(installDir, binDir, desktopDir)
	if err := extractAntigravityTarballWithRetry(archivePath, installDir); err != nil {
		cleanupBrokenLinuxArtifacts(installDir, binDir, desktopDir)
		return err
	}

	return completeLinuxDeployment(installDir, binDir, desktopDir)
}

func installAntigravityDesktopPlatform(opts installOptions) error {
	url := getAntigravityDesktopDownloadUrl("linux")
	tempArchive := filepath.Join(tempdir.RepoTempDir("downloads"), "Antigravity.tar.gz")
	defer os.Remove(tempArchive)

	if err := downloadFileToDest(url, tempArchive); err != nil {
		return apperror.WrapSimple(err, "download.antigravity")
	}

	reqInstall, reqBin, reqDesktop := resolveLinuxInstallDirs()

	return deployAntigravityDesktopLinux(reqInstall, reqBin, reqDesktop, tempArchive)
}

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
