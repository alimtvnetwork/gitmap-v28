package cmd

// clonetermverifyexit_test.go — exercises the run-tail exit policy
// added by --verify-cmd-faithful-exit-on-mismatch.
//
// We swap cmdFaithfulExiter for a recorder so the test can assert
// on the exit code without terminating the test binary, then drive
// runCmdFaithfulCheck through both paths (match / mismatch) under
// every relevant combination of the verify + exit-on-mismatch flags.
//
// resetCmdFaithfulState is called at the top of every test (and in a
// t.Cleanup) so the package-level globals from clonetermverifystate.go
// don't leak across cases — Go test ordering is not guaranteed and a
// stale sticky bit would silently flip an unrelated test from green
// to red.

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// withRecordingExiter installs a stub that captures the exit code
// instead of terminating the process. Returns a pointer the caller
// reads after triggering the run-tail; restoration runs via t.Cleanup
// so a failing assertion can't leave the global pointing at the stub.
func withRecordingExiter(t *testing.T) *int {
	t.Helper()
	prev := cmdFaithfulExiter
	captured := -1
	cmdFaithfulExiter = func(code int) { captured = code }
	t.Cleanup(func() { cmdFaithfulExiter = prev })

	// Reset state inside the helper too so every call site gets a
	// clean slate without remembering a separate setup line.
	resetCmdFaithfulState()
	t.Cleanup(resetCmdFaithfulState)

	return &captured
}

// fixtureMismatchInput builds a CloneTermBlockInput / executor argv
// pair that's guaranteed to diverge at index 3 (branch token). Reused
// across cases so the tests stay focused on policy, not on contriving
// a divergent input each time.
func fixtureMismatchInput() (CloneTermBlockInput, []string) {
	in := CloneTermBlockInput{
		Name:        "scripts-fixer",
		OriginalURL: "https://example.com/scripts-fixer.git",
		TargetURL:   "https://example.com/scripts-fixer.git",
		Branch:      "main",
		CmdBranch:   "main",
		Dest:        "scripts-fixer",
	}
	executorArgv := []string{"clone", "-b", "develop",
		"https://example.com/scripts-fixer.git", "scripts-fixer"}

	return in, executorArgv
}

// TestMaybeExitOnCmdFaithfulMismatch_DefaultIsNoOp guarantees that the
// new flag is opt-in: when neither toggle is set, the run tail does
// NOT exit even if a mismatch is recorded externally. This protects
// callers (clone, clone-now, …) that always invoke the helper.
func TestMaybeExitOnCmdFaithfulMismatch_DefaultIsNoOp(t *testing.T) {
	captured := withRecordingExiter(t)

	// Simulate a stale mismatch bit being set by some earlier code path.
	cmdFaithfulHadMismatch.Store(true)

	maybeExitOnCmdFaithfulMismatch()

	if *captured != -1 {
		t.Fatalf("expected no exit when flag is off, got code %d", *captured)
	}
}

// TestMaybeExitOnCmdFaithfulMismatch_NoMismatchNoExit guarantees that
// enabling the flag alone is harmless: if the verifier ran and found
// nothing, the run tail still exits 0 (i.e. doesn't call the exiter).
func TestMaybeExitOnCmdFaithfulMismatch_NoMismatchNoExit(t *testing.T) {
	captured := withRecordingExiter(t)
	setCmdFaithfulExitOnMismatch(true)

	maybeExitOnCmdFaithfulMismatch()

	if *captured != -1 {
		t.Fatalf("expected no exit on clean run, got code %d", *captured)
	}
}

// TestMaybeExitOnCmdFaithfulMismatch_FiresOnMismatch is the main
// success path: when the user opts in AND a divergence is recorded,
// the helper exits with constants.CloneVerifyCmdFaithfulExitCode.
func TestMaybeExitOnCmdFaithfulMismatch_FiresOnMismatch(t *testing.T) {
	captured := withRecordingExiter(t)
	setCmdFaithfulExitOnMismatch(true)

	in, argv := fixtureMismatchInput()
	runCmdFaithfulCheck(in, argv) // records the mismatch sticky bit
	maybeExitOnCmdFaithfulMismatch()

	if *captured != constants.CloneVerifyCmdFaithfulExitCode {
		t.Fatalf("expected exit code %d, got %d",
			constants.CloneVerifyCmdFaithfulExitCode, *captured)
	}
}

// TestSetCmdFaithfulExitOnMismatch_ImpliesVerify documents the
// dependency: passing only --verify-cmd-faithful-exit-on-mismatch
// (without --verify-cmd-faithful) MUST still run the verifier,
// otherwise the sticky bit can never flip and the exit code is dead.
func TestSetCmdFaithfulExitOnMismatch_ImpliesVerify(t *testing.T) {
	captured := withRecordingExiter(t)

	// Set ONLY the exit-on-mismatch toggle; do not call setCmdFaithfulVerify.
	setCmdFaithfulExitOnMismatch(true)
	if !cmdFaithfulVerifyEnabled() {
		t.Fatal("setCmdFaithfulExitOnMismatch(true) must imply verify-enabled")
	}

	in, argv := fixtureMismatchInput()
	runCmdFaithfulCheck(in, argv)
	maybeExitOnCmdFaithfulMismatch()

	if *captured != constants.CloneVerifyCmdFaithfulExitCode {
		t.Fatalf("expected exit code %d, got %d",
			constants.CloneVerifyCmdFaithfulExitCode, *captured)
	}
}
