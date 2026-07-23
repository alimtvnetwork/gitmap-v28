package clonenow

// Per-row error-message tests. Pin that schema validation reports
// the EXACT 1-based row / data-row number that failed so users can
// jump straight to the offending line in their manifest instead of
// hunting through a blob of "row N is bad" with no N.

import (
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// TestParseFile_JSONRowNotObject verifies a non-object element in
// the top-level array (here: a bare string at index 1) fails with
// the 1-based row number AND the observed JSON kind.
func TestParseFile_JSONRowNotObject(t *testing.T) {
	body := `[
		{"httpsUrl":"https://x/a.git","relativePath":"a"},
		"this is not an object",
		{"httpsUrl":"https://x/c.git","relativePath":"c"}
	]`
	path := writeTemp(t, ".json", body)
	_, err := ParseFile(path, "", constants.CloneNowModeHTTPS,
		constants.CloneNowOnExistsSkip)
	if err == nil {
		t.Fatal("ParseFile: want non-object-row error, got nil")
	}
	msg := err.Error()
	if !strings.Contains(msg, "row 2") {
		t.Errorf("error %q missing 'row 2' (1-based row number)", msg)
	}
	if !strings.Contains(msg, "string") {
		t.Errorf("error %q missing observed JSON kind 'string'", msg)
	}
}

// TestParseFile_JSONRowNumberSecondRow verifies the row index is
// truly per-row (not always 1) by placing the bad row second.
func TestParseFile_JSONRowNumberSecondRow(t *testing.T) {
	body := `[
		{"httpsUrl":"https://x/a.git","relativePath":"a"},
		{"repoName":"b","relativePath":"b"}
	]`
	path := writeTemp(t, ".json", body)
	_, err := ParseFile(path, "", constants.CloneNowModeHTTPS,
		constants.CloneNowOnExistsSkip)
	if err == nil {
		t.Fatal("ParseFile: want missing-url error, got nil")
	}
	if !strings.Contains(err.Error(), "row 2") {
		t.Errorf("error %q missing 'row 2'", err.Error())
	}
}

// TestParseFile_CSVDataRowMissingURL verifies a data row with empty
// URL columns is reported with its 1-based DATA row number (header
// is row 0; first data row = 1, matching what the user counts after
// the header line).
func TestParseFile_CSVDataRowMissingURL(t *testing.T) {
	body := "httpsUrl,sshUrl,relativePath\r\n" +
		"https://x/a.git,,a\r\n" +
		",,b\r\n" +
		"https://x/c.git,,c\r\n"
	path := writeTemp(t, ".csv", body)
	_, err := ParseFile(path, "", constants.CloneNowModeHTTPS,
		constants.CloneNowOnExistsSkip)
	if err == nil {
		t.Fatal("ParseFile: want missing-url-row error, got nil")
	}
	msg := err.Error()
	if !strings.Contains(msg, "data row 2") {
		t.Errorf("error %q missing 'data row 2'", msg)
	}
}

// TestParseFile_CSVDataRowMissingURLFirst pins that the first data
// row is reported as "data row 1" (not 0, not 2 — header is excluded
// from the count, but the first data row gets number 1).
func TestParseFile_CSVDataRowMissingURLFirst(t *testing.T) {
	body := "httpsUrl,sshUrl,relativePath\r\n" +
		",,a\r\n"
	path := writeTemp(t, ".csv", body)
	_, err := ParseFile(path, "", constants.CloneNowModeHTTPS,
		constants.CloneNowOnExistsSkip)
	if err == nil {
		t.Fatal("ParseFile: want missing-url-row error, got nil")
	}
	if !strings.Contains(err.Error(), "data row 1") {
		t.Errorf("error %q missing 'data row 1'", err.Error())
	}
}
