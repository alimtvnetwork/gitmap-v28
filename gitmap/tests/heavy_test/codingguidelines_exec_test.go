package heavy_test

import (
	"bytes"
	"errors"
	"os/exec"
	"runtime"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cmd"
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

func TestRunCodingGuidelinesInstall_SuccessViaFakeRunner(t *testing.T) {
	t.Parallel()
	isWindows := runtime.GOOS == "windows"
	if isWindows {
		t.Skip("dispatcher branch covered by unix path in CI")
	}

	var stdout, stderr bytes.Buffer
	opts := cmd.CodingGuidelinesOpts{Runner: fakeTrueRunner, Stdout: &stdout, Stderr: &stderr}
	if err := cmd.RunCodingGuidelinesInstall(opts); err != nil {
		t.Fatalf("expected success, got err=%v; stderr=%q", err, stderr.String())
	}

	assertCGBanners(t, stderr.String())
}

func assertCGMissingShell(t *testing.T, err error, stderr string) {
	t.Helper()
	isExpectedErr := errors.Is(err, cmd.ErrCGShellNotFound)
	if !isExpectedErr {
		t.Fatalf("expected ErrCGShellNotFound, got %v", err)
	}

	hasUnixFallback := strings.Contains(stderr, "curl -fsSL")
	hasWindowsFallback := strings.Contains(stderr, "irm ")
	if !hasUnixFallback && !hasWindowsFallback {
		t.Fatalf("stderr missing manual fallback recipe: %q", stderr)
	}
}

func TestRunCodingGuidelinesInstall_ShellMissing(t *testing.T) {
	t.Parallel()
	var stderr bytes.Buffer
	opts := cmd.CodingGuidelinesOpts{
		LookPath: func(file string) (string, error) {
			return "", exec.ErrNotFound
		},
		Stderr: &stderr,
	}

	err := cmd.RunCodingGuidelinesInstall(opts)
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

func TestRunCodingGuidelinesInstall_ExitCodePropagates(t *testing.T) {
	t.Parallel()
	isWindows := runtime.GOOS == "windows"
	if isWindows {
		t.Skip("dispatcher branch covered by unix path in CI")
	}

	var stderr bytes.Buffer
	err := cmd.RunCodingGuidelinesInstall(cmd.CodingGuidelinesOpts{Runner: fakeFalseRunner, Stderr: &stderr})
	assertCGErrorAndBanner(t, err, stderr.String())
}
