package cmd

import (
	"testing"
)

func TestIsFixAgyRequest_Variations(t *testing.T) {
	testCases := []struct {
		name     string
		args     []string
		expected bool
	}{
		{
			name:     "fix_agy",
			args:     []string{"agy"},
			expected: true,
		},
		{
			name:     "fix_errors_agy",
			args:     []string{"errors", "agy"},
			expected: true,
		},
		{
			name:     "fix_agy_errors",
			args:     []string{"agy", "errors"},
			expected: true,
		},
		{
			name:     "fix_aef",
			args:     []string{"aef"},
			expected: true,
		},
		{
			name:     "fix_agy_errors_fix",
			args:     []string{"agy-errors-fix"},
			expected: true,
		},
		{
			name:     "fix_normal_repo",
			args:     []string{"my-repo"},
			expected: false,
		},
		{
			name:     "empty_args",
			args:     []string{},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := isFixAgyRequest(tc.args)
			if actual != tc.expected {
				t.Fatalf("expected %v for %v, got %v", tc.expected, tc.args, actual)
			}
		})
	}
}

func TestIsFixAllRequested_Variations(t *testing.T) {
	testCases := []struct {
		name     string
		args     []string
		expected bool
	}{
		{"bare_all", []string{"all"}, true},
		{"flag_all", []string{"--all"}, true},
		{"short_a", []string{"-a"}, true},
		{"single_dash_all", []string{"-all"}, true},
		{"all_with_action", []string{"all", "wip"}, true},
		{"specific_repo", []string{"my-repo"}, false},
		{"empty", []string{}, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := isFixAllRequested(tc.args)
			if actual != tc.expected {
				t.Fatalf("expected %v for %v, got %v", tc.expected, tc.args, actual)
			}
		})
	}
}

func TestIsFixPromptRequested_Variations(t *testing.T) {
	testCases := []struct {
		name     string
		args     []string
		expected bool
	}{
		{"double_dash_prompt", []string{"--prompt"}, true},
		{"short_p", []string{"-p"}, true},
		{"word_prompt", []string{"prompt"}, true},
		{"double_dash_interactive", []string{"--interactive"}, true},
		{"short_i", []string{"-i"}, true},
		{"word_interactive", []string{"interactive"}, true},
		{"all_with_prompt", []string{"all", "--prompt"}, true},
		{"bare_all", []string{"all"}, false},
		{"empty", []string{}, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := isFixPromptRequested(tc.args)
			if actual != tc.expected {
				t.Fatalf("expected %v for %v, got %v", tc.expected, tc.args, actual)
			}
		})
	}
}

func TestResolveFixAllAction_Variations(t *testing.T) {
	testCases := []struct {
		name          string
		args          []string
		aliasOverride string
		expected      string
	}{
		{"default_stash", []string{"all"}, "", "stash"},
		{"explicit_wip", []string{"all", "wip"}, "", "wip"},
		{"explicit_2", []string{"all", "2"}, "", "2"},
		{"explicit_discard", []string{"all", "discard"}, "", "discard"},
		{"alias_override_takes_precedence", []string{"all"}, "wip", "wip"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := resolveFixAllAction(tc.args, tc.aliasOverride)
			if actual != tc.expected {
				t.Fatalf("expected %q, got %q", tc.expected, actual)
			}
		})
	}
}
