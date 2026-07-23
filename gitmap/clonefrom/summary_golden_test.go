package clonefrom

// summary_golden_test.go — locks the EXACT bytes of clone-from's CSV
// report (writeReportRows in summary.go) so any drift in column set,
// column order, CRLF line endings, status spelling, duration
// formatting, or quoting is caught in CI before it ships.
//
// Why golden vs. inline assertions: the report is consumed by
// downstream tooling (spreadsheets, jq-on-CSV, custom dashboards
// that ingest .gitmap/clone-from-report-*.csv). A silent reorder
// or rename of a column breaks those consumers in ways that don't
// surface in the gitmap test suite — so we pin the bytes.
//
// Two fixtures mirror the scan-side contract (formatter/scangolden_
// contract_test.go):
//
//   - empty: header row only — verifies the header itself is stable
//     and that zero-result runs still produce a valid CSV file.
//   - canonical: 3 rows covering ok / skipped / failed so every
//     status-dependent rendering branch is exercised.
//
// Regenerate after deliberate schema changes:
//
//	GITMAP_UPDATE_GOLDEN=1 go test ./gitmap/clonefrom/ -run \
//	  TestCloneFromReport_Golden
//
// then commit the regenerated files under clonefrom/testdata/ and
// call out the consumer-visible change in CHANGELOG.md.

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/goldenguard"
)

// canonicalReportResults builds a deterministic 3-row fixture that
// hits every status branch (ok, skipped, failed) AND every nullable
// column (Detail empty for ok-without-detail, Branch/Depth zero for
// the "no overrides" row). Hand-constructed (vs. derived from a
// real Execute run) so the test has zero external dependencies —
// no git, no filesystem, no network.
func canonicalReportResults() []Result {

	return []Result{
		{
			Row: Row{URL: "https://github.com/acme/widget.git",
				Branch: "main", Depth: 1},
			Dest:     "widget",
			Status:   constants.CloneFromStatusOK,
			Detail:   "",
			Duration: 1234 * time.Millisecond,
		},
		{
			Row:      Row{URL: "https://github.com/acme/gadget.git"},
			Dest:     "gadget",
			Status:   constants.CloneFromStatusSkipped,
			Detail:   "dest exists",
			Duration: 0,
		},
		{
			Row:      Row{URL: "git@github.com:acme/missing.git"},
			Dest:     "missing",
			Status:   constants.CloneFromStatusFailed,
			Detail:   "fatal: repository not found",
			Duration: 89 * time.Millisecond,
		},
	}
}

// TestCloneFromReport_Golden_Empty pins the header-only output for
// an empty result set. Header row is always emitted (even with zero
// results) so downstream tools can self-discover columns; CRLF is
// enforced for cross-platform byte-identical output.
func TestCloneFromReport_Golden_Empty(t *testing.T) {
	var buf bytes.Buffer
	if err := writeReportRows(&buf, nil); err != nil {
		t.Fatalf("writeReportRows: %v", err)
	}
	assertReportGolden(t, "clonefrom_report_empty.csv", buf.Bytes())
}

// TestCloneFromReport_Golden_Canonical pins the bytes for the 3-row
// canonical fixture. Catches drift in column set, column ORDER,
// status spelling, duration_seconds float format (must be %.3f),
// CRLF line endings, and CSV quoting behavior.
func TestCloneFromReport_Golden_Canonical(t *testing.T) {
	var buf bytes.Buffer
	if err := writeReportRows(&buf, canonicalReportResults()); err != nil {
		t.Fatalf("writeReportRows: %v", err)
	}
	assertReportGolden(t, "clonefrom_report_canonical.csv", buf.Bytes())
}

// assertReportGolden mirrors formatter.assertScanGolden — duplicated
// (rather than shared via a testutil package) because Go test helpers
// can't cross package boundaries without exporting them, and the
// assertion itself is small enough that a copy beats the API surface.
func assertReportGolden(t *testing.T, name string, got []byte) {
	t.Helper()
	path := filepath.Join("testdata", name)
	trigger := os.Getenv("GITMAP_UPDATE_GOLDEN") == "1"
	if goldenguard.AllowUpdate(t, trigger) {
		writeReportGolden(t, path, normalizeGoldenBytes(got))

		return
	}
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read golden %s: %v "+
			"(run with GITMAP_UPDATE_GOLDEN=1 and "+
			"GITMAP_ALLOW_GOLDEN_UPDATE=1 to create)", path, err)
	}
	// Normalize CRLF→LF on both sides — see formatter/scangolden_contract_test.go
	// for the full rationale. The production CRLF contract is still enforced
	// separately by csvcrlf_contract_test.go in cmd/.
	gotN := normalizeGoldenBytes(got)
	wantN := normalizeGoldenBytes(want)
	if !bytes.Equal(gotN, wantN) {
		t.Fatalf("golden mismatch for %s\n"+
			"--- want (%d bytes)\n%s\n--- got (%d bytes)\n%s",
			name, len(wantN), string(wantN), len(gotN), string(gotN))
	}
}

// normalizeGoldenBytes strips \r before \n so CRLF and LF fixtures
// compare equal regardless of platform / git autocrlf state.
func normalizeGoldenBytes(b []byte) []byte {
	return bytes.ReplaceAll(b, []byte("\r\n"), []byte("\n"))
}

// writeReportGolden persists a regenerated fixture and FAILS the test
// loudly so a CI run can never silently pass on a regenerate cycle —
// the regenerate path must be a deliberate two-step (run with the
// env var, then commit, then re-run without it).
func writeReportGolden(t *testing.T, path string, got []byte) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir testdata: %v", err)
	}
	if err := os.WriteFile(path, got, 0o644); err != nil {
		t.Fatalf("write golden %s: %v", path, err)
	}
	t.Fatalf("regenerated golden %s — re-run "+
		"without GITMAP_UPDATE_GOLDEN to confirm", path)
}
