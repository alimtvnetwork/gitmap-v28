package cmd

import (
	"testing"
)

func TestDetermineSSHCommand(t *testing.T) {
	cases := []struct {
		osType           string
		args             []string
		wantShell        string
		wantCommand      string
		wantDelegate     bool
	}{
		{"unix", []string{"mkdir", "-p"}, "", "", true},
		{"windows", []string{"cat", "foo.txt"}, "", "", true},
		{"linux", []string{"ssh", "create"}, "", "", true},
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
