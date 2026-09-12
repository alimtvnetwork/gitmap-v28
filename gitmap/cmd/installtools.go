package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// installTool dispatches to the platform-specific installer.
func installTool(opts installOptions) {
	manager := resolvePackageManager(opts.Manager, opts.Tool)
	installCmd := buildInstallCommand(manager, opts.Tool, opts.Version)
	versionLabel := resolveDisplayVersion(opts.Version)

	printInstallPlan(opts.Tool, versionLabel, manager, installCmd)

	isDryRun := handleDryRunInstall(opts.DryRun, manager, installCmd)

	if isDryRun {
		return
	}

	executeInstallSteps(installCmd, opts, manager, versionLabel)
}

func printInstallPlan(tool, version, manager string, installCmd []string) {
	fmt.Printf("\n  +-- Install Plan ---------------------\n")
	fmt.Printf("  | Tool:    %s\n", tool)
	fmt.Printf("  | Version: %s\n", version)
	fmt.Printf("  | Manager: %s\n", manager)
	fmt.Printf("  | Command: %s\n", strings.Join(installCmd, " "))
	fmt.Printf("  +--------------------------------------\n\n")
}

func handleDryRunInstall(dryRun bool, manager string, installCmd []string) bool {
	if !dryRun {
		return false
	}

	if manager == constants.PkgMgrApt {
		fmt.Printf(constants.MsgInstallDryCmd, "sudo apt-get update")
	}

	fmt.Printf(constants.MsgInstallDryCmd, strings.Join(installCmd, " "))

	return true
}

func resolveTotalInstallSteps(manager string) int {
	if manager == constants.PkgMgrApt {
		return 4
	}

	return 3
}

func executeInstallSteps(installCmd []string, opts installOptions, manager, versionLabel string) {
	totalSteps := resolveTotalInstallSteps(manager)
	step := 1

	if manager == constants.PkgMgrApt {
		fmt.Printf("  [%d/%d] Updating package index...\n", step, totalSteps)
		runAptUpdate(opts.Verbose)
		step++
	}

	runStepInstall(installCmd, opts, manager, versionLabel, step, totalSteps)
}

func runStepInstall(
	installCmd []string, opts installOptions,
	manager, versionLabel string, step, totalSteps int,
) error {
	fmt.Printf("  [%d/%d] Installing %s v%s via %s...\n", step, totalSteps, opts.Tool, versionLabel, manager)
	runInstallCommand(installCmd, opts)
	step++

	fmt.Printf("  [%d/%d] Verifying installation...\n", step, totalSteps)
	verifyInstallation(opts.Tool)
	step++

	fmt.Printf("  [%d/%d] Recording installation...\n", step, totalSteps)
	recordInstallation(opts.Tool, manager)

	return nil
}

func buildInstallCommand(manager, tool, version string) []string {
	pkg := resolvePackageName(manager, tool)
	cmd, isUnix := buildUnixInstallCommand(manager, tool, pkg, version)

	if isUnix {
		return cmd
	}

	return buildWindowsInstallCommand(manager, pkg, version)
}

func buildUnixInstallCommand(manager, tool, pkg, version string) ([]string, bool) {
	switch manager {
	case constants.PkgMgrApt:

		return buildAptCommand(pkg, version), true
	case constants.PkgMgrBrew:

		return buildBrewCommand(tool, pkg), true
	case constants.PkgMgrSnap:

		return buildSnapCommand(pkg), true
	}

	return nil, false
}

func buildWindowsInstallCommand(manager, pkg, version string) []string {
	if manager == constants.PkgMgrWinget {
		return buildWingetCommand(pkg, version)
	}

	return buildChocoCommand(pkg, version)
}

func buildChocoCommand(pkg, version string) []string {
	baseArgs := []string{"choco", "install", pkg, "-y", "--no-progress"}

	return appendVersionFlag(baseArgs, version)
}

func buildWingetCommand(pkg, version string) []string {
	baseArgs := []string{"winget", "install", pkg, "--accept-package-agreements", "--accept-source-agreements", "--silent"}

	return appendVersionFlag(baseArgs, version)
}

func appendVersionFlag(baseArgs []string, version string) []string {
	if version != "" {
		return append(baseArgs, "--version", version)
	}

	return baseArgs
}

func buildAptCommand(pkg, version string) []string {
	target := resolveAptTarget(pkg, version)

	return []string{"sudo", "apt", "install", "-y", target}
}

func resolveAptTarget(pkg, version string) string {
	if version != "" {
		return pkg + "=" + version
	}

	return pkg
}

func buildBrewCommand(tool, pkg string) []string {
	if isBrewCaskTool(tool) {
		return []string{"brew", "install", "--cask", pkg}
	}

	return []string{"brew", "install", pkg}
}

func buildSnapCommand(pkg string) []string {
	if pkg == constants.CmdCode {
		return []string{"sudo", "snap", "install", pkg, "--classic"}
	}

	return []string{"sudo", "snap", "install", pkg}
}

func isBrewCaskTool(tool string) bool {
	switch tool {
	case constants.ToolVSCode, constants.ToolGitHubDesktop, constants.ToolPowerShell,
		constants.ToolDbeaver, constants.ToolOBS, constants.ToolVLC,
		constants.ToolFlameshot, constants.ToolDocker, constants.ToolFlutter,
		constants.ToolQBittorrent, constants.ToolUTorrent:

		return true
	default:

		return false
	}
}

func runAptUpdate(verbose bool) error {
	fmt.Print(constants.MsgInstallAptUpdate)
	cmd := exec.Command("sudo", "apt-get", "update")

	if verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrInstallAptUpdateFailed, err)

		return nil
	}

	fmt.Print(constants.MsgInstallAptUpdateDone)

	return nil
}

func recordInstallation(tool, manager string) {
	version := detectInstalledVersion(tool)
	splitDB, err := store.OpenInstallationSplitDB()

	if err != nil {
		return
	}

	defer splitDB.Close()

	if err := splitDB.SaveInstalledTool(tool, version, manager); err == nil && version != "" {
		fmt.Printf(constants.MsgInstallRecorded, tool, version)
	}

	_ = splitDB.RecordLog(store.InstallationLogRecord{
		Tool:           tool,
		Action:         "install",
		Version:        version,
		PackageManager: manager,
		IsSuccess:      true,
	})
}
