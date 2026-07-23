package cmd

// Locks the --debug-paths trace contract:
//
//   1. canonicalizePMPath is silent on stderr when GITMAP_DEBUG_PATHS
//      is unset (production default — zero diagnostic noise).
//   2. When GITMAP_DEBUG_PATHS=1, every canonicalize call emits one
//      line matching constants.MsgDebugPathsTrace with raw / cleaned
//      / resolved fields, on BOTH the resolved-ok and soft-fail paths.
//   3. applyDebugPathsEnv only sets the env var when the flag is on,
//      and never clears a pre-existing value (so CI presets survive).
//
// Stderr is captured via os.Pipe so the assertions stay hermetic and
// don't depend on log redirection. t.Setenv handles teardown.

import (
	"os"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

func TestCanonicalizePMPathSilentWhenDebugPathsOff(t *testing.T) {
	t.Setenv(constants.EnvDebugPaths, "")

	dir := t.TempDir()

	out := captureStderr(t, func() {
		_ = canonicalizePMPath(dir)
	})

	if out != "" {
		t.Fatalf("expected no stderr output when %s is unset, got %q",
			constants.EnvDebugPaths, out)
	}
}

func TestCanonicalizePMPathTracesWhenDebugPathsOn(t *testing.T) {
	t.Setenv(constants.EnvDebugPaths, constants.EnvDebugPathsOn)

	dir := t.TempDir()

	out := captureStderr(t, func() {
		_ = canonicalizePMPath(dir)
	})

	if !strings.Contains(out, "[debug-paths]") {
		t.Fatalf("expected [debug-paths] prefix in stderr, got %q", out)
	}

	if !strings.Contains(out, "in=") || !strings.Contains(out, "clean=") ||
		!strings.Contains(out, "resolved=") {
		t.Fatalf("expected in=/clean=/resolved= fields, got %q", out)
	}
}

func TestCanonicalizePMPathTracesOnSoftFail(t *testing.T) {
	t.Setenv(constants.EnvDebugPaths, constants.EnvDebugPathsOn)

	missing := "/nonexistent-gitmap-debug-paths-test-path-xyz"

	out := captureStderr(t, func() {
		_ = canonicalizePMPath(missing)
	})

	if !strings.Contains(out, "[debug-paths]") {
		t.Fatalf("expected trace on soft-fail path, got %q", out)
	}
}

func TestApplyDebugPathsEnvSetsVarWhenOn(t *testing.T) {
	t.Setenv(constants.EnvDebugPaths, "")

	applyDebugPathsEnv(true)

	if got := os.Getenv(constants.EnvDebugPaths); got != constants.EnvDebugPathsOn {
		t.Fatalf("expected env=%q, got %q", constants.EnvDebugPathsOn, got)
	}
}

func TestApplyDebugPathsEnvPreservesPresetWhenOff(t *testing.T) {
	t.Setenv(constants.EnvDebugPaths, constants.EnvDebugPathsOn)

	applyDebugPathsEnv(false)

	if got := os.Getenv(constants.EnvDebugPaths); got != constants.EnvDebugPathsOn {
		t.Fatalf("expected pre-set env to survive, got %q", got)
	}
}
