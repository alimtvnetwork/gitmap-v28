package formatter

// CSV header & column-order contract tests (v3.136.0+).
//
// These tests pin the exact byte layout of the CSV that `gitmap scan`
// emits, so accidental reordering, renaming, or dropping a column
// fails CI loudly instead of silently breaking downstream consumers
// that diff CSVs byte-for-byte.
//
// Two layers of assertion:
//
//  1. Header contract — `constants.ScanCSVHeaders` must equal an
//     immutable expected slice in this file. Bumping the schema
//     requires a deliberate edit to BOTH lists.
//  2. Row contract — `WriteCSV` of a fixed record must produce the
//     exact byte sequence expected, including header line, comma
//     separators, and CRLF line endings emitted by encoding/csv.
//
// The fixed record uses only ASCII values without commas/quotes so
// the expected bytes are easy to read in source. A separate
// quoting test covers the encoding/csv escaping rules.

import (
	"bytes"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// expectedScanCSVHeaders is the locked-in canonical header list. Any
// change here is intentional and must be paired with a constants edit
// AND a documented schema bump (CSV consumers diff on these names).
var expectedScanCSVHeaders = []string{
	"repoName", "httpsUrl", "sshUrl", "branch", "branchSource",
	"relativePath", "absolutePath", "cloneInstruction", "notes", "depth",
	"repoId", "discoveredUrl", "transport",
}

// expectedLatestBranchCSVHeaders pins the latest-branch CSV layout.
var expectedLatestBranchCSVHeaders = []string{
	"branch", "remote", "sha", "commitDate", "subject", "ref",
}

// TestScanCSVHeaders_ExactOrder asserts the scan-CSV header slice
// matches the locked expectation byte-for-byte AND in order. Catches:
// reordering, renaming, insertion, deletion.
func TestScanCSVHeaders_ExactOrder(t *testing.T) {
	got := constants.ScanCSVHeaders
	if len(got) != len(expectedScanCSVHeaders) {
		t.Fatalf("ScanCSVHeaders length: got %d, want %d (%v)",
			len(got), len(expectedScanCSVHeaders), expectedScanCSVHeaders)
	}
	for i, want := range expectedScanCSVHeaders {
		if got[i] != want {
			t.Errorf("ScanCSVHeaders[%d]: got %q, want %q", i, got[i], want)
		}
	}
}

// TestLatestBranchCSVHeaders_ExactOrder mirrors the scan check for
// the latest-branch CSV emitted by `gitmap latest-branch --output csv`.
func TestLatestBranchCSVHeaders_ExactOrder(t *testing.T) {
	got := constants.LatestBranchCSVHeaders
	if len(got) != len(expectedLatestBranchCSVHeaders) {
		t.Fatalf("LatestBranchCSVHeaders length: got %d, want %d (%v)",
			len(got), len(expectedLatestBranchCSVHeaders), expectedLatestBranchCSVHeaders)
	}
	for i, want := range expectedLatestBranchCSVHeaders {
		if got[i] != want {
			t.Errorf("LatestBranchCSVHeaders[%d]: got %q, want %q", i, got[i], want)
		}
	}
}

// TestWriteCSV_ExactBytes pins the full byte output of WriteCSV for a
// single deterministic record. encoding/csv emits CRLF (\r\n) line
// endings per RFC 4180, so the expected blob uses \r\n. Any drift in
// header text, column order, or separator/line-terminator choice
// fails this test.
func TestWriteCSV_ExactBytes(t *testing.T) {
	rec := model.ScanRecord{
		RepoName: "repo-a", HTTPSUrl: "https://example.com/u/repo-a.git",
		SSHUrl: "git@example.com:u/repo-a.git", Branch: "main",
		BranchSource: "head", RelativePath: "p/repo-a",
		AbsolutePath: "/p/repo-a", CloneInstruction: "git clone X",
		Notes: "n", Depth: 3,
		RepoID:        "example.com/u/repo-a",
		DiscoveredURL: "https://example.com/u/repo-a.git",
		Transport:     "https",
	}
	var buf bytes.Buffer
	if err := WriteCSV(&buf, []model.ScanRecord{rec}); err != nil {
		t.Fatalf("WriteCSV: %v", err)
	}
	want := "repoName,httpsUrl,sshUrl,branch,branchSource,relativePath,absolutePath,cloneInstruction,notes,depth,repoId,discoveredUrl,transport\r\n" +
		"repo-a,https://example.com/u/repo-a.git,git@example.com:u/repo-a.git,main,head,p/repo-a,/p/repo-a,git clone X,n,3,example.com/u/repo-a,https://example.com/u/repo-a.git,https\r\n"
	if got := buf.String(); got != want {
		t.Errorf("WriteCSV bytes mismatch.\n got: %q\nwant: %q", got, want)
	}
}

// TestWriteCSV_HeaderIsFirstLine guards against any future change
// that interleaves preamble / BOM / blank lines before the header.
// Downstream parsers use `head -n1` to grab the header — keep it on
// line 1, no exceptions.
func TestWriteCSV_HeaderIsFirstLine(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteCSV(&buf, nil); err != nil {
		t.Fatalf("WriteCSV: %v", err)
	}
	first, _, _ := strings.Cut(buf.String(), "\r\n")
	want := strings.Join(expectedScanCSVHeaders, ",")
	if first != want {
		t.Errorf("first CSV line: got %q, want %q", first, want)
	}
}

// TestWriteCSV_ColumnCountMatchesHeader catches a row builder that
// writes more or fewer fields than the header advertises — a class
// of bug that makes downstream `csv.Reader` reject the file with
// `wrong number of fields`.
func TestWriteCSV_ColumnCountMatchesHeader(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteCSV(&buf, []model.ScanRecord{{RepoName: "x"}}); err != nil {
		t.Fatalf("WriteCSV: %v", err)
	}
	lines := strings.Split(strings.TrimRight(buf.String(), "\r\n"), "\r\n")
	if len(lines) < 2 {
		t.Fatalf("expected header + 1 row, got %d lines", len(lines))
	}
	hdrCols := strings.Count(lines[0], ",") + 1
	rowCols := strings.Count(lines[1], ",") + 1
	if hdrCols != rowCols {
		t.Errorf("column count drift: header=%d row=%d", hdrCols, rowCols)
	}
	if hdrCols != len(expectedScanCSVHeaders) {
		t.Errorf("header column count: got %d, want %d", hdrCols, len(expectedScanCSVHeaders))
	}
}
