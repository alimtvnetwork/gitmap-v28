// Package cmd — contract tests for output filenames + extensions.
//
// Downstream tooling (CI parsers, IDE pickers, the docs site, the
// PowerShell loaders shipped under .gitmap/output/, and external
// scripts that grep for `gitmap.json` / `gitmap.csv`) all rely on
// the EXACT names and extensions produced by `gitmap scan` and
// `gitmap clone`. These tests lock in:
//
//  1. The basenames of every file emitted by `writeAllOutputs`.
//  2. The extension on each (.csv / .json / .md / .ps1 / .txt).
//  3. The shorthand → path mapping consumed by `gitmap clone json`
//     / `gitmap clone csv` / `gitmap clone text`.
//  4. The timestamped report names from `clonefrom.WriteReport` and
//     `errreport.WriteIfAny` — both must end in the documented
//     extension (.csv / .json) so post-run pipelines globbing for
//     `clone-from-report-*.csv` / `errors-*.json` keep matching.
//  5. The `resolveOutFile` precedence: explicit override wins,
//     otherwise `<outputDir>/<defaultName>` (no extension mutation).
//
// If any of these change intentionally, update this test AND the
// docs page that advertises the filenames in the same PR.

package cmd

import (
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/clonefrom"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/errreport"
)

// TestScanOutputFilenames_Contract pins every default scan/clone
// artifact basename. A change here breaks downstream globbers.
func TestScanOutputFilenames_Contract(t *testing.T) {
	cases := []struct {
		label    string
		got      string
		wantName string
		wantExt  string
	}{
		{"csv", constants.DefaultCSVFile, "gitmap.csv", ".csv"},
		{"json", constants.DefaultJSONFile, "gitmap.json", ".json"},
		{"text", constants.DefaultTextFile, "gitmap.txt", ".txt"},
		{"structure", constants.DefaultStructureFile, "folder-structure.md", ".md"},
		{"clone-script", constants.DefaultCloneScript, "clone.ps1", ".ps1"},
		{"direct-clone", constants.DefaultDirectCloneScript, "direct-clone.ps1", ".ps1"},
		{"direct-clone-ssh", constants.DefaultDirectCloneSSHScript, "direct-clone-ssh.ps1", ".ps1"},
		{"desktop", constants.DefaultDesktopScript, "register-desktop.ps1", ".ps1"},
	}
	for _, tc := range cases {
		t.Run(tc.label, func(t *testing.T) {
			if tc.got != tc.wantName {
				t.Fatalf("basename drift: got %q, want %q", tc.got, tc.wantName)
			}
			if filepath.Ext(tc.got) != tc.wantExt {
				t.Fatalf("extension drift on %s: got %q, want %q",
					tc.got, filepath.Ext(tc.got), tc.wantExt)
			}
		})
	}
}

// TestCloneShorthand_Contract verifies `gitmap clone json|csv|text`
// resolve to the canonical default paths under DefaultOutputFolder.
// We don't call resolveCloneShorthand directly because it Stat()s
// the path and exits on miss; we verify the mapping by reproducing
// it from the same constants the production code uses, AND assert
// the constants themselves haven't drifted.
func TestCloneShorthand_Contract(t *testing.T) {
	if constants.ShorthandJSON != "json" {
		t.Fatalf("ShorthandJSON drift: %q", constants.ShorthandJSON)
	}
	if constants.ShorthandCSV != "csv" {
		t.Fatalf("ShorthandCSV drift: %q", constants.ShorthandCSV)
	}
	if constants.ShorthandText != "text" {
		t.Fatalf("ShorthandText drift: %q", constants.ShorthandText)
	}

	wantJSON := filepath.Join(constants.DefaultOutputFolder, "gitmap.json")
	wantCSV := filepath.Join(constants.DefaultOutputFolder, "gitmap.csv")
	wantText := filepath.Join(constants.DefaultOutputFolder, "gitmap.txt")

	gotJSON := filepath.Join(constants.DefaultOutputFolder, constants.DefaultJSONFile)
	gotCSV := filepath.Join(constants.DefaultOutputFolder, constants.DefaultCSVFile)
	gotText := filepath.Join(constants.DefaultOutputFolder, constants.DefaultTextFile)

	if gotJSON != wantJSON {
		t.Fatalf("clone json shorthand drift: got %q, want %q", gotJSON, wantJSON)
	}
	if gotCSV != wantCSV {
		t.Fatalf("clone csv shorthand drift: got %q, want %q", gotCSV, wantCSV)
	}
	if gotText != wantText {
		t.Fatalf("clone text shorthand drift: got %q, want %q", gotText, wantText)
	}
}

// TestResolveOutFile_Precedence locks the override semantics:
// explicit --out beats the default; empty falls back to <dir>/<name>.
// Any extension mutation here would silently rename downstream files.
func TestResolveOutFile_Precedence(t *testing.T) {
	cases := []struct {
		label       string
		outFile     string
		outputDir   string
		defaultName string
		want        string
	}{
		{"override-wins", "/abs/custom.csv", "ignored", "gitmap.csv", "/abs/custom.csv"},
		{"empty-uses-default-csv", "", "outdir", "gitmap.csv", filepath.Join("outdir", "gitmap.csv")},
		{"empty-uses-default-json", "", "outdir", "gitmap.json", filepath.Join("outdir", "gitmap.json")},
		{"override-keeps-extension", "report.tsv", "outdir", "gitmap.csv", "report.tsv"},
	}
	for _, tc := range cases {
		t.Run(tc.label, func(t *testing.T) {
			got := resolveOutFile(tc.outFile, tc.outputDir, tc.defaultName)
			if got != tc.want {
				t.Fatalf("resolveOutFile(%q,%q,%q) = %q, want %q",
					tc.outFile, tc.outputDir, tc.defaultName, got, tc.want)
			}
		})
	}
}

// reCloneFromReport pins the clone-from report name shape.
// `clone-from-report-<unix-seconds>.csv`. Note: only digits between
// the prefix and `.csv`. Anything else (e.g. `.csv.tmp`, `.CSV`,
// hyphenated suffixes) breaks the docs glob and the published spec.
var reCloneFromReport = regexp.MustCompile(`^clone-from-report-\d+\.csv$`)

// TestCloneFromReportName_Contract calls the real WriteReport and
// asserts the produced basename matches the documented pattern.
// Uses an empty results slice → the writer still emits the header
// row, which is fine; we only care about the filename here.
func TestCloneFromReportName_Contract(t *testing.T) {
	t.Chdir(t.TempDir()) // WriteReport writes under ./.gitmap

	before := time.Now().Unix()
	abs, err := clonefrom.WriteReport(nil)
	after := time.Now().Unix()
	if err != nil {
		t.Fatalf("WriteReport: %v", err)
	}
	base := filepath.Base(abs)
	if !reCloneFromReport.MatchString(base) {
		t.Fatalf("report basename %q does not match %q", base, reCloneFromReport)
	}
	if filepath.Ext(base) != ".csv" {
		t.Fatalf("report extension drift: got %q, want %q", filepath.Ext(base), ".csv")
	}
	// Parent directory must be `.gitmap` (no `output/` nesting —
	// that would change the documented path users grep for).
	parent := filepath.Base(filepath.Dir(abs))
	if parent != ".gitmap" {
		t.Fatalf("report parent dir drift: got %q, want %q", parent, ".gitmap")
	}
	// Sanity: timestamp inside the name falls within [before,after].
	tsStr := strings.TrimSuffix(strings.TrimPrefix(base, "clone-from-report-"), ".csv")
	if len(tsStr) == 0 {
		t.Fatalf("missing timestamp in %q", base)
	}
	// Convert via time math: regex already proved digits-only.
	var ts int64
	for _, c := range tsStr {
		ts = ts*10 + int64(c-'0')
	}
	if ts < before || ts > after {
		t.Fatalf("timestamp %d outside [%d,%d]", ts, before, after)
	}
}

// reErrorsReport pins the error-report name shape:
// `errors-<unix-seconds>.json`.
var reErrorsReport = regexp.MustCompile(`^errors-\d+\.json$`)

// TestErrorsReportName_Contract drives errreport.Collector through
// a single failure so WriteIfAny actually produces a file, then
// validates the basename + extension + parent dir layout.
func TestErrorsReportName_Contract(t *testing.T) {
	dir := t.TempDir()
	c := errreport.New("v0.0.0-test", "clone-next")
	c.Add(errreport.PhaseClone, errreport.Entry{
		RepoPath:  "/tmp/example",
		RemoteURL: "https://example.com/x.git",
		Step:      "clone",
		Error:     "clone failed",
	})

	abs, err := c.WriteIfAny(dir)
	if err != nil {
		t.Fatalf("WriteIfAny: %v", err)
	}
	if abs == "" {
		t.Fatalf("expected a path, got empty (failure was registered)")
	}
	base := filepath.Base(abs)
	if !reErrorsReport.MatchString(base) {
		t.Fatalf("errors report basename %q does not match %q", base, reErrorsReport)
	}
	if filepath.Ext(base) != ".json" {
		t.Fatalf("errors report extension drift: got %q, want %q",
			filepath.Ext(base), ".json")
	}
	wantParent := filepath.Join(dir, ".gitmap", "reports")
	if filepath.Dir(abs) != wantParent {
		t.Fatalf("errors report dir drift: got %q, want %q",
			filepath.Dir(abs), wantParent)
	}
}
