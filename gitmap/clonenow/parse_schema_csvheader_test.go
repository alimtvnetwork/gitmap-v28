package clonenow

// Robustness tests for CSV header parsing. Real-world scan exports
// (Excel "Save As CSV", PowerShell `Out-File -Encoding utf8`,
// hand-edited fixtures with stray spaces) routinely add cosmetic
// noise to the header row. None of those quirks should make
// clone-now reject a header whose column names are otherwise valid;
// only a TRULY unknown name (typo, wrong schema) must fail.
//
// Cases pinned here:
//   - UTF-8 BOM on the first column header (Excel / PowerShell)
//   - Double-quoted header names (some exporters quote unconditionally)
//   - Leading / trailing whitespace inside cells
//   - Trailing empty column (stray comma at end of header line)
//   - Plus a negative case: an actually-unknown column still fails.

import (
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// TestParseFile_CSVHeader_BOMTolerated verifies the UTF-8 BOM that
// Excel prefixes to "CSV UTF-8" exports does not turn the first
// header column into an unknown field.
func TestParseFile_CSVHeader_BOMTolerated(t *testing.T) {
	body := "\ufeffhttpsUrl,relativePath\r\n" +
		"https://x/a.git,a\r\n"
	path := writeTemp(t, ".csv", body)
	if _, err := ParseFile(path, "", constants.CloneNowModeHTTPS,
		constants.CloneNowOnExistsSkip); err != nil {
		t.Fatalf("ParseFile rejected BOM-prefixed header: %v", err)
	}
}

// TestParseFile_CSVHeader_QuotedNames verifies surrounding double-
// quotes around header names (a quoting habit of some exporters)
// are stripped before lookup.
func TestParseFile_CSVHeader_QuotedNames(t *testing.T) {
	// csv.Reader already unwraps standard quoted fields, so we use
	// inner literal quotes via doubled-quote escapes to simulate an
	// exporter that quoted the column NAMES themselves.
	body := "\"\"\"httpsUrl\"\"\",\"\"\"relativePath\"\"\"\r\n" +
		"https://x/a.git,a\r\n"
	path := writeTemp(t, ".csv", body)
	if _, err := ParseFile(path, "", constants.CloneNowModeHTTPS,
		constants.CloneNowOnExistsSkip); err != nil {
		t.Fatalf("ParseFile rejected quoted header names: %v", err)
	}
}

// TestParseFile_CSVHeader_PaddedNames verifies surrounding whitespace
// inside header cells does not trigger an unknown-column error.
func TestParseFile_CSVHeader_PaddedNames(t *testing.T) {
	body := "  httpsUrl  ,  relativePath  \r\n" +
		"https://x/a.git,a\r\n"
	path := writeTemp(t, ".csv", body)
	if _, err := ParseFile(path, "", constants.CloneNowModeHTTPS,
		constants.CloneNowOnExistsSkip); err != nil {
		t.Fatalf("ParseFile rejected padded header names: %v", err)
	}
}

// TestParseFile_CSVHeader_TrailingEmptyColumn verifies a stray
// trailing comma on the header line (producing one empty column) is
// tolerated. The corresponding empty payload column is harmless.
func TestParseFile_CSVHeader_TrailingEmptyColumn(t *testing.T) {
	body := "httpsUrl,relativePath,\r\n" +
		"https://x/a.git,a,\r\n"
	path := writeTemp(t, ".csv", body)
	if _, err := ParseFile(path, "", constants.CloneNowModeHTTPS,
		constants.CloneNowOnExistsSkip); err != nil {
		t.Fatalf("ParseFile rejected trailing-empty header column: %v", err)
	}
}

// TestParseFile_CSVHeader_TrulyUnknownStillFails is the negative
// control: after all the cosmetic tolerance above, a real typo
// MUST still fail loudly with the offending name + the known list.
func TestParseFile_CSVHeader_TrulyUnknownStillFails(t *testing.T) {
	body := "\ufeff  \"https_url\"  ,relativePath\r\n" +
		"https://x/a.git,a\r\n"
	path := writeTemp(t, ".csv", body)
	_, err := ParseFile(path, "", constants.CloneNowModeHTTPS,
		constants.CloneNowOnExistsSkip)
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
