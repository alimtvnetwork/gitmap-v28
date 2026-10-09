package cmdinstall

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func runInstallAgManagerWithOpts(opts installOptions) error {
	return executeAgManagerOneLiner(opts, false)
}

func runUpdateAgManagerWithOpts(opts installOptions) error {
	return executeAgManagerOneLiner(opts, true)
}

func executeAgManagerOneLiner(opts installOptions, isUpdate bool) error {
	if isAgManagerInstallSkipped(opts, isUpdate) {
		return nil
	}
	if opts.DryRun {
		fmt.Printf("  [dry-run] Would run: %s (version: %s)\n", resolveAgManagerCommandForOS(), opts.Version)
		return nil
	}
	return runAgManagerScriptWithFeedback(isUpdate, opts.Version)
}

func isAgManagerInstallSkipped(opts installOptions, isUpdate bool) bool {
	ver, isFound := isAgManagerInstalled()
	if isFound && !opts.Force && !isUpdate {
		fmt.Printf("  ✓ Antigravity Manager is already installed (%s). Use 'gitmap agm update' or --force to reinstall.\n", ver)
		return true
	}
	return false
}

func runAgManagerScriptWithFeedback(isUpdate bool, version string) error {
	if version == "" {
		version = resolveLatestAgManagerReleaseVersion()
	}
	renderer := newAgmUpdateRenderer(resolveAgManagerActionName(isUpdate), resolveAgManagerPlatformName(), version)
	if err := dispatchAgManagerScriptWithVersion(renderer, version); err != nil {
		renderer.renderFailure(err)
		verifyAgManagerOnFailure()
		return apperror.WrapSimple(err, "Antigravity Manager execution failed")
	}
	renderer.renderSuccess()
	recordAgManagerInstalled(resolveInstalledVerName(version))
	return nil
}

func verifyAgManagerOnFailure() {
	if _, isFound := isAgManagerInstalled(); !isFound {
		reportVerificationFailure(constants.ToolAgManager, "agm-alim")
	}
}

func formatAgManagerVerLabel(version string) string {
	if version != "" {
		return " (" + version + ")"
	}
	return ""
}

func resolveInstalledVerName(version string) string {
	if version != "" {
		return version
	}
	return "latest"
}

func resolveAgManagerActionName(isUpdate bool) string {
	if isUpdate {
		return "Updating"
	}
	return "Installing"
}

func resolveAgManagerCommandForOS() string {
	if runtime.GOOS == "windows" {
		return constants.AgManagerWindowsInstallCmd
	}
	return constants.AgManagerUnixInstallCmd
}

func resolveAgManagerPlatformName() string {
	if runtime.GOOS == "windows" {
		return "PowerShell"
	}
	return "curl | bash"
}

func resolveLatestAgManagerReleaseVersion() string {
	client := &http.Client{Timeout: 3 * time.Second}
	if tag := fetchLatestTagFromGitHubAPI(client); tag != "" {
		return tag
	}
	return fetchLatestTagFromRedirect(client)
}

func fetchLatestTagFromGitHubAPI(client *http.Client) string {
	resp, err := client.Get("https://api.github.com/repos/alimtvnetwork/Antigravity-Manager/releases/latest")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	var payload struct {
		TagName string `json:"tag_name"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&payload)
	return strings.TrimPrefix(strings.TrimSpace(payload.TagName), "v")
}

func fetchLatestTagFromRedirect(client *http.Client) string {
	resp, err := client.Get("https://github.com/alimtvnetwork/Antigravity-Manager/releases/latest")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	parts := strings.Split(strings.TrimSuffix(resp.Request.URL.Path, "/"), "/")
	return resolveCleanReleaseTag(parts[len(parts)-1])
}

func resolveCleanReleaseTag(tag string) string {
	if tag == "" || tag == "latest" || tag == "releases" {
		return ""
	}
	return strings.TrimPrefix(tag, "v")
}

func dispatchAgManagerScriptWithVersion(renderer *agmUpdateRenderer, version string) error {
	if runtime.GOOS == "windows" {
		return dispatchAgManagerWindowsWithVersion(renderer, version)
	}
	return dispatchAgManagerUnixWithVersion(renderer, version)
}

func buildWindowsAgManagerInstallCmd(version string) string {
	scriptURL := fmt.Sprintf("https://raw.githubusercontent.com/alimtvnetwork/Antigravity-Manager/main/install.ps1?cb=%d", time.Now().Unix())
	if version != "" {
		return fmt.Sprintf("& ([scriptblock]::Create((irm '%s'))) -Version '%s' -Update -NoLaunch", scriptURL, strings.TrimPrefix(version, "v"))
	}
	return fmt.Sprintf("& ([scriptblock]::Create((irm '%s'))) -Update -NoLaunch", scriptURL)
}

func dispatchAgManagerWindowsWithVersion(renderer *agmUpdateRenderer, version string) error {
	pwsh := resolvePowerShellBinary()
	if pwsh == "" {
		return apperror.NewSimple("PowerShell not found on PATH. Run manually:\n  "+constants.AgManagerWindowsInstallCmd, "E9000")
	}
	if version == "" {
		version = resolveLatestAgManagerReleaseVersion()
	}
	installCmd := buildWindowsAgManagerInstallCmd(version)
	env := append(os.Environ(), "AGM_VERSION="+strings.TrimPrefix(version, "v"))
	return runAgmInstallerStreamed(pwsh, []string{"-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", installCmd}, env, renderer)
}

func dispatchAgManagerUnixWithVersion(renderer *agmUpdateRenderer, version string) error {
	installCmd := resolveAgmUnixInstallCmd(version)
	return runAgmInstallerStreamed("bash", []string{"-c", installCmd}, nil, renderer)
}

func recordAgManagerInstalled(ver string) {
	splitDB, err := store.OpenInstallationSplitDB()
	if err != nil {
		return
	}
	defer splitDB.Close()
	_ = splitDB.SaveInstalledTool("ag-manager", ver, "github-release")
}
