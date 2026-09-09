package cmd

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func resolveGlobalNpmCommand(pkg string) []string {
	if isCommandAvailable("sudo") {

		return []string{"sudo", "npm", "install", "-g", pkg}
	}

	return []string{"npm", "install", "-g", pkg}
}

func tryInstallNpmGlobal(tool, pkg string, verbose bool) bool {
	if !isCommandAvailable("npm") {

		return false
	}
	cmd := resolveGlobalNpmCommand(pkg)
	err := runPhaseWithAudit(tool, "npm", "latest", cmd, pkg+"-npm-install", verbose)

	return err == nil
}

func printPnpmInstallPlan(tool string) {
	fmt.Printf("\n  +-- Install Plan ---------------------\n")
	fmt.Printf("  | Tool:    %s\n", tool)
	fmt.Printf("  | Version: latest\n")
	fmt.Printf("  | Manager: npm / standalone\n")
	fmt.Printf("  | Command: sudo npm install -g pnpm (or curl standalone script)\n")
	fmt.Printf("  +--------------------------------------\n\n")
}

func isPnpmInstallPreempted(opts installOptions) bool {
	if isInstalled := alreadyInstalled(opts.Tool); isInstalled {

		return true
	}
	if opts.Manager == constants.PkgMgrBrew {
		executeGenericInstall(opts)

		return true
	}

	return false
}

func runInstallPnpmLinux(opts installOptions) error {
	if isPnpmInstallPreempted(opts) {

		return nil
	}
	printPnpmInstallPlan(opts.Tool)
	if opts.DryRun {
		fmt.Printf(constants.MsgInstallDryCmd, "sudo npm install -g pnpm")

		return nil
	}

	return executePnpmLinuxInstall(opts)
}

func executePnpmLinuxInstall(opts installOptions) error {
	err := runPnpmInstallPhase(opts)
	if err != nil {

		return err
	}
	postPnpmLinuxSetup()
	verifyInstallation(opts.Tool)
	recordInstallation(opts.Tool, "npm")

	return nil
}

func runPnpmInstallPhase(opts installOptions) error {
	if isSuccess := tryInstallNpmGlobal(opts.Tool, "pnpm", opts.Verbose); isSuccess {

		return nil
	}
	fmt.Printf("  Falling back to official standalone pnpm installer...\n")
	curlCmd := []string{"sh", "-c", "curl -fsSL https://get.pnpm.io/install.sh | sh -"}

	return runPhaseWithAudit(opts.Tool, "curl", "latest", curlCmd, "pnpm-curl-install", opts.Verbose)
}

func printYarnInstallPlan(tool string) {
	fmt.Printf("\n  +-- Install Plan ---------------------\n")
	fmt.Printf("  | Tool:    %s\n", tool)
	fmt.Printf("  | Version: latest\n")
	fmt.Printf("  | Manager: npm / corepack\n")
	fmt.Printf("  | Command: sudo npm install -g yarn (or corepack enable)\n")
	fmt.Printf("  +--------------------------------------\n\n")
}

func isYarnInstallPreempted(opts installOptions) bool {
	if isInstalled := alreadyInstalled(opts.Tool); isInstalled {

		return true
	}
	if opts.Manager == constants.PkgMgrBrew {
		executeGenericInstall(opts)

		return true
	}

	return false
}

func runInstallYarnLinux(opts installOptions) error {
	if isYarnInstallPreempted(opts) {

		return nil
	}
	printYarnInstallPlan(opts.Tool)
	if opts.DryRun {
		fmt.Printf(constants.MsgInstallDryCmd, "sudo npm install -g yarn")

		return nil
	}

	return executeYarnLinuxInstall(opts)
}

func executeYarnLinuxInstall(opts installOptions) error {
	err := runYarnInstallPhase(opts)
	if err != nil {

		return err
	}
	verifyInstallation(opts.Tool)
	recordInstallation(opts.Tool, "npm")

	return nil
}

func runYarnInstallPhase(opts installOptions) error {
	if isSuccess := tryInstallNpmGlobal(opts.Tool, "yarn", opts.Verbose); isSuccess {

		return nil
	}
	fmt.Printf("  Falling back to corepack yarn installer...\n")
	corepackCmd := []string{"sh", "-c", "corepack enable && corepack prepare yarn@stable --activate"}

	return runPhaseWithAudit(opts.Tool, "corepack", "latest", corepackCmd, "yarn-corepack-install", opts.Verbose)
}

func printBunInstallPlan(tool string) {
	fmt.Printf("\n  +-- Install Plan ---------------------\n")
	fmt.Printf("  | Tool:    %s\n", tool)
	fmt.Printf("  | Version: latest\n")
	fmt.Printf("  | Manager: standalone curl script\n")
	fmt.Printf("  | Command: curl -fsSL https://bun.sh/install | bash\n")
	fmt.Printf("  +--------------------------------------\n\n")
}

func runInstallBunLinux(opts installOptions) error {
	if isInstalled := alreadyInstalled(opts.Tool); isInstalled {

		return nil
	}
	printBunInstallPlan(opts.Tool)
	if opts.DryRun {
		fmt.Printf(constants.MsgInstallDryCmd, "curl -fsSL https://bun.sh/install | bash")

		return nil
	}

	return executeBunLinuxInstall(opts)
}

func executeBunLinuxInstall(opts installOptions) error {
	cmd := []string{"sh", "-c", "curl -fsSL https://bun.sh/install | bash"}
	err := runPhaseWithAudit(opts.Tool, "curl", "latest", cmd, "bun-curl-install", opts.Verbose)
	if err != nil {

		return err
	}
	verifyInstallation(opts.Tool)
	recordInstallation(opts.Tool, "curl")

	return nil
}
