package cluster

import (
	"context"
	"crypto/tls"
	"net"
	"net/rpc"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/db"
)

type AgentExecArgs struct {
	Command string
}

type AgentExecReply struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

const dialAgentTimeout = 500 * time.Millisecond

// Dispatch routes the subcommand to the correct executor for the node.
// This is currently a stub for executors.

func Dispatch(ctx context.Context, node ClusterNode, subCmd ClusterSubCommand) db.ClusterExecResult {
	start := time.Now()
	res := db.ClusterExecResult{
		NodeId:       node.ID,
		SubCommand:   subCmd.Kind.String(),
		ResultStatus: db.ResultStatusPending,
		StartedAt:    &start,
	}

	raw := subCmd.RawArg
	res.CommandText = &raw

	switch subCmd.Kind {
	case db.CommandKindPsCommand:
		dispatchPSCommand(ctx, node, subCmd.RawArg, &res)
	case db.CommandKindCmdCommand:
		dispatchCmdCommand(ctx, node, subCmd.RawArg, &res)
	case db.CommandKindInstall:
		res.ResultStatus = db.ResultStatusSkipped
		msg := "ExecInstall stubbed"
		res.ErrorMessage = &msg
	case db.CommandKindGitPull, db.CommandKindGitPush, db.CommandKindGitCommit, db.CommandKindGitStatus:
		res.ResultStatus = db.ResultStatusSkipped
		msg := "ExecGit stubbed"
		res.ErrorMessage = &msg
	case db.CommandKindProjRun, db.CommandKindProjCreateCICD:
		res.ResultStatus = db.ResultStatusSkipped
		msg := "ExecProj stubbed"
		res.ErrorMessage = &msg
	case db.CommandKindRestart, db.CommandKindShutdown, db.CommandKindLogoff:
		res.ResultStatus = db.ResultStatusSkipped
		msg := "ExecLifecycle stubbed"
		res.ErrorMessage = &msg
	default:
		res.ResultStatus = db.ResultStatusFailed
		msg := "unknown command kind"
		res.ErrorMessage = &msg
	}

	end := time.Now()
	res.FinishedAt = &end
	durMs := int(end.Sub(start).Milliseconds())
	res.DurationMs = &durMs

	return res
}

func dispatchPSCommand(ctx context.Context, node ClusterNode, rawArg string, res *db.ClusterExecResult) {
	res.ResultStatus = db.ResultStatusSucceeded
	if !node.IsServer {
		dispatchRemotePS(node, rawArg, res)
		return
	}
	dispatchLocalPS(ctx, node, rawArg, res)
}

func dispatchRemotePS(node ClusterNode, rawArg string, res *db.ClusterExecResult) {
	conf := &tls.Config{InsecureSkipVerify: true}
	dialer := &net.Dialer{Timeout: dialAgentTimeout}
	client, err := tls.DialWithDialer(dialer, "tcp", node.IP+":8081", conf)
	if err != nil {
		setCommandError(res, err.Error())
		return
	}
	defer client.Close()
	rpcClient := rpc.NewClient(client)
	args := &AgentExecArgs{Command: rawArg}
	var reply AgentExecReply
	if err := rpcClient.Call("Agent.ExecPS", args, &reply); err != nil {
		setCommandError(res, err.Error())
		return
	}
	applyReply(res, reply)
}

func dispatchLocalPS(ctx context.Context, node ClusterNode, rawArg string, res *db.ClusterExecResult) {
	stdout, stderr, exitCode, err := ExecPS(ctx, node, rawArg)
	res.Stdout = &stdout
	res.Stderr = &stderr
	res.ExitCode = &exitCode
	if err != nil {
		setCommandError(res, err.Error())
		return
	}
	if exitCode != 0 {
		res.ResultStatus = db.ResultStatusFailed
	}
}

func dispatchCmdCommand(ctx context.Context, node ClusterNode, rawArg string, res *db.ClusterExecResult) {
	res.ResultStatus = db.ResultStatusSucceeded
	if !node.IsServer {
		dispatchRemoteCmd(node, rawArg, res)
		return
	}
	dispatchLocalCmd(ctx, node, rawArg, res)
}

func dispatchRemoteCmd(node ClusterNode, rawArg string, res *db.ClusterExecResult) {
	conf := &tls.Config{InsecureSkipVerify: true}
	dialer := &net.Dialer{Timeout: dialAgentTimeout}
	client, err := tls.DialWithDialer(dialer, "tcp", node.IP+":8081", conf)
	if err != nil {
		setCommandError(res, err.Error())
		return
	}
	defer client.Close()
	rpcClient := rpc.NewClient(client)
	args := &AgentExecArgs{Command: rawArg}
	var reply AgentExecReply
	if err := rpcClient.Call("Agent.ExecCmd", args, &reply); err != nil {
		setCommandError(res, err.Error())
		return
	}
	applyReply(res, reply)
}

func dispatchLocalCmd(ctx context.Context, node ClusterNode, rawArg string, res *db.ClusterExecResult) {
	stdout, stderr, exitCode, err := ExecCmd(ctx, node, rawArg)
	res.Stdout = &stdout
	res.Stderr = &stderr
	res.ExitCode = &exitCode
	if err != nil {
		setCommandError(res, err.Error())
		return
	}
	if exitCode != 0 {
		res.ResultStatus = db.ResultStatusFailed
	}
}

func setCommandError(res *db.ClusterExecResult, msg string) {
	res.ResultStatus = db.ResultStatusFailed
	res.ErrorMessage = &msg
}

func applyReply(res *db.ClusterExecResult, reply AgentExecReply) {
	if reply.ExitCode != 0 {
		res.ResultStatus = db.ResultStatusFailed
	}
	res.Stdout = &reply.Stdout
	res.Stderr = &reply.Stderr
	res.ExitCode = &reply.ExitCode
}
