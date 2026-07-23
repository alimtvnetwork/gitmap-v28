package cmd

// Report-writing helpers for `gitmap clone-from --execute`. Split
// from clonefrom.go to keep that file under the project's 200-line
// cap. Always writes CSV (unless --no-report); also writes JSON when
// --output terminal is set so the terminal summary can surface both
// paths.
//
// Pre-flight validation: before handing the result slice to either
// writer we call validateCloneFromResults. Missing required fields
// (URL, Dest, Status) on any row mean the report would be
// ambiguous (a CSV row with no URL/dest is just noise) so we
// refuse to write and surface a single actionable error instead.

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/cliexit"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/clonefrom"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// writeCloneFromReports persists the CSV report (always, unless
// --no-report) and additionally the JSON report when --output
// terminal is set. Returns ("", "") when --no-report skips both,
// or when the result set fails the pre-flight field check.
//
// Failures are logged to stderr but never abort — clones already
// happened, the reports are bonus.
func writeCloneFromReports(results []clonefrom.Result, cfg cloneFromFlags) (string, string) {
	if cfg.noReport {
		return "", ""
	}
	if err := validateCloneFromResults(results); err != nil {
		cliexit.Reportf(constants.CmdCloneFrom, "validate-results", cfg.file, err)

		return "", ""
	}
	csvPath := writeCloneFromCSV(results, cfg)
	if cfg.output != constants.OutputTerminal {
		return csvPath, ""
	}

	return csvPath, writeCloneFromJSON(results, cfg)
}

// writeCloneFromCSV invokes the CSV writer and reports the failure
// path (without aborting) on error. Extracted so writeCloneFromReports
// stays declarative and under the function-length cap.
func writeCloneFromCSV(results []clonefrom.Result, cfg cloneFromFlags) string {
	path, err := clonefrom.WriteReport(results)
	if err == nil {
		return path
	}
	cliexit.Reportf(constants.CmdCloneFrom, "write-csv-report", cfg.file, err)

	return ""
}

// writeCloneFromJSON invokes the JSON writer and reports the failure
// path (without aborting) on error. Mirror of writeCloneFromCSV so
// future changes to the failure-reporting shape happen in one place
// per format.
func writeCloneFromJSON(results []clonefrom.Result, cfg cloneFromFlags) string {
	path, err := clonefrom.WriteReportJSON(results)
	if err == nil {
		return path
	}
	cliexit.Reportf(constants.CmdCloneFrom, "write-json-report", cfg.file, err)

	return ""
}

// validateCloneFromResults walks the result set and returns nil when
// every row carries the fields required for a meaningful report
// (URL, Dest, Status). Otherwise returns one error that lists every
// offending row index together with the field name(s) it was missing
// — single error means a single stderr line at the call site.
func validateCloneFromResults(results []clonefrom.Result) error {
	indices, summaries := collectMissingResultFields(results)
	if len(indices) == 0 {
		return nil
	}

	return fmt.Errorf(constants.ErrCloneFromReportMissingFields,
		len(indices),
		strings.Join(indices, ","),
		strings.Join(summaries, "; "))
}

// collectMissingResultFields returns the row-index strings and the
// per-row missing-field summaries for any row that fails the field
// check. Two slices kept in lock-step (same length, same order) so
// the caller can format them into one message.
func collectMissingResultFields(results []clonefrom.Result) (indices, summaries []string) {
	for i, r := range results {
		missing := missingResultFields(r)
		if len(missing) == 0 {
			continue
		}
		indices = append(indices, fmt.Sprintf("%d", i))
		summaries = append(summaries, fmt.Sprintf("%d:%s", i, strings.Join(missing, "+")))
	}

	return indices, summaries
}

// missingResultFields lists the field names absent from a single
// result row. Order matches the canonical field list (URL, Dest,
// Status) so the test suite can assert on stable text.
func missingResultFields(r clonefrom.Result) []string {
	out := make([]string, 0, 3)
	if strings.TrimSpace(r.Row.URL) == "" {
		out = append(out, constants.CloneFromReportFieldURL)
	}
	if strings.TrimSpace(r.Dest) == "" {
		out = append(out, constants.CloneFromReportFieldDest)
	}
	if strings.TrimSpace(r.Status) == "" {
		out = append(out, constants.CloneFromReportFieldStatus)
	}

	return out
}
