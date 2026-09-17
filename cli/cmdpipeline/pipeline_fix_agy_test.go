package cmdpipeline

import (
	"testing"
)

func TestIsPipelineFixAgyArgs_UserRequestVariations(t *testing.T) {
	testCases := []struct {
		name     string
		args     []string
		expected bool
	}{
		{
			name:     "pipeline_fix_errors_agy",
			args:     []string{"fix", "errors", "agy"},
			expected: true,
		},
		{
			name:     "pipeline_fix_errors_agy_with_flags",
			args:     []string{"fix", "errors", "agy", "--force", "-v"},
			expected: true,
		},
		{
			name:     "pipeline_fix_subcommand",
			args:     []string{"pipeline-fix", "errors", "agy"},
			expected: true,
		},
		{
			name:     "pipeline_agy_errors_fix",
			args:     []string{"agy", "errors", "fix"},
			expected: true,
		},
		{
			name:     "pipeline_agy_errors_fix_kebab",
			args:     []string{"agy-errors-fix"},
			expected: true,
		},
		{
			name:     "pipeline_aef_shorthand",
			args:     []string{"aef"},
			expected: true,
		},
		{
			name:     "pipeline_fix_agy",
			args:     []string{"fix", "agy"},
			expected: true,
		},
		{
			name:     "pipeline_fix_errors",
			args:     []string{"fix", "errors"},
			expected: true,
		},
		{
			name:     "pipeline_errors_agy",
			args:     []string{"errors", "agy"},
			expected: true,
		},
		{
			name:     "status_command_negative",
			args:     []string{"status"},
			expected: false,
		},
		{
			name:     "logs_command_negative",
			args:     []string{"logs"},
			expected: false,
		},
		{
			name:     "empty_args_negative",
			args:     []string{},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := IsPipelineFixAgyArgs(tc.args)
			if got != tc.expected {
				t.Fatalf("expected IsPipelineFixAgyArgs(%v) = %v, got %v", tc.args, tc.expected, got)
			}
		})
	}
}

func TestRunPipeline_DispatchesToPipelineAgyFixRunner(t *testing.T) {
	origRunner := PipelineAgyFixRunner
	defer func() { PipelineAgyFixRunner = origRunner }()

	called := false
	capturedArgs := []string{}
	PipelineAgyFixRunner = func(args []string) error {
		called = true
		capturedArgs = args
		return nil
	}

	testArgs := []string{"fix", "errors", "agy", "-f"}
	err := runPipeline(testArgs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !called {
		t.Fatalf("expected PipelineAgyFixRunner to be called")
	}

	if len(capturedArgs) != len(testArgs) {
		t.Fatalf("expected captured args %v, got %v", testArgs, capturedArgs)
	}
}
