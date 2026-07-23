// Package cmd — Phase 3 of the self-update handoff chain.
//
// Phase 1 (update.go): active gitmap.exe → handoff copy (gitmap-update-<pid>.exe)
// Phase 2 (update.go): handoff copy runs build/deploy via run.ps1
// Phase 3 (this file): handoff copy spawns the freshly-deployed gitmap.exe
//
//	(a different file with no lock) detached, with a small delay, to
//	run `update-cleanup`. Only the deployed binary can safely remove
//	the still-locked handoff copy and the just-renamed *.exe.old.
//
// See spec/08-generic-update/06-cleanup.md and
// spec/03-general/02f-self-update-orchestration.md for the full sequence.
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/verbose"
)

// scheduleDeployedCleanupHandoff hands off cleanup work to the freshly
// deployed gitmap binary. The deployed binary lives at a different path
// than the running handoff copy, so it can delete both:
//   - the handoff copy (gitmap-update-<pid>.exe) — locked by us
//   - the *.exe.old backup — sometimes briefly held by AV/Explorer
//
// On Windows we launch the deployed binary directly in a hidden process
// and let `update-cleanup` sleep briefly before removal starts, so the
// current handoff process has time to exit and release its file lock.
// On Unix we just exec it inline since no lock conflicts exist.
//
// Best-effort cleanup remains non-fatal, but launch failures are printed
// to stderr and verbose logs so the user can see what went wrong.
func scheduleDeployedCleanupHandoff() {
	dumpDebugWindowsHeader("phase-3 handoff (update-runner)")
	defer dumpDebugWindowsFooter()
	fmt.Printf(constants.MsgUpdatePhase3LogFile, handoffLogPath())

	deployed, source := resolveDeployedBinaryPath()
	if len(deployed) == 0 {
		fmt.Fprint(os.Stderr, constants.ErrUpdatePhase3TargetMissing)
		logUpdatePhase3(constants.UpdatePhase3LogTargetMissing)
		logHandoffEvent("phase-3", "target_missing", nil)
		dumpDebugWindowsNote("target missing — no cleanup child will be spawned")

		return
	}

	self, err := os.Executable()
	if err == nil && normalizeCleanupPath(self) == normalizeCleanupPath(deployed) {
		// We *are* the deployed binary (Unix in-place update). Just
		// run cleanup directly — no handoff needed.
		logUpdatePhase3(constants.UpdatePhase3LogInline, deployed)
		logHandoffEvent("phase-3", "inline",
			map[string]string{"target": deployed, "source": source})
		dumpDebugWindowsNote("inline cleanup — self == deployed (%s)", deployed)
		runUpdateCleanup()

		return
	}

	if runtime.GOOS != constants.OSWindows {
		spawnDeployedCleanupUnix(deployed, source)

		return
	}

	spawnDeployedCleanupWindows(deployed, source)
}

// resolveDeployedBinaryPath returns the path to the freshly-deployed
// gitmap binary plus the resolution source label used in logs.
//
// Resolution order matters:
//  1. Config-declared deployed binary (powershell.json deployPath)
//  2. Sibling gitmap(.exe) next to the handoff copy
//  3. PATH lookup as a last resort only
//
// PATH is intentionally last because duplicate/stale gitmap.exe installs can
// linger on Windows and point cleanup at the wrong binary after an update.
func resolveDeployedBinaryPath() (string, string) {
	deployed, _ := resolveDeployedAndConfigPaths()
	if len(deployed) > 0 {
		return deployed, constants.UpdateCleanupSourceConfig
	}

	self, err := os.Executable()
	if err == nil {
		candidate := filepath.Join(filepath.Dir(self), deployedBinaryName())
		if _, err := os.Stat(candidate); err == nil {
			return candidate, constants.UpdateCleanupSourceSibling
		}
	}

	path, err := exec.LookPath(constants.GitMapBin)
	if err != nil {
		return "", constants.UpdateCleanupSourceUnknown
	}
	resolved, evalErr := filepath.EvalSymlinks(path)
	if evalErr == nil {
		return resolved, constants.UpdateCleanupSourcePath
	}

	return path, constants.UpdateCleanupSourcePath
}

// deployedBinaryName returns the platform-specific deployed binary filename.
func deployedBinaryName() string {
	if runtime.GOOS == constants.OSWindows {
		return constants.GitMapBin + ".exe"
	}

	return constants.GitMapBin
}

// spawnDeployedCleanupWindows launches the deployed binary directly in a
// detached hidden process. We avoid `cmd.exe /C start ...` entirely because
// its quoting rules are brittle when combined with Go's Windows argument
// escaping and can surface GUI popups like "Windows cannot find '\\'".
func spawnDeployedCleanupWindows(deployed, source string) {
	fmt.Printf(constants.MsgUpdatePhase3Handoff, filepath.Base(deployed))
	fmt.Printf(constants.MsgUpdatePhase3Resolve, source)
	fmt.Printf(constants.MsgUpdatePhase3Target, deployed)
	logUpdatePhase3(constants.UpdatePhase3LogResolve, source, deployed)
	logHandoffEvent("phase-3", "resolve",
		map[string]string{"source": source, "target": deployed})

	childArgs := buildCleanupChildArgs()
	cmd := exec.Command(deployed, childArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = nil
	setHiddenProcessAttr(cmd)
	cmd.Env = buildCleanupChildEnv()
	dumpDebugWindowsHandoff(source, deployed, append([]string{deployed}, childArgs...))
	dumpDebugWindowsCommandPlan(deployed, childArgs)
	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrUpdatePhase3Handoff, deployed, err)
		logUpdatePhase3(constants.UpdatePhase3LogStartFail, deployed, err)
		logHandoffEvent("phase-3", "start_fail", map[string]string{
			"target": deployed, "err": err.Error(),
		})

		return
	}

	fmt.Printf(constants.MsgUpdatePhase3Started, cmd.Process.Pid)
	logUpdatePhase3(constants.UpdatePhase3LogStarted, cmd.Process.Pid, deployed)
	logHandoffEvent("phase-3", "start_ok", map[string]string{
		"target": deployed,
		"pid":    fmt.Sprintf("%d", cmd.Process.Pid),
	})
	dumpDebugWindowsChildPID(cmd.Process.Pid)
}

// spawnDeployedCleanupUnix invokes the deployed binary's update-cleanup
// directly. No lock conflicts exist on Unix, so we don't need detachment.
func spawnDeployedCleanupUnix(deployed, source string) {
	fmt.Printf(constants.MsgUpdatePhase3Handoff, filepath.Base(deployed))
	fmt.Printf(constants.MsgUpdatePhase3Resolve, source)
	fmt.Printf(constants.MsgUpdatePhase3Target, deployed)
	logUpdatePhase3(constants.UpdatePhase3LogResolve, source, deployed)
	logHandoffEvent("phase-3", "resolve",
		map[string]string{"source": source, "target": deployed})

	childArgs := buildCleanupChildArgs()
	cmd := exec.Command(deployed, childArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = buildCleanupChildEnv()
	dumpDebugWindowsHandoff(source, deployed, append([]string{deployed}, childArgs...))
	dumpDebugWindowsCommandPlan(deployed, childArgs)
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrUpdatePhase3Handoff, deployed, err)
		logUpdatePhase3(constants.UpdatePhase3LogStartFail, deployed, err)
		logHandoffEvent("phase-3", "run_fail", map[string]string{
			"target": deployed, "err": err.Error(),
		})

		return
	}
	logHandoffEvent("phase-3", "run_ok",
		map[string]string{"target": deployed})
}

// buildCleanupChildArgs returns the argv for the detached `update-cleanup`
// child. --debug-windows and --debug-windows-json are forwarded so the
// child runs the same diagnostics as the parent.
func buildCleanupChildArgs() []string {
	args := []string{constants.CmdUpdateCleanup}
	if isDebugWindowsRequested() {
		args = append(args, constants.FlagDebugWindows)
	}
	if isDebugWindowsJSONRequested() {
		args = append(args, constants.FlagDebugWindowsJSON)
	}

	return args
}

// buildCleanupChildEnv returns the env for the detached cleanup child.
// EnvUpdateCleanupDelayMS is always set to give locks time to release.
// EnvDebugWindows and EnvDebugWindowsJSON are forwarded so the dump
// (and the JSON sink path) survive even if argv is dropped by an
// intermediate launcher.
func buildCleanupChildEnv() []string {
	env := append(os.Environ(), constants.EnvUpdateCleanupDelayMS+"=1500")
	if isDebugWindowsRequested() {
		env = append(env, constants.EnvDebugWindows+"=1")
	}
	if path := debugWindowsJSONPath(); len(path) > 0 {
		env = append(env, constants.EnvDebugWindowsJSON+"="+path)
	}

	return env
}

// logUpdatePhase3 writes handoff diagnostics to the shared verbose logger.
func logUpdatePhase3(format string, args ...interface{}) {
	log := verbose.Get()
	if log != nil {
		log.Log(format, args...)
	}
}
