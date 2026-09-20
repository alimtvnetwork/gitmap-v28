// Package cmd — author_sponsor_test.go tests author and sponsor commands.
package cmd

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestRunAuthor(t *testing.T) {
	err := runAuthor([]string{})
	if err != nil {
		t.Fatalf("runAuthor failed: %v", err)
	}
}

func TestRunSponsor(t *testing.T) {
	err := runSponsor([]string{})
	if err != nil {
		t.Fatalf("runSponsor failed: %v", err)
	}
}

func TestRunCredits(t *testing.T) {
	err := runCredits([]string{})
	if err != nil {
		t.Fatalf("runCredits failed: %v", err)
	}
}

func TestRunSponsorContent(t *testing.T) {
	out := captureSponsorOutput(t)
	if !strings.Contains(out, "RISEUP ASIA LLC") {
		t.Errorf("expected RISEUP ASIA LLC in sponsor output, got: %s", out)
	}
	if !strings.Contains(out, "https://riseup-asia.com") {
		t.Errorf("expected https://riseup-asia.com in sponsor output, got: %s", out)
	}
}

func captureSponsorOutput(t *testing.T) string {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	_ = runSponsor([]string{})
	_ = w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String()
}
