package cmd

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
)

func TestFormatStepCommand_StripsCDir(t *testing.T) {
	testCases := []struct {
		name     string
		step     gitutil.RemediationStep
		expected string
	}{
		{
			name: "git_with_c_flag",
			step: gitutil.RemediationStep{
				Name: "git",
				Args: []string{"-C", "D:/work/my-repo", "stash", "-u"},
			},
			expected: "git stash -u",
		},
		{
			name: "git_pull_with_c_flag",
			step: gitutil.RemediationStep{
				Name: "git",
				Args: []string{"-C", "D:/work/my-repo", "pull"},
			},
			expected: "git pull",
		},
		{
			name: "git_without_c_flag",
			step: gitutil.RemediationStep{
				Name: "git",
				Args: []string{"status", "--short"},
			},
			expected: "git status --short",
		},
		{
			name: "git_only_c_flag",
			step: gitutil.RemediationStep{
				Name: "git",
				Args: []string{"-C"},
			},
			expected: "git -C",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := formatStepCommand(tc.step)
			if actual != tc.expected {
				t.Fatalf("expected %q, got %q", tc.expected, actual)
			}
		})
	}
}
