package cmd

import (
	"bytes"
	"errors"
	"os/exec"
	"runtime"
	"strings"
	"testing"
)

func fakeTrueRunner(name string, args ...string) *exec.Cmd {
	return exec.Command("true")
}

func fakeFalseRunner(name string, args ...string) *exec.Cmd {
	return exec.Command("false")
}

func assertCGBanners(t *testing.T, stderr string) {
	t.Helper()
	hasRunning := strings.Contains(stderr, "Installing coding guidelines")
	if !hasRunning {
		t.Fatalf("stderr missing running banner: %q", stderr)
	}
	hasDone := strings.Contains(stderr, "OK Coding guidelines")
	if !hasDone {
		t.Fatalf("stderr missing done banner: %q", stderr)
	}
}

// TestRunCodingGuidelinesInstall_SuccessViaFakeRunner verifies the
// dispatcher wires the OS-appropriate installer, streams stdio, and
// reports success without shelling out to the network. The injected
// Runner swaps every command for a no-op shell that exits 0.
func TestRunCodingGuidelinesInstall_SuccessViaFakeRunner(t *testing.T) {
	t.Parallel()
	isWindows := runtime.GOOS == "windows"
	if isWindows {
		t.Skip("dispatcher branch covered by unix path in CI")
	}

	var stdout, stderr bytes.Buffer
	opts := CodingGuidelinesOpts{Runner: fakeTrueRunner, Stdout: &stdout, Stderr: &stderr}
	if err := RunCodingGuidelinesInstall(opts); err != nil {
		t.Fatalf("expected success, got err=%v; stderr=%q", err, stderr.String())
	}
	assertCGBanners(t, stderr.String())
}

func assertCGMissingShell(t *testing.T, err error, stderr string) {
	t.Helper()
	isExpectedErr := errors.Is(err, ErrCGShellNotFound)
	if !isExpectedErr {
		t.Fatalf("expected ErrCGShellNotFound, got %v", err)
	}
	hasUnixFallback := strings.Contains(stderr, "curl -fsSL")
	hasWindowsFallback := strings.Contains(stderr, "irm ")
	if !hasUnixFallback && !hasWindowsFallback {
		t.Fatalf("stderr missing manual fallback recipe: %q", stderr)
	}
}

// TestRunCodingGuidelinesInstall_ShellMissing verifies the dispatcher
// returns ErrCGShellNotFound (and prints an actionable manual fallback)
// when the required shell is absent from PATH. Uses injected LookPath
// to avoid mutating process-wide PATH.
func TestRunCodingGuidelinesInstall_ShellMissing(t *testing.T) {
	t.Parallel()
	var stderr bytes.Buffer
	opts := CodingGuidelinesOpts{
		LookPath: func(file string) (string, error) {
			return "", exec.ErrNotFound
		},
		Stderr: &stderr,
	}
	err := RunCodingGuidelinesInstall(opts)
	assertCGMissingShell(t, err, stderr.String())
}

func assertCGErrorAndBanner(t *testing.T, err error, stderr string) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected non-nil error from failing installer")
	}
	hasContext := strings.Contains(err.Error(), "coding-guidelines install")
	if !hasContext {
		t.Fatalf("error missing context prefix: %v", err)
	}
	hasFailBanner := strings.Contains(stderr, "install failed")
	if !hasFailBanner {
		t.Fatalf("stderr missing failure banner: %q", stderr)
	}
}

// TestRunCodingGuidelinesInstall_ExitCodePropagates verifies non-zero
// installer exit codes bubble up wrapped so callers can errors.Is /
// unwrap them per the zero-swallow error policy.
func TestRunCodingGuidelinesInstall_ExitCodePropagates(t *testing.T) {
	t.Parallel()
	isWindows := runtime.GOOS == "windows"
	if isWindows {
		t.Skip("dispatcher branch covered by unix path in CI")
	}

	var stderr bytes.Buffer
	err := RunCodingGuidelinesInstall(CodingGuidelinesOpts{Runner: fakeFalseRunner, Stderr: &stderr})
	assertCGErrorAndBanner(t, err, stderr.String())
}

func TestPatchCGArithmeticIncrements(t *testing.T) {
	t.Parallel()

	in := "((WROTE_NEW++))\n((COPIED++))\n((count + 1))\n"
	got := patchCGArithmeticIncrements(in)
	want := "((WROTE_NEW+=1))\n((COPIED+=1))\n((count + 1))\n"
	isMismatch := got != want
	if isMismatch {
		t.Fatalf("patched script mismatch:\nwant %q\n got %q", want, got)
	}
}

func assertCGNotes(t *testing.T, stderr string) {
	t.Helper()
	for _, want := range []string{"Note: --no-commit set", "Note: --no-push set"} {
		hasNote := strings.Contains(stderr, want)
		if !hasNote {
			t.Fatalf("stderr missing %q: %q", want, stderr)
		}
	}
}

func TestCommitCodingGuidelinesNoCommitNoPushPrintsBothNotes(t *testing.T) {
	t.Parallel()

	var stdout, stderr bytes.Buffer
	err := CommitCodingGuidelines(CGCommitOpts{NoCommit: true, NoPush: true, Stdout: &stdout, Stderr: &stderr})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	assertCGNotes(t, stderr.String())
}

