package cmdinstall

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

const (
	tortoiseGitNupkgURL = "https://packages.chocolatey.org/TortoiseGit.2.4.0.2.nupkg"
	tortoiseGitVersion  = "2.4.0.2"
)

func runInstallTortoiseGit(opts installOptions) error {
	if !constants.IsWindowsPlatform() {
		fmt.Println("  ℹ TortoiseGit is a Windows Explorer shell extension and only available on Windows.")
		return nil
	}

	return runInstallTortoiseGitWindows(opts)
}

func runInstallTortoiseGitWindows(opts installOptions) error {
	if opts.Check {
		return checkTortoiseGitWindows()
	}

	if isTortoiseGitInstalledWindows() {
		fmt.Printf("  ✓ TortoiseGit %s is already installed.\n", tortoiseGitVersion)
		return nil
	}

	announceTortoiseGitPlan("Chocolatey / MSI (Windows)")

	if opts.DryRun {
		return nil
	}

	isChocoSuccess := tryChocoInstallTortoiseGit()

	if isChocoSuccess {
		recordInstallation(constants.ToolTortoiseGit, "choco")
		fmt.Printf("  ✓ TortoiseGit %s installed successfully via Chocolatey.\n", tortoiseGitVersion)
		return nil
	}

	return executeTortoiseGitMsiInstall()
}

func checkTortoiseGitWindows() error {
	if isTortoiseGitInstalledWindows() {
		fmt.Printf("  ✓ TortoiseGit %s is installed.\n", tortoiseGitVersion)
		return nil
	}

	fmt.Printf("  ✗ TortoiseGit %s is not installed.\n", tortoiseGitVersion)
	return nil
}

func announceTortoiseGitPlan(manager string) {
	fmt.Printf("\n  +-- Install Plan ---------------------\n")
	fmt.Printf("  | Tool:     TortoiseGit %s\n", tortoiseGitVersion)
	fmt.Printf("  | Manager:  %s\n", manager)
	fmt.Printf("  | Mode:     Silent /qn /norestart\n")
	fmt.Printf("  +--------------------------------------\n\n")
}

func isTortoiseGitInstalledWindows() bool {
	paths := getWindowsTortoiseGitPaths()

	for _, p := range paths {

		if isFileExist(p) {
			return true
		}
	}

	return isBinaryOnPath("TortoiseGitProc.exe")
}

func getWindowsTortoiseGitPaths() []string {
	return []string{
		`C:\Program Files\TortoiseGit\bin\TortoiseGitProc.exe`,
		`C:\Program Files (x86)\TortoiseGit\bin\TortoiseGitProc.exe`,
	}
}

func tryChocoInstallTortoiseGit() bool {
	if !isBinaryOnPath("choco") {
		return false
	}

	fmt.Printf("Installing TortoiseGit %s via Chocolatey...\n", tortoiseGitVersion)
	cmd := exec.Command("choco", "install", "tortoisegit", "--version", tortoiseGitVersion, "-y")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	return err == nil
}

func executeTortoiseGitMsiInstall() error {
	tempZip := filepath.Join(os.TempDir(), "TortoiseGit.2.4.0.2.nupkg.zip")
	extractDir := filepath.Join(os.TempDir(), "TortoiseGitExtract")
	defer func() {
		_ = os.Remove(tempZip)
		_ = os.RemoveAll(extractDir)
	}()

	fmt.Printf("Downloading TortoiseGit %s package...\n", tortoiseGitVersion)

	if err := downloadMultiTier(tortoiseGitNupkgURL, tempZip); err != nil {
		return apperror.WrapSimple(err, "download TortoiseGit package")
	}

	msiPath, err := extractTortoiseGitMsi(tempZip, extractDir)

	if err != nil {
		return err
	}

	return runMsiSilentInstall(msiPath)
}

func extractTortoiseGitMsi(zipPath, targetDir string) (string, error) {
	reader, err := zip.OpenReader(zipPath)

	if err != nil {
		return "", apperror.WrapSimple(err, "open TortoiseGit nupkg zip")
	}

	defer reader.Close()

	is64Bit := runtime.GOARCH == "amd64"
	targetSuffix := "_x32.msi"

	if is64Bit {
		targetSuffix = "_x64.msi"
	}

	return findAndExtractMsiEntry(reader, targetDir, targetSuffix)
}

func findAndExtractMsiEntry(reader *zip.ReadCloser, targetDir, targetSuffix string) (string, error) {
	for _, f := range reader.File {

		if strings.HasSuffix(f.Name, targetSuffix) {
			return extractSingleMsiFile(f, targetDir)
		}
	}

	return "", apperror.NewSimple("target MSI not found in TortoiseGit package", "E_MSI_NOT_FOUND")
}

func extractSingleMsiFile(f *zip.File, targetDir string) (string, error) {
	rc, err := f.Open()

	if err != nil {
		return "", apperror.WrapSimple(err, "open MSI entry in zip")
	}

	defer rc.Close()

	_ = os.MkdirAll(targetDir, 0755)
	destPath := filepath.Join(targetDir, filepath.Base(f.Name))
	outFile, err := os.Create(destPath)

	if err != nil {
		return "", apperror.WrapSimple(err, "create MSI file on disk")
	}

	defer outFile.Close()

	if _, err := io.Copy(outFile, rc); err != nil {
		return "", apperror.WrapSimple(err, "copy MSI content")
	}

	return destPath, nil
}

func runMsiSilentInstall(msiPath string) error {
	args := []string{"/i", msiPath, "/quiet", "/qn", "/norestart", "REBOOT=ReallySuppress"}
	cmd := exec.Command("msiexec.exe", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Running silent MSI installation for TortoiseGit %s...\n", tortoiseGitVersion)

	if err := cmd.Run(); err != nil {
		return apperror.WrapSimple(err, "run TortoiseGit MSI installer")
	}

	recordInstallation(constants.ToolTortoiseGit, "msi")
	fmt.Printf("  ✓ TortoiseGit %s installed successfully via MSI.\n", tortoiseGitVersion)

	return nil
}
