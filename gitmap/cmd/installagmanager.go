package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

func runInstallAgManager() error {

	return runInstallAgManagerWithOpts(installOptions{})
}

func runInstallAgManagerWithOpts(opts installOptions) error {
	fmt.Println("Fetching release for Antigravity-Manager...")
	assetURL, ver, err := resolveAgManagerAssetURL(opts.Version)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching release: %v\n", err)

		return nil
	}
	if opts.DryRun {
		fmt.Printf("  [dry-run] Would download %s (version: %s) and execute installer\n", assetURL, ver)

		return nil
	}
	fmt.Printf("Downloading %s (version: %s)...\n", assetURL, ver)
	performAgManagerDownloadAndInstall(assetURL, ver)

	return nil
}

func performAgManagerDownloadAndInstall(assetURL, ver string) {
	tmpPath, err := downloadAgManagerFile(assetURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error downloading file: %v\n", err)

		return
	}
	fmt.Printf("Installing %s...\n", filepath.Base(tmpPath))
	if err := executeAgManagerInstaller(tmpPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error installing: %v\n", err)

		return
	}
	recordAgManagerInstalled(ver)
	fmt.Println(constants.ColorGreen + "✓" + constants.ColorReset + " Antigravity Manager installed successfully.")
}

func recordAgManagerInstalled(ver string) {
	splitDB, err := store.OpenInstallationSplitDB()
	if err != nil {

		return
	}
	defer splitDB.Close()
	_ = splitDB.SaveInstalledTool("ag-manager", ver, "github-release")
}

func downloadAgManagerFile(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {

		return "", err
	}
	defer resp.Body.Close()
	parts := strings.Split(url, "/")
	name := parts[len(parts)-1]
	tmpPath := filepath.Join(os.TempDir(), name)
	out, err := os.Create(tmpPath)
	if err != nil {

		return "", err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)

	return tmpPath, err
}

func executeAgManagerInstaller(path string) error {
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
	os.Chmod(path, 0755)

	return exec.Command(path)
}
