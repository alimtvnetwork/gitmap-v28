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
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runUpdateRemoteInstall is the v5.52.0+ remote-installer flow:
//
//  1. Resolve the latest gitmap-vN repo slug natively in Go via the
//     20-parallel sibling probe (spec/01-app/111-update-remote-probe.md).
//  2. Download THAT repo's install.{ps1,sh} from raw.githubusercontent.com.
//  3. Exec it with inherited stdio.
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
	return executeRemoteUpdateWorkflow(url, currentVersion, targetVersion)
}

func executeRemoteUpdateWorkflow(url, currentVersion, targetVersion string) bool {
	scriptPath, err := downloadRemoteInstaller(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrUpdateRemoteDownload, err)
		return false
	}
	defer os.Remove(scriptPath)
	fmt.Printf(constants.MsgUpdateVersionCompare, currentVersion, targetVersion)
	fmt.Printf(constants.MsgUpdateRemoteRun, scriptPath)
	errRun := runRemoteInstaller(scriptPath)
	if errRun != nil {
		handleRemoteInstallerError(errRun)
		return false
	}
	fmt.Printf(constants.MsgUpdateSummaryDetail, currentVersion, targetVersion, url)
	return true
}

func handleRemoteInstallerError(errRun error) {
	var exitErr *exec.ExitError
	if errors.As(errRun, &exitErr) {
		appErr := apperror.NewWithDetails(
			"cmd.updateremoteinstall.run",
			"E1153",
			fmt.Sprintf("remote installer exited with code %d", exitErr.ExitCode()),
			"cmd.updateremoteinstall",
			apperror.ErrorTypeExecution,
			apperror.SeverityError,
			map[string]any{"exitCode": exitErr.ExitCode()},
		)
		cliexit.HandleError(appErr, exitErr.ExitCode())
		return
	}
	fmt.Fprintf(os.Stderr, constants.ErrUpdateRemoteRun, errRun)
}

// fetchRemoteTargetVersion fetches the latest version from GitHub releases or version.json.
func fetchRemoteTargetVersion(slug string) string {
	ghVer := fetchGitHubLatestReleaseVersion(slug)
	if len(ghVer) > 0 && ghVer != "unknown" {
		return ghVer
	}
	return fetchVersionJSON(slug)
}

func fetchGitHubLatestReleaseVersion(slug string) string {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", constants.UpdateRepoOwner, slug)
	client := &http.Client{Timeout: 6 * time.Second}
	req, reqErr := http.NewRequest(http.MethodGet, url, nil)
	if reqErr != nil {
		return ""
	}
	req.Header.Set("User-Agent", "gitmap-updater")
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		closeResponse(resp)
		return ""
	}
	defer resp.Body.Close()
	return decodeReleaseTagName(resp.Body)
}

func closeResponse(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		_ = resp.Body.Close()
	}
}

func decodeReleaseTagName(body io.Reader) string {
	var releaseData struct {
		TagName string `json:"tag_name"`
		Name    string `json:"name"`
	}
	if err := json.NewDecoder(body).Decode(&releaseData); err != nil {
		return ""
	}
	tag := strings.TrimPrefix(releaseData.TagName, "v")
	if len(tag) > 0 {
		return tag
	}
	return strings.TrimPrefix(releaseData.Name, "v")
}

func fetchVersionJSON(slug string) string {
	url := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/main/version.json", constants.UpdateRepoOwner, slug)
	client := &http.Client{Timeout: 6 * time.Second}
	resp, err := client.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		closeResponse(resp)
		return "unknown"
	}
	defer resp.Body.Close()
	return decodeVersionFromMap(resp.Body)
}

func decodeVersionFromMap(body io.Reader) string {
	var rawMap map[string]interface{}
	if err := json.NewDecoder(body).Decode(&rawMap); err != nil {
		return "unknown"
	}
	if v, ok := rawMap["Version"].(string); ok && len(v) > 0 {
		return v
	}
	if v, ok := rawMap["version"].(string); ok && len(v) > 0 {
		return v
	}
	return "unknown"
}

// resolveTargetSlug returns the repo slug to install from.
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

// downloadRemoteInstaller fetches url into a platform-appropriate temp file.
func downloadRemoteInstaller(url string) (string, error) {
	resp, err := http.Get(url) //nolint:gosec
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d from %s", resp.StatusCode, url)
	}
	return writeInstallerTempFile(resp.Body)
}

func writeInstallerTempFile(body io.Reader) (string, error) {
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
	if _, copyErr := io.Copy(tmp, body); copyErr != nil {
		tmp.Close()
		os.Remove(tmp.Name())
		return "", copyErr
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
	} else {
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
