package cmdgit

import (
	"reflect"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func TestBuildToolArguments(t *testing.T) {
	tests := []struct {
		name     string
		tool     string
		dirs     []string
		expected []string
	}{
		{
			name:     "vscode two files with diff flag",
			tool:     constants.CompareToolCode,
			dirs:     []string{"left.go", "right.go"},
			expected: []string{"--wait", "--diff", "left.go", "right.go"},
		},
		{
			name:     "vscode single dir",
			tool:     constants.CompareToolCode,
			dirs:     []string{"repo"},
			expected: []string{"--wait", "repo"},
		},
		{
			name:     "bcomp preserves dirs directly",
			tool:     constants.CompareToolBComp,
			dirs:     []string{"left.go", "right.go"},
			expected: []string{"left.go", "right.go"},
		},
		{
			name:     "meld preserves dirs directly",
			tool:     constants.CompareToolMeld,
			dirs:     []string{"left.go", "right.go"},
			expected: []string{"left.go", "right.go"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := buildToolArguments(tt.tool, tt.dirs)
			if !reflect.DeepEqual(args, tt.expected) {
				t.Errorf("buildToolArguments() = %v, want %v", args, tt.expected)
			}
		})
	}
}

func TestGetBeyondCompareCandidates(t *testing.T) {
	candidates := getBeyondCompareCandidates()
	if len(candidates) == 0 {
		t.Fatal("expected at least one beyond compare candidate")
	}

	if candidates[0] != "bcomp" || candidates[1] != "bcompare" {
		t.Errorf("expected bcomp and bcompare as first candidates, got %v", candidates[:2])
	}
}
