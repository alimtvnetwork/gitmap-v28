// Package cmd — macro_add_mkdir.go: in-builder mkdir helper for interactive macro creation.
package cmdmacro

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/macro"
)

func isMkdirCmd(line string) bool {
	low := strings.ToLower(line)

	return low == "mkdir" || low == ":mkdir" || low == "gitmap mkdir" ||
		strings.HasPrefix(low, "mkdir ") || strings.HasPrefix(low, ":mkdir ") ||
		strings.HasPrefix(low, "gitmap mkdir ")
}

func handleMkdirHelper(line string, state *interactiveSessionState) bool {
	parts := strings.Fields(line)
	targetDir := extractMkdirTarget(parts)
	if targetDir == "" {
		fmt.Printf("  %s▲ Usage: mkdir [-p] <directory>%s\n\n", constants.ColorYellow, constants.ColorReset)

		return true
	}

	state.lastInspectedCmd = line

	return executeInteractiveMkdir(targetDir, line)
}

func extractMkdirTarget(parts []string) string {
	for i := len(parts) - 1; i >= 0; i-- {
		token := parts[i]
		if token != "gitmap" && token != "mkdir" && token != ":mkdir" && token != "-p" {
			return token
		}
	}

	return ""
}

func executeInteractiveMkdir(targetDir, line string) bool {
	expanded := macro.ExpandPathAndEnv(strings.Trim(targetDir, "\"'"))
	absPath, err := filepath.Abs(expanded)
	if err != nil {
		fmt.Printf("  %s▲ mkdir %s: %v%s\n\n", constants.ColorRed, targetDir, err, constants.ColorReset)

		return true
	}

	if err := os.MkdirAll(absPath, 0755); err != nil {
		fmt.Printf("  %s▲ mkdir %s: %v%s\n\n", constants.ColorRed, targetDir, err, constants.ColorReset)

		return true
	}

	printInteractiveMkdirSuccess(absPath, line)

	return true
}

func printInteractiveMkdirSuccess(absPath, line string) {
	fmt.Printf("  %s✓ Created directory: %s%s\n", constants.ColorGreen, absPath, constants.ColorReset)
	fmt.Printf("  %s(created directory live. Enter command for Step, or '+add' to record '%s')%s\n\n",
		constants.ColorDim, line, constants.ColorReset)
}
