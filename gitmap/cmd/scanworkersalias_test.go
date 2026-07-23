package cmd

import (
	"flag"
	"io"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// newScanWorkersFS builds a throwaway flag set wired to the same
// names parseScanFlags uses, parses argv, and returns the pointers
// + flag set so each table-driven case can call resolveScanWorkers
// directly and inspect the result.
func newScanWorkersFS(t *testing.T, argv []string) (*flag.FlagSet, *int, *int) {
	t.Helper()
	fs := flag.NewFlagSet("scan-test", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	workers := fs.Int(constants.FlagScanWorkers,
		constants.DefaultScanWorkers, constants.FlagDescScanWorkers)
	conc := fs.Int(constants.FlagScanWorkersConcurrencyAlias,
		constants.DefaultScanWorkers,
		constants.FlagDescScanWorkersConcurrencyAlias)
	if err := fs.Parse(argv); err != nil {
		t.Fatalf("parse %v: %v", argv, err)
	}

	return fs, workers, conc
}

// TestResolveScanWorkers_Canonical asserts that --workers is
// honored verbatim and emits no deprecation notice.
func TestResolveScanWorkers_Canonical(t *testing.T) {
	fs, w, c := newScanWorkersFS(t, []string{"--workers", "8"})
	var got int
	stderr := captureStderr(t, func() { got = resolveScanWorkers(fs, w, c) })
	if got != 8 {
		t.Fatalf("workers: got %d, want 8", got)
	}
	if stderr != "" {
		t.Fatalf("expected silent stderr, got %q", stderr)
	}
}

// TestResolveScanWorkers_DeprecatedAlias asserts that --concurrency
// is honored AND emits the documented one-line deprecation notice.
func TestResolveScanWorkers_DeprecatedAlias(t *testing.T) {
	fs, w, c := newScanWorkersFS(t, []string{"--concurrency", "12"})
	var got int
	stderr := captureStderr(t, func() { got = resolveScanWorkers(fs, w, c) })
	if got != 12 {
		t.Fatalf("workers: got %d, want 12 (from --concurrency)", got)
	}
	if !strings.Contains(stderr, "deprecated") ||
		!strings.Contains(stderr, "--workers") {
		t.Fatalf("expected deprecation notice mentioning --workers, got %q",
			stderr)
	}
}

// TestResolveScanWorkers_CanonicalWinsOverAlias asserts that when
// both flags are passed, --workers wins and NO deprecation notice
// fires (the user is clearly already on the canonical spelling).
func TestResolveScanWorkers_CanonicalWinsOverAlias(t *testing.T) {
	fs, w, c := newScanWorkersFS(t,
		[]string{"--workers", "4", "--concurrency", "12"})
	var got int
	stderr := captureStderr(t, func() { got = resolveScanWorkers(fs, w, c) })
	if got != 4 {
		t.Fatalf("workers: got %d, want 4 (canonical wins)", got)
	}
	if stderr != "" {
		t.Fatalf("expected silent stderr when canonical is set, got %q",
			stderr)
	}
}

// TestResolveScanWorkers_NeitherSet asserts the auto/default value
// flows through unchanged when no related flag is passed.
func TestResolveScanWorkers_NeitherSet(t *testing.T) {
	fs, w, c := newScanWorkersFS(t, []string{})
	var got int
	stderr := captureStderr(t, func() { got = resolveScanWorkers(fs, w, c) })
	if got != constants.DefaultScanWorkers {
		t.Fatalf("workers: got %d, want default %d",
			got, constants.DefaultScanWorkers)
	}
	if stderr != "" {
		t.Fatalf("expected silent stderr at defaults, got %q", stderr)
	}
}
