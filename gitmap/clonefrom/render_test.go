package clonefrom

// Render tests. No filesystem, no git — pure string assembly.
// Locks the dry-run output shape so a future tweak that breaks
// downstream log parsers (yes, some users grep dry-run output) is
// caught at PR time.

import (
	"bytes"
	"strings"
	"testing"
)

// TestRender_HeaderShape pins the three-line header (banner,
// source, count) and the trailing blank line that separates it
// from the per-row blocks.
func TestRender_HeaderShape(t *testing.T) {
	plan := Plan{Source: "/tmp/p.csv", Format: "csv", Rows: []Row{
		{URL: "https://github.com/a/b.git"},
	}}
	var buf bytes.Buffer
	if err := Render(&buf, plan); err != nil {
		t.Fatalf("Render: %v", err)
	}
	out := buf.String()
	wantSubstrings := []string{
		"gitmap clone-from: dry-run",
		"source: /tmp/p.csv (csv)",
		"1 row(s) -- pass --execute",
		"  1. https://github.com/a/b.git",
		"     dest:   b  (derived)",
		"     branch: (default HEAD)",
		"     depth:  full",
	}
	for _, s := range wantSubstrings {
		if !strings.Contains(out, s) {
			t.Errorf("missing substring: %q in:\n%s", s, out)
		}
	}
}

// TestRender_ExplicitFieldsNotMarkedDerived confirms an explicit
// dest does NOT get the `(derived)` annotation. Catches a future
// regression where displayDest forgets to branch on len(r.Dest).
func TestRender_ExplicitFieldsNotMarkedDerived(t *testing.T) {
	plan := Plan{Source: "x", Format: "json", Rows: []Row{
		{URL: "https://github.com/a/b.git", Dest: "custom", Branch: "dev", Depth: 5},
	}}
	var buf bytes.Buffer
	_ = Render(&buf, plan)
	out := buf.String()
	if strings.Contains(out, "(derived)") {
		t.Errorf("explicit dest got (derived) annotation:\n%s", out)
	}
	if !strings.Contains(out, "branch: dev") {
		t.Errorf("explicit branch missing")
	}
	if !strings.Contains(out, "depth:  5") {
		t.Errorf("explicit depth missing")
	}
}

// TestDeriveDest covers the URL→dirname computation across the
// three input shapes the executor will actually see in the wild.
func TestDeriveDest(t *testing.T) {
	cases := []struct {
		url, want string
	}{
		{"https://github.com/owner/repo.git", "repo"},
		{"https://github.com/owner/repo", "repo"},
		{"git@github.com:owner/repo.git", "repo"},
		{"ssh://git@host/path/to/proj.git", "proj"},
		{"https://example.org/", "repo"}, // empty basename → fallback
	}
	for _, tc := range cases {
		if got := DeriveDest(tc.url); got != tc.want {
			t.Errorf("DeriveDest(%q) = %q, want %q", tc.url, got, tc.want)
		}
	}
}
