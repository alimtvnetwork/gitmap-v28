// Package macro — safe_rm_test.go tests idempotent command adaptation for macro steps.
package macro

import (
	"runtime"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func TestAdaptCommandForPlatform_NonWindows(t *testing.T) {
	if runtime.GOOS == constants.OSWindows {
		t.Skip("skipping non-Windows check on Windows")
	}

	cmd := "rm -rf test"
	adapted := AdaptCommandForPlatform(cmd)
	if adapted != cmd {
		t.Fatalf("expected unchanged command on non-Windows, got %q", adapted)
	}
}

func TestAdaptCommandForPlatform_Windows(t *testing.T) {
	if runtime.GOOS != constants.OSWindows {
		t.Skip("skipping Windows-specific check on non-Windows")
	}

	tests := []struct {
		name       string
		input      string
		wantSubstr string
		isChanged  bool
	}{
		{
			name:       "rm single target",
			input:      "rm test",
			wantSubstr: "Test-Path -LiteralPath $__target",
			isChanged:  true,
		},
		{
			name:       "rm with flags and multiple targets",
			input:      "rm -rf dir1 dir2",
			wantSubstr: "@('dir1', 'dir2')",
			isChanged:  true,
		},
		{
			name:       "rmdir with slash flags",
			input:      "rmdir /s /q somedir",
			wantSubstr: "@('somedir')",
			isChanged:  true,
		},
		{
			name:       "Remove-Item target",
			input:      "Remove-Item C:\\Temp\\dummy",
			wantSubstr: "Remove-Item -Recurse -Force",
			isChanged:  true,
		},
		{
			name:       "del single target",
			input:      "del test.txt",
			wantSubstr: "@('test.txt')",
			isChanged:  true,
		},
		{
			name:       "erase single target",
			input:      "erase dummy.log",
			wantSubstr: "@('dummy.log')",
			isChanged:  true,
		},
		{
			name:       "non-removal command echo",
			input:      "echo hello world",
			wantSubstr: "echo hello world",
			isChanged:  false,
		},
		{
			name:       "non-removal command mkdir",
			input:      "mkdir -p test",
			wantSubstr: "mkdir -p test",
			isChanged:  false,
		},
		{
			name:       "empty input",
			input:      "",
			wantSubstr: "",
			isChanged:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapted := AdaptCommandForPlatform(tt.input)
			if tt.isChanged && adapted == tt.input {
				t.Fatalf("expected command %q to be adapted, but remained unchanged", tt.input)
			}
			if !tt.isChanged && adapted != tt.input {
				t.Fatalf("expected command %q to remain unchanged, got %q", tt.input, adapted)
			}
			if tt.wantSubstr != "" && !strings.Contains(adapted, tt.wantSubstr) {
				t.Fatalf("expected adapted command to contain %q, got %q", tt.wantSubstr, adapted)
			}
		})
	}
}

func TestIsWindowsRemovalCmd(t *testing.T) {
	cases := []struct {
		cmd      string
		expected bool
	}{
		{"rm foo", true},
		{"rm\tbar", true},
		{"rmdir foo", true},
		{"rmdir\tbar", true},
		{"remove-item foo", true},
		{"rd foo", true},
		{"del foo.txt", true},
		{"del\tfoo.txt", true},
		{"erase foo.txt", true},
		{"erase\tfoo.txt", true},
		{"echo rm foo", false},
		{"mkdir foo", false},
		{"git status", false},
	}

	for _, tc := range cases {
		actual := isWindowsRemovalCmd(tc.cmd)
		if actual != tc.expected {
			t.Errorf("isWindowsRemovalCmd(%q) = %v; want %v", tc.cmd, actual, tc.expected)
		}
	}
}

func TestExtractRemovalTargets(t *testing.T) {
	args := []string{"-rf", `"-Force"`, "/s", "/q", "folder1", "'folder2'", `"folder3"`}
	targets := extractRemovalTargets(args)

	if len(targets) != 3 {
		t.Fatalf("expected 3 targets, got %d: %v", len(targets), targets)
	}

	expected := []string{"'folder1'", "'folder2'", "'folder3'"}
	for i, exp := range expected {
		if targets[i] != exp {
			t.Errorf("target[%d] = %q, want %q", i, targets[i], exp)
		}
	}
}
