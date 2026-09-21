package cmdinstall

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

const (
	bc4WinURL = "https://www.scootersoftware.com/files/BCompare-4.4.7.28397.exe"
	bc5WinURL = "https://www.scootersoftware.com/files/BCompare-5.0.3.30064.exe"
	bc4DebURL = "https://www.scootersoftware.com/files/bcompare-4.4.7.28397_amd64.deb"
	bc5DebURL = "https://www.scootersoftware.com/files/bcompare-5.0.3.30064_amd64.deb"
	bc4TarURL = "https://www.scootersoftware.com/files/bcompare-4.4.7.28397_x86_64.tar.gz"
	bc5TarURL = "https://www.scootersoftware.com/files/bcompare-5.0.3.30064_x86_64.tar.gz"
)

func handleBeyondCompareInstall(opts installOptions, defaultVersion int) {
	err := runInstallBeyondCompare(defaultVersion, opts)
	if err != nil {
		cliexit.HandleError(apperror.WrapSimple(err, "cmdinstall.installBeyondCompare"), 1)
	}
}

func runInstallBeyondCompare(defaultVersion int, opts installOptions) error {
	version := resolveBCVersion(opts.Version, defaultVersion)
	if runtime.GOOS == "windows" {
		return runInstallBeyondCompareWindows(version, opts)
	}

	return runInstallBeyondCompareLinux(version, opts)
}

func resolveBCVersion(versionStr string, defaultVersion int) int {
	if strings.Contains(versionStr, "4") {
		return 4
	}

	if strings.Contains(versionStr, "5") {
		return 5
	}

	return defaultVersion
}

func resolveBCToolName(version int) string {
	if version == 4 {
		return constants.ToolBeyondCompare4
	}

	return constants.ToolBeyondCompare5
}

func runInstallBeyondCompareWindows(version int, opts installOptions) error {
	if isBCInstalledWindows(version) {
		fmt.Printf("  ✓ Beyond Compare %d is already installed.\n", version)
		return nil
	}

	announceBCInstallPlan(version, "Inno Setup (Windows)")
	if opts.DryRun {
		return nil
	}

	return executeBCInstallWindows(version)
}

func announceBCInstallPlan(version int, manager string) {
	fmt.Printf("\n  +-- Install Plan ---------------------\n")
	fmt.Printf("  | Tool:     Beyond Compare %d\n", version)
	fmt.Printf("  | Manager:  %s\n", manager)
	fmt.Printf("  | Mode:     Silent /VERYSILENT\n")
	fmt.Printf("  +--------------------------------------\n\n")
}

func executeBCInstallWindows(version int) error {
	installerURL := resolveBCWindowsURL(version)
	destPath := filepath.Join(os.TempDir(), fmt.Sprintf("BCompare-%d-Setup.exe", version))
	defer os.Remove(destPath)

	fmt.Printf("Downloading Beyond Compare %d installer...\n", version)
	if err := downloadMultiTier(installerURL, destPath); err != nil {
		return apperror.WrapSimple(err, "download Beyond Compare installer")
	}

	return runBCInnoSetup(destPath, version)
}

func resolveBCWindowsURL(version int) string {
	if version == 4 {
		return bc4WinURL
	}

	return bc5WinURL
}

func runBCInnoSetup(installerPath string, version int) error {
	args := []string{"/VERYSILENT", "/NORESTART", "/SP-", "/SUPPRESSMSGBOXES"}
	cmd := exec.Command(installerPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Running silent Inno Setup for Beyond Compare %d...\n", version)
	if err := cmd.Run(); err != nil {
		return apperror.WrapSimple(err, "run Beyond Compare installer")
	}

	recordInstallation(resolveBCToolName(version), "inno-setup")
	fmt.Printf("  ✓ Beyond Compare %d installed successfully.\n", version)
	return nil
}

func isBCInstalledWindows(version int) bool {
	paths := getWindowsBCInstalledPaths(version)
	for _, p := range paths {
		if isFileExist(p) {
			return true
		}
	}

	return false
}

func getWindowsBCInstalledPaths(version int) []string {
	if version == 4 {
		return []string{
			`C:\Program Files\Beyond Compare 4\BCompare.exe`,
			`C:\Program Files (x86)\Beyond Compare 4\BCompare.exe`,
		}
	}

	return []string{
		`C:\Program Files\Beyond Compare 5\BCompare.exe`,
	}
}

func runInstallBeyondCompareLinux(version int, opts installOptions) error {
	if isBCInstalledLinux(version) {
		fmt.Printf("  ✓ Beyond Compare %d is already installed.\n", version)
		return nil
	}

	announceBCInstallPlan(version, "APT / DEB (Linux)")
	if opts.DryRun {
		return nil
	}

	return executeBCInstallLinux(version)
}

func isBCInstalledLinux(version int) bool {
	binName := fmt.Sprintf("bcompare%d", version)
	if isBinaryOnPath(binName) || isBinaryOnPath("bcompare") {
		return true
	}

	return isFileExist(fmt.Sprintf("/opt/beyondcompare%d/bcompare", version))
}

func isFileExist(path string) bool {
	if path == "" {
		return false
	}

	info, err := os.Stat(path)

	return err == nil && !info.IsDir()
}

func executeBCInstallLinux(version int) error {
	if hasAptPackageMgr() {
		return installBCDebLinux(version)
	}

	return installBCTarLinux(version)
}

func hasAptPackageMgr() bool {
	return isBinaryOnPath("apt-get")
}

func installBCDebLinux(version int) error {
	debURL := resolveBCDebURL(version)
	destPath := filepath.Join("/tmp", fmt.Sprintf("bcompare-%d.deb", version))
	defer os.Remove(destPath)

	fmt.Printf("Downloading Beyond Compare %d deb...\n", version)
	if err := downloadMultiTier(debURL, destPath); err != nil {
		return apperror.WrapSimple(err, "download Beyond Compare deb")
	}

	return runAptInstallDeb(destPath, version)
}

func resolveBCDebURL(version int) string {
	if version == 4 {
		return bc4DebURL
	}

	return bc5DebURL
}

func runAptInstallDeb(debPath string, version int) error {
	cmd := exec.Command("sudo", "apt-get", "install", "-y", debPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return apperror.WrapSimple(err, "apt install Beyond Compare deb")
	}

	recordInstallation(resolveBCToolName(version), constants.PkgMgrApt)
	fmt.Printf("  ✓ Beyond Compare %d installed successfully.\n", version)
	return nil
}

func installBCTarLinux(version int) error {
	tarURL := resolveBCTarURL(version)
	destPath := filepath.Join("/tmp", fmt.Sprintf("bcompare-%d.tar.gz", version))
	defer os.Remove(destPath)

	fmt.Printf("Downloading Beyond Compare %d tar.gz...\n", version)
	if err := downloadMultiTier(tarURL, destPath); err != nil {
		return apperror.WrapSimple(err, "download Beyond Compare tar")
	}

	return extractAndDeployBCTar(destPath, version)
}

func resolveBCTarURL(version int) string {
	if version == 4 {
		return bc4TarURL
	}

	return bc5TarURL
}

func extractAndDeployBCTar(tarPath string, version int) error {
	targetDir := fmt.Sprintf("/opt/beyondcompare%d", version)
	_ = os.MkdirAll(targetDir, 0755)

	cmd := exec.Command("tar", "-xzf", tarPath, "-C", targetDir, "--strip-components=1")
	if err := cmd.Run(); err != nil {
		return apperror.WrapSimple(err, "extract Beyond Compare tar")
	}

	_ = createBCSymlinks(targetDir, version)
	recordInstallation(resolveBCToolName(version), "tar")
	fmt.Printf("  ✓ Beyond Compare %d installed successfully to %s.\n", version, targetDir)
	return nil
}

func createBCSymlinks(targetDir string, version int) error {
	binPath := filepath.Join(targetDir, "bcompare")
	linkPath := fmt.Sprintf("/usr/local/bin/bcompare%d", version)
	_ = os.Remove(linkPath)

	return os.Symlink(binPath, linkPath)
}
