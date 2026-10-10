package cmdautofix

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ---------------------------------------------------------------------------
// Reports
// ---------------------------------------------------------------------------

const maxHumanViolations = 50

// categoryStat is one row of the mandatory summary table.
type categoryStat struct {
	name         string
	filesFlagged int
	fixes        int
	reportOnly   bool
}

func summarizeCategories(result *ScanResult) []categoryStat {
	byCat := map[string]*categoryStat{}
	seen := map[string]map[string]bool{}
	order := []string{}
	for _, v := range result.Violations {
		st, ok := byCat[v.Category]
		if !ok {
			st = &categoryStat{name: v.Category, reportOnly: !categoryHasFix(v.Category)}
			byCat[v.Category] = st
			seen[v.Category] = map[string]bool{}
			order = append(order, v.Category)
		}
		st.fixes++
		seen[v.Category][v.Path] = true
	}
	out := make([]categoryStat, 0, len(order))
	for _, name := range order {
		st := byCat[name]
		st.filesFlagged = len(seen[name])
		out = append(out, *st)
	}
	return out
}

// printScanSummary prints the MANDATORY summary: files scanned, files
// modified (0 before apply), per-category fix counts, time taken — always,
// even when zero findings.
func printScanSummary(opts Options, result *ScanResult) {
	fmt.Printf("Scan complete in %.2fs — %d files scanned (%d from cache).\n",
		result.Elapsed.Seconds(), result.FilesScanned, result.FilesFromCache)
	fmt.Printf("%-11s %-14s %s\n", "CATEGORY", "FILES FLAGGED", "FIXES")
	totalFiles, totalFixes := 0, 0
	for _, st := range summarizeCategories(result) {
		fixes := fmt.Sprintf("%d", st.fixes)
		if st.reportOnly {
			fixes = "0 (report-only)"
		}
		fmt.Printf("%-11s %-14d %s\n", st.name, st.filesFlagged, fixes)
		totalFiles += st.filesFlagged
		totalFixes += st.fixes
	}
	fmt.Printf("%-11s %-14d %d\n", "TOTAL", totalFiles, totalFixes)
}

func printViolationList(result *ScanResult) {
	shown := result.Violations
	truncated := 0
	if len(shown) > maxHumanViolations {
		truncated = len(shown) - maxHumanViolations
		shown = shown[:maxHumanViolations]
	}
	for _, v := range shown {
		loc := v.Path
		if v.Line > 0 {
			loc = fmt.Sprintf("%s:%d", v.Path, v.Line)
		}
		fmt.Printf("  %s [%s] %s\n", loc, v.Category, v.Detail)
	}
	if truncated > 0 {
		fmt.Printf("  ... and %d more (use --json for the full list)\n", truncated)
	}
}

// jsonFixReport is the --json machine-readable report shape.
type jsonFixReport struct {
	Command        string        `json:"command"`
	Subcommand     string        `json:"subcommand"`
	Path           string        `json:"path"`
	FilesScanned   int           `json:"files_scanned"`
	FilesFromCache int           `json:"files_from_cache"`
	FilesModified  int           `json:"files_modified"`
	ScanSeconds    float64       `json:"scan_seconds"`
	ApplySeconds   float64       `json:"apply_seconds"`
	Categories     []jsonCatStat `json:"categories"`
	Violations     []Violation   `json:"violations"`
	FixedFiles     []string      `json:"fixed_files"`
	ExitCode       int           `json:"exit_code"`
}

type jsonCatStat struct {
	Name         string `json:"name"`
	FilesFlagged int    `json:"files_flagged"`
	Fixes        int    `json:"fixes"`
	ReportOnly   bool   `json:"report_only"`
}

func printJSONReport(opts Options, result *ScanResult, applied *ApplyResult) {
	exitCode := 0
	if len(result.Violations) > 0 {
		exitCode = 1
	}
	catStats := []jsonCatStat{}
	for _, st := range summarizeCategories(result) {
		catStats = append(catStats, jsonCatStat{
			Name: st.name, FilesFlagged: st.filesFlagged,
			Fixes: st.fixes, ReportOnly: st.reportOnly,
		})
	}
	violations := result.Violations
	if violations == nil {
		violations = []Violation{}
	}
	report := jsonFixReport{
		Command:        "fix",
		Subcommand:     subcommandName(opts),
		Path:           absPathOf(opts.Root),
		FilesScanned:   result.FilesScanned,
		FilesFromCache: result.FilesFromCache,
		ScanSeconds:    round2(result.Elapsed.Seconds()),
		Categories:     catStats,
		Violations:     violations,
		FixedFiles:     []string{},
		ExitCode:       exitCode,
	}
	if applied != nil {
		report.FilesModified = applied.FilesModified
		report.ApplySeconds = round2(applied.Elapsed.Seconds())
		report.FixedFiles = applied.ModifiedFiles
		if report.FixedFiles == nil {
			report.FixedFiles = []string{}
		}
		if !unfixableRemain(result, applied) {
			report.ExitCode = 0
		}
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(report)
}

func subcommandName(opts Options) string {
	if len(opts.Categories) == len(categoryRegistry) {
		return "all"
	}
	if len(opts.Categories) == 1 {
		return opts.Categories[0]
	}
	return strings.Join(opts.Categories, ",")
}

func absPathOf(root string) string {
	abs, err := filepath.Abs(root)
	if err != nil {
		return root
	}
	return abs
}

func round2(f float64) float64 {
	return float64(int(f*100+0.5)) / 100
}
