package cmd

import (
	"strings"
	"testing"
)

func TestBuildPowerShellCmd_Basic(t *testing.T) {
	cmd := buildPowerShellCmd("pwsh", []string{"Write-Host", "hello"})
	if cmd == nil {
		t.Fatal("expected non-nil cmd")
	}

	joined := strings.Join(cmd.Args, " ")
	if !strings.Contains(joined, "Write-Host hello") {
		t.Fatalf("expected command to contain script, got: %s", joined)
	}
}

func TestBuildPowerShellCmd_WithFlagC(t *testing.T) {
	cmd := buildPowerShellCmd("pwsh", []string{"-c", "$a = 123; Write-Host $a"})
	if cmd == nil {
		t.Fatal("expected non-nil cmd")
	}

	joined := strings.Join(cmd.Args, " ")
	if !strings.Contains(joined, "-Command") {
		t.Fatalf("expected -Command in args, got: %s", joined)
	}

	if !strings.Contains(joined, "$a = 123; Write-Host $a") {
		t.Fatalf("expected script to be preserved as command, got: %s", joined)
	}
}

func TestBuildPowerShellCmd_Empty(t *testing.T) {
	cmd := buildPowerShellCmd("pwsh", nil)
	if cmd == nil {
		t.Fatal("expected non-nil cmd")
	}

	if len(cmd.Args) != 1 || cmd.Args[0] != "pwsh" {
		t.Fatalf("expected single argument pwsh, got: %v", cmd.Args)
	}
}
