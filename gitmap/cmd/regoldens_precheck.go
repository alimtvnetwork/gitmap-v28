package cmd

// Determinism pre-check pass for `gitmap regoldens --determinism`.
//
// Strategy: run `go test` with the per-test trigger ON but the
// allow-update gate var deliberately OFF. Any test using
// goldenguard.AllowUpdateAfterDeterminism will exercise its writer
// determinismRunCount times BEFORE consulting the gate. Three
// outcomes are possible per test:
//
//  1. Writer is deterministic, then AllowUpdate Fatalfs because the
//     gate is off. This is the EXPECTED success path — no fixture is
//     written, the test "fails" but only because of the gate.
//  2. Writer is non-deterministic — AssertWriterDeterministic
//     Fatalfs first, with a message containing the marker substring
//     "is non-deterministic" (centralized in goldenguard).
//  3. Test does not use AllowUpdateAfterDeterminism — the
//     determinism check is a no-op for that test (correct: nothing
//     to assert).
//
// We distinguish (1) from (2) by scanning combined go-test output
// for the marker. If found, pass 1 is NOT run and the CLI exits 1.
// Otherwise the pre-check is reported as passed and pass 1 proceeds.

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// runDeterminismPrecheck executes the trigger-only pass and returns
// true on success (proceed to pass 1). On non-determinism it logs
// the failure and exits 1; pass 1 is never reached.
func runDeterminismPrecheck(cfg regoldensFlags) {
	fmt.Fprint(os.Stderr, constants.MsgRegoldensPrecheckHeader)
	captured := runPrecheckGoTest(cfg)
	if precheckFoundNonDeterminism(captured) {
		fmt.Fprintln(os.Stderr, constants.ErrRegoldensPrecheckFailed)
		os.Exit(1)
	}
	fmt.Fprint(os.Stderr, constants.MsgRegoldensPrecheckPass)
}

// runPrecheckGoTest runs `go test` with trigger ON / allow OFF and
// streams output to the user's terminal while ALSO capturing it for
// post-hoc marker scanning. Returns the captured combined bytes.
func runPrecheckGoTest(cfg regoldensFlags) []byte {
	argv := goTestArgv(cfg)
	cmd := exec.Command(argv[0], argv[1:]...)
	var buf bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &buf)
	cmd.Stderr = io.MultiWriter(os.Stderr, &buf)
	cmd.Env = buildPrecheckEnv()
	_ = cmd.Run() // exit code is intentionally ignored; we read the buffer
	return buf.Bytes()
}

// buildPrecheckEnv returns the child env: parent env with both gate
// vars stripped, then the trigger var (only) re-added. Mirrors the
// gate-strip done by buildPassEnv for safety against leaked exports.
func buildPrecheckEnv() []string {
	out := stripGoldenGateVars(os.Environ())
	return append(out,
		goTestUpdateTriggerEnv+"="+goTestUpdateEnvValue,
	)
}

// precheckFoundNonDeterminism reports whether the captured output
// contains the goldenguard non-determinism marker. Substring match
// is sufficient — the marker is a long, distinctive phrase chosen
// for exactly this purpose.
func precheckFoundNonDeterminism(captured []byte) bool {
	return strings.Contains(string(captured), constants.RegoldensNonDetMarker)
}
