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
//nolint:revive
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

	// Stub routing logic
	switch subCmd.Kind {
	case db.CommandKindPsCommand:
		res.ResultStatus = db.ResultStatusSucceeded
		if !node.IsServer {
			// Dial the remote agent
			conf := &tls.Config{InsecureSkipVerify: true}
			dialer := &net.Dialer{Timeout: dialAgentTimeout}
			client, err := tls.DialWithDialer(dialer, "tcp", node.IP+":8081", conf)
			if err != nil {
				res.ResultStatus = db.ResultStatusFailed
				msg := err.Error()
				res.ErrorMessage = &msg
			} else {
				defer client.Close()
				rpcClient := rpc.NewClient(client)
				args := &AgentExecArgs{Command: subCmd.RawArg}
				var reply AgentExecReply
				if err := rpcClient.Call("Agent.ExecPS", args, &reply); err != nil {
					res.ResultStatus = db.ResultStatusFailed
					msg := err.Error()
					res.ErrorMessage = &msg
				} else {
					if reply.ExitCode != 0 {
						res.ResultStatus = db.ResultStatusFailed
					}
					res.Stdout = &reply.Stdout
					res.Stderr = &reply.Stderr
					res.ExitCode = &reply.ExitCode
				}
			}
		} else {
			// Run locally
			stdout, stderr, exitCode, err := ExecPS(ctx, node, subCmd.RawArg)
			res.Stdout = &stdout
			res.Stderr = &stderr
			res.ExitCode = &exitCode
			if err != nil || exitCode != 0 {
				res.ResultStatus = db.ResultStatusFailed
				if err != nil {
					msg := err.Error()
					res.ErrorMessage = &msg
				}
			}
		}

	case db.CommandKindCmdCommand:
		res.ResultStatus = db.ResultStatusSucceeded
		if !node.IsServer {
			conf := &tls.Config{InsecureSkipVerify: true}
			dialer := &net.Dialer{Timeout: dialAgentTimeout}
			client, err := tls.DialWithDialer(dialer, "tcp", node.IP+":8081", conf)
			if err != nil {
				res.ResultStatus = db.ResultStatusFailed
				msg := err.Error()
				res.ErrorMessage = &msg
			} else {
				defer client.Close()
				rpcClient := rpc.NewClient(client)
				args := &AgentExecArgs{Command: subCmd.RawArg}
				var reply AgentExecReply
				if err := rpcClient.Call("Agent.ExecCmd", args, &reply); err != nil {
					res.ResultStatus = db.ResultStatusFailed
					msg := err.Error()
					res.ErrorMessage = &msg
				} else {
					if reply.ExitCode != 0 {
						res.ResultStatus = db.ResultStatusFailed
					}
					res.Stdout = &reply.Stdout
					res.Stderr = &reply.Stderr
					res.ExitCode = &reply.ExitCode
				}
			}
		} else {
			// Run locally
			stdout, stderr, exitCode, err := ExecCmd(ctx, node, subCmd.RawArg)
			res.Stdout = &stdout
			res.Stderr = &stderr
			res.ExitCode = &exitCode
			if err != nil || exitCode != 0 {
				res.ResultStatus = db.ResultStatusFailed
				if err != nil {
					msg := err.Error()
					res.ErrorMessage = &msg
				}
			}
		}
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
