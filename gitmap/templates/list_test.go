package templates

import (
	"os"
	"path/filepath"
	"testing"
)

func TestListIncludesEmbeddedCorpus(t *testing.T) {
	withTempHome(t) // ensures no overlay leaks in
	entries, err := List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) == 0 {
		t.Fatal("expected at least the embedded corpus")
	}
	hasGoIgnore := false
	for _, e := range entries {
		if e.Kind == kindIgnore && e.Lang == "go" {
			hasGoIgnore = true
			if e.Source != SourceEmbed {
				t.Errorf("ignore/go should be SourceEmbed, got %v", e.Source)
			}
		}
	}
	if !hasGoIgnore {
		t.Fatal("ignore/go missing from List()")
	}
}

func TestListSortedByKindThenLang(t *testing.T) {
	withTempHome(t)
	entries, err := List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	for i := 1; i < len(entries); i++ {
		prev, cur := entries[i-1], entries[i]
		pk, ck := kindRank(prev.Kind), kindRank(cur.Kind)
		if pk > ck {
			t.Fatalf("kind order broken at %d: %s before %s", i, prev.Kind, cur.Kind)
		}
		if pk == ck && prev.Lang > cur.Lang {
			t.Fatalf("lang order broken at %d: %s before %s", i, prev.Lang, cur.Lang)
		}
	}
}

func TestListUserOverlayShadowsEmbed(t *testing.T) {
	withTempHome(t)
	dir, err := EnsureUserDir()
	if err != nil {
		t.Fatalf("EnsureUserDir: %v", err)
	}
	overridePath := filepath.Join(dir, kindIgnore, "go"+templateExtIgnore)
	if mkErr := os.MkdirAll(filepath.Dir(overridePath), 0o755); mkErr != nil {
		t.Fatalf("mkdir: %v", mkErr)
	}
	if wErr := os.WriteFile(overridePath, []byte("# overridden\n"), 0o644); wErr != nil {
		t.Fatalf("write override: %v", wErr)
	}

	entries, err := List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	for _, e := range entries {
		if e.Kind == kindIgnore && e.Lang == "go" {
			if e.Source != SourceUser {
				t.Fatalf("ignore/go should be SourceUser after overlay, got %v", e.Source)
			}
			if e.Path != overridePath {
				t.Fatalf("ignore/go path = %s, want %s", e.Path, overridePath)
			}

			return
		}
	}
	t.Fatal("ignore/go missing after overlay")
}

func TestParseRelTemplatePath(t *testing.T) {
	cases := []struct {
		in       string
		wantKind string
		wantLang string
		wantOK   bool
	}{
		{"ignore/go.gitignore", "ignore", "go", true},
		{"attributes/node.gitattributes", "attributes", "node", true},
		{"lfs/common.gitattributes", "lfs", "common", true},
		{"ignore/go.gitattributes", "", "", false}, // wrong ext for kind
		{"unknown/foo.gitignore", "", "", false},
		{"toplevel.gitignore", "", "", false},
	}
	for _, c := range cases {
		k, l, ok := parseRelTemplatePath(c.in)
		if ok != c.wantOK || k != c.wantKind || l != c.wantLang {
			t.Errorf("parseRelTemplatePath(%q) = (%q,%q,%v), want (%q,%q,%v)",
				c.in, k, l, ok, c.wantKind, c.wantLang, c.wantOK)
		}
	}
}
