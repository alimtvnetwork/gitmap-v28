package cmd

import (
	"fmt"
	"os"
	"path/filepath"
)

// runInstalledDir prints the directory and full path of the active gitmap binary.
func runInstalledDir() {
	selfPath, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ Could not resolve executable path: %v\n", err)
		os.Exit(1)
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
}
