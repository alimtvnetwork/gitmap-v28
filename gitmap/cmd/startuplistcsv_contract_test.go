package cmd

// CSV column-stability contract test for `gitmap startup-list
// --format=csv`. Sibling to startuplistjson_contract_test.go (JSON)
// and startuplisttable_contract_test.go (table) — together they pin
// every consumer-facing surface so script authors can rely on a
// stable contract regardless of which format they pipe into awk /
// jq / Excel.
//
// What csvcrlf_contract_test.go ALREADY pins for us:
//   - CRLF line endings everywhere (RFC 4180 conformance)
//   - Comma separator in the header
//   - No bare LFs anywhere in the output
//
// What this file ADDS on top:
//   - Exact header line: "name,path,exec\r\n" (catches rename / reorder)
//   - Column count = 3 in header AND every data row (catches add/remove)
//   - Field order matches the JSON schema (cross-format stability so
//     a script consuming both formats sees the same field order)
//   - Empty input produces header-only output (1 row, 3 columns, no
//     data) — guarantees `tail -n +2` works on every output, even
//     when there are zero entries.
//
// To intentionally regenerate after a deliberate schema change,
// bump cmd/testdata/schemas/startup-list.vN.json (or run with
// GITMAP_UPDATE_SCHEMA=startup-list). The CSV header check below
// reads the same registry entry, so JSON+CSV cannot drift apart.

import (
	"bytes"
	"encoding/csv"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/startup"
)

// expectedStartupListCSVHeader is the EXACT first-line bytes scripts
// rely on. Any change here means a downstream consumer's column
// indexes change — bump the consumer-facing changelog before
// touching this constant.
const expectedStartupListCSVHeader = "name,path,exec\r\n"

// TestStartupListCSVContract_HeaderIsExact pins the header bytes
// verbatim. Catches: column rename ("exec" → "command"), column
// reorder ("name,path,exec" → "name,exec,path"), separator drift
// (comma → semicolon), line-ending drift (CRLF → LF).
func TestStartupListCSVContract_HeaderIsExact(t *testing.T) {
	var buf bytes.Buffer
	if err := encodeStartupListCSV(&buf, nil); err != nil {
		t.Fatalf("encode: %v", err)
	}
	got := buf.String()
	if !strings.HasPrefix(got, expectedStartupListCSVHeader) {
		t.Fatalf("header drift\n  want prefix: %q\n  got:         %q",
			expectedStartupListCSVHeader, got)
	}
}

// TestStartupListCSVContract_EmptyIsHeaderOnly guarantees the
// "always emit header even when zero entries" rule that lets
// scripts use `tail -n +2 | wc -l` for entry count without
// special-casing the empty list. The whole output must be exactly
// the header line — no extra blank line, no `null`, no nothing.
func TestStartupListCSVContract_EmptyIsHeaderOnly(t *testing.T) {
	var buf bytes.Buffer
	if err := encodeStartupListCSV(&buf, nil); err != nil {
		t.Fatalf("encode: %v", err)
	}
	if got := buf.String(); got != expectedStartupListCSVHeader {
		t.Fatalf("empty CSV must be header-only\n  want: %q\n  got:  %q",
			expectedStartupListCSVHeader, got)
	}
}

// TestStartupListCSVContract_EveryRowHas3Columns parses the output
// with encoding/csv (the same parser downstream Go consumers would
// use) and asserts every row — header AND data — has exactly 3
// fields. Catches: a value containing an unescaped comma slipping
// past csv.Writer's quoting (would surface as a >3-field row), or
// a future "summary row" being appended with a different shape.
func TestStartupListCSVContract_EveryRowHas3Columns(t *testing.T) {
	entries := []startup.Entry{
		{Name: "gitmap-a", Path: "/p/a.desktop", Exec: "/bin/a"},
		// Value with comma + quote — proves csv.Writer's quoting
		// keeps the row at exactly 3 fields when parsed back.
		{Name: "gitmap-b", Path: "/p/b.desktop", Exec: `/bin/b -arg "x,y"`},
		{Name: "gitmap-c", Path: "/p/c.desktop", Exec: ""},
	}
	var buf bytes.Buffer
	if err := encodeStartupListCSV(&buf, entries); err != nil {
		t.Fatalf("encode: %v", err)
	}
	rows := mustParseCSV(t, buf.Bytes())
	wantRows := 1 + len(entries) // header + data
	if len(rows) != wantRows {
		t.Fatalf("row count: want %d (header + %d data), got %d",
			wantRows, len(entries), len(rows))
	}
	for i, row := range rows {
		if len(row) != 3 {
			t.Errorf("row[%d] column count: want 3, got %d (%v)", i, len(row), row)
		}
	}
}

// TestStartupListCSVContract_HeaderMatchesJSONSchema is the
// cross-format stability check. CSV column names and order MUST
// match the JSON schema so a script consuming `--format=json | jq`
// and a script consuming `--format=csv | awk` see the same field
// names in the same positions. Without this test, the two schemas
// could silently drift apart.
func TestStartupListCSVContract_HeaderMatchesJSONSchema(t *testing.T) {
	var buf bytes.Buffer
	if err := encodeStartupListCSV(&buf, nil); err != nil {
		t.Fatalf("encode: %v", err)
	}
	rows := mustParseCSV(t, buf.Bytes())
	if len(rows) == 0 {
		t.Fatalf("empty parse — expected at least header row")
	}
	header := rows[0]
	wantHeader := assertSchemaKeysSlice(t, "startup-list")
	if !equalStringSlices(header, wantHeader) {
		t.Fatalf("CSV header drifted from JSON schema\n  json: %v\n  csv:  %v",
			wantHeader, header)
	}
}

// mustParseCSV runs encoding/csv.Reader over `raw` and returns all
// rows. Failures fatal-out with a labeled error so a malformed
// emitter doesn't bubble up as a confusing nil-deref later.
func mustParseCSV(t *testing.T, raw []byte) [][]string {
	t.Helper()
	r := csv.NewReader(bytes.NewReader(raw))
	rows, err := r.ReadAll()
	if err != nil {
		t.Fatalf("csv parse: %v\nraw: %q", err, string(raw))
	}

	return rows
}
