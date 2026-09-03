package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestAppendSSHHelp(t *testing.T) {
	cmd := &cobra.Command{}
	var buf bytes.Buffer
	err := appendSSHHelp(cmd, []string{}, &buf)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "gitmap ssh m1") {
		t.Errorf("expected output to contain 'gitmap ssh m1', got %s", output)
	}

	if !strings.Contains(output, "Alias Resolution") {
		t.Errorf("expected output to contain 'Alias Resolution', got %s", output)
	}
}
