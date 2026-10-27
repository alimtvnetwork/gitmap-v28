package cmdswap

import (
	"testing"
)

func TestParseSwapFlags(t *testing.T) {
	args := []string{
		"--from", "./cli/gitmap.exe",
		"--to", "C:\\bin\\gitmap.exe",
		"--pid", "1234",
		"--kill",
		"--restart",
		"--force",
		"--script-only",
		"--json",
	}

	opts, err := parseSwapFlags(args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if opts.FromPath != "./cli/gitmap.exe" {
		t.Errorf("expected fromPath './cli/gitmap.exe', got '%s'", opts.FromPath)
	}

	if opts.ToPath != "C:\\bin\\gitmap.exe" {
		t.Errorf("expected toPath 'C:\\bin\\gitmap.exe', got '%s'", opts.ToPath)
	}

	if opts.TargetPID != 1234 {
		t.Errorf("expected targetPID 1234, got %d", opts.TargetPID)
	}

	if !opts.HasKill || !opts.HasRestart || !opts.HasForce || !opts.IsScriptOnly || !opts.IsJSONOutput {
		t.Errorf("boolean flags mismatch: %+v", opts)
	}
}

func TestParseSwapFlags_InvalidPID(t *testing.T) {
	args := []string{"--pid", "invalid"}
	_, err := parseSwapFlags(args)
	if err == nil {
		t.Fatal("expected error for invalid PID, got nil")
	}
}
