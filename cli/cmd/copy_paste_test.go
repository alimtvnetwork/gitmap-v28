package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFilterFlagArgs(t *testing.T) {
	args := []string{"foo", "--file", "sample.txt", "bar"}
	got := filterFlagArgs(args)
	if len(got) != 2 || got[0] != "foo" || got[1] != "bar" {
		t.Fatalf("expected [foo bar], got %v", got)
	}
}

func TestExtractCopyFilePath(t *testing.T) {
	if got := extractCopyFilePath([]string{"--file", "test.txt"}); got != "test.txt" {
		t.Fatalf("expected test.txt, got %s", got)
	}

	if got := extractCopyFilePath([]string{"-f", "test.txt"}); got != "test.txt" {
		t.Fatalf("expected test.txt, got %s", got)
	}

	if got := extractCopyFilePath([]string{"--file=inline.txt"}); got != "inline.txt" {
		t.Fatalf("expected inline.txt, got %s", got)
	}
}

func TestExtractPasteOutFile(t *testing.T) {
	if got := extractPasteOutFile([]string{"--file", "out.txt"}); got != "out.txt" {
		t.Fatalf("expected out.txt, got %s", got)
	}

	if got := extractPasteOutFile([]string{"-o", "out.txt"}); got != "out.txt" {
		t.Fatalf("expected out.txt, got %s", got)
	}

	if got := extractPasteOutFile([]string{"--out=out.txt"}); got != "out.txt" {
		t.Fatalf("expected out.txt, got %s", got)
	}
}

func TestSaveAndFetchMemoryContent(t *testing.T) {
	tempHome := t.TempDir()
	t.Setenv("HOME", tempHome)
	t.Setenv("USERPROFILE", tempHome)

	testText := "gitmap persistent clipboard content test"
	if err := saveContentToMemoryAndClipboard(testText); err != nil {
		t.Fatalf("saveContentToMemoryAndClipboard failed: %v", err)
	}

	memPath, err := resolveMemoryFilePath()
	if err != nil {
		t.Fatalf("resolveMemoryFilePath failed: %v", err)
	}

	data, readErr := os.ReadFile(memPath)
	if readErr != nil || string(data) != testText {
		t.Fatalf("expected %q, got %q (err: %v)", testText, string(data), readErr)
	}

	outPath := filepath.Join(tempHome, "pasted_output.txt")
	if err := writePastedContentToFile(outPath, testText); err != nil {
		t.Fatalf("writePastedContentToFile failed: %v", err)
	}
}
