package cmd

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func runInstallGitHubDesktopLinux(opts installOptions) error {
	isInstalled := alreadyInstalled(opts.Tool)

	if isInstalled {

		return nil
	}

	manager := resolvePackageManager(opts.Manager, opts.Tool)

	if manager != constants.PkgMgrApt {
		executeGenericInstall(opts)

		return nil
	}

	return executeGitHubDesktopAptPipeline(opts)
}

func executeGitHubDesktopAptPipeline(opts installOptions) error {
	printGhdInstallPlan(opts.Tool)
	isDryRun := handleDryRunInstall(opts.DryRun, "apt", []string{"sudo", "apt", "install", "-y", "github-desktop"})

	if isDryRun {

		return nil
	}

	return executeGhdAptPhases(opts)
}

func printGhdInstallPlan(tool string) {
	fmt.Printf("\n  +-- Install Plan ---------------------\n")
	fmt.Printf("  | Tool:    %s\n", tool)
	fmt.Printf("  | Version: shiftkey/desktop\n")
	fmt.Printf("  | Manager: apt (shiftkey repo)\n")
	fmt.Printf("  | Command: add shiftkey apt repository & sudo apt install github-desktop\n")
	fmt.Printf("  +--------------------------------------\n\n")
}

func executeGhdAptPhases(opts installOptions) error {
	err1 := runGhdGpgKeyStep(opts)

	if err1 != nil {

		return err1
	}

	err2 := runGhdRepoStep(opts)

	if err2 != nil {

		return err2
	}

	runAptUpdate(opts.Verbose)
	runInstallCommand([]string{"sudo", "apt", "install", "-y", "github-desktop"}, opts)
	recordInstallation(opts.Tool, "apt")
	fmt.Printf("  !Done\n")

	return nil
}

func runGhdGpgKeyStep(opts installOptions) error {
	fmt.Printf("  [1/3] Adding shiftkey/desktop APT repository...\n")
	cmd := []string{"sh", "-c", "wget -qO - https://mirror.mwt.me/ghd/gpgkey | sudo tee /etc/apt/keyrings/shiftkey-packages.asc > /dev/null"}

	return runPhaseWithAudit(opts.Tool, "apt", "shiftkey/desktop", cmd, "ghd-gpg-key", opts.Verbose)
}

func runGhdRepoStep(opts installOptions) error {
	cmd := []string{"sh", "-c", "sudo sh -c 'echo \"deb [arch=amd64 signed-by=/etc/apt/keyrings/shiftkey-packages.asc] https://mirror.mwt.me/ghd/deb/ any main\" > /etc/apt/sources.list.d/shiftkey-packages.list'"}

	return runPhaseWithAudit(opts.Tool, "apt", "shiftkey/desktop", cmd, "ghd-apt-repo", opts.Verbose)
}

func runInstallVSCodeLinux(opts installOptions) error {
	isInstalled := alreadyInstalled(opts.Tool)

	if isInstalled {

		return nil
	}

	manager := resolvePackageManager(opts.Manager, opts.Tool)

	if manager != constants.PkgMgrApt {
		executeGenericInstall(opts)

		return nil
	}

	return executeVSCodeAptPipeline(opts)
}

func executeVSCodeAptPipeline(opts installOptions) error {
	printVSCodeInstallPlan(opts.Tool)
	isDryRun := handleDryRunInstall(opts.DryRun, "apt", []string{"sudo", "apt", "install", "-y", "code"})

	if isDryRun {

		return nil
	}

	return executeVSCodeAptPhases(opts)
}

func printVSCodeInstallPlan(tool string) {
	fmt.Printf("\n  +-- Install Plan ---------------------\n")
	fmt.Printf("  | Tool:    %s\n", tool)
	fmt.Printf("  | Version: latest\n")
	fmt.Printf("  | Manager: apt (Microsoft repo)\n")
	fmt.Printf("  | Command: add Microsoft apt repository & sudo apt install code\n")
	fmt.Printf("  +--------------------------------------\n\n")
}

func executeVSCodeAptPhases(opts installOptions) error {
	err1 := runVSCodeGpgKeyStep(opts)

	if err1 != nil {

		return err1
	}

	err2 := runVSCodeRepoStep(opts)

	if err2 != nil {

		return err2
	}

	err3 := runVSCodeUpdateStep(opts)

	if err3 != nil {

		return err3
	}

	return runVSCodeInstallStep(opts)
}

func runVSCodeGpgKeyStep(opts installOptions) error {
	fmt.Printf("  [1/4] Adding Microsoft GPG key...\n")
	cmd := []string{"sh", "-c", "wget -qO - https://packages.microsoft.com/keys/microsoft.asc | sudo tee /etc/apt/keyrings/microsoft.asc > /dev/null"}

	return runPhaseWithAudit(opts.Tool, "apt", "latest", cmd, "vscode-gpg-key", opts.Verbose)
}

func runVSCodeRepoStep(opts installOptions) error {
	fmt.Printf("  [2/4] Adding VS Code repository...\n")
	cmd := []string{"sh", "-c", "echo \"deb [arch=amd64,arm64,armhf signed-by=/etc/apt/keyrings/microsoft.asc] https://packages.microsoft.com/repos/code stable main\" | sudo tee /etc/apt/sources.list.d/vscode.list > /dev/null"}

	return runPhaseWithAudit(opts.Tool, "apt", "latest", cmd, "vscode-apt-repo", opts.Verbose)
}

func runVSCodeUpdateStep(opts installOptions) error {
	fmt.Printf("  [3/4] Updating package index...\n")
	cmd := []string{"sudo", "apt-get", "update"}

	return runPhaseWithAudit(opts.Tool, "apt", "latest", cmd, "vscode-apt-update", opts.Verbose)
}

func runVSCodeInstallStep(opts installOptions) error {
	fmt.Printf("  [4/4] Installing code (VS Code)...\n")
	cmd := []string{"sudo", "apt", "install", "-y", "code"}
	err := runPhaseWithAudit(opts.Tool, "apt", "latest", cmd, "vscode-apt-install", opts.Verbose)

	if err != nil {

		return err
	}

	recordInstallation(opts.Tool, "apt")
	fmt.Println("  ✓ VS Code installed successfully.")

	return nil
}
