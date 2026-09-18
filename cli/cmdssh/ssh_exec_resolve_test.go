package cmdssh

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

func mockConnections() []db.SSHConnection {
	return []db.SSHConnection{
		{Alias: "devbox", IPAddress: "192.168.1.14"},
		{Alias: "worker-1", IPAddress: "192.168.1.20"},
	}
}

func TestResolveExecTargetAndArgs_SingleTarget(t *testing.T) {
	conns := mockConnections()
	opts := seOptions{Args: []string{"devbox", "uptime"}}
	resConns, args := resolveExecTargetAndArgs(conns, opts)
	if len(resConns) != 1 || resConns[0].Alias != "devbox" || len(args) != 1 || args[0] != "uptime" {
		t.Fatalf("unexpected result: conns=%+v args=%+v", resConns, args)
	}
}

func TestResolveExecTargetAndArgs_CommaMultiTarget(t *testing.T) {
	conns := mockConnections()
	opts := seOptions{Args: []string{"devbox,worker-1", "free -m"}}
	resConns, args := resolveExecTargetAndArgs(conns, opts)
	if len(resConns) != 2 || len(args) != 1 || args[0] != "free -m" {
		t.Fatalf("unexpected result: conns=%+v args=%+v", resConns, args)
	}
}

func TestResolveExecTargetAndArgs_SpaceMultiTarget(t *testing.T) {
	conns := mockConnections()
	opts := seOptions{Args: []string{"devbox", "worker-1", "df -h"}}
	resConns, args := resolveExecTargetAndArgs(conns, opts)
	if len(resConns) != 2 || len(args) != 1 || args[0] != "df -h" {
		t.Fatalf("unexpected result: conns=%+v args=%+v", resConns, args)
	}
}

func TestResolveExecTargetAndArgs_NoTargetGitmapCmd(t *testing.T) {
	conns := mockConnections()
	opts := seOptions{Args: []string{"status", "--json"}}
	resConns, args := resolveExecTargetAndArgs(conns, opts)
	if len(resConns) != 2 || len(args) != 2 || args[0] != "status" {
		t.Fatalf("unexpected result: conns=%+v args=%+v", resConns, args)
	}
}

func TestResolveExecTargetAndArgs_IPTarget(t *testing.T) {
	conns := mockConnections()
	opts := seOptions{Args: []string{"192.168.1.20", "uname -a"}}
	resConns, args := resolveExecTargetAndArgs(conns, opts)
	if len(resConns) != 1 || resConns[0].IPAddress != "192.168.1.20" || len(args) != 1 {
		t.Fatalf("unexpected result: conns=%+v args=%+v", resConns, args)
	}
}
