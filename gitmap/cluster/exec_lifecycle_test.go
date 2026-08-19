package cluster

import (
	"context"
	"os/exec"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func TestExecLifecycle_ServerNodeGuard(t *testing.T) {
	node := ClusterNode{NodeRole: "server", IsServer: true}
	_, _, _, err := ExecRestart(context.Background(), node, true, "")
	if err == nil || err.Error() != constants.ErrClusterServerProtected {
		t.Errorf("expected server protected error, got %v", err)
	}
}

func TestExecLifecycle_ForceLifecycleGuard(t *testing.T) {
	node := ClusterNode{NodeRole: "node", IsServer: false}
	_, _, _, err := ExecRestart(context.Background(), node, false, "")
	if err == nil || err.Error() != constants.ErrClusterLifecycleRequiresForce {
		t.Errorf("expected force lifecycle error, got %v", err)
	}
}

func TestExecLifecycle_Success(t *testing.T) {
	origRunCmd := runCmdFunc
	defer func() { runCmdFunc = origRunCmd }()
	runCmdFunc = func(cmd *exec.Cmd) error { return nil }

	node := ClusterNode{NodeRole: "node", IsServer: false}
	_, _, code, err := ExecRestart(context.Background(), node, true, "")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
}
