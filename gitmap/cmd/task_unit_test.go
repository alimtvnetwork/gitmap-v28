package cmd

import "testing"

func TestIsGitignoreCommentClassifies(t *testing.T) {
	cases := []struct {
		line string
		want bool
	}{
		{"", true},
		{"# comment", true},
		{"#trailing", true},
		{"node_modules/", false},
		{"*.log", false},
		{"src/foo.go", false},
	}
	for _, c := range cases {
		got := isGitignoreComment(c.line)
		if got != c.want {
			t.Errorf("isGitignoreComment(%q) = %v, want %v", c.line, got, c.want)
		}
	}
}

func TestParseGitignoreLinesStripsCommentsAndBlanks(t *testing.T) {
	content := "# header\n\nnode_modules/\n  *.log  \n# trailing\nbuild/\n"
	got := parseGitignoreLines(content)
	want := []string{"node_modules/", "*.log", "build/"}
	if len(got) != len(want) {
		t.Fatalf("len = %d, want %d (%v)", len(got), len(want), got)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestParseGitignoreLinesEmptyInput(t *testing.T) {
	if got := parseGitignoreLines(""); len(got) != 0 {
		t.Errorf("expected empty, got %v", got)
	}
}

func TestMatchGlobBasename(t *testing.T) {
	cases := []struct {
		path, pattern string
		want          bool
	}{
		{"src/app.log", "*.log", true},
		{"src/app.go", "*.log", false},
		{"node_modules", "node_modules", true},
		{"src/main.go", "main.go", true},
		{"src/main.go", "[invalid", false}, // bad pattern -> false
	}
	for _, c := range cases {
		got := matchGitignoreGlob(c.path, c.pattern)
		if got != c.want {
			t.Errorf("matchGitignoreGlob(%q,%q) = %v, want %v", c.path, c.pattern, got, c.want)
		}
	}
}

func TestMatchesPatternDirOnlyVsFile(t *testing.T) {
	cases := []struct {
		name    string
		relPath string
		isDir   bool
		pattern string
		want    bool
	}{
		{"dir-pattern matches dir", "node_modules", true, "node_modules/", true},
		{"dir-pattern skips file", "node_modules", false, "node_modules/", false},
		{"file-pattern matches file", "app.log", false, "*.log", true},
		{"file-pattern matches dir basename", "logs", true, "*.log", false},
	}
	for _, c := range cases {
		got := matchesPattern(c.relPath, c.isDir, c.pattern)
		if got != c.want {
			t.Errorf("%s: matchesPattern(%q, isDir=%v, %q) = %v, want %v",
				c.name, c.relPath, c.isDir, c.pattern, got, c.want)
		}
	}
}

func TestIsIgnoredEmptyPatternsShortCircuit(t *testing.T) {
	if isIgnored("anything", false, nil) {
		t.Error("nil patterns should not match")
	}
	if isIgnored("anything", true, []string{}) {
		t.Error("empty patterns should not match")
	}
}

func TestIsIgnoredAnyMatchWins(t *testing.T) {
	patterns := []string{"*.tmp", "*.log", "node_modules/"}
	if !isIgnored("debug.log", false, patterns) {
		t.Error("expected debug.log to be ignored")
	}
	if !isIgnored("node_modules", true, patterns) {
		t.Error("expected node_modules dir to be ignored")
	}
	if isIgnored("src/main.go", false, patterns) {
		t.Error("src/main.go should not be ignored")
	}
}
