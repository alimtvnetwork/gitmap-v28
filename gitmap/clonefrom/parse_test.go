package clonefrom

// Parser tests covering both formats + the validation rules.
// Pure-function tests — no filesystem fixtures, every input is
// inline so a future refactor doesn't have to keep test data in
// sync with code in two places.

import (
	"path/filepath"
	"strings"
	"testing"
)

// TestParseFile_JSON exercises the JSON path with one row that
// hits every optional field, plus one minimal row to confirm
// defaults flow through.
func TestParseFile_JSON(t *testing.T) {
	path := writeTemp(t, "plan.json", `[
  {"url": "https://github.com/a/b.git", "dest": "bb", "branch": "main", "depth": 1},
  {"url": "git@github.com:c/d.git"}
]`)
	plan, err := ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}
	if plan.Format != "json" {
		t.Errorf("format = %q, want json", plan.Format)
	}
	if len(plan.Rows) != 2 {
		t.Fatalf("rows = %d, want 2", len(plan.Rows))
	}
	want0 := Row{URL: "https://github.com/a/b.git", Dest: "bb", Branch: "main", Depth: 1}
	if plan.Rows[0] != want0 {
		t.Errorf("row0 = %+v, want %+v", plan.Rows[0], want0)
	}
	if plan.Rows[1].URL != "git@github.com:c/d.git" {
		t.Errorf("row1 URL = %q", plan.Rows[1].URL)
	}
}

// TestParseFile_CSV exercises the CSV path with case-mismatched
// header, ragged rows, and a quoted field — the three things real
// spreadsheet exports throw at us.
func TestParseFile_CSV(t *testing.T) {
	body := "URL,Dest,Branch,Depth\n" +
		"https://github.com/a/b.git,,,\n" +
		"git@github.com:c/d.git,my-d\n" + // ragged: missing branch+depth
		"\"https://example.org/with,comma.git\",,main,5\n"
	path := writeTemp(t, "plan.csv", body)
	plan, err := ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}
	if len(plan.Rows) != 3 {
		t.Fatalf("rows = %d, want 3", len(plan.Rows))
	}
	if plan.Rows[1].Dest != "my-d" {
		t.Errorf("row1 dest = %q, want my-d", plan.Rows[1].Dest)
	}
	if plan.Rows[2].URL != "https://example.org/with,comma.git" {
		t.Errorf("row2 URL = %q (quoted-comma not preserved)", plan.Rows[2].URL)
	}
	if plan.Rows[2].Depth != 5 {
		t.Errorf("row2 depth = %d, want 5", plan.Rows[2].Depth)
	}
}

// TestParseFile_RejectsBadURL pins the error message shape so
// users grepping for `clone-from:` in CI logs find row-pointing
// failures consistently.
func TestParseFile_RejectsBadURL(t *testing.T) {
	path := writeTemp(t, "bad.csv", "url\nowner/repo\n")
	_, err := ParseFile(path)
	if err == nil {
		t.Fatal("expected error for bare owner/repo")
	}
	if !strings.Contains(err.Error(), "row 2") {
		t.Errorf("error %q does not point at row 2", err.Error())
	}
}

// TestParseFile_DedupMergesLaterFields confirms re-listing a URL
// with a more-specific branch overrides the earlier default —
// the documented spreadsheet workflow.
func TestParseFile_DedupMergesLaterFields(t *testing.T) {
	body := "url,dest,branch,depth\n" +
		"https://github.com/a/b.git,,,\n" +
		"https://github.com/a/b.git,,main,1\n"
	plan, err := ParseFile(writeTemp(t, "dup.csv", body))
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}
	if len(plan.Rows) != 1 {
		t.Fatalf("rows = %d, want 1 (deduped)", len(plan.Rows))
	}
	if plan.Rows[0].Branch != "main" || plan.Rows[0].Depth != 1 {
		t.Errorf("merged row = %+v, want branch=main depth=1", plan.Rows[0])
	}
}

// TestParseFile_MissingURLColumn pins the friendlier-than-default
// error when the CSV header has no `url`. Without this check the
// parser would accept the file and emit "url is empty" for every
// row, which is much harder to debug.
func TestParseFile_MissingURLColumn(t *testing.T) {
	_, err := ParseFile(writeTemp(t, "no-url.csv", "name,branch\nfoo,main\n"))
	if err == nil {
		t.Fatal("expected error for missing url column")
	}
	if !strings.Contains(err.Error(), "url") {
		t.Errorf("error %q does not mention url", err.Error())
	}
}

// writeTemp drops a fixture file in a per-test tempdir and
// returns the absolute path. Centralized so the parser tests
// stay declarative.
func writeTemp(t *testing.T, name, body string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, name)
	if err := writeFile(path, body); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}

	return path
}
