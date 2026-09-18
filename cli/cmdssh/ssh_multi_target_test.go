package cmdssh

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

func createSampleConns() []db.SSHConnection {
	return []db.SSHConnection{
		{Alias: "devbox", IPAddress: "192.168.1.10", Username: "dev"},
		{Alias: "worker-1", IPAddress: "192.168.1.20", Username: "worker"},
		{Alias: "worker-2", IPAddress: "192.168.1.30", Username: "worker"},
	}
}

func TestFilterConnectionsByTarget_SingleAndAll(t *testing.T) {
	conns := createSampleConns()

	allConns := filterConnectionsByTarget(conns, "all")
	if len(allConns) != 3 {
		t.Errorf("expected 3 connections for 'all', got %d", len(allConns))
	}

	single := filterConnectionsByTarget(conns, "devbox")
	if len(single) != 1 || single[0].Alias != "devbox" {
		t.Errorf("expected single match devbox, got %+v", single)
	}

	byIP := filterConnectionsByTarget(conns, "192.168.1.20")
	if len(byIP) != 1 || byIP[0].Alias != "worker-1" {
		t.Errorf("expected single match by IP, got %+v", byIP)
	}
}

func TestFilterConnectionsByTarget_MultiComma(t *testing.T) {
	conns := createSampleConns()

	multi := filterConnectionsByTarget(conns, "devbox,worker-2")
	if len(multi) != 2 {
		t.Fatalf("expected 2 matches for comma list, got %d", len(multi))
	}

	aliases := []string{multi[0].Alias, multi[1].Alias}
	if aliases[0] != "devbox" || aliases[1] != "worker-2" {
		t.Errorf("unexpected matches: %v", aliases)
	}
}

func TestFilterConnectionsByTarget_MultiMixedAliasIP(t *testing.T) {
	conns := createSampleConns()

	mixed := filterConnectionsByTarget(conns, "devbox, 192.168.1.30")
	if len(mixed) != 2 {
		t.Fatalf("expected 2 matches for mixed alias and IP, got %d", len(mixed))
	}

	none := filterConnectionsByTarget(conns, "nonexistent, 10.0.0.99")
	if len(none) != 0 {
		t.Errorf("expected 0 matches for unknown targets, got %d", len(none))
	}
}
