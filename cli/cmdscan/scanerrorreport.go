package cmdscan

// Helpers wiring `--report-errors` into the scan command. Split out
// of scan.go to keep that file under the 200-line per-file budget
// while still expressing the (small) glue between the scanner /
// probe-runner callbacks and the diag.Collector.

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/diag"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/scanpipe"
	"os"
)

// newScanCollector returns a ready-to-use Collector when reportErrors
// is true, otherwise nil. Returning nil (rather than an empty
// collector) lets every downstream Add be a no-op without touching
// the hot path — the collector's nil-receiver methods short-circuit.
func newScanCollector(reportErrors bool) *diag.Collector {
	if !reportErrors {
		return nil
	}

	return diag.New(constants.Version, "scan")
}

// scanDirErrorCallback returns the OnDirError closure passed to
// scanpipe.ScanOptions. Captures `c` so the goroutine-safe
// diag.Collector receives one entry per failed ReadDir.
func scanDirErrorCallback(c *diag.Collector) func(string, error) {
	if c == nil {
		return nil
	}

	return func(path string, err error) {
		c.Add(diag.PhaseScan, diag.Entry{
			RepoPath: path,
			Step:     "readdir",
			Error:    err.Error(),
		})
	}
}

// installProbeFailureHook wires the background-probe runner's
// per-failure callback into the collector. No-op when either side is
// nil so callers can chain unconditionally.
func installProbeFailureHook(runner *scanpipe.BackgroundRunner, c *diag.Collector) {
	if runner == nil || c == nil {
		return
	}

	runner.SetFailureHook(func(rec model.ScanRecord, res scanpipe.ProbeResult) {
		c.Add(diag.PhaseScan, diag.Entry{
			RepoPath:  rec.RelativePath,
			RemoteURL: rec.HTTPSUrl,
			Step:      "probe-" + res.Method,
			Error:     res.Error,
		})
	})
}

// finalizeErrorReport writes the report to disk if any failures were
// recorded. Errors writing the report are logged to stderr but never
// fail the parent command — the report is auxiliary, not load-bearing.
func finalizeErrorReport(c *diag.Collector, quiet bool) {
	if c == nil {
		return
	}

	scan, clone := c.Count()
	if scan+clone == 0 {
		return
	}

	path, err := c.WriteIfAny(resolveBinaryDir())
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ failed to write error report: %v\n", err)

		return
	}

	if !quiet && len(path) > 0 {
		fmt.Printf("  📝 %d failure(s) recorded → %s\n", scan+clone, path)
	}
}
