package cluster

import (
	"context"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// PeerDelegationPayload encapsulates parameters for delegating execution to a peer node.
type PeerDelegationPayload struct {
	TaskID     string            `json:"taskId,omitempty"`
	Action     string            `json:"action"`
	Command    string            `json:"command"`
	Args       []string          `json:"args,omitempty"`
	WorkDir    string            `json:"workDir,omitempty"`
	TargetRepo string            `json:"targetRepo,omitempty"`
	Requester  string            `json:"requester,omitempty"`
	Env        map[string]string `json:"env,omitempty"`
	TimeoutSec int               `json:"timeoutSec,omitempty"`
	IsAsync    bool              `json:"isAsync,omitempty"`
}

// PeerDelegationResponse represents telemetry returned from peer node execution.
type PeerDelegationResponse struct {
	TaskID       string `json:"taskId,omitempty"`
	NodeID       string `json:"nodeId,omitempty"`
	NodeAlias    string `json:"nodeAlias,omitempty"`
	IsSuccess    bool   `json:"isSuccess"`
	ExitCode     int    `json:"exitCode"`
	Stdout       string `json:"stdout,omitempty"`
	Stderr       string `json:"stderr,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
	DurationMs   int64  `json:"durationMs"`
}

type peerDelegationRunnerFunc func(ctx context.Context, node ClusterNode, cmd string) (string, string, int, error)

var peerDelegationRunner peerDelegationRunnerFunc = defaultPeerDelegationRunner

func defaultPeerDelegationRunner(ctx context.Context, node ClusterNode, cmd string) (string, string, int, error) {
	isWin := node.OS == "windows"
	if isWin {
		return ExecPS(ctx, node, cmd)
	}
	return ExecCmd(ctx, node, cmd)
}

// DispatchPeerDelegation dispatches a task execution to a cluster node and collects results.
func DispatchPeerDelegation(
	ctx context.Context,
	node ClusterNode,
	payload PeerDelegationPayload,
) (*PeerDelegationResponse, error) {
	cmdStr := resolvePeerDelegationCommand(payload)
	start := time.Now()

	execCtx, cancel := resolvePeerContext(ctx, payload.TimeoutSec)
	defer cancel()

	stdout, stderr, exitCode, err := peerDelegationRunner(execCtx, node, cmdStr)
	duration := time.Since(start).Milliseconds()

	return buildPeerDelegationResponse(payload.TaskID, node, stdout, stderr, exitCode, duration, err)
}

func resolvePeerContext(ctx context.Context, timeoutSec int) (context.Context, context.CancelFunc) {
	hasTimeout := timeoutSec > 0
	if hasTimeout {
		return context.WithTimeout(ctx, time.Duration(timeoutSec)*time.Second)
	}
	return context.WithCancel(ctx)
}

func resolvePeerDelegationCommand(payload PeerDelegationPayload) string {
	hasCommand := len(payload.Command) > 0
	if hasCommand {
		return payload.Command
	}
	hasAction := len(payload.Action) > 0
	if hasAction {
		return "gitmap " + payload.Action
	}
	return "gitmap status"
}

func buildPeerDelegationResponse(
	taskID string,
	node ClusterNode,
	stdout string,
	stderr string,
	exitCode int,
	durationMs int64,
	err error,
) (*PeerDelegationResponse, error) {
	isSuccess := exitCode == 0 && err == nil
	resp := &PeerDelegationResponse{
		TaskID:     taskID,
		NodeID:     node.ID,
		NodeAlias:  node.Alias,
		ExitCode:   exitCode,
		Stdout:     stdout,
		Stderr:     stderr,
		DurationMs: durationMs,
		IsSuccess:  isSuccess,
	}

	if err != nil {
		resp.ErrorMessage = err.Error()
		return resp, apperror.WrapSimple(err, "peer delegation failed:")
	}

	return resp, nil
}
