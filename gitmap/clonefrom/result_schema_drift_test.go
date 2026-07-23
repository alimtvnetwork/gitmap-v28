package clonefrom

// Compile-time / reflect-time guard that keeps clonefrom.Result
// (and its embedded Row) in sync with the on-disk report schema
// emitted by WriteReport (CSV) and WriteReportJSON (JSON envelope
// of reportRowJSON).
//
// Why this test exists:
//
//   The CSV writer in summary.go (writeReportRows) and the JSON
//   row struct (reportRowJSON in summary.go, also written by
//   writeReportRowsJSON) reach into Result / Row by FIELD NAME.
//   If someone removes a field from Result -- or renames one --
//   the writers stop compiling and we'd notice immediately.
//   But if someone ADDS a field to Result (a new piece of per-
//   row state worth surfacing), nothing in the build forces them
//   to either:
//
//     a) wire the field into the report writers, or
//     b) explicitly mark the field as report-exempt.
//
//   The cmd-side helper /gitmap/cmd/clonefrom_reports.go just
//   delegates to WriteReport / WriteReportJSON, so the contract
//   that has to stay synchronized is REALLY:
//
//     Result + embedded Row  <-->  reportRowJSON / CSV header
//
//   This test enumerates the expected field set on each side
//   using REFLECT (so a rename breaks the build because the
//   referenced field literal stops resolving) and fails when
//   the two sets drift. When you add a Result/Row field you
//   must update exactly one of:
//
//     - reportedResultFields / reportedRowFields  (if it's
//       surfaced in the on-disk report), OR
//     - exemptResultFields  / exemptRowFields     (if it's
//       intentionally NOT surfaced -- include a one-line WHY).
//
//   The test is intentionally noisy on failure: it prints the
//   drift in both directions so the fix is mechanical.

import (
	"reflect"
	"sort"
	"testing"
	"time"
)

// reportedResultFields enumerates Result-level fields that the CSV
// and JSON report writers consume. Mirrors writeReportRows in
// summary.go and the reportRowJSON struct. Update both ends
// together when adding or removing a reported field.
var reportedResultFields = []string{
	"Row",      // embedded; checked separately against reportedRowFields
	"Dest",     // resolved dest path -- written as "dest" column / json
	"Status",   // ok | skipped | failed
	"Detail",   // human-readable context
	"Duration", // emitted as duration_seconds
}

// reportedRowFields enumerates Row-level fields the report writers
// reach into via Result.Row.<Field>. Anything else on Row is
// per-row INPUT state that the report intentionally omits (and is
// listed in exemptRowFields below with the rationale).
var reportedRowFields = []string{
	"URL",    // -> "url" column / "url" json
	"Branch", // -> "branch" column / "branch" json
	"Depth",  // -> "depth" column / "depth" json
}

// exemptResultFields lists Result fields that are intentionally
// NOT surfaced in the report. Empty today; kept as a hook so a
// future "internal tracing" field can opt out without forcing a
// schema bump. Each entry MUST carry a one-line "// why:" comment
// when added.
var exemptResultFields = []string{}

// exemptRowFields lists Row fields that are deliberately not
// surfaced in the on-disk report. Each entry carries the rationale
// inline so reviewers don't have to dig.
var exemptRowFields = []string{
	"Dest",     // why: report uses Result.Dest (RESOLVED path) instead.
	"Checkout", // why: per-row INPUT toggle, not an outcome to record.
}

// TestResult_FieldsMatchReportSchema is the drift guard. It builds
// the expected field set from (reported + exempt) on each struct
// and compares it against reflect.TypeOf(...).Field names. Any
// mismatch -- in either direction -- fails with a diff that names
// the missing / surplus fields.
func TestResult_FieldsMatchReportSchema(t *testing.T) {
	checkStructFields(t, "Result", reflect.TypeOf(Result{}),
		reportedResultFields, exemptResultFields)
	checkStructFields(t, "Row", reflect.TypeOf(Row{}),
		reportedRowFields, exemptRowFields)
}

// checkStructFields asserts that the actual exported field set on
// `typ` equals the union of `reported` and `exempt`. Reports any
// drift via t.Errorf with a precise diff so the fix is mechanical.
func checkStructFields(t *testing.T, name string, typ reflect.Type,
	reported, exempt []string) {
	t.Helper()
	want := unionSet(reported, exempt)
	got := exportedFieldNames(typ)
	if missing := diffSorted(want, got); len(missing) > 0 {
		t.Errorf("%s: declared in test but not on struct: %v "+
			"(remove from reported/exempt list, the field is gone)",
			name, missing)
	}
	if surplus := diffSorted(got, want); len(surplus) > 0 {
		t.Errorf("%s: new field(s) %v not declared in the report-"+
			"schema test. Add each to either reportedFields (if "+
			"the report writer surfaces it) or exemptFields (with "+
			"a one-line // why: rationale).", name, surplus)
	}
}

// exportedFieldNames returns the sorted set of exported field names
// on `typ`. Unexported fields are skipped because the report
// writers can't reach them anyway.
func exportedFieldNames(typ reflect.Type) []string {
	out := make([]string, 0, typ.NumField())
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		if f.IsExported() {
			out = append(out, f.Name)
		}
	}
	sort.Strings(out)

	return out
}

// unionSet returns the sorted, de-duplicated union of two string
// slices. Used to combine reported + exempt into the "expected"
// universe of field names for the drift comparison.
func unionSet(a, b []string) []string {
	seen := make(map[string]bool, len(a)+len(b))
	for _, s := range a {
		seen[s] = true
	}
	for _, s := range b {
		seen[s] = true
	}
	out := make([]string, 0, len(seen))
	for k := range seen {
		out = append(out, k)
	}
	sort.Strings(out)

	return out
}

// diffSorted returns elements of `want` not present in `got`.
// Both inputs MUST be sorted; result is sorted by construction.
func diffSorted(want, got []string) []string {
	have := make(map[string]bool, len(got))
	for _, s := range got {
		have[s] = true
	}
	out := make([]string, 0)
	for _, s := range want {
		if !have[s] {
			out = append(out, s)
		}
	}

	return out
}

// TestResult_ReportFieldTypesMatch pins the TYPE of every reported
// field so an accidental int->int64 / string->[]string change
// breaks the build (the literal field accesses below stop
// compiling) AND fails the test (the type assertion fires).
//
// This is the "compile-time" leg the user asked for: removing or
// renaming any of these fields makes the file fail to BUILD
// because the field selectors below resolve at compile time.
func TestResult_ReportFieldTypesMatch(t *testing.T) {
	r := Result{
		Row:      Row{URL: "u", Branch: "b", Depth: 1},
		Dest:     "d",
		Status:   "ok",
		Detail:   "x",
		Duration: time.Second,
	}
	// Compile-time field-existence checks: any rename / removal
	// of these field selectors fails `go build ./clonefrom/...`
	// before this test ever runs. Keep one selector per reported
	// field so the failure points at the offender.
	_ = r.Row.URL
	_ = r.Row.Branch
	_ = r.Row.Depth
	_ = r.Dest
	_ = r.Status
	_ = r.Detail
	_ = r.Duration
	// Runtime type guards mirror the CSV / reportRowJSON columns.
	// Update both ends together if any column ever changes type.
	assertType(t, "Row.URL", r.Row.URL, "")
	assertType(t, "Row.Branch", r.Row.Branch, "")
	assertType(t, "Row.Depth", r.Row.Depth, 0)
	assertType(t, "Dest", r.Dest, "")
	assertType(t, "Status", r.Status, "")
	assertType(t, "Detail", r.Detail, "")
	assertType(t, "Duration", r.Duration, time.Duration(0))
}

// assertType fails the test when `got` and `want` differ in
// reflect.Kind / Type. Tiny helper so the per-field assertions
// above stay one line each.
func assertType(t *testing.T, name string, got, want any) {
	t.Helper()
	if reflect.TypeOf(got) != reflect.TypeOf(want) {
		t.Errorf("%s: type drifted to %s (report writer expects %s)",
			name, reflect.TypeOf(got), reflect.TypeOf(want))
	}
}
