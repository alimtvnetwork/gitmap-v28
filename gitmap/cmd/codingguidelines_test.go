package cmd

import (
	"bytes"
	"errors"
	"os/exec"
	"runtime"
	"strings"
	"testing"
)

// TestRunCodingGuidelinesInstall_SuccessViaFakeRunner verifies the
// dispatcher wires the OS-appropriate installer, streams stdio, and
// reports success without shelling out to the network. The injected
// Runner swaps every command for a no-op shell that exits 0.
func TestRunCodingGuidelinesInstall_SuccessViaFakeRunner(t *testing.T) {
	t.Parallel()

	if runtime.GOOS == "windows" {
		// Skip on Windows: dispatchCGWindows requires a real PowerShell
		// binary discovered via resolvePowerShellBinary; the Unix path
		// gives us equivalent coverage for the runner contract.
		t.Skip("dispatcher branch covered by unix path in CI")
	}

	var stdout, stderr bytes.Buffer
	fake := func(name string, args ...string) *exec.Cmd {
		// Ignore the real (bash / pwsh) command and args: exercise the
		// wiring, not the network install. `true` exits 0 on every Unix.
		return exec.Command("true")
	}

	err := RunCodingGuidelinesInstall(CodingGuidelinesOpts{
		Runner: fake,
		Stdout: &stdout,
		Stderr: &stderr,
	})
	if err != nil {
		t.Fatalf("expected success, got err=%v; stderr=%q", err, stderr.String())
	}
	if !strings.Contains(stderr.String(), "Installing coding guidelines") {
		t.Fatalf("stderr missing running banner: %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "OK Coding guidelines") {
		t.Fatalf("stderr missing done banner: %q", stderr.String())
	}
}

// TestRunCodingGuidelinesInstall_ShellMissing verifies the dispatcher
// returns ErrCGShellNotFound (and prints an actionable manual fallback)
// when the required shell is absent from PATH. Simulated by emptying
// PATH for the duration of the test.
func TestRunCodingGuidelinesInstall_ShellMissing(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("PATH manipulation for pwsh discovery is platform-specific")
	}
	t.Setenv("PATH", "")

	var stderr bytes.Buffer
	err := RunCodingGuidelinesInstall(CodingGuidelinesOpts{Stderr: &stderr})
	if !errors.Is(err, ErrCGShellNotFound) {
		t.Fatalf("expected ErrCGShellNotFound, got %v", err)
	}
	if !strings.Contains(stderr.String(), "curl -fsSL") {
		t.Fatalf("stderr missing manual fallback recipe: %q", stderr.String())
	}
}

// TestRunCodingGuidelinesInstall_ExitCodePropagates verifies non-zero
// installer exit codes bubble up wrapped so callers can errors.Is /
// unwrap them per the zero-swallow error policy.
func TestRunCodingGuidelinesInstall_ExitCodePropagates(t *testing.T) {
	t.Parallel()
	if runtime.GOOS == "windows" {
		t.Skip("dispatcher branch covered by unix path in CI")
	}

	fake := func(name string, args ...string) *exec.Cmd {
		return exec.Command("false") // exits 1 on every Unix
	}

	var stderr bytes.Buffer
	err := RunCodingGuidelinesInstall(CodingGuidelinesOpts{Runner: fake, Stderr: &stderr})
	if err == nil {
		t.Fatalf("expected non-nil error from failing installer")
	}
	if !strings.Contains(err.Error(), "coding-guidelines install") {
		t.Fatalf("error missing context prefix: %v", err)
	}
	if !strings.Contains(stderr.String(), "install failed") {
		t.Fatalf("stderr missing failure banner: %q", stderr.String())
	}
}

func TestPatchCGArithmeticIncrements(t *testing.T) {
	t.Parallel()

	in := "((WROTE_NEW++))\n((COPIED++))\n((count + 1))\n"
	got := patchCGArithmeticIncrements(in)
	want := "((WROTE_NEW+=1))\n((COPIED+=1))\n((count + 1))\n"
	if got != want {
		t.Fatalf("patched script mismatch:\nwant %q\n got %q", want, got)
	}
}

func TestCommitCodingGuidelinesNoCommitNoPushPrintsBothNotes(t *testing.T) {
	t.Parallel()

	var stderr bytes.Buffer
	err := CommitCodingGuidelines(CGCommitOpts{NoCommit: true, NoPush: true, Stderr: &stderr})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	for _, want := range []string{"Note: --no-commit set", "Note: --no-push set"} {
		if !strings.Contains(stderr.String(), want) {
			t.Fatalf("stderr missing %q: %q", want, stderr.String())
		}
	}
}
