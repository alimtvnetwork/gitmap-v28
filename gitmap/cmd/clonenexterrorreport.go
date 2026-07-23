package cmd

// `--errors-report` glue for `gitmap cn` batch mode (v3.130.0+).
// Split out of clonenextbatch.go to keep that file under the
// 200-line per-file budget. Mirrors the scan-side helpers in
// scanerrorreport.go: collector → JSON file at `.gitmap/reports/
// errors-<unixts>.json` next to the binary, only when failures
// exist.

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/errreport"
)

// writeCNErrorReport converts every BatchStatusFailed row in
// `results` into an errreport entry and writes the consolidated JSON
// report. No-op when reportErrors is false OR no failures occurred —
// clean runs leave nothing on disk. Any I/O error from the writer is
// logged to stderr by finalizeErrorReport but does NOT fail the
// command (the report is auxiliary).
func writeCNErrorReport(reportErrors bool, results []batchRowResult) {
	if !reportErrors {
		return
	}
	c := errreport.New(constants.Version, "clone-next")
	for _, r := range results {
		if r.Status != constants.BatchStatusFailed {
			continue
		}
		c.Add(errreport.PhaseClone, errreport.Entry{
			RepoPath: r.RepoPath,
			Step:     "clone",
			Error:    r.Detail,
		})
	}
	finalizeErrorReport(c, false)
}
