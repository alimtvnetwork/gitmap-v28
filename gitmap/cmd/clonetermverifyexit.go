package cmd

// clonetermverifyexit.go — implements the run-tail exit hook for
// --verify-cmd-faithful-exit-on-mismatch.
//
// Why a separate file: the state file (clonetermverifystate.go) holds
// the package-level toggles + the per-row check, and is approaching
// the project's 200-line ceiling. Splitting the "decide whether to
// exit" decision into its own file keeps both files small and gives
// the exit policy an obvious search target.
//
// Stream + exit-code contract (mirrors PrintCmdFaithfulReport's
// stderr stream choice — diagnostics never pollute the per-repo
// terminal block on stdout):
//
//   - exit code: constants.CloneVerifyCmdFaithfulExitCode (3) so CI
//     scripts can disambiguate from 1 (runtime / git failure) and 2
//     (bad CLI usage).
//   - stderr summary line: a one-line "FAIL" notice so the failure is
//     obvious in CI logs that surface the LAST stderr line as the
//     failure reason, even when the per-row mismatch reports scrolled
//     off the top.
//
// Timing: every clone dispatcher MUST call this AFTER its work is
// done so the FULL list of mismatches is printed before exit. Calling
// it earlier would short-circuit later rows and hide divergences.

import (
	"fmt"
	"os"
	"sync"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

var (
	cmdFaithfulExiterMu sync.RWMutex
	cmdFaithfulExiter   func(int) = os.Exit
)

func getCmdFaithfulExiter() func(int) {
	cmdFaithfulExiterMu.RLock()
	defer cmdFaithfulExiterMu.RUnlock()

	return cmdFaithfulExiter
}

func setCmdFaithfulExiter(fn func(int)) {
	cmdFaithfulExiterMu.Lock()
	defer cmdFaithfulExiterMu.Unlock()

	cmdFaithfulExiter = fn
}

// maybeExitOnCmdFaithfulMismatch is the single integration point for
// the exit-on-mismatch policy. Called from every clone dispatcher's
// terminal path (runClone, runCloneNow, runCloneFrom, runClonePick,
// runCloneNext) AFTER all per-row work is finished.
//
// No-op when either:
//   - the user did not pass --verify-cmd-faithful-exit-on-mismatch, or
//   - the verifier ran but found no divergences.
//
// Otherwise it prints a one-line stderr summary and exits with
// constants.CloneVerifyCmdFaithfulExitCode. The summary is a single
// fmt.Fprintln so it appears as the LAST stderr line — many CI UIs
// surface that line as the headline failure reason.
func maybeExitOnCmdFaithfulMismatch() {
	if !cmdFaithfulExitOnMismatchEnabled() {
		return
	}
	if !cmdFaithfulHadMismatchSet() {
		return
	}
	fmt.Fprintln(os.Stderr, constants.MsgCloneVerifyCmdFaithfulExit)
	getCmdFaithfulExiter()(constants.CloneVerifyCmdFaithfulExitCode)
}
