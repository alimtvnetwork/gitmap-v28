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

	params := buildTestNodeParams()
	_, _, code, err := ExecRestart(params)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
}

func TestExecShutdown_Success(t *testing.T) {
	origRunCmd := runCmdFunc
	defer func() { runCmdFunc = origRunCmd }()

	var capturedCmd *exec.Cmd
	runCmdFunc = func(cmd *exec.Cmd) error {
		capturedCmd = cmd
		return nil
	}

	params := buildTestNodeParams()
	_, _, code, err := ExecShutdown(params)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if code != 0 || capturedCmd == nil {
		t.Errorf("expected exit 0 and captured cmd, got code %d cmd %v", code, capturedCmd)
	}
}

func TestExecLogoff_Success(t *testing.T) {
	origRunCmd := runCmdFunc
	defer func() { runCmdFunc = origRunCmd }()

	var capturedCmd *exec.Cmd
	runCmdFunc = func(cmd *exec.Cmd) error {
		capturedCmd = cmd
		return nil
	}

	params := buildTestNodeParams()
	_, _, code, err := ExecLogoff(params)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if code != 0 || capturedCmd == nil {
		t.Errorf("expected exit 0 and captured cmd, got code %d cmd %v", code, capturedCmd)
	}
}

func buildTestNodeParams() LifecycleExecParams {
	return LifecycleExecParams{
		Context:          context.Background(),
		Node:             ClusterNode{NodeRole: "node", IsServer: false},
		IsForceLifecycle: true,
		ProvidedPassword: "",
	}
}
