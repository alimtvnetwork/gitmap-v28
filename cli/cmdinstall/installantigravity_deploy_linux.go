//go:build linux

package cmdinstall

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/assets"
)

func isPrivilegedLinuxUser() bool {
	return os.Geteuid() == 0
}

func resolveFallbackLinuxDirs() (string, string, string) {
	home, _ := os.UserHomeDir()

	return filepath.Join(home, ".local", "share", "antigravity-ide"),
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
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return apperror.WrapSimple(err, "mkdir.destDir")
	}

	cmd := exec.Command("tar", "-xzf", archivePath, "-C", destDir, "--strip-components=1")
	out, err := cmd.CombinedOutput()
	if err == nil {
		return nil
	}

	fallbackCmd := exec.Command("tar", "-xzf", archivePath, "-C", destDir)
	fallbackOut, fallbackErr := fallbackCmd.CombinedOutput()
	if fallbackErr == nil {
		return nil
	}

	return fmt.Errorf("tar extraction failed: %w (strip: %s, direct: %s)", fallbackErr, string(out), string(fallbackOut))
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

	if err := symlinkAntigravityBinary(binFile, filepath.Join(binDir, "agy")); err != nil {
		return err
	}

	return symlinkAntigravityBinary(binFile, filepath.Join(binDir, "antigravity-ide"))
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

func deployAntigravityDesktopIcons(installDir string) {
	opts := IconDeployOptions{
		IconName:   AntigravityCanonicalIconName,
		SourceDir:  installDir,
		RawBytes:   assets.AntigravityDefaultIcon,
		IsUserOnly: !isPrivilegedLinuxUser(),
	}

	_ = DeployMultiResolutionIcons(opts)
}

func buildLauncherDeployConfig(binPath, desktopDir string) AntigravityLauncherDeployConfig {
	home, _ := os.UserHomeDir()
	userAppDir := filepath.Join(home, ".local", "share", "applications")
	if desktopDir != "" && desktopDir != "/usr/share/applications" {
		userAppDir = desktopDir
	}

	return AntigravityLauncherDeployConfig{
		BinPath:      binPath,
		UserAppDir:   userAppDir,
		DesktopDir:   filepath.Join(home, "Desktop"),
		IconName:     AntigravityCanonicalIconName,
		IsPrivileged: isPrivilegedLinuxUser(),
	}
}

func installLinuxDesktopLauncher(binPath, installDir, desktopDir string) error {
	deployAntigravityDesktopIcons(installDir)
	cfg := buildLauncherDeployConfig(binPath, desktopDir)
	_ = os.MkdirAll(cfg.UserAppDir, 0755)

	return DeployCanonicalAntigravityLauncher(cfg)
}

func linkAndLaunchLinuxDesktop(binFile, installDir, binDir, desktopDir string) error {
	activeBinDir, errSym := createDualLinuxSymlinksWithFallback(binFile, binDir)
	if errSym != nil {
		return errSym
	}

	ensureDirInPath(activeBinDir)

	return installLinuxDesktopLauncher(binFile, installDir, desktopDir)
}

func buildLauncherWrapperScript(installDir string) string {
	return fmt.Sprintf(`#!/bin/bash
REAL_SOURCE="$(readlink -f "${BASH_SOURCE[0]}" 2>/dev/null || realpath "${BASH_SOURCE[0]}" 2>/dev/null || echo "${BASH_SOURCE[0]}")"
DIR="$(cd "$(dirname "$REAL_SOURCE")" && pwd)"

if [ "$DIR" = "$HOME/.local/bin" ] || [ "$DIR" = "/usr/local/bin" ] || [ "$DIR" = "/usr/bin" ]; then
    [ -d "%s" ] && DIR="%s"
fi

export LD_LIBRARY_PATH="$DIR:$DIR/lib:${LD_LIBRARY_PATH:-}"
EXEC="$DIR/antigravity"
[ -f "$EXEC" ] || EXEC="$DIR/Antigravity"

export ELECTRON_OZONE_PLATFORM_HINT="auto"
export DONT_PROMPT_WSL_INSTALL=1

# Unconfined tarball Electron builds require --no-sandbox on modern Linux (AppArmor userns restriction)
exec "$EXEC" --no-sandbox "$@"
`, installDir, installDir)
}

func ensureLinuxLauncherWrapper(installDir string) (string, error) {
	wrapperPath := filepath.Join(installDir, "antigravity-ide.run")
	runnerSymlink := filepath.Join(installDir, "antigravity-runner.sh")
	runSymlink := filepath.Join(installDir, "antigravity.run")

	content := buildLauncherWrapperScript(installDir)
	if err := os.WriteFile(wrapperPath, []byte(content), 0755); err != nil {
		return "", apperror.WrapSimple(err, "write.launcherWrapper")
	}

	_ = os.Chmod(wrapperPath, 0755)
	_ = os.Remove(runnerSymlink)
	_ = os.Symlink(wrapperPath, runnerSymlink)
	_ = os.Remove(runSymlink)
	_ = os.Symlink(wrapperPath, runSymlink)

	return wrapperPath, nil
}

func ensureLinuxDirectorySymlink(installDir string) {
	home, _ := os.UserHomeDir()
	aliasDir := filepath.Join(home, ".local", "share", "antigravity")
	if aliasDir == installDir {
		return
	}

	if fi, err := os.Lstat(aliasDir); err == nil {
		if fi.Mode()&os.ModeSymlink != 0 {
			_ = os.Remove(aliasDir)
		} else if fi.IsDir() {
			entries, _ := os.ReadDir(aliasDir)
			if len(entries) == 0 {
				_ = os.Remove(aliasDir)
			}
		}
	}

	_ = os.Symlink(installDir, aliasDir)
}

func completeLinuxDeployment(installDir, binDir, desktopDir string) error {
	binFile, err := locateExtractedLinuxBinary(installDir)
	if err != nil {
		cleanupBrokenLinuxArtifacts(installDir, binDir, desktopDir)
		return err
	}

	_ = configureLinuxBinaryPermissions(binFile)
	configureChromeSandbox(installDir)

	targetToLaunch := binFile
	if wrapperPath, errWrapper := ensureLinuxLauncherWrapper(installDir); errWrapper == nil && wrapperPath != "" {
		targetToLaunch = wrapperPath
	}

	ensureLinuxDirectorySymlink(installDir)

	return linkAndLaunchLinuxDesktop(targetToLaunch, installDir, binDir, desktopDir)
}

func deployAntigravityDesktopLinux(reqInstall, reqBin, reqDesktop, archivePath string) error {
	installDir, binDir, desktopDir, errDir := prepareLinuxTargetDirs(reqInstall, reqBin, reqDesktop)
	if errDir != nil {
		return errDir
	}

	if err := extractAntigravityTarballWithRetry(archivePath, installDir); err != nil {
		return err
	}

	return completeLinuxDeployment(installDir, binDir, desktopDir)
}

func fetchLinuxTarballArchive(opts installOptions) (string, error) {
	params := ArchiveDownloadParams{
		URL:            getAntigravityDesktopDownloadUrl("linux"),
		IsDownloadMust: opts.IsDownloadMust,
		Verbose:        opts.Verbose,
	}

	return FetchOrReuseArchive(params)
}

func installAntigravityDesktopPlatform(opts installOptions) error {
	archivePath, err := fetchLinuxTarballArchive(opts)
	if err != nil {
		return apperror.WrapSimple(err, "fetch.linuxTarball")
	}

	reqInstall, reqBin, reqDesktop := resolveLinuxInstallDirs()
	if opts.Prefix != "" {
		reqInstall = opts.Prefix
	}

	return deployAntigravityDesktopLinux(reqInstall, reqBin, reqDesktop, archivePath)
}

func getAntigravityLinuxCandidatePaths() []string {
	home, _ := os.UserHomeDir()

	return []string{
		"/opt/antigravity/antigravity",
		filepath.Join(home, ".local", "share", "antigravity-ide", "antigravity-ide.run"),
		filepath.Join(home, ".local", "share", "antigravity-ide", "antigravity"),
		filepath.Join(home, ".local", "share", "antigravity", "antigravity-ide.run"),
		filepath.Join(home, ".local", "share", "antigravity", "antigravity"),
		filepath.Join(home, ".local", "share", "antigravity", "Antigravity"),
		filepath.Join(home, ".local", "bin", "antigravity-ide"),
		filepath.Join(home, ".local", "bin", "antigravity"),
		filepath.Join(home, ".local", "bin", "agy"),
		"/usr/local/bin/antigravity",
		"/usr/local/bin/agy",
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
