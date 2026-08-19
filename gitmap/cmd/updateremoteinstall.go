package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runUpdateRemoteInstall is the v5.52.0+ remote-installer flow:
//
//  1. Resolve the latest gitmap-vN repo slug natively in Go via the
//     20-parallel sibling probe (spec/01-app/111-update-remote-probe.md).
//  2. Download THAT repo's install.{ps1,sh} from raw.githubusercontent.com.
//  3. Exec it with inherited stdio.
//
// Returns true if the install completed (exit 0) so the caller can
// short-circuit any legacy source-rebuild fallback.
//
// Flags: --probe-only prints the resolution and exits; --no-probe skips
// the probe and installs from the current repo only.
func runUpdateRemoteInstall() bool {
	slug, source, err := resolveTargetSlug()
	if err != nil {
		return false
	}
	if hasFlag(constants.FlagProbeOnly) {
		fmt.Printf(constants.MsgUpdateProbeOnly, slug, source)
		return true
	}

	currentVersion := constants.Version
	targetVersion := fetchRemoteTargetVersion(slug)

	url := installerURLFor(slug)
	fmt.Printf(constants.MsgUpdateRemoteFetch, url)

	scriptPath, err := downloadRemoteInstaller(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrUpdateRemoteDownload, err)
		return false
	}
	defer os.Remove(scriptPath)

	fmt.Printf(constants.MsgUpdateVersionCompare, currentVersion, targetVersion)
	fmt.Printf(constants.MsgUpdateRemoteRun, scriptPath)
	errRun := runRemoteInstaller(scriptPath)
	var exitErr *exec.ExitError
	if errRun != nil && errors.As(errRun, &exitErr) == true {
		os.Exit(exitErr.ExitCode())
	}
	if errRun != nil {
		fmt.Fprintf(os.Stderr, constants.ErrUpdateRemoteRun, errRun)
		return false
	}

	fmt.Printf(constants.MsgUpdateSummaryDetail, currentVersion, targetVersion, url)
	return true
}

// fetchRemoteTargetVersion fetches version.json from the remote repo.
func fetchRemoteTargetVersion(slug string) string {
	url := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/main/version.json", constants.UpdateRepoOwner, slug)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url) //nolint:gosec
	if err != nil {
		return "unknown"
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "unknown"
	}

	var data struct {
		Version string `json:"version"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "unknown"
	}
	if data.Version == "" {
		return "unknown"
	}
	return data.Version
}

// resolveTargetSlug returns the repo slug to install from, honoring
// --no-probe (skip probe, use current slug).
func resolveTargetSlug() (string, string, error) {
	if hasFlag(constants.FlagNoProbe) {
		fmt.Printf(constants.MsgUpdateProbeSkipped, constants.UpdateCurrentRepoSlug)
		return constants.UpdateCurrentRepoSlug, constants.UpdateProbeSourceMain, nil
	}
	return resolveLatestRepoSlug(newProbeClient())
}

// installerURLFor builds the raw.githubusercontent installer URL for slug.
func installerURLFor(slug string) string {
	name := constants.UpdateInstallerNameBash
	if runtime.GOOS == "windows" {
		name = constants.UpdateInstallerNamePwsh
	}
	return fmt.Sprintf(constants.UpdateRawInstallerTmpl,
		constants.UpdateRepoOwner, slug, name)
}

// downloadRemoteInstaller fetches url into a platform-appropriate temp
// file (.ps1 on Windows, .sh elsewhere) and returns the path.
func downloadRemoteInstaller(url string) (string, error) {
	resp, err := http.Get(url) //nolint:gosec // URL built from constants
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d from %s", resp.StatusCode, url)
	}

	ext := ".sh"
	if runtime.GOOS == "windows" {
		ext = ".ps1"
	}
	tmp, err := os.CreateTemp("", "gitmap-update-*"+ext)
	if err != nil {
		return "", err
	}
	if runtime.GOOS == "windows" {
		_, _ = tmp.Write([]byte{0xEF, 0xBB, 0xBF})
	}
	if _, err := io.Copy(tmp, resp.Body); err != nil {
		tmp.Close()
		os.Remove(tmp.Name())
		return "", err
	}
	tmp.Close()

	if runtime.GOOS != "windows" {
		_ = os.Chmod(tmp.Name(), 0o755)
	}
	return tmp.Name(), nil
}

// runRemoteInstaller exec's the downloaded script with the right shell.
func runRemoteInstaller(scriptPath string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("powershell",
			"-ExecutionPolicy", "Bypass",
			"-NoProfile", "-NoLogo",
			"-File", scriptPath)
	}
	if runtime.GOOS != "windows" {
		cmd = exec.Command(getUnixShell(), scriptPath)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Dir = filepath.Dir(scriptPath)
	return cmd.Run()
}

func getUnixShell() string {
	shell := "bash"
	_, errLook := exec.LookPath(shell)
	if errLook != nil {
		shell = "sh"
	}
	return shell
}
