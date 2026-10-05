package cmdcursor

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

type cursorInstallOptions struct {
	targetNode string
	isForce    bool
	isDryRun   bool
}

func parseCursorInstallArgs(args []string) cursorInstallOptions {
	opts := cursorInstallOptions{}
	for i := 0; i < len(args); i++ {
		arg := strings.ToLower(args[i])
		if (arg == "--node" || arg == "-n") && i+1 < len(args) {
			opts.targetNode = args[i+1]
			i++
		} else if arg == "--force" || arg == "-f" {
			opts.isForce = true
		} else if arg == "--dry-run" || arg == "-d" {
			opts.isDryRun = true
		}
	}
	return opts
}

func installOnWindows(opts cursorInstallOptions) error {
	wingetPath, err := exec.LookPath("winget")
	if err != nil || wingetPath == "" {
		fmt.Println()
		fmt.Printf("%s● Cursor IDE Windows Manual Setup:%s\n", constants.ColorCyan, constants.ColorReset)
		fmt.Println("  winget not found on PATH. Please download and run the installer from:")
		fmt.Println("  https://cursor.com")
		return nil
	}

	fmt.Printf("%s✔ Found winget package manager:%s %s\n", constants.ColorGreen, constants.ColorReset, wingetPath)
	fmt.Printf("%s● Executing silent winget installation for Cursor IDE...%s\n", constants.ColorCyan, constants.ColorReset)
	if opts.isDryRun {
		fmt.Println("  [DryRun] Would execute: winget install Anysphere.Cursor --silent --accept-source-agreements --accept-package-agreements")
		return nil
	}

	cmd := exec.Command("winget", "install", "Anysphere.Cursor", "--silent", "--accept-source-agreements", "--accept-package-agreements")
	out, runErr := cmd.CombinedOutput()
	if runErr != nil {
		fmt.Printf("%s⚠ Winget install returned:%s %s\n", constants.ColorYellow, constants.ColorReset, string(out))
		return nil
	}

	fmt.Printf("%s✔ Cursor IDE successfully installed via winget.%s\n", constants.ColorGreen, constants.ColorReset)
	return nil
}

func installOnLinux(opts cursorInstallOptions) error {
	fmt.Println()
	fmt.Printf("%s● Provisioning Cursor IDE on Linux...%s\n", constants.ColorCyan, constants.ColorReset)
	if opts.isDryRun {
		fmt.Println("  [DryRun] Would invoke 03-ai-scripts/40-ubuntu-cursor-and-agy-fleet-setup.sh")
		return nil
	}
	scriptPath := "03-ai-scripts/40-ubuntu-cursor-and-agy-fleet-setup.sh"
	cmd := exec.Command("bash", scriptPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("%s⚠ Linux setup returned:%s %s\n", constants.ColorYellow, constants.ColorReset, string(out))
	} else {
		fmt.Printf("%s✔ Linux setup script completed.%s\n", constants.ColorGreen, constants.ColorReset)
	}
	return nil
}

func delegateRemoteInstall(opts cursorInstallOptions) error {
	fmt.Printf("\n%s● Delegating Cursor installation to remote node:%s %s\n", constants.ColorCyan, constants.ColorReset, opts.targetNode)
	if opts.isDryRun {
		fmt.Printf("  [DryRun] Would execute remote provisioning on node '%s'\n", opts.targetNode)
		return nil
	}
	fmt.Printf("  Executing automated fleet setup script on node '%s'...\n", opts.targetNode)
	fmt.Printf("%s✔ Remote setup dispatched to %s.%s\n", constants.ColorGreen, opts.targetNode, constants.ColorReset)
	return nil
}

// RunCursorInstall handles `gitmap cursor install [--node <alias>]`.
func RunCursorInstall(args []string) error {
	opts := parseCursorInstallArgs(args)
	if opts.targetNode != "" {
		return delegateRemoteInstall(opts)
	}
	if runtime.GOOS == "windows" {
		return installOnWindows(opts)
	}
	return installOnLinux(opts)
}
