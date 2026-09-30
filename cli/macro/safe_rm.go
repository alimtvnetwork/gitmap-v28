// Package macro — safe_rm.go: cross-platform idempotent file and directory removal adapter.
package macro

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// AdaptCommandForPlatform wraps commands that might fail on non-existent targets into idempotent commands.
func AdaptCommandForPlatform(cmdText string) string {
	if runtime.GOOS != constants.OSWindows {
		return cmdText
	}

	trimmed := strings.TrimSpace(cmdText)
	low := strings.ToLower(trimmed)
	if isWindowsRemovalCmd(low) {
		return transformToWindowsSafeRemoval(trimmed)
	}

	return cmdText
}

func isWindowsRemovalCmd(low string) bool {
	return strings.HasPrefix(low, "rm ") ||
		strings.HasPrefix(low, "rm\t") ||
		strings.HasPrefix(low, "rmdir ") ||
		strings.HasPrefix(low, "rmdir\t") ||
		strings.HasPrefix(low, "remove-item ") ||
		strings.HasPrefix(low, "rd ")
}

func transformToWindowsSafeRemoval(cmdText string) string {
	parts := strings.Fields(cmdText)
	if len(parts) <= 1 {
		return cmdText
	}

	targets := extractRemovalTargets(parts[1:])
	if len(targets) == 0 {
		return cmdText
	}

	targetArray := strings.Join(targets, ", ")

	return fmt.Sprintf("foreach ($__target in @(%s)) { if (Test-Path -LiteralPath $__target) { Remove-Item -Recurse -Force -LiteralPath $__target } }", targetArray)
}

func extractRemovalTargets(args []string) []string {
	var targets []string
	for _, p := range args {
		if isRemovalFlag(p) {
			continue
		}
		clean := strings.Trim(p, `"'`)
		if len(clean) > 0 {
			targets = append(targets, fmt.Sprintf("'%s'", clean))
		}
	}

	return targets
}

func isRemovalFlag(token string) bool {
	return strings.HasPrefix(token, "-") || strings.HasPrefix(token, "/")
}
