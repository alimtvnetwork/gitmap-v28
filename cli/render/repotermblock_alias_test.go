package render

import (
	"bytes"
	"strings"
	"testing"
)

func TestRenderRepoTermBlock_PrimaryAliasBracket(t *testing.T) {
	var buf bytes.Buffer
	block := RepoTermBlock{
		Index:        1,
		Name:         "my-package",
		PrimaryAlias: "mp",
	}
	if err := RenderRepoTermBlock(&buf, block); err != nil {
		t.Fatalf("RenderRepoTermBlock failed: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "  1. my-package [mp]\n") {
		t.Errorf("expected header to contain bracketed alias, got:\n%s", out)
	}
}

func TestRenderRepoTermBlock_MultiAliasTree(t *testing.T) {
	var buf bytes.Buffer
	block := RepoTermBlock{
		Index:        2,
		Name:         "my-package",
		PrimaryAlias: "mp",
		Aliases:      []string{"mp", "my-pkg"},
	}
	if err := RenderRepoTermBlock(&buf, block); err != nil {
		t.Fatalf("RenderRepoTermBlock failed: %v", err)
	}

	out := buf.String()
	assertMultiAliasOutput(t, out)
}

func assertMultiAliasOutput(t *testing.T, out string) {
	t.Helper()
	if !strings.Contains(out, "aliases:\n") {
		t.Errorf("expected aliases block in body, got:\n%s", out)
	}
	if !strings.Contains(out, "├── mp (primary)") {
		t.Errorf("expected primary branch line, got:\n%s", out)
	}
	if !strings.Contains(out, "└── my-pkg (secondary)") {
		t.Errorf("expected secondary branch line, got:\n%s", out)
	}
}

func TestRenderRepoTermBlock_NoAlias(t *testing.T) {
	var buf bytes.Buffer
	block := RepoTermBlock{Index: 3, Name: "clean-repo"}
	if err := RenderRepoTermBlock(&buf, block); err != nil {
		t.Fatalf("RenderRepoTermBlock failed: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "  3. clean-repo\n") {
		t.Errorf("expected clean header without brackets, got:\n%s", out)
	}
	if strings.Contains(out, "aliases:\n") {
		t.Errorf("expected no aliases section, got:\n%s", out)
	}
}
