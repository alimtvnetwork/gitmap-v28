package termtable

import (
	"strings"
	"testing"
)

func TestRenderTable_Empty(t *testing.T) {
	rendered := RenderTable(TableConfig{})
	if rendered != "" {
		t.Fatalf("expected empty string, got %q", rendered)
	}
}

func TestRenderTable_Basic(t *testing.T) {
	cfg := TableConfig{
		Columns: []Column{
			{Title: "NODE", MinWidth: 10},
			{Title: "STATUS", MinWidth: 8},
		},
		Rows: []Row{
			{Cells: []string{"node-1", "ONLINE"}},
			{Cells: []string{"node-2", "OFFLINE"}},
		},
	}
	rendered := RenderTable(cfg)
	if !strings.Contains(rendered, "NODE") || !strings.Contains(rendered, "STATUS") {
		t.Fatalf("missing expected headers in rendered table:\n%s", rendered)
	}
	if !strings.Contains(rendered, "node-1") || !strings.Contains(rendered, "ONLINE") {
		t.Fatalf("missing expected row contents in rendered table:\n%s", rendered)
	}
}

func TestTruncateMiddle_ShortAndExact(t *testing.T) {
	if TruncateMiddle("hello", 10, "...") != "hello" {
		t.Fatalf("expected hello")
	}
	if TruncateMiddle("hello", 5, "...") != "hello" {
		t.Fatalf("expected hello")
	}
}

func TestTruncateMiddle_Long(t *testing.T) {
	truncated := TruncateMiddle("abcdefghijklmnop", 9, "...")
	if len(truncated) != 9 {
		t.Fatalf("expected length 9, got %d (%q)", len(truncated), truncated)
	}
	if !strings.Contains(truncated, "...") {
		t.Fatalf("expected ellipsis in %q", truncated)
	}
}

func TestTruncateMiddle_EdgeCases(t *testing.T) {
	if TruncateMiddle("abc", 0, "...") != "" {
		t.Fatalf("expected empty string for maxWidth 0")
	}
	if TruncateMiddle("abcdef", 2, "...") != ".." {
		t.Fatalf("expected .. for maxWidth 2")
	}
}
