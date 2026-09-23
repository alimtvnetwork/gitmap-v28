package cmdssh

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

func TestDetermineSSHCommand(t *testing.T) {
	cases := []struct {
		osType       string
		args         []string
		wantShell    string
		wantCommand  string
		wantDelegate bool
	}{
		{"unix", []string{"mkdir", "-p"}, "", "gitmap mkdir -p", true},
		{"windows", []string{"cat", "foo.txt"}, "", "gitmap cat foo.txt", true},
		{"linux", []string{"ssh", "create"}, "", "gitmap ssh create", true},
		{"unix", []string{"ps", "echo test"}, "ps", "echo test", false},
		{"windows", []string{"bash", "ls", "-l"}, "bash", "ls -l", false},
		{"windows", []string{"echo", "test"}, "ps", "echo test", false},
		{"unix", []string{"echo", "test"}, "bash", "echo test", false},
	}

	for _, tc := range cases {
		gotShell, gotCommand, gotDelegate := determineSSHCommand(tc.osType, tc.args)
		if gotShell != tc.wantShell || gotCommand != tc.wantCommand || gotDelegate != tc.wantDelegate {
			t.Errorf("determineSSHCommand(%q, %q) = %q, %q, %v; want %q, %q, %v",
				tc.osType, tc.args, gotShell, gotCommand, gotDelegate, tc.wantShell, tc.wantCommand, tc.wantDelegate)
		}
	}
}

func TestResolveIPCommandArgs(t *testing.T) {
	args := resolveIPCommandArgs([]string{"ip"})
	if len(args) != 1 || args[0] != "ip" {
		t.Fatalf("expected preserved ip command args for delegation, got: %v", args)
	}

	regularArgs := resolveIPCommandArgs([]string{"gitmap", "--version"})
	if len(regularArgs) != 2 || regularArgs[0] != "gitmap" {
		t.Fatalf("expected untouched args for non-ip command, got: %v", regularArgs)
	}
}

func TestResolveExecTargetAndArgs(t *testing.T) {
	conns := []db.SSHConnection{
		{Alias: "w1", IPAddress: "192.168.1.8"},
		{Alias: "w2", IPAddress: "192.168.1.9"},
	}

	opts := seOptions{Args: []string{"w1", "uptime"}}
	resConns, resArgs := resolveExecTargetAndArgs(conns, opts)
	if len(resConns) != 1 || resConns[0].Alias != "w1" {
		t.Fatalf("expected target w1 filtered, got %d conns", len(resConns))
	}
	if len(resArgs) != 1 || resArgs[0] != "uptime" {
		t.Fatalf("expected resArgs to be ['uptime'], got %v", resArgs)
	}
}

func TestIsInteractiveMacroAdd_TrueBasic(t *testing.T) {
	cases := [][]string{
		{"macro", "add", "my-macro"},
		{"macro add my-macro"},
		{"gitmap", "macro", "add", "my-macro"},
		{"gitmap macro add my-macro"},
	}
	assertAllInteractiveMacroAdd(t, cases, true)
}

func TestIsInteractiveMacroAdd_TrueFlags(t *testing.T) {
	cases := [][]string{
		{"devbox", "macro", "add", "my-macro"},
		{"macro", "add", "my-macro", "--desc", "some desc"},
		{"macro", "add", "my-macro", "--pwd"},
		{"macro", "add", "my-macro", "--no-exec"},
	}
	assertAllInteractiveMacroAdd(t, cases, true)
}

func TestIsInteractiveMacroAdd_FalseCommands(t *testing.T) {
	cases := [][]string{
		{"macro", "add", "my-macro", "echo 1"},
		{"macro add my-macro echo 1"},
		{"gitmap", "macro", "add", "my-macro", "uptime"},
		{"macro", "add", "my-macro", "--desc", "some desc", "go build"},
	}
	assertAllInteractiveMacroAdd(t, cases, false)
}

func TestIsInteractiveMacroAdd_FalseOther(t *testing.T) {
	cases := [][]string{
		{"macro", "add", "--help"},
		{"macro", "add", "-h"},
		{"status"},
		{},
	}
	assertAllInteractiveMacroAdd(t, cases, false)
}

func assertAllInteractiveMacroAdd(t *testing.T, cases [][]string, expected bool) {
	for _, c := range cases {
		if isInteractiveMacroAdd(c) != expected {
			t.Errorf("isInteractiveMacroAdd(%v) = %v; want %v", c, !expected, expected)
		}
	}
}

func TestRunSSHExec_InterceptsInteractiveMacroAdd(t *testing.T) {
	err := runSSHExec([]string{"macro", "add", "my-macro"})
	if err == nil {
		t.Fatalf("expected runSSHExec to intercept interactive macro add and return error")
	}
}
