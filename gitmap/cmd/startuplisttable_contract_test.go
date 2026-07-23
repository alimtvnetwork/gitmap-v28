package cmd

// Table (default human-readable) format contract test for
// `gitmap startup-list`. Sibling to startuplistjson_contract_test.go
// (JSON) and startuplistcsv_contract_test.go (CSV) — together they
// pin every consumer-facing surface so script authors who grep the
// table output can rely on the bullet+arrow shape staying stable.
//
// Why test the human-readable format at all?
//
// Plenty of shell scripts in the wild parse `--help`-style output
// because that's the path of least resistance ("just pipe to grep").
// gitmap's startup-list table has been stable since v3.133.0 and
// downstream automation has started to depend on:
//
//   - The header line being the first line and including the
//     directory path (so scripts can confirm which autostart dir
//     was scanned).
//   - The empty-state literal `(none — no gitmap-managed autostart
//     entries found)` (so scripts can branch on "nothing to do"
//     without parsing the count line).
//   - The per-entry shape `  • <name>  →  <exec>` (so scripts can
//     extract names with a single awk pattern).
//   - The footer line `Total: N entry(ies).` (so scripts can sanity-
//     check the count without re-counting bullet lines).
//
// Pinning these here means any future i18n / styling change to
// constants.MsgStartupList* fires a clear test failure naming the
// exact line that drifted, with a pointer to bump the consumer-
// facing changelog before regenerating.

import (
	"bytes"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/startup"
)

// TestStartupListTableContract_HeaderIncludesDir pins the first line
// shape: it must start with "Linux/Unix autostart entries managed
// by gitmap" and contain the scanned directory in parentheses. This
// is the exact substring scripts grep to confirm scope.
func TestStartupListTableContract_HeaderIncludesDir(t *testing.T) {
	got := mustRenderTable(t, "/home/user/.config/autostart", nil)
	first := firstLine(got)
	const wantPrefix = "Linux/Unix autostart entries managed by gitmap"
	if !strings.HasPrefix(first, wantPrefix) {
		t.Errorf("header prefix drift\n  want prefix: %q\n  got line:    %q", wantPrefix, first)
	}
	if !strings.Contains(first, "/home/user/.config/autostart") {
		t.Errorf("header missing scanned dir, got: %q", first)
	}
}

// TestStartupListTableContract_EmptyMarkerIsStable pins the
// machine-grep-able literal scripts use to detect "no entries"
// without parsing the count line. Any change to the parenthesized
// text breaks downstream automation — the test will fail loudly.
func TestStartupListTableContract_EmptyMarkerIsStable(t *testing.T) {
	got := mustRenderTable(t, "/x", nil)
	const wantMarker = "(none — no gitmap-managed autostart entries found)"
	if !strings.Contains(got, wantMarker) {
		t.Errorf("empty-state marker drift\n  want substring: %q\n  got: %q", wantMarker, got)
	}
}

// TestStartupListTableContract_RowShape pins the per-entry shape
// `  • <name>  →  <exec>`. Asserts every entry produces exactly one
// line matching that prefix and containing the arrow separator.
// Catches: bullet character change, indentation change, separator
// change ("→" → "->"), or a newline being inserted mid-row.
func TestStartupListTableContract_RowShape(t *testing.T) {
	entries := []startup.Entry{
		{Name: "gitmap-a", Path: "/p/a.desktop", Exec: "/bin/a"},
		{Name: "gitmap-b", Path: "/p/b.desktop", Exec: "/bin/b --flag"},
		{Name: "gitmap-c", Path: "/p/c.desktop", Exec: ""},
	}
	got := mustRenderTable(t, "/x", entries)
	rows := matchingRows(got, "  • ")
	if len(rows) != len(entries) {
		t.Fatalf("row count: want %d, got %d\n  output: %q", len(entries), len(rows), got)
	}
	for i, row := range rows {
		if !strings.Contains(row, "  →  ") {
			t.Errorf("row[%d] missing arrow separator: %q", i, row)
		}
		if !strings.Contains(row, entries[i].Name) {
			t.Errorf("row[%d] missing name %q: %q", i, entries[i].Name, row)
		}
	}
}

// TestStartupListTableContract_EmptyExecRendersPlaceholder verifies
// that an entry with no Exec line shows the documented placeholder
// "(no Exec line)" instead of the entry being silently dropped or
// the row showing as `name  →  ` with nothing after the arrow.
// Scripts that count entries by counting bullet lines depend on
// this — a missing row would throw their accounting off.
func TestStartupListTableContract_EmptyExecRendersPlaceholder(t *testing.T) {
	entries := []startup.Entry{{Name: "gitmap-c", Path: "/p/c.desktop", Exec: ""}}
	got := mustRenderTable(t, "/x", entries)
	if !strings.Contains(got, "(no Exec line)") {
		t.Errorf("missing Exec placeholder for empty Exec, got: %q", got)
	}
}

// TestStartupListTableContract_FooterPinsCount pins the footer
// shape `Total: N entry(ies).` for both populated and edge counts.
// Scripts that assert "expected at least 5 entries" use this line
// directly via grep + awk. Including the literal "entry(ies)"
// catches well-meaning future grammar fixes ("entries" / "1 entry").
func TestStartupListTableContract_FooterPinsCount(t *testing.T) {
	entries := []startup.Entry{
		{Name: "a", Path: "/p/a.desktop", Exec: "/bin/a"},
		{Name: "b", Path: "/p/b.desktop", Exec: "/bin/b"},
		{Name: "c", Path: "/p/c.desktop", Exec: "/bin/c"},
	}
	got := mustRenderTable(t, "/x", entries)
	const want = "Total: 3 entry(ies)."
	if !strings.Contains(got, want) {
		t.Errorf("footer drift\n  want substring: %q\n  got: %q", want, got)
	}
}

// mustRenderTable captures renderStartupListTable's output into a
// buffer. Relies on the io.Writer parameter added in v3.154.0 — the
// table renderer used to hardcode os.Stdout, which is why this
// contract file didn't exist before. Returns the buffer as a string
// since every assertion below is substring / line based.
func mustRenderTable(t *testing.T, dir string, entries []startup.Entry) string {
	t.Helper()
	var buf bytes.Buffer
	if err := renderStartupListTable(&buf, dir, entries); err != nil {
		t.Fatalf("renderStartupListTable: %v", err)
	}

	return buf.String()
}

// firstLine returns the first newline-terminated line of s without
// the trailing "\n". Used by header assertions where only the first
// line matters and trailing-content drift is covered elsewhere.
func firstLine(s string) string {
	if idx := strings.IndexByte(s, '\n'); idx >= 0 {

		return s[:idx]
	}

	return s
}

// matchingRows returns every line in s that starts with `prefix`
// (after splitting on "\n"). Used by the row-shape test to extract
// just the bullet lines and ignore the header / footer.
func matchingRows(s, prefix string) []string {
	var out []string
	for _, line := range strings.Split(s, "\n") {
		if strings.HasPrefix(line, prefix) {
			out = append(out, line)
		}
	}

	return out
}
