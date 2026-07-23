package cmd

// clonestream_integration_test.go — end-to-end-style test that drives
// the production per-row print pipeline (printCloneFromTermBlockRow)
// the same way clone-from's executor does at runtime, captures the
// real os.Stdout and os.Stderr streams, and golden-compares both.
//
// What this test asserts that unit tests don't:
//
//  1. STREAMING ORDER: blocks for rows 1..N appear on stdout in the
//     same order the rows were submitted (catches a future regression
//     where buffering or goroutine fan-out reorders output).
//
//  2. STREAM SEPARATION: terminal blocks land on stdout ONLY; verifier
//     reports and --print-clone-argv dumps land on stderr ONLY. This
//     is the contract that makes `gitmap clone-from … | jq` and
//     similar pipelines work — a single byte of diagnostics leaking
//     to stdout would corrupt downstream parsers.
//
//  3. FLAG INDEPENDENCE: stdout bytes are byte-identical whether
//     --verify-cmd-faithful and --print-clone-argv are on or off.
//     Two scenarios (both flags on vs. both flags off) share the same
//     stdout golden to make this guarantee explicit in CI.
//
// Why drive printCloneFromTermBlockRow directly (vs. spawning a real
// `gitmap` process): a subprocess test would need git on PATH, network
// access to ls-remote a real repo, and a clean working directory —
// none of which the sandboxed CI runners can reliably provide. This
// in-process harness exercises the SAME function the executor's
// BeforeRow hook calls (clonefrom/execute_hooks.go invokes it with the
// same signature), so the streaming + separation contract is verified
// without external dependencies.
//
// Network avoidance: every test row sets row.Branch to a non-empty
// value so printCloneFromTermBlockRow's `len(branch) == 0` guard
// short-circuits — detectRemoteHEAD (which would shell out to
// `git ls-remote`) is never reached. Pinned by the assertion that
// stderr matches the golden EXACTLY, which leaves no room for stray
// network-error lines.

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/clonefrom"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/goldenguard"
)

// streamCaptureRows is the fixed 3-row input used by every scenario
// in this file. Realistic URLs + dest names so the golden reads as
// a plausible production transcript. Branches are pinned (non-empty)
// to avoid network calls — see file header.
func streamCaptureRows() []clonefrom.Row {

	return []clonefrom.Row{
		{URL: "https://github.com/acme/widget.git",
			Branch: "main", Depth: 1},
		{URL: "git@github.com:acme/gadget.git",
			Branch: "develop"},
		{URL: "https://github.com/acme/sprocket.git",
			Branch: "release/v2", Depth: 5},
	}
}

// streamCaptureDests pairs 1:1 with streamCaptureRows. Pulled out
// (vs. inlined) so the row table reads cleanly and the test loop
// can index both with a single i.
func streamCaptureDests() []string {

	return []string{"widget", "gadget", "sprocket"}
}

// TestCloneFromStream_Integration_FlagsOff is the baseline: no
// verifier, no argv dump. Stdout MUST contain 3 blocks (15 lines
// total: 5 per block); stderr MUST be empty (proving no diagnostic
// bytes leak when the flags are off — including no stray network
// errors from accidental ls-remote calls).
func TestCloneFromStream_Integration_FlagsOff(t *testing.T) {
	defer resetCloneStreamFlags(setCloneStreamFlags(false, false))
	out, errBuf := captureStreamedRows(t)
	assertStreamGolden(t, "clonestream_blocks_3rows.stdout.golden", out)
	assertStreamGolden(t, "clonestream_diag_empty.stderr.golden", errBuf)
}

// TestCloneFromStream_Integration_FlagsOn enables BOTH the verifier
// and the argv dump. Stdout golden is REUSED from the flags-off case
// — the separation contract requires byte-identical stdout regardless
// of stderr-side knobs. Stderr now carries 3 argv dumps (verifier
// stays silent because the displayed cmd: matches the real argv).
func TestCloneFromStream_Integration_FlagsOn(t *testing.T) {
	defer resetCloneStreamFlags(setCloneStreamFlags(true, true))
	out, errBuf := captureStreamedRows(t)
	assertStreamGolden(t, "clonestream_blocks_3rows.stdout.golden", out)
	assertStreamGolden(t, "clonestream_diag_argv_3rows.stderr.golden", errBuf)
}

// captureStreamedRows swaps os.Stdout and os.Stderr for pipes,
// streams the 3 fixed rows through printCloneFromTermBlockRow exactly
// as clonefrom.Execute's BeforeRow hook does, then returns each
// stream's bytes. Pipe-based capture (vs. a bytes.Buffer) is
// mandatory because the production code writes directly to os.Stdout
// — a Buffer swap wouldn't intercept those writes.
func captureStreamedRows(t *testing.T) ([]byte, []byte) {
	t.Helper()
	rows := streamCaptureRows()
	dests := streamCaptureDests()
	restore, outCh, errCh := redirectStdStreams(t)
	for i, row := range rows {
		printCloneFromTermBlockRow(i+1, len(rows), row, dests[i])
	}
	restore() // closes the pipe writers so the drain goroutines exit

	return <-outCh, <-errCh
}

// redirectStdStreams replaces os.Stdout/os.Stderr with pipe writers
// and spawns a drain goroutine per stream. Returns a restore closure
// that puts the originals back AND closes the writers so the drains
// finish. The two channels deliver the captured bytes once draining
// completes — read them AFTER calling restore.
func redirectStdStreams(t *testing.T) (func(), <-chan []byte, <-chan []byte) {
	t.Helper()
	origOut, origErr := os.Stdout, os.Stderr
	outR, outW := mustPipe(t)
	errR, errW := mustPipe(t)
	os.Stdout, os.Stderr = outW, errW
	outCh := drainPipe(outR)
	errCh := drainPipe(errR)
	restore := func() {
		_ = outW.Close()
		_ = errW.Close()
		os.Stdout, os.Stderr = origOut, origErr
	}

	return restore, outCh, errCh
}

// mustPipe wraps os.Pipe with a t.Fatalf failure path so callers
// get one-line "happy path" code without an err shuffle.
func mustPipe(t *testing.T) (*os.File, *os.File) {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}

	return r, w
}

// drainPipe reads r to EOF in a goroutine and delivers the bytes
// on the returned channel. EOF arrives when the corresponding
// writer is closed by restoreStdStreams.
func drainPipe(r *os.File) <-chan []byte {
	ch := make(chan []byte, 1)
	go func() {
		defer close(ch)
		b, _ := io.ReadAll(r)
		ch <- b
	}()

	return ch
}

// setCloneStreamFlags flips both request-scoped knobs and returns
// their PREVIOUS values so the caller can restore via
// resetCloneStreamFlags(defer). Two-value return keeps the test
// signature one defer line.
func setCloneStreamFlags(verify, argv bool) (bool, bool) {
	prevVerify := cmdFaithfulVerifyEnabled()
	prevArgv := cmdPrintArgvEnabled()
	setCmdFaithfulVerify(verify)
	setCmdPrintArgv(argv)

	return prevVerify, prevArgv
}

// resetCloneStreamFlags is the symmetric restore for
// setCloneStreamFlags. Splitting setter/restore lets the test use
// the standard `defer restore(setter())` pattern.
func resetCloneStreamFlags(prev0, prev1 bool) {
	setCmdFaithfulVerify(prev0)
	setCmdPrintArgv(prev1)
}

// assertStreamGolden mirrors the golden helpers used by the other
// cmd-package goldens but routes through goldenguard.AllowUpdate so
// regenerate is dual-gated (per-test trigger + GITMAP_ALLOW_GOLDEN_
// UPDATE=1). Per-test trigger here is GITMAP_UPDATE_GOLDEN=1 to
// match the existing convention used by jsoncontract_helpers_test.go.
func assertStreamGolden(t *testing.T, name string, got []byte) {
	t.Helper()
	path := filepath.Join("testdata", name)
	trigger := os.Getenv("GITMAP_UPDATE_GOLDEN") == "1"
	if goldenguard.AllowUpdate(t, trigger) {
		writeStreamGolden(t, path, got)

		return
	}
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read golden %s: %v (run with GITMAP_UPDATE_GOLDEN=1 "+
			"and GITMAP_ALLOW_GOLDEN_UPDATE=1 to create)", path, err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("golden mismatch for %s\n--- want (%d bytes) ---\n"+
			"%s\n--- got (%d bytes) ---\n%s",
			name, len(want), string(want), len(got), string(got))
	}
}

// writeStreamGolden persists a regenerated fixture and FAILS the
// test loudly so a CI run can never silently pass on a regenerate.
// Same contract as the other golden writers in this package.
func writeStreamGolden(t *testing.T, path string, got []byte) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir testdata: %v", err)
	}
	if err := os.WriteFile(path, got, 0o644); err != nil {
		t.Fatalf("write golden %s: %v", path, err)
	}
	t.Fatalf("regenerated golden %s — re-run without "+
		"GITMAP_UPDATE_GOLDEN to confirm", path)
}
