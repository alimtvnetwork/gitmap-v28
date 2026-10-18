package cmdgit

import (
	"reflect"
	"testing"
)

func TestParseGitDirAndArgs(t *testing.T) {
	tests := []struct {
		name         string
		input        []string
		expectedDir  string
		expectedArgs []string
	}{
		{
			name:         "empty args",
			input:        []string{},
			expectedDir:  "",
			expectedArgs: nil,
		},
		{
			name:         "no dir flag",
			input:        []string{"HEAD~1", "--", "readme.md"},
			expectedDir:  "",
			expectedArgs: []string{"HEAD~1", "--", "readme.md"},
		},
		{
			name:         "flag -C with path",
			input:        []string{"-C", "some/repo", "HEAD"},
			expectedDir:  "some/repo",
			expectedArgs: []string{"HEAD"},
		},
		{
			name:         "flag --dir with path",
			input:        []string{"--dir", "other/repo", "-n", "5"},
			expectedDir:  "other/repo",
			expectedArgs: []string{"-n", "5"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir, args := ParseGitDirAndArgs(tt.input)
			if dir != tt.expectedDir {
				t.Errorf("ParseGitDirAndArgs dir = %q, want %q", dir, tt.expectedDir)
			}
			if !reflect.DeepEqual(args, tt.expectedArgs) {
				t.Errorf("ParseGitDirAndArgs args = %v, want %v", args, tt.expectedArgs)
			}
		})
	}
}

func TestBuildLogArgs(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:  "empty input adds default -n 10 and format",
			input: []string{},
			expected: []string{
				"log",
				"-n", "10",
				"--date=short",
				"--pretty=format:%C(yellow)%h%C(reset) %C(green)%cd%C(reset) %C(cyan)%an%C(reset) %s%C(auto)%d%C(reset)",
			},
		},
		{
			name:  "custom limit preserved",
			input: []string{"-n", "5"},
			expected: []string{
				"log",
				"--date=short",
				"--pretty=format:%C(yellow)%h%C(reset) %C(green)%cd%C(reset) %C(cyan)%an%C(reset) %s%C(auto)%d%C(reset)",
				"-n", "5",
			},
		},
		{
			name:     "oneline flag passes through without default pretty format",
			input:    []string{"--oneline"},
			expected: []string{"log", "-n", "10", "--oneline"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := buildLogArgs(tt.input)
			if !reflect.DeepEqual(args, tt.expected) {
				t.Errorf("buildLogArgs() = %v, want %v", args, tt.expected)
			}
		})
	}
}
