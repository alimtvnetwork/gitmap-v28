package formatter

// Golden-file contract tests for the scan-record JSON and CSV
// emitters (formatter.WriteJSON and formatter.WriteCSV).
//
// Strategy: feed a fixed, hand-authored slice of model.ScanRecord
// values into each writer, capture the bytes into a buffer, and
// byte-compare against committed fixtures under formatter/testdata/.
// Two fixtures per format pin the most common consumer shapes:
//
//   - empty list  → headers/array shell only (jq + Excel sanity check)
//   - canonical 2 → multi-row payload covering every ScanRecord field
//
// Both writers also call emitWriteSummary which writes to the
// validation sink (default stderr). Tests redirect that sink to a
// throwaway buffer so summary lines don't pollute the golden output.
//
// To regenerate after a deliberate schema change:
//
//   GITMAP_UPDATE_GOLDEN=1 go test ./formatter/ -run ScanGolden
//
// Then commit the regenerated files under formatter/testdata/ and
// bump the consumer-facing changelog.

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/goldenguard"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// canonicalScanRecords builds a deterministic 2-row fixture covering
// every ScanRecord field with non-empty values. ID is set so the
// JSON `id` field renders as a stable integer; CSV ignores ID
// because ScanCSVHeaders intentionally omits it (legacy contract).
func canonicalScanRecords() []model.ScanRecord {

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
			Transport:        "https",
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
			// Depth 4 = at the DefaultMaxDepth boundary; this is exactly
			// the row a user would inspect to decide whether to widen
			// --max-depth before the next scan.
			Depth:     4,
			Transport: "ssh",
		},
	}
}

// silenceValidationSink swaps the package-level validation sink for
// a throwaway buffer so emitWriteSummary lines (e.g. "gitmap: csv:
// wrote 2 record(s), 0 validation issue(s)") do NOT contaminate the
// captured output buffer that we golden-compare. Returns a restore
// closure suitable for `defer`.
func silenceValidationSink(t *testing.T) func() {
	t.Helper()
	prev := SetValidationSink(&bytes.Buffer{})

	return func() { SetValidationSink(prev) }
}

// TestScanGolden_JSONEmpty pins the JSON shell for an empty record
// list. Empty input must encode as `null\n` because WriteJSON passes
// the slice straight to json.Encoder without a non-nil-init step —
// THIS IS THE CURRENT, DELIBERATELY-LOCKED-IN BEHAVIOR. If a future
// change normalizes empty → `[]`, the fixture regenerates and the
// changelog must call out the consumer-visible shape change.
func TestScanGolden_JSONEmpty(t *testing.T) {
	defer silenceValidationSink(t)()
	var buf bytes.Buffer
	if err := WriteJSON(&buf, nil); err != nil {
		t.Fatalf("WriteJSON: %v", err)
	}
	assertScanGolden(t, "scan_empty.json", buf.Bytes())
}

// TestScanGolden_JSONCanonical pins the bytes for the 2-row canonical
// fixture. Catches drift in field set, JSON tag names, declaration
// order, indentation, and trailing-newline behavior.
func TestScanGolden_JSONCanonical(t *testing.T) {
	defer silenceValidationSink(t)()
	var buf bytes.Buffer
	if err := WriteJSON(&buf, canonicalScanRecords()); err != nil {
		t.Fatalf("WriteJSON: %v", err)
	}
	assertScanGolden(t, "scan_canonical.json", buf.Bytes())
}

// TestScanGolden_CSVEmpty pins the CSV header-only shell for empty
// input. Header row is always emitted (even with zero records) so
// downstream tools can self-discover columns; CRLF is enforced.
func TestScanGolden_CSVEmpty(t *testing.T) {
	defer silenceValidationSink(t)()
	var buf bytes.Buffer
	if err := WriteCSV(&buf, nil); err != nil {
		t.Fatalf("WriteCSV: %v", err)
	}
	assertScanGolden(t, "scan_empty.csv", buf.Bytes())
}

// TestScanGolden_CSVCanonical pins the bytes for the 2-row CSV
// fixture. Catches drift in column set, column ORDER (legacy 9-col
// layout — no id/slug), CRLF line endings, and quoting behavior.
func TestScanGolden_CSVCanonical(t *testing.T) {
	defer silenceValidationSink(t)()
	var buf bytes.Buffer
	if err := WriteCSV(&buf, canonicalScanRecords()); err != nil {
		t.Fatalf("WriteCSV: %v", err)
	}
	assertScanGolden(t, "scan_canonical.csv", buf.Bytes())
}

// assertScanGolden is the package-local twin of cmd's assertGoldenBytes.
// Duplicated rather than shared because Go test helpers cannot cross
// package boundaries without exporting them — and the assertion is
// small enough that a copy is cheaper than the API surface.
func assertScanGolden(t *testing.T, name string, got []byte) {
	t.Helper()
	path := filepath.Join("testdata", name)
	trigger := os.Getenv("GITMAP_UPDATE_GOLDEN") == "1"
	if goldenguard.AllowUpdate(t, trigger) {
		writeScanGolden(t, path, normalizeGoldenBytes(got))

		return
	}
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read golden %s: %v (run with "+
			"GITMAP_UPDATE_GOLDEN=1 and "+
			"GITMAP_ALLOW_GOLDEN_UPDATE=1 to create)", path, err)
	}
	// Normalize CRLF→LF on both sides: the production CSV writer emits
	// CRLF (csv.Writer.UseCRLF=true, pinned by csvcrlf_contract_test.go)
	// but the on-disk golden is checked out as LF on every platform via
	// `.gitattributes` (testdata/** text eol=lf) so the fixture stays
	// hand-readable and cross-platform consistent. Without this
	// normalization Linux+Windows CI both fail with a 1-byte-per-line
	// drift. The CRLF contract is still enforced separately by
	// TestCSVCRLF_WriteCSV in the same package — that test inspects the
	// in-memory bytes, not the fixture.
	gotN := normalizeGoldenBytes(got)
	wantN := normalizeGoldenBytes(want)
	if !bytes.Equal(gotN, wantN) {
		t.Fatalf("golden mismatch for %s\n--- want (%d bytes)\n%s\n--- got (%d bytes)\n%s",
			name, len(wantN), string(wantN), len(gotN), string(gotN))
	}
}

// normalizeGoldenBytes strips \r before \n so CRLF and LF fixtures
// compare equal. JSON goldens have no \r so this is a no-op for them.
func normalizeGoldenBytes(b []byte) []byte {
	return bytes.ReplaceAll(b, []byte("\r\n"), []byte("\n"))
}

// writeScanGolden persists a regenerated fixture and FAILS the test
// loudly so a CI run can never silently pass on a regenerate cycle.
func writeScanGolden(t *testing.T, path string, got []byte) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir testdata: %v", err)
	}
	if err := os.WriteFile(path, got, 0o644); err != nil {
		t.Fatalf("write golden %s: %v", path, err)
	}
	t.Fatalf("regenerated golden %s — re-run without GITMAP_UPDATE_GOLDEN to confirm", path)
}
