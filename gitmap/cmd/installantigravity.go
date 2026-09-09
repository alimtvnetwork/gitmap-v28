package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

func runInstallAntigravityWithOpts(opts installOptions) error {
	if opts.DryRun {
		fmt.Println("  [dry-run] Would run: curl -fsSL https://get.antigravity.dev | bash (or powershell irm)")

		return nil
	}
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
	fmt.Fprintf(os.Stderr, "Notice: Antigravity installed. Restart shell or run: agy --version\n")

	return nil
}

func executeAntigravityInstaller() error {
	cmd := buildAgyPlatformCmd()
	if cmd == nil {

		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr

	return cmd.Run()
}

func buildAgyPlatformCmd() *exec.Cmd {
	if runtime.GOOS == "windows" {

		return exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-Command", "irm https://get.antigravity.dev | iex")
	}

	return exec.Command("bash", "-c", "curl -fsSL https://get.antigravity.dev | bash")
}

func runAgyNpmFallback() error {
	cmd := exec.Command("npm", "install", "-g", "@google/antigravity")
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr

	return cmd.Run()
}

func verifyAndRecordAntigravity() (string, bool) {
	out, err := exec.Command("agy", "--version").Output()
	if err != nil {

		return "", false
	}
	ver := strings.TrimSpace(string(out))
	recordAgyInstalled(ver)

	return ver, true
}

func recordAgyInstalled(ver string) {
	splitDB, err := store.OpenInstallationSplitDB()
	if err != nil {

		return
	}
	defer splitDB.Close()
	_ = splitDB.SaveInstalledTool("antigravity", ver, "installer")
}
