package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

const (
	chromeDebURL          = "https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb"
	chromeDebStage        = "/tmp/google-chrome-stable_current_amd64.deb"
	phaseChromeUpdate     = "chrome-apt-update"
	phaseChromeFetchUtils = "chrome-fetch-utils"
	phaseChromeDownload   = "chrome-download-deb"
	phaseChromeInstall    = "chrome-apt-install-deb"
	phaseChromeCleanup    = "chrome-cleanup-deb"
	phaseChromeVerify     = "chrome-verify-version"
)

type chromePhaseStep struct {
	PhaseName string
	Command   []string
	Notice    string
}

func runInstallChromeLinux(opts installOptions) error {
	if isChromeInstalledLinux() {
		fmt.Println("  ✓ Google Chrome is already installed.")

		return nil
	}

	if checkChromeDryRun(opts) {

		return nil
	}

	tool := resolveChromeTool(opts.Tool)
	version := resolveChromeVersion(opts.Version)

	return executeChromePipeline(tool, version, opts.Verbose)
}

func checkChromeDryRun(opts installOptions) bool {
	announceChromeInstallPlan()

	return handleDryRunInstall(opts.DryRun, constants.PkgMgrApt, buildChromeInstallCmd())
}

func resolveChromeTool(tool string) string {
	if tool != "" {

		return tool
	}

	return constants.CmdChrome
}

func resolveChromeVersion(version string) string {
	if version != "" {

		return version
	}

	return "official-deb"
}

func announceChromeInstallPlan() {
	fmt.Printf("\n  +-- Install Plan ---------------------\n")
	fmt.Printf("  | Tool:    chrome\n")
	fmt.Printf("  | Version: official direct deb\n")
	fmt.Printf("  | Manager: apt (Google official repo)\n")
	fmt.Printf("  | Pipeline: 6-step deb (update, fetch utils, download, apt install, cleanup, verify)\n")
	fmt.Printf("  +--------------------------------------\n\n")
}

func executeChromePipeline(tool, version string, verbose bool) error {
	err := runChromePhases(tool, version, verbose)

	if err != nil {

		return err
	}

	finishChromeInstall(tool)

	return nil
}

func runChromePhases(tool, version string, verbose bool) error {
	steps := getChromePipelineSteps()

	for _, step := range steps {
		err := executeSingleChromeStep(tool, version, step, verbose)

		if err != nil {

			return err
		}
	}

	return nil
}

func executeSingleChromeStep(tool, version string, step chromePhaseStep, verbose bool) error {
	fmt.Println(step.Notice)
	err := runPhaseWithAudit(tool, constants.PkgMgrApt, version, step.Command, step.PhaseName, verbose)

	if err != nil {

		return err
	}

	return nil
}

func finishChromeInstall(tool string) {
	_ = os.Remove(chromeDebStage)
	recordInstallation(tool, constants.PkgMgrApt)
	fmt.Println("  ✓ Google Chrome installed successfully via 6-step deb pipeline.")
}

func getChromePipelineSteps() []chromePhaseStep {

	return []chromePhaseStep{
		{phaseChromeUpdate, buildChromeUpdateCmd(), "  [1/6] Updating APT package repositories..."},
		{phaseChromeFetchUtils, buildChromeFetchUtilsCmd(), "  [2/6] Ensuring fetch utilities (wget, curl)..."},
		{phaseChromeDownload, buildChromeDownloadCmd(), "  [3/6] Downloading official Chrome deb package..."},
		{phaseChromeInstall, buildChromeInstallCmd(), "  [4/6] Installing Chrome deb and resolving dependencies..."},
		{phaseChromeCleanup, buildChromeCleanupCmd(), "  [5/6] Cleaning up scratch package..."},
		{phaseChromeVerify, buildChromeVerifyCmd(), "  [6/6] Verifying Google Chrome binary version..."},
	}
}

func buildChromeUpdateCmd() []string {

	return []string{"sudo", "apt-get", "update"}
}

func buildChromeFetchUtilsCmd() []string {

	return []string{"sudo", "apt-get", "install", "-y", "wget", "curl"}
}

func buildChromeDownloadCmd() []string {
	if isBinaryOnPath("wget") {

		return []string{"wget", "-q", "-O", chromeDebStage, chromeDebURL}
	}

	return []string{"curl", "-fSL", "-o", chromeDebStage, chromeDebURL}
}

func buildChromeInstallCmd() []string {

	return []string{"sudo", "apt-get", "install", "-y", chromeDebStage}
}

func buildChromeCleanupCmd() []string {

	return []string{"rm", "-f", chromeDebStage}
}

func buildChromeVerifyCmd() []string {

	return []string{"google-chrome", "--version"}
}

func isChromeInstalledLinux() bool {
	bins := []string{"google-chrome", "google-chrome-stable", "chromium-browser", "chromium"}

	for _, bin := range bins {
		if isBinaryOnPath(bin) {

			return true
		}
	}

	return false
}

func isBinaryOnPath(bin string) bool {
	_, err := exec.LookPath(bin)

	return err == nil
}
