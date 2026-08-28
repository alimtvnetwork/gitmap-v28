package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runOpen is the entrypoint for `gitmap open` / `o`.
// It opens the specified directory (or the repo root, or cwd) using the OS's default opener.
func runOpen(args []string) error {
	checkHelp(constants.CmdOpen, args)

	target, err := resolveOpenTarget(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrOpenResolveCwd, err)
		return err
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", target)
	case "darwin":
		cmd = exec.Command("open", target)
	default:
		cmd = exec.Command("xdg-open", target)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to open target %q: %w", target, err)
	}

	return cmd.Wait()
}

// resolveOpenTarget picks the directory to open. Prefers args[0] if provided,
// else the git toplevel (so running `open` from a subfolder still opens the repo
// root), and falls back to plain cwd when git isn't available or the folder isn't a repo.
func resolveOpenTarget(args []string) (string, error) {
	if len(args) > 0 && args[0] != "" {
		resolved := resolveEndpointString(args[0])
		return filepath.Abs(resolved)
	}

	if root, err := gitTopLevel(); err == nil && len(root) > 0 {
		return root, nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return filepath.Abs(cwd)
}
