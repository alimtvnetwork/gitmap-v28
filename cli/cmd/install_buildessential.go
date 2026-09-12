package cmd

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

var buildEssentialPackages = []string{
	"build-essential", "vim", "wget", "nano", "curl", "file",
	"git", "zlib1g", "zlib1g-dev", "libssl-dev", "git-core",
	"sshpass", "zsh", "git-lfs", "snapd",
}

func isBuildEssentialAlias(tool string) bool {
	norm := strings.ToLower(strings.TrimSpace(tool))

	return norm == "build-essential" || norm == "buildessential" || norm == "be" || norm == "ubuntu-common" || norm == "ub-common"
}

func getBuildEssentialPackages() []string {
	cp := make([]string, len(buildEssentialPackages))
	copy(cp, buildEssentialPackages)

	return cp
}

func runInstallBuildEssential(opts installOptions) error {
	if opts.DryRun {
		fmt.Printf("Dry-run: would execute: sudo apt-get update && sudo apt-get install -y %s\n", strings.Join(buildEssentialPackages, " "))

		return nil
	}

	if runtime.GOOS != "linux" {
		fmt.Println("  ℹ Ubuntu common 'build-essential' profile is optimized for Debian/Ubuntu Linux.")

		return nil
	}

	if err := executeAptInstall(buildEssentialPackages); err != nil {
		return err
	}

	recordBuildEssentialTools()
	printProfileInstallSummary("build-essential")

	return nil
}

func executeAptInstall(packages []string) error {
	fmt.Println("▶ Updating apt package repositories...")
	_ = exec.Command("sudo", "apt-get", "update").Run()

	fmt.Printf("▶ Installing %d packages...\n", len(packages))
	args := append([]string{"apt-get", "install", "-y"}, packages...)
	cmd := exec.Command("sudo", args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		return apperror.NewWithDetails("cmd.install_buildessential", "E4005", fmt.Sprintf("apt-get install failed: %s (%v)", string(out), err), "cmd.install_buildessential", apperror.ErrorTypeExecution, apperror.SeverityError, nil)
	}

	return nil
}

func recordBuildEssentialTools() {
	splitDB, err := store.OpenInstallationSplitDB()
	if err != nil {
		return
	}

	defer splitDB.Close()

	for _, tool := range buildEssentialPackages {
		version := probeToolVersion(tool)
		_ = splitDB.SaveInstalledTool(tool, version, "apt")
	}

	for _, tool := range []string{"gcc", "g++", "make"} {
		version := probeToolVersion(tool)
		_ = splitDB.SaveInstalledTool(tool, version, "apt")
	}

	_ = splitDB.RecordLog(store.InstallationLogRecord{
		Tool:           "build-essential",
		Action:         "profile-install",
		PackageManager: "apt",
		IsSuccess:      true,
	})
}

func probeToolVersion(tool string) string {
	out, err := exec.Command("dpkg-query", "-W", "-f=${Version}", tool).Output()
	if err == nil && len(strings.TrimSpace(string(out))) > 0 {
		return strings.TrimSpace(string(out))
	}

	outCmd, errCmd := exec.Command(tool, "--version").Output()
	lines := strings.Split(string(outCmd), "\n")
	if errCmd == nil && len(lines) > 0 {
		return strings.TrimSpace(lines[0])
	}

	return "installed"
}
