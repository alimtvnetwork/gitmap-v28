package clonefrom

// Tests for the `--output terminal` summary block + URL-scheme
// classifier added in summary_terminal.go / summary_scheme.go.
// Pure string-assembly tests — no filesystem, no git.

import (
	"bytes"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// TestClassifyScheme_AllBuckets pins one URL per bucket so a future
// refactor that drops a prefix from the table fails loudly.
func TestClassifyScheme_AllBuckets(t *testing.T) {
	cases := []struct{ url, want string }{
		{"https://github.com/a/b.git", constants.CloneFromSchemeHTTPS},
		{"http://example.org/x.git", constants.CloneFromSchemeHTTP},
		{"ssh://git@host/x.git", constants.CloneFromSchemeSSH},
		{"git://host/x.git", constants.CloneFromSchemeGit},
		{"file:///tmp/x.git", constants.CloneFromSchemeFile},
		{"git@github.com:owner/repo.git", constants.CloneFromSchemeSCP},
		{"weird-thing-no-colon", constants.CloneFromSchemeOther},
	}
	for _, tc := range cases {
		if got := ClassifyScheme(tc.url); got != tc.want {
			t.Errorf("ClassifyScheme(%q) = %q, want %q", tc.url, got, tc.want)
		}
	}
}

// TestRenderSummaryTerminal_FullBlock asserts every section appears
// (banner, found count, by-mode header, per-scheme rows, status
// tally, both report paths) and that zero-count schemes are omitted.
func TestRenderSummaryTerminal_FullBlock(t *testing.T) {
	results := []Result{
		{Status: constants.CloneFromStatusOK, Row: Row{URL: "https://x/a.git"}},
		{Status: constants.CloneFromStatusOK, Row: Row{URL: "https://x/b.git"}},
		{Status: constants.CloneFromStatusSkipped, Row: Row{URL: "git@h:o/c.git"}},
		{Status: constants.CloneFromStatusFailed, Row: Row{URL: "ssh://h/d.git"}, Detail: "boom"},
	}
	var buf bytes.Buffer
	if err := RenderSummaryTerminal(&buf, results, "/r/x.csv", "/r/x.json"); err != nil {
		t.Fatalf("RenderSummaryTerminal: %v", err)
	}
	out := buf.String()
	want := []string{
		"gitmap clone-from: summary",
		"found:    4 repo(s)",
		"by mode:",
		"https   2",
		"ssh     1",
		"scp     1",
		"status:   2 ok, 1 skipped, 1 failed (4 total)",
		"report csv : /r/x.csv",
		"report json: /r/x.json",
	}
	for _, s := range want {
		if !strings.Contains(out, s) {
			t.Errorf("missing %q in:\n%s", s, out)
		}
	}
	// Zero-count schemes must NOT appear.
	for _, absent := range []string{"http ", "git ", "file ", "other "} {
		if strings.Contains(out, "    "+absent) {
			t.Errorf("zero-count scheme %q rendered in:\n%s", absent, out)
		}
	}
}

// TestRenderSummaryTerminal_NoReportPlaceholder confirms that when
// both paths are empty the renderer still emits the predictable
// placeholder line so log scrapers can pin the section's presence.
func TestRenderSummaryTerminal_NoReportPlaceholder(t *testing.T) {
	var buf bytes.Buffer
	_ = RenderSummaryTerminal(&buf, nil, "", "")
	if !strings.Contains(buf.String(), "(skipped") {
		t.Errorf("missing placeholder line:\n%s", buf.String())
	}
}
