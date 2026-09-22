package cmd

import "testing"

func TestIsPatternMatch_Wildcard(t *testing.T) {
	if !isPatternMatch("foo/bar.md", "BAR.md", "*") {
		t.Errorf("expected * to match BAR.md")
	}
	if !isPatternMatch("foo/bar.txt", "BAR.txt", "*.*") {
		t.Errorf("expected *.* to match BAR.txt")
	}
}

func TestIsPatternMatch_Extension(t *testing.T) {
	if !isPatternMatch("docs/README.md", "README.md", "*.md") {
		t.Errorf("expected *.md to match README.md")
	}
	if !isPatternMatch("docs/README.md", "README.md", "*md") {
		t.Errorf("expected *md to match README.md")
	}
	if isPatternMatch("docs/README.txt", "README.txt", "*.md") {
		t.Errorf("expected *.md not to match README.txt")
	}
}

func TestIsPatternMatch_Prefix(t *testing.T) {
	if !isPatternMatch("skills/SKILL.md", "SKILL.md", "SKILL*") {
		t.Errorf("expected SKILL* to match SKILL.md")
	}
	if !isPatternMatch("README.md", "README.md", "readme*") {
		t.Errorf("expected readme* to match README.md")
	}
}

func TestIsPatternMatch_RelativePath(t *testing.T) {
	if !isPatternMatch("docs/API.md", "API.md", "docs/*.md") {
		t.Errorf("expected docs/*.md to match docs/API.md")
	}
	if !isPatternMatch("docs/sub/API.md", "API.md", "docs") {
		t.Errorf("expected docs directory prefix to match docs/sub/API.md")
	}
	if !isPatternMatch("docs/sub/API.md", "API.md", ".") {
		t.Errorf("expected . to match docs/sub/API.md")
	}
	if !isPatternMatch("docs/sub/API.md", "API.md", "**/*.md") {
		t.Errorf("expected **/*.md to match docs/sub/API.md")
	}
}

func TestIsPatternMatch_Substring(t *testing.T) {
	if !isPatternMatch("docs/MY_DOC_FILE.txt", "MY_DOC_FILE.txt", "*doc*") {
		t.Errorf("expected *doc* to match MY_DOC_FILE.txt")
	}
	if isPatternMatch("docs/OTHER.txt", "OTHER.txt", "*doc*") {
		t.Errorf("expected *doc* not to match OTHER.txt")
	}
}

func TestFindMatchingPattern(t *testing.T) {
	matched, p := findMatchingPattern("docs/README.md", "README.md", []string{"*.txt", "*.md"})
	if !matched || p != "*.md" {
		t.Errorf("expected matched=true with pattern=*.md, got %v, %s", matched, p)
	}
}
