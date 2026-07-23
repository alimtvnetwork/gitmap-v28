package clonefrom

// Summary + report helpers. Two outputs after Execute:
//
//   1. RenderSummary — human-readable end-of-batch table written
//      to stdout. Pass/fail counts + per-row status grid.
//
//   2. WriteReport — CSV file at .gitmap/clone-from-report-
//      <unixts>.csv with one row per Result. Mirrors the
//      `cn --csv` report convention so existing tooling that
//      ingests gitmap CSV reports works unchanged.

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// RenderSummary writes a human-readable summary to w. Format:
//
//	gitmap clone-from: 5 ok, 1 skipped, 2 failed (8 total)
//	report: .gitmap/clone-from-report-1735000000.csv
//
//	  ok       https://github.com/a/b.git    (1.2s)
//	  skipped  https://github.com/c/d.git    dest exists
//	  failed   https://github.com/e/f.git    fatal: repository not found
//	  ...
func RenderSummary(w io.Writer, results []Result, reportPath string) error {
	ok, skipped, failed := tallyResults(results)
	header := fmt.Sprintf(constants.MsgCloneFromSummaryHeader,
		ok, skipped, failed, len(results))
	if _, err := io.WriteString(w, header); err != nil {
		return err
	}
	if err := writeTransportLine(w, results); err != nil {
		return err
	}
	if len(reportPath) > 0 {
		if _, err := fmt.Fprintf(w, "report: %s\n\n", reportPath); err != nil {
			return err
		}
	}
	for _, r := range results {
		if _, err := io.WriteString(w, formatSummaryRow(r)); err != nil {
			return err
		}
	}

	return nil
}

// tallyResults counts each status category once. Three returns
// (ok, skipped, failed) are clearer at the call site than a
// map[string]int since the caller already knows the keys.
func tallyResults(results []Result) (int, int, int) {
	var ok, skipped, failed int
	for _, r := range results {
		switch r.Status {
		case constants.CloneFromStatusOK:
			ok++
		case constants.CloneFromStatusSkipped:
			skipped++
		case constants.CloneFromStatusFailed:
			failed++
		}
	}

	return ok, skipped, failed
}

// formatSummaryRow renders one result line. Status is left-padded
// to a fixed width so URLs line up vertically; detail (if any) is
// appended in a trailing column. Duration is shown only for `ok`
// rows — for skipped/failed it's noise.
func formatSummaryRow(r Result) string {
	tail := r.Detail
	if r.Status == constants.CloneFromStatusOK {
		tail = fmt.Sprintf("(%.1fs)", r.Duration.Seconds())
	}
	if len(tail) > 0 {
		return fmt.Sprintf("  %-7s  %s    %s\n", r.Status, r.Row.URL, tail)
	}

	return fmt.Sprintf("  %-7s  %s\n", r.Status, r.Row.URL)
}

// WriteReport persists the result set as CSV under .gitmap/. The
// timestamp suffix (Unix seconds) lets users keep a history of
// runs without one overwriting another. Returns the absolute path
// for the caller to surface in the summary header. On directory-
// create failure, returns "" + the error so the caller can decide
// whether the failure is fatal (it isn't — clones already
// happened, the report is bonus).
func WriteReport(results []Result) (string, error) {
	dir := ".gitmap"
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf(constants.ErrCloneFromReportMkdir, dir, err)
	}
	name := fmt.Sprintf("clone-from-report-%d.csv", time.Now().Unix())
	full := filepath.Join(dir, name)
	f, err := os.Create(full)
	if err != nil {
		return "", fmt.Errorf(constants.ErrCloneFromReportCreate, full, err)
	}
	defer f.Close()
	if err := writeReportRows(f, results); err != nil {
		return "", err
	}
	abs, _ := filepath.Abs(full)

	return abs, nil
}

// writeReportRows is the CSV-emit half of WriteReport. Split out
// so WriteReport stays under the function-length budget.
func writeReportRows(w io.Writer, results []Result) error {
	cw := csv.NewWriter(w)
	cw.UseCRLF = true // match other gitmap CSV reports per csvcrlf_contract_test.go
	if err := cw.Write([]string{"url", "dest", "branch", "depth", "status", "detail", "duration_seconds"}); err != nil {
		return err
	}
	for _, r := range results {
		rec := []string{r.Row.URL, r.Dest, r.Row.Branch,
			fmt.Sprintf("%d", r.Row.Depth),
			r.Status, r.Detail,
			fmt.Sprintf("%.3f", r.Duration.Seconds())}
		if err := cw.Write(rec); err != nil {
			return err
		}
	}
	cw.Flush()

	return cw.Error()
}

// JSON-report emit logic (reportRowJSON, reportEnvelopeJSON,
// transportTallyJSON, provenanceEntryJSON, writeReportRowsJSON,
// buildProvenanceEntries) lives in summary_json.go to keep this
// file under the 200-line per-file cap (mem://style/code-constraints).
