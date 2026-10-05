package cmdcursor

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
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
			opts.targetNode, i = args[i+1], i+1
		} else if arg == "--force" || arg == "-f" {
			opts.isForce = true
		} else if arg == "--dry-run" || arg == "-d" {
			opts.isDryRun = true
		}
	}
	return opts
}
func installOnWindows(opts cursorInstallOptions) error {
	if _, err := exec.LookPath("winget"); err != nil {
		fmt.Printf("\n%s● Windows Setup:%s winget not found. Visit https://cursor.com\n", constants.ColorCyan, constants.ColorReset)
		return nil
	}
	if opts.isDryRun {
		fmt.Println("  [DryRun] Would execute: winget install Anysphere.Cursor --silent")
		return nil
	}
	if out, err := exec.Command("winget", "install", "Anysphere.Cursor", "--silent", "--accept-source-agreements", "--accept-package-agreements").CombinedOutput(); err != nil {
		fmt.Printf("%s⚠ Winget install returned:%s %s\n", constants.ColorYellow, constants.ColorReset, string(out))
		return nil
	}
	fmt.Printf("%s✔ Cursor IDE successfully installed via winget.%s\n", constants.ColorGreen, constants.ColorReset)
	return nil
}

func installOnLinux(opts cursorInstallOptions) error {
	fmt.Printf("\n%s● Provisioning Cursor IDE on Linux...%s\n", constants.ColorCyan, constants.ColorReset)
	if opts.isDryRun {
		fmt.Println("  [DryRun] Would invoke repo-secrets/05-scripts/setup-cursor-ubuntu.py")
		return nil
	}
	cmd := exec.Command("python3", "repo-secrets/05-scripts/setup-cursor-ubuntu.py")
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Printf("%s⚠ Linux setup returned:%s %s\n", constants.ColorYellow, constants.ColorReset, string(out))
		return nil
	}
	fmt.Printf("%s✔ Linux setup script completed.%s\n", constants.ColorGreen, constants.ColorReset)
	return nil
}

func buildRemoteSetupCmd(opts cursorInstallOptions) string {
	cmd := "python3 repo-secrets/05-scripts/setup-cursor-ubuntu.py --node " + opts.targetNode + " || bash repo-secrets/05-scripts/setup-cursor-ubuntu.sh " + opts.targetNode
	if opts.isForce {
		return "FORCE=true " + cmd
	}
	return cmd
}

func delegateRemoteInstall(opts cursorInstallOptions) error {
	fmt.Printf("\n%s● Delegating Cursor installation to remote node:%s %s\n", constants.ColorCyan, constants.ColorReset, opts.targetNode)
	if opts.isDryRun {
		fmt.Printf("  [DryRun] Would execute remote provisioning on node '%s'\n", opts.targetNode)
		return nil
	}
	remoteCmd := buildRemoteSetupCmd(opts)
	if err := cmdssh.RunSSHExec([]string{opts.targetNode, remoteCmd}); err != nil {
		fmt.Printf("%s⚠ Note: Remote setup execution notice: %v%s\n", constants.ColorYellow, err, constants.ColorReset)
		return nil
	}
	fmt.Printf("%s✔ Remote setup completed on %s.%s\n", constants.ColorGreen, opts.targetNode, constants.ColorReset)
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
