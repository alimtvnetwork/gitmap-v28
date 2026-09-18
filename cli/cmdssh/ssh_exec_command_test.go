package cmdssh

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
)

func TestDetermineSSHCommand_GitmapDelegation(t *testing.T) {
	cases := []struct {
		osType      string
		args        []string
		wantCommand string
	}{
		{"unix", []string{"status"}, "gitmap status"},
		{"unix", []string{"gitmap", "status"}, "gitmap status"},
		{"windows", []string{"pipeline", "run"}, "gitmap pipeline run"},
		{"linux", []string{"clone", "https://github.com/repo"}, "gitmap clone https://github.com/repo"},
		{"linux", []string{"agy", "fix-pipeline"}, "gitmap agy fix-pipeline"},
		{"linux", []string{"aef"}, "gitmap aef"},
	}

	for _, tc := range cases {
		shell, cmd, isDelegate := determineSSHCommand(tc.osType, tc.args)
		if !isDelegate || shell != "" || cmd != tc.wantCommand {
			t.Errorf("determineSSHCommand(%q, %v) = (%q, %q, %v), want ('', %q, true)",
				tc.osType, tc.args, shell, cmd, isDelegate, tc.wantCommand)
		}
	}
}

func TestDetermineSSHCommand_NativeShells(t *testing.T) {
	cases := []struct {
		osType      string
		args        []string
		wantShell   string
		wantCommand string
	}{
		{"unix", []string{"uname", "-a"}, "bash", "uname -a"},
		{"windows", []string{"Get-Process"}, "ps", "Get-Process"},
		{"linux", []string{"bash", "echo", "hi"}, "bash", "echo hi"},
		{"windows", []string{"cmd", "dir"}, "cmd", "dir"},
		{"unix", []string{"sh", "uptime"}, "sh", "uptime"},
		{"unix", []string{"ps", "aux"}, "ps", "aux"},
	}

	for _, tc := range cases {
		shell, cmd, isDelegate := determineSSHCommand(tc.osType, tc.args)
		if isDelegate || shell != tc.wantShell || cmd != tc.wantCommand {
			t.Errorf("determineSSHCommand(%q, %v) = (%q, %q, %v), want (%q, %q, false)",
				tc.osType, tc.args, shell, cmd, isDelegate, tc.wantShell, tc.wantCommand)
		}
	}
}

func TestDetermineSSHCommand_ChainedAndEmpty(t *testing.T) {
	shell, cmd, isDelegate := determineSSHCommand("unix", []string{})
	if isDelegate || shell != "" || cmd != "" {
		t.Errorf("empty args should return zero values, got (%q, %q, %v)", shell, cmd, isDelegate)
	}

	shell, cmd, isDelegate = determineSSHCommand("unix", []string{"uname -a && df -h"})
	if isDelegate || shell != "bash" || cmd != "uname -a && df -h" {
		t.Errorf("chained command failed, got (%q, %q, %v)", shell, cmd, isDelegate)
	}
}

func TestIsGitmapCommand(t *testing.T) {
	validCommands := []string{
		"gitmap", "status", "pipeline", "pipe", "pl", "clone", "pull",
		"sync", "push", "clean", "log", "branch", "diff", "storage",
		"macro", "install", "update", "setup", "chrome", "vscode",
		"vsc", "vhost", "zip", "service", "os", "schedule", "schedules",
		"agy", "ag", "antigravity", "aef", "fix-pipeline", "pt",
		"ssh", "se", "sj", "cluster", "sc", "mkdir", "cat",
	}

	for _, cmd := range validCommands {
		if !isGitmapCommand(cmd) {
			t.Errorf("expected isGitmapCommand(%q) to be true", cmd)
		}
	}

	invalidCommands := []string{"uptime", "whoami", "curl", "python", "node", "dir"}
	for _, cmd := range invalidCommands {
		if isGitmapCommand(cmd) {
			t.Errorf("expected isGitmapCommand(%q) to be false", cmd)
		}
	}
}

func TestResolveGitmapCommandString(t *testing.T) {
	if got := resolveGitmapCommandString([]string{}); got != "gitmap" {
		t.Errorf("resolveGitmapCommandString([]) = %q, want 'gitmap'", got)
	}

	withPrefix := []string{"gitmap", "status", "--json"}
	if got := resolveGitmapCommandString(withPrefix); got != "gitmap status --json" {
		t.Errorf("resolveGitmapCommandString(%v) = %q, want 'gitmap status --json'", withPrefix, got)
	}

	withoutPrefix := []string{"status", "--json"}
	if got := resolveGitmapCommandString(withoutPrefix); got != "gitmap status --json" {
		t.Errorf("resolveGitmapCommandString(%v) = %q, want 'gitmap status --json'", withoutPrefix, got)
	}
}

func TestParseSEFlags_Options(t *testing.T) {
	args := []string{"--target", "node1", "--exclude", "worker2,worker3", "--ip", "10.0.0.5", "uptime"}
	opts := parseSEFlags(args)

	if opts.Target != "node1" || opts.Exclude != "worker2,worker3" || opts.IP != "10.0.0.5" {
		t.Errorf("parseSEFlags options mismatch: %+v", opts)
	}

	if len(opts.Args) != 1 || opts.Args[0] != "uptime" {
		t.Errorf("parseSEFlags args mismatch: %v", opts.Args)
	}
}

func TestValidateSEArgs_EmptyExits(t *testing.T) {
	var exitCode int
	prevExit := cliexit.SetExitFunc(func(code int) {
		exitCode = code
	})
	defer cliexit.SetExitFunc(prevExit)

	validateSEArgs([]string{})
	if exitCode != 1 {
		t.Errorf("expected exit code 1 on empty args, got %d", exitCode)
	}
}

func TestExtractFirstToken(t *testing.T) {
	if got := extractFirstToken("gitmap status"); got != "gitmap" {
		t.Errorf("expected 'gitmap', got %q", got)
	}
	if got := extractFirstToken("  status -v  "); got != "status" {
		t.Errorf("expected 'status', got %q", got)
	}
	if got := extractFirstToken(""); got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestDetermineSSHCommand_QuotedGitmapMultiCommand(t *testing.T) {
	cases := []struct {
		input       string
		wantCommand string
	}{
		{"gitmap status && gitmap pipeline", "gitmap status && gitmap pipeline"},
		{"status --json", "gitmap status --json"},
		{"pipeline errors agy fix", "gitmap pipeline errors agy fix"},
	}

	for _, tc := range cases {
		shell, cmd, isDelegate := determineSSHCommand("linux", []string{tc.input})
		if !isDelegate || shell != "" || cmd != tc.wantCommand {
			t.Errorf("determineSSHCommand(%q) = (%q, %q, %v), want ('', %q, true)",
				tc.input, shell, cmd, isDelegate, tc.wantCommand)
		}
	}
}

func TestDetermineSSHCommand_QuotedExplicitShell(t *testing.T) {
	shell, cmd, isDelegate := determineSSHCommand("linux", []string{"bash echo hello"})
	if isDelegate || shell != "bash" || cmd != "echo hello" {
		t.Errorf("quoted explicit shell failed, got (%q, %q, %v)", shell, cmd, isDelegate)
	}
}

