// Package scanclone_test holds end-to-end tests that wire the
// scan-output writers to the clone-now parser/renderer to prove the
// pipeline `scan -> JSON -> clone-now -> dry-run plan` stays byte-
// identical across runs and across input formats.
//
// Why a fresh package rather than living under cmd_test or one of
// the producer packages: this test deliberately depends on BOTH
// `formatter` (the writer side) and `clonenow` (the consumer side)
// at the same time. Hosting it under either package would create a
// one-directional coupling that doesn't reflect the e2e contract.
// A neutral _test-suffixed package keeps the dependency arrows
// pointing the right way for a real consumer (a downstream user
// who imports both packages from outside the module).
package scanclone_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/clonenow"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/formatter"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/goldenguard"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// e2eScanRecords is the synthetic "scan output" the pipeline replays.
// Hand-authored (vs. invoking scanner.ScanDir on a real tree) so the
// test stays hermetic — no git, no fs walk, no clock — and so the
// e2e contract is decoupled from scanner internals. Schema mirrors
// the canonical scan-side fixture so a future investigator can diff
// the two by eye.
func e2eScanRecords() []model.ScanRecord {
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

// writeScanArtifact materializes one format file in tmpDir using the
// real `gitmap scan` writer. Returns the absolute path the clone-now
// parser will pick up via extension auto-detect.
func writeScanArtifact(t *testing.T, tmpDir, name string,
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

// renderPlanNormalized renders the dry-run plan and replaces the
// absolute source path AND the format token with stable placeholders
// so the byte comparison is portable across machines AND across
// input formats. The header's source-path and format-name slots are
// envelope metadata that legitimately differ per run/per format;
// the meat being pinned is the row count, the per-row blocks, and
// the overall layout — which MUST stay identical regardless of
// which format the user fed in.
func renderPlanNormalized(t *testing.T, plan clonenow.Plan, marker string) []byte {
	t.Helper()
	var buf bytes.Buffer
	if err := clonenow.Render(&buf, plan); err != nil {
		t.Fatalf("Render: %v", err)
	}
	out := strings.ReplaceAll(buf.String(), plan.Source, marker)
	// Format token lives in the dry-run header as `(<format>, mode=...)`.
	// Replace the parenthesized prefix only — bare `csv` / `json`
	// substrings could legitimately appear inside repo URLs or
	// branch names in other test fixtures.
	out = strings.Replace(out,
		"("+plan.Format+", mode=", "(<FORMAT>, mode=", 1)

	return []byte(out)
}

// scanThenClone is the e2e helper: write `recs` as one of the scan
// output formats, parse it back with clonenow.ParseFile, and return
// the rendered, source-path-normalized dry-run plan bytes.
func scanThenClone(t *testing.T, tmpDir, fileName string,
	recs []model.ScanRecord,
	write func(*os.File) error,
) []byte {
	t.Helper()
	path := writeScanArtifact(t, tmpDir, fileName, write)
	plan, err := clonenow.ParseFile(path, "",
		constants.CloneNowModeHTTPS, constants.CloneNowOnExistsSkip)
	if err != nil {
		t.Fatalf("ParseFile(%s): %v", fileName, err)
	}

	return renderPlanNormalized(t, plan, "<SOURCE>")
}

// TestScanClone_E2E_JSONRoundTrip exercises the full pipeline:
//
//	in-memory ScanRecords --(formatter.WriteJSON)--> .gitmap/output/gitmap.json
//	gitmap.json --(clonenow.ParseFile)--> Plan
//	Plan --(clonenow.Render)--> dry-run text
//
// then runs the SAME pipeline a second time on the SAME records and
// asserts the rendered bytes are identical run-over-run. Catches
// non-determinism introduced by map iteration, time.Now() leakage,
// or any future ordering bug in the writer/parser/renderer chain.
func TestScanClone_E2E_JSONRoundTrip(t *testing.T) {
	tmp := t.TempDir()
	recs := e2eScanRecords()

	first := scanThenClone(t, tmp, "gitmap.json", recs, func(f *os.File) error {
		return formatter.WriteJSON(f, recs)
	})
	second := scanThenClone(t, tmp, "gitmap2.json", recs, func(f *os.File) error {
		return formatter.WriteJSON(f, recs)
	})

	if !bytes.Equal(first, second) {
		t.Fatalf("non-deterministic e2e plan\n--- run 1 (%d bytes)\n%s\n--- run 2 (%d bytes)\n%s",
			len(first), string(first), len(second), string(second))
	}
	assertE2EGolden(t, "scanclone_e2e_plan.golden", first)
}

// TestScanClone_E2E_CSVMatchesJSON proves a user can swap which
// scan artifact they hand to clone-now (.csv vs .json) and get the
// IDENTICAL re-clone plan. This is the contract that lets shell
// scripts pick whichever format is more convenient without changing
// downstream behavior.
func TestScanClone_E2E_CSVMatchesJSON(t *testing.T) {
	tmp := t.TempDir()
	recs := e2eScanRecords()

	jsonPlan := scanThenClone(t, tmp, "gitmap.json", recs, func(f *os.File) error {
		return formatter.WriteJSON(f, recs)
	})
	csvPlan := scanThenClone(t, tmp, "gitmap.csv", recs, func(f *os.File) error {
		return formatter.WriteCSV(f, recs)
	})

	if !bytes.Equal(jsonPlan, csvPlan) {
		t.Fatalf("CSV/JSON e2e plan drift\n--- json (%d bytes)\n%s\n--- csv (%d bytes)\n%s",
			len(jsonPlan), string(jsonPlan),
			len(csvPlan), string(csvPlan))
	}
	// Both formats share the JSON-roundtrip golden — anything else
	// would let drift sneak in via a fixture mismatch.
	assertE2EGolden(t, "scanclone_e2e_plan.golden", csvPlan)
}

// assertE2EGolden mirrors the assertReportGolden / assertScanGolden
// helpers in producer packages — duplicated rather than shared
// because Go test helpers can't cross package boundaries without
// being exported, and the assertion is small enough that a copy
// beats the API surface.
func assertE2EGolden(t *testing.T, name string, got []byte) {
	t.Helper()
	path := filepath.Join("testdata", name)
	trigger := os.Getenv("GITMAP_UPDATE_GOLDEN") == "1"
	if goldenguard.AllowUpdate(t, trigger) {
		writeE2EGolden(t, path, got)

		return
	}
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read golden %s: %v "+
			"(run with GITMAP_UPDATE_GOLDEN=1 and "+
			"GITMAP_ALLOW_GOLDEN_UPDATE=1 to create)", path, err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("golden mismatch for %s\n--- want (%d bytes)\n%s\n--- got (%d bytes)\n%s",
			name, len(want), string(want), len(got), string(got))
	}
}

// writeE2EGolden persists a regenerated fixture and FAILS the test
// loudly so a CI run can never silently pass on a regenerate cycle.
func writeE2EGolden(t *testing.T, path string, got []byte) {
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
