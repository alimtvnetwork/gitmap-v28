package cmd

// Helpers wiring `--report-errors` into the scan command. Split out
// of scan.go to keep that file under the 200-line per-file budget
// while still expressing the (small) glue between the scanner /
// probe-runner callbacks and the errreport.Collector.

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/errreport"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/probe"
)

// newScanCollector returns a ready-to-use Collector when reportErrors
// is true, otherwise nil. Returning nil (rather than an empty
// collector) lets every downstream Add be a no-op without touching
// the hot path — the collector's nil-receiver methods short-circuit.
func newScanCollector(reportErrors bool) *errreport.Collector {
	if !reportErrors {
		return nil
	}

	return errreport.New(constants.Version, "scan")
}

// scanDirErrorCallback returns the OnDirError closure passed to
// scanner.ScanOptions. Captures `c` so the goroutine-safe
// errreport.Collector receives one entry per failed ReadDir.
func scanDirErrorCallback(c *errreport.Collector) func(string, error) {
	if c == nil {
		return nil
	}

	return func(path string, err error) {
		c.Add(errreport.PhaseScan, errreport.Entry{
			RepoPath: path,
			Step:     "readdir",
			Error:    err.Error(),
		})
	}
}

// installProbeFailureHook wires the background-probe runner's
// per-failure callback into the collector. No-op when either side is
// nil so callers can chain unconditionally.
func installProbeFailureHook(runner *probe.BackgroundRunner, c *errreport.Collector) {
	if runner == nil || c == nil {
		return
	}
	runner.SetFailureHook(func(rec model.ScanRecord, res probe.Result) {
		c.Add(errreport.PhaseScan, errreport.Entry{
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
func finalizeErrorReport(c *errreport.Collector, quiet bool) {
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
