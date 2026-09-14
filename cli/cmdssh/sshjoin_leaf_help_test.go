package cmdssh

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestGetLeafHelp_Subcommands(t *testing.T) {
	commands := []string{"add", "scan", "status", "ls", "rm", "add-auth", "history"}
	for _, cmd := range commands {
		help, hasHelp := GetLeafHelp(cmd)
		if !hasHelp || help == "" {
			t.Errorf("expected help for command %s, got none", cmd)
		}
		if !strings.Contains(help, "Usage:") || !strings.Contains(help, "Examples:") {
			t.Errorf("command %s help missing Usage or Examples section: %s", cmd, help)
		}
	}
}

func TestGetLeafHelp_Aliases(t *testing.T) {
	aliases := []string{"find", "discover", "probe", "ping", "health", "check", "list", "auth", "hist"}
	for _, a := range aliases {
		help, hasHelp := GetLeafHelp(a)
		if !hasHelp || help == "" {
			t.Errorf("expected help for alias %s, got none", a)
		}
	}
}

func TestGetLeafHelp_Unknown(t *testing.T) {
	help, hasHelp := GetLeafHelp("nonexistent-cmd")
	if hasHelp || help != "" {
		t.Errorf("expected false for unknown command, got %v: %s", hasHelp, help)
	}
}

func TestPrintLeafHelp_Valid(t *testing.T) {
	var buf bytes.Buffer
	err := PrintLeafHelp(&buf, "scan")
	if err != nil {
		t.Fatalf("unexpected error printing leaf help: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "gitmap sj scan") {
		t.Errorf("expected scan usage in output: %s", out)
	}
}

func TestPrintLeafHelp_Unknown(t *testing.T) {
	var buf bytes.Buffer
	err := PrintLeafHelp(&buf, "invalid-cmd")
	if err == nil {
		t.Fatal("expected error for invalid command, got nil")
	}
}

func TestAttachLeafHelp_Integration(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	AttachLeafHelp(cmd, "status")
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.HelpFunc()(cmd, nil)
	out := buf.String()
	if !strings.Contains(out, "gitmap sj status") {
		t.Errorf("expected attached help to render status usage, got: %s", out)
	}
}
