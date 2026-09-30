package cmdos

import (
	"testing"
)

func TestParseChangePasswordArgs(t *testing.T) {
	opts, exit, err := parseChangePasswordArgs([]string{"admin", "pass123", "-y", "--dry-run"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exit {
		t.Fatalf("expected exit false")
	}
	if opts.Username != "admin" {
		t.Errorf("expected username 'admin', got '%s'", opts.Username)
	}
	if opts.Password != "pass123" {
		t.Errorf("expected password 'pass123', got '%s'", opts.Password)
	}
	if !opts.IsYes {
		t.Errorf("expected IsYes true")
	}
	if !opts.IsDryRun {
		t.Errorf("expected IsDryRun true")
	}
}

func TestParseChangePasswordArgsHelp(t *testing.T) {
	_, exit, err := parseChangePasswordArgs([]string{"--help"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !exit {
		t.Fatalf("expected exit true for --help")
	}
}

func TestRunChangePasswordCLIDryRun(t *testing.T) {
	err := RunChangePasswordCLI([]string{"testuser", "testpass", "-y", "--dry-run"})
	if err != nil {
		t.Fatalf("dry run failed: %v", err)
	}
}
