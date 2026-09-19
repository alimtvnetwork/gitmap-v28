package cmdinstall

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

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
		fmt.Printf("  [dry-run] Would run: %s\n", resolveAgManagerCommandForOS())
		return nil
	}
	return runAgManagerScriptWithFeedback(isUpdate)
}

func isAgManagerInstallSkipped(opts installOptions, isUpdate bool) bool {
	ver, isFound := isAgManagerInstalled()
	if isFound && !opts.Force && !isUpdate {
		fmt.Printf("  ✓ Antigravity Manager is already installed (%s). Use 'gitmap agm update' or --force to reinstall.\n", ver)
		return true
	}
	return false
}

func runAgManagerScriptWithFeedback(isUpdate bool) error {
	action := resolveAgManagerActionName(isUpdate)
	fmt.Printf("%s Antigravity Manager via %s...\n", action, resolveAgManagerPlatformName())
	if err := dispatchAgManagerScript(); err != nil {
		reportVerificationFailure(constants.ToolAgManager, "ag-manager")
		return apperror.WrapSimple(err, "Antigravity Manager execution failed")
	}
	recordAgManagerInstalled("latest")
	fmt.Printf("%s✓%s Antigravity Manager completed successfully.\n", constants.ColorGreen, constants.ColorReset)
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

func dispatchAgManagerScript() error {
	if runtime.GOOS == "windows" {
		return dispatchAgManagerWindows()
	}
	return dispatchAgManagerUnix()
}

func dispatchAgManagerWindows() error {
	pwsh := resolvePowerShellBinary()
	if pwsh == "" {
		return apperror.NewSimple("PowerShell not found on PATH. Run manually:\n  "+constants.AgManagerWindowsInstallCmd, "E9000")
	}
	cmd := exec.Command(pwsh, "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", constants.AgManagerWindowsInstallCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func dispatchAgManagerUnix() error {
	cmd := exec.Command("bash", "-c", constants.AgManagerUnixInstallCmd)
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
