package macro

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestParseOpenCommand(t *testing.T) {
	tests := getParseOpenTestCases()
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

type parseOpenCase struct {
	cmd        string
	wantOpen   bool
	wantTarget string
}

func getParseOpenTestCases() []parseOpenCase {
	cases := []parseOpenCase{
		{"open chrome", true, "chrome"}, {"open \"chrome\"", true, "chrome"},
		{"open 'chrome'", true, "chrome"}, {"open \"linkedin.com\"", true, "linkedin.com"},
		{"open https://github.com", true, "https://github.com"}, {"open", true, "."},
		{"  open   ", true, "."}, {"open .", true, "."}, {"open /tmp", true, "/tmp"},
		{"openapi generate", false, ""}, {"echo open", false, ""},
		{"gitmap open", false, ""}, {"", false, ""},
	}

	return cases
}

func TestParseURLTarget(t *testing.T) {
	tests := getURLTargetTestCases()
	for _, tt := range tests {
		isURL, out := parseURLTarget(tt.target)
		if isURL != tt.wantURL {
			t.Errorf("parseURLTarget(%q) isURL = %v, want %v", tt.target, isURL, tt.wantURL)
		}

		if isURL && out != tt.wantOut {
			t.Errorf("parseURLTarget(%q) out = %q, want %q", tt.target, out, tt.wantOut)
		}
	}
}

type urlTestCase struct {
	target  string
	wantURL bool
	wantOut string
}

func getURLTargetTestCases() []urlTestCase {
	return []urlTestCase{
		{"https://linkedin.com", true, "https://linkedin.com"},
		{"http://localhost:3000", true, "http://localhost:3000"},
		{"linkedin.com", true, "https://linkedin.com"},
		{"www.google.com", true, "https://www.google.com"},
		{"chrome", false, ""},
		{"sample.txt", false, ""},
		{"sample.csv", false, ""},
		{"main.go", false, ""},
		{"readme.txt", false, ""},
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

	m := buildTestOpenMacro()
	err := Execute(context.Background(), m, ExecOptions{DryRun: false})
	if err != nil || len(launchedTargets) != 2 {
		t.Fatalf("Execute failed: %v, targets: %v", err, launchedTargets)
	}

	assertLaunchedTargets(t, launchedTargets)
}

func buildTestOpenMacro() *Macro {
	return &Macro{
		Name: "test-open-macro",
		Steps: []MacroStep{
			{StepNum: 1, CommandLine: "open chrome"},
			{StepNum: 2, CommandLine: "open \"linkedin.com\""},
		},
	}
}

func assertLaunchedTargets(t *testing.T, targets []string) {
	if targets[0] != "chrome" {
		t.Errorf("target 0 = %q, want chrome", targets[0])
	}

	if targets[1] != "linkedin.com" {
		t.Errorf("target 1 = %q, want linkedin.com", targets[1])
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

func TestDefaultOpenLauncher_LocalFilePrecedence(t *testing.T) {
	tmpDir := t.TempDir()
	localFile := filepath.Join(tmpDir, "readme.txt")
	_ = os.WriteFile(localFile, []byte("hello"), 0644)
	openedPath := mockLaunchPath(t)

	err := defaultOpenLauncher(context.Background(), "readme.txt", tmpDir)
	if err != nil || *openedPath != localFile {
		t.Fatalf("openedPath = %q, want %q, err = %v", *openedPath, localFile, err)
	}
}

func TestDefaultOpenLauncher_DomainLikeLocalFilePrecedence(t *testing.T) {
	tmpDir := t.TempDir()
	domainFile := filepath.Join(tmpDir, "example.com")
	_ = os.WriteFile(domainFile, []byte("local domain content"), 0644)
	openedPath := mockLaunchPath(t)

	err := defaultOpenLauncher(context.Background(), "example.com", tmpDir)
	if err != nil || *openedPath != domainFile {
		t.Fatalf("openedPath = %q, want %q, err = %v", *openedPath, domainFile, err)
	}
}

func mockLaunchPath(t *testing.T) *string {
	var opened string
	orig := launchPathFn
	t.Cleanup(func() { launchPathFn = orig })
	launchPathFn = func(ctx context.Context, p string) error {
		opened = p

		return nil
	}

	return &opened
}
