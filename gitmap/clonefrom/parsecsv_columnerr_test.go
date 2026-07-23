package clonefrom

// Tests pinning the column-aware CSV row error format. A failing
// row must report the 1-indexed CSV row number AND the offending
// column name so an operator editing a large spreadsheet can jump
// directly to the bad cell instead of re-reading every field to
// guess which one tripped validation.
//
// Each subtest crafts a 2-row CSV (header + one data row), parses
// it, and asserts the returned error string contains both the row
// number and the expected column name. We assert on substrings (not
// the full message) so a future tightening of the inner validator
// wording does not force a fixture rewrite.

import (
	"strings"
	"testing"
)

func TestParseCSV_RowError_NamesOffendingColumn(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name      string
		csv       string
		wantRow   string // row number must appear as a literal substring
		wantCol   string // column name must appear (quoted) as a substring
		wantInErr string // a fragment from the inner validator message
	}{
		{
			name:      "bad depth",
			csv:       "url,depth\nhttps://x/y.git,not-a-number\n",
			wantRow:   "row 2",
			wantCol:   `"depth"`,
			wantInErr: "not a valid integer",
		},
		{
			name:      "bad checkout",
			csv:       "url,checkout\nhttps://x/y.git,bogus\n",
			wantRow:   "row 2",
			wantCol:   `"checkout"`,
			wantInErr: "auto",
		},
		{
			name:      "bad url",
			csv:       "url\nnot a url\n",
			wantRow:   "row 2",
			wantCol:   `"url"`,
			wantInErr: "does not look like a git URL",
		},
		{
			name:      "branch with leading dash",
			csv:       "url,branch\nhttps://x/y.git,-rf\n",
			wantRow:   "row 2",
			wantCol:   `"branch"`,
			wantInErr: "valid git ref",
		},
		{
			name:      "branch with whitespace",
			csv:       "url,branch\nhttps://x/y.git,my branch\n",
			wantRow:   "row 2",
			wantCol:   `"branch"`,
			wantInErr: "valid git ref",
		},
		{
			name:      "negative depth on row 3",
			csv:       "url,depth\nhttps://x/y.git,1\nhttps://a/b.git,-5\n",
			wantRow:   "row 3",
			wantCol:   `"depth"`,
			wantInErr: "negative",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := parseCSV(strings.NewReader(tc.csv))
			if err == nil {
				t.Fatalf("expected parse error, got nil")
			}
			msg := err.Error()
			for _, want := range []string{tc.wantRow, tc.wantCol, tc.wantInErr} {
				if !strings.Contains(msg, want) {
					t.Errorf("error %q missing substring %q", msg, want)
				}
			}
		})
	}
}

// TestParseCSV_ValidRows_NoError guards against the column-aware
// wrapper accidentally rejecting rows that should pass — every
// shape exercised by TestParseCSV_RowError above must have a
// passing counterpart so the validator rules cannot drift to
// "reject everything" without being noticed.
func TestParseCSV_ValidRows_NoError(t *testing.T) {
	t.Parallel()
	csv := "url,dest,branch,depth,checkout\n" +
		"https://x/y.git,out/y,main,1,auto\n" +
		"git@github.com:owner/repo.git,,feature/x,,skip\n"
	rows, err := parseCSV(strings.NewReader(csv))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("got %d rows, want 2", len(rows))
	}
}
