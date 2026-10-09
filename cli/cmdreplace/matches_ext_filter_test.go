package cmdreplace

import (
	"testing"
)

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

func TestMatchesExtFilter_EmptySlice(t *testing.T) {
	if !matchesExtFilter("anything.xyz", []string{}, true) {
		t.Error("empty slice should match all files")
	}
}
