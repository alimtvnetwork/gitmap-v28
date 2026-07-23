package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// RunWorker is the hidden command that performs the actual update.
// It is invoked by the handoff copy: gitmap-updater-tmp update-worker <version>
func RunWorker() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: gitmap-updater update-worker <version>")
		os.Exit(1)
	}

	version := os.Args[2]

	scriptPath, err := downloadInstaller(version)
	if err != nil {
		fmt.Fprintf(os.Stderr, ErrDownload, err)
		os.Exit(1)
	}
	defer os.Remove(scriptPath)

	fmt.Print(MsgRunningInstall)
	if err := runInstaller(scriptPath); err != nil {
		fmt.Fprintf(os.Stderr, ErrRunInstaller, err)
		os.Exit(1)
	}

	verifyUpdate(version)
	cleanupSelf()
	fmt.Print(MsgDone)
}

// downloadInstaller downloads install.ps1 from the release assets.
func downloadInstaller(version string) (string, error) {
	url := fmt.Sprintf(ReleaseInstallURL, version)
	fmt.Printf(MsgDownloading, version)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d from %s", resp.StatusCode, url)
	}

	tmpFile, err := os.CreateTemp(os.TempDir(), "gitmap-install-*.ps1")
	if err != nil {
		return "", err
	}

	// Write UTF-8 BOM for PowerShell compatibility
	tmpFile.Write([]byte{0xEF, 0xBB, 0xBF})

	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())

		return "", err
	}

	tmpFile.Close()

	return tmpFile.Name(), nil
}

// runInstaller executes the downloaded install.ps1 via PowerShell.
func runInstaller(scriptPath string) error {
	cmd := exec.Command(PSBin, PSExecPolicy, PSBypass,
		PSNoProfile, PSNoLogo, PSFile, scriptPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

// verifyUpdate checks that gitmap now reports the expected version.
func verifyUpdate(expectedVersion string) {
	actual, err := getInstalledVersion()
	if err != nil {
		return
	}

	if normalizeVersion(actual) != normalizeVersion(expectedVersion) {
		fmt.Fprintf(os.Stderr, MsgVerifyFail, expectedVersion, actual)
	}
}

// cleanupSelf removes leftover updater temp copies from the directory.
func cleanupSelf() {
	selfPath, err := os.Executable()
	if err != nil {
		return
	}

	dir := filepath.Dir(selfPath)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	for _, e := range entries {
		if strings.HasPrefix(e.Name(), "gitmap-updater-tmp-") {
			os.Remove(filepath.Join(dir, e.Name()))
		}
	}
}
