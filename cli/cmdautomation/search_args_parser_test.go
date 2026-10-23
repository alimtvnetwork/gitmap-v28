package cmdautomation

import (
	"testing"
)

func TestParseSearchPositionalArgs(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		wantPattern string
		wantDir     string
	}{
		{
			name:        "single arg pattern",
			args:        []string{"hello"},
			wantPattern: "hello",
			wantDir:     "",
		},
		{
			name:        "powershell split escaped all",
			args:        []string{`\`, `all\`},
			wantPattern: "all",
			wantDir:     "",
		},
		{
			name:        "powershell split escaped all with dir",
			args:        []string{`\`, `all\`, "cli"},
			wantPattern: "all",
			wantDir:     "cli",
		},
		{
			name:        "powershell split multi-word pattern with target file",
			args:        []string{`\`, "Candidate", `Response\`, "src/components/runner/FormRunner.tsx"},
			wantPattern: "Candidate Response",
			wantDir:     "src/components/runner/FormRunner.tsx",
		},
		{
			name:        "multi-word unquoted pattern with target file",
			args:        []string{"Candidate", "Response", "search.go"},
			wantPattern: "Candidate Response",
			wantDir:     "search.go",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotPattern, gotDir := parseSearchPositionalArgs(tc.args)
			if gotPattern != tc.wantPattern {
				t.Errorf("parseSearchPositionalArgs() gotPattern = %q, want %q", gotPattern, tc.wantPattern)
			}
			if gotDir != tc.wantDir {
				t.Errorf("parseSearchPositionalArgs() gotDir = %q, want %q", gotDir, tc.wantDir)
			}
		})
	}
}

func TestResolveSearchArgs(t *testing.T) {
	defer resetSearchFlags()

	t.Run("resolves explicit path flag with quotes", func(t *testing.T) {
		resetSearchFlags()
		searchPathFlag = `"02-spec"`
		pat, dir := resolveSearchArgs([]string{`"letterly"`})
		if pat != "letterly" {
			t.Errorf("expected pattern 'letterly', got %q", pat)
		}
		if dir != "02-spec" {
			t.Errorf("expected dir '02-spec', got %q", dir)
		}
	})

	t.Run("resolves explicit dir flag", func(t *testing.T) {
		resetSearchFlags()
		searchDirFlag = "cli"
		pat, dir := resolveSearchArgs([]string{"hello", "world"})
		if pat != "hello world" {
			t.Errorf("expected pattern 'hello world', got %q", pat)
		}
		if dir != "cli" {
			t.Errorf("expected dir 'cli', got %q", dir)
		}
	})

	t.Run("resolves explicit root and query flags without positional args", func(t *testing.T) {
		resetSearchFlags()
		searchRootFlag = `"cli/cmd"`
		searchQueryFlag = `"cmdpipeline"`
		pat, dir := resolveSearchArgs(nil)
		if pat != "cmdpipeline" {
			t.Errorf("expected pattern 'cmdpipeline', got %q", pat)
		}
		if dir != "cli/cmd" {
			t.Errorf("expected dir 'cli/cmd', got %q", dir)
		}
	})

	t.Run("resolves explicit query flag with positional dir", func(t *testing.T) {
		resetSearchFlags()
		searchQueryFlag = "cmdpipeline"
		pat, dir := resolveSearchArgs([]string{"cli/cmd"})
		if pat != "cmdpipeline" {
			t.Errorf("expected pattern 'cmdpipeline', got %q", pat)
		}
		if dir != "cli/cmd" {
			t.Errorf("expected dir 'cli/cmd', got %q", dir)
		}
	})

	t.Run("resolves explicit root flag with positional pattern", func(t *testing.T) {
		resetSearchFlags()
		searchRootFlag = "cli/cmd"
		pat, dir := resolveSearchArgs([]string{"cmdpipeline"})
		if pat != "cmdpipeline" {
			t.Errorf("expected pattern 'cmdpipeline', got %q", pat)
		}
		if dir != "cli/cmd" {
			t.Errorf("expected dir 'cli/cmd', got %q", dir)
		}
	})

	t.Run("resets flags properly", func(t *testing.T) {
		searchPathFlag = "some-path"
		searchDirFlag = "some-dir"
		searchRootFlag = "some-root"
		searchQueryFlag = "some-query"
		searchPatternFlag = "some-pattern"
		resetSearchFlags()
		if searchPathFlag != "" || searchDirFlag != "" || searchRootFlag != "" || searchQueryFlag != "" || searchPatternFlag != "" {
			t.Errorf("expected empty flags after reset, got root=%q query=%q", searchRootFlag, searchQueryFlag)
		}
	})
}
