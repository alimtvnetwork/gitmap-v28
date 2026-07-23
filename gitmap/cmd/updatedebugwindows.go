// Package cmd — `--debug-windows` diagnostics for the self-update
// Phase 3 cleanup handoff.
//
// The flag is opt-in and prints a structured dump to os.Stderr on
// every relevant lifecycle event. It propagates across the handoff
// boundary via two channels:
//
//  1. Argv — `--debug-windows` is forwarded into the handoff copy
//     (Phase 2) and the detached cleanup child (Phase 3).
//  2. Env  — `GITMAP_DEBUG_WINDOWS=1` is set on the cleanup child so
//     even processes spawned without an inherited argv (e.g. future
//     re-execs) keep printing the dump.
//
// Either signal alone activates the dump; users can flip the env var
// manually to enable the dump on a single run without rebuilding.
//
// The dump is intentionally cross-platform (works on Unix too) so the
// same flag can debug Linux/macOS handoffs, even though the original
// motivation was the Windows update-cleanup loop tracked in Issue #10.
package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/fsutil"
)

// isDebugWindowsRequested returns true when --debug-windows is on argv
// OR the GITMAP_DEBUG_WINDOWS env var is set to a truthy value. The env
// fallback is what carries the signal across the detached cleanup spawn.
func isDebugWindowsRequested() bool {
	for _, arg := range os.Args[1:] {
		if arg == constants.FlagDebugWindows {
			return true
		}
	}

	return isDebugWindowsEnvOn()
}

// isDebugWindowsEnvOn reports whether the env-bridge is enabled.
// Accepts the common truthy spellings so users can flip it manually.
func isDebugWindowsEnvOn() bool {
	v := os.Getenv(constants.EnvDebugWindows)
	switch v {
	case "1", "true", "TRUE", "yes", "YES", "on", "ON":
		return true
	}

	return false
}

// dumpDebugWindowsHeader prints the shared header lines (phase, GOOS,
// self executable, self pid, parent pid). Call from each lifecycle
// hook before the phase-specific lines.
func dumpDebugWindowsHeader(phase string) {
	if !isDebugWindowsRequested() {
		return
	}
	self, _ := os.Executable()
	fmt.Fprint(os.Stderr, constants.MsgDebugWinHeader)
	fmt.Fprintf(os.Stderr, constants.MsgDebugWinPhase, phase)
	fmt.Fprintf(os.Stderr, constants.MsgDebugWinGOOS, runtime.GOOS)
	fmt.Fprintf(os.Stderr, constants.MsgDebugWinSelf, self)
	fmt.Fprintf(os.Stderr, constants.MsgDebugWinPID, os.Getpid())
	fmt.Fprintf(os.Stderr, constants.MsgDebugWinPPID, os.Getppid())
	fmt.Fprintf(os.Stderr, constants.MsgDebugWinLogFile, handoffLogPath())
	emitDebugWindowsJSON("header", map[string]any{"phase": phase,
		"handoff_log": handoffLogPath()})
}

// dumpDebugWindowsFooter closes a dump block. Symmetric with
// dumpDebugWindowsHeader so the on/off state is identical.
func dumpDebugWindowsFooter() {
	if !isDebugWindowsRequested() {
		return
	}
	fmt.Fprint(os.Stderr, constants.MsgDebugWinFooter)
	emitDebugWindowsJSON("footer", nil)
}

// dumpDebugWindowsHandoff prints the resolution + child-launch details
// for the Phase 3 handoff. Called from spawnDeployedCleanup{Windows,Unix}
// before `cmd.Start()` so users can see what's about to happen even if
// the spawn fails immediately afterwards.
func dumpDebugWindowsHandoff(source, target string, childArgv []string) {
	if !isDebugWindowsRequested() {
		return
	}
	fmt.Fprintf(os.Stderr, constants.MsgDebugWinSource, source)
	fmt.Fprintf(os.Stderr, constants.MsgDebugWinTarget, target)
	fmt.Fprintf(os.Stderr, constants.MsgDebugWinTargetExists, fsutil.FileOrDirExists(target))
	fmt.Fprintf(os.Stderr, constants.MsgDebugWinChildArgv, childArgv)
	dumpDebugWindowsRelevantEnv()
	emitDebugWindowsJSON("handoff", map[string]any{
		"source": source, "target": target,
		"target_exists": fsutil.FileOrDirExists(target),
		"child_argv":    childArgv,
	})
}

// dumpDebugWindowsRelevantEnv prints the env vars that influence the
// update-cleanup handoff. Keep this list small and explicit — we never
// dump the full process environment because it can leak secrets.
func dumpDebugWindowsRelevantEnv() {
	keys := []string{
		constants.EnvDebugWindows,
		constants.EnvUpdateCleanupDelayMS,
		constants.EnvDebugRepoDetect,
		constants.EnvReportErrorsFormat,
		constants.EnvReportErrorsFile,
		"PATH",
		"GITMAP_DEPLOY_PATH",
	}
	for _, k := range keys {
		fmt.Fprintf(os.Stderr, constants.MsgDebugWinChildEnv,
			k, os.Getenv(k))
	}
}

// dumpDebugWindowsChildPID prints the spawned child's PID after a
// successful Start(). Skipped on failure (the start-fail message
// already carries the error).
func dumpDebugWindowsChildPID(pid int) {
	if !isDebugWindowsRequested() {
		return
	}
	fmt.Fprintf(os.Stderr, constants.MsgDebugWinChildPID, pid)
	emitDebugWindowsJSON("child_pid", map[string]any{"pid": pid})
}

// dumpDebugWindowsNote prints a freeform context line — used for
// "inline cleanup, no spawn needed" and "target missing" branches
// that don't have a child to report.
func dumpDebugWindowsNote(format string, args ...interface{}) {
	if !isDebugWindowsRequested() {
		return
	}
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(os.Stderr, constants.MsgDebugWinNote, msg)
	emitDebugWindowsJSON("note", map[string]any{"message": msg})
}

// (fileExistsLoose was removed in v3.113.0 — its contract now lives in
// fsutil.FileOrDirExists, which provides the same empty-short-circuit +
// directory-treated-as-existing semantics. Centralizing the predicate
// removes the redeclaration footgun that bit this package twice; the
// rename pin in updatedebugwindows_rename_test.go enforces the next
// step of the migration.)
