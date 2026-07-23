package cmd

// E2E-style CSV serialization tests for `gitmap cn --all` / `--csv`
// batch mode under concurrency (v3.126.0+). These tests exercise
// the production CSV writers (writeReportRows + writeBatchReport)
// to prove byte-identical output regardless of pool size.
//
// Stub helpers + helpers live in
// clonenextbatchconcurrent_e2e_helpers_test.go.

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// TestE2E_BatchConcurrency_ByteIdenticalAcrossPoolSizes is the
// strongest determinism guarantee: regardless of worker count the
// CSV bytes produced by writeReportRows are identical.
func TestE2E_BatchConcurrency_ByteIdenticalAcrossPoolSizes(t *testing.T) {
	installStubProcessor(t)
	repos := makeRepoPaths(50)

	baseline := runAndSerialize(t, repos, 1)
	// Pool sizes cover powers of two (1/2/4/8/16) AND odd/prime
	// values (3/5/7/9) — the latter catch edge cases where jobs
	// don't divide evenly across workers, exercising the trailing
	// partial batch and uneven-drain paths in collectBatchResults.
	for _, workers := range []int{2, 3, 4, 5, 7, 8, 9, 16} {
		got := runAndSerialize(t, repos, workers)
		if !bytes.Equal(baseline, got) {
			t.Fatalf("CSV bytes differ at workers=%d (sequential vs parallel)\n--- want ---\n%s\n--- got ---\n%s",
				workers, baseline, got)
		}
	}
}

// runAndSerialize runs the concurrent pool then writes the CSV via
// the production writeReportRows into a temp file (exercising the
// real *os.File path).
func runAndSerialize(t *testing.T, repos []string, workers int) []byte {
	t.Helper()
	results := processBatchReposConcurrent(repos, workers, nil)

	tmp, err := os.CreateTemp(t.TempDir(), "cn-batch-*.csv")
	if err != nil {
		t.Fatalf("temp csv: %v", err)
	}
	writeReportRows(tmp, results)
	tmp.Close()

	data, err := os.ReadFile(tmp.Name())
	if err != nil {
		t.Fatalf("read back csv: %v", err)
	}
	return data
}

// TestE2E_BatchConcurrency_FullWriteBatchReport drives the
// production writeBatchReport (the real cn entrypoint helper) end
// to end: temp CWD, real file creation with the unix-second name,
// real CSV bytes. Asserts header + one row per repo + input
// ordering preserved.
func TestE2E_BatchConcurrency_FullWriteBatchReport(t *testing.T) {
	installStubProcessor(t)
	t.Chdir(t.TempDir())

	repos := makeRepoPaths(20)
	results := processBatchReposConcurrent(repos, 4, nil)

	reportPath := writeBatchReport(results)
	if reportPath == "" {
		t.Fatal("writeBatchReport returned empty path")
	}

	data, err := os.ReadFile(reportPath)
	if err != nil {
		t.Fatalf("read report: %v", err)
	}
	lines := strings.Split(strings.TrimRight(string(data), "\n"), "\n")
	if len(lines) != len(repos)+1 {
		t.Fatalf("line count: got %d, want %d (1 header + %d rows)",
			len(lines), len(repos)+1, len(repos))
	}
	if !strings.HasPrefix(lines[0], "repo,from,to,status,detail") {
		t.Fatalf("header line: got %q", lines[0])
	}
	for i := 0; i < len(repos); i++ {
		want := fmt.Sprintf("repo-%d", i)
		if !strings.Contains(lines[i+1], want) {
			t.Fatalf("row %d should contain %q, got %q", i, want, lines[i+1])
		}
	}
}

// TestE2E_BatchConcurrency_CSVStatusCountsMatchInMemory closes the
// loop on the report-writing pipeline: the in-memory tallyBatch
// counts (the values printed in the 1-line summary) MUST equal the
// counts a downstream consumer would compute by re-parsing the CSV
// from disk. A drift between these two views — e.g. a writer that
// silently drops a row, mis-quotes the status field, or changes
// column order — would corrupt every downstream dashboard while
// leaving the in-process summary looking healthy.
func TestE2E_BatchConcurrency_CSVStatusCountsMatchInMemory(t *testing.T) {
	installStubProcessor(t)
	t.Chdir(t.TempDir())

	repos := makeRepoPaths(50)
	results := processBatchReposConcurrent(repos, 4, nil)

	wantOK, wantFailed, wantSkipped := tallyBatch(results)

	reportPath := writeBatchReport(results)
	if reportPath == "" {
		t.Fatal("writeBatchReport returned empty path")
	}
	gotOK, gotFailed, gotSkipped := tallyCSVStatuses(t, reportPath)

	if gotOK != wantOK || gotFailed != wantFailed || gotSkipped != wantSkipped {
		t.Fatalf("CSV status counts diverged from in-memory tally:\n"+
			"  in-memory: ok=%d failed=%d skipped=%d\n"+
			"  csv-parse: ok=%d failed=%d skipped=%d",
			wantOK, wantFailed, wantSkipped, gotOK, gotFailed, gotSkipped)
	}
}

// tallyCSVStatuses re-parses the produced report with encoding/csv
// (so quoting, escaping, and column order are all exercised) and
// returns the bucket counts derived purely from the on-disk bytes.
// Reads the status column by header name to stay resilient to
// future column reorderings within the same header set.
func tallyCSVStatuses(t *testing.T, path string) (ok, failed, skipped int) {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open report: %v", err)
	}
	defer f.Close()

	rows, err := csv.NewReader(f).ReadAll()
	if err != nil {
		t.Fatalf("parse csv: %v", err)
	}
	if len(rows) < 1 {
		t.Fatalf("csv has no header row")
	}
	statusCol := indexOfHeader(rows[0], "status")
	if statusCol < 0 {
		t.Fatalf("csv header missing 'status' column: %v", rows[0])
	}
	for _, row := range rows[1:] {
		ok, failed, skipped = bumpStatusBucket(row[statusCol], ok, failed, skipped)
	}
	return ok, failed, skipped
}

// indexOfHeader returns the column index of name in header, or -1.
func indexOfHeader(header []string, name string) int {
	for i, h := range header {
		if h == name {
			return i
		}
	}
	return -1
}

// bumpStatusBucket mirrors tallyBatch's switch but operates on the
// raw status string read from the CSV (no enum coupling beyond the
// shared constants). Unknown statuses are deliberately ignored so a
// future bucket addition surfaces as a count mismatch in the caller.
func bumpStatusBucket(status string, ok, failed, skipped int) (int, int, int) {
	switch status {
	case constants.BatchStatusOK:
		ok++
	case constants.BatchStatusFailed:
		failed++
	case constants.BatchStatusSkipped:
		skipped++
	}
	return ok, failed, skipped
}
