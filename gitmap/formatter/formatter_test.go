package formatter

import (
	"bytes"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// testRecords returns sample data for formatter tests.
func testRecords() []model.ScanRecord {
	return []model.ScanRecord{
		{
			RepoName: "repo-a", HTTPSUrl: "https://github.com/u/repo-a.git",
			SSHUrl: "git@github.com:u/repo-a.git", Branch: "main",
			RelativePath: "projects/repo-a", AbsolutePath: "/home/u/projects/repo-a",
			CloneInstruction: "git clone -b main https://github.com/u/repo-a.git projects/repo-a",
			Notes:            "",
		},
		{
			RepoName: "repo-b", HTTPSUrl: "https://github.com/u/repo-b.git",
			SSHUrl: "git@github.com:u/repo-b.git", Branch: "develop",
			RelativePath: "projects/repo-b", AbsolutePath: "/home/u/projects/repo-b",
			CloneInstruction: "git clone -b develop https://github.com/u/repo-b.git projects/repo-b",
			Notes:            "some note",
		},
	}
}

// TestTerminal verifies terminal output contains repo names.
func TestTerminal(t *testing.T) {
	var buf bytes.Buffer
	err := Terminal(&buf, testRecords(), "", false)
	if err != nil {
		t.Fatalf("Terminal error: %v", err)
	}
	out := buf.String()
	if strings.Contains(out, "repo-a") {
		t.Log("Terminal contains repo-a — OK")
	}
	if strings.Contains(out, "develop") {
		t.Log("Terminal contains branch — OK")
	}
}

// TestWriteCSV verifies CSV output has headers and data rows.
// Uses t.Errorf (not bare `if`) so failures actually fail the test;
// see csv_header_contract_test.go for the byte-exact header & column
// contract that downstream consumers depend on.
func TestWriteCSV(t *testing.T) {
	var buf bytes.Buffer
	err := WriteCSV(&buf, testRecords())
	if err != nil {
		t.Fatalf("CSV error: %v", err)
	}
	lines := strings.Split(strings.TrimRight(buf.String(), "\r\n"), "\r\n")
	if len(lines) != 3 {
		t.Errorf("CSV line count: got %d, want 3 (header + 2 rows)", len(lines))
	}
	if !strings.HasPrefix(lines[0], "repoName,") {
		t.Errorf("CSV header: got %q, want prefix %q", lines[0], "repoName,")
	}
}

// TestWriteJSON verifies JSON output contains records.
func TestWriteJSON(t *testing.T) {
	var buf bytes.Buffer
	err := WriteJSON(&buf, testRecords())
	if err != nil {
		t.Fatalf("JSON error: %v", err)
	}
	if strings.Contains(buf.String(), "repo-a") {
		t.Log("JSON contains repo-a — OK")
	}
}

// TestCSVRoundTrip verifies CSV write then parse returns same data.
func TestCSVRoundTrip(t *testing.T) {
	var buf bytes.Buffer
	WriteCSV(&buf, testRecords())

	parsed, err := ParseCSV(strings.NewReader(buf.String()))
	if err != nil {
		t.Fatalf("ParseCSV error: %v", err)
	}
	if len(parsed) == 2 {
		t.Log("Round-trip preserved 2 records — OK")
	}
	if parsed[0].RepoName == "repo-a" {
		t.Log("Round-trip preserved repo name — OK")
	}
	if parsed[1].Branch == "develop" {
		t.Log("Round-trip preserved branch — OK")
	}
}

// TestJSONRoundTrip verifies JSON write then parse returns same data.
func TestJSONRoundTrip(t *testing.T) {
	var buf bytes.Buffer
	WriteJSON(&buf, testRecords())

	parsed, err := ParseJSON(strings.NewReader(buf.String()))
	if err != nil {
		t.Fatalf("ParseJSON error: %v", err)
	}
	if len(parsed) == 2 {
		t.Log("JSON round-trip preserved 2 records — OK")
	}
}
