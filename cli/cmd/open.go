package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// runOpen is the entrypoint for `gitmap open` / `o`.
// It opens the specified directory (or the repo root, or cwd) using the OS's default opener.
func runOpen(args []string) error {
	checkHelp(constants.CmdOpen, args)

	if len(args) > 0 && isEditorTarget(args[0]) {
		return launchEditor(args)
	}

	if len(args) > 0 && isWebTarget(args[0]) {
		return launchNativeOpener(normalizeWebUrl(args[0]))
	}

	target, err := resolveOpenTarget(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrOpenResolveCwd, err)

		return err
	}

	return launchNativeOpener(target)
}

// resolveOpenTarget picks the directory to open. Prefers args[0] if provided,
// else the git toplevel (so running `open` from a subfolder still opens the repo
// root), and falls back to plain cwd when git isn't available or the folder isn't a repo.
func resolveOpenTarget(args []string) (string, error) {
	if len(args) > 0 && args[0] != "" {
		expanded := expandTilde(args[0])
		resolved := resolveEndpointString(expanded)

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
