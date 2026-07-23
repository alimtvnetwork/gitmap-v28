package cmd

import (
	"strings"
	"testing"
)

// TestNormalizeLangs proves the user's order is preserved while
// duplicates and the implicit `common` are stripped. Without this
// guarantee, `add ignore go go common node` would either crash
// resolveAddTemplates with a duplicate-resolve race or silently
// double-prepend `common` and bloat the marker block.
func TestNormalizeLangs(t *testing.T) {
	cases := []struct {
		name string
		in   []string
		want []string
	}{
		{"empty", nil, []string{}},
		{"single", []string{"go"}, []string{"go"}},
		{"strip common", []string{"common", "go"}, []string{"go"}},
		{"dedupe", []string{"go", "GO", " go "}, []string{"go"}},
		{"preserve order", []string{"node", "go", "rust"}, []string{"node", "go", "rust"}},
		{"mixed case + blanks", []string{"", "Node", " GO ", "rust"}, []string{"node", "go", "rust"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := normalizeLangs(tc.in)
			if len(got) != len(tc.want) {
				t.Fatalf("len: want %d, got %d (%v)", len(tc.want), len(got), got)
			}
			for i, v := range got {
				if v != tc.want[i] {
					t.Errorf("[%d]: want %q, got %q", i, tc.want[i], v)
				}
			}
		})
	}
}

// TestBuildAddTagSorted proves "go,node" and "node,go" share a tag.
// This is the load-bearing invariant for marker-block stability —
// without it, `add ignore go node` followed by `add ignore node go`
// would write two parallel blocks instead of refreshing one.
func TestBuildAddTagSorted(t *testing.T) {
	a := buildAddTag("ignore", []string{"go", "node"})
	b := buildAddTag("ignore", []string{"node", "go"})
	if a != b {
		t.Fatalf("tag order mismatch: %q vs %q", a, b)
	}
	if a != "ignore/go+node" {
		t.Errorf("want ignore/go+node, got %q", a)
	}
	if buildAddTag("attributes", nil) != "attributes/common" {
		t.Errorf("empty langs should yield <kind>/common")
	}
}

// TestDedupeLinesPreservesBlanks pins the rule that blank lines are
// visual separators (kept) while non-blank duplicates collapse to
// the first occurrence. Without this, merging common+go would lose
// every section spacer and become an unreadable wall of rules.
func TestDedupeLinesPreservesBlanks(t *testing.T) {
	in := []byte("# section A\n*.exe\n\n# section A\n*.exe\n*.dll\n\n*.exe\n")
	got := string(dedupeLines(in))
	want := "# section A\n*.exe\n\n\n*.dll\n\n\n"
	if got != want {
		t.Fatalf("dedupe mismatch:\nwant: %q\ngot:  %q", want, got)
	}
}

// TestConcatTemplateBodiesAddsLangBanners proves each resolved template
// is prefixed with a `# ── <lang> ──` banner so multi-lang merges stay
// scannable in the rendered .gitignore.
func TestConcatTemplateBodiesAddsLangBanners(t *testing.T) {
	body := concatTemplateBodies(nil)
	if len(body) != 0 {
		t.Fatalf("nil resolved should yield empty body, got %q", body)
	}
	// Smoke test on real embedded templates — common always exists.
	got := concatTemplateBodies(mustResolveForTest(t, "ignore", "common"))
	if !strings.Contains(string(got), "# ── common ──") {
		t.Errorf("expected lang banner in merged body, got:\n%s", got)
	}
}
