package cluster

import (
	"context"
	"os/exec"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func TestExecLifecycle_ServerNodeGuard(t *testing.T) {
	node := ClusterNode{NodeRole: "server", IsServer: true}
	params := LifecycleExecParams{
		Context:          context.Background(),
		Node:             node,
		IsForceLifecycle: true,
		ProvidedPassword: "",
	}
	_, _, _, err := ExecRestart(params)
	if err == nil || !strings.Contains(err.Error(), constants.ErrClusterServerProtected) {
		t.Errorf("expected server protected error, got %v", err)
	}
}

func TestExecLifecycle_ForceLifecycleGuard(t *testing.T) {
	node := ClusterNode{NodeRole: "node", IsServer: false}
	params := LifecycleExecParams{
		Context:          context.Background(),
		Node:             node,
		IsForceLifecycle: false,
		ProvidedPassword: "",
	}
	_, _, _, err := ExecRestart(params)
	if err == nil || !strings.Contains(err.Error(), constants.ErrClusterLifecycleRequiresForce) {
		t.Errorf("expected force lifecycle error, got %v", err)
	}
}

func TestExecLifecycle_Success(t *testing.T) {
	origRunCmd := runCmdFunc
	defer func() { runCmdFunc = origRunCmd }()
	runCmdFunc = func(cmd *exec.Cmd) error { return nil }

	node := ClusterNode{NodeRole: "node", IsServer: false}
	params := LifecycleExecParams{
		Context:          context.Background(),
		Node:             node,
		IsForceLifecycle: true,
		ProvidedPassword: "",
	}
	_, _, code, err := ExecRestart(params)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
}
