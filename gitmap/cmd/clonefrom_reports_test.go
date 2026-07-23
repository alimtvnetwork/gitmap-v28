package cmd

// Unit tests for the pre-flight result validator added to
// writeCloneFromReports. The validator's contract:
//
//   - empty result set is OK (writers handle len==0)
//   - every fully-populated row passes
//   - any missing required field (URL / Dest / Status) on any row
//     produces a single error that names every offending row
//     index plus which fields were missing
//
// We assert against the constants the production code emits so a
// future rename has to update both sides.

import (
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/clonefrom"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// TestValidateCloneFromResults_PassesWhenComplete locks the happy
// path: nil error when every row carries URL+Dest+Status.
func TestValidateCloneFromResults_PassesWhenComplete(t *testing.T) {
	t.Parallel()
	results := []clonefrom.Result{
		newCloneFromResult("https://x/a.git", "a", constants.CloneFromStatusOK),
		newCloneFromResult("https://x/b.git", "b", constants.CloneFromStatusFailed),
	}
	if err := validateCloneFromResults(results); err != nil {
		t.Fatalf("expected nil for complete results, got: %v", err)
	}
}

// TestValidateCloneFromResults_EmptyIsOK guards against a
// regression that would refuse a zero-row report (the executor can
// legitimately produce an empty slice for an empty manifest).
func TestValidateCloneFromResults_EmptyIsOK(t *testing.T) {
	t.Parallel()
	if err := validateCloneFromResults(nil); err != nil {
		t.Fatalf("empty result set should validate, got: %v", err)
	}
}

// TestValidateCloneFromResults_FlagsMissingFields drives the
// validator with a deliberately-broken row and asserts the error
// names the row index plus each missing field.
func TestValidateCloneFromResults_FlagsMissingFields(t *testing.T) {
	t.Parallel()
	results := []clonefrom.Result{
		newCloneFromResult("https://x/a.git", "a", constants.CloneFromStatusOK),
		// row 1: missing url + status
		{Row: clonefrom.Row{URL: "  "}, Dest: "b", Status: ""},
		// row 2: missing dest
		{Row: clonefrom.Row{URL: "https://x/c.git"}, Dest: "", Status: constants.CloneFromStatusFailed},
	}
	err := validateCloneFromResults(results)
	if err == nil {
		t.Fatalf("expected validation error, got nil")
	}
	msg := err.Error()
	for _, want := range []string{
		"rows=[1,2]",
		"1:" + constants.CloneFromReportFieldURL + "+" + constants.CloneFromReportFieldStatus,
		"2:" + constants.CloneFromReportFieldDest,
	} {
		if !strings.Contains(msg, want) {
			t.Fatalf("error missing %q\nfull error: %s", want, msg)
		}
	}
}

// newCloneFromResult is a tiny constructor so the table cases above
// stay readable. Duration is left zero — the validator does not
// inspect it.
func newCloneFromResult(url, dest, status string) clonefrom.Result {
	return clonefrom.Result{
		Row:    clonefrom.Row{URL: url},
		Dest:   dest,
		Status: status,
	}
}
