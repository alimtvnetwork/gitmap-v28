package cmd

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestRunNoArgsDoesNotPanic(t *testing.T) {
	origArgs := os.Args
	defer func() {
		os.Args = origArgs
	}()

	os.Args = []string{"gitmap"}

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Run() panicked with no arguments: %v", r)
		}
	}()

	Run()
}

func TestUsageCompactLineCountAndSuggestions(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe() failed: %v", err)
	}

	origStdout := os.Stdout
	os.Stdout = w

	printUsageCompact()

	w.Close()
	os.Stdout = origStdout

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	lines := strings.Split(strings.TrimRight(output, "\r\n"), "\n")
	lineCount := len(lines)

	if lineCount > 18 {
		t.Errorf("expected compact root help line count <= 18, got %d:\n%s", lineCount, output)
	}

	if !strings.Contains(output, "Suggestions:") {
		t.Errorf("expected output to contain 'Suggestions:', got:\n%s", output)
	}
	if !strings.Contains(output, "gitmap help") {
		t.Errorf("expected output to contain 'gitmap help', got:\n%s", output)
	}
	if !strings.Contains(output, "gitmap llm train") {
		t.Errorf("expected output to contain 'gitmap llm train', got:\n%s", output)
	}
	if !strings.Contains(output, "Mandatory") {
		t.Errorf("expected output to contain 'Mandatory', got:\n%s", output)
	}
}
