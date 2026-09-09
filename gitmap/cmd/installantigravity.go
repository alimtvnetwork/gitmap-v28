package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func runInstallAntigravityWithOpts(opts installOptions) error {
	ver, isFound := verifyAndRecordAntigravity()
	if isFound {
		fmt.Printf("  ✓ Antigravity is already installed (%s)\n", ver)

		return nil
	}
	if opts.DryRun {
		fmt.Println("  [dry-run] Would download and install Antigravity CLI (agy)")

		return nil
	}

	return performAntigravityInstall()
}

func performAntigravityInstall() error {
	fmt.Println("Installing Antigravity (agy) CLI...")
	if err := executeAntigravityInstaller(); err != nil {
		fmt.Fprintf(os.Stderr, "Installer failed: %v, trying npm fallback...\n", err)
		_ = runAgyNpmFallback()
	}
	ver, isFound := verifyAndRecordAntigravity()
	if isFound {
		fmt.Printf(constants.ColorGreen+"✓"+constants.ColorReset+" Antigravity installed: %s\n", ver)

		return nil
	}
	reportVerificationFailure(constants.ToolAntigravity, "agy")

	return fmt.Errorf("antigravity installation verification failed")
}

func executeAntigravityInstaller() error {
	scriptPath, err := downloadAndValidateAgyScript()
	if err != nil {

		return err
	}
	cmd := buildScriptExecutionCmd(scriptPath)
	if cmd == nil {

		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr

	return cmd.Run()
}

func buildScriptExecutionCmd(scriptPath string) *exec.Cmd {
	if runtime.GOOS == "windows" {

		return exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-File", scriptPath)
	}

	return exec.Command("bash", scriptPath)
}

func verifyAndRecordAntigravity() (string, bool) {
	bin := resolveToolBinaryPath("agy")
	if bin == "" {
		bin = resolveToolBinaryPath("antigravity")
	}
	if bin == "" {

		return "", false
	}
	out, err := exec.Command(bin, "--version").Output()
	if err != nil {

		return "", false
	}
	ver := parseVersionFromOutput(string(out))
	if ver == "" {
		ver = strings.TrimSpace(string(out))
	}
	recordAgyInstalled(ver)

	return ver, true
}
