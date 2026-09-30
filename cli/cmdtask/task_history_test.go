// Package cmdtask — task_history_test.go tests history argument parsing and slice bounds.
package cmdtask

import (
	"testing"
)

func TestParseHistoryArgs(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantLimit  int
		wantOffset int
		wantSec    string
		wantNeg    bool
	}{
		{
			name:       "default values",
			args:       []string{},
			wantLimit:  50,
			wantOffset: 0,
			wantSec:    "all",
			wantNeg:    false,
		},
		{
			name:       "limit and offset",
			args:       []string{"25", "10"},
			wantLimit:  25,
			wantOffset: 10,
			wantSec:    "all",
			wantNeg:    false,
		},
		{
			name:       "negative offset",
			args:       []string{"20", "-5"},
			wantLimit:  20,
			wantOffset: 5,
			wantSec:    "all",
			wantNeg:    true,
		},
		{
			name:       "section flag separate",
			args:       []string{"--section", "macro", "15"},
			wantLimit:  15,
			wantOffset: 0,
			wantSec:    "macro",
			wantNeg:    false,
		},
		{
			name:       "section flag equals",
			args:       []string{"--section=ssh", "30"},
			wantLimit:  30,
			wantOffset: 0,
			wantSec:    "ssh",
			wantNeg:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			win := ParseHistoryArgs(tt.args)
			if win.Limit != tt.wantLimit {
				t.Errorf("Limit = %d, want %d", win.Limit, tt.wantLimit)
			}
			if win.Offset != tt.wantOffset {
				t.Errorf("Offset = %d, want %d", win.Offset, tt.wantOffset)
			}
			if win.Section != tt.wantSec {
				t.Errorf("Section = %q, want %q", win.Section, tt.wantSec)
			}
			if win.IsNegativeOffset != tt.wantNeg {
				t.Errorf("IsNegativeOffset = %v, want %v", win.IsNegativeOffset, tt.wantNeg)
			}
		})
	}
}

func TestCalculateSliceBounds(t *testing.T) {
	total := 100

	// Positive bounds
	winPos := TaskHistoryWindow{Limit: 10, Offset: 20, IsNegativeOffset: false}
	start, end := calculateSliceBounds(total, winPos)
	if start != 20 || end != 30 {
		t.Errorf("positive bounds: got start=%d, end=%d, want 20, 30", start, end)
	}

	// Positive overflow
	winOverflow := TaskHistoryWindow{Limit: 50, Offset: 80, IsNegativeOffset: false}
	start, end = calculateSliceBounds(total, winOverflow)
	if start != 80 || end != 100 {
		t.Errorf("overflow bounds: got start=%d, end=%d, want 80, 100", start, end)
	}

	// Negative bounds
	winNeg := TaskHistoryWindow{Limit: 10, Offset: 15, IsNegativeOffset: true}
	start, end = calculateSliceBounds(total, winNeg)
	if start != 85 || end != 95 {
		t.Errorf("negative bounds: got start=%d, end=%d, want 85, 95", start, end)
	}
}
