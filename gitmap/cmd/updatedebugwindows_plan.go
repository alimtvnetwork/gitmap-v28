// Package cmd — extended `--debug-windows` output that prints the
// exact commands and filesystem operations the update-cleanup handoff
// will perform, so the user can audit and reproduce them.
//
// Two functions are exported within the package:
//
//	dumpDebugWindowsCommandPlan — renders the Phase 3 spawn command
//	line as a copy-pastable shell invocation (proper quoting), and
//	prints an explicit note that no `git` subprocess is launched.
//
//	dumpDebugWindowsCleanupPlan — enumerates the filepath.Glob
//	patterns the deployed binary will scan and the actual file
//	matches that will be passed to os.Remove / os.RemoveAll. Called
//	from runUpdateCleanup BEFORE any deletion happens so the output
//	reflects intent, not outcome.
//
// These complement the structured-event log written by
// updatehandofflog.go (which is post-hoc) with a pre-flight view.
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// dumpDebugWindowsCommandPlan renders the exact spawn invocation as
// a single shell-friendly line: `"<deployed>" "<arg1>" "<arg2>" ...`.
// Always quotes for cross-platform copy-paste safety. Followed by an
// explicit note that the cleanup itself is pure Go syscalls — no
// `git` or other subprocess is launched.
func dumpDebugWindowsCommandPlan(deployed string, childArgs []string) {
	if !isDebugWindowsRequested() {
		return
	}
	full := append([]string{deployed}, childArgs...)
	cmdLine := renderShellCommand(full)
	fmt.Fprintf(os.Stderr, constants.MsgDebugWinCmdLine, cmdLine)
	fmt.Fprint(os.Stderr, constants.MsgDebugWinCmdNote)
	emitDebugWindowsJSON("command_plan", map[string]any{
		"deployed": deployed, "child_args": childArgs,
		"command_line": cmdLine,
	})
}

// renderShellCommand quotes each token with double quotes (escaping
// any embedded `"`) and joins with spaces. Cross-platform safe for
// copy-paste — works in PowerShell, cmd.exe, bash, and zsh.
func renderShellCommand(tokens []string) string {
	parts := make([]string, 0, len(tokens))
	for _, t := range tokens {
		parts = append(parts, quoteShellToken(t))
	}

	return strings.Join(parts, " ")
}

// quoteShellToken wraps a token in double quotes, escaping inner `"`
// as `\"`. Tokens with no whitespace or special chars are still
// quoted for visual consistency in the dump.
func quoteShellToken(t string) string {
	escaped := strings.ReplaceAll(t, `"`, `\"`)

	return `"` + escaped + `"`
}

// dumpDebugWindowsCleanupPlan enumerates every filesystem operation
// the deployed binary will attempt during update-cleanup. Called
// from runUpdateCleanup AFTER loadUpdateCleanupContext but BEFORE
// the actual cleanupTempArtifacts/cleanupBackupArtifacts/etc. calls
// so the user sees the *plan*, not the outcome.
//
// Operations enumerated:
//  1. Temp-artifact globs (gitmap-update-* in TEMP and deploy dirs)
//  2. Backup-artifact globs (*.old next to deployed binary)
//  3. Clone-swap dir globs (*.gitmap-tmp-*)
//  4. Drive-root shim candidate (Windows only)
func dumpDebugWindowsCleanupPlan(ctx updateCleanupContext) {
	if !isDebugWindowsRequested() {
		return
	}
	fmt.Fprint(os.Stderr, constants.MsgDebugWinCleanHdr)
	tempOps := dumpPlannedRemovals(ctx.tempPatterns, ctx.selfPath,
		constants.MsgDebugWinCleanMatch)
	backupOps := dumpPlannedRemovals(ctx.backupPatterns, ctx.selfPath,
		constants.MsgDebugWinCleanMatch)
	swapOps := dumpPlannedSwapDirs(ctx)
	shim, shimVerdict := dumpPlannedDriveRootShim(ctx)
	fmt.Fprint(os.Stderr, constants.MsgDebugWinCleanFooter)
	emitDebugWindowsJSON("cleanup_plan", map[string]any{
		"temp_removals":   tempOps,
		"backup_removals": backupOps,
		"swap_dirs":       swapOps,
		"drive_root_shim": map[string]any{
			"path": shim, "verdict": shimVerdict,
		},
	})
}

// dumpPlannedRemovals prints each glob pattern and the matches that
// would be passed to os.Remove. Skips matches equal to selfPath
// because removeCleanupMatch will skip them too.
func dumpPlannedRemovals(patterns []string, selfPath, matchFmt string) []map[string]any {
	ops := make([]map[string]any, 0, len(patterns))
	for _, pattern := range patterns {
		fmt.Fprintf(os.Stderr, constants.MsgDebugWinCleanGlob, pattern)
		matches, err := filepath.Glob(pattern)
		if err != nil || len(matches) == 0 {
			fmt.Fprint(os.Stderr, constants.MsgDebugWinCleanEmpty)
			ops = append(ops, map[string]any{"glob": pattern, "matches": []string{}})

			continue
		}
		printed := collectAndPrintMatches(matches, selfPath, matchFmt)
		if len(printed) == 0 {
			fmt.Fprint(os.Stderr, constants.MsgDebugWinCleanEmpty)
		}
		ops = append(ops, map[string]any{"glob": pattern, "matches": printed})
	}

	return ops
}

// collectAndPrintMatches filters the glob result the same way
// removeCleanupMatch does and prints + returns the survivors.
func collectAndPrintMatches(matches []string, selfPath, matchFmt string) []string {
	printed := make([]string, 0, len(matches))
	for _, m := range matches {
		if isActiveCleanupPath(normalizeCleanupPath(m), selfPath) {
			continue
		}
		fmt.Fprintf(os.Stderr, matchFmt, m)
		printed = append(printed, m)
	}

	return printed
}

// dumpPlannedSwapDirs prints every *.gitmap-tmp-* dir the cleanup
// pass will pass to os.RemoveAll. Mirrors removeCloneSwapDirsIn so
// the dump cannot drift from the live behavior. Returns the matched
// dirs so the JSON sink can record them too.
func dumpPlannedSwapDirs(ctx updateCleanupContext) []map[string]any {
	dirs := uniqueParentDirs(ctx.tempPatterns, ctx.backupPatterns)
	ops := make([]map[string]any, 0, len(dirs))
	for _, dir := range dirs {
		pattern := filepath.Join(dir, cloneSwapDirGlob)
		fmt.Fprintf(os.Stderr, constants.MsgDebugWinCleanGlob, pattern)
		matches, err := filepath.Glob(pattern)
		if err != nil || len(matches) == 0 {
			fmt.Fprint(os.Stderr, constants.MsgDebugWinCleanEmpty)
			ops = append(ops, map[string]any{"glob": pattern, "matches": []string{}})

			continue
		}
		ops = append(ops, map[string]any{"glob": pattern,
			"matches": collectSwapDirMatches(matches)})
	}

	return ops
}

// collectSwapDirMatches stat-filters glob hits to directories and
// prints + returns those that are real swap dirs.
func collectSwapDirMatches(matches []string) []string {
	out := make([]string, 0, len(matches))
	for _, m := range matches {
		info, statErr := os.Stat(m)
		if statErr != nil || !info.IsDir() {
			continue
		}
		fmt.Fprintf(os.Stderr, constants.MsgDebugWinCleanSwap, m)
		out = append(out, m)
	}

	return out
}

// dumpPlannedDriveRootShim reports the Windows-only drive-root shim
// candidate and whether cleanup would actually delete it. Returns
// (path, verdict) so the JSON sink can record the same fact.
func dumpPlannedDriveRootShim(ctx updateCleanupContext) (string, string) {
	shim := resolveDriveRootShimPath(ctx.selfPath)
	if len(shim) == 0 {
		return "", ""
	}
	verdict := constants.MsgDebugWinCleanShimSkip
	if isRemovableDriveRootShim(shim, ctx.selfPath) {
		verdict = constants.MsgDebugWinCleanShimDel
	}
	fmt.Fprintf(os.Stderr, constants.MsgDebugWinCleanShim, shim, verdict)

	return shim, verdict
}
