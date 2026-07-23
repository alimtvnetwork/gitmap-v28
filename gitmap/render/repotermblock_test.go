package render

import (
	"bytes"
	"strings"
	"testing"
)

// TestRenderRepoTermBlock_FullBlock verifies the canonical 5-line
// block layout when every field is populated.
func TestRenderRepoTermBlock_FullBlock(t *testing.T) {
	var buf bytes.Buffer
	err := RenderRepoTermBlock(&buf, RepoTermBlock{
		Index:        1,
		Name:         "scripts-fixer",
		Branch:       "main",
		BranchSource: "HEAD",
		OriginalURL:  "https://github.com/owner/scripts-fixer.git",
		TargetURL:    "https://github.com/owner/scripts-fixer.git",
		CloneCommand: "git clone https://github.com/owner/scripts-fixer.git scripts-fixer",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := strings.Join([]string{
		"  1. scripts-fixer",
		"     branch:    main (HEAD)",
		"     transport: https",
		"     https:     https://github.com/owner/scripts-fixer.git",
		"     ssh:       (unknown)",
		"     from:      https://github.com/owner/scripts-fixer.git",
		"     to:        https://github.com/owner/scripts-fixer.git",
		"     command:   git clone https://github.com/owner/scripts-fixer.git scripts-fixer",
		"",
	}, "\n")
	if got := buf.String(); got != want {
		t.Fatalf("block mismatch\n--- want ---\n%s\n--- got ---\n%s", want, got)
	}
}

// TestRenderRepoTermBlock_UnknownPlaceholders ensures empty fields
// render as "(unknown)" so the block shape stays stable.
func TestRenderRepoTermBlock_UnknownPlaceholders(t *testing.T) {
	var buf bytes.Buffer
	err := RenderRepoTermBlock(&buf, RepoTermBlock{Index: 2, Name: "x"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	for _, line := range []string{
		"     branch:    (unknown)",
		"     transport: other",
		"     https:     (unknown)",
		"     ssh:       (unknown)",
		"     from:      (unknown)",
		"     to:        (unknown)",
		"     command:   (unknown)",
	} {
		if !strings.Contains(got, line) {
			t.Fatalf("missing line %q in:\n%s", line, got)
		}
	}
}

// TestRenderRepoTermBlock_BranchWithoutSource omits the parenthesized
// source when only the branch name is known.
func TestRenderRepoTermBlock_BranchWithoutSource(t *testing.T) {
	var buf bytes.Buffer
	_ = RenderRepoTermBlock(&buf, RepoTermBlock{
		Index: 1, Name: "n", Branch: "develop",
	})
	if !strings.Contains(buf.String(), "     branch:    develop\n") {
		t.Fatalf("branch-only line missing:\n%s", buf.String())
	}
}

// TestRenderRepoTermBlocks_OrderPreserved verifies multi-block rendering
// emits blocks in the input order.
func TestRenderRepoTermBlocks_OrderPreserved(t *testing.T) {
	var buf bytes.Buffer
	err := RenderRepoTermBlocks(&buf, []RepoTermBlock{
		{Index: 1, Name: "a"},
		{Index: 2, Name: "b"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	ai := strings.Index(out, "1. a")
	bi := strings.Index(out, "2. b")
	if ai < 0 || bi < 0 || ai > bi {
		t.Fatalf("order mismatch:\n%s", out)
	}
}

// TestFormatBranch_EdgeCases exercises the branch helper directly so
// callers can rely on its contract without going through the writer.
func TestFormatBranch_EdgeCases(t *testing.T) {
	cases := []struct {
		branch, source, want string
	}{
		{"", "", fieldUnknown},
		{"  ", "HEAD", fieldUnknown},
		{"main", "", "main"},
		{"main", "HEAD", "main (HEAD)"},
		{"main", " HEAD ", "main (HEAD)"},
	}
	for _, c := range cases {
		if got := formatBranch(c.branch, c.source); got != c.want {
			t.Fatalf("formatBranch(%q,%q) = %q, want %q",
				c.branch, c.source, got, c.want)
		}
	}
}
