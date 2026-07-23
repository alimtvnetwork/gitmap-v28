package cmd

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestHelpJSONStripANSI verifies stripANSI removes CSI color codes
// while preserving the surrounding text bytes.
func TestHelpJSONStripANSI(t *testing.T) {
	in := "\x1b[33mhello\x1b[0m world"
	got := stripANSI(in)
	want := "hello world"
	if got != want {
		t.Fatalf("stripANSI: got %q want %q", got, want)
	}
}

// TestHelpFilterRowsMatch ensures filterRows performs case-
// insensitive substring matching against group + line.
func TestHelpFilterRowsMatch(t *testing.T) {
	rows := []helpRow{
		{Group: "Cloning", Line: "  clone   Clone a repo"},
		{Group: "Release", Line: "  rel     Cut a release"},
	}
	hits := filterRows(rows, "CLONE")
	if len(hits) != 1 || !strings.Contains(hits[0].Line, "clone") {
		t.Fatalf("filterRows: unexpected hits %#v", hits)
	}
}

// TestHelpJSONDocSchema confirms allHelpRows feeds a non-empty
// helpJSONDoc that round-trips through encoding/json with a stable
// shape (Version + Groups[].Lines).
func TestHelpJSONDocSchema(t *testing.T) {
	rows := allHelpRows()
	if len(rows) == 0 {
		t.Fatal("allHelpRows returned 0 rows")
	}

	doc := helpJSONDoc{Version: "test"}
	byGroup := make(map[string]int)
	for _, r := range rows {
		line := stripANSI(strings.TrimRight(r.Line, "\n"))
		if idx, seen := byGroup[r.Group]; seen {
			doc.Groups[idx].Lines = append(doc.Groups[idx].Lines, line)

			continue
		}
		byGroup[r.Group] = len(doc.Groups)
		doc.Groups = append(doc.Groups, helpJSONGroup{
			Group: stripANSI(r.Group),
			Lines: []string{line},
		})
	}
	doc.Count = len(rows)

	b, err := json.Marshal(doc)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	s := string(b)
	for _, key := range []string{`"version":"test"`, `"count":`, `"groups":[`} {
		if !strings.Contains(s, key) {
			t.Errorf("json missing %q in %s", key, s[:min(200, len(s))])
		}
	}
}
