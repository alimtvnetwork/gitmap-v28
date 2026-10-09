package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestAppendPRHelp(t *testing.T) {
	var buf bytes.Buffer
	err := appendPRHelp(nil, nil, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "## PR (Pull Request)") {
		t.Errorf("expected PR header, got %s", out)
	}
	if !strings.Contains(out, "gitmap pr LEFT RIGHT") {
		t.Errorf("expected usage command, got %s", out)
	}
	if !strings.Contains(out, "Final Snapshot Synchronization") {
		t.Errorf("expected snapshot sync section, got %s", out)
	}
}
