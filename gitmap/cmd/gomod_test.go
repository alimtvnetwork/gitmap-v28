package cmd

import (
	"testing"
)

func TestParseExtFlag_Empty(t *testing.T) {
	exts := parseExtFlag("")
	if exts != nil {
		t.Errorf("expected nil, got %v", exts)
	}
}

func TestParseExtFlag_Single(t *testing.T) {
	exts := parseExtFlag("*.go")
	if len(exts) != 1 || exts[0] != ".go" {
		t.Errorf("expected [.go], got %v", exts)
	}
}

func TestParseExtFlag_Multiple(t *testing.T) {
	exts := parseExtFlag("*.go,*.md,*.txt")
	if len(exts) != 3 {
		t.Fatalf("expected 3 exts, got %d", len(exts))
	}
	expected := []string{".go", ".md", ".txt"}
	for i, e := range expected {
		if exts[i] != e {
			t.Errorf("exts[%d]: expected %q, got %q", i, e, exts[i])
		}
	}
}

func TestParseExtFlag_WithSpaces(t *testing.T) {
	exts := parseExtFlag("*.go , *.md")
	if len(exts) != 2 || exts[0] != ".go" || exts[1] != ".md" {
		t.Errorf("expected [.go .md], got %v", exts)
	}
}

func TestParseExtFlag_NoStar(t *testing.T) {
	exts := parseExtFlag(".go,.md")
	if len(exts) != 2 || exts[0] != ".go" || exts[1] != ".md" {
		t.Errorf("expected [.go .md], got %v", exts)
	}
}

func TestMatchesExtFilter_EmptyAllows(t *testing.T) {
	if !matchesExtFilter("main.go", nil, true) {
		t.Error("empty exts should match all files")
	}
}

func TestMatchesExtFilter_Match(t *testing.T) {
	if !matchesExtFilter("main.go", []string{".go", ".md"}, true) {
		t.Error("expected .go to match")
	}
}

func TestMatchesExtFilter_NoMatch(t *testing.T) {
	if matchesExtFilter("style.css", []string{".go", ".md"}, true) {
		t.Error("expected .css to not match")
	}
}

func TestMatchesExtFilter_NoExtension(t *testing.T) {
	if matchesExtFilter("Makefile", []string{".go"}, true) {
		t.Error("expected file without extension to not match")
	}
}

func TestParseGoModFlags_Defaults(t *testing.T) {
	opts := parseGoModFlags([]string{"github.com/new/path"})
	if opts.newPath != "github.com/new/path" {
		t.Errorf("expected newPath=github.com/new/path, got %q", opts.newPath)
	}
	if opts.dryRun || opts.noMerge || opts.noTidy || opts.verbose {
		t.Error("expected all flags false by default")
	}
	if opts.exts != nil {
		t.Errorf("expected nil exts, got %v", opts.exts)
	}
}

func TestParseGoModFlags_AllFlags(t *testing.T) {
	opts := parseGoModFlags([]string{
		"--dry-run", "--no-merge", "--no-tidy", "--verbose",
		"--ext", "*.go,*.md",
		"github.com/new/path",
	})
	if !opts.dryRun || !opts.noMerge || !opts.noTidy || !opts.verbose {
		t.Error("expected all flags true")
	}
	if len(opts.exts) != 2 {
		t.Errorf("expected 2 exts, got %v", opts.exts)
	}
	if opts.newPath != "github.com/new/path" {
		t.Errorf("expected newPath=github.com/new/path, got %q", opts.newPath)
	}
}

func TestParseGoModFlags_NoArgs(t *testing.T) {
	opts := parseGoModFlags([]string{})
	if opts.newPath != "" {
		t.Errorf("expected empty newPath, got %q", opts.newPath)
	}
}
