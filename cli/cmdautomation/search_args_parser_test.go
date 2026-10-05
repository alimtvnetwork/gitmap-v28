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
