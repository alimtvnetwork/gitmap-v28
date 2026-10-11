package cmdexclude

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

func TestEvaluatePurgeReason(t *testing.T) {
	tests := []struct {
		name     string
		record   model.ScanRecord
		expected string
	}{
		{
			name: "empty path",
			record: model.ScanRecord{
				AbsolutePath: "",
			},
			expected: "empty path",
		},
		{
			name: "non existent path",
			record: model.ScanRecord{
				AbsolutePath: "C:/non/existent/path/for/sure",
			},
			expected: "path does not exist",
		},
		{
			name: "cursor plugin cache",
			record: model.ScanRecord{
				AbsolutePath: "C:/Users/Administrator/.cursor/plugins/cache/cursor-public/figma/f04932",
			},
			expected: "path does not exist",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			reason := evaluatePurgeReason(tc.record)
			if reason != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, reason)
			}
		})
	}
}
