package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	chromeDebURL   = "https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb"
	chromeDebStage = "/tmp/google-chrome-stable_current_amd64.deb"
)

func runInstallChromeLinux(opts installOptions) error {
	if isChromeInstalledLinux() {
		fmt.Println("  ✓ Google Chrome is already installed.")

		return nil
	}

	announceChromeInstallPlan()
	if handleDryRunInstall(opts.DryRun, "apt", []string{"sudo", "apt", "install", "-y", chromeDebStage}) {
		return nil
	}

	if err := downloadChromeDebLinux(opts.Verbose); err != nil {
		fmt.Printf("  ⚠ Direct deb download failed (%v); falling back to repository setup...\n", err)

		return runInstallChromeRepoFallbackLinux(opts)
	}

	if err := installChromeDebViaApt(opts); err != nil {
		fmt.Printf("  ⚠ Local deb apt installation failed (%v); falling back to repository setup...\n", err)

		return runInstallChromeRepoFallbackLinux(opts)
	}

	_ = os.Remove(chromeDebStage)
	recordInstallation(opts.Tool, "apt")
	fmt.Println("  ✓ Google Chrome installed successfully via official Debian package.")

	return nil
}

func isChromeInstalledLinux() bool {
	for _, bin := range []string{"google-chrome", "google-chrome-stable", "chromium-browser", "chromium"} {
		if _, err := exec.LookPath(bin); err == nil {
			return true
		}
	}

	return false
}

func announceChromeInstallPlan() {
	fmt.Printf("\n  +-- Install Plan ---------------------\n")
	fmt.Printf("  | Tool:    chrome\n")
	fmt.Printf("  | Version: official direct deb\n")
	fmt.Printf("  | Manager: apt (Google official repo)\n")
	fmt.Printf("  | Command: curl -fSL %s -o %s && sudo apt install -y %s\n", chromeDebURL, chromeDebStage, chromeDebStage)
	fmt.Printf("  +--------------------------------------\n\n")
}

func downloadChromeDebLinux(verbose bool) error {
	fmt.Println("  [1/3] Downloading official Google Chrome .deb package...")
	dlCmd := exec.Command("curl", "-fSL", "-o", chromeDebStage, chromeDebURL)
	if _, err := exec.LookPath("curl"); err != nil {
		dlCmd = exec.Command("wget", "-qO", chromeDebStage, chromeDebURL)
	}
	if verbose {
		dlCmd.Stdout = os.Stdout
		dlCmd.Stderr = os.Stderr
	}
	if err := dlCmd.Run(); err != nil {
		return apperror.WrapSimple(err, "install_chrome_linux.download")
	}

	return nil
}

func installChromeDebViaApt(opts installOptions) error {
	fmt.Println("  [2/3] Installing Google Chrome and resolving dependencies via apt...")
	installArgs := []string{"sudo", "apt", "install", "-y", chromeDebStage}
	output, err := execInstallCommand(installArgs, opts.Verbose)
	if err != nil {
		return apperror.NewWithDetails(
			"install_chrome_linux.apt_install",
			"E9000",
			fmt.Sprintf("apt install failed: %s (%v)", string(output), err),
			"cmd/install_chrome_linux",
			apperror.ErrorTypeExecution,
			apperror.SeverityError,
			map[string]any{"command": installArgs, "output": string(output)},
		)
	}

	return nil
}

func runInstallChromeRepoFallbackLinux(opts installOptions) error {
	fmt.Println("  [Fallback 1/3] Adding Google official APT keyring...")
	keyringCmd := exec.Command("sh", "-c", "wget -q -O - https://dl.google.com/linux/linux_signing_key.pub | sudo gpg --dearmor -o /etc/apt/keyrings/google-chrome.gpg")
	if out, err := keyringCmd.CombinedOutput(); err != nil {
		return apperror.NewWithDetails("install_chrome_linux.keyring", "E9000", fmt.Sprintf("keyring setup failed: %s (%v)", string(out), err), "cmd/install_chrome_linux", apperror.ErrorTypeExecution, apperror.SeverityError, nil)
	}

	fmt.Println("  [Fallback 2/3] Adding Google Chrome repository source...")
	repoCmd := exec.Command("sh", "-c", "echo \"deb [arch=amd64 signed-by=/etc/apt/keyrings/google-chrome.gpg] http://dl.google.com/linux/chrome/deb/ stable main\" | sudo tee /etc/apt/sources.list.d/google-chrome.list > /dev/null")
	if out, err := repoCmd.CombinedOutput(); err != nil {
		return apperror.NewWithDetails("install_chrome_linux.repo", "E9000", fmt.Sprintf("repo setup failed: %s (%v)", string(out), err), "cmd/install_chrome_linux", apperror.ErrorTypeExecution, apperror.SeverityError, nil)
	}

	fmt.Println("  [Fallback 3/3] Updating apt and installing google-chrome-stable...")
	runAptUpdate(opts.Verbose)
	runInstallCommand([]string{"sudo", "apt", "install", "-y", "google-chrome-stable"}, opts)
	recordInstallation(opts.Tool, "apt")

	return nil
}
