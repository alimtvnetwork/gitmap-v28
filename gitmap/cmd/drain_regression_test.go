package cmd

// Regression test for the Windows drain fix (v6.74.0).
//
// Bug: when the theme/glyphs pipe wrappers were installed (safe mode),
// the last stdout line written just before the process exited could be
// lost on Windows because the forwarder goroutine was descheduled
// before the runtime flushed the pipe. The CI installer smoke test
// caught this only after every release cut. This test reproduces the
// failure locally without any release dependency by forcing
// GITMAP_GLYPHS=safe (which activates the pipe wrapper) and asserting
// that `gitmap version` still prints its expected single-line output
// on a clean exit.
//
// Regression guard: if a future change removes the deferred Drain
// calls from runDispatch — or introduces a new entry point that calls
// dispatch directly — this test fails on Windows CI with empty stdout.

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// TestVersion_FlushedOnCleanExit runs `gitmap version` under
// GITMAP_GLYPHS=safe so the pipe wrappers are actually installed
// (hermeticEnv() pins them off for other tests). Asserts stdout is
// non-empty and contains the pinned constants.Version.
func TestVersion_FlushedOnCleanExit(t *testing.T) {
	t.Parallel()

	bin := ensureGitmapBinary(t)
	tmp := t.TempDir()
	stdoutPath := filepath.Join(tmp, "stdout")
	stderrPath := filepath.Join(tmp, "stderr")

	stdoutF, err := os.Create(stdoutPath)
	if err != nil {
		t.Fatalf("create stdout capture: %v", err)
	}
	stderrF, err := os.Create(stderrPath)
	if err != nil {
		t.Fatalf("create stderr capture: %v", err)
	}

	cmd := exec.Command(bin, "version")
	cmd.Dir = t.TempDir()
	cmd.Env = envSwapGlyphs("safe")
	cmd.Stdout = stdoutF
	cmd.Stderr = stderrF

	runErr := cmd.Run()
	stdoutF.Close()
	stderrF.Close()

	stdout, _ := os.ReadFile(stdoutPath)
	stderr, _ := os.ReadFile(stderrPath)
	code := extractTestExitCode(runErr)

	if code != 0 {
		t.Fatalf("version exited %d\nstdout=%q\nstderr=%q", code, stdout, stderr)
	}
	if len(strings.TrimSpace(string(stdout))) == 0 {
		t.Fatalf("version produced empty stdout under glyphs=safe — drain regression\nstderr=%q", stderr)
	}
	if !strings.Contains(string(stdout), constants.Version) {
		t.Fatalf("version stdout %q missing pinned version %q", stdout, constants.Version)
	}
}

// envSwapGlyphs returns hermeticEnv() with GITMAP_GLYPHS forced to
// mode (replacing the pinned "rich" default). Keeps the swap logic
// in one place so future glyph-mode tests are cheap to add.
func envSwapGlyphs(mode string) []string {
	base := hermeticEnv()
	out := make([]string, 0, len(base)+1)
	for _, kv := range base {
		if strings.HasPrefix(kv, "GITMAP_GLYPHS=") {
			continue
		}
		out = append(out, kv)
	}
	out = append(out, "GITMAP_GLYPHS="+mode)

	return out
}
