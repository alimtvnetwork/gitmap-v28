package cmd_test

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/cmd"
)

func TestFormatSeq_TwoDigits(t *testing.T) {
	tests := []struct {
		name   string
		seq    int
		digits int
		want   string
	}{
		{"single digit padded", 1, 2, "01"},
		{"double digit", 12, 2, "12"},
		{"max two digit", 99, 2, "99"},
		{"three digit pad", 5, 3, "005"},
		{"four digit pad", 42, 4, "0042"},
		{"no padding needed", 100, 3, "100"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cmd.FormatSeq(tt.seq, tt.digits)
			if got != tt.want {
				t.Errorf("FormatSeq(%d, %d) = %q, want %q", tt.seq, tt.digits, got, tt.want)
			}
		})
	}
}

func TestParseVersionPattern_Valid(t *testing.T) {
	tests := []struct {
		name       string
		pattern    string
		wantPrefix string
		wantDigits int
	}{
		{"two dollar", "v1.$$", "v1.", 2},
		{"three dollar", "v2.$$$", "v2.", 3},
		{"four dollar", "v1.0.$$$$", "v1.0.", 4},
		{"single dollar", "v$", "v", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prefix, digits := cmd.ParseVersionPatternSafe(tt.pattern)
			if prefix != tt.wantPrefix {
				t.Errorf("prefix = %q, want %q", prefix, tt.wantPrefix)
			}
			if digits != tt.wantDigits {
				t.Errorf("digits = %d, want %d", digits, tt.wantDigits)
			}
		})
	}
}

func TestParseVersionPattern_NoDollar(t *testing.T) {
	_, digits := cmd.ParseVersionPatternSafe("v1.0")
	if digits != 0 {
		t.Errorf("expected 0 digits for pattern without $, got %d", digits)
	}
}

func TestResolveTRBranch_WithPrefix(t *testing.T) {
	got := cmd.ResolveTRBranchExported("temp-release/v1.05")
	if got != "temp-release/v1.05" {
		t.Errorf("got %q, want unchanged", got)
	}
}

func TestResolveTRBranch_WithoutPrefix(t *testing.T) {
	got := cmd.ResolveTRBranchExported("v1.05")
	if got != "temp-release/v1.05" {
		t.Errorf("got %q, want temp-release/v1.05", got)
	}
}

func TestCheckSequenceRange_Valid(t *testing.T) {
	err := cmd.CheckSequenceRange(1, 10, 2)
	if err != nil {
		t.Errorf("expected valid range, got error: %v", err)
	}
}

func TestCheckSequenceRange_Overflow(t *testing.T) {
	err := cmd.CheckSequenceRange(90, 15, 2)
	if err == nil {
		t.Error("expected overflow error, got nil")
	}
}

func TestCheckSequenceRange_ExactMax(t *testing.T) {
	err := cmd.CheckSequenceRange(90, 10, 2)
	if err != nil {
		t.Errorf("expected valid at boundary, got error: %v", err)
	}
}
