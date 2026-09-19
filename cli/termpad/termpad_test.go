package termpad

import (
	"bytes"
	"sync"
	"testing"
)

func TestFormatPadded_SingleAndMultiLine(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"hello", "  hello"},
		{"line1\nline2", "  line1\n  line2"},
		{"  already\n  padded", "  already\n  padded"},
		{"line1\n\nline2", "  line1\n\n  line2"},
	}

	for _, tc := range cases {
		got := FormatPadded(tc.input)
		if got != tc.want {
			t.Errorf("FormatPadded(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestPaddingApplied_State(t *testing.T) {
	SetPaddingApplied(false)
	if IsPaddingApplied() {
		t.Error("expected paddingApplied to be false")
	}

	SetPaddingApplied(true)
	if !IsPaddingApplied() {
		t.Error("expected paddingApplied to be true")
	}

	SetPaddingApplied(false)
}

func TestEnsureBottomPadding_Coverage(t *testing.T) {
	SetPaddingApplied(true)
	EnsureBottomPadding("test")

	SetPaddingApplied(false)
	EnsureBottomPadding("test\n\n")
	EnsureBottomPadding("test\n")
	EnsureBottomPadding("test")
}

func TestSmartPaddingWriter_Basic(t *testing.T) {
	var buf bytes.Buffer
	writer := NewSmartPaddingWriter(&buf)

	n, err := writer.Write([]byte("hello\nworld\n"))
	if err != nil || n != len("hello\nworld\n") {
		t.Fatalf("unexpected write error: %v, n=%d", err, n)
	}

	want := "  hello\n  world\n"
	if buf.String() != want {
		t.Errorf("got %q, want %q", buf.String(), want)
	}

	writer.Flush()
	if !IsPaddingApplied() {
		t.Error("expected IsPaddingApplied to be true after flush")
	}
}

func TestSmartPaddingWriter_EmptyWrite(t *testing.T) {
	var buf bytes.Buffer
	writer := NewSmartPaddingWriter(&buf)

	n, err := writer.Write([]byte{})
	if err != nil || n != 0 {
		t.Errorf("expected 0 bytes written, got %d, err=%v", n, err)
	}
}

func TestSmartPaddingWriter_Concurrent(t *testing.T) {
	var buf bytes.Buffer
	writer := NewSmartPaddingWriter(&buf)
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = writer.Write([]byte("concurrent line\n"))
		}()
	}
	wg.Wait()

	if buf.Len() == 0 {
		t.Error("expected non-empty buffer after concurrent writes")
	}
}

func TestIsBoundaryRule(t *testing.T) {
	if !IsBoundaryRule("──────") {
		t.Error("expected ────── to be boundary rule")
	}
	if !IsBoundaryRule("\033[1;94m──────\033[0m") {
		t.Error("expected colored rule to be boundary rule")
	}
	if !IsBoundaryRule("========") {
		t.Error("expected ======== to be boundary rule")
	}
	if IsBoundaryRule("normal text line") {
		t.Error("expected normal text line NOT to be boundary rule")
	}
}

func TestHasExistingMargin(t *testing.T) {
	if !HasExistingMargin("  indented") {
		t.Error("expected 2-space indented to have margin")
	}
	if !HasExistingMargin("\tindented") {
		t.Error("expected tab indented to have margin")
	}
	if !HasExistingMargin("\033[1;94m  colored-indented") {
		t.Error("expected colored indented to have margin")
	}
	if HasExistingMargin("unindented") {
		t.Error("expected unindented NOT to have margin")
	}
}

func TestStripAnsi(t *testing.T) {
	colored := "\033[1;94mhello world\033[0m"
	plain := StripAnsi(colored)
	if plain != "hello world" {
		t.Errorf("StripAnsi(%q) = %q, want 'hello world'", colored, plain)
	}
}

func TestFormatPadded_SuppressRepeatedRules(t *testing.T) {
	input := "header\n──────────\n──────────\nfooter"
	got := FormatPadded(input)
	want := "  header\n  ──────────\n  footer"
	if got != want {
		t.Errorf("FormatPadded rule suppression got %q, want %q", got, want)
	}
}

func TestPrintBlueRule_Suppression(t *testing.T) {
	ResetPaddingState()
	PrintBlueRule(10)
	PrintBlueRule(10)
	ResetPaddingState()
}

func TestPrintSeparator_Suppression(t *testing.T) {
	ResetPaddingState()
	PrintSeparator("──────────")
	PrintSeparator("──────────")
	ResetPaddingState()
}
