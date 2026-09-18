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
	if len(args) < 2 || args[0] != "sh" {
		t.Fatalf("expected resolved shell command for ip, got: %v", args)
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
