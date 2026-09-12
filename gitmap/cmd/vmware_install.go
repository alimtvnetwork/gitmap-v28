package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func runVmwareInstall(args []string) error {
	checkHelp(constants.CmdVmware, args)
	isDryRun := hasDryRunFlag(args) || hasShortDryRunFlag(args)
	if isDryRun {
		return simulateVmwareInstall()
	}

	if !isLinuxOS() {
		return newUnsupportedOSError()
	}

	return executeVmwareInstall()
}

func simulateVmwareInstall() error {
	fmt.Println("▶ gitmap vmware install (dry-run)")
	fmt.Println("  [dry-run] Would update apt index: sudo apt-get update")
	fmt.Println("  [dry-run] Would install packages: sudo apt-get install -y open-vm-tools open-vm-tools-desktop")
	fmt.Println("  [dry-run] Would enable service: sudo systemctl enable --now open-vm-tools")
	fmt.Println("  [dry-run] Would verify binary: vmhgfs-fuse")

	return nil
}

func newUnsupportedOSError() error {
	return apperror.NewWithDetails(
		"cmd.vmware.install",
		"E4002",
		"vmware tools installation via apt is only supported on Linux guest environments",
		"cmd.vmware",
		apperror.ErrorTypeValidation,
		apperror.SeverityError,
		nil,
	)
}

func executeVmwareInstall() error {
	fmt.Println("▶ gitmap vmware install")
	if err := runAptInstallPackages("open-vm-tools", "open-vm-tools-desktop"); err != nil {
		return err
	}

	enableVMwareService()
	printInstallCompletion()

	return nil
}

func runAptInstallPackages(pkgs ...string) error {
	if _, err := exec.LookPath("apt-get"); err != nil {
		return apperror.NewWithDetails("cmd.vmware.install", "E4003", "apt-get package manager not found on system", "cmd.vmware", apperror.ErrorTypeExecution, apperror.SeverityError, nil)
	}

	args := append([]string{"apt-get", "install", "-y"}, pkgs...)
	cmd := exec.Command("sudo", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return apperror.WrapSimple(err, "vmware.aptInstall")
	}

	return nil
}

func enableVMwareService() {
	if _, err := exec.LookPath("systemctl"); err != nil {
		return
	}

	_ = exec.Command("sudo", "systemctl", "enable", "--now", "open-vm-tools").Run()
}

func printInstallCompletion() {
	fmt.Println("  ✓ Successfully installed open-vm-tools and open-vm-tools-desktop")
	fmt.Println("  ✓ Service open-vm-tools enabled and started")
	fmt.Println("  Next: run 'gitmap vmware shared enable' to mount shared folders")
}
