package cmdlatestbranch

// CSV CRLF + comma contract test for the latest-branch CSV emitter:
//
//   - latest-branch --csv → encodeLatestBranchCSV
//
// The formatter package's WriteCSV (scan records) is covered by a
// sibling test in gitmap/formatter/csvcrlf_contract_test.go, and the
// startup-list CSV emitter is covered in cmdstartup.
//
// The contract is deliberately strict: every newline in the output
// — header AND every data row, including the trailing newline at
// end-of-file — must be CRLF ("\r\n"), and the field separator must
// be comma. This pins RFC 4180 conformance so downstream Excel /
// PowerShell / curl-based pipelines get byte-identical output on
// Linux, macOS, and Windows runs.
//
// NOTE: this test lives in cmdlatestbranch (not cmd) because the
// encoder (encodeLatestBranchCSV) is unexported here.

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
)

// TestCSVCRLF_LatestBranch asserts CRLF + comma for the latest-branch
// CSV emitter. Uses a deterministic single item so the test isn't
// time / format dependent — only line endings and separators are
// being validated here.
func TestCSVCRLF_LatestBranch(t *testing.T) {
	items := []gitutil.RemoteBranchInfo{
		{
			RemoteRef:  "refs/remotes/origin/main",
			CommitDate: time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC),
			Sha:        "abc1234567890",
			Subject:    "Initial commit",
		},
	}

	var buf bytes.Buffer
	if err := encodeLatestBranchCSV(&buf, items, "origin", 1); err != nil {
		t.Fatalf("encode: %v", err)
	}

	assertCSVCommaCRLF(t, buf.Bytes())
}

// assertCSVCommaCRLF is the shared assertion shape:
//
//  1. Output must contain at least one "\r\n" (proves CRLF is on).
//  2. Output must NOT contain any bare "\n" — i.e. every "\n" must
//     be preceded by "\r". A bare "\n" means encoding/csv reverted
//     to its default \n-only line ending and the contract is broken.
//  3. The header row (first line) must contain a comma — pins the
//     field separator to the RFC 4180 default. We only check the
//     header because data values can legitimately omit commas.
//
// Implemented as a single helper rather than three separate
// assertions so failure messages name the exact violated rule.
func assertCSVCommaCRLF(t *testing.T, got []byte) {
	t.Helper()
	s := string(got)
	if !strings.Contains(s, "\r\n") {
		t.Fatalf("expected CRLF line endings, got none in: %q", s)
	}

	if hasBareLF(s) {
		t.Fatalf("found bare LF (not preceded by CR) — encoding/csv UseCRLF likely off: %q", s)
	}

	header, _, ok := strings.Cut(s, "\r\n")
	if !ok {
		t.Fatalf("output missing CRLF-terminated header: %q", s)
	}

	if !strings.Contains(header, ",") {
		t.Fatalf("expected comma separator in header, got: %q", header)
	}
}

// hasBareLF returns true if any "\n" in s is not preceded by "\r".
// Walks the string once instead of allocating a regexp — keeps the
// test fast and the dependency surface zero.
func hasBareLF(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] != '\n' {
			continue
		}

		if i == 0 || s[i-1] != '\r' {
			return true
		}
	}

	return false
}
