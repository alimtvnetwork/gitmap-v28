package cmd

// CSV report writing + 1-line summary helpers for `gitmap cn` batch
// mode. Split out of clonenextbatch.go to keep that file under the
// 200-line per-file budget.

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// writeBatchReport emits cn-batch-<unixts>.csv with one row per repo and
// returns the absolute path to the report. A write failure is logged and
// the function returns "" so the caller can decide how loud to be.
func writeBatchReport(results []batchRowResult) string {
	name := fmt.Sprintf("cn-batch-%d.csv", time.Now().Unix())
	abs, err := filepath.Abs(name)
	if err != nil {
		abs = name
	}

	file, err := os.Create(abs)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.WarnCloneNextBatchReport, err)

		return ""
	}
	defer file.Close()

	writeReportRows(file, results)

	return abs
}

// writeReportRows formats and writes the header + one row per result.
func writeReportRows(file *os.File, results []batchRowResult) {
	fmt.Fprintln(file, "repo,from,to,status,detail")
	for _, r := range results {
		fmt.Fprintf(file, "%q,%s,%s,%s,%q\n",
			r.RepoPath, r.FromVersion, r.ToVersion, r.Status, r.Detail)
	}
}

// printBatchSummary prints a 1-line tally + the report path.
func printBatchSummary(results []batchRowResult, reportPath string) {
	ok, failed, skipped := tallyBatch(results)
	fmt.Printf(constants.MsgCloneNextBatchSummary, ok, failed, skipped)
	if len(reportPath) > 0 {
		fmt.Printf(constants.MsgCloneNextBatchReport, reportPath)
	}
}

// tallyBatch counts each status bucket.
func tallyBatch(results []batchRowResult) (ok, failed, skipped int) {
	for _, r := range results {
		switch r.Status {
		case constants.BatchStatusOK:
			ok++
		case constants.BatchStatusFailed:
			failed++
		case constants.BatchStatusSkipped:
			skipped++
		}
	}

	return ok, failed, skipped
}
