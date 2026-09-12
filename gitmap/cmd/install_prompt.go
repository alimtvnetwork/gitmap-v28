package cmd

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func isInstallApproved(opts installOptions, installName string) bool {
	if alreadyInstalled(installName) {
		return false
	}

	if opts.Check {
		fmt.Printf(constants.MsgInstallNotFound, installName)

		return false
	}

	return true
}

func runToolInstallation(opts installOptions, originalTool, installName string) error {
	opts.Tool = installName
	manager := resolvePackageManager(opts.Manager, opts.Tool)
	announceInstallPlan(opts.Version, manager)
	if confirmInstallIfNeeded(opts, installName, manager) {
		installTool(opts)
		postInstallSettingsSync(originalTool)
	}

	return nil
}

func executeGenericInstall(opts installOptions) {
	originalTool := opts.Tool
	installName := resolveNppInstallName(opts.Tool)
	if !isInstallApproved(opts, installName) {
		return
	}

	runToolInstallation(opts, originalTool, installName)
}

func alreadyInstalled(installName string) bool {
	fmt.Printf(constants.MsgInstallChecking, installName)
	existingVersion := detectInstalledVersion(installName)
	if existingVersion == "" {
		return false
	}

	fmt.Printf(constants.MsgInstallFound, installName, existingVersion)

	return true
}

func announceInstallPlan(version, manager string) {
	if version != "" {
		fmt.Printf(constants.MsgInstallVersion, version)
	} else {
		fmt.Print(constants.MsgInstallVersionLabel)
	}

	fmt.Printf(constants.MsgInstallManager, manager)
}

func confirmInstallIfNeeded(opts installOptions, installName, manager string) bool {
	if opts.Yes || opts.DryRun {
		return true
	}

	if confirmInstall(installName, opts.Version, manager) {
		return true
	}

	fmt.Print(constants.MsgInstallAborted)

	return false
}

func postInstallSettingsSync(originalTool string) {
	switch originalTool {
	case constants.ToolNpp:
		runNppSettings()
	case constants.ToolNppInstall:
		fmt.Print(constants.MsgInstallNppSkipSet)
	}
}

func confirmInstall(tool, version, manager string) bool {
	if version != "" {
		fmt.Printf(constants.MsgInstallPrompt, tool, version, manager)
	} else {
		fmt.Printf(constants.MsgInstallPromptNoVer, tool, manager)
	}

	var answer string
	fmt.Scanln(&answer)

	return answer == "y" || answer == "Y"
}
