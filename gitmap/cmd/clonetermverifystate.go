package cmd

// clonetermverifystate.go — request-scoped knob that turns on the
// --verify-cmd-faithful checker for the current command invocation.
//
// Why a package-level variable (vs. plumbing a bool through every
// hook signature): the existing print-row helpers
// (printCloneNowTermBlockRow / printCloneFromTermBlockRow /
// printCloneTermBlockForURL) are reached from 5+ call sites with
// fixed signatures dictated by clonenow.BeforeRowHook /
// clonefrom.BeforeRowHook. Adding a parameter to each signature
// would force a churning change in two executor packages plus every
// dry-run path. A package-level toggle set ONCE per `gitmap …`
// invocation (CLI is single-threaded at the dispatcher level) keeps
// the executor packages oblivious to verification AND keeps the
// hooks' signatures stable for tests/contracts.
//
// Concurrency note: clone batch modes use a worker pool, but each
// row-hook reads (never writes) cmdFaithfulVerify; the variable is
// set before Execute starts and never mutated again. Reads under a
// fixed-after-startup invariant don't need synchronization (Go's
// happens-before guarantees the dispatcher's set happens-before the
// goroutines' reads via the goroutine launch).

import (
	"os"
	"sync/atomic"
)

// cmdFaithfulVerifyEnabled is the request-scoped flag. The package's
// CLI dispatchers (runClone / runCloneNow / runCloneFrom / runClonePick
// / runCloneNext) flip it to true via setCmdFaithfulVerify when
// --verify-cmd-faithful is parsed, and the print-row helpers consult
// it via cmdFaithfulVerifyEnabled() before running the checker.
//
// Default false so existing behavior is byte-identical when the flag
// is absent.
var cmdFaithfulVerify bool

// cmdFaithfulExitOnMismatch toggles the hard-fail companion behavior
// of --verify-cmd-faithful-exit-on-mismatch. When true, ANY mismatch
// recorded by runCmdFaithfulCheck flips cmdFaithfulHadMismatch, and
// the per-command run tail calls maybeExitOnCmdFaithfulMismatch which
// terminates with constants.CloneVerifyCmdFaithfulExitCode.
//
// Independent of cmdFaithfulVerify in storage so flag-parse code can
// set both without ordering caveats; setCmdFaithfulExitOnMismatch
// enforces the implied dependency by also flipping cmdFaithfulVerify.
var cmdFaithfulExitOnMismatch bool

// cmdFaithfulHadMismatch is the run-scoped sticky bit set by every
// per-row check that finds at least one divergence. Atomic because
// batch clone modes drive concurrent goroutines through the row hook
// (clonefrom/clonenow worker pools) — a plain bool would be a data
// race even though the WRITE is monotonic (false → true).
//
// Read via cmdFaithfulHadMismatchSet at the END of the dispatcher to
// decide whether maybeExitOnCmdFaithfulMismatch should exit non-zero.
var cmdFaithfulHadMismatch atomic.Bool

// setCmdFaithfulVerify enables (or disables) the verifier for the
// remainder of the current process. Safe to call multiple times —
// last write wins, which matches the "set once at dispatcher" usage.
func setCmdFaithfulVerify(on bool) { cmdFaithfulVerify = on }

// setCmdFaithfulExitOnMismatch enables (or disables) the hard-fail
// companion. Implies --verify-cmd-faithful: when on we ALSO flip the
// verifier so users opting into the exit code don't have to re-type
// the prerequisite flag in CI invocations.
func setCmdFaithfulExitOnMismatch(on bool) {
	cmdFaithfulExitOnMismatch = on
	if on {
		cmdFaithfulVerify = true
	}
}

// cmdFaithfulVerifyEnabled returns true when the verifier should run.
// Predicate (vs. exposing the var) so a future move to atomic.Bool
// or a context-bound state stays a one-line refactor.
func cmdFaithfulVerifyEnabled() bool { return cmdFaithfulVerify }

// cmdFaithfulExitOnMismatchEnabled mirrors cmdFaithfulVerifyEnabled
// for the exit-on-mismatch toggle. Used by maybeExitOnCmdFaithfulMismatch.
func cmdFaithfulExitOnMismatchEnabled() bool { return cmdFaithfulExitOnMismatch }

// cmdFaithfulHadMismatchSet reports whether ANY per-row check has
// recorded a divergence since process start.
func cmdFaithfulHadMismatchSet() bool { return cmdFaithfulHadMismatch.Load() }

// resetCmdFaithfulState zeroes every verifier knob. Tests use this to
// keep the package-level globals from leaking across test cases; the
// production CLI never calls it (one process = one invocation).
func resetCmdFaithfulState() {
	cmdFaithfulVerify = false
	cmdFaithfulExitOnMismatch = false
	cmdFaithfulHadMismatch.Store(false)
}

// runCmdFaithfulCheck is the single integration point used by every
// per-row print helper. No-op when the flag is off so callers can
// invoke it unconditionally on the hot path.
//
// On mismatch it prints a structured report to stderr AND flips the
// sticky cmdFaithfulHadMismatch bit so the run tail can decide to
// exit non-zero when --verify-cmd-faithful-exit-on-mismatch is set.
// It does NOT abort here — every row gets a chance to print its
// block before the run-end exit fires.
func runCmdFaithfulCheck(in CloneTermBlockInput, executorArgv []string) {
	if !cmdFaithfulVerifyEnabled() {
		return
	}
	report := VerifyCmdFaithful(in, executorArgv)
	if report.HasMismatch() {
		cmdFaithfulHadMismatch.Store(true)
	}
	if err := PrintCmdFaithfulReport(os.Stderr, report); err != nil {
		// Zero-swallow policy: surface the write failure but don't
		// abort the clone — the verifier is purely informational.
		_, _ = os.Stderr.WriteString(
			"  Warning: --verify-cmd-faithful: failed to write report: " +
				err.Error() + "\n")
	}
}
