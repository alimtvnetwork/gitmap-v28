package cmd

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/templates"
)

// fixtureEntries returns a deterministic three-row corpus for filter
// tests. We don't read the embedded FS so the tests stay hermetic.
func fixtureEntries() []templates.Entry {
	return []templates.Entry{
		{Kind: "ignore", Lang: "go", Source: templates.SourceEmbed, Path: "assets/ignore/go.gitignore"},
		{Kind: "ignore", Lang: "node", Source: templates.SourceUser, Path: "/home/me/.gitmap/templates/ignore/node.gitignore"},
		{Kind: "attributes", Lang: "go", Source: templates.SourceEmbed, Path: "assets/attributes/go.gitattributes"},
		{Kind: "lfs", Lang: "common", Source: templates.SourceEmbed, Path: "assets/lfs/common.gitattributes"},
	}
}

// TestFilterTemplatesNoFilters is the load-bearing identity case: with
// both filters empty, the input slice MUST come back unchanged so the
// default `templates list` invocation keeps working.
func TestFilterTemplatesNoFilters(t *testing.T) {
	got := filterTemplates(fixtureEntries(), "", "")
	if len(got) != 4 {
		t.Fatalf("want 4 rows, got %d", len(got))
	}
}

// TestFilterTemplatesByKind narrows to a single kind and proves both
// the matching rows survive and the non-matching kinds get dropped.
func TestFilterTemplatesByKind(t *testing.T) {
	got := filterTemplates(fixtureEntries(), "ignore", "")
	if len(got) != 2 {
		t.Fatalf("want 2 ignore rows, got %d", len(got))
	}
	for _, e := range got {
		if e.Kind != "ignore" {
			t.Errorf("unexpected kind in result: %s", e.Kind)
		}
	}
}

// TestFilterTemplatesByLang narrows to a single language across kinds
// — the cross-kind matching is the whole reason --lang exists.
func TestFilterTemplatesByLang(t *testing.T) {
	got := filterTemplates(fixtureEntries(), "", "go")
	if len(got) != 2 {
		t.Fatalf("want 2 go rows (ignore + attributes), got %d", len(got))
	}
	kinds := map[string]bool{}
	for _, e := range got {
		kinds[e.Kind] = true
	}
	if !kinds["ignore"] || !kinds["attributes"] {
		t.Errorf("expected both ignore + attributes for go, got %v", kinds)
	}
}

// TestFilterTemplatesByKindAndLang proves the two filters AND together,
// not OR. Without this guard, --kind ignore --lang go could leak
// attributes/go rows.
func TestFilterTemplatesByKindAndLang(t *testing.T) {
	got := filterTemplates(fixtureEntries(), "ignore", "go")
	if len(got) != 1 {
		t.Fatalf("want exactly one row, got %d", len(got))
	}
	if got[0].Kind != "ignore" || got[0].Lang != "go" {
		t.Errorf("want ignore/go, got %s/%s", got[0].Kind, got[0].Lang)
	}
}

// TestFilterTemplatesEmptyResult covers the no-match path so the
// downstream "(no templates match the requested filter)" message has
// a pinned trigger condition.
func TestFilterTemplatesEmptyResult(t *testing.T) {
	got := filterTemplates(fixtureEntries(), "ignore", "rust")
	if len(got) != 0 {
		t.Fatalf("want empty result, got %d rows", len(got))
	}
}

// TestIsValidKindFilter pins the allow-list. Adding a new kind here
// forces a sibling change in templates/constants.go (kind* constants)
// — a deliberate co-evolution.
func TestIsValidKindFilter(t *testing.T) {
	cases := map[string]bool{
		"":           true,
		"ignore":     true,
		"attributes": true,
		"lfs":        true,
		"foo":        false,
		"Ignore":     false, // already lowered upstream
	}
	for kind, want := range cases {
		if got := isValidKindFilter(kind); got != want {
			t.Errorf("isValidKindFilter(%q) = %v, want %v", kind, got, want)
		}
	}
}

// TestParseTemplatesListFlagsLowersValues guards the case-folding rule:
// users typing `--kind IGNORE --lang GO` should match the same rows as
// the lowercase forms. Trim is also asserted via the leading space.
func TestParseTemplatesListFlagsLowersValues(t *testing.T) {
	kind, lang := parseTemplatesListFlags([]string{"--kind", " IGNORE ", "--lang", "GO"})
	if kind != "ignore" {
		t.Errorf("kind = %q, want ignore", kind)
	}
	if lang != "go" {
		t.Errorf("lang = %q, want go", lang)
	}
}
