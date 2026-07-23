package cmd

// clonetermverifytestbanner.go — test-only convention helper that
// brackets a mismatch report with a sentinel banner so log readers
// can tell INTENTIONAL test-driven mismatches apart from REAL
// production verifier failures.
//
// Why this lives in non-_test.go: it's used by both the unit tests
// in clonetermverify_test.go and the integration tests in
// clonestream_integration_test.go. Putting it in a _test.go file
// would scope it to a single test binary and force a duplicate, so
// we publish it as an internal helper and document the "tests only"
// rule in the function comment + the file header.
//
// Why a separate file (vs. inline in clonetermverify.go): the
// production printer is at the project's 200-line ceiling, and the
// test wrapper is a different concern (developer-experience banner
// vs. user-facing report format). Splitting keeps both files small
// and makes the test-only intent obvious from the filename alone.

import (
	"io"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// PrintCmdFaithfulReportForTest wraps PrintCmdFaithfulReport in a
// "--- expected mismatch ---" banner so test runs that intentionally
// drive a divergent input don't bury real `[FAIL]` lines in identical-
// looking simulated ones. Production code MUST NOT call this — only
// tests that deliberately exercise the mismatch print path. The
// banner emits even when the report is empty so the section is always
// paired open/close in captured output, making golden diffs and human
// scans deterministic.
func PrintCmdFaithfulReportForTest(w io.Writer, r CmdFaithfulReport) error {
	if _, err := io.WriteString(w, constants.CmdFaithfulReportTestPrefix); err != nil {
		return err
	}
	if err := PrintCmdFaithfulReport(w, r); err != nil {
		return err
	}
	if _, err := io.WriteString(w, constants.CmdFaithfulReportTestSuffix); err != nil {
		return err
	}

	return nil
}
