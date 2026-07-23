package cmd

import (
	"testing"
)

func TestParseModuleLine_Valid(t *testing.T) {
	content := "module github.com/user/repo\n\ngo 1.21\n"
	got := parseModuleLine(content)
	if got != "github.com/user/repo" {
		t.Errorf("expected github.com/user/repo, got %q", got)
	}
}

func TestParseModuleLine_WithSpaces(t *testing.T) {
	content := "  module   github.com/org/pkg  \n"
	got := parseModuleLine(content)
	if got != "github.com/org/pkg" {
		t.Errorf("expected github.com/org/pkg, got %q", got)
	}
}

func TestParseModuleLine_MultipleLines(t *testing.T) {
	content := "// comment\nmodule github.com/test/app\n\nrequire (\n)\n"
	got := parseModuleLine(content)
	if got != "github.com/test/app" {
		t.Errorf("expected github.com/test/app, got %q", got)
	}
}

func TestMatchesExtFilter_EmptySlice(t *testing.T) {
	if !matchesExtFilter("anything.xyz", []string{}, true) {
		t.Error("empty slice should match all files")
	}
}

func TestIsExcludedDir_Git(t *testing.T) {
	if !isExcludedDir(".git") {
		t.Error("expected .git to be excluded")
	}
}

func TestIsExcludedDir_Vendor(t *testing.T) {
	if !isExcludedDir("vendor") {
		t.Error("expected vendor to be excluded")
	}
}

func TestIsExcludedDir_NodeModules(t *testing.T) {
	if !isExcludedDir("node_modules") {
		t.Error("expected node_modules to be excluded")
	}
}

func TestIsExcludedDir_Regular(t *testing.T) {
	if isExcludedDir("src") {
		t.Error("expected src to not be excluded")
	}
}

func TestFileContains_NotFound(t *testing.T) {
	if fileContains("nonexistent_file_xyz.go", "anything") {
		t.Error("expected false for nonexistent file")
	}
}
