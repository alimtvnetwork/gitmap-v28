// Package cmd — extra cleanup passes for update-cleanup.
//
// These complement the pattern-based pass in updatecleanup_remove.go by
// targeting two artifact classes that don't fit the simple-glob model:
//  1. The obsolete v2.90.0 drive-root forwarding shim
//     (e.g. E:\gitmap.exe sitting at the literal drive root, NOT
//     inside a gitmap\ subfolder).
//  2. *.gitmap-tmp-* swap directories left by interrupted clones.
//
// Both passes follow the spec/04-generic-cli/22-data-folder-deploy-and-cleanup.md
// contract (DFD-6, DFD-7).
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

const (
	driveRootShimMaxBytes = 5 * 1024 * 1024
	cloneSwapDirGlob      = "*.gitmap-tmp-*"
)

// cleanupDriveRootShim removes the obsolete drive-root forwarding shim if present.
// Safe-by-default: only deletes when the candidate sits at the literal drive root
// AND is under the size cap AND is not inside a gitmap/ folder.
func cleanupDriveRootShim(ctx updateCleanupContext) int {
	if runtime.GOOS != "windows" {
		return 0
	}

	shimPath := resolveDriveRootShimPath(ctx.selfPath)
	if len(shimPath) == 0 {
		return 0
	}

	if !isRemovableDriveRootShim(shimPath, ctx.selfPath) {
		return 0
	}

	return removeDriveRootShim(shimPath)
}

// resolveDriveRootShimPath returns <drive>:\<binaryName> derived from selfPath.
func resolveDriveRootShimPath(selfPath string) string {
	if len(selfPath) == 0 {
		return ""
	}

	drive := filepath.VolumeName(selfPath)
	if len(drive) == 0 {
		return ""
	}

	binaryName := filepath.Base(selfPath)

	return filepath.Join(drive+`\`, binaryName)
}

// isRemovableDriveRootShim returns true only when the candidate is safe to delete.
func isRemovableDriveRootShim(shimPath, selfPath string) bool {
	if normalizeCleanupPath(shimPath) == normalizeCleanupPath(selfPath) {
		return false
	}

	parent := filepath.Dir(shimPath)
	if !isLiteralDriveRoot(parent) {
		return false
	}

	info, err := os.Stat(shimPath)
	if err != nil || info.IsDir() {
		return false
	}
	if info.Size() > driveRootShimMaxBytes {
		fmt.Fprintf(os.Stderr, "  !! [cleanup] skipping drive-root shim %s (size %d > 5MB)\n", shimPath, info.Size())
		logHandoffEvent("cleanup", "drive_root_skip", map[string]string{
			"path":   shimPath,
			"reason": "size_guard",
			"bytes":  fmt.Sprintf("%d", info.Size()),
		})
		emitDebugWindowsJSON("drive_root_skip", map[string]any{
			"path":   shimPath,
			"reason": "size_guard",
			"bytes":  info.Size(),
		})

		return false
	}

	return true
}

// isLiteralDriveRoot returns true when path is the literal drive root (e.g. "E:\").
func isLiteralDriveRoot(path string) bool {
	clean := strings.TrimRight(path, `\/`)

	return len(clean) == 2 && clean[1] == ':'
}

// removeDriveRootShim deletes the candidate and logs the outcome.
func removeDriveRootShim(shimPath string) int {
	if err := os.Remove(shimPath); err != nil {
		fmt.Fprintf(os.Stderr, "  !! [cleanup] could not remove drive-root shim %s: %v\n", shimPath, err)
		logHandoffEvent("cleanup", "drive_root_remove_fail", map[string]string{
			"path": shimPath,
			"err":  err.Error(),
		})
		emitDebugWindowsJSON("drive_root_remove_fail", map[string]any{
			"path": shimPath,
			"err":  err.Error(),
		})

		return 0
	}

	fmt.Printf("  → Removed obsolete drive-root shim: %s\n", shimPath)
	logHandoffEvent("cleanup", "drive_root_remove_ok", map[string]string{
		"path": shimPath,
	})

	return 1
}

// cleanupCloneSwapDirs removes *.gitmap-tmp-* directories left by interrupted clones.
// Scans every cleanup directory we already resolved.
func cleanupCloneSwapDirs(ctx updateCleanupContext) int {
	dirs := uniqueParentDirs(ctx.tempPatterns, ctx.backupPatterns)
	removed := 0
	for _, dir := range dirs {
		removed += removeCloneSwapDirsIn(dir)
	}

	return removed
}

// uniqueParentDirs extracts the unique parent directories from glob patterns.
func uniqueParentDirs(patternGroups ...[]string) []string {
	seen := map[string]bool{}
	out := make([]string, 0)
	for _, group := range patternGroups {
		for _, pattern := range group {
			dir := filepath.Dir(pattern)
			key := normalizeCleanupPath(dir)
			if seen[key] {
				continue
			}
			seen[key] = true
			out = append(out, dir)
		}
	}

	return out
}

// removeCloneSwapDirsIn removes every *.gitmap-tmp-* dir directly under base.
func removeCloneSwapDirsIn(base string) int {
	matches, err := filepath.Glob(filepath.Join(base, cloneSwapDirGlob))
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrUpdateCleanupGlob, base, err)
		logHandoffEvent("cleanup", "swap_glob_error", map[string]string{
			"base": base,
			"err":  err.Error(),
		})
		emitDebugWindowsJSON("swap_glob_error", map[string]any{
			"base": base,
			"err":  err.Error(),
		})

		return 0
	}

	removed := 0
	for _, match := range matches {
		info, statErr := os.Stat(match)
		if statErr != nil || !info.IsDir() {
			continue
		}
		if err := os.RemoveAll(match); err != nil {
			fmt.Fprintf(os.Stderr, "  !! [cleanup] could not remove swap dir %s: %v\n", match, err)
			logHandoffEvent("cleanup", "swap_remove_fail", map[string]string{
				"path": match,
				"err":  err.Error(),
			})
			emitDebugWindowsJSON("swap_remove_fail", map[string]any{
				"path": match,
				"err":  err.Error(),
			})

			continue
		}
		fmt.Printf("  → Removed swap dir: %s\n", match)
		logHandoffEvent("cleanup", "swap_remove_ok", map[string]string{
			"path": match,
		})
		removed++
	}

	return removed
}
