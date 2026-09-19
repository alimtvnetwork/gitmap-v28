package cmdssh

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

func TestNormalizeMultiCommands(t *testing.T) {
	cases := []struct {
		input     string
		isWindows bool
		expected  string
	}{
		{"cmd1,cmd2,cmd3", false, "cmd1 && cmd2 && cmd3"},
		{"cmd1, cmd2, cmd3", true, "cmd1; cmd2; cmd3"},
		{"single-cmd", false, "single-cmd"},
		{"single-cmd", true, "single-cmd"},
		{"a, b , c ", false, "a && b && c"},
	}

	for _, tc := range cases {
		actual := normalizeMultiCommands(tc.input, tc.isWindows)
		if actual != tc.expected {
			t.Errorf("normalizeMultiCommands(%q, %v) = %q, want %q", tc.input, tc.isWindows, actual, tc.expected)
		}
	}
}

func TestFilterSSHConns_ExceptAndExclude(t *testing.T) {
	conns := []db.SSHConnection{
		{ID: 1, Alias: "main", IPAddress: "192.168.1.16", Username: "root"},
		{ID: 2, Alias: "worker-1", IPAddress: "192.168.1.14", Username: "ubuntu"},
		{ID: 3, Alias: "worker-2", IPAddress: "192.168.1.8", Username: "admin"},
	}

	filtered := filterSSHConns(conns, "worker-1,192.168.1.8")
	if len(filtered) != 1 || filtered[0].Alias != "main" {
		t.Fatalf("expected 1 connection (main), got %d", len(filtered))
	}

	byId := filterSSHConns(conns, "1")
	if len(byId) != 2 || byId[0].Alias != "worker-1" {
		t.Fatalf("expected exclusion by ID=1, got %d", len(byId))
	}
}
