package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// cleanupRemoveMaxAttempts caps the retry loop for transient Windows file locks
// (e.g. a freshly-renamed .old still settling, or AV scanners holding a handle).
const (
	cleanupRemoveMaxAttempts = 5
	cleanupRemoveRetryDelay  = 200 * time.Millisecond
)

// cleanupTempArtifacts removes update handoff copies and generated scripts.
func cleanupTempArtifacts(ctx updateCleanupContext) int {
	return removeCleanupPatterns(ctx.tempPatterns, ctx.selfPath, constants.MsgUpdateTempRemoved)
}

// cleanupBackupArtifacts removes .old binaries left by deploy and PATH sync.
func cleanupBackupArtifacts(ctx updateCleanupContext) int {
	return removeCleanupPatterns(ctx.backupPatterns, ctx.selfPath, constants.MsgUpdateOldRemoved)
}

// removeCleanupPatterns removes every file matched by the provided glob list.
func removeCleanupPatterns(patterns []string, selfPath, successMsg string) int {
	seen := map[string]bool{}
	cleaned := 0
	for _, pattern := range patterns {
		cleaned += removeCleanupPattern(pattern, selfPath, seen, successMsg)
	}

	return cleaned
}

// removeCleanupPattern removes files for a single glob pattern.
func removeCleanupPattern(pattern, selfPath string, seen map[string]bool, successMsg string) int {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		logUpdateCleanupGlobError(pattern, err)
		logHandoffEvent("cleanup", "glob_error", map[string]string{
			"pattern": pattern,
			"err":     err.Error(),
		})
		emitDebugWindowsJSON("glob_error", map[string]any{
			"pattern": pattern,
			"err":     err.Error(),
		})

		return 0
	}

	cleaned := 0
	for _, match := range matches {
		if removeCleanupMatch(match, selfPath, seen, successMsg) {
			cleaned++
		}
	}

	return cleaned
}

// removeCleanupMatch removes a single cleanup candidate once.
func removeCleanupMatch(match, selfPath string, seen map[string]bool, successMsg string) bool {
	cleanPath := filepath.Clean(match)
	normalizedPath := normalizeCleanupPath(cleanPath)
	if hasSeenCleanupPath(seen, normalizedPath) {
		return false
	}
	if isActiveCleanupPath(normalizedPath, selfPath) {
		return false
	}

	return removeCleanupFile(match, cleanPath, successMsg)
}

// hasSeenCleanupPath reports whether this cleanup path was already processed.
func hasSeenCleanupPath(seen map[string]bool, normalizedPath string) bool {
	if seen[normalizedPath] {
		return true
	}

	seen[normalizedPath] = true

	return false
}

// isActiveCleanupPath reports whether the candidate points to the active binary.
func isActiveCleanupPath(normalizedPath, selfPath string) bool {
	return len(selfPath) > 0 && normalizedPath == normalizeCleanupPath(selfPath)
}

// removeCleanupFile removes a cleanup candidate and prints the success message.
// Retries a few times with a small delay because Windows can briefly hold a
// handle on a file we just renamed (e.g. gitmap.exe.old) or on a binary that
// our handoff process is still releasing.
func removeCleanupFile(match, cleanPath, successMsg string) bool {
	var lastErr error
	for attempt := 1; attempt <= cleanupRemoveMaxAttempts; attempt++ {
		err := os.Remove(match)
		if err == nil {
			fmt.Printf(successMsg, filepath.Base(match))
			logHandoffEvent("cleanup", "remove_ok", map[string]string{
				"path": cleanPath,
			})

			return true
		}

		lastErr = err
		if attempt < cleanupRemoveMaxAttempts {
			logHandoffEvent("cleanup", "remove_retry", map[string]string{
				"path":    cleanPath,
				"attempt": fmt.Sprintf("%d", attempt),
				"err":     err.Error(),
			})
			emitDebugWindowsJSON("remove_retry", map[string]any{
				"path":    cleanPath,
				"attempt": attempt,
				"err":     err.Error(),
			})
			time.Sleep(cleanupRemoveRetryDelay)
		}
	}

	logUpdateCleanupRemoveError(cleanPath, lastErr)
	logHandoffEvent("cleanup", "remove_fail", map[string]string{
		"path": cleanPath,
		"err":  lastErr.Error(),
	})
	emitDebugWindowsJSON("remove_fail", map[string]any{
		"path": cleanPath,
		"err":  lastErr.Error(),
	})

	return false
}

// logUpdateCleanupExecutableError reports os.Executable failures.
func logUpdateCleanupExecutableError(err error) {
	fmt.Fprintf(os.Stderr, constants.ErrUpdateCleanupExecPath, err)
}

// logUpdateCleanupConfigReadError reports powershell.json read failures.
func logUpdateCleanupConfigReadError(path string, err error) {
	fmt.Fprintf(os.Stderr, constants.ErrUpdateCleanupConfigRead, path, err)
}

// logUpdateCleanupGlobError reports filepath.Glob failures.
func logUpdateCleanupGlobError(path string, err error) {
	fmt.Fprintf(os.Stderr, constants.ErrUpdateCleanupGlob, path, err)
}

// logUpdateCleanupRemoveError reports os.Remove failures.
func logUpdateCleanupRemoveError(path string, err error) {
	fmt.Fprintf(os.Stderr, constants.ErrUpdateCleanupRemove, path, err)
}
