package clonenow

// Schema validation tests for clone-now JSON and CSV inputs.
//
// These pin the contract that clone-now rejects malformed manifests
// loudly (clear error mentioning the offending field/row + the full
// list of known fields) instead of silently dropping rows. See
// parse_schema.go for the rationale.

import (
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// TestParseFile_JSONUnknownField verifies that a typo'd JSON key
// (e.g. "https_url" instead of "httpsUrl") fails ParseFile with a
// message that names the offending field and lists known fields.
func TestParseFile_JSONUnknownField(t *testing.T) {
	body := `[{"https_url":"https://x/a.git","relativePath":"a"}]`
	path := writeTemp(t, ".json", body)
	_, err := ParseFile(path, "", constants.CloneNowModeHTTPS, constants.CloneNowOnExistsSkip)
	if err == nil {
		t.Fatal("ParseFile: want unknown-field error, got nil")
	}
	msg := err.Error()
	if !strings.Contains(msg, "https_url") {
		t.Errorf("error %q missing offending field name", msg)
	}
	if !strings.Contains(msg, "httpsUrl") {
		t.Errorf("error %q missing known-field list", msg)
	}
}

// TestParseFile_JSONMissingURL verifies that a JSON row with neither
// httpsUrl nor sshUrl is rejected with a row-numbered error.
func TestParseFile_JSONMissingURL(t *testing.T) {
	body := `[{"repoName":"a","relativePath":"a"}]`
	path := writeTemp(t, ".json", body)
	_, err := ParseFile(path, "", constants.CloneNowModeHTTPS, constants.CloneNowOnExistsSkip)
	if err == nil {
		t.Fatal("ParseFile: want missing-url error, got nil")
	}
	if !strings.Contains(err.Error(), "row 1") {
		t.Errorf("error %q missing row index", err.Error())
	}
}

// TestParseFile_JSONNotArray verifies that a top-level JSON object
// (instead of an array) is rejected with a shape error.
func TestParseFile_JSONNotArray(t *testing.T) {
	body := `{"httpsUrl":"https://x/a.git","relativePath":"a"}`
	path := writeTemp(t, ".json", body)
	_, err := ParseFile(path, "", constants.CloneNowModeHTTPS, constants.CloneNowOnExistsSkip)
	if err == nil {
		t.Fatal("ParseFile: want shape error, got nil")
	}
	if !strings.Contains(err.Error(), "array") {
		t.Errorf("error %q missing 'array' phrasing", err.Error())
	}
}

// TestParseFile_CSVUnknownColumn verifies that a CSV with a header
// column not in ScanRecord's tag set is rejected loudly.
func TestParseFile_CSVUnknownColumn(t *testing.T) {
	body := "repoName,https_url,sshUrl,branch,branchSource,relativePath,absolutePath,cloneInstruction,notes,depth\r\n" +
		"a,https://x/a.git,,main,HEAD,a,/abs/a,git clone https://x/a.git a,,0\r\n"
	path := writeTemp(t, ".csv", body)
	_, err := ParseFile(path, "", constants.CloneNowModeHTTPS, constants.CloneNowOnExistsSkip)
	if err == nil {
		t.Fatal("ParseFile: want unknown-column error, got nil")
	}
	msg := err.Error()
	if !strings.Contains(msg, "https_url") {
		t.Errorf("error %q missing offending column name", msg)
	}
	if !strings.Contains(msg, "httpsUrl") {
		t.Errorf("error %q missing known-column list", msg)
	}
}

// TestParseFile_CSVMissingURLColumn verifies that a CSV header with
// neither httpsUrl nor sshUrl is rejected.
func TestParseFile_CSVMissingURLColumn(t *testing.T) {
	body := "repoName,branch,relativePath\r\n" +
		"a,main,a\r\n"
	path := writeTemp(t, ".csv", body)
	_, err := ParseFile(path, "", constants.CloneNowModeHTTPS, constants.CloneNowOnExistsSkip)
	if err == nil {
		t.Fatal("ParseFile: want missing-url-column error, got nil")
	}
	if !strings.Contains(err.Error(), "httpsUrl") {
		t.Errorf("error %q missing httpsUrl mention", err.Error())
	}
}

// TestParseFile_CSVEmpty verifies that a completely empty CSV file
// is rejected with the empty-csv error rather than the generic
// zero-rows error (so the user knows the FILE is bad, not just empty).
func TestParseFile_CSVEmpty(t *testing.T) {
	path := writeTemp(t, ".csv", "")
	_, err := ParseFile(path, "", constants.CloneNowModeHTTPS, constants.CloneNowOnExistsSkip)
	if err == nil {
		t.Fatal("ParseFile: want empty-csv error, got nil")
	}
	if !strings.Contains(err.Error(), "empty") {
		t.Errorf("error %q missing 'empty' phrasing", err.Error())
	}
}

// TestParseFile_JSONAllKnownFieldsAccepted is a positive-control
// test: a row carrying every documented field must parse cleanly.
// This guards against the validator being over-eager and rejecting
// the very fields that gitmap scan emits.
func TestParseFile_JSONAllKnownFieldsAccepted(t *testing.T) {
	body := `[{
		"id": 1, "slug": "a", "repoId": "r", "repoName": "a",
		"httpsUrl": "https://x/a.git", "sshUrl": "git@x:a.git",
		"discoveredUrl": "https://x/a.git",
		"branch": "main", "branchSource": "HEAD",
		"relativePath": "a", "absolutePath": "/abs/a",
		"cloneInstruction": "git clone https://x/a.git a",
		"notes": "", "depth": 0
	}]`
	path := writeTemp(t, ".json", body)
	plan, err := ParseFile(path, "", constants.CloneNowModeHTTPS, constants.CloneNowOnExistsSkip)
	if err != nil {
		t.Fatalf("ParseFile rejected scan-shaped input: %v", err)
	}
	if len(plan.Rows) != 1 {
		t.Fatalf("rows = %d, want 1", len(plan.Rows))
	}
}
