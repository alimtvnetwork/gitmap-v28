package clonefrom

// summary_jsongolden_test.go — locks the EXACT bytes of clone-from's
// JSON report (writeReportRowsJSON in summary.go) so any drift in
// field set, field order, indentation, status spelling, or
// duration_seconds float format is caught in CI before it ships.
//
// Why golden vs. inline assertions: the JSON report is consumed by
// downstream tooling (jq pipelines, dashboards, custom CI gates). A
// silent rename or reorder of a key — or a switch from `[]` to
// `null` for empty results — silently breaks those consumers.
//
// Two fixtures mirror the CSV-side contract (summary_golden_test.go):
//
//   - empty: `[]\n` only — verifies the empty-array convention
//     (NEVER `null`) and trailing-newline so consumers can treat
//     the file as an unconditional JSON array.
//   - canonical: 3 rows covering ok / skipped / failed so every
//     status branch and every nullable field is exercised.
//
// Regenerate after deliberate schema changes:
//
//	GITMAP_UPDATE_GOLDEN=1 GITMAP_ALLOW_GOLDEN_UPDATE=1 \
//	  go test ./gitmap/clonefrom/ -run TestCloneFromReportJSON_Golden
//
// then commit the regenerated files under clonefrom/testdata/ and
// call out the consumer-visible change in CHANGELOG.md.

import (
	"bytes"
	"testing"
)

// TestCloneFromReportJSON_Golden_Empty pins the empty-array output.
// Contract: zero results must produce `[]\n`, NEVER `null` —
// downstream `jq '. | length'` pipelines depend on the array shape.
func TestCloneFromReportJSON_Golden_Empty(t *testing.T) {
	var buf bytes.Buffer
	if err := writeReportRowsJSON(&buf, nil); err != nil {
		t.Fatalf("writeReportRowsJSON: %v", err)
	}
	assertReportGolden(t, "clonefrom_report_empty.json", buf.Bytes())
}

// TestCloneFromReportJSON_Golden_Canonical pins the bytes for the
// 3-row canonical fixture (shared with the CSV golden test). Catches
// drift in field set, field ORDER, status spelling, duration_seconds
// numeric format, JSON indentation, and HTML-escape behavior.
func TestCloneFromReportJSON_Golden_Canonical(t *testing.T) {
	var buf bytes.Buffer
	if err := writeReportRowsJSON(&buf, canonicalReportResults()); err != nil {
		t.Fatalf("writeReportRowsJSON: %v", err)
	}
	assertReportGolden(t, "clonefrom_report_canonical.json", buf.Bytes())
}
