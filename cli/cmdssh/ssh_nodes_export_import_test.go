package cmdssh

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

func TestFilterSSHConnectionsByExcept_IDIPAndAlias(t *testing.T) {
	conns := []db.SSHConnection{
		{Alias: "node-1", IPAddress: "192.168.1.10"},
		{Alias: "node-2", IPAddress: "192.168.1.20"},
		{Alias: "node-3", IPAddress: "192.168.1.30"},
	}
	res := FilterSSHConnectionsByExcept(conns, "worker-1,192.168.1.30")
	if len(res) != 1 || res[0].Alias != "node-2" {
		t.Fatalf("expected only node-2 to remain, got %+v", res)
	}
}
