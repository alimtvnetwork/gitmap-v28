package cmd

// Byte-exact contract for `gitmap startup-list --format=jsonl`.
// Sibling to startuplistjson_bytes_contract_test.go (pretty JSON
// array). The two formats share buildStartupListJSONItems so any
// drift in field NAMES or ORDER fails BOTH tests in lockstep —
// catching a column rename in one diff instead of two.
//
// What this contract pins:
//
//   1. Empty input → ZERO bytes (NOT `\n`, NOT `[]\n`). JSONL's
//      line-oriented contract means `wc -l` must equal entry count.
//   2. Single entry → exactly one line, compact, trailing `\n`.
//   3. Multi-entry → N lines, every line independently parseable as
//      a JSON object, every line terminated by `\n` (including the
//      last) so concatenation of two outputs stays valid.
//   4. Key order within each line matches the JSON pretty contract:
//      name, path, exec — verified by parsing each line back into
//      an *ordered* decoder and asserting the key sequence.
//   5. Special characters are escaped identically to --format=json
//      (delegated to encoding/json via stablejson) — verified by a
//      cross-format round-trip rather than by a separate golden so
//      the escape contract has a single owner.

import (
	"bufio"
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/startup"
)

// TestStartupListJSONL_EmptyEmitsNothing asserts the empty-list
// contract: no bytes at all. A stray `\n` here would make `wc -l`
// report 1 for an empty list — silently breaking any pipeline that
// uses line count as a record count.
func TestStartupListJSONL_EmptyEmitsNothing(t *testing.T) {
	var buf bytes.Buffer
	if err := encodeStartupListJSONL(&buf, nil); err != nil {
		t.Fatalf("encode: %v", err)
	}
	if buf.Len() != 0 {
		t.Fatalf("empty list must emit zero bytes, got %d: %q",
			buf.Len(), buf.String())
	}
}

// TestStartupListJSONL_SingleEntryByteExact pins the exact bytes
// for a one-entry list. Catches drift in inter-key whitespace
// (compact must be `,"` not `, "`) and the trailing-newline rule.
func TestStartupListJSONL_SingleEntryByteExact(t *testing.T) {
	entries := []startup.Entry{
		{Name: "gitmap-a", Path: "/p/a.desktop", Exec: "/bin/a"},
	}
	var buf bytes.Buffer
	if err := encodeStartupListJSONL(&buf, entries); err != nil {
		t.Fatalf("encode: %v", err)
	}
	want := `{"name":"gitmap-a","path":"/p/a.desktop","exec":"/bin/a"}` + "\n"
	if got := buf.String(); got != want {
		t.Fatalf("byte mismatch\nwant: %q\ngot:  %q", want, got)
	}
}

// TestStartupListJSONL_MultiEntryLineCount asserts that a 3-entry
// list produces exactly 3 lines, each terminated by `\n`, and that
// each line is independently parseable. This is the core JSONL
// contract — a downstream consumer must be able to split on `\n`
// and feed each chunk to a JSON parser without any framing logic.
func TestStartupListJSONL_MultiEntryLineCount(t *testing.T) {
	entries := []startup.Entry{
		{Name: "gitmap-a", Path: "/p/a.desktop", Exec: "/bin/a"},
		{Name: "gitmap-b", Path: "/p/b.desktop", Exec: "/bin/b --flag"},
		{Name: "gitmap-c", Path: "/p/c.desktop", Exec: ""},
	}
	var buf bytes.Buffer
	if err := encodeStartupListJSONL(&buf, entries); err != nil {
		t.Fatalf("encode: %v", err)
	}
	if !bytes.HasSuffix(buf.Bytes(), []byte{'\n'}) {
		t.Fatalf("output must end with newline: %q", buf.String())
	}
	lines := strings.Split(strings.TrimSuffix(buf.String(), "\n"), "\n")
	if len(lines) != len(entries) {
		t.Fatalf("want %d lines, got %d: %q",
			len(entries), len(lines), buf.String())
	}
	for i, line := range lines {
		var got map[string]any
		if err := json.Unmarshal([]byte(line), &got); err != nil {
			t.Fatalf("line %d not valid JSON: %v\nline: %q", i, err, line)
		}
		if got["name"] != entries[i].Name {
			t.Fatalf("line %d name: want %q got %v",
				i, entries[i].Name, got["name"])
		}
	}
}

// TestStartupListJSONL_KeyOrderStable verifies the per-line key
// order matches the JSON pretty contract (name, path, exec). Uses
// json.Decoder.Token() to walk keys in stream order — a regular
// map[string]any decode would lose order.
func TestStartupListJSONL_KeyOrderStable(t *testing.T) {
	entries := []startup.Entry{
		{Name: "gitmap-a", Path: "/p/a.desktop", Exec: "/bin/a"},
		{Name: "gitmap-b", Path: "/p/b.desktop", Exec: "/bin/b"},
	}
	var buf bytes.Buffer
	if err := encodeStartupListJSONL(&buf, entries); err != nil {
		t.Fatalf("encode: %v", err)
	}
	wantKeys := assertSchemaKeysSlice(t, "startup-list")
	scanner := bufio.NewScanner(&buf)
	lineNo := 0
	for scanner.Scan() {
		assertJSONKeyOrder(t, lineNo, scanner.Bytes(), wantKeys)
		lineNo++
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("scan: %v", err)
	}
	if lineNo != len(entries) {
		t.Fatalf("want %d lines, scanned %d", len(entries), lineNo)
	}
}

// assertJSONKeyOrder walks one JSON object via Decoder.Token and
// confirms the keys appear in the given order. Split out so the
// test function above stays under the 15-line code-style budget.
func assertJSONKeyOrder(t *testing.T, lineNo int, line []byte, want []string) {
	t.Helper()
	dec := json.NewDecoder(bytes.NewReader(line))
	if _, err := dec.Token(); err != nil { // opening `{`
		t.Fatalf("line %d: open: %v", lineNo, err)
	}
	for i, wantKey := range want {
		tok, err := dec.Token()
		if err != nil {
			t.Fatalf("line %d: key %d: %v", lineNo, i, err)
		}
		if tok != wantKey {
			t.Fatalf("line %d: key %d want %q got %v",
				lineNo, i, wantKey, tok)
		}
		if _, err := dec.Token(); err != nil { // value
			t.Fatalf("line %d: value %d: %v", lineNo, i, err)
		}
	}
}

// TestStartupListJSONL_SpecialCharsMatchJSON cross-checks that the
// escape sequences in JSONL match those in --format=json for the
// same input. Same Field slice → same escaped values, since both
// formats route values through encoding/json via stablejson.
func TestStartupListJSONL_SpecialCharsMatchJSON(t *testing.T) {
	entries := []startup.Entry{
		{Name: "gitmap-中文", Path: "/p/\"q\"\\b.desktop", Exec: "a\tb\nc\u0001d"},
	}
	var jsonlBuf, jsonBuf bytes.Buffer
	if err := encodeStartupListJSONL(&jsonlBuf, entries); err != nil {
		t.Fatalf("jsonl encode: %v", err)
	}
	if err := encodeStartupListJSON(&jsonBuf, entries); err != nil {
		t.Fatalf("json encode: %v", err)
	}
	// Re-parse both and compare semantic content. We don't compare
	// bytes (formatting differs by design) — we compare that the
	// VALUES survive the round-trip identically in both encodings.
	var fromJSONL map[string]any
	jsonlLine := bytes.TrimSuffix(jsonlBuf.Bytes(), []byte{'\n'})
	if err := json.Unmarshal(jsonlLine, &fromJSONL); err != nil {
		t.Fatalf("jsonl reparse: %v", err)
	}
	var fromJSON []map[string]any
	if err := json.Unmarshal(jsonBuf.Bytes(), &fromJSON); err != nil {
		t.Fatalf("json reparse: %v", err)
	}
	if len(fromJSON) != 1 {
		t.Fatalf("want 1 json record, got %d", len(fromJSON))
	}
	for _, k := range assertSchemaKeysSlice(t, "startup-list") {
		if fromJSONL[k] != fromJSON[0][k] {
			t.Fatalf("key %q diverges: jsonl=%v json=%v",
				k, fromJSONL[k], fromJSON[0][k])
		}
	}
}
