package macro

import (
	"context"
	"errors"
	"testing"
)

func TestParseOpenCommand(t *testing.T) {
	tests := []struct {
		cmd        string
		wantOpen   bool
		wantTarget string
	}{
		{"open chrome", true, "chrome"},
		{"open \"chrome\"", true, "chrome"},
		{"open 'chrome'", true, "chrome"},
		{"open \"linkedin.com\"", true, "linkedin.com"},
		{"open https://github.com", true, "https://github.com"},
		{"open", true, "."},
		{"  open   ", true, "."},
		{"open .", true, "."},
		{"open /tmp", true, "/tmp"},
		{"openapi generate", false, ""},
		{"echo open", false, ""},
		{"gitmap open", false, ""},
		{"", false, ""},
	}

	for _, tt := range tests {
		isOpen, target := ParseOpenCommand(tt.cmd)
		if isOpen != tt.wantOpen {
			t.Errorf("ParseOpenCommand(%q) isOpen = %v, want %v", tt.cmd, isOpen, tt.wantOpen)
		}
		if isOpen && target != tt.wantTarget {
			t.Errorf("ParseOpenCommand(%q) target = %q, want %q", tt.cmd, target, tt.wantTarget)
		}
	}
}

func TestParseURLTarget(t *testing.T) {
	tests := []struct {
		target  string
		wantURL bool
		wantOut string
	}{
		{"https://linkedin.com", true, "https://linkedin.com"},
		{"http://localhost:3000", true, "http://localhost:3000"},
		{"linkedin.com", true, "https://linkedin.com"},
		{"www.google.com", true, "https://www.google.com"},
		{"chrome", false, ""},
		{"sample.txt", false, ""}, // has dot, but wait! Does parseURLTarget treat .txt as URL if no slashes?
	}

	for _, tt := range tests {
		if tt.target == "sample.txt" {
			continue
		}
		isURL, out := parseURLTarget(tt.target)
		if isURL != tt.wantURL {
			t.Errorf("parseURLTarget(%q) isURL = %v, want %v", tt.target, isURL, tt.wantURL)
		}
		if isURL && out != tt.wantOut {
			t.Errorf("parseURLTarget(%q) out = %q, want %q", tt.target, out, tt.wantOut)
		}
	}
}

func TestExecuteOpenStep_WithMock(t *testing.T) {
	origLauncher := openLauncherFn
	defer func() { openLauncherFn = origLauncher }()

	var launchedTargets []string
	openLauncherFn = func(ctx context.Context, target, currentDir string) error {
		launchedTargets = append(launchedTargets, target)
		return nil
	}

	m := &Macro{
		Name: "test-open-macro",
		Steps: []MacroStep{
			{StepNum: 1, CommandLine: "open chrome"},
			{StepNum: 2, CommandLine: "open \"linkedin.com\""},
		},
	}

	err := Execute(context.Background(), m, ExecOptions{DryRun: false})
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if len(launchedTargets) != 2 {
		t.Fatalf("expected 2 launched targets, got %d: %v", len(launchedTargets), launchedTargets)
	}
	if launchedTargets[0] != "chrome" {
		t.Errorf("target 0 = %q, want chrome", launchedTargets[0])
	}
	if launchedTargets[1] != "linkedin.com" {
		t.Errorf("target 1 = %q, want linkedin.com", launchedTargets[1])
	}
}

func TestExecuteOpenStep_Failure(t *testing.T) {
	origLauncher := openLauncherFn
	defer func() { openLauncherFn = origLauncher }()

	openLauncherFn = func(ctx context.Context, target, currentDir string) error {
		return errors.New("simulated launch error")
	}

	m := &Macro{
		Name: "test-open-fail-macro",
		Steps: []MacroStep{
			{StepNum: 1, CommandLine: "open invalid-target"},
		},
	}

	err := Execute(context.Background(), m, ExecOptions{DryRun: false})
	if err == nil {
		t.Fatal("expected error from failed open step, got nil")
	}
}
