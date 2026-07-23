package cmd

// Integration test: scan → export JSON/CSV → clone-from for both
// https and ssh "modes". Drives the production functions end-to-end
// against real on-disk artifacts:
//
//   1. Build a tiny bare repo + two fake worktrees that point at it.
//   2. scanner.ScanDirWithOptions discovers the worktrees.
//   3. mapper.BuildRecords (mode=https / mode=ssh) populates the
//      HTTPSUrl + SSHUrl columns from the real git remote.
//   4. formatter.WriteJSON / formatter.WriteCSV serialize the
//      records to disk via a real *os.File.
//   5. The exported scan file is transformed into a clone-from
//      manifest by reading the column the mode selected
//      (httpsUrl for ModeHTTPS, sshUrl for ModeSSH) — this is the
//      round-trip the test guards.
//   6. clonefrom.ParseFile + clonefrom.Execute re-clone every row
//      and the assertion is "every row Status == ok".
//
// Both "modes" use file:// URLs underneath because CI has no
// reachable git server. The mode wiring still gets exercised: the
// transform step picks a different column per mode and a regression
// in mapper.selectCloneURL or in either parser would surface as a
// failed Result row. Fixture helpers live in
// scan_export_clonefrom_fixture_test.go to keep this file under the
// 200-line per-file cap.

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/clonefrom"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/formatter"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/mapper"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/scanner"
)

// TestScanExportCloneFrom_HTTPSAndSSH_RoundTrips runs the full
// pipeline once per mode × format. Sub-tests share the bare-repo
// fixture so the slow path (git init + commit + bare clone) only
// runs once per top-level test invocation.
func TestScanExportCloneFrom_HTTPSAndSSH_RoundTrips(t *testing.T) {
	requireGitForIntegration(t)
	bare := makeIntegrationBareRepo(t)
	scanRoot := seedScanTree(t, bare)

	for _, mode := range []string{constants.ModeHTTPS, constants.ModeSSH} {
		t.Run("json/"+mode, func(t *testing.T) {
			runRoundTrip(t, scanRoot, mode, "json")
		})
		t.Run("csv/"+mode, func(t *testing.T) {
			runRoundTrip(t, scanRoot, mode, "csv")
		})
	}
}

// runRoundTrip executes one scan → export(format) → transform →
// clone-from cycle and asserts every executed row landed ok.
func runRoundTrip(t *testing.T, scanRoot, mode, format string) {
	t.Helper()
	records := scanAndBuildRecords(t, scanRoot, mode)
	exportPath := exportRecords(t, records, format)
	manifest := writeCloneFromManifest(t, records, mode, format)
	executePlanAndAssertOK(t, manifest, exportPath)
}

// scanAndBuildRecords runs the real scanner against the seeded
// worktree tree, then converts the RepoInfo slice to ScanRecords
// the same way `gitmap scan` does at runtime.
func scanAndBuildRecords(t *testing.T, root, mode string) []model.ScanRecord {
	t.Helper()
	repos, err := scanner.ScanDirWithOptions(root, scanner.ScanOptions{})
	if err != nil {
		t.Fatalf("scanner.ScanDirWithOptions: %v", err)
	}
	if len(repos) < 2 {
		t.Fatalf("scanner found %d repos, want >=2 (root=%s)", len(repos), root)
	}

	return mapper.BuildRecords(repos, mode, "")
}

// exportRecords writes the records via the production WriteJSON /
// WriteCSV writers to a tempdir file and returns the absolute path.
// The path is logged so a failure mid-pipeline points at the bytes
// that actually got serialized.
func exportRecords(t *testing.T, records []model.ScanRecord, format string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "gitmap."+format)
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create export file: %v", err)
	}
	defer f.Close()
	if err := pickFormatterWriter(format)(f, records); err != nil {
		t.Fatalf("write %s: %v", format, err)
	}

	return path
}

// pickFormatterWriter returns the formatter writer matching format.
// Centralized so exportRecords stays under the per-function budget.
func pickFormatterWriter(format string) func(*os.File, []model.ScanRecord) error {
	if format == "json" {
		return func(f *os.File, r []model.ScanRecord) error { return formatter.WriteJSON(f, r) }
	}

	return func(f *os.File, r []model.ScanRecord) error { return formatter.WriteCSV(f, r) }
}

// writeCloneFromManifest builds a clone-from input file (always
// JSON to keep the parser path simple and the dest column explicit).
// The URL column read from the records depends on `mode` — this is
// the round-trip step the integration test exists to guard.
func writeCloneFromManifest(t *testing.T, records []model.ScanRecord, mode, originFormat string) string {
	t.Helper()
	destRoot := t.TempDir()
	rows := make([]map[string]any, 0, len(records))
	for i, rec := range records {
		url := pickURLForMode(rec, mode)
		if url == "" {
			t.Fatalf("record %d has empty URL for mode %s", i, mode)
		}
		rows = append(rows, map[string]any{
			"url":  url,
			"dest": filepath.Join(destRoot, originFormat+"-"+mode+"-"+rec.RepoName),
		})
	}
	path := filepath.Join(t.TempDir(), "clone-from."+originFormat+"."+mode+".json")
	writeJSONFile(t, path, rows)

	return path
}

// pickURLForMode is the column-selection rule the round-trip exists
// to test. ModeHTTPS reads HTTPSUrl; ModeSSH reads SSHUrl. Any new
// mode would need a corresponding branch here AND in selectCloneURL.
func pickURLForMode(rec model.ScanRecord, mode string) string {
	if mode == constants.ModeSSH {
		return rec.SSHUrl
	}

	return rec.HTTPSUrl
}

// writeJSONFile encodes rows to path as JSON. Tiny helper so the
// caller stays under the per-function budget.
func writeJSONFile(t *testing.T, path string, rows []map[string]any) {
	t.Helper()
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create manifest: %v", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(rows); err != nil {
		t.Fatalf("encode manifest: %v", err)
	}
}

// executePlanAndAssertOK parses the manifest with the real
// clonefrom.ParseFile, executes it, and fails with the exported
// scan path in the message so debugging starts at the bytes the
// formatter wrote.
func executePlanAndAssertOK(t *testing.T, manifest, exportPath string) {
	t.Helper()
	plan, err := clonefrom.ParseFile(manifest)
	if err != nil {
		t.Fatalf("ParseFile(%s): %v (export=%s)", manifest, err, exportPath)
	}
	results := clonefrom.Execute(plan, "", os.Stderr)
	if len(results) != len(plan.Rows) {
		t.Fatalf("results=%d, want %d", len(results), len(plan.Rows))
	}
	for i, r := range results {
		if r.Status != constants.CloneFromStatusOK {
			t.Fatalf("row %d status=%q detail=%q (export=%s)",
				i, r.Status, r.Detail, exportPath)
		}
	}
}
