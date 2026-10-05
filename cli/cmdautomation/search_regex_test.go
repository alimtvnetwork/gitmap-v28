package cmdautomation

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNormalizeRegexPattern(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "candidate response with escaped pipe",
			input:    `Candidate Response\|CANDIDATE RESPONSE`,
			expected: `Candidate Response|CANDIDATE RESPONSE`,
		},
		{
			name:     "timer multiple escaped pipes",
			input:    `timer\|Clock\|8:47\|timerBox\|border.*timer`,
			expected: `timer|Clock|8:47|timerBox|border.*timer`,
		},
		{
			name:     "slide layout escaped pipes",
			input:    `slide.*layout\|slideMode\|isSlide\|slide-mode\|effectiveLayout`,
			expected: `slide.*layout|slideMode|isSlide|slide-mode|effectiveLayout`,
		},
		{
			name:     "layout mode escaped pipes",
			input:    `layoutMode\|effectiveLayout\|answer.*placement\|lg:grid-cols`,
			expected: `layoutMode|effectiveLayout|answer.*placement|lg:grid-cols`,
		},
		{
			name:     "escaped backslash before pipe preserved",
			input:    `foo\\|bar`,
			expected: `foo\\|bar`,
		},
		{
			name:     "plain string without pipe",
			input:    `Candidate Response`,
			expected: `Candidate Response`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := normalizeRegexPattern(tc.input)
			if actual != tc.expected {
				t.Errorf("normalizeRegexPattern(%q) = %q, want %q", tc.input, actual, tc.expected)
			}
		})
	}
}

func TestAutoPromoteRegex(t *testing.T) {
	opts1 := SearchOptions{
		Pattern: `Candidate Response\|CANDIDATE RESPONSE`,
		IsRegex: false,
	}
	autoPromoteRegex(&opts1)
	if !opts1.IsRegex {
		t.Errorf("expected IsRegex to be true after auto-promotion")
	}
	if opts1.Pattern != "Candidate Response|CANDIDATE RESPONSE" {
		t.Errorf("expected pattern %q, got %q", "Candidate Response|CANDIDATE RESPONSE", opts1.Pattern)
	}

	opts2 := SearchOptions{
		Pattern: `timer\|Clock\|8:47\|timerBox\|border.*timer`,
		IsRegex: true,
	}
	autoPromoteRegex(&opts2)
	if opts2.Pattern != "timer|Clock|8:47|timerBox|border.*timer" {
		t.Errorf("expected normalized pattern, got %q", opts2.Pattern)
	}

	opts3 := SearchOptions{
		Pattern: `Plain Word`,
		IsRegex: false,
	}
	autoPromoteRegex(&opts3)
	if opts3.IsRegex {
		t.Errorf("expected IsRegex to remain false for plain string")
	}
}

func TestResolveSearchCandidateFiles(t *testing.T) {
	tempDir := t.TempDir()
	subDir := filepath.Join(tempDir, "components", "runner")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	formRunnerPath := filepath.Join(subDir, "FormRunner.tsx")
	if err := os.WriteFile(formRunnerPath, []byte("export const FormRunner = () => null;"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	// Case 1: Exact file path
	files1, found1 := resolveSearchCandidateFiles(formRunnerPath)
	if !found1 || len(files1) != 1 || files1[0] != formRunnerPath {
		t.Errorf("case 1 exact path failed: %v, %v", files1, found1)
	}

	// Case 2: Missing extension (FormRunner)
	noExtPath := filepath.Join(subDir, "FormRunner")
	files2, found2 := resolveSearchCandidateFiles(noExtPath)
	if !found2 || len(files2) != 1 || files2[0] != formRunnerPath {
		t.Errorf("case 2 missing ext failed: %v, %v", files2, found2)
	}

	// Case 3: Truncated extension (FormRunner.t)
	truncPath := filepath.Join(subDir, "FormRunner.t")
	files3, found3 := resolveSearchCandidateFiles(truncPath)
	if !found3 || len(files3) != 1 || files3[0] != formRunnerPath {
		t.Errorf("case 3 truncated ext failed: %v, %v", files3, found3)
	}
}

func TestRunSearchRegexScenarios(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "runner.tsx")
	content := []byte("const CANDIDATE_RESPONSE = 'ok';\nconst timer = 10;\nconst isSlide = true;\nconst layoutMode = 'grid';\n")
	if err := os.WriteFile(testFile, content, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	scenarios := []struct {
		name    string
		pattern string
		isRegex bool
		isCI    bool
	}{
		{
			name:    "candidate response escaped pipe with auto-regex",
			pattern: `Candidate Response\|CANDIDATE_RESPONSE`,
			isRegex: false,
			isCI:    true,
		},
		{
			name:    "timer escaped pipe with explicit regex",
			pattern: `timer\|Clock\|8:47\|timerBox\|border.*timer`,
			isRegex: true,
			isCI:    false,
		},
		{
			name:    "slide layout escaped pipe without regex flag",
			pattern: `slide.*layout\|slideMode\|isSlide\|slide-mode\|effectiveLayout`,
			isRegex: false,
			isCI:    false,
		},
		{
			name:    "layout mode escaped pipe without regex flag",
			pattern: `layoutMode\|effectiveLayout\|answer.*placement\|lg:grid-cols`,
			isRegex: false,
			isCI:    false,
		},
	}

	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			opts := SearchOptions{
				Pattern:           sc.pattern,
				Dir:               testFile,
				IsRegex:           sc.isRegex,
				IsCaseInsensitive: sc.isCI,
			}
			res, appErr := RunSearch(opts)
			if appErr != nil {
				t.Fatalf("RunSearch failed: %v", appErr)
			}
			if len(res.Matches) == 0 {
				t.Errorf("scenario %s expected match, got 0 matches", sc.name)
			}
		})
	}
}
