package cmdagy

import (
	"strings"
	"testing"
)

func TestCleanWorkspaceURI(t *testing.T) {
	raw := `["file:///d%3A/work/cat-my"]`
	expected := "d:/work/cat-my"
	actual := cleanWorkspaceURI(raw)
	if actual != expected {
		t.Errorf("cleanWorkspaceURI(%q) = %q, expected %q", raw, actual, expected)
	}
}

func TestParseAgySummariesOptions(t *testing.T) {
	args := []string{"25", "-a", "-j", "-q", "fix"}
	opts := parseAgySummariesOptions(args)
	if opts.limit != 25 {
		t.Errorf("expected limit 25, got %d", opts.limit)
	}
	if !opts.isActive {
		t.Errorf("expected isActive true")
	}
	if !opts.isJSON {
		t.Errorf("expected isJSON true")
	}
	if opts.query != "fix" {
		t.Errorf("expected query 'fix', got %q", opts.query)
	}
}

func TestBuildAgySummariesQuery(t *testing.T) {
	opts := agySummariesOptions{limit: 5, isActive: true, query: "test"}
	sql, params := buildAgySummariesQuery(opts)
	if !strings.Contains(sql, "WHERE") {
		t.Errorf("expected WHERE in sql: %s", sql)
	}
	if !strings.Contains(sql, "LIMIT ?") {
		t.Errorf("expected LIMIT in sql: %s", sql)
	}
	if len(params) != 4 {
		t.Errorf("expected 4 params, got %d", len(params))
	}
}

func TestTruncateSummaryText(t *testing.T) {
	short := "Hello"
	if truncateSummaryText(short, 10) != "Hello" {
		t.Errorf("expected unchanged short text")
	}
	long := "This is a very long text that must be truncated"
	trunc := truncateSummaryText(long, 15)
	if len(trunc) > 15 || !strings.HasSuffix(trunc, "...") {
		t.Errorf("unexpected truncated text: %q", trunc)
	}
}

func TestIsAgySummariesSubcommand(t *testing.T) {
	valid := []string{"summaries", "summary", "convs", "conversations", "agydb", "conv-db"}
	for _, v := range valid {
		if !isAgySummariesSubcommand(v) {
			t.Errorf("expected %q to be recognized", v)
		}
	}
	if isAgySummariesSubcommand("unknown") {
		t.Errorf("expected 'unknown' to be false")
	}
}
