package clonenow

// crossformat_golden_test.go — verifies that clone-now's CSV and JSON
// input parsers produce IDENTICAL Row slices when fed the SAME
// underlying scan records.
//
// Why this matters: `gitmap scan` emits BOTH .csv and .json next to
// each other under .gitmap/output/. Users routinely point clone-now
// at whichever they prefer (jq fans pick .json, spreadsheet users
// pick .csv). If the two formats ever diverge — a CSV column gets
// trimmed differently, JSON adds a field the CSV parser drops, etc.
// — the same scan would produce different clone trees depending on
// which file the user passed. This test fails the moment that
// happens.
//
// Strategy:
//
//  1. Build a hand-authored []model.ScanRecord covering both URL
//     fields, branch overrides, and a depth-bearing row.
//  2. Write it to disk via formatter.WriteCSV + formatter.WriteJSON
//     (the EXACT writers `gitmap scan` uses) into a temp dir.
//  3. ParseFile each one through clone-now's public entry point.
//  4. Render both Row slices to a stable, line-oriented canonical
//     form (one row per line, fields | -separated) and golden-pin
//     it. ONE golden file is shared between the two assertions —
//     if either format drifts the byte-compare fires immediately.
//
// Regenerate after deliberate schema changes:
//
//	GITMAP_UPDATE_GOLDEN=1 GITMAP_ALLOW_GOLDEN_UPDATE=1 \
//	  go test ./gitmap/clonenow/ -run TestCloneNowCrossFormat_Golden
//
// then commit clonenow/testdata/crossformat_rows.golden and call
// out the consumer-visible change in CHANGELOG.md.

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/formatter"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/goldenguard"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// crossFormatScanRecords mirrors formatter.canonicalScanRecords but
// is duplicated here (small struct, cheap to keep local) so this test
// remains decoupled from the formatter package's internal helpers.
// Same shape as the scan-side golden so a future investigator can
// diff the two by eye.
func crossFormatScanRecords() []model.ScanRecord {
	return []model.ScanRecord{
		{
			ID: 1, Slug: "acme/widget", RepoName: "widget",
			RepoID:           "github.com/acme/widget",
			HTTPSUrl:         "https://github.com/acme/widget.git",
			SSHUrl:           "git@github.com:acme/widget.git",
			DiscoveredURL:    "https://github.com/acme/widget.git",
			Branch:           "main",
			BranchSource:     "remote-head",
			RelativePath:     "acme/widget",
			AbsolutePath:     "/repos/acme/widget",
			CloneInstruction: "git clone https://github.com/acme/widget.git",
			Notes:            "",
			Depth:            1,
		},
		{
			ID: 2, Slug: "acme/gadget", RepoName: "gadget",
			RepoID:           "github.com/acme/gadget",
			HTTPSUrl:         "https://github.com/acme/gadget.git",
			SSHUrl:           "git@github.com:acme/gadget.git",
			DiscoveredURL:    "git@github.com:acme/gadget.git",
			Branch:           "develop",
			BranchSource:     "config",
			RelativePath:     "acme/gadget",
			AbsolutePath:     "/repos/acme/gadget",
			CloneInstruction: "git clone -b develop https://github.com/acme/gadget.git",
			Notes:            "pinned to develop",
			Depth:            4,
		},
	}
}

// canonicalRows renders a []Row to a stable, byte-comparable form.
// Format per row: `RepoName|HTTPSUrl|SSHUrl|Branch|RelativePath\n`.
// Pipe-separated rather than tab/CSV so the canonical form is
// trivially diff-readable AND can never be confused with one of the
// real input formats it's testing. Trailing newline matches POSIX.
func canonicalRows(rows []Row) []byte {
	var buf bytes.Buffer
	for _, r := range rows {
		fmt.Fprintf(&buf, "%s|%s|%s|%s|%s\n",
			r.RepoName, r.HTTPSUrl, r.SSHUrl, r.Branch, r.RelativePath)
	}

	return buf.Bytes()
}

// writeFormatFile materializes one format under tmpDir using the
// real formatter writers `gitmap scan` calls. Returns the absolute
// path so ParseFile can pick the format up via extension auto-detect.
func writeFormatFile(t *testing.T, tmpDir, name string,
	write func(*os.File) error,
) string {
	t.Helper()
	path := filepath.Join(tmpDir, name)
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create %s: %v", path, err)
	}
	defer f.Close()
	if err := write(f); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}

	return path
}

// TestCloneNowCrossFormat_Golden writes the canonical scan-record
// fixture as both CSV and JSON, parses each via ParseFile, and
// asserts the resulting Row slices serialize to byte-identical
// canonical form — pinned by a single golden file.
func TestCloneNowCrossFormat_Golden(t *testing.T) {
	tmp := t.TempDir()
	recs := crossFormatScanRecords()

	csvPath := writeFormatFile(t, tmp, "scan.csv", func(f *os.File) error {
		return formatter.WriteCSV(f, recs)
	})
	jsonPath := writeFormatFile(t, tmp, "scan.json", func(f *os.File) error {
		return formatter.WriteJSON(f, recs)
	})

	csvPlan, err := ParseFile(csvPath, "",
		constants.CloneNowModeHTTPS, constants.CloneNowOnExistsSkip)
	if err != nil {
		t.Fatalf("ParseFile(csv): %v", err)
	}
	jsonPlan, err := ParseFile(jsonPath, "",
		constants.CloneNowModeHTTPS, constants.CloneNowOnExistsSkip)
	if err != nil {
		t.Fatalf("ParseFile(json): %v", err)
	}

	csvBytes := canonicalRows(csvPlan.Rows)
	jsonBytes := canonicalRows(jsonPlan.Rows)

	// First-line-of-defense: the two formats must agree EXACTLY.
	// Reported separately from the golden mismatch so a CI failure
	// tells the author "CSV/JSON diverged" vs "schema drifted".
	if !bytes.Equal(csvBytes, jsonBytes) {
		t.Fatalf("CSV/JSON cross-format drift\n"+
			"--- csv (%d bytes)\n%s\n--- json (%d bytes)\n%s",
			len(csvBytes), string(csvBytes),
			len(jsonBytes), string(jsonBytes))
	}

	// Second-line-of-defense: pin the agreed bytes so a coordinated
	// drift (where BOTH parsers change in lockstep) still trips CI.
	assertCrossFormatGolden(t, "crossformat_rows.golden", csvBytes)
}

// assertCrossFormatGolden mirrors the assertReportGolden / assertScanGolden
// helpers in sibling packages — duplicated rather than shared because
// Go test helpers can't cross package boundaries without exporting them.
func assertCrossFormatGolden(t *testing.T, name string, got []byte) {
	t.Helper()
	path := filepath.Join("testdata", name)
	trigger := os.Getenv("GITMAP_UPDATE_GOLDEN") == "1"
	if goldenguard.AllowUpdate(t, trigger) {
		writeCrossFormatGolden(t, path, got)

		return
	}
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read golden %s: %v "+
			"(run with GITMAP_UPDATE_GOLDEN=1 and "+
			"GITMAP_ALLOW_GOLDEN_UPDATE=1 to create)", path, err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("golden mismatch for %s\n"+
			"--- want (%d bytes)\n%s\n--- got (%d bytes)\n%s",
			name, len(want), string(want), len(got), string(got))
	}
}

// writeCrossFormatGolden persists a regenerated fixture and FAILS
// the test loudly so a CI run can never silently pass on a
// regenerate cycle.
func writeCrossFormatGolden(t *testing.T, path string, got []byte) {
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
