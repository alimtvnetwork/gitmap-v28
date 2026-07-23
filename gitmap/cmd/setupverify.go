package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/completion"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// verifyShellWrapper checks if the shell wrapper is active after setup.
func verifyShellWrapper(dryRun bool) {
	if dryRun {
		return
	}

	shell := completion.DetectShell()
	fmt.Printf("\n  %s■ Shell Wrapper Verify —%s\n", constants.ColorYellow, constants.ColorReset)

	if isWrapperActive() {
		fmt.Printf(constants.MsgWrapperVerifyOK, constants.ColorGreen, constants.ColorReset)

		return
	}

	printWrapperReloadTip(shell)
}

// isWrapperActive returns true if the command wrapper env var is set.
func isWrapperActive() bool {
	return os.Getenv(constants.EnvGitmapCommandWrapper) == constants.EnvGitmapWrapperVal
}

// printWrapperReloadTip prints reload instructions for the detected shell.
func printWrapperReloadTip(shell string) {
	fmt.Printf(constants.MsgWrapperVerifyTip,
		constants.ColorYellow, constants.ColorReset,
		constants.ColorCyan, constants.ColorReset,
		constants.ColorCyan, constants.ColorReset,
		constants.ColorCyan, constants.ColorReset,
	)
}

// warnIfNoWrapper prints a stderr warning when cd is called without wrapper.
// As of v5.18.0 it ALSO auto-runs the full `gitmap setup` (idempotent,
// marker-guarded) so the very first `gitmap cd` after a fresh install
// self-heals: shell wrapper + completions get installed, the user just
// needs to reload their profile / open a new terminal. We still print
// the reload tip so the user knows why the next `cd` will actually
// move the parent shell.
func warnIfNoWrapper() {
	if isWrapperActive() {
		return
	}

	autoRunSetupForCD()
	printNoWrapperWarning()
}

// autoRunSetupForCD invokes `gitmap setup` as a best-effort self-heal
// from inside `gitmap cd` when the shell wrapper isn't loaded. All
// failures are non-fatal — the cd path was already printed to stdout
// before we got here, so the user can still copy-paste it manually.
func autoRunSetupForCD() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "  (auto-setup skipped: %v)\n", r)
		}
	}()
	fmt.Fprintln(os.Stderr, "  → Shell wrapper not detected — auto-running 'gitmap setup' to install it...")
	runSetup(nil)
}

func installWrapperForCurrentShell() {
	shell := completion.DetectShell()
	if err := completion.InstallCDFunction(shell); err != nil {
		fmt.Fprintf(os.Stderr, constants.WarnWrapperInstallFmt, err)
	}
}

func printNoWrapperWarning() {
	fmt.Fprintf(os.Stderr, constants.MsgWrapperNotLoaded,
		constants.ColorYellow, constants.ColorReset,
		constants.ColorCyan, constants.ColorReset,
		constants.ColorCyan, constants.ColorReset,
		constants.ColorCyan, constants.ColorReset,
	)
}
