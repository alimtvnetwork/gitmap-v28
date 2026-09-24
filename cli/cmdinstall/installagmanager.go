package cmdinstall

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

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
	action := resolveAgManagerActionName(isUpdate)
	verLabel := ""
	if version != "" {
		verLabel = " (" + version + ")"
	}
	fmt.Printf("%s Antigravity Manager%s via %s...\n", action, verLabel, resolveAgManagerPlatformName())
	if err := dispatchAgManagerScriptWithVersion(version); err != nil {
		reportVerificationFailure(constants.ToolAgManager, "ag-manager")
		return apperror.WrapSimple(err, "Antigravity Manager execution failed")
	}
	installedVer := "latest"
	if version != "" {
		installedVer = version
	}
	recordAgManagerInstalled(installedVer)
	fmt.Printf("%s✓%s Antigravity Manager%s completed successfully.\n", constants.ColorGreen, constants.ColorReset, verLabel)
	return nil
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

func dispatchAgManagerScriptWithVersion(version string) error {
	if runtime.GOOS == "windows" {
		return dispatchAgManagerWindowsWithVersion(version)
	}
	return dispatchAgManagerUnixWithVersion(version)
}

func dispatchAgManagerWindowsWithVersion(version string) error {
	pwsh := resolvePowerShellBinary()
	if pwsh == "" {
		return apperror.NewSimple("PowerShell not found on PATH. Run manually:\n  "+constants.AgManagerWindowsInstallCmd, "E9000")
	}
	installCmd := constants.AgManagerWindowsInstallCmd
	if version != "" {
		clean := strings.TrimPrefix(version, "v")
		installCmd = fmt.Sprintf(`& ([scriptblock]::Create((irm https://raw.githubusercontent.com/alimtvnetwork/Antigravity-Manager/main/install.ps1))) -Version '%s'`, clean)
	}
	cmd := exec.Command(pwsh, "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", installCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func dispatchAgManagerUnixWithVersion(version string) error {
	installCmd := constants.AgManagerUnixInstallCmd
	if version != "" {
		clean := strings.TrimPrefix(version, "v")
		installCmd = fmt.Sprintf(`curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/Antigravity-Manager/main/install.sh | bash -s -- --version '%s'`, clean)
	}
	cmd := exec.Command("bash", "-c", installCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func recordAgManagerInstalled(ver string) {
	splitDB, err := store.OpenInstallationSplitDB()
	if err != nil {
		return
	}
	defer splitDB.Close()
	_ = splitDB.SaveInstalledTool("ag-manager", ver, "github-release")
}
