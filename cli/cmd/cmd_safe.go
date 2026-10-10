// Package cmd — safe_rm_cmd.go: idempotent safe file and directory removal command.
package cmd

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/fsutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/macro"
	"os"
	"strings"
)

// runSafeRmCLI implements the `gitmap safe-rm` / `gitmap rm-safe` CLI command.
// Safely removes files or directories, exiting with 0 even if target does not exist.
func RunSafeRmCLI(args []string) error {
	targets := filterSafeRmTargets(args)
	if len(targets) == 0 {
		printSafeRmUsage()

		return nil
	}

	for _, target := range targets {
		executeSafeRmTarget(target)
	}

	return nil
}

func filterSafeRmTargets(args []string) []string {
	var targets []string
	for _, a := range args {
		clean := strings.TrimSpace(a)
		if strings.HasPrefix(clean, "-") {
			continue
		}
		if len(clean) > 0 {
			targets = append(targets, clean)
		}
	}

	return targets
}

func executeSafeRmTarget(target string) {
	expanded := macro.ExpandPathAndEnv(target)
	if !isSafeRmPathExists(expanded) {
		fmt.Printf("  %s✓ Target already absent: %s%s\n", constants.ColorDim, expanded, constants.ColorReset)

		return
	}

	if err := fsutil.SafeRemoveAll(expanded); err != nil {
		fmt.Fprintf(os.Stderr, "  %s▲ Warning removing %s: %v%s\n", constants.ColorYellow, expanded, err, constants.ColorReset)

		return
	}

	fmt.Printf("  %s✔ Safely removed: %s%s\n", constants.ColorGreen, expanded, constants.ColorReset)
}

func isSafeRmPathExists(path string) bool {
	_, err := os.Stat(path)

	return err == nil
}

func printSafeRmUsage() {
	fmt.Println("Usage: gitmap safe-rm <path...> [--force]")
	fmt.Println()
	fmt.Println("Idempotently and safely removes files or directories without failing if targets do not exist.")
	fmt.Println()
}
