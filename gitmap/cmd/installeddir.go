package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

// runInstalledDir prints the directory and full path of the active gitmap binary.
func runInstalledDir() error {
	selfPath, err := os.Executable()
	if err != nil {
		return apperror.Wrap(err, "✗ Could not resolve executable path:", nil)
	}

	resolved, err := filepath.EvalSymlinks(selfPath)
	if err != nil {
		resolved = selfPath
	}

	absPath, err := filepath.Abs(resolved)
	if err != nil {
		absPath = resolved
	}

	dir := filepath.Dir(absPath)

	fmt.Printf("\n  📂 Installed directory\n\n")
	fmt.Printf("  Binary:    %s\n", absPath)
	fmt.Printf("  Directory: %s\n\n", dir)
	return nil
}
