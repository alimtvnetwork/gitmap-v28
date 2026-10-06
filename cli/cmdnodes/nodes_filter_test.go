package cmdnodes

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

func TestFilterFleetNodes_DefaultExcludeMain(t *testing.T) {
	conns := []db.SSHConnection{
		{Alias: "local", IPAddress: "127.0.0.1"},
		{Alias: "main", IPAddress: "10.254.1.10"},
		{Alias: "w1", IPAddress: "10.254.1.11"},
		{Alias: "w2", IPAddress: "10.254.1.12"},
	}

	opts := NodeFilterOptions{}
	filtered := FilterFleetNodes(conns, opts)

	if len(filtered) != 2 {
		t.Fatalf("expected 2 nodes (w1, w2), got %d", len(filtered))
	}

	for _, c := range filtered {
		if c.Alias == "main" {
			t.Errorf("main node was not excluded by default")
		}

		if c.Alias == "local" {
			t.Errorf("local node was not excluded")
		}
	}
}

func TestFilterFleetNodes_IncludeMain(t *testing.T) {
	conns := []db.SSHConnection{
		{Alias: "local", IPAddress: "127.0.0.1"},
		{Alias: "main", IPAddress: "10.254.1.10"},
		{Alias: "w1", IPAddress: "10.254.1.11"},
	}

	opts := ParseNodeFilterOptions([]string{"--include-main"})
	if !opts.IncludeMain {
		t.Fatalf("expected IncludeMain to be true")
	}

	filtered := FilterFleetNodes(conns, opts)
	if len(filtered) != 2 {
		t.Fatalf("expected 2 nodes (main, w1), got %d", len(filtered))
	}

	hasMain := false
	for _, c := range filtered {
		if c.Alias == "main" {
			hasMain = true
		}
	}

	if !hasMain {
		t.Errorf("expected main node to be kept when --include-main is set")
	}
}

func TestFilterFleetNodes_ExceptCommaSeparated(t *testing.T) {
	conns := []db.SSHConnection{
		{Alias: "w1", IPAddress: "10.254.1.11"},
		{Alias: "w2", IPAddress: "10.254.1.12"},
		{Alias: "w3", IPAddress: "10.254.1.13"},
	}

	opts := ParseNodeFilterOptions([]string{"--except", "w1,w2"})
	if len(opts.Except) != 2 {
		t.Fatalf("expected 2 except tokens, got %d", len(opts.Except))
	}

	filtered := FilterFleetNodes(conns, opts)
	if len(filtered) != 1 {
		t.Fatalf("expected 1 node (w3), got %d", len(filtered))
	}

	if filtered[0].Alias != "w3" {
		t.Errorf("expected remaining node to be w3, got %s", filtered[0].Alias)
	}
}

func TestFilterFleetNodes_Target(t *testing.T) {
	conns := []db.SSHConnection{
		{Alias: "u1", IPAddress: "10.254.1.21"},
		{Alias: "u2", IPAddress: "10.254.1.22"},
	}

	opts := ParseNodeFilterOptions([]string{"--target", "u1"})
	if opts.Target != "u1" {
		t.Fatalf("expected Target to be u1, got %s", opts.Target)
	}

	filtered := FilterFleetNodes(conns, opts)
	if len(filtered) != 1 {
		t.Fatalf("expected 1 node, got %d", len(filtered))
	}

	if filtered[0].Alias != "u1" {
		t.Errorf("expected target u1, got %s", filtered[0].Alias)
	}
}

func TestFilterFleetNodes_TargetPositionalAndShort(t *testing.T) {
	optsShort := ParseNodeFilterOptions([]string{"-t", "u1"})
	if optsShort.Target != "u1" {
		t.Errorf("expected -t to set Target=u1, got %s", optsShort.Target)
	}

	optsPos := ParseNodeFilterOptions([]string{"u1"})
	if optsPos.Target != "u1" {
		t.Errorf("expected positional argument to set Target=u1, got %s", optsPos.Target)
	}
}

func TestFilterFleetNodes_SkipLocalMachine(t *testing.T) {
	conns := []db.SSHConnection{
		{Alias: "local", IPAddress: "127.0.0.1"},
		{Alias: "my-local", IPAddress: "localhost"},
		{Alias: "remote-worker", IPAddress: "192.168.10.20"},
	}

	filtered := FilterFleetNodes(conns, NodeFilterOptions{})
	if len(filtered) != 1 {
		t.Fatalf("expected only 1 remote node, got %d", len(filtered))
	}

	if filtered[0].Alias != "remote-worker" {
		t.Errorf("expected remote-worker, got %s", filtered[0].Alias)
	}
}

func TestFilterFleetNodes_IncludeWhitelist(t *testing.T) {
	conns := []db.SSHConnection{
		{Alias: "w1", IPAddress: "10.254.1.11"},
		{Alias: "w2", IPAddress: "10.254.1.12"},
		{Alias: "w3", IPAddress: "10.254.1.13"},
	}

	opts := ParseNodeFilterOptions([]string{"--include", "w1,w3"})
	filtered := FilterFleetNodes(conns, opts)

	if len(filtered) != 2 {
		t.Fatalf("expected 2 whitelisted nodes, got %d", len(filtered))
	}

	for _, c := range filtered {
		if c.Alias == "w2" {
			t.Errorf("non-whitelisted node w2 should not be present")
		}
	}
}

func TestParseNodeFilterOptions_AllFlags(t *testing.T) {
	args := []string{
		"-t", "w1",
		"-e", "w2,w3",
		"--accept", "w1,w4",
		"--include-main",
		"--open-only",
	}

	opts := ParseNodeFilterOptions(args)

	if opts.Target != "w1" {
		t.Errorf("expected Target w1, got %s", opts.Target)
	}

	if len(opts.Except) != 2 || opts.Except[0] != "w2" || opts.Except[1] != "w3" {
		t.Errorf("expected Except [w2 w3], got %v", opts.Except)
	}

	if len(opts.Include) != 2 || opts.Include[0] != "w1" || opts.Include[1] != "w4" {
		t.Errorf("expected Include [w1 w4], got %v", opts.Include)
	}

	if !opts.IncludeMain {
		t.Errorf("expected IncludeMain to be true")
	}

	if !opts.OpenOnly {
		t.Errorf("expected OpenOnly to be true")
	}
}

func TestFilterFleetNodes_OpenOnly_Offline(t *testing.T) {
	conns := []db.SSHConnection{
		{Alias: "offline-node", IPAddress: "192.0.2.1:54321"},
	}

	opts := NodeFilterOptions{OpenOnly: true}
	filtered := FilterFleetNodes(conns, opts)

	if len(filtered) != 0 {
		t.Errorf("expected unreachable node to be excluded with OpenOnly=true, got %d", len(filtered))
	}
}

func TestRunNodesScan_Help(t *testing.T) {
	err := RunNodesScan([]string{"--help"})
	if err != nil {
		t.Errorf("expected no error for --help, got: %v", err)
	}
}

func TestRunNodesRescan_Help(t *testing.T) {
	err := RunNodesRescan([]string{"--help"})
	if err != nil {
		t.Errorf("expected no error for --help, got: %v", err)
	}
}
