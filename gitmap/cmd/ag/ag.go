package ag

import (
	"os"
	"os/exec"
	"strings"

	"github.com/pterm/pterm"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

// Run executes the ag command.
func Run(args []string) error {
	if len(args) > 0 && strings.EqualFold(args[0], "install") {
		return runAgInstall()
	}

	// Just launch ag in CWD.
	pterm.Info.Println("Launching Antigravity...")
	cmd := exec.Command("ag", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return apperror.WrapSimple(err, "Run: failed to launch antigravity")
	}
	return nil
}

func runAgInstall() error {
	pterm.Info.Println("Running gitmap install antigravity...")
	cmd := exec.Command("gitmap", "install", "ag-ctx")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return apperror.WrapSimple(err, "runAgInstall: install failed")
	}
	return nil
}
