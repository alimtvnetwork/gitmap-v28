package cmd

// Two-pass orchestration for `gitmap regoldens`. Split from
// regoldens.go to respect the 200-line file cap and to keep each
// orchestration step under the 15-line function cap. The diff
// summary (when --diff is set) runs between pass 1 and pass 2 and
// fires regardless of pass-1 success so contributors can see what
// pass 1 wrote even when it failed mid-way.

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// executeRegoldens is the top-level orchestrator. Each step is a
// dedicated helper so the control flow reads top-to-bottom and the
// per-function line budget stays comfortably under 15.
func executeRegoldens(cfg regoldensFlags) {
	if cfg.determinism {
		runDeterminismPrecheck(cfg)
	}
	pass1Code := runRegoldensPassCapture(cfg, true, constants.MsgRegoldensPass1Header)
	maybeEmitDiffSummary(cfg)
	exitOnPass1Failure(cfg, pass1Code)
	if handleSkipVerify(cfg) {
		return
	}
	runPass2AndAnnounce(cfg)
}

// runRegoldensPassCapture prints the header and runs one pass,
// returning the exit code instead of exiting. Used so downstream
// work (diff summary) runs whether the pass succeeded or failed.
func runRegoldensPassCapture(cfg regoldensFlags, withGate bool, header string) int {
	fmt.Fprint(os.Stderr, header)
	return runGoTestPass(cfg, withGate)
}

// maybeEmitDiffSummary fires the post-pass-1 diff report when the
// user passed --diff=short|full. The mode is forwarded to the diff
// emitter so it can pick the per-line and totals format.
func maybeEmitDiffSummary(cfg regoldensFlags) {
	if cfg.hasDiff() {
		emitGoldenDiffSummary(cfg.diffMode)
	}
}

// exitOnPass1Failure logs the pass-1 error and exits 1 when pass 1
// returned non-zero. When --diff is enabled, also emits the explicit
// final line declaring that pass 2 did not run, so contributors get
// an unambiguous status even on the failure path.
func exitOnPass1Failure(cfg regoldensFlags, code int) {
	if code == 0 {
		return
	}
	fmt.Fprintf(os.Stderr, constants.ErrRegoldensPass1Failed, code)
	fmt.Fprintln(os.Stderr)
	if cfg.hasDiff() {
		fmt.Fprintf(os.Stderr, constants.MsgRegoldensPass2NotRun, code)
	}
	os.Exit(1)
}

// handleSkipVerify emits the skip-verify success path and returns
// true when --skip-verify is set. With --diff enabled it also prints
// the explicit pass-2-not-run final line for symmetry with the
// failure path.
func handleSkipVerify(cfg regoldensFlags) bool {
	if !cfg.skipVerify {
		return false
	}
	fmt.Fprint(os.Stderr, constants.MsgRegoldensSkipVerify)
	if cfg.hasDiff() {
		fmt.Fprint(os.Stderr, constants.MsgRegoldensPass2NotRunSkip)
	}
	fmt.Fprintf(os.Stdout, constants.MsgRegoldensSuccessNoVeri,
		cfg.pattern, cfg.pkg)
	return true
}

// runPass2AndAnnounce executes the determinism verification pass
// and prints the final success line. Pass 2 exits 1 internally on
// failure (via runRegoldensPass), so reaching the final Fprintf
// implies both passes succeeded. With --diff enabled it also emits
// the explicit "Pass 2 ran and PASSED" final line.
func runPass2AndAnnounce(cfg regoldensFlags) {
	runRegoldensPass(cfg, false,
		constants.MsgRegoldensPass2Header,
		constants.ErrRegoldensPass2Failed)
	if cfg.hasDiff() {
		fmt.Fprint(os.Stderr, constants.MsgRegoldensPass2Ran)
	}
	fmt.Fprintf(os.Stdout, constants.MsgRegoldensSuccess,
		cfg.pattern, cfg.pkg)
}
