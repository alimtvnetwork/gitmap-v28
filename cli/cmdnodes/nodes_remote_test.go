package cmdnodes

import (
	"context"
	"errors"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/cluster"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

type mockRemoteSSHExecutor struct {
	output string
	err    error
}

func (m *mockRemoteSSHExecutor) Execute(conn db.SSHConnection, command string, shell string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.output, nil
}

func TestParseNodesRemoteOptions_WithTargetAndCommand(t *testing.T) {
	args := []string{"worker-1", "gitmap status", "--json"}
	opts := parseNodesRemoteOptions(args)

	if opts.FilterOpts.Target != "worker-1" {
		t.Fatalf("expected target worker-1, got %s", opts.FilterOpts.Target)
	}
	if opts.Command != "gitmap status" {
		t.Fatalf("expected command 'gitmap status', got %s", opts.Command)
	}
	if !opts.IsJSON {
		t.Fatalf("expected isJSON true")
	}
}

func TestParseNodesRemoteOptions_FlagTarget(t *testing.T) {
	args := []string{"--target=node-a", "uptime", "--except=node-b"}
	opts := parseNodesRemoteOptions(args)

	if opts.FilterOpts.Target != "node-a" {
		t.Fatalf("expected target node-a, got %s", opts.FilterOpts.Target)
	}
	if opts.Command != "uptime" {
		t.Fatalf("expected command uptime, got %s", opts.Command)
	}
	if len(opts.FilterOpts.Except) != 1 || opts.FilterOpts.Except[0] != "node-b" {
		t.Fatalf("expected except node-b, got %v", opts.FilterOpts.Except)
	}
}

func TestResolveRemoteExecShell(t *testing.T) {
	if resolveRemoteExecShell("windows") != "ps" {
		t.Fatalf("expected ps for windows")
	}
	if resolveRemoteExecShell("linux") != "bash" {
		t.Fatalf("expected bash for linux")
	}
	if resolveRemoteExecShell("darwin") != "bash" {
		t.Fatalf("expected bash for darwin")
	}
}

func TestExecuteSingleRemoteNode_Success(t *testing.T) {
	oldExec := currentSSHExecutor
	defer func() { currentSSHExecutor = oldExec }()

	currentSSHExecutor = &mockRemoteSSHExecutor{output: "v1.0.0"}
	conn := db.SSHConnection{
		Alias:     "node-1",
		IPAddress: "192.168.1.100",
		OS:        "linux",
	}

	res := executeSingleRemoteNode(conn, "gitmap version")
	if !res.IsSuccess {
		t.Fatalf("expected success")
	}
	if res.Stdout != "v1.0.0" {
		t.Fatalf("expected v1.0.0, got %s", res.Stdout)
	}
	if res.ExitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", res.ExitCode)
	}
}

func TestExecuteSingleRemoteNode_Error(t *testing.T) {
	oldExec := currentSSHExecutor
	defer func() { currentSSHExecutor = oldExec }()

	currentSSHExecutor = &mockRemoteSSHExecutor{err: errors.New("connection timed out")}
	conn := db.SSHConnection{
		Alias:     "node-2",
		IPAddress: "192.168.1.101",
		OS:        "windows",
	}

	res := executeSingleRemoteNode(conn, "gitmap status")
	if res.IsSuccess {
		t.Fatalf("expected failure")
	}
	if res.ExitCode != 1 {
		t.Fatalf("expected exit code 1, got %d", res.ExitCode)
	}
	if res.ErrorMessage != "connection timed out" {
		t.Fatalf("expected error message 'connection timed out', got %s", res.ErrorMessage)
	}
}

func TestDispatchNodesRemoteConcurrent(t *testing.T) {
	oldExec := currentSSHExecutor
	defer func() { currentSSHExecutor = oldExec }()

	currentSSHExecutor = &mockRemoteSSHExecutor{output: "ok"}
	targets := []db.SSHConnection{
		{Alias: "n1", IPAddress: "10.0.0.1", OS: "linux"},
		{Alias: "n2", IPAddress: "10.0.0.2", OS: "linux"},
	}

	results := dispatchNodesRemoteConcurrent(targets, "echo hi")
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if !r.IsSuccess {
			t.Fatalf("expected node %s to succeed", r.NodeAlias)
		}
	}
}

func TestDispatchPeerDelegation(t *testing.T) {
	node := cluster.ClusterNode{
		ID:    "node-a",
		Alias: "worker-a",
		OS:    "linux",
	}
	payload := cluster.PeerDelegationPayload{
		TaskID:  "task-123",
		Action:  "status",
		Command: "gitmap status",
	}

	resp, _ := cluster.DispatchPeerDelegation(context.Background(), node, payload)
	if resp == nil {
		t.Fatalf("expected non-nil response")
	}
	if resp.TaskID != "task-123" {
		t.Fatalf("expected task-123, got %s", resp.TaskID)
	}
	if resp.NodeAlias != "worker-a" {
		t.Fatalf("expected worker-a, got %s", resp.NodeAlias)
	}
}

