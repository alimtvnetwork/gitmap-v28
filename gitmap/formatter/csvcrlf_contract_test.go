package formatter

// CSV CRLF + comma contract test for formatter.WriteCSV (the scan-
// record emitter used by `gitmap scan` and friends). Sibling of
// gitmap/cmd/csvcrlf_contract_test.go which covers the cmd-package
// CSV emitters; together they pin every CSV-producing path in the
// codebase to RFC 4180 (comma separator, CRLF line endings).
//
// Rationale for CRLF: Excel and PowerShell-based downstream
// pipelines historically misinterpret bare-LF CSVs on Windows;
// forcing CRLF makes the bytes identical across Linux/macOS/Windows
// CI runs and matches the de-facto industry default.

import (
	"bytes"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// TestCSVCRLF_WriteCSV asserts CRLF + comma for the scan-record
// emitter across empty and populated inputs. Validation warnings are
// silenced via SetValidationSink so they don't pollute test output.
func TestCSVCRLF_WriteCSV(t *testing.T) {
	prev := SetValidationSink(&bytes.Buffer{})
	defer SetValidationSink(prev)

	cases := []struct {
		name    string
		records []model.ScanRecord
	}{
		{name: "empty", records: nil},
		{name: "populated", records: []model.ScanRecord{
			{
				Slug: "ok", RepoName: "ok",
				HTTPSUrl: "https://x/ok.git", RelativePath: "ok",
			},
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			if err := WriteCSV(&buf, tc.records); err != nil {
				t.Fatalf("WriteCSV: %v", err)
			}
			assertScanCSVCommaCRLF(t, buf.Bytes())
		})
	}
}

// assertScanCSVCommaCRLF mirrors the helper in cmd/csvcrlf_contract_test.go
// — duplicated rather than shared because the two test files live in
// different packages and the assertion is small / unlikely to drift.
// See that file for the full rule rationale.
func assertScanCSVCommaCRLF(t *testing.T, got []byte) {
	t.Helper()
	s := string(got)
	if !strings.Contains(s, "\r\n") {
		t.Fatalf("expected CRLF line endings, got none in: %q", s)
	}
	if hasBareLF(s) {
		t.Fatalf("found bare LF (not preceded by CR) — UseCRLF likely off: %q", s)
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
